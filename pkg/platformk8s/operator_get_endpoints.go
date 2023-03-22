package platformk8s

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformk8s/kube"
)

// GetEndpoints implements operator.Operator.
func (op Operator) GetEndpoints(ctx context.Context, res *model.ApplicationResource) ([]types.ApplicationResourceEndpoint, error) {
	if res == nil || res.DeployerType != types.DeployerTypeTF {
		return nil, nil
	}

	client, err := kubernetes.NewForConfig(op.RestConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes core client: %w", err)
	}

	return kube.GetEndpoints(ctx, client, res.Type, res.Name)
}
