package k8s

import (
	"context"
	"fmt"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/operator/k8s/intercept"
	"github.com/seal-io/walrus/pkg/operator/k8s/kubeendpoint"
)

// GetEndpoints implements operator.Operator.
func (op Operator) GetEndpoints(
	ctx context.Context,
	res *model.ResourceComponent,
) ([]types.ResourceComponentEndpoint, error) {
	if res == nil {
		return nil, nil
	}

	rs, err := parseResources(ctx, op, res, intercept.Accessible())
	if err != nil {
		if !isResourceParsingError(err) {
			return nil, err
		}
		// Warn out if got above errors.
		op.Logger.Warn(err)

		return nil, nil
	}

	var eps []types.ResourceComponentEndpoint

	for _, r := range rs {
		switch r.Resource {
		case "services":
			svc, err := op.CoreCli.Services(r.Namespace).
				Get(ctx, r.Name, meta.GetOptions{ResourceVersion: "0"})
			if err != nil {
				return nil, fmt.Errorf("error getting kubernetes service endpoints %s/%s: %w",
					r.Namespace, r.Name, err)
			}

			accessHostname := op.ServeHostname
			if !op.IsEmbedded {
				accessHostname, err = kubeendpoint.GetNodeIP(ctx, op.CoreCli, svc)
				if err != nil {
					return nil, err
				}
			}

			endpoints := kubeendpoint.GetServiceEndpoints(
				accessHostname,
				svc,
			)

			eps = append(eps, endpoints...)
		case "ingresses":
			endpoints, err := kubeendpoint.GetIngressEndpoints(
				ctx,
				op.NetworkingCli,
				r.Namespace,
				r.Name,
			)
			if err != nil {
				return nil, fmt.Errorf("error getting kubernetes ingress endpoints %s/%s: %w",
					r.Namespace, r.Name, err)
			}

			eps = append(eps, endpoints...)
		}
	}

	return eps, nil
}
