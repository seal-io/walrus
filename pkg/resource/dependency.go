package resource

import (
	"context"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/terraform/parser"
)

// GetDependencyOutputs gets the dependency outputs of the resource.
func GetDependencyOutputs(
	ctx context.Context,
	client model.ClientSet,
	dependencyResourceIDs []object.ID,
	dependOutputs map[string]string,
) (map[string]types.OutputValue, error) {
	dependencyResources, err := client.Resources().Query().
		Where(resource.IDIn(dependencyResourceIDs...)).
		WithState().
		All(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make(map[string]types.OutputValue)

	var p parser.StateParser

	for _, r := range dependencyResources {
		osm, err := p.GetOutputMap(r.Edges.State.Data)
		if err != nil {
			return nil, err
		}

		for n, o := range osm {
			if _, ok := dependOutputs[n]; !ok {
				continue
			}

			outputs[n] = types.OutputValue{
				Value:     o.Value,
				Type:      o.Type,
				Sensitive: o.Sensitive,
			}
		}
	}

	return outputs, nil
}
