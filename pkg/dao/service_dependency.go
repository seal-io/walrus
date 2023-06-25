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

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/servicedependency"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/utils/strs"
)

var serviceRegexp = regexp.MustCompile(`\${service\.([^.\s]+)\.[^}]+}`)

func GetDependencyNames(entity *model.Service) []string {
	var names []string

	for _, d := range entity.Attributes {
		matches := serviceRegexp.FindAllSubmatch(d, -1)
		for _, m := range matches {
			names = append(names, string(m[1]))
		}
	}

	return names
}

// GetNewDependencies returns the new dependencies of the given service.
func GetNewDependencies(
	ctx context.Context,
	mc model.ClientSet,
	entity *model.Service,
) ([]*model.ServiceDependency, error) {
	var (
		dependencies []*model.ServiceDependency
		serviceNames = GetDependencyNames(entity)
	)

	// Get the service IDs of the service names in same project and environment.
	serviceIDs, err := mc.Services().Query().
		Where(
			service.NameIn(serviceNames...),
			service.EnvironmentID(entity.EnvironmentID),
			service.ProjectID(entity.ProjectID),
		).
		IDs(ctx)
	if err != nil {
		return nil, err
	}

	parentDependencies, err := mc.ServiceDependencies().Query().
		Where(servicedependency.ServiceIDIn(serviceIDs...)).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	// TODO (alex): handle the case for user add explicit dependency.
	if len(parentDependencies) == 0 {
		for _, id := range serviceIDs {
			dependencies = append(dependencies, &model.ServiceDependency{
				ServiceID:   entity.ID,
				DependentID: id,
				Path:        []oid.ID{id, entity.ID},
				Type:        types.ServiceDependencyTypeImplicit,
			})
		}
	} else {
		for _, d := range parentDependencies {
			path := d.Path
			path = append(path, entity.ID)

			denpendency := &model.ServiceDependency{
				ServiceID:   entity.ID,
				DependentID: d.ServiceID,
				Type:        types.ServiceDependencyTypeImplicit,
				Path:        path,
			}

			// Check if there is a dependency cycle.
			if existCycle := CheckDependencyCycle(denpendency); existCycle {
				return nil, errors.New("service dependency contains cycle")
			}
			dependencies = append(dependencies, denpendency)
		}
	}

	return dependencies, nil
}

// UpdateDependants updates the service dependants dependencies.
func UpdateDependants(ctx context.Context, mc model.ClientSet, s *model.Service) error {
	// Dependant services of the service.
	oldDependants, err := mc.ServiceDependencies().Query().
		Where(
			servicedependency.ServiceIDNEQ(s.ID),
			func(sq *sql.Selector) {
				sq.Where(sqljson.ValueContains(servicedependency.FieldPath, s.ID.String()))
			}).
		All(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	var (
		newDependantPath = sets.NewString()
		newDependants    = make([]*model.ServiceDependency, 0, len(s.Edges.Dependencies))
	)

	if len(s.Edges.Dependencies) == 0 {
		for _, sd := range oldDependants {
			path := getRightPath(s.ID, sd.Path)

			newDependants = append(newDependants, &model.ServiceDependency{
				ServiceID:   sd.ServiceID,
				DependentID: s.ID,
				Type:        types.ServiceDependencyTypeImplicit,
				Path:        path,
			})
		}
	} else {
		for _, sd := range s.Edges.Dependencies {
			for _, osd := range oldDependants {
				path := sd.Path[:len(sd.Path)-1]
				path = append(path, getRightPath(s.ID, osd.Path)...)

				newDependant := &model.ServiceDependency{
					ServiceID:   osd.ServiceID,
					DependentID: sd.ServiceID,
					Type:        types.ServiceDependencyTypeImplicit,
					Path:        path,
				}

				if existCycle := CheckDependencyCycle(newDependant); existCycle {
					return errors.New("service dependency contains cycle")
				}

				newDependants = append(newDependants, newDependant)
			}
		}
	}

	for _, d := range newDependants {
		newDependantPath.Insert(getCompleteDependencyPath(d))

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
			return err
		}
	}

	// Delete old dependants.
	for _, d := range oldDependants {
		if newDependantPath.Has(getCompleteDependencyPath(d)) {
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
					s.Where(sqljson.ValueEQ(servicedependency.FieldPath, paths))
				},
			).Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func getCompleteDependencyPath(d *model.ServiceDependency) string {
	ids := make([]string, 0, len(d.Path))
	for _, id := range d.Path {
		ids = append(ids, id.String())
	}

	return strs.Join("/", ids...)
}

// CheckDependencyCycle checks if the path contains a dependency cycle.
func CheckDependencyCycle(s *model.ServiceDependency) bool {
	idMap := make(map[oid.ID]any)
	for _, id := range s.Path {
		if _, ok := idMap[id]; ok {
			return true
		}
		idMap[id] = nil
	}

	return false
}

func getRightPath(serviceID oid.ID, path []oid.ID) []oid.ID {
	for i := 0; i < len(path)-1; i++ {
		if path[i] == serviceID {
			return path[i:]
		}
	}

	return nil
}
