package resource

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerun"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/terraform/parser"
)

// GetDependencyOutputs gets the dependency outputs of the resource.
func GetDependencyOutputs(
	ctx context.Context,
	client model.ClientSet,
	dependencyResourceIDs []object.ID,
	dependOutputs map[string]string,
) (map[string]parser.OutputState, error) {
	dependencyRuns, err := client.ResourceRuns().Query().
		Select(
			resourcerun.FieldID,
			resourcerun.FieldAttributes,
			resourcerun.FieldOutput,
			resourcerun.FieldResourceID,
			resourcerun.FieldProjectID,
		).
		Where(func(s *sql.Selector) {
			sq := s.Clone().
				AppendSelectExprAs(
					sql.RowNumber().
						PartitionBy(resourcerun.FieldResourceID).
						OrderBy(sql.Desc(resourcerun.FieldCreateTime)),
					"row_number",
				).
				Where(s.P()).
				From(s.Table()).
				As(resourcerun.Table)

			// Query the latest run of the resource.
			s.Where(sql.EQ(s.C("row_number"), 1)).
				From(sq)
		}).
		Where(resourcerun.ResourceIDIn(dependencyResourceIDs...)).
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

	outputs := make(map[string]parser.OutputState)

	var p parser.ResourceRunParser

	for _, r := range dependencyRuns {
		osm, err := p.GetOutputsMap(r)
		if err != nil {
			return nil, err
		}

		for n, o := range osm {
			if _, ok := dependOutputs[n]; !ok {
				continue
			}

			// FIXME(thxCode): migrate parser.OutputState to types.OutputValue.
			outputs[n] = parser.OutputState{
				Value:     o.Value,
				Type:      o.Type,
				Sensitive: o.Sensitive,
			}
		}
	}

	return outputs, nil
}
