package dao

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
)

// WrappedProjectCreate is a wrapper for model.ProjectCreate
// to process the relationship with model.Role.
type WrappedProjectCreate struct {
	*model.ProjectCreate

	entity *model.Project
}

func (sc *WrappedProjectCreate) Save(ctx context.Context) (created *model.Project, err error) {
	mc := sc.ProjectCreate.Mutation().Client()

	// Save entity.
	created, err = sc.ProjectCreate.Save(ctx)
	if err != nil {
		return
	}

	// Construct relationships.
	newRss := sc.entity.Edges.SubjectRoles
	createRss := make([]*model.SubjectRoleRelationshipCreate, len(newRss))

	for i, rs := range newRss {
		if rs == nil {
			return nil, errors.New("invalid input: nil relationship")
		}

		// Required.
		c := mc.SubjectRoleRelationships().Create().
			SetProjectID(created.ID).
			SetSubjectID(rs.SubjectID).
			SetRoleID(rs.RoleID)

		createRss[i] = c
	}

	// Save relationships.
	newRss, err = mc.SubjectRoleRelationships().CreateBulk(createRss...).
		Save(ctx)
	if err != nil {
		return
	}
	created.Edges.SubjectRoles = newRss

	return created, nil
}

func (sc *WrappedProjectCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

func ProjectCreates(mc model.ClientSet, input ...*model.Project) ([]*WrappedProjectCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*WrappedProjectCreate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.Projects().Create().
			SetName(r.Name)

		// Optional.
		c.SetDescription(r.Description)

		if r.Labels != nil {
			c.SetLabels(r.Labels)
		}
		rrs[i] = &WrappedProjectCreate{
			ProjectCreate: c,
			entity:        input[i],
		}
	}

	return rrs, nil
}

func ProjectUpdates(mc model.ClientSet, input ...*model.Project) ([]*model.ProjectUpdate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.ProjectUpdate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Predicated.
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

		// Conditional.
		c := mc.Projects().Update().
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
