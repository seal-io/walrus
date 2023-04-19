package platformk8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
)

func TestGetConfig(t *testing.T) {
	t.Run("v1", func(t *testing.T) {
		var dummyKubeconfigText = `
apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: ""
    server: "https://127.0.0.1:6443"
  name: dummy
contexts:
- context:
    cluster: dummy
    user: dummy
  name: dummy
current-context: dummy
users:
- name: dummy
  user:
    client-certificate-data: ""
    client-key-data: ""
`
		var _, err = GetConfig(model.Connector{
			Name:          "test",
			ConfigVersion: "v1",
			ConfigData: crypto.Properties{
				"kubeconfig_": crypto.StringProperty(dummyKubeconfigText),
			},
		})
		assert.EqualError(t, err, "error load config from connector test: not found \"kubeconfig\"")

		config, err := GetConfig(model.Connector{
			ConfigVersion: "v1",
			ConfigData: crypto.Properties{
				"kubeconfig": crypto.StringProperty(dummyKubeconfigText),
			},
		})
		if assert.NoError(t, err, "unexpected error") {
			assert.Equal(t, &rest.Config{
				Host: "https://127.0.0.1:6443",
			}, config)
		}

		apiConfig, _, err := LoadApiConfig(model.Connector{
			ConfigVersion: "v1",
			ConfigData: crypto.Properties{
				"kubeconfig": crypto.StringProperty(dummyKubeconfigText),
			},
		})
		if assert.NoError(t, err, "unexpected error") {
			assert.Equal(t, &api.Config{
				Clusters: map[string]*api.Cluster{
					"dummy": {
						CertificateAuthorityData: []byte{},
						Server:                   "https://127.0.0.1:6443",
						Extensions:               map[string]runtime.Object{},
					},
				},
				Contexts: map[string]*api.Context{
					"dummy": {
						Cluster:    "dummy",
						AuthInfo:   "dummy",
						Extensions: map[string]runtime.Object{},
					},
				},
				CurrentContext: "dummy",
				AuthInfos: map[string]*api.AuthInfo{
					"dummy": {
						ClientCertificateData: []byte{},
						ClientKeyData:         []byte{},
						Extensions:            map[string]runtime.Object{},
					},
				},
				Extensions: map[string]runtime.Object{},
				Preferences: api.Preferences{
					Colors:     false,
					Extensions: map[string]runtime.Object{},
				},
			}, apiConfig)
		}
	})
}
