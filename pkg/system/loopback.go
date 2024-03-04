package system

import (
	"github.com/seal-io/utils/varx"
	"k8s.io/client-go/rest"
	ctrlcache "sigs.k8s.io/controller-runtime/pkg/cache"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/seal-io/walrus/pkg/clients/clientset"
)

var (
	// LoopbackKubeInside is a flag that indicates whether the system runs inside the loopback Kubernetes cluster.
	LoopbackKubeInside varx.Once[bool]

	// LoopbackKubeNearby is a flag that indicates whether the system runs nearby the loopback Kubernetes cluster.
	// If the system runs nearby, it can connect to the loopback Kubernetes cluster even if it is not inside the cluster.
	LoopbackKubeNearby varx.Once[bool]

	// LoopbackKubeClientConfigPath is the path to the loopback Kubernetes client configuration file.
	LoopbackKubeClientConfigPath varx.Once[string]

	// LoopbackKubeClientConfig is the loopback Kubernetes client configuration.
	LoopbackKubeClientConfig varx.Once[rest.Config]

	// LoopbackKubeClient is the loopback Kubernetes client.
	LoopbackKubeClient varx.Once[clientset.Interface]
)

// ConfigureLoopbackKube configures the loopback Kubernetes.
func ConfigureLoopbackKube(inside, nearby bool, configPath string, config rest.Config, client clientset.Interface) {
	LoopbackKubeInside.Configure(inside)
	LoopbackKubeNearby.Configure(nearby)
	LoopbackKubeClientConfigPath.Configure(configPath)
	LoopbackKubeClientConfig.Configure(config)
	LoopbackKubeClient.Configure(client)
}

var (
	// LoopbackCtrlClient is the controller client for the loopback Kubernetes cluster.
	//
	// LoopbackCtrlClient is similar to LoopbackKubeClient,
	// but it has a self-manager cache,
	// which means we don't need to handle list/watch manually.
	LoopbackCtrlClient varx.Once[ctrlcli.Client]

	// LoopbackCtrlCache is the controller cache for the loopback Kubernetes cluster.
	LoopbackCtrlCache varx.Once[ctrlcache.Cache]
)

// ConfigureLoopbackCtrlRuntime configures the loopback Kubernetes controller runtime.
func ConfigureLoopbackCtrlRuntime(client ctrlcli.Client, cache ctrlcache.Cache) {
	LoopbackCtrlClient.Configure(client)
	LoopbackCtrlCache.Configure(cache)
}
