package kubeendpoint

import (
	"context"
	"errors"
	"net"
	"net/url"
	"strconv"

	apicorev1 "k8s.io/api/core/v1"
	apinetworkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"

	"github.com/seal-io/seal/pkg/dao/types"
)

func GetServiceEndpoints(ctx context.Context, kubeCli *kubernetes.Clientset, ns, n string) ([]types.ApplicationResourceEndpoint, error) {
	svc, err := kubeCli.CoreV1().Services(ns).
		Get(ctx, n, metav1.GetOptions{ResourceVersion: "0"})
	if err != nil {
		return nil, err
	}

	var (
		resourceSubKind = string(svc.Spec.Type)
		endpoints       []string
	)
	switch svc.Spec.Type {
	case apicorev1.ServiceTypeNodePort:
		accessIP, err := nodeIP(ctx, kubeCli, svc)
		if err != nil {
			return nil, err
		}
		for _, port := range svc.Spec.Ports {
			nodePort := strconv.FormatInt(int64(port.NodePort), 10)
			endpoints = append(endpoints, net.JoinHostPort(accessIP, nodePort))
		}
	case apicorev1.ServiceTypeLoadBalancer:
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
	return []types.ApplicationResourceEndpoint{
		{
			EndpointType: resourceSubKind,
			Endpoints:    endpoints,
		},
	}, nil
}

func nodeIP(ctx context.Context, kubeCli *kubernetes.Clientset, svc *apicorev1.Service) (string, error) {
	list, err := kubeCli.CoreV1().Nodes().
		List(ctx, metav1.ListOptions{ResourceVersion: "0"})
	if err != nil {
		return "", err
	}

	if len(list.Items) == 0 {
		return "", errors.New("node list is empty")
	}

	var nodes = list.Items
	if svc.Spec.ExternalTrafficPolicy == apicorev1.ServiceExternalTrafficPolicyTypeLocal {
		k8sEndpoints, err := kubeCli.CoreV1().Endpoints(svc.Namespace).
			Get(ctx, svc.Name, metav1.GetOptions{ResourceVersion: "0"})
		if err != nil {
			return "", err
		}

		var nameSet = sets.Set[string]{}
		for _, v := range k8sEndpoints.Subsets {
			for _, addr := range v.Addresses {
				nameSet.Insert(*addr.NodeName)
			}
		}

		var filtered []apicorev1.Node
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

	// prefer external ip.
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

func serviceLoadBalancerIP(svc apicorev1.Service) string {
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

func GetIngressEndpoints(ctx context.Context, kubeCli *kubernetes.Clientset, ns, n string) ([]types.ApplicationResourceEndpoint, error) {
	ing, err := kubeCli.NetworkingV1().Ingresses(ns).
		Get(ctx, n, metav1.GetOptions{ResourceVersion: "0"})
	if err != nil {
		return nil, err
	}
	return []types.ApplicationResourceEndpoint{
		{
			Endpoints: ingressEndpoints(*ing),
		},
	}, nil
}

func ingressEndpoints(ing apinetworkingv1.Ingress) []string {
	var lbAddr string
	for _, ig := range ing.Status.LoadBalancer.Ingress {
		lbAddr = ig.Hostname
		if lbAddr == "" {
			lbAddr = ig.IP
		}
	}

	var tlsHostSet = sets.Set[string]{}
	for _, v := range ing.Spec.TLS {
		tlsHostSet.Insert(v.Hosts...)
	}

	var endpoints []string
	for _, v := range ing.Spec.Rules {
		var scheme = "http"
		if tlsHostSet.Has(v.Host) {
			scheme = "https"
		}

		var host = lbAddr
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
