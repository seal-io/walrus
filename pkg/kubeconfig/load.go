package kubeconfig

import (
	"errors"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// LoadRestConfigNonInteractive loads a rest config according to the following rules.
//
//  1. assume that running as a Pod and try to connect to
//     the Kubernetes cluster with the mounted ServiceAccount.
//  2. load from recommended home file if none of the above conditions are met.
func LoadRestConfigNonInteractive() (cfgPath string, restCfg *rest.Config, inside bool, err error) {
	// Try the in-cluster config.
	restCfg, err = rest.InClusterConfig()
	switch {
	case err == nil:
		return "", restCfg, true, nil
	case err != nil && !errors.Is(err, rest.ErrNotInCluster):
		return "", nil, false, err
	}

	// Try the recommended config.
	var (
		ld = &clientcmd.ClientConfigLoadingRules{
			Precedence: []string{clientcmd.RecommendedHomeFile},
		}
		od = &clientcmd.ConfigOverrides{}
	)
	restCfg, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(ld, od).ClientConfig()
	return clientcmd.RecommendedHomeFile, restCfg, false, err
}

// LoadClientConfig loads a client config from the specified path,
// the given path must exist.
func LoadClientConfig(path string) (clientcmd.ClientConfig, error) {
	if path == "" {
		return nil, errors.New("blank kubeconfig path")
	}

	var (
		ld = &clientcmd.ClientConfigLoadingRules{
			ExplicitPath: path,
		}
		od = &clientcmd.ConfigOverrides{}
	)

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(ld, od), nil
}

// LoadRestConfig loads a rest config from the specified path,
// the given path must exist.
func LoadRestConfig(path string) (*rest.Config, error) {
	cc, err := LoadClientConfig(path)
	if err != nil {
		return nil, err
	}

	return cc.ClientConfig()
}
