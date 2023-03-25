package util

import (
	"fmt"
)

// GetK8sSecretName returns the secret config name for the given connector.
// used for kubernetes connector. or other connectors which need to store the kubeconfig in secret.
func GetK8sSecretName(connectorID string) string {
	return fmt.Sprintf("config%s", connectorID)
}
