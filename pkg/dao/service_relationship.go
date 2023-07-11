package dao

import (
	"context"
	stdsql "database/sql"
	"errors"
	"fmt"
	"regexp"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/servicerelationship"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/utils/strs"
)

var serviceRegexp = regexp.MustCompile(`\${service\.([^.\s]+)\.[^}]+}`)

func serviceRelationshipCreate(ctx context.Context, mc model.ClientSet, input *model.ServiceRelationship) error {
	c := mc.ServiceRelationships().Create().
		SetServiceID(input.ServiceID).
		SetDependencyID(input.DependencyID).
		SetType(input.Type).
		SetPath(input.Path)

	err := c.OnConflict(
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

	return nil
}

// serviceRelationshipGetDependantPath returns the path from the service to the end of the path.
func serviceRelationshipGetDependantPath(serviceID oid.ID, path []oid.ID) []oid.ID {
	for i := 0; i < len(path)-1; i++ {
		if path[i] == serviceID {
			return path[i:]
		}
	}

	return nil
}

func serviceRelationshipGetCompletePath(d *model.ServiceRelationship) string {
	ids := make([]string, 0, len(d.Path))
	for _, id := range d.Path {
		ids = append(ids, id.String())
	}

	return strs.Join("/", ids...)
}

// ServiceRelationshipGetDependencyNames gets dependency service names of the given service.
func ServiceRelationshipGetDependencyNames(entity *model.Service) []string {
	names := sets.NewString()

	for _, d := range entity.Attributes {
		matches := serviceRegexp.FindAllSubmatch(d, -1)
		for _, m := range matches {
			names.Insert(string(m[1]))
		}
	}

	return names.List()
}

// GetServiceDependantNames gets names of services that depends on the given service.
func GetServiceDependantNames(
	ctx context.Context,
	modelClient model.ClientSet,
	entity *model.Service,
) ([]string, error) {
	var names []string

	err := modelClient.ServiceRelationships().Query().
		Where(
			servicerelationship.ServiceIDNEQ(entity.ID),
			servicerelationship.DependencyID(entity.ID),
		).
		Modify(func(s *sql.Selector) {
			t := sql.Table(service.Table).As("s")
			s.LeftJoin(t).
				On(t.C(service.FieldID), servicerelationship.FieldServiceID).
				Select(service.FieldName)
		}).
		Scan(ctx, &names)
	if err != nil {
		return nil, runtime.Errorw(err, "failed to get service relationships")
	}

	return names, nil
}

// GetServiceDependantIDs gets IDs of services that depend on the given services.
func GetServiceDependantIDs(ctx context.Context, mc model.ClientSet, serviceIDs ...oid.ID) ([]oid.ID, error) {
	var ids []oid.ID

	err := mc.ServiceRelationships().Query().
		Where(
			servicerelationship.ServiceIDNotIn(serviceIDs...),
			servicerelationship.DependencyIDIn(serviceIDs...),
		).
		Select(servicerelationship.FieldServiceID).
		Scan(ctx, &ids)

	return ids, err
}

// serviceRelationshipGetDependencies returns the new dependencies of the given service.
func serviceRelationshipGetDependencies(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.Service,
) ([]*model.ServiceRelationship, error) {
	preServiceNames := ServiceRelationshipGetDependencyNames(entity)

	// Get the service IDs of the service names in same project and environment.
	dependencyServices, err := mc.Services().Query().
		Where(
			service.ProjectID(entity.ProjectID),
			service.EnvironmentID(entity.EnvironmentID),
			service.NameIn(preServiceNames...),
		).
		WithDependencies().
		All(ctx)
	if err != nil {
		return nil, err
	}

	var dependencyServiceRelationships []*model.ServiceRelationship
	for _, preSvc := range dependencyServices {
		dependencyServiceRelationships = append(dependencyServiceRelationships, preSvc.Edges.Dependencies...)
	}

	dependencies := make([]*model.ServiceRelationship, 0, len(dependencyServiceRelationships)+len(dependencyServices)+1)
	// Add dependency path to the service itself.
	dependencies = append(dependencies, &model.ServiceRelationship{
		ServiceID:    entity.ID,
		DependencyID: entity.ID,
		Path:         []oid.ID{entity.ID},
		Type:         types.ServiceRelationshipTypeImplicit,
	})

	// TODO (alex): handle the case for user add explicit dependency.
	for _, preSvc := range dependencyServices {
		for _, d := range preSvc.Edges.Dependencies {
			path := make([]oid.ID, len(d.Path), len(d.Path)+1)
			copy(path, d.Path)
			path = append(path, entity.ID)

			denpendency := &model.ServiceRelationship{
				ServiceID:    entity.ID,
				DependencyID: d.ServiceID,
				Type:         types.ServiceRelationshipTypeImplicit,
				Path:         path,
			}

			// Check if there is a dependency cycle.
			if existCycle := ServiceRelationshipCheckCycle(denpendency); existCycle {
				return nil, errors.New("service dependency contains cycle")
			}

			dependencies = append(dependencies, denpendency)
		}
	}

	return dependencies, nil
}

// serviceRelationshipUpdateDependants updates dependant service dependencies of the service.
// E.G. Service C depends on service B, service B depends on service A,
// given service A as input, service B and C's dependencies will be updated.
func serviceRelationshipUpdateDependants(ctx context.Context, mc model.ClientSet, s *model.Service) error {
	// Dependant services' dependencies of the service.
	postDependencies, err := mc.ServiceRelationships().Query().
		Where(
			servicerelationship.ServiceIDNEQ(s.ID),
			func(sq *sql.Selector) {
				sq.Where(sqljson.ValueContains(servicerelationship.FieldPath, s.ID.String()))
			}).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	newPostDependencies := make([]*model.ServiceRelationship, 0, len(s.Edges.Dependencies))

	for _, sd := range s.Edges.Dependencies {
		for _, osd := range postDependencies {
			path := make([]oid.ID, len(sd.Path)-1, len(sd.Path)+len(osd.Path))
			copy(path, sd.Path[:len(sd.Path)-1])
			path = append(path, serviceRelationshipGetDependantPath(s.ID, osd.Path)...)

			newDependency := &model.ServiceRelationship{
				ServiceID:    osd.ServiceID,
				DependencyID: osd.DependencyID,
				Type:         types.ServiceRelationshipTypeImplicit,
				Path:         path,
			}

			if existCycle := ServiceRelationshipCheckCycle(newDependency); existCycle {
				return errors.New("service dependency contains cycle")
			}

			newPostDependencies = append(newPostDependencies, newDependency)
		}
	}

	newPostDependencyPath := sets.NewString()
	// Create new dependencies and record new paths for dependant services.
	for _, d := range newPostDependencies {
		newPostDependencyPath.Insert(serviceRelationshipGetCompletePath(d))

		if err = serviceRelationshipCreate(ctx, mc, d); err != nil {
			return err
		}
	}

	// Delete old dependant services dependencies.
	for _, d := range postDependencies {
		if newPostDependencyPath.Has(serviceRelationshipGetCompletePath(d)) {
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
					s.Where(sqljson.ValueEQ(servicerelationship.FieldPath, paths))
				},
			).Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// ServiceRelationshipCheckCycle checks if the path contains a dependency cycle.
func ServiceRelationshipCheckCycle(s *model.ServiceRelationship) bool {
	return sets.New[oid.ID](s.Path...).Len() != len(s.Path)
}
