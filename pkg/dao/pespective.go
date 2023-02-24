package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/perspective"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

func PerspectiveCreates(mc model.ClientSet, input ...*model.Perspective) ([]*model.PerspectiveCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*model.PerspectiveCreate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// required.
		var c = mc.Perspectives().Create().
			SetName(r.Name).
			SetStartTime(r.StartTime).
			SetEndTime(r.EndTime).
			SetBuiltin(r.Builtin)

		// optional.
		if len(r.AllocationQueries) != 0 {
			c.SetAllocationQueries(r.AllocationQueries)
		}

		rrs[i] = c
	}
	return rrs, nil
}

func PerspectiveUpdates(mc model.ClientSet, input ...*model.Perspective) ([]*model.PerspectiveUpdate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*model.PerspectiveUpdate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}
		// predicated.
		var ps = []predicate.Perspective{
			perspective.ID(r.ID),
		}
		var c = mc.Perspectives().Update().
			Where(ps...)
		if r.StartTime != "" {
			c.SetStartTime(r.StartTime)
		}
		if r.EndTime != "" {
			c.SetStartTime(r.EndTime)
		}
		if len(r.AllocationQueries) != 0 {
			c.SetAllocationQueries(r.AllocationQueries)
		}
		rrs[i] = c
	}
	return rrs, nil
}
