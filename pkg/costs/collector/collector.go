package collector

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/opencost/opencost/pkg/costmodel"
	"github.com/opencost/opencost/pkg/kubecost"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/costs/deployer"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/json"
)

// cost endpoint access path.
var (
	pathServiceProxy         = fmt.Sprintf("/api/v1/namespaces/%s/services/http:%s:9003/proxy", deployer.NamespaceSeal, deployer.NameOpencost)
	pathAllocation           = "/allocation/compute"
	pathClusterCost          = "/clusterCosts"
	pathPrometheusQueryRange = "/prometheusQueryRange"
)

// prometheus expression.
const (
	// exprClusterMgmtHrCost defined expression for management cost.
	exprClusterMgmtHrCost = "avg(avg_over_time(kubecost_cluster_management_cost[1h:5m]))"
)

type Collector struct {
	clusterName   string
	clusterClient *http.Client
	restCfg       *rest.Config
	conn          *model.Connector
}

func NewCollector(restCfg *rest.Config, conn *model.Connector, clusterName string) (*Collector, error) {
	client, err := rest.HTTPClientFor(restCfg)
	if err != nil {
		return nil, err
	}
	return &Collector{
		clusterName:   clusterName,
		clusterClient: client,
		restCfg:       restCfg,
		conn:          conn,
	}, nil
}

func (c *Collector) K8sCosts(startTime, endTime *time.Time, step time.Duration) (
	[]*model.ClusterCost, []*model.AllocationCost, error) {
	cc, err := c.clusterCosts(startTime, endTime, step)
	if err != nil {
		return nil, nil, err
	}

	if len(cc) == 0 {
		return nil, nil, nil
	}

	ac, err := c.allocationResourceCosts(startTime, endTime, step)
	if err != nil {
		return nil, nil, err
	}

	c.applyIdleCost(cc, ac)
	return cc, ac, nil
}

// allocationResourceCosts get cost for allocation resources.
func (c *Collector) allocationResourceCosts(startTime, endTime *time.Time, step time.Duration) ([]*model.AllocationCost, error) {
	window := fmt.Sprintf("%d,%d", startTime.Unix(), endTime.Unix())
	queries := url.Values{
		// each AllocationSet would be a container.
		"aggregate": []string{"container"},
		// accumulate sums each AllocationSet in the given range, just returning a single cumulative.
		"accumulate": []string{"false"},
		// e.g. "1586822400,1586908800"
		"window": []string{window},
		// e.g. "1h"
		"step": []string{step.String()},
	}

	u, err := url.Parse(c.restCfg.Host)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, pathServiceProxy, pathAllocation)
	u.RawQuery = queries.Encode()

	ac := &AllocationComputeResponse{}
	if err = c.getRequest(u.String(), ac); err != nil {
		return nil, err
	}

	if len(ac.Data) == 0 {
		return nil, nil
	}

	var costs []*model.AllocationCost
	for _, data := range ac.Data {
		for _, v := range data {
			ka := v.kubecostAllocation()
			cost := &model.AllocationCost{
				ConnectorID:         c.conn.ID,
				StartTime:           v.Window.Start,
				EndTime:             v.Window.End,
				Minutes:             ka.Minutes(),
				Name:                v.Name,
				ClusterName:         c.clusterName,
				Namespace:           v.Properties.Namespace,
				Node:                v.Properties.Node,
				Controller:          v.Properties.Controller,
				ControllerKind:      v.Properties.ControllerKind,
				Pod:                 v.Properties.Pod,
				Container:           v.Properties.Container,
				Pvs:                 toPVs(v.PVs),
				Labels:              v.Properties.Labels,
				TotalCost:           ka.TotalCost(),
				CpuCost:             ka.CPUTotalCost(),
				CpuCoreRequest:      ka.CPUCoreRequestAverage,
				RamCost:             ka.RAMTotalCost(),
				RamByteRequest:      v.RAMBytesRequestAverage,
				PvCost:              ka.PVTotalCost(),
				PvBytes:             ka.PVBytes(),
				CpuCoreUsageAverage: v.CPUCoreUsageAverage,
				RamByteUsageAverage: v.RAMBytesUsageAverage,
			}

			if v.RawAllocationOnly != nil {
				cost.CpuCoreUsageMax = v.RawAllocationOnly.CPUCoreUsageMax
				cost.RamByteUsageMax = v.RawAllocationOnly.RAMBytesUsageMax
			}

			costs = append(costs, cost)
		}
	}
	return costs, nil
}

// clusterCosts get costs for cluster.
func (c *Collector) clusterCosts(startTime, endTime *time.Time, step time.Duration) ([]*model.ClusterCost, error) {
	var costs []*model.ClusterCost
	var stepStart = *startTime
	for endTime.After(stepStart) {
		stepEnd := stepStart.Add(step)
		cc, err := c.clusterCostsWithinRange(&stepStart, &stepEnd)
		if err != nil {
			return nil, err
		}

		mgmtCost, err := c.clusterManagementCost(&stepStart, &stepEnd)
		if err != nil {
			return nil, err
		}

		stepStart = stepStart.Add(step)
		switch {
		case cc == nil:
			continue
		default:
			cc.ManagementCost = mgmtCost
			cc.TotalCost += mgmtCost
			costs = append(costs, cc)
		}
	}
	return costs, nil
}

// getClusterCostWithinRange get cluster cost within range.
func (c *Collector) clusterCostsWithinRange(startTime, endTime *time.Time) (*model.ClusterCost, error) {
	offset := math.Ceil(time.Since(*endTime).Minutes())
	if offset < 0 {
		return nil, nil
	}

	window := math.Ceil(endTime.Sub(*startTime).Minutes())
	queries := url.Values{
		// e.g. 1h
		"window": []string{fmt.Sprintf("%.0fm", window)},
		// e.g. 1h
		"offset": []string{fmt.Sprintf("%.0fm", offset)},
	}

	u, err := url.Parse(c.restCfg.Host)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, pathServiceProxy, pathClusterCost)
	u.RawQuery = queries.Encode()

	cc := &ClusterCostResponse{}
	if err = c.getRequest(u.String(), cc); err != nil {
		return nil, err
	}

	var clusterCost *costmodel.ClusterCosts
	for _, v := range cc.Data {
		if v != nil {
			clusterCost = v
		}
		break
	}

	if clusterCost == nil {
		return nil, nil
	}

	return &model.ClusterCost{
		ConnectorID: c.conn.ID,
		StartTime:   *startTime,
		EndTime:     *endTime,
		Minutes:     window,
		ClusterName: c.clusterName,
		TotalCost:   clusterCost.TotalCumulative,
		CpuCost:     clusterCost.CPUCumulative,
		RamCost:     clusterCost.RAMCumulative,
		StorageCost: clusterCost.StorageCumulative,
	}, nil
}

func (c *Collector) applyIdleCost(ccs []*model.ClusterCost, acs []*model.AllocationCost) {
	var allocationCosts = make(map[string]float64)

	for _, v := range acs {
		key := fmt.Sprintf("%s-%s", v.StartTime.String(), v.EndTime.String())
		allocationCosts[key] += v.TotalCost
	}

	for i, v := range ccs {
		key := fmt.Sprintf("%s-%s", v.StartTime.String(), v.EndTime.String())
		if ac, ok := allocationCosts[key]; ok {
			ccs[i].AllocationCost = ac
			ccs[i].IdleCost = ccs[i].TotalCost - ccs[i].ManagementCost - ac
		}
	}
}

// clusterManagementCost get cluster management cost.
func (c *Collector) clusterManagementCost(startTime, endTime *time.Time) (float64, error) {
	layout := "2006-01-02T15:04:05.000Z"
	queries := url.Values{
		// e.g "2006-01-02T15:04:05.000Z"
		"start": []string{startTime.Format(layout)},
		// e.g "2006-01-02T15:04:05.000Z"
		"end": []string{endTime.Format(layout)},
		// prometheus query step
		"duration": []string{"1h"},
		// prometheus query expression
		"query": []string{exprClusterMgmtHrCost},
	}

	u, err := url.Parse(c.restCfg.Host)
	if err != nil {
		return 0, err
	}

	u.Path = path.Join(u.Path, pathServiceProxy, pathPrometheusQueryRange)
	u.RawQuery = queries.Encode()

	obj := &PrometheusQueryRangeResult{}
	if err = c.getRequest(u.String(), obj); err != nil {
		return 0, err
	}

	if len(obj.Data.Result) == 0 || len(obj.Data.Result[0].Values) == 0 || len(obj.Data.Result[0].Values[0]) < 2 {
		return 0, nil
	}

	value := obj.Data.Result[0].Values[0][1]
	mgntCost, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
	if err != nil {
		return 0, err
	}
	return mgntCost, nil
}

func (c *Collector) getRequest(url string, obj interface{}) error {
	resp, err := c.clusterClient.Get(url)
	if err != nil {
		return fmt.Errorf("request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response from %s: %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response from %s, code: %d, body: %s", url, resp.StatusCode, string(body))
	}

	if err = json.Unmarshal(body, obj); err != nil {
		return fmt.Errorf("decode response from %s: %w", url, err)
	}
	return nil
}

func toPVs(pvAlloc kubecost.PVAllocations) map[string]types.PVCost {
	if pvAlloc == nil {
		return nil
	}

	pvs := make(map[string]types.PVCost)
	for k, v := range pvAlloc {
		pvs[k.Name] = types.PVCost{
			Cost:  v.Cost,
			Bytes: v.ByteHours,
		}
	}
	return pvs
}
