package dao

import (
	"context"
	stdsql "database/sql"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/servicedependency"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/utils/strs"
)

func ServiceCreates(
	mc model.ClientSet,
	input ...*model.Service,
) ([]*WrappedServiceCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*WrappedServiceCreate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.Services().Create().
			SetProjectID(r.ProjectID).
			SetEnvironmentID(r.EnvironmentID).
			SetTemplate(r.Template).
			SetName(r.Name)

		// Optional.
		c.SetDescription(r.Description)

		if r.Labels != nil {
			c.SetLabels(r.Labels)
		}

		if r.Attributes != nil {
			c.SetAttributes(r.Attributes)
		}

		status.ServiceStatusDeployed.Unknown(r, "Deploying service")
		r.Status.SetSummary(status.WalkService(&r.Status))
		c.SetStatus(r.Status)

		rrs[i] = &WrappedServiceCreate{
			ServiceCreate: c,
			entity:        input[i],
		}
	}

	return rrs, nil
}

func ServiceUpdate(
	mc model.ClientSet,
	input *model.Service,
) (*WrappedServiceUpdate, error) {
	if input == nil {
		return nil, errors.New("invalid input: nil entity")
	}

	if input.ID == "" {
		return nil, errors.New("invalid input: illegal predicates")
	}
	ps := []predicate.Service{
		service.ID(input.ID),
	}

	c := mc.Services().UpdateOne(input).
		SetTemplate(input.Template)

	c.SetDescription(input.Description)

	if input.Labels != nil {
		c.SetLabels(input.Labels)
	}

	if input.Attributes != nil {
		c.SetAttributes(input.Attributes)
	}

	input.Status.SetSummary(status.WalkService(&input.Status))

	if input.Status.Changed() {
		c.SetStatus(input.Status)
	}

	return &WrappedServiceUpdate{
		ServiceUpdateOne: c,
		entity:           input,
		entityPredicates: ps,
	}, nil
}

// ServiceStatusUpdate updates the status of the given service.
// TODO (alex): unify the status update logic for all services.
func ServiceStatusUpdate(
	mc model.ClientSet,
	input *model.Service,
) (*model.ServiceUpdateOne, error) {
	if input == nil {
		return nil, errors.New("invalid input: nil entity")
	}

	if input.ID == "" {
		return nil, errors.New("invalid input: illegal predicates")
	}

	c := mc.Services().UpdateOne(input)
	input.Status.SetSummary(status.WalkService(&input.Status))

	if input.Status.Changed() {
		c.SetStatus(input.Status)
	}

	return c, nil
}

type WrappedServiceCreate struct {
	*model.ServiceCreate

	entity *model.Service
}

func (sc *WrappedServiceCreate) Save(ctx context.Context) (created *model.Service, err error) {
	mc := sc.ServiceCreate.Mutation().Client()

	// Save entity.
	created, err = sc.ServiceCreate.Save(ctx)
	if err != nil {
		return
	}

	// Construct dependencies.
	// New dependencies.
	dependencies, err := GetNewDependencies(ctx, mc, created)
	if err != nil {
		return nil, err
	}

	for _, d := range dependencies {
		c := mc.ServiceDependencies().Create().
			SetServiceID(d.ServiceID).
			SetDependentID(d.DependentID).
			SetType(d.Type).
			SetPath(d.Path)

		err = c.OnConflict(
			sql.ConflictColumns(
				servicedependency.FieldServiceID,
				servicedependency.FieldDependentID,
				servicedependency.FieldPath,
			)).
			DoNothing().
			Exec(ctx)
		if err != nil && !errors.Is(err, stdsql.ErrNoRows) {
			return nil, err
		}
	}

	created.Edges.Dependencies = dependencies
	if err = UpdateDependants(ctx, mc, created); err != nil {
		return nil, err
	}

	return created, nil
}

func (sc *WrappedServiceCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

type WrappedServiceUpdate struct {
	*model.ServiceUpdateOne

	entity           *model.Service
	entityPredicates []predicate.Service
}

func (su *WrappedServiceUpdate) Save(ctx context.Context) (*model.Service, error) {
	mc := su.ServiceUpdateOne.Mutation().Client()

	// Update entity.
	update, err := su.ServiceUpdateOne.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Get old dependencies.
	oldDependencies, err := mc.ServiceDependencies().Query().
		Where(servicedependency.ServiceID(update.ID)).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	// New dependencies.
	dependencies, err := GetNewDependencies(ctx, mc, update)
	if err != nil {
		return nil, err
	}

	newDependencyPath := sets.NewString()
	for _, d := range dependencies {
		newDependencyPath.Insert(getCompleteDependencyPath(d))

		c := mc.ServiceDependencies().Create().
			SetServiceID(d.ServiceID).
			SetDependentID(d.DependentID).
			SetType(d.Type).
			SetPath(d.Path)

		err = c.OnConflict(
			sql.ConflictColumns(
				servicedependency.FieldServiceID,
				servicedependency.FieldDependentID,
				servicedependency.FieldPath,
			)).
			DoNothing().
			Exec(ctx)

		if err != nil && !errors.Is(err, stdsql.ErrNoRows) {
			return nil, err
		}
	}

	// Delete old dependencies.
	for _, d := range oldDependencies {
		if newDependencyPath.Has(getCompleteDependencyPath(d)) {
			continue
		}

		var ids []string
		for _, id := range d.Path {
			ids = append(ids, fmt.Sprintf("%q", id.String()))
		}
		paths := fmt.Sprintf("[%s]", strs.Join(",", ids...))

		_, err := mc.ServiceDependencies().Delete().
			Where(
				servicedependency.ServiceID(d.ServiceID),
				servicedependency.DependentID(d.DependentID),
				servicedependency.Type(d.Type),
				func(s *sql.Selector) {
					s.Where(sqljson.ValueIn(servicedependency.FieldPath, []any{paths}))
				},
			).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	update.Edges.Dependencies = dependencies

	if err = UpdateDependants(ctx, mc, update); err != nil {
		return nil, err
	}

	return update, nil
}
