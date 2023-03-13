package platformk8s

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformk8s/intercept"
	"github.com/seal-io/seal/pkg/platformk8s/kube"
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

// parseOperableResources parse the given model.ApplicationResource to operable resource list.
func parseOperableResources(ctx context.Context, op Operator, res *model.ApplicationResource) ([]resource, error) {
	if res.DeployerType != types.DeployerTypeTF {
		return nil, resourceParsingError("unknown deployer type: " + res.DeployerType)
	}

	if res.Type == "helm_release" {
		return parseOperableResourcesOfHelm(ctx, op, res)
	}

	var gvr, ok = intercept.Terraform().GetGVR(res.Type)
	if !ok {
		return nil, nil
	}
	// only allow operable resource.
	if !intercept.Operable().AllowGVR(gvr) {
		return nil, nil
	}
	var ns, n = kube.ParseNamespacedName(res.Name)

	var rs = make([]resource, 0, 1)
	rs = append(rs, resource{
		GroupVersionResource: gvr,
		Namespace:            ns,
		Name:                 n,
	})
	return rs, nil
}
