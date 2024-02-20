package resourcestate

import (
	"context"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resourcestate"
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
	states, err := client.ResourceStates().Query().
		Where(resourcestate.ResourceIDIn(dependencyResourceIDs...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make(map[string]types.OutputValue)

	var p parser.StateParser

	for _, s := range states {
		osm, err := p.GetOutputMap(s.Data)
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
