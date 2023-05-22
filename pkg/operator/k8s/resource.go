package k8s

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/operator/k8s/intercept"
	"github.com/seal-io/seal/pkg/operator/k8s/kube"
)

// resourceParsingError emits if the given model.ApplicationResource got deployer error.
type resourceParsingError string

func (e resourceParsingError) Error() string {
	return "error resource parsing: " + string(e)
}

// isResourceParsingError returns true if the given error is resourceParsingError.
func isResourceParsingError(err error) bool {
	var e resourceParsingError
	return errors.As(err, &e)
}

// resource holds the GVR, namespace and name of a Kubernetes resource.
type resource struct {
	schema.GroupVersionResource

	Namespace string
	Name      string
}

// parseOperableResources parse the given model.ApplicationResource,
// and keeps resource item which matches enforcer validation.
func parseResources(
	ctx context.Context,
	op Operator,
	res *model.ApplicationResource,
	enforcer intercept.Enforcer,
) ([]resource, error) {
	if res.DeployerType != types.DeployerTypeTF {
		return nil, resourceParsingError("unknown deployer type: " + res.DeployerType)
	}

	if res.Type == "helm_release" {
		return parseResourcesOfHelm(ctx, op, res, enforcer.AllowGVK)
	}

	gvr, ok := intercept.Terraform().GetGVR(res.Type)
	if !ok {
		return nil, nil
	}

	if !enforcer.AllowGVR(gvr) {
		return nil, nil
	}
	ns, n := kube.ParseNamespacedName(res.Name)

	rs := make([]resource, 0, 1)
	rs = append(rs, resource{
		GroupVersionResource: gvr,
		Namespace:            ns,
		Name:                 n,
	})

	return rs, nil
}
