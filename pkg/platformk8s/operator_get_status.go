package platformk8s

import (
	"context"
	"fmt"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	dynamicclient "k8s.io/client-go/dynamic"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformk8s/helm"
	"github.com/seal-io/seal/pkg/platformk8s/intercept"
	"github.com/seal-io/seal/pkg/platformk8s/kube"
	"github.com/seal-io/seal/pkg/platformk8s/kubestatus"
	"github.com/seal-io/seal/utils/strs"
)

// GetStatus implements operator.Operator.
func (op Operator) GetStatus(ctx context.Context, res *model.ApplicationResource) (string, error) {
	if res == nil {
		return kubestatus.GeneralStatusUnknown, nil
	}

	if res.DeployerType != types.DeployerTypeTF {
		op.Logger.Warn("error resource stating: unknown deployer type: " + res.DeployerType)
		return kubestatus.GeneralStatusUnknown, nil
	}

	if res.Type == "helm_release" {
		var opts = helm.GetReleaseOptions{
			RESTClientGetter: helm.IncompleteRestClientGetter(*op.RestConfig),
			Log:              op.Logger.Debugf,
		}
		var hr, err = helm.GetRelease(ctx, res, opts)
		if err != nil {
			return kubestatus.GeneralStatusUnknown,
				fmt.Errorf("error getting helm release %s, %w", res.Name, err)
		}
		if hr.Info == nil {
			return kubestatus.GeneralStatusUnready, nil
		}
		// select one from "Unknown", "Deployed", "Uninstalled",
		//   "Superseded", "Failed", "Uninstalling",
		//   "PendingInstall", "PendingUpgrade", "PendingRollback".
		return strs.Camelize(string(hr.Info.Status)), nil
	}

	var gvr, ok = intercept.Terraform().GetGVR(res.Type)
	if !ok {
		// mark ready if it's unresolved type.
		return kubestatus.GeneralStatusReady, nil
	}
	var ns, n = kube.ParseNamespacedName(res.Name)

	// fetch label selector with dynamic client.
	dynamicCli, err := dynamicclient.NewForConfig(op.RestConfig)
	if err != nil {
		return kubestatus.GeneralStatusUnknown,
			fmt.Errorf("error creating kubernetes dynamic client: %w", err)
	}
	o, err := dynamicCli.Resource(gvr).Namespace(ns).
		Get(ctx, n, meta.GetOptions{ResourceVersion: "0"}) // non quorum read
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return kubestatus.GeneralStatusUnknown, fmt.Errorf("error getting kubernetes %s %s/%s: %w",
				gvr.Resource, ns, n, err)
		}
		// mark unknown if not found.
		return kubestatus.GeneralStatusUnknown, nil
	}

	var opts = kubestatus.GetOptions{
		IgnorePaused: true,
	}
	os, err := kubestatus.Get(ctx, dynamicCli, o, opts)
	if err != nil {
		return kubestatus.GeneralStatusUnknown,
			fmt.Errorf("error stating status of kubernetes %s %s/%s: %w", gvr.Resource, ns, n, err)
	}
	return os, nil
}
