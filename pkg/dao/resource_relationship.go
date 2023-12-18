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

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerelationship"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/strs"
)

// Match both ${svc.name.attr} and ${res.name.attr}.
var resourceRegexp = regexp.MustCompile(`\${(svc|res)\.([^.\s]+)\.[^}]+}`)

func resourceRelationshipCreate(ctx context.Context, mc model.ClientSet, input *model.ResourceRelationship) error {
	err := mc.ResourceRelationships().Create().
		Set(input).
		OnConflict(
			sql.ConflictColumns(
				resourcerelationship.FieldResourceID,
				resourcerelationship.FieldDependencyID,
				resourcerelationship.FieldPath,
			)).
		DoNothing().
		Exec(ctx)
	if err != nil && !errors.Is(err, stdsql.ErrNoRows) {
		return err
	}

	return nil
}

// resourceRelationshipGetDependantPath returns the path from the resource to the end of the path.
func resourceRelationshipGetDependantPath(resourceID object.ID, path []object.ID) []object.ID {
	for i := 0; i < len(path)-1; i++ {
		if path[i] == resourceID {
			return path[i:]
		}
	}

	return nil
}

func resourceRelationshipGetCompletePath(d *model.ResourceRelationship) string {
	ids := make([]string, 0, len(d.Path))
	for _, id := range d.Path {
		ids = append(ids, id.String())
	}

	return strs.Join("/", ids...)
}

// ResourceRelationshipGetDependencyNames gets dependency resource names of the given resource.
func ResourceRelationshipGetDependencyNames(entity *model.Resource) []string {
	names := sets.NewString()

	for _, d := range entity.Attributes {
		matches := resourceRegexp.FindAllSubmatch(d, -1)
		for _, m := range matches {
			names.Insert(string(m[2]))
		}
	}

	return names.List()
}

// GetResourceDependantNames gets names of resources that depends on the given resource.
// ExcludeStatus is an optional string slice of dependant resource status to exclude.
func GetResourceDependantNames(
	ctx context.Context,
	modelClient model.ClientSet,
	entity *model.Resource,
	excludeStatus ...any,
) ([]string, error) {
	var names []string

	err := modelClient.ResourceRelationships().Query().
		Where(
			resourcerelationship.ResourceIDNEQ(entity.ID),
			resourcerelationship.DependencyID(entity.ID),
		).
		Modify(func(s *sql.Selector) {
			t := sql.Table(resource.Table).As("s")
			s.LeftJoin(t).
				On(t.C(resource.FieldID), resourcerelationship.FieldResourceID).
				Select(resource.FieldName).
				Distinct()

			if len(excludeStatus) > 0 {
				s.Where(sqljson.ValueNotIn(
					resource.FieldStatus,
					excludeStatus,
					sqljson.Path("summaryStatus"),
				))
			}
		}).
		Scan(ctx, &names)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource relationships: %w", err)
	}

	return names, nil
}

// GetResourceDependantIDs gets IDs of resources that depend on the given resources.
func GetResourceDependantIDs(ctx context.Context, mc model.ClientSet, resourceIDs ...object.ID) ([]object.ID, error) {
	var ids []object.ID

	err := mc.ResourceRelationships().Query().
		Where(
			resourcerelationship.ResourceIDNotIn(resourceIDs...),
			resourcerelationship.DependencyIDIn(resourceIDs...),
		).
		Select(resourcerelationship.FieldResourceID).
		Scan(ctx, &ids)

	return ids, err
}

// GetNonStoppedResourceDependantIDs gets IDs of resources that depend on the given resources
// and is not in stopped status.
func GetNonStoppedResourceDependantIDs(
	ctx context.Context,
	mc model.ClientSet,
	resourceIDs ...object.ID,
) ([]object.ID, error) {
	var ids []object.ID

	err := mc.ResourceRelationships().Query().
		Where(
			resourcerelationship.ResourceIDNotIn(resourceIDs...),
			resourcerelationship.DependencyIDIn(resourceIDs...),
		).
		Modify(func(s *sql.Selector) {
			t := sql.Table(resource.Table).As("s")
			s.RightJoin(t).
				On(t.C(resource.FieldID), resourcerelationship.FieldResourceID).
				Distinct().
				Where(sqljson.ValueNEQ(
					resource.FieldStatus,
					status.ResourceStatusStopped.String(),
					sqljson.Path("summaryStatus"),
				))
		}).
		Select(resourcerelationship.FieldResourceID).
		Scan(ctx, &ids)

	return ids, err
}

// resourceRelationshipGetDependencies returns the new dependencies of the given resource.
func resourceRelationshipGetDependencies(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.Resource,
) ([]*model.ResourceRelationship, error) {
	preServiceNames := ResourceRelationshipGetDependencyNames(entity)

	// Get the resource IDs of the resource names in same project and environment.
	dependencyResources, err := mc.Resources().Query().
		Where(
			resource.ProjectID(entity.ProjectID),
			resource.EnvironmentID(entity.EnvironmentID),
			resource.NameIn(preServiceNames...),
		).
		WithDependencies().
		All(ctx)
	if err != nil {
		return nil, err
	}

	var dependencyResourceRelationships []*model.ResourceRelationship
	for _, preSvc := range dependencyResources {
		dependencyResourceRelationships = append(dependencyResourceRelationships, preSvc.Edges.Dependencies...)
	}

	dependencies := make(
		[]*model.ResourceRelationship,
		0,
		len(dependencyResourceRelationships)+len(dependencyResources)+1,
	)
	// Add dependency path to the resource itself.
	dependencies = append(dependencies, &model.ResourceRelationship{
		ResourceID:   entity.ID,
		DependencyID: entity.ID,
		Path:         []object.ID{entity.ID},
		Type:         types.ResourceRelationshipTypeImplicit,
	})

	// TODO (alex): handle the case for user add explicit dependency.
	for _, preSvc := range dependencyResources {
		for _, d := range preSvc.Edges.Dependencies {
			path := make([]object.ID, len(d.Path), len(d.Path)+1)
			copy(path, d.Path)
			path = append(path, entity.ID)

			dependency := &model.ResourceRelationship{
				ResourceID:   entity.ID,
				DependencyID: d.ResourceID,
				Type:         types.ResourceRelationshipTypeImplicit,
				Path:         path,
			}

			// Check if there is a dependency cycle.
			if existCycle := ResourceRelationshipCheckCycle(dependency); existCycle {
				return nil, errorx.New("resource dependency contains cycle")
			}

			dependencies = append(dependencies, dependency)
		}
	}

	return dependencies, nil
}

// resourceRelationshipUpdateDependants updates dependant resource dependencies of the resource.
// E.G. Resource C depends on resource B, resource B depends on resource A,
// given resource A as input, resource B and C's dependencies will be updated.
func resourceRelationshipUpdateDependants(ctx context.Context, mc model.ClientSet, s *model.Resource) error {
	// Dependant resources' dependencies of the resource.
	postDependencies, err := mc.ResourceRelationships().Query().
		Where(
			resourcerelationship.ResourceIDNEQ(s.ID),
			func(sq *sql.Selector) {
				sq.Where(sqljson.ValueContains(resourcerelationship.FieldPath, s.ID))
			}).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	newPostDependencies := make([]*model.ResourceRelationship, 0, len(s.Edges.Dependencies))

	for _, sd := range s.Edges.Dependencies {
		for _, osd := range postDependencies {
			path := make([]object.ID, len(sd.Path)-1, len(sd.Path)+len(osd.Path))
			copy(path, sd.Path[:len(sd.Path)-1])
			path = append(path, resourceRelationshipGetDependantPath(s.ID, osd.Path)...)

			newDependency := &model.ResourceRelationship{
				ResourceID:   osd.ResourceID,
				DependencyID: osd.DependencyID,
				Type:         types.ResourceRelationshipTypeImplicit,
				Path:         path,
			}

			if existCycle := ResourceRelationshipCheckCycle(newDependency); existCycle {
				return errorx.New("resource dependency contains cycle")
			}

			newPostDependencies = append(newPostDependencies, newDependency)
		}
	}

	newPostDependencyPath := sets.NewString()
	// Create new dependencies and record new paths for dependant resources.
	for _, d := range newPostDependencies {
		newPostDependencyPath.Insert(resourceRelationshipGetCompletePath(d))

		if err = resourceRelationshipCreate(ctx, mc, d); err != nil {
			return err
		}
	}

	// Delete old dependant resources dependencies.
	for _, d := range postDependencies {
		if newPostDependencyPath.Has(resourceRelationshipGetCompletePath(d)) {
			continue
		}

		var ids []string
		for _, id := range d.Path {
			ids = append(ids, fmt.Sprintf("%q", id.String()))
		}
		paths := fmt.Sprintf("[%s]", strs.Join(",", ids...))

		_, err := mc.ResourceRelationships().Delete().
			Where(
				resourcerelationship.ResourceID(d.ResourceID),
				resourcerelationship.DependencyID(d.DependencyID),
				resourcerelationship.Type(d.Type),
				func(s *sql.Selector) {
					s.Where(sqljson.ValueEQ(resourcerelationship.FieldPath, paths))
				},
			).Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// ResourceRelationshipCheckCycle checks if the path contains a dependency cycle.
func ResourceRelationshipCheckCycle(s *model.ResourceRelationship) bool {
	return sets.New[object.ID](s.Path...).Len() != len(s.Path)
}
