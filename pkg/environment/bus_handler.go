package environment

import (
	"context"
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	envbus "github.com/seal-io/seal/pkg/bus/environment"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/object"
	opk8s "github.com/seal-io/seal/pkg/operator/k8s"
)

// SyncManagedKubernetesNamespace creates/deletes Seal managed kubernetes namespace
// for environments using kubernetes connector.
func SyncManagedKubernetesNamespace(ctx context.Context, m envbus.BusMessage) error {
	for _, e := range m.Refers {
		if len(e.Edges.Connectors) == 0 {
			continue
		}

		if e.Edges.Project == nil {
			return errors.New("invalid edge: empty project")
		}

		nsName := GetManagedNamespaceName(e)
		if nsName == "" {
			continue
		}

		connectorIDs := make([]object.ID, len(e.Edges.Connectors))
		for i := range e.Edges.Connectors {
			connectorIDs[i] = e.Edges.Connectors[i].ConnectorID
		}

		connectors, err := m.TransactionalModelClient.Connectors().Query().
			Where(connector.IDIn(connectorIDs...)).
			All(ctx)
		if err != nil {
			return err
		}

		for _, c := range connectors {
			switch {
			case c.Type == types.ConnectorTypeK8s && m.Event == envbus.EventDelete:
				if err = deleteNamespace(ctx, c, nsName); err != nil {
					return err
				}
			case c.Type == types.ConnectorTypeK8s && m.Event != envbus.EventDelete:
				if err = createNamespaceIfNotExist(ctx, c, nsName); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func createNamespaceIfNotExist(ctx context.Context, connector *model.Connector, name string) error {
	restCfg, err := opk8s.GetConfig(*connector)
	if err != nil {
		return err
	}

	corev1Client, err := corev1.NewForConfig(restCfg)
	if err != nil {
		return fmt.Errorf("error creating kubernetes core client: %w", err)
	}

	if _, err = corev1Client.Namespaces().Get(ctx, name, metav1.GetOptions{}); !kerrors.IsNotFound(err) {
		return nil
	}

	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				types.LabelSealManaged: "true",
			},
		},
	}

	_, err = corev1Client.Namespaces().Create(ctx, namespace, metav1.CreateOptions{})

	return err
}

func deleteNamespace(ctx context.Context, connector *model.Connector, name string) error {
	restCfg, err := opk8s.GetConfig(*connector)
	if err != nil {
		return err
	}

	corev1Client, err := corev1.NewForConfig(restCfg)
	if err != nil {
		return fmt.Errorf("error creating kubernetes core client: %w", err)
	}

	current, err := corev1Client.Namespaces().Get(ctx, name, metav1.GetOptions{})
	if kerrors.IsNotFound(err) || current.Labels[types.LabelSealManaged] != "true" {
		return nil
	}

	return corev1Client.Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
}
