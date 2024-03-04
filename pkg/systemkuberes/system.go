package systemkuberes

import (
	"context"
	"fmt"
	"slices"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/seal-io/walrus/pkg/clients/clientset"
	"github.com/seal-io/walrus/pkg/kubeclientset"
	"github.com/seal-io/walrus/pkg/kubeclientset/review"
	"github.com/seal-io/walrus/pkg/system"
)

// SystemNamespaceName is the name indicates which Kubernetes Namespace storing system resources.
const SystemNamespaceName = "walrus-system"

// InstallSystemNamespace creates the system namespace.
func InstallSystemNamespace(ctx context.Context, cli clientset.Interface) error {
	err := review.CanDoCreate(ctx,
		cli.AuthorizationV1().SelfSubjectAccessReviews(),
		review.Simples{
			{
				Group:    core.SchemeGroupVersion.Group,
				Version:  core.SchemeGroupVersion.Version,
				Resource: "namespaces",
			},
		},
	)
	if err != nil {
		return err
	}

	nsCli := cli.CoreV1().Namespaces()
	ns := &core.Namespace{
		ObjectMeta: meta.ObjectMeta{
			Name: SystemNamespaceName,
		},
	}

	_, err = kubeclientset.Create(ctx, nsCli, ns)
	if err != nil {
		return fmt.Errorf("install namespace %q: %w", ns.GetName(), err)
	}

	return nil
}

// SystemRoutingServiceName is the name indicates which Kubernetes Service routing system access.
const SystemRoutingServiceName = "walrus"

// InstallFakeSystemRoutingService creates the fake routing service/endpoint for system.
//
// The service points to the PrimaryIP of the system.
func InstallFakeSystemRoutingService(ctx context.Context, cli clientset.Interface, port int) error {
	err := InstallSystemNamespace(ctx, cli)
	if err != nil {
		return err
	}

	err = review.CanDoCreate(ctx,
		cli.AuthorizationV1().SelfSubjectAccessReviews(),
		review.Simples{
			{
				Group:    core.SchemeGroupVersion.Group,
				Version:  core.SchemeGroupVersion.Version,
				Resource: "services",
			},
			{
				Group:    core.SchemeGroupVersion.Group,
				Version:  core.SchemeGroupVersion.Version,
				Resource: "endpoints",
			},
		},
		review.WithRecreateIfDuplicated(),
		review.WithUpdateIfExisted(),
	)
	if err != nil {
		return err
	}

	svcCli := cli.CoreV1().Services(SystemNamespaceName)
	eSvc := &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name: SystemRoutingServiceName,
			Labels: map[string]string{
				"walrus.seal.io/fake-routing": "true",
			},
		},
		Spec: core.ServiceSpec{
			Type:      core.ServiceTypeClusterIP,
			ClusterIP: core.ClusterIPNone,
			Ports: []core.ServicePort{
				{
					Port: int32(port),
				},
			},
		},
	}
	svcCompareFn := func(aSvc *core.Service) bool {
		return aSvc.Spec.Type == eSvc.Spec.Type &&
			aSvc.Spec.ClusterIP == eSvc.Spec.ClusterIP &&
			slices.ContainsFunc(aSvc.Spec.Ports, func(ap core.ServicePort) bool {
				return ap.Port == eSvc.Spec.Ports[0].Port
			})
	}

	_, err = kubeclientset.Create(ctx, svcCli, eSvc,
		kubeclientset.WithRecreateIfDuplicated(svcCompareFn))
	if err != nil {
		return fmt.Errorf("install fake rounting service %q: %w", eSvc.GetName(), err)
	}

	epCli := cli.CoreV1().Endpoints(SystemNamespaceName)
	eEp := &core.Endpoints{
		ObjectMeta: meta.ObjectMeta{
			Name: eSvc.GetName(),
			Labels: map[string]string{
				"walrus.seal.io/fake-routing": "true",
			},
		},
		Subsets: []core.EndpointSubset{
			{
				Addresses: []core.EndpointAddress{
					{
						IP: system.PrimaryIP.Get(),
					},
				},
				Ports: []core.EndpointPort{
					{
						Port: int32(port),
					},
				},
			},
		},
	}
	epAlignFn := func(aEp *core.Endpoints) (*core.Endpoints, bool, error) {
		var found bool
		for i := range aEp.Subsets {
			for j := range aEp.Subsets[i].Addresses {
				if aEp.Subsets[i].Addresses[j].IP == eEp.Subsets[0].Addresses[0].IP {
					found = true
					break
				}
			}
			if found {
				found = false
				for j := range aEp.Subsets[i].Ports {
					if aEp.Subsets[i].Ports[j].Port == eEp.Subsets[0].Ports[0].Port {
						found = true
						break
					}
				}
			}
			if found {
				break
			}
		}
		if found {
			return nil, true, nil
		}

		// Append the existing subsets.
		aEp.Subsets = append(aEp.Subsets, eEp.Subsets...)
		return eEp, false, nil
	}

	_, err = kubeclientset.Create(ctx, epCli, eEp,
		kubeclientset.WithUpdateIfExisted(epAlignFn))
	if err != nil {
		return fmt.Errorf("install fake routing service %q: %w", eEp.GetName(), err)
	}

	return nil
}
