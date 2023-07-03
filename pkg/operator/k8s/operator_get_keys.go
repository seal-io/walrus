package k8s

import (
	"context"

	"k8s.io/utils/pointer"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/operator/k8s/intercept"
	"github.com/seal-io/seal/pkg/operator/k8s/kube"
)

// GetKeys implements operator.Operator.
func (op Operator) GetKeys(
	ctx context.Context,
	res *model.ServiceResource,
) (*types.ServiceResourceOperationKeys, error) {
	if res == nil {
		return nil, nil
	}

	// Parse operable resources.
	rs, err := parseResources(ctx, op, res, intercept.Operable())
	if err != nil {
		if !isResourceParsingError(err) {
			return nil, err
		}
		// Warn out if got above errors.
		op.Logger.Warn(err)

		return nil, nil
	}

	if len(rs) == 0 {
		return nil, nil
	}

	// Get Pod of resources.
	p, err := op.getPod(ctx, rs[0].Namespace, rs[0].Name)
	if err != nil {
		return nil, err
	}

	var (
		running = kube.IsPodRunning(p)
		states  = kube.GetContainerStates(p)
	)

	ks := make([]types.ServiceResourceOperationKey, len(states))
	for i := 0; i < len(states); i++ {
		ks[i] = types.ServiceResourceOperationKey{
			Name:       states[i].Name,
			Value:      states[i].String(),
			Loggable:   pointer.Bool(states[i].State > kube.ContainerStateUnknown),
			Executable: pointer.Bool(running && states[i].State == kube.ContainerStateRunning),
		}
	}

	// {
	//      "labels": ["Container"],
	//      "keys":   [
	//          {
	//              "name": "<container name>",
	//              "value": "<key>",
	//              ...
	//          }
	//      ]
	// }.
	return &types.ServiceResourceOperationKeys{
		Labels: []string{"Container"},
		Keys:   ks,
	}, nil
}
