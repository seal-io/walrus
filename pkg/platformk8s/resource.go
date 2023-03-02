package platformk8s

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformk8s/intercept"
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
	gvr schema.GroupVersionResource
	ns  string
	n   string
}

// parseResources parse the given model.ApplicationResource to resource list.
func (op Operator) parseResources(ctx context.Context, res *model.ApplicationResource) ([]resource, error) {
	// TODO(thxCode): get deployer of the application resource.
	var dt = types.DeployerTypeTF
	if dt != types.DeployerTypeTF {
		return nil, resourceParsingError("unknown deployer type: " + dt)
	}

	if res.Type == "helm_release" {
		return parseResourcesOfHelm(ctx, op, res)
	}

	var gvr, ok = intercept.Terraform().GetGVR(res.Type)
	if !ok {
		return nil, nil
	}
	ns, n, ok := parseNamespacedName(res.Name)
	if !ok {
		return nil, fmt.Errorf("failed to parse given resource name: %q", res.Name)
	}

	var rs = make([]resource, 0, 1)
	rs = append(rs, resource{
		gvr: gvr,
		ns:  ns,
		n:   n,
	})
	return rs, nil
}

// parseNamespacedName parses the given string into {namespace, name},
// returns false if not a valid namespaced name, e.g. kube-system/coredns.
func parseNamespacedName(s string) (namespace, name string, ok bool) {
	var ss = strings.SplitN(s, "/", 2)
	ok = len(ss) == 2
	if !ok {
		return
	}
	namespace = ss[0]
	name = ss[1]
	return
}
