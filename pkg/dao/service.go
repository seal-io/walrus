package dao

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/servicerelationship"
	"github.com/seal-io/seal/pkg/dao/model/servicerevision"
	"github.com/seal-io/seal/pkg/dao/types/oid"
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

		if r.Annotations != nil {
			c.SetAnnotations(r.Annotations)
		}

		status.ServiceStatusDeployed.Unknown(r, "")
		r.Status.SetSummary(status.WalkService(&r.Status))
		c.SetStatus(r.Status)

		rrs[i] = &WrappedServiceCreate{
			ServiceCreate: c,
			entity:        input[i],
		}
	}

	return rrs, nil
}

// ServiceUpdates returns a slice of wrapped service update builder.
func ServiceUpdates(
	mc model.ClientSet,
	input ...*model.Service,
) ([]*WrappedServiceUpdate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*WrappedServiceUpdate, len(input))

	for i, r := range input {
		if r.ID == "" {
			return nil, errors.New("invalid input: illegal predicates")
		}
		ps := []predicate.Service{
			service.ID(r.ID),
		}

		c := mc.Services().UpdateOne(r).
			SetTemplate(r.Template)

		c.SetDescription(r.Description)

		if r.Labels != nil {
			c.SetLabels(r.Labels)
		}

		if r.Annotations != nil {
			c.SetAnnotations(r.Annotations)
		}

		if r.Attributes != nil {
			c.SetAttributes(r.Attributes)
		}

		r.Status.SetSummary(status.WalkService(&r.Status))

		if r.Status.Changed() {
			c.SetStatus(r.Status)
		}

		rrs[i] = &WrappedServiceUpdate{
			ServiceUpdateOne: c,
			entity:           input[i],
			entityPredicates: ps,
		}
	}

	return rrs, nil
}

// ServiceUpdate returns a wrapped service update builder.
func ServiceUpdate(
	mc model.ClientSet,
	input *model.Service,
) (*WrappedServiceUpdate, error) {
	updates, err := ServiceUpdates(mc, input)
	if err != nil {
		return nil, err
	}

	return updates[0], err
}

func GetLatestRevisions(
	ctx context.Context,
	modelClient model.ClientSet,
	serviceIDs ...oid.ID,
) ([]*model.ServiceRevision, error) {
	// Get the latest revisions of given services by the following sql:
	// SELECT service_revisions.*
	// FROM service_revisions
	// JOIN (
	// 	 SELECT service_id, MAX(create_time) AS create_time FROM service_revisions GROUP BY service_id
	// ) t
	// ON service_revisions.service_id=t.service_id
	// AND service_revisions.create_time=t.create_time
	// WHERE service_revisions.service_id IN (...)
	ids := make([]any, len(serviceIDs))
	for i := range serviceIDs {
		ids[i] = serviceIDs[i]
	}

	return modelClient.ServiceRevisions().Query().
		Modify(func(s *sql.Selector) {
			t := sql.Select(
				servicerevision.FieldServiceID,
				sql.As(sql.Max(servicerevision.FieldCreateTime), servicerevision.FieldCreateTime),
			).
				From(sql.Table(servicerevision.Table)).
				GroupBy(servicerevision.FieldServiceID).
				As("t")
			s.Join(t).
				OnP(
					sql.And(
						sql.ColumnsEQ(
							s.C(servicerevision.FieldServiceID),
							t.C(servicerevision.FieldServiceID),
						),
						sql.ColumnsEQ(
							s.C(servicerevision.FieldCreateTime),
							t.C(servicerevision.FieldCreateTime),
						),
					),
				).
				Where(
					sql.In(s.C(servicerevision.FieldServiceID), ids...),
				)
		}).
		WithService(func(sq *model.ServiceQuery) {
			sq.Select(
				service.FieldName,
			)
		}).
		All(ctx)
}

func GetServiceNamesByIDs(ctx context.Context, modelClient model.ClientSet, serviceIDs ...oid.ID) ([]string, error) {
	var names []string
	err := modelClient.Services().Query().
		Where(service.IDIn(serviceIDs...)).
		Select(service.FieldName).
		Scan(ctx, &names)

	return names, err
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
	dependencies, err := serviceRelationshipGetDependencies(ctx, mc, created)
	if err != nil {
		return nil, err
	}

	for _, d := range dependencies {
		if err = serviceRelationshipCreate(ctx, mc, d); err != nil {
			return nil, err
		}
	}

	created.Edges.Dependencies = dependencies
	if err = serviceRelationshipUpdateDependants(ctx, mc, created); err != nil {
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
	oldDependencies, err := mc.ServiceRelationships().Query().
		Where(servicerelationship.ServiceID(update.ID)).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	// New dependencies.
	dependencies, err := serviceRelationshipGetDependencies(ctx, mc, update)
	if err != nil {
		return nil, err
	}

	newDependencyPath := sets.NewString()
	// Create new dependencies.
	for _, d := range dependencies {
		newDependencyPath.Insert(serviceRelationshipGetCompletePath(d))

		if err = serviceRelationshipCreate(ctx, mc, d); err != nil {
			return nil, err
		}
	}

	// Delete old dependencies.
	for _, d := range oldDependencies {
		if newDependencyPath.Has(serviceRelationshipGetCompletePath(d)) {
			continue
		}

		var ids []string
		for _, id := range d.Path {
			ids = append(ids, fmt.Sprintf("%q", id.String()))
		}
		paths := fmt.Sprintf("[%s]", strs.Join(",", ids...))

		_, err := mc.ServiceRelationships().Delete().
			Where(
				servicerelationship.ServiceID(d.ServiceID),
				servicerelationship.DependencyID(d.DependencyID),
				servicerelationship.Type(d.Type),
				func(s *sql.Selector) {
					s.Where(sqljson.ValueIn(servicerelationship.FieldPath, []any{paths}))
				},
			).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	update.Edges.Dependencies = dependencies

	if err = serviceRelationshipUpdateDependants(ctx, mc, update); err != nil {
		return nil, err
	}

	return update, nil
}
