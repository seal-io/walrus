package apis

import (
	"context"
	"fmt"

	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	apireg "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	walrus "github.com/seal-io/walrus/pkg/apis/walrus/v1"
	walruscore "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
	"github.com/seal-io/walrus/pkg/clients/clientset"
	"github.com/seal-io/walrus/pkg/kubeclientset"
	"github.com/seal-io/walrus/pkg/kubeclientset/review"
)

// NB(thxCode): Register APIs below.
var (
	crdGetters = []_CustomResourceDefinitionsGetter{
		walruscore.GetCustomResourceDefinitions,
	}
	apiSvcGetters = []_APIServiceGetter{
		walrus.GetAPIService,
	}
)

type _CustomResourceDefinitionsGetter func() map[string]*apiext.CustomResourceDefinition

// GetCustomResourceDefinitions returns the registered custom resource definitions.
func GetCustomResourceDefinitions() []*apiext.CustomResourceDefinition {
	// Merge all the CRDs from the getters.
	var (
		ret = make([]map[string]*apiext.CustomResourceDefinition, len(crdGetters))
		csc int
	)
	for i, get := range crdGetters {
		ret[i] = get()
		csc += len(ret[i])
	}

	crds := make([]*apiext.CustomResourceDefinition, 0, csc)
	for i := range ret {
		if ret[i] == nil {
			continue
		}
		for _, n := range sets.List(sets.KeySet(ret[i])) {
			crds = append(crds, ret[i][n])
		}
	}

	return crds
}

// InstallCustomResourceDefinitions installs the custom resource definitions.
func InstallCustomResourceDefinitions(ctx context.Context, cli clientset.Interface) error {
	err := review.CanDoUpdate(ctx,
		cli.AuthorizationV1().SelfSubjectAccessReviews(),
		review.Simples{
			{
				Group:    apiext.SchemeGroupVersion.Group,
				Version:  apiext.SchemeGroupVersion.Version,
				Resource: "customresourcedefinitions",
			},
		},
		review.WithCreate(),
	)
	if err != nil {
		return err
	}

	crdCli := cli.ApiextensionsV1().CustomResourceDefinitions()

	crds := GetCustomResourceDefinitions()
	for i := range crds {
		_, err = kubeclientset.Update(ctx, crdCli, crds[i],
			kubeclientset.UpdateOrCreate[*apiext.CustomResourceDefinition]())
		if err != nil {
			return fmt.Errorf("install custom resource definition %q: %w",
				crds[i].GetName(), err)
		}
	}

	return nil
}

type _APIServiceGetter func(apireg.ServiceReference, []byte) *apireg.APIService

// GetAPIServices returns the registered api services.
func GetAPIServices(svc apireg.ServiceReference, ca []byte) []*apireg.APIService {
	ret := make([]*apireg.APIService, 0, len(apiSvcGetters))
	for i := range apiSvcGetters {
		r := apiSvcGetters[i](svc, ca)
		if r != nil {
			ret = append(ret, r)
		}
	}
	return ret
}

// InstallAPIServices installs the api services.
func InstallAPIServices(ctx context.Context, cli clientset.Interface, svc apireg.ServiceReference, ca []byte) error {
	err := review.CanDoUpdate(ctx,
		cli.AuthorizationV1().SelfSubjectAccessReviews(),
		review.Simples{
			{
				Group:    apireg.SchemeGroupVersion.Group,
				Version:  apireg.SchemeGroupVersion.Version,
				Resource: "apiservices",
			},
		},
		review.WithCreate(),
	)
	if err != nil {
		return err
	}

	svcCli := cli.ApiregistrationV1().APIServices()

	svcs := GetAPIServices(svc, ca)
	for i := range svcs {
		_, err = kubeclientset.Update(ctx, svcCli, svcs[i],
			kubeclientset.UpdateOrCreate[*apireg.APIService]())
		if err != nil {
			return fmt.Errorf("install api service %q: %w",
				svcs[i].Name, err)
		}
	}

	return nil
}
