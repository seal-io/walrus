package dao

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resourcecomponent"
	"github.com/seal-io/walrus/pkg/dao/model/resourcecomponentrelationship"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/utils/strs"
)

// ResourceComponentInstancesEdgeSave saves the edge instances of model.ResourceComponent entity.
func ResourceComponentInstancesEdgeSave(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.ResourceComponent,
) error {
	_, err := ResourceComponentInstancesEdgeSaveWithResult(ctx, mc, entity)
	return err
}

// ResourceComponentInstancesEdgeSaveWithResult saves the edge instances of model.ResourceComponent entity
// and return new created components.
func ResourceComponentInstancesEdgeSaveWithResult(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.ResourceComponent,
) ([]*model.ResourceComponent, error) {
	if entity.Edges.Instances == nil {
		return nil, nil
	}

	// Delete stale items.
	_, err := mc.ResourceComponents().Delete().
		Where(resourcecomponent.ClassID(entity.ID)).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Add new items.
	newItems := entity.Edges.Instances
	for i := range newItems {
		if newItems[i] == nil {
			return nil, errors.New("invalid input: nil relationship")
		}
		newItems[i].ClassID = entity.ID
	}

	newItems, err = mc.ResourceComponents().CreateBulk().
		Set(newItems...).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	entity.Edges.Instances = newItems // Feedback.

	return newItems, nil
}

// ResourceComponentShapeClassQuery wraps the given model.ResourceComponent query
// to select all shape class resources and the owned components and dependencies of them.
func ResourceComponentShapeClassQuery(query *model.ResourceComponentQuery) *model.ResourceComponentQuery {
	order := model.Desc(resourcecomponent.FieldCreateTime)

	return query.
		Where(
			resourcecomponent.Shape(types.ResourceComponentShapeClass),
			resourcecomponent.Mode(types.ResourceComponentModeManaged)).
		Order(order).
		WithInstances(func(iq *model.ResourceComponentQuery) {
			iq.
				Order(order).
				WithComponents(func(cq *model.ResourceComponentQuery) {
					cq.
						Order(order)
				})
		}).
		WithDependencies(func(rrq *model.ResourceComponentRelationshipQuery) {
			rrq.Select(
				resourcecomponentrelationship.FieldResourceComponentID,
				resourcecomponentrelationship.FieldDependencyID,
				resourcecomponentrelationship.FieldType)
		})
}

// ResourceComponentToMap recursive set a map of resource components indexed by its unique index.
func ResourceComponentToMap(resources []*model.ResourceComponent) map[string]*model.ResourceComponent {
	m := make(map[string]*model.ResourceComponent)

	stack := make([]*model.ResourceComponent, 0)
	stack = append(stack, resources...)

	for len(stack) > 0 {
		res := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		key := ResourceComponentGetUniqueKey(res)
		if _, ok := m[key]; ok {
			continue
		}
		m[key] = res

		stack = append(stack, res.Edges.Components...)
		stack = append(stack, res.Edges.Instances...)
	}

	return m
}

// ResourceComponentGetUniqueKey returns the unique index key of the given model.ResourceComponent.
func ResourceComponentGetUniqueKey(r *model.ResourceComponent) string {
	// Align to schema definition.
	return strs.Join("-", string(r.ConnectorID), r.Shape, r.Mode, r.Type, r.Name)
}

// GetCleanResourceComponents get the components that exist in the database.
// It is used to filter the components that do not exist in the database.
func GetCleanResourceComponents(
	ctx context.Context,
	mc model.ClientSet,
	components model.ResourceComponents,
) (model.ResourceComponents, error) {
	ids := make([]object.ID, 0, len(components))
	for _, c := range components {
		ids = append(ids, c.ID)
	}

	// Get exist component IDs.
	existComponentIDs, err := mc.ResourceComponents().Query().
		Where(resourcecomponent.IDIn(ids...)).
		IDs(ctx)
	if err != nil {
		return nil, err
	}

	existComponentIDSet := sets.New[object.ID](existComponentIDs...)

	// Clean components.
	cleanComponents := make(model.ResourceComponents, 0)

	for _, c := range components {
		if existComponentIDSet.Has(c.ID) {
			cleanComponents = append(cleanComponents, c)
		}
	}

	return cleanComponents, nil
}
