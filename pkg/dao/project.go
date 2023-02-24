package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
)

func ProjectCreates(mc model.ClientSet, input ...*model.Project) ([]*model.ProjectCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*model.ProjectCreate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// required.
		var c = mc.Projects().Create().
			SetName(r.Name)

		// optional.
		c.SetDescription(r.Description)
		if r.Labels != nil {
			c.SetLabels(r.Labels)
		}
		rrs[i] = c
	}
	return rrs, nil
}

func ProjectUpdates(mc model.ClientSet, input ...*model.Project) ([]*model.ProjectUpdate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*model.ProjectUpdate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// predicated.
		var ps []predicate.Project
		switch {
		case r.ID.IsNaive():
			ps = append(ps, project.ID(r.ID))
		case r.Name != "":
			ps = append(ps, project.Name(r.Name))
		}
		if len(ps) == 0 {
			return nil, errors.New("invalid input: illegal predicates")
		}

		// conditional.
		var c = mc.Projects().Update().
			Where(ps...).
			SetDescription(r.Description)
		if r.Name != "" {
			c.SetName(r.Name)
		}
		if r.Labels != nil {
			c.SetLabels(r.Labels)
		}
		rrs[i] = c
	}
	return rrs, nil
}
