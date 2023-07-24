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

// ServiceResourceInstancesEdgeSave saves the edge instances of model.ServiceResource entity.
func ServiceResourceInstancesEdgeSave(ctx context.Context, mc model.ClientSet, entity *model.ServiceResource) error {
	if entity.Edges.Instances == nil {
		return nil
	}

	// Delete stale items.
	_, err := mc.ServiceResources().Delete().
		Where(serviceresource.ClassID(entity.ID)).
		Exec(ctx)
	if err != nil {
		return err
	}

	// Add new items.
	newItems := entity.Edges.Instances
	for i := range newItems {
		if newItems[i] == nil {
			return errors.New("invalid input: nil relationship")
		}
		newItems[i].ClassID = entity.ID
	}

	newItems, err = mc.ServiceResources().CreateBulk().
		Set(newItems...).
		Save(ctx)
	if err != nil {
		return err
	}

	entity.Edges.Instances = newItems // Feedback.

	return nil
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
