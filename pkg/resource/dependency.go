package resource

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerelationship"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/terraform/parser"
)

// GetServiceDependencyOutputsByID gets the dependency outputs of the resource by resource id.
func GetServiceDependencyOutputsByID(
	ctx context.Context,
	client model.ClientSet,
	resourceID object.ID,
	dependOutputs map[string]string,
) (map[string]parser.OutputState, error) {
	entity, err := client.Resources().Query().
		Where(resource.ID(resourceID)).
		WithDependencies(func(sq *model.ResourceRelationshipQuery) {
			sq.Where(func(s *sql.Selector) {
				s.Where(sql.ColumnsNEQ(resourcerelationship.FieldResourceID, resourcerelationship.FieldDependencyID))
			})
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	dependencyResourceIDs := make([]object.ID, 0, len(entity.Edges.Dependencies))

	for _, d := range entity.Edges.Dependencies {
		if d.Type != types.ResourceRelationshipTypeImplicit {
			continue
		}

		dependencyResourceIDs = append(dependencyResourceIDs, d.DependencyID)
	}

	return GetResourceDependencyOutputs(ctx, client, dependencyResourceIDs, dependOutputs)
}

// GetResourceDependencyOutputs gets the dependency outputs of the resource.
func GetResourceDependencyOutputs(
	ctx context.Context,
	client model.ClientSet,
	dependencyResourceIDs []object.ID,
	dependOutputs map[string]string,
) (map[string]parser.OutputState, error) {
	dependencyRevisions, err := client.ResourceRevisions().Query().
		Select(
			resourcerevision.FieldID,
			resourcerevision.FieldAttributes,
			resourcerevision.FieldOutput,
			resourcerevision.FieldResourceID,
			resourcerevision.FieldProjectID,
		).
		Where(func(s *sql.Selector) {
			sq := s.Clone().
				AppendSelectExprAs(
					sql.RowNumber().
						PartitionBy(resourcerevision.FieldResourceID).
						OrderBy(sql.Desc(resourcerevision.FieldCreateTime)),
					"row_number",
				).
				Where(s.P()).
				From(s.Table()).
				As(resourcerevision.Table)

			// Query the latest revision of the resource.
			s.Where(sql.EQ(s.C("row_number"), 1)).
				From(sq)
		}).
		Where(resourcerevision.ResourceIDIn(dependencyResourceIDs...)).
		WithResource(func(rq *model.ResourceQuery) {
			rq.Select(
				resource.FieldTemplateID,
				resource.FieldResourceDefinitionID,
			)
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make(map[string]parser.OutputState, 0)

	for _, r := range dependencyRevisions {
		revisionOutput, err := parser.ParseStateOutputRawMap(r)
		if err != nil {
			return nil, err
		}

		for n, o := range revisionOutput {
			outputRefKind, ok := dependOutputs[n]
			if !ok {
				continue
			}

			if outputRefKind == "service" && !IsService(r.Edges.Resource) {
				// An output reference in format ${service.name.output_name} only applies to resource of service type.
				continue
			}

			outputs[n] = o
		}
	}

	return outputs, nil
}
