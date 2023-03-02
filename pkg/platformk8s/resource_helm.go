package platformk8s

import (
	"context"
	"fmt"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/releaseutil"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/platformk8s/intercept"
	"github.com/seal-io/seal/pkg/platformk8s/polymorphic"
	"github.com/seal-io/seal/utils/strs"
)

// parseResourcesOfHelm parses the given `helm_release` model.ApplicationResource to resource list.
func parseResourcesOfHelm(_ context.Context, op Operator, res *model.ApplicationResource) ([]resource, error) {
	// TODO(thxCode): there are several drivers of Helm backend,
	//  https://registry.terraform.io/providers/hashicorp/helm/latest/docs#helm_driver,
	//  get driver of the `helm` provider.
	var dr = strings.ToLower(driver.SecretsDriverName)
	if dr != strings.ToLower(driver.SecretsDriverName) &&
		dr != strings.ToLower(driver.ConfigMapsDriverName) {
		return nil, resourceParsingError("unresolved helm driver: " + dr)
	}

	// get helm release with namespace.
	var hrn = res.Name
	var hc action.Configuration
	var err = hc.Init(incompleteRestClientGetter(*op.RestConfig), "", dr, op.Logger.Debugf)
	if err != nil {
		return nil, fmt.Errorf("error initing helm config: %w", err)
	}
	var hg = action.NewGet(&hc)
	hr, err := hg.Run(hrn)
	if err != nil {
		return nil, fmt.Errorf("error getting helm release %s, %w", hrn, err)
	}

	// parse helm release.
	var rs []resource
	var hms = releaseutil.SplitManifests(hr.Manifest)
	if len(hms) == 0 {
		return nil, nil
	}
	var hs = polymorphic.YamlSerializer()
	for k := range hms {
		var (
			obj unstructured.Unstructured
			gvk *schema.GroupVersionKind
		)
		var ms = hms[k]
		_, gvk, err = hs.Decode(strs.ToBytes(&ms), nil, &obj)
		if err != nil {
			op.Logger.Warnf("error decoding helm release resource: %v", err)
			continue
		}
		// only append resource that we are interested in.
		if !intercept.Helm().AllowGVK(*gvk) {
			continue
		}
		var gvr, _ = meta.UnsafeGuessKindToResource(*gvk)
		rs = append(rs, resource{
			gvr: gvr,
			ns:  obj.GetNamespace(),
			n:   obj.GetName(),
		})
	}
	return rs, nil
}

// incompleteRestClientGetter doesn't completely implement the genericclioptions.RESTClientGetter below k8s.io/cli-runtime/pkg,
// it looks like the ToRESTConfig function is enough for kube.Client below helm.sh/helm/v3/pkg.
type incompleteRestClientGetter rest.Config

func (g incompleteRestClientGetter) ToRESTConfig() (*rest.Config, error) {
	var r = rest.Config(g)
	return &r, nil
}

func (g incompleteRestClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	panic("incomplete k8s rest client getter")
}

func (g incompleteRestClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	panic("incomplete k8s rest client getter")
}

func (g incompleteRestClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	panic("incomplete k8s rest client getter")
}
