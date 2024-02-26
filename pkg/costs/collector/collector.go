package collector

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/opencost/opencost/pkg/kubecost"
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/costs/deployer"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/utils/json"
)

// cost endpoint access path.
var (
	pathServiceProxy = fmt.Sprintf("/api/v1/namespaces/%s/services/http:%s:9003/proxy",
		types.WalrusSystemNamespace, deployer.NameOpencost)
	pathAllocation           = "/allocation/compute"
	pathPrometheusQueryRange = "/prometheusQueryRange"
)

// labelMapping indicate the relation between opencost converted label and original label.
var (
	labelMapping = map[string]string{
		"walrus_seal_io_project_name":     types.LabelWalrusProjectName,
		"walrus_seal_io_environment_name": types.LabelWalrusEnvironmentName,
		"walrus_seal_io_service_name":     types.LabelWalrusResourceName,
	}
)

// prometheus expression.
const (
	// ExprClusterMgmtHrCost defined expression for management cost.
	exprClusterMgmtHrCost = "avg(avg_over_time(kubecost_cluster_management_cost[1h:5m]))"
)

type Collector struct {
	clusterName   string
	clusterClient *http.Client
	restCfg       *rest.Config
	conn          *model.Connector
}

func NewCollector(
	restCfg *rest.Config,
	conn *model.Connector,
	clusterName string,
) (*Collector, error) {
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

func (c *Collector) K8sCosts(
	startTime, endTime *time.Time,
	step time.Duration,
) ([]*model.CostReport, error) {
	ac, err := c.allocationResourceCosts(startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("error get allocation resource costs: %w", err)
	}

	mgntCost, err := c.clusterManagementCost(startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error get cluster management cost: %w", err)
	}

	if mgntCost != nil && mgntCost.TotalCost > 0 {
		return append(ac, mgntCost), nil
	}

	return ac, nil
}

// allocationResourceCosts get cost for allocation resources.
func (c *Collector) allocationResourceCosts(
	startTime, endTime *time.Time,
	step time.Duration,
) ([]*model.CostReport, error) {
	window := fmt.Sprintf("%d,%d", startTime.Unix(), endTime.Unix())
	queries := url.Values{
		// Each AllocationSet would be a container, use pod,
		// container as aggregate key to prevent containers with same name.
		"aggregate": []string{"pod,container"},
		// Accumulate sums each AllocationSet in the given range, just returning a single cumulative.
		"accumulate": []string{"false"},
		// E.g. "1586822400,1586908800".
		"window": []string{window},
		// E.g. "1h".
		"step": []string{step.String()},
		// Include idle cost.
		"includeIdle": []string{"true"},
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

	var costs []*model.CostReport

	for _, data := range ac.Data {
		for _, v := range data {
			var (
				name = v.Name
				ka   = v.kubecostAllocation()
			)

			cost := &model.CostReport{
				ConnectorID:         c.conn.ID,
				StartTime:           v.Window.Start,
				EndTime:             v.Window.End,
				Minutes:             ka.Minutes(),
				Name:                name,
				ClusterName:         c.clusterName,
				Namespace:           v.Properties.Namespace,
				Node:                v.Properties.Node,
				Controller:          v.Properties.Controller,
				ControllerKind:      v.Properties.ControllerKind,
				Pod:                 v.Properties.Pod,
				Container:           v.Properties.Container,
				Pvs:                 toPVs(v.PVs),
				Labels:              toLabels(v.Properties.Labels),
				TotalCost:           ka.TotalCost(),
				CPUCost:             ka.CPUTotalCost(),
				CPUCoreRequest:      ka.CPUCoreRequestAverage,
				RAMCost:             ka.RAMTotalCost(),
				RAMByteRequest:      v.RAMBytesRequestAverage,
				PVCost:              ka.PVTotalCost(),
				PVBytes:             ka.PVBytes(),
				LoadBalancerCost:    v.LoadBalancerCost,
				CPUCoreUsageAverage: v.CPUCoreUsageAverage,
				RAMByteUsageAverage: v.RAMBytesUsageAverage,
			}

			if types.IsIdleCost(name) {
				cost.Namespace = name
				cost.Node = name
				cost.Controller = name
				cost.ControllerKind = name
				cost.Pod = name
				cost.Container = name
			}

			if v.RawAllocationOnly != nil {
				cost.CPUCoreUsageMax = v.RawAllocationOnly.CPUCoreUsageMax
				cost.RAMByteUsageMax = v.RawAllocationOnly.RAMBytesUsageMax
			}

			costs = append(costs, cost)
		}
	}

	return costs, nil
}

// clusterManagementCost get cluster management cost.
func (c *Collector) clusterManagementCost(startTime, endTime *time.Time) (*model.CostReport, error) {
	layout := "2006-01-02T15:04:05.000Z"
	queries := url.Values{
		// E.g "2006-01-02T15:04:05.000Z".
		"start": []string{startTime.Format(layout)},
		// E.g "2006-01-02T15:04:05.000Z".
		"end": []string{endTime.Format(layout)},
		// Prometheus query step.
		"duration": []string{"1h"},
		// Prometheus query expression.
		"query": []string{exprClusterMgmtHrCost},
	}

	u, err := url.Parse(c.restCfg.Host)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, pathServiceProxy, pathPrometheusQueryRange)
	u.RawQuery = queries.Encode()

	obj := &PrometheusQueryRangeResult{}
	if err = c.getRequest(u.String(), obj); err != nil {
		return nil, err
	}

	if len(obj.Data.Result) == 0 || len(obj.Data.Result[0].Values) == 0 ||
		len(obj.Data.Result[0].Values[0]) < 2 {
		return nil, nil
	}

	value := obj.Data.Result[0].Values[0][1]

	mgntCost, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
	if err != nil {
		return nil, err
	}

	name := types.ManagementCostItemName

	return &model.CostReport{
		ConnectorID:    c.conn.ID,
		StartTime:      *startTime,
		EndTime:        *endTime,
		Minutes:        endTime.Sub(*startTime).Minutes(),
		Name:           name,
		Namespace:      name,
		Node:           name,
		Controller:     name,
		ControllerKind: name,
		Pod:            name,
		Container:      name,
		ClusterName:    c.clusterName,
		TotalCost:      mgntCost,
	}, nil
}

func (c *Collector) getRequest(url string, obj any) error {
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
		return fmt.Errorf(
			"response from %s, code: %d, body: %s",
			url,
			resp.StatusCode,
			string(body),
		)
	}

	if len(body) == 0 {
		return nil
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

func toLabels(origin map[string]string) map[string]string {
	labels := origin

	for k, v := range origin {
		mapping, ok := labelMapping[k]
		if ok {
			labels[mapping] = v
		}
	}

	proj, ok1 := labels[types.LabelWalrusProjectName]
	env, ok2 := labels[types.LabelWalrusEnvironmentName]
	svc, ok3 := labels[types.LabelWalrusResourceName]

	if ok1 && ok2 && ok3 {
		labels[types.LabelWalrusEnvironmentPath] = fmt.Sprintf("%s/%s", proj, env)
		labels[types.LabelWalrusResourcePath] = fmt.Sprintf("%s/%s/%s", proj, env, svc)
	}

	return labels
}
