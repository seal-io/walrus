package kube

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"

	apicorev1 "k8s.io/api/core/v1"
	apinetworkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
)

type ResourceEndpoint struct {
	// ResourceID is the namespaced name
	ResourceID string `json:"resourceID,omitempty"`
	// ResourceKind be Ingress or Service
	ResourceKind string `json:"resourceKind,omitempty"`
	// ResourceSubKind is the sub kind for endpoint, like nodePort, loadBalance
	ResourceSubKind string `json:"resourceSubKind,omitempty"`
	// Endpoints is access endpoints
	Endpoints []string `json:"endpoints,omitempty"`
}

type EndpointGetter interface {
	Endpoints(ctx context.Context, resourceID ...string) ([]ResourceEndpoint, error)
}

func ServiceEndpointGetter(clientSet *kubernetes.Clientset) EndpointGetter {
	return &serviceEndpointGetter{
		clientSet: clientSet,
	}
}

type serviceEndpointGetter struct {
	clientSet *kubernetes.Clientset
}

func (s *serviceEndpointGetter) Endpoints(ctx context.Context, resourceIDs ...string) ([]ResourceEndpoint, error) {
	var eps []ResourceEndpoint
	for _, v := range resourceIDs {
		ep, err := s.endpoint(ctx, v)
		if err != nil {
			return nil, err
		}

		if ep != nil {
			eps = append(eps, *ep)
		}
	}
	return eps, nil
}

func (s *serviceEndpointGetter) endpoint(ctx context.Context, resourceID string) (*ResourceEndpoint, error) {
	var rn = strings.SplitN(resourceID, "/", 2)
	if len(rn) != 2 {
		return nil, fmt.Errorf("invalid service namespaced name: %s", rn)
	}

	var (
		ns   = rn[0]
		name = rn[1]
	)
	svc, err := s.clientSet.CoreV1().Services(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var (
		resourceKind    = "Service"
		resourceSubKind = string(svc.Spec.Type)
		endpoints       []string
	)
	switch svc.Spec.Type {
	case apicorev1.ServiceTypeNodePort:
		nodes, err := s.clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		if len(nodes.Items) == 0 {
			return nil, fmt.Errorf("node list is empty")
		}

		accessIP := nodeIP(nodes.Items)
		for _, port := range svc.Spec.Ports {
			nodePort := fmt.Sprint(port.NodePort)
			endpoints = append(endpoints, net.JoinHostPort(accessIP, nodePort))
		}
	case apicorev1.ServiceTypeLoadBalancer:
		accessIP := serviceLoadBalancerIP(*svc)
		if accessIP != "" {
			endpoints = append(endpoints, accessIP)
		}
	}

	if len(endpoints) == 0 {
		return nil, nil
	}
	return &ResourceEndpoint{
		ResourceID:      resourceID,
		ResourceKind:    resourceKind,
		ResourceSubKind: resourceSubKind,
		Endpoints:       endpoints,
	}, nil
}

func nodeIP(nodes []apicorev1.Node) string {
	var (
		externalIP string
		internalIP string
	)

	// prefer external ip
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
			return externalIP
		}
	}

	return internalIP
}

func serviceLoadBalancerIP(svc apicorev1.Service) string {
	var lbIP string
	for _, ing := range svc.Status.LoadBalancer.Ingress {
		lbIP = ing.IP
	}

	if lbIP != "" {
		return lbIP
	}
	return svc.Spec.LoadBalancerIP
}

func IngressEndpointGetter(clientSet *kubernetes.Clientset) EndpointGetter {
	return &ingressEndpointGetter{
		clientSet: clientSet,
	}
}

type ingressEndpointGetter struct {
	clientSet *kubernetes.Clientset
}

func (ig *ingressEndpointGetter) Endpoints(ctx context.Context, resourceIDs ...string) ([]ResourceEndpoint, error) {
	var eps []ResourceEndpoint
	for _, v := range resourceIDs {
		ep, err := ig.endpoint(ctx, v)
		if err != nil {
			return nil, err
		}

		if ep != nil {
			eps = append(eps, *ep)
		}
	}
	return eps, nil
}

func (ig *ingressEndpointGetter) endpoint(ctx context.Context, resourceID string) (*ResourceEndpoint, error) {
	var rn = strings.SplitN(resourceID, "/", 2)
	if len(rn) != 2 {
		return nil, fmt.Errorf("invalid ingress namespaced name: %s", rn)
	}

	var (
		ns   = rn[0]
		name = rn[1]
	)
	ing, err := ig.clientSet.NetworkingV1().Ingresses(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	endpoints := ingressEndpoints(*ing)
	if len(endpoints) == 0 {
		return nil, nil
	}
	return &ResourceEndpoint{
		ResourceID:   resourceID,
		ResourceKind: "Ingress",
		Endpoints:    endpoints,
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
