package k8s

import (
	"context"
	"sort"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/operator/k8s/intercept"
	"github.com/seal-io/walrus/pkg/operator/k8s/kube"
	"github.com/seal-io/walrus/utils/pointer"
)

// GetKeys implements operator.Operator.
func (op Operator) GetKeys(
	ctx context.Context,
	res *model.ResourceComponent,
) (*types.ResourceComponentOperationKeys, error) {
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

	sort.Slice(states, func(i, j int) bool {
		ti, tj := states[i].Type, states[j].Type

		return kube.ContainerTypeOrderMap[ti] < kube.ContainerTypeOrderMap[tj]
	})

	ks := make([]types.ResourceComponentOperationKey, len(states))
	for i := 0; i < len(states); i++ {
		ks[i] = types.ResourceComponentOperationKey{
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
	return &types.ResourceComponentOperationKeys{
		Labels: []string{"Container"},
		Keys:   ks,
	}, nil
}
