package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/strs"
)

func AllocationCostCreates(mc model.ClientSet, input ...*model.AllocationCost) ([]*model.AllocationCostCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*model.AllocationCostCreate, len(input))
	for i := range input {
		r := input[i]
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		c := mc.AllocationCosts().Create().
			// required.
			SetStartTime(r.StartTime).
			SetEndTime(r.EndTime).
			SetMinutes(r.Minutes).
			SetClusterName(r.ClusterName).
			SetConnectorID(r.ConnectorID).
			// optional.
			SetName(r.Name).
			SetClusterName(r.ClusterName).
			SetNamespace(r.Namespace).
			SetNode(r.Node).
			SetControllerKind(r.ControllerKind).
			SetController(r.Controller).
			SetPod(r.Pod).
			SetContainer(r.Container).
			SetPvs(r.Pvs).
			SetLabels(r.Labels).
			SetFingerprint(allocationCostFingerprint(r)).
			SetTotalCost(r.TotalCost).
			SetCurrency(r.Currency).
			SetCpuCost(r.CpuCost).
			SetCpuCoreRequest(r.CpuCoreRequest).
			SetGpuCost(r.GpuCost).
			SetGpuCount(r.GpuCount).
			SetRamCost(r.RamCost).
			SetRamByteRequest(r.RamByteRequest).
			SetPvCost(r.PvCost).
			SetPvBytes(r.PvBytes).
			SetCpuCoreUsageAverage(r.CpuCoreUsageAverage).
			SetCpuCoreUsageMax(r.CpuCoreUsageMax).
			SetRamByteUsageAverage(r.RamByteUsageAverage).
			SetRamByteUsageMax(r.RamByteUsageMax).
			SetConnectorID(r.ConnectorID).
			SetTotalCost(r.TotalCost).
			SetCurrency(r.Currency).
			SetCpuCost(r.CpuCost).
			SetGpuCost(r.GpuCost).
			SetRamCost(r.RamCost).
			SetLoadBalancerCost(r.LoadBalancerCost)

		rrs[i] = c
	}
	return rrs, nil
}

func allocationCostFingerprint(cost *model.AllocationCost) string {
	// TODO: support pv
	return strs.Join("/", cost.ClusterName, cost.Node, cost.Namespace, cost.Pod, cost.Container)
}
