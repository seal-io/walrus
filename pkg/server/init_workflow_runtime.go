package server

import (
	"context"
	"fmt"
	"reflect"

	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/workflow/installer/argoworkflows"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/pointer"
)

// setupDeployerRuntime configures the deployer runtime at initialization phase,
// like Namespace, RBAC, etc.
func (r *Server) setupWorkflowRuntime(ctx context.Context, opts initOptions) error {
	cli, err := kubernetes.NewForConfig(opts.K8sConfig)
	if err != nil {
		return fmt.Errorf("failed to create client via cfg: %w", err)
	}

	cs := []func(context.Context, *kubernetes.Clientset) error{
		applyWorkflowWorkspace,
		applyWorkflowPermission,
	}

	for i := range cs {
		err = cs[i](ctx, cli)
		if err != nil {
			return fmt.Errorf("failed to execute preparation: %w", err)
		}
	}

	applyWorkflowDeployment(opts.ModelClient, opts.K8sConfig)

	return nil
}

// applyWorkflowWorkspace applies the Kubernetes Namespace to run workflow at next.
func applyWorkflowWorkspace(ctx context.Context, cli *kubernetes.Clientset) error {
	if !isCreationAllowed(ctx, cli, "namespaces") {
		return nil
	}

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

// applyWorkflowPermission applies the Kubernetes RBAC resources for workflow running.
func applyWorkflowPermission(ctx context.Context, cli *kubernetes.Clientset) error {
	if !isCreationAllowed(ctx, cli, "serviceaccounts", "roles", "rolebindings", "secrets") {
		return nil
	}

	sa := core.ServiceAccount{
		ObjectMeta: meta.ObjectMeta{
			Namespace: types.WalrusSystemNamespace,
			Name:      types.WorkflowServiceAccountName,
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
			Name:      types.WorkflowServiceAccountName,
		},
		Rules: []rbac.PolicyRule{
			// The below rules are used for workflow execution.
			{
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{"get", "watch", "patch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"pods/logs"},
				Verbs:     []string{"get", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{"argoproj.io"},
				Resources: []string{"workflowtasksets"},
				Verbs:     []string{"watch", "list"},
			},
			{
				APIGroups: []string{"argoproj.io"},
				Resources: []string{"workflowtaskresults"},
				Verbs:     []string{"create", "patch"},
			},
			{
				APIGroups: []string{"argoproj.io"},
				Resources: []string{"workflowtasksets/status"},
				Verbs:     []string{"patch"},
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
			Name:      types.WorkflowServiceAccountName,
		},
		Subjects: []rbac.Subject{
			{
				Kind: rbac.ServiceAccountKind,
				Name: types.WorkflowServiceAccountName,
			},
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "Role",
			Name:     r.Name,
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
				PropagationPolicy: pointer.Ref(meta.DeletePropagationForeground),
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

	token := core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name: fmt.Sprintf("%s.service-account-token", types.WorkflowServiceAccountName),
			Annotations: map[string]string{
				"kubernetes.io/service-account.name": types.WorkflowServiceAccountName,
			},
		},
		Type: core.SecretTypeServiceAccountToken,
	}

	_, err = cli.CoreV1().
		Secrets(types.WalrusSystemNamespace).
		Create(ctx, &token, meta.CreateOptions{})
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

func applyWorkflowDeployment(mc model.ClientSet, config *rest.Config) {
	gopool.Go(func() {
		err := argoworkflows.Install(context.Background(), mc, config)
		if err != nil {
			log.WithName("workflow").Errorf("failed to install argo workflows: %v", err)
		}
	})
}
