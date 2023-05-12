package platformk8s

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformk8s/intercept"
	"github.com/seal-io/seal/pkg/platformk8s/kubeendpoint"
)

// GetEndpoints implements operator.Operator.
func (op Operator) GetEndpoints(ctx context.Context, res *model.ApplicationResource) ([]types.ApplicationResourceEndpoint, error) {
	if res == nil {
		return nil, nil
	}

	var rs, err = parseResources(ctx, op, res, intercept.Accessible())
	if err != nil {
		if !isResourceParsingError(err) {
			return nil, err
		}
		// Warn out if got above errors.
		op.Logger.Warn(err)
		return nil, nil
	}

	client, err := kubernetes.NewForConfig(op.RestConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes core client: %w", err)
	}

	var eps []types.ApplicationResourceEndpoint
	for _, r := range rs {
		switch r.Resource {
		case "services":
			var endpoints, err = kubeendpoint.GetServiceEndpoints(ctx, client, r.Namespace, r.Name)
			if err != nil {
				return nil, fmt.Errorf("error getting kubernetes service endpoints %s/%s: %w",
					r.Namespace, r.Name, err)
			}
			eps = append(eps, endpoints...)
		case "ingresses":
			var endpoints, err = kubeendpoint.GetIngressEndpoints(ctx, client, r.Namespace, r.Name)
			if err != nil {
				return nil, fmt.Errorf("error getting kubernetes ingress endpoints %s/%s: %w",
					r.Namespace, r.Name, err)
			}
			eps = append(eps, endpoints...)
		}
	}
	return eps, nil
}
