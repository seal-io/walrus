package resourcestate

import (
	"context"
	"strings"

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
		WithResource().
		All(ctx)
	if err != nil {
		return nil, err
	}

	var p parser.StateParser

	// Get the outputs of the dependency resources.
	resToOutputs := make(map[string]map[string]types.OutputValue)
	for _, s := range states {
		resToOutputs[s.Edges.Resource.Name], err = p.GetOutputMap(s.Data)
		if err != nil {
			return nil, err
		}
	}

	outputs := make(map[string]types.OutputValue)

	for k := range dependOutputs {
		split := strings.Split(k, "_")
		res, output := split[0], split[1]

		o := resToOutputs[res][output]

		outputs[k] = types.OutputValue{
			Value:     o.Value,
			Type:      o.Type,
			Sensitive: o.Sensitive,
		}
	}

	return outputs, nil
}
