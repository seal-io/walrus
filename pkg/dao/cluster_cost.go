package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
)

func ClusterCostCreates(mc model.ClientSet, input ...*model.ClusterCost) ([]*model.ClusterCostCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.ClusterCostCreate, len(input))

	for i := range input {
		r := input[i]
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		c := mc.ClusterCosts().Create().
			// Required.
			SetStartTime(r.StartTime).
			SetEndTime(r.EndTime).
			SetMinutes(r.Minutes).
			SetClusterName(r.ClusterName).
			SetConnectorID(r.ConnectorID).
			SetTotalCost(r.TotalCost).
			SetCurrency(r.Currency).
			// Optional.
			SetManagementCost(r.ManagementCost).
			SetIdleCost(r.IdleCost).
			SetAllocationCost(r.AllocationCost)

		rrs[i] = c
	}

	return rrs, nil
}
