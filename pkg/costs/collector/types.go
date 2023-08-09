package collector

import (
	"time"

	"github.com/opencost/opencost/pkg/costmodel"
	"github.com/opencost/opencost/pkg/kubecost"
)

type (
	AllocationComputeResponse struct {
		Code    int                     `json:"code"`
		Status  string                  `json:"status"`
		Data    []map[string]Allocation `json:"data"`
		Message string                  `json:"message,omitempty"`
		Warning string                  `json:"warning,omitempty"`
	}
	Allocation struct {
		Name                       string                          `json:"name"`
		Properties                 *kubecost.AllocationProperties  `json:"properties,omitempty"`
		Window                     Window                          `json:"window"`
		Start                      time.Time                       `json:"start"`
		End                        time.Time                       `json:"end"`
		CPUCoreHours               float64                         `json:"cpuCoreHours"`
		CPUCoreRequestAverage      float64                         `json:"cpuCoreRequestAverage"`
		CPUCoreUsageAverage        float64                         `json:"cpuCoreUsageAverage"`
		CPUCost                    float64                         `json:"cpuCost"`
		CPUCostAdjustment          float64                         `json:"cpuCostAdjustment"`
		GPUHours                   float64                         `json:"gpuHours"`
		GPUCost                    float64                         `json:"gpuCost"`
		GPUCostAdjustment          float64                         `json:"gpuCostAdjustment"`
		NetworkTransferBytes       float64                         `json:"networkTransferBytes"`
		NetworkReceiveBytes        float64                         `json:"networkReceiveBytes"`
		NetworkCost                float64                         `json:"networkCost"`
		NetworkCrossZoneCost       float64                         `json:"networkCrossZoneCost"`
		NetworkCrossRegionCost     float64                         `json:"networkCrossRegionCost"`
		NetworkInternetCost        float64                         `json:"networkInternetCost"`
		NetworkCostAdjustment      float64                         `json:"networkCostAdjustment"`
		LoadBalancerCost           float64                         `json:"loadBalancerCost"`
		LoadBalancerCostAdjustment float64                         `json:"loadBalancerCostAdjustment"`
		PVs                        kubecost.PVAllocations          `json:"pvs"`
		PVCostAdjustment           float64                         `json:"pvCostAdjustment"`
		RAMByteHours               float64                         `json:"ramByteHours"`
		RAMBytesRequestAverage     float64                         `json:"ramByteRequestAverage"`
		RAMBytesUsageAverage       float64                         `json:"ramByteUsageAverage"`
		RAMCost                    float64                         `json:"ramCost"`
		RAMCostAdjustment          float64                         `json:"ramCostAdjustment"`
		SharedCost                 float64                         `json:"sharedCost"`
		ExternalCost               float64                         `json:"externalCost"`
		RawAllocationOnly          *kubecost.RawAllocationOnlyData `json:"rawAllocationOnly"`
	}
	Window struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}
)

type ClusterCostResponse struct {
	Code    int                                `json:"code"`
	Status  string                             `json:"status"`
	Data    map[string]*costmodel.ClusterCosts `json:"data"`
	Message string                             `json:"message,omitempty"`
	Warning string                             `json:"warning,omitempty"`
}

type (
	PrometheusQueryRangeResult struct {
		Status string                         `json:"status"`
		Data   PrometheusQueryRangeResultData `json:"data"`
	}
	PrometheusQueryRangeResultData struct {
		ResultType string        `json:"resultType"`
		Result     []QueryResult `json:"result"`
	}
	QueryResult struct {
		Metric map[string]any `json:"metric"`
		Values [][]any        `json:"values"`
	}
	Vector struct {
		Timestamp float64 `json:"timestamp"`
		Value     any     `json:"value"`
	}
)

func (a *Allocation) kubecostAllocation() kubecost.Allocation {
	return kubecost.Allocation{
		Name:                       a.Name,
		Properties:                 a.Properties,
		Window:                     kubecost.NewWindow(&a.Window.Start, &a.Window.End),
		Start:                      a.Start,
		End:                        a.End,
		CPUCoreHours:               a.CPUCoreHours,
		CPUCoreRequestAverage:      a.CPUCoreRequestAverage,
		CPUCoreUsageAverage:        a.CPUCoreUsageAverage,
		CPUCost:                    a.CPUCost,
		CPUCostAdjustment:          a.CPUCostAdjustment,
		GPUHours:                   a.GPUHours,
		GPUCost:                    a.GPUCost,
		GPUCostAdjustment:          a.GPUCostAdjustment,
		NetworkTransferBytes:       a.NetworkTransferBytes,
		NetworkReceiveBytes:        a.NetworkReceiveBytes,
		NetworkCost:                a.NetworkCost,
		NetworkCrossZoneCost:       a.NetworkCrossZoneCost,
		NetworkCrossRegionCost:     a.NetworkCrossRegionCost,
		NetworkInternetCost:        a.NetworkInternetCost,
		NetworkCostAdjustment:      a.NetworkCostAdjustment,
		LoadBalancerCost:           a.LoadBalancerCost,
		LoadBalancerCostAdjustment: a.LoadBalancerCostAdjustment,
		PVs:                        a.PVs,
		PVCostAdjustment:           a.PVCostAdjustment,
		RAMByteHours:               a.RAMByteHours,
		RAMBytesRequestAverage:     a.RAMBytesRequestAverage,
		RAMBytesUsageAverage:       a.RAMBytesUsageAverage,
		RAMCost:                    a.RAMCost,
		RAMCostAdjustment:          a.RAMCostAdjustment,
		SharedCost:                 a.SharedCost,
		ExternalCost:               a.ExternalCost,
		RawAllocationOnly:          a.RawAllocationOnly,
	}
}
