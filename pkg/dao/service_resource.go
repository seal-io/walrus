package dao

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/model/serviceresourcerelationship"
	"github.com/seal-io/seal/utils/strs"
)

func ServiceResourceCreates(
	mc model.ClientSet,
	input ...*model.ServiceResource,
) ([]*WrappedServiceResourceCreate, error) {
	creates, err := serviceResourceCreates(mc, input...)
	if err != nil {
		return nil, err
	}

	rrs := make([]*WrappedServiceResourceCreate, len(creates))

	for i, c := range creates {
		rrs[i] = &WrappedServiceResourceCreate{
			ServiceResourceCreate: c,
			entity:                input[i],
		}
	}

	return rrs, nil
}

func serviceResourceCreates(
	mc model.ClientSet,
	input ...*model.ServiceResource,
) ([]*model.ServiceResourceCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.ServiceResourceCreate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.ServiceResources().Create().
			SetProjectID(r.ProjectID).
			SetServiceID(r.ServiceID).
			SetConnectorID(r.ConnectorID).
			SetName(r.Name).
			SetType(r.Type).
			SetMode(r.Mode).
			SetDeployerType(r.DeployerType).
			SetShape(r.Shape)

		// Optional.
		if r.CompositionID.Valid(0) {
			c.SetCompositionID(r.CompositionID)
		}

		if r.ClassID.Valid(0) {
			c.SetClassID(r.ClassID)
		}

		rrs[i] = c
	}

	return rrs, nil
}

type WrappedServiceResourceCreate struct {
	*model.ServiceResourceCreate

	entity *model.ServiceResource
}

func (r *WrappedServiceResourceCreate) Save(ctx context.Context) (created *model.ServiceResource, err error) {
	mc := r.ServiceResourceCreate.Mutation().Client()

	created, err = r.ServiceResourceCreate.Save(ctx)
	if err != nil {
		return
	}

	// Create instance services.
	if len(r.entity.Edges.Instances) > 0 {
		instances := make(model.ServiceResources, len(r.entity.Edges.Instances))

		for i, r := range r.entity.Edges.Instances {
			instances[i] = r
			instances[i].ClassID = created.ID
		}

		var instanceCreates []*model.ServiceResourceCreate

		instanceCreates, err = serviceResourceCreates(mc, instances...)
		if err != nil {
			return nil, err
		}

		for i, c := range instanceCreates {
			var instance *model.ServiceResource

			instance, err = c.Save(ctx)
			if err != nil {
				return
			}
			instances[i] = instance
		}
		created.Edges.Instances = instances
	}

	return
}

// ServiceResourceShapeClassQuery returns a query that selects a shape class service resource,
// components and dependencies.
func ServiceResourceShapeClassQuery(query *model.ServiceResourceQuery, fields ...string) *model.ServiceResourceQuery {
	var (
		order  = model.Desc(serviceresource.FieldCreateTime)
		wcOpts = func(q *model.ConnectorQuery) {
			q.Select(
				connector.FieldName,
				connector.FieldType,
				connector.FieldCategory,
				connector.FieldConfigVersion,
				connector.FieldConfigData,
			)
		}
	)

	return query.Select(fields...).
		WithInstances(func(srq *model.ServiceResourceQuery) {
			srq.WithConnector(wcOpts).
				WithComponents(func(q *model.ServiceResourceQuery) {
					q.Order(order).
						Select(fields...).
						WithConnector(wcOpts)
				})
		}).WithDependencies(func(rrq *model.ServiceResourceRelationshipQuery) {
		rrq.Select(
			serviceresourcerelationship.FieldServiceResourceID,
			serviceresourcerelationship.FieldDependencyID,
			serviceresourcerelationship.FieldType,
		)
	})
}

// ServiceResourceToMap recursive set a map of service resources indexed by its unique index.
func ServiceResourceToMap(resources []*model.ServiceResource) map[string]*model.ServiceResource {
	m := make(map[string]*model.ServiceResource)

	stack := make([]*model.ServiceResource, 0)
	stack = append(stack, resources...)

	for len(stack) > 0 {
		res := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		key := ServiceResourceGetUniqueKey(res)
		if _, ok := m[key]; ok {
			continue
		}
		m[key] = res

		stack = append(stack, res.Edges.Components...)
		stack = append(stack, res.Edges.Instances...)
	}

	return m
}

// ServiceResourceGetUniqueKey returns the unique index key of the given model.ServiceResource.
func ServiceResourceGetUniqueKey(r *model.ServiceResource) string {
	// Align to schema definition.
	return strs.Join("-", string(r.ConnectorID), r.Shape, r.Mode, r.Type, r.Name)
}
