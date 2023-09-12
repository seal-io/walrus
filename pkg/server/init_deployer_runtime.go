package server

import (
	"context"
	"fmt"
	"reflect"

	batch "k8s.io/api/batch/v1"
	coordination "k8s.io/api/coordination/v1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/seal-io/walrus/pkg/dao/types"
)

// setupDeployerRuntime configures the deployer runtime at initialization phase,
// like Namespace, RBAC, etc.
func (r *Server) setupDeployerRuntime(ctx context.Context, opts initOptions) error {
	cli, err := kubernetes.NewForConfig(opts.K8sConfig)
	if err != nil {
		return fmt.Errorf("failed to create client via cfg: %w", err)
	}

	cs := []func(context.Context, *kubernetes.Clientset) error{
		applyDeployWorkspace,
		applyDeployPermission,
	}

	for i := range cs {
		err = cs[i](ctx, cli)
		if err != nil {
			return fmt.Errorf("failed to execute preparation: %w", err)
		}
	}

	return nil
}

// applyDeployWorkspace applies the Kubernetes Namespace to run deployer at next.
func applyDeployWorkspace(ctx context.Context, cli *kubernetes.Clientset) error {
	ns := core.Namespace{
		ObjectMeta: meta.ObjectMeta{
			Name: types.WalrusSystemNamespace,
		},
	}

	_, err := cli.CoreV1().
		Namespaces().
		Create(ctx, &ns, meta.CreateOptions{})
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

// applyDeployPermission applies the Kubernetes RBAC resources for deployer running.
func applyDeployPermission(ctx context.Context, cli *kubernetes.Clientset) error {
	sa := core.ServiceAccount{
		ObjectMeta: meta.ObjectMeta{
			Namespace: types.WalrusSystemNamespace,
			Name:      types.DeployerServiceAccountName,
		},
	}

	_, err := cli.CoreV1().
		ServiceAccounts(sa.Namespace).
		Create(ctx, &sa, meta.CreateOptions{})
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return err
	}

	r := rbac.Role{
		ObjectMeta: meta.ObjectMeta{
			Namespace: types.WalrusSystemNamespace,
			Name:      types.DeployerServiceAccountName,
		},
		Rules: []rbac.PolicyRule{
			{
				APIGroups: []string{batch.GroupName},
				Resources: []string{"jobs"},
				Verbs:     []string{rbac.VerbAll},
			},
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"secrets", "pods", "pods/log", "events"},
				Verbs:     []string{rbac.VerbAll},
			},
			{
				APIGroups: []string{coordination.GroupName},
				Resources: []string{"leases"},
				Verbs:     []string{rbac.VerbAll},
			},
		},
	}

	r_, err := cli.RbacV1().
		Roles(r.Namespace).
		Get(ctx, r.Name, meta.GetOptions{
			ResourceVersion: "0",
		})
	if err != nil && !kerrors.IsNotFound(err) {
		return err
	}

	switch {
	case r_ == nil || r_.Name == "" || r_.DeletionTimestamp != nil:
		// Create.
		_, err = cli.RbacV1().
			Roles(r.Namespace).
			Create(ctx, &r, meta.CreateOptions{})
		if err != nil {
			return err
		}
	case !reflect.DeepEqual(r.Rules, r_.Rules):
		// Update.
		r.Labels = r_.Labels
		r.Annotations = r_.Annotations
		r.ResourceVersion = r_.ResourceVersion

		_, err = cli.RbacV1().
			Roles(r.Namespace).
			Update(ctx, &r, meta.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	rb := rbac.RoleBinding{
		ObjectMeta: meta.ObjectMeta{
			Namespace: types.WalrusSystemNamespace,
			Name:      types.DeployerServiceAccountName,
		},
		Subjects: []rbac.Subject{
			{
				Kind: rbac.ServiceAccountKind,
				Name: types.DeployerServiceAccountName,
			},
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "Role",
			Name:     types.DeployerServiceAccountName,
		},
	}

	rb_, err := cli.RbacV1().
		RoleBindings(rb.Namespace).
		Get(ctx, rb.Name, meta.GetOptions{
			ResourceVersion: "0",
		})
	if err != nil && !kerrors.IsNotFound(err) {
		return err
	}

	switch {
	case rb_ == nil || rb_.Name == "" || rb_.DeletionTimestamp != nil:
		// Create.
		_, err = cli.RbacV1().
			RoleBindings(rb.Namespace).
			Create(ctx, &rb, meta.CreateOptions{})
		if err != nil {
			return err
		}
	case !reflect.DeepEqual(rb.RoleRef, rb_.RoleRef):
		// Delete.
		err = cli.RbacV1().
			RoleBindings(rb.Namespace).
			Delete(ctx, rb.Name, meta.DeleteOptions{
				PropagationPolicy: point(meta.DeletePropagationForeground),
			})
		if err != nil && !kerrors.IsNotFound(err) {
			return err
		}

		// Recreate.
		_, err = cli.RbacV1().
			RoleBindings(rb.Namespace).
			Create(ctx, &rb, meta.CreateOptions{})
		if err != nil {
			return err
		}
	case !reflect.DeepEqual(rb.Subjects, rb_.Subjects):
		// Update.
		rb.Labels = rb_.Labels
		rb.Annotations = rb_.Annotations
		rb.ResourceVersion = rb_.ResourceVersion

		_, err = cli.RbacV1().
			RoleBindings(rb.Namespace).
			Update(ctx, &rb, meta.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func point[T ~string | int | int64 | bool](s T) *T {
	return &s
}
