package helm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platformk8s/kube"
	"github.com/seal-io/seal/pkg/platformk8s/kubestatus"
	"github.com/seal-io/seal/utils/strs"
)

type (
	Release           = release.Release
	GetReleaseOptions = action.Configuration
)

func GetRelease(ctx context.Context, res *model.ApplicationResource, opts GetReleaseOptions) (*Release, error) {
	// TODO(thxCode): there are several drivers of Operable backend,
	//  https://registry.terraform.io/providers/hashicorp/helm/latest/docs#helm_driver,
	//  get driver of the `helm` provider.
	var dr = strings.ToLower(driver.SecretsDriverName)
	if dr != strings.ToLower(driver.SecretsDriverName) &&
		dr != strings.ToLower(driver.ConfigMapsDriverName) {
		return nil, errors.New("unresolved helm driver: " + dr)
	}

	// get helm release with namespace.
	var hrns, hrn = kube.ParseNamespacedName(res.Name)
	var restConfig, err = opts.RESTClientGetter.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("error getting helm rest config: %w", err)
	}
	if opts.KubeClient == nil || opts.Releases == nil {
		var clientGetter = IncompleteRestClientGetter(*restConfig)
		err = opts.Init(clientGetter, hrns, dr, opts.Log)
		if err != nil {
			return nil, fmt.Errorf("error initing helm config: %w", err)
		}
	}
	var hg = action.NewGet(&opts)
	hr, err := hg.Run(hrn)
	if err != nil {
		return nil, fmt.Errorf("error getting helm release %s, %w", hrn, err)
	}

	return hr, nil
}

func GetReleaseStatus(ctx context.Context, res *model.ApplicationResource, opts GetReleaseOptions) (*status.Status, error) {
	var hr, err = GetRelease(ctx, res, opts)
	if err != nil {
		return kubestatus.StatusError(err.Error()), err
	}
	if hr.Info == nil {
		return &kubestatus.GeneralStatusReadyTransitioning, nil
	}

	var isErr, isTransitioning bool
	switch hr.Info.Status {
	case release.StatusFailed:
		isErr = true
	case release.StatusUnknown,
		release.StatusUninstalling,
		release.StatusPendingInstall,
		release.StatusPendingUpgrade,
		release.StatusPendingRollback:
		isTransitioning = true
	default:
		// release.StatusDeployed,
		// release.StatusUninstalled,
		// release.StatusSuperseded
	}
	return &status.Status{
		Summary: status.Summary{
			SummaryStatus: strs.Camelize(string(hr.Info.Status)),
			Error:         isErr,
			Transitioning: isTransitioning,
		},
	}, nil
}

// IncompleteRestClientGetter doesn't completely implement the genericclioptions.RESTClientGetter below k8s.io/cli-runtime/pkg,
// it looks like the ToRESTConfig function is enough for kube.Client below helm.sh/helm/v3/pkg.
type IncompleteRestClientGetter rest.Config

func (g IncompleteRestClientGetter) ToRESTConfig() (*rest.Config, error) {
	var r = rest.Config(g)
	return &r, nil
}

func (g IncompleteRestClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	panic("incomplete k8s rest client getter")
}

func (g IncompleteRestClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	panic("incomplete k8s rest client getter")
}

func (g IncompleteRestClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	panic("incomplete k8s rest client getter")
}
