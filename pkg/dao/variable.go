package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/variable"
)

func VariableCreates(mc model.ClientSet, input ...*model.Variable) ([]*model.VariableCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.VariableCreate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.Variables().Create().
			SetName(r.Name).
			SetValue(r.Value).
			SetSensitive(r.Sensitive)

		// Optional.
		if r.ProjectID.IsNaive() {
			c.SetProjectID(r.ProjectID)
		}

		if r.EnvironmentID.IsNaive() {
			c.SetEnvironmentID(r.EnvironmentID)
		}

		if r.Description != "" {
			c.SetDescription(r.Description)
		}
		rrs[i] = c
	}

	return rrs, nil
}

func VariableUpdates(mc model.ClientSet, input ...*model.Variable) ([]*model.VariableUpdate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.VariableUpdate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Predicated.
		var ps []predicate.Variable

		switch {
		case r.ID.IsNaive():
			ps = append(ps, variable.ID(r.ID))
		case r.Name != "":
			ps = append(ps, variable.Name(r.Name))
			if r.ProjectID != "" {
				ps = append(ps, variable.ProjectID(r.ProjectID))

				if r.EnvironmentID != "" {
					ps = append(ps, variable.EnvironmentID(r.EnvironmentID))
				}
			} else {
				ps = append(ps, variable.ProjectIDIsNil())
			}
		}

		if len(ps) == 0 {
			return nil, errors.New("invalid input: illegal predicates")
		}

		// Required.
		c := mc.Variables().Update().
			Where(ps...).
			SetValue(r.Value).
			SetSensitive(r.Sensitive)

		// Optional.
		if r.Description != "" {
			c.SetDescription(r.Description)
		}

		rrs[i] = c
	}

	return rrs, nil
}
