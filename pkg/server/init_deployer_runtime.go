package server

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/seal-io/seal/pkg/dao/types"
)

const (
	roleBindingName = "deployer-rolebinding"
)

func (r *Server) initDeployerRuntime(ctx context.Context, opts initOptions) error {
	clientSet, err := kubernetes.NewForConfig(opts.K8sConfig)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client set: %w", err)
	}

	if err = createNamespace(ctx, clientSet); err != nil {
		return fmt.Errorf("failed to create deployer namespace: %w", err)
	}

	if err = createServiceAccount(ctx, clientSet); err != nil {
		return fmt.Errorf("failed to create deployer service account: %w", err)
	}

	if err = createRoleBinding(ctx, clientSet); err != nil {
		return fmt.Errorf("failed to create deployer role binding: %w", err)
	}

	return nil
}

// createNamespace create a job namespace if not exist.
func createNamespace(ctx context.Context, clientSet *kubernetes.Clientset) error {
	createNamespace := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: types.SealSystemNamespace,
		},
	}
	// create namespace if not exist.
	_, err := clientSet.CoreV1().
		Namespaces().
		Create(ctx, &createNamespace, metav1.CreateOptions{})
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

// createServiceAccount create a service account for deployer if not exist.
func createServiceAccount(ctx context.Context, clientSet *kubernetes.Clientset) error {
	sa := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: types.DeployerServiceAccountName,
		},
	}
	// create service account if not exist.
	_, err := clientSet.CoreV1().
		ServiceAccounts(types.SealSystemNamespace).
		Create(ctx, &sa, metav1.CreateOptions{})
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

// createRoleBinding create a role binding for deployer if not exist.
func createRoleBinding(ctx context.Context, clientSet *kubernetes.Clientset) error {
	rb := rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: roleBindingName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      types.DeployerServiceAccountName,
				Namespace: types.SealSystemNamespace,
			},
		},
		// TODO minimum role
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     "admin",
			APIGroup: rbacv1.GroupName,
		},
	}
	// create role binding if not exist.
	_, err := clientSet.RbacV1().
		RoleBindings(types.SealSystemNamespace).
		Create(ctx, &rb, metav1.CreateOptions{})
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}
