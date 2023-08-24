package k8s

import (
	"context"
	"fmt"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/seal-io/walrus/pkg/operator/k8s/intercept"
	"github.com/seal-io/walrus/pkg/operator/k8s/kube"
	"github.com/seal-io/walrus/pkg/operator/k8s/kubelabel"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
)

// Label implements operator.Operator.
func (op Operator) Label(ctx context.Context, res *model.ServiceResource, labels map[string]string) error {
	if res == nil {
		return nil
	}

	if res.DeployerType != types.DeployerTypeTF {
		op.Logger.Warn("error resource label: unknown deployer type: " + res.DeployerType)
		return nil
	}

	gvr, ok := intercept.Terraform().GetGVR(res.Type)
	if !ok {
		return nil
	}

	ns, n := kube.ParseNamespacedName(res.Name)

	obj, err := op.DynamicCli.Resource(gvr).Namespace(ns).
		Get(ctx, n, meta.GetOptions{ResourceVersion: "0"})
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return fmt.Errorf("error getting kubernetes %s %s/%s: %w",
				gvr.Resource, ns, n, err)
		}

		return nil
	}

	return kubelabel.Apply(ctx, op.DynamicCli, obj, labels)
}
