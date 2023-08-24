package kubeendpoint

import (
	"context"
	"errors"
	"net"
	"net/url"
	"strconv"

	core "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	networkingclient "k8s.io/client-go/kubernetes/typed/networking/v1"

	"github.com/seal-io/walrus/pkg/dao/types"
)

func GetServiceEndpoints(
	ctx context.Context,
	coreCli *coreclient.CoreV1Client,
	ns, n string,
) ([]types.ServiceResourceEndpoint, error) {
	svc, err := coreCli.Services(ns).
		Get(ctx, n, meta.GetOptions{ResourceVersion: "0"})
	if err != nil {
		return nil, err
	}

	var (
		resourceSubKind = string(svc.Spec.Type)
		endpoints       []string
	)

	switch svc.Spec.Type {
	case core.ServiceTypeNodePort:
		accessIP, err := nodeIP(ctx, coreCli, svc)
		if err != nil {
			return nil, err
		}

		for _, port := range svc.Spec.Ports {
			nodePort := strconv.FormatInt(int64(port.NodePort), 10)
			endpoints = append(endpoints, net.JoinHostPort(accessIP, nodePort))
		}
	case core.ServiceTypeLoadBalancer:
		accessIP := serviceLoadBalancerIP(*svc)
		if accessIP != "" {
			for _, port := range svc.Spec.Ports {
				targetPort := strconv.FormatInt(int64(port.Port), 10)
				endpoints = append(endpoints, net.JoinHostPort(accessIP, targetPort))
			}
		}
	}

	if len(endpoints) == 0 {
		return nil, nil
	}

	return []types.ServiceResourceEndpoint{
		{
			EndpointType: resourceSubKind,
			Endpoints:    endpoints,
		},
	}, nil
}

func nodeIP(ctx context.Context, coreCli *coreclient.CoreV1Client, svc *core.Service) (string, error) {
	list, err := coreCli.Nodes().
		List(ctx, meta.ListOptions{ResourceVersion: "0"})
	if err != nil {
		return "", err
	}

	if len(list.Items) == 0 {
		return "", errors.New("node list is empty")
	}

	nodes := list.Items

	if svc.Spec.ExternalTrafficPolicy == core.ServiceExternalTrafficPolicyTypeLocal {
		k8sEndpoints, err := coreCli.Endpoints(svc.Namespace).
			Get(ctx, svc.Name, meta.GetOptions{ResourceVersion: "0"})
		if err != nil {
			return "", err
		}

		nameSet := sets.Set[string]{}

		for _, v := range k8sEndpoints.Subsets {
			for _, addr := range v.Addresses {
				nameSet.Insert(*addr.NodeName)
			}
		}

		var filtered []core.Node

		for _, node := range nodes {
			if nameSet.Has(node.Name) {
				filtered = append(filtered, node)
			}
		}

		if len(filtered) == 0 {
			return "", errors.New("node list from k8s endpoints is empty")
		}
		nodes = filtered
	}

	var (
		externalIP string
		internalIP string
	)

	// Prefer external ip.
	for _, node := range nodes {
		for _, ip := range node.Status.Addresses {
			if ip.Type == "ExternalIP" && ip.Address != "" {
				externalIP = ip.Address
				break
			} else if ip.Type == "InternalIP" && ip.Address != "" {
				internalIP = ip.Address
			}
		}

		if externalIP != "" {
			return externalIP, nil
		}
	}

	return internalIP, nil
}

func serviceLoadBalancerIP(svc core.Service) string {
	for _, ing := range svc.Status.LoadBalancer.Ingress {
		if ing.Hostname != "" {
			return ing.Hostname
		}

		if ing.IP != "" {
			return ing.IP
		}
	}

	return svc.Spec.LoadBalancerIP
}

func GetIngressEndpoints(
	ctx context.Context,
	networkCli *networkingclient.NetworkingV1Client,
	ns, n string,
) ([]types.ServiceResourceEndpoint, error) {
	ing, err := networkCli.Ingresses(ns).
		Get(ctx, n, meta.GetOptions{ResourceVersion: "0"})
	if err != nil {
		return nil, err
	}

	return []types.ServiceResourceEndpoint{
		{
			Endpoints: ingressEndpoints(*ing),
		},
	}, nil
}

func ingressEndpoints(ing networking.Ingress) []string {
	var lbAddr string
	for _, ig := range ing.Status.LoadBalancer.Ingress {
		lbAddr = ig.Hostname
		if lbAddr == "" {
			lbAddr = ig.IP
		}
	}

	tlsHostSet := sets.Set[string]{}
	for _, v := range ing.Spec.TLS {
		tlsHostSet.Insert(v.Hosts...)
	}

	var endpoints []string

	for _, v := range ing.Spec.Rules {
		scheme := "http"
		if tlsHostSet.Has(v.Host) {
			scheme = "https"
		}

		host := lbAddr
		if v.Host != "" {
			host = v.Host
		}

		if host == "" {
			continue
		}

		if v.HTTP != nil {
			for _, httpPath := range v.HTTP.Paths {
				ep := url.URL{
					Host:   host,
					Path:   httpPath.Path,
					Scheme: scheme,
				}
				endpoints = append(endpoints, ep.String())
			}
		}
	}

	return endpoints
}
