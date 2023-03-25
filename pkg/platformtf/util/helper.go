package util

import (
	"fmt"
)

// GetSecretK8sConfigName returns the secret config name for the given connector.
// used for kubernetes connector. or other connectors which need to store the kubeconfig in secret.
func GetSecretK8sConfigName(connectorID string) string {
	return fmt.Sprintf("config%s", connectorID)
}
