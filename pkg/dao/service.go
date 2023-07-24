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
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/servicerelationship"
	"github.com/seal-io/seal/pkg/dao/model/servicerevision"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/utils/strs"
)

// ServiceDependenciesEdgeSave saves the edge dependencies of model.Service entity.
func ServiceDependenciesEdgeSave(ctx context.Context, mc model.ClientSet, entity *model.Service) error {
	// Default new items and create key set for new items.
	newItems, err := serviceRelationshipGetDependencies(ctx, mc, entity)
	if err != nil {
		return err
	}

	newItemKeySet := sets.New[string]()

	for i := range newItems {
		if newItems[i] == nil {
			return errors.New("invalid input: nil relationship")
		}

		newItemKeySet.Insert(serviceRelationshipGetCompletePath(newItems[i]))
	}

	// Add/Update new items.
	if len(newItems) != 0 {
		err = mc.ServiceRelationships().CreateBulk().
			Set(newItems...).
			OnConflict(
				sql.ConflictColumns(
					servicerelationship.FieldServiceID,
					servicerelationship.FieldDependencyID,
					servicerelationship.FieldPath,
				)).
			DoNothing().
			Exec(ctx)
		if err != nil && !errors.Is(err, stdsql.ErrNoRows) {
			return err
		}

		entity.Edges.Dependencies = newItems // Feedback.
	}

	// Delete stale items.
	oldItems, err := mc.ServiceRelationships().Query().
		Where(servicerelationship.ServiceID(entity.ID)).
		All(ctx)
	if err != nil {
		return err
	}

	for i := range oldItems {
		if newItemKeySet.Has(serviceRelationshipGetCompletePath(oldItems[i])) {
			continue
		}

		prefixPath := fmt.Sprintf(`["%s"]`, strs.Join(`","`, oldItems[i].Path...))

		_, err = mc.ServiceRelationships().Delete().
			Where(
				func(s *sql.Selector) {
					s.Where(sqljson.ValueContains(servicerelationship.FieldPath, prefixPath))
				}).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	// Recreate dependants.
	return serviceRelationshipUpdateDependants(ctx, mc, entity)
}

func GetLatestRevisions(
	ctx context.Context,
	modelClient model.ClientSet,
	serviceIDs ...object.ID,
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

func GetServiceNamesByIDs(ctx context.Context, modelClient model.ClientSet, serviceIDs ...object.ID) ([]string, error) {
	var names []string
	err := modelClient.Services().Query().
		Where(service.IDIn(serviceIDs...)).
		Select(service.FieldName).
		Scan(ctx, &names)

	return names, err
}
