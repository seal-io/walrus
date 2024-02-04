package k8s

import (
	"context"
	"fmt"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/operator/k8s/helm"
	"github.com/seal-io/walrus/pkg/operator/k8s/intercept"
	"github.com/seal-io/walrus/pkg/operator/k8s/kube"
	"github.com/seal-io/walrus/pkg/operator/k8s/kubestatus"
)

// GetStatus implements operator.Operator.
func (op Operator) GetStatus(ctx context.Context, res *model.ResourceComponent) (*status.Status, error) {
	if res == nil {
		return kubestatus.StatusError(""), nil
	}

	if res.DeployerType != types.DeployerTypeTF {
		op.Logger.Warn("error resource stating: unknown deployer type: " + res.DeployerType)
		return kubestatus.StatusError("unknown deployer type"), nil
	}

	if res.Type == "helm_release" {
		opts := helm.GetReleaseOptions{
			RESTClientGetter: helm.IncompleteRestClientGetter(*op.RestConfig),
			Log:              op.Logger.Debugf,
		}

		return helm.GetReleaseStatus(ctx, res, opts)
	}

	gvr, ok := intercept.Terraform().GetGVR(res.Type)
	if !ok {
		// Mark ready if it's unresolved type.
		return &kubestatus.GeneralStatusReady, nil
	}
	ns, n := kube.ParseNamespacedName(res.Name)

	// Fetch label selector with dynamic client.
	o, err := op.DynamicCli.Resource(gvr).Namespace(ns).
		Get(ctx, n, meta.GetOptions{ResourceVersion: "0"}) // Non quorum read.
	if err != nil {
		if !kerrors.IsNotFound(err) {
			err = fmt.Errorf("error getting kubernetes %s %s/%s: %w", gvr.Resource, ns, n, err)
			return kubestatus.StatusError(err.Error()), err
		}
		// Mark unknown if not found.
		return kubestatus.StatusError("resource not found"), nil
	}

	os, err := kubestatus.Get(ctx, op.DynamicCli, o)
	if err != nil {
		err = fmt.Errorf("error stating status of kubernetes %s %s/%s: %w", gvr.Resource, ns, n, err)
		return kubestatus.StatusError(err.Error()), err
	}

	return os, nil
}
