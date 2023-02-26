package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

func EnvironmentCreates(mc model.ClientSet, input ...*model.Environment) ([]*model.EnvironmentCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*model.EnvironmentCreate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// required.
		var c = mc.Environments().Create().
			SetName(r.Name)

		// optional.
		c.SetDescription(r.Description).
			SetVariables(r.Variables)
		if r.Labels != nil {
			c.SetLabels(r.Labels)
		}

		c.AddConnectors(r.Edges.Connectors...)

		rrs[i] = c
	}
	return rrs, nil
}

func EnvironmentUpdates(mc model.ClientSet, input ...*model.Environment) ([]*model.EnvironmentUpdate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*model.EnvironmentUpdate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// predicated.
		var ps []predicate.Environment
		switch {
		case r.ID.IsNaive():
			ps = append(ps, environment.ID(r.ID))
		case r.Name != "":
			ps = append(ps, environment.Name(r.Name))
		}
		if len(ps) == 0 {
			return nil, errors.New("invalid input: illegal predicates")
		}

		// conditional.
		var c = mc.Environments().Update().
			Where(ps...).
			SetDescription(r.Description).
			SetVariables(r.Variables)
		if r.Name != "" {
			c.SetName(r.Name)
		}
		if r.Labels != nil {
			c.SetLabels(r.Labels)
		}

		if r.Edges.Connectors != nil {
			c.ClearConnectors()
			c.AddConnectors(r.Edges.Connectors...)
		}

		rrs[i] = c
	}
	return rrs, nil
}
