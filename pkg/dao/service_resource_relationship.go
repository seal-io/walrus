package dao

import (
	"context"
	stdsql "database/sql"
	"errors"

	"entgo.io/ent/dialect/sql"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/serviceresourcerelationship"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

// ServiceResourceRelationshipUpdateWithDependencies updates the relationship with dependencies and resources.
func ServiceResourceRelationshipUpdateWithDependencies(
	ctx context.Context,
	mc model.ClientSet,
	dependencies map[string][]string,
	recordResources,
	createResources model.ServiceResources,
) error {
	logger := log.WithName("dao").WithName("service-revision")

	resourceMap := ServiceResourceToMap(append(recordResources, createResources...))

	recordResourceIDs := make([]object.ID, 0, len(recordResources))
	for _, r := range recordResources {
		recordResourceIDs = append(recordResourceIDs, r.ID)
	}

	recordResourceRelationships, err := mc.ServiceResourceRelationships().Query().
		Where(serviceresourcerelationship.ServiceResourceIDIn(recordResourceIDs...)).
		All(ctx)
	if err != nil {
		return err
	}

	relationshipSets := sets.NewString()
	for _, r := range recordResourceRelationships {
		relationshipSets.Insert(strs.Join("/", r.ServiceResourceID.String(), r.DependencyID.String(), r.Type))
	}

	// Create resource relationships with dependencies.
	createResourceRelationships := make(model.ServiceResourceRelationships, 0, len(dependencies))

	for k, deps := range dependencies {
		r, ok := resourceMap[k]
		if !ok {
			logger.Warnf("resource not found, resource index: %s", k)
			continue
		}

		for _, k := range deps {
			dr, ok := resourceMap[k]
			if !ok {
				logger.Warnf("dependant resource not found, resource index: %s", k)
				continue
			}

			// Skip if relationship already exists.
			if relationshipSets.Has(strs.Join("/", r.ID.String(), dr.ID.String(), dr.Type)) {
				continue
			}

			// Skip composition or realization relationship.
			if dr.ID == r.CompositionID || dr.ID == r.ClassID {
				continue
			}

			createResourceRelationships = append(
				createResourceRelationships,
				&model.ServiceResourceRelationship{
					ServiceResourceID: r.ID,
					DependencyID:      dr.ID,
					Type:              types.ServiceResourceRelationshipTypeDependency,
				})
		}
	}

	if len(createResourceRelationships) > 0 {
		// Create relationships.
		err = mc.ServiceResourceRelationships().CreateBulk().
			Set(createResourceRelationships...).
			OnConflict(
				sql.ConflictColumns(
					serviceresourcerelationship.FieldServiceResourceID,
					serviceresourcerelationship.FieldDependencyID,
					serviceresourcerelationship.FieldType,
				)).
			DoNothing().
			Exec(ctx)
		if err != nil && !errors.Is(err, stdsql.ErrNoRows) {
			return err
		}
	}

	return nil
}
