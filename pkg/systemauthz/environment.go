package systemauthz

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"

	rbac "k8s.io/api/rbac/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"

	walrus "github.com/seal-io/walrus/pkg/apis/walrus/v1"
	"github.com/seal-io/walrus/pkg/kubeclientset"
	"github.com/seal-io/walrus/pkg/kubemeta"
)

// CreateEnvironmentSpace adds the given environment to the project space cluster roles.
func CreateEnvironmentSpace(ctx context.Context, cli ctrlcli.Client, env *walrus.Environment) error {
	ref := kubemeta.GetOwnerRefOfNoCopy(env, walrus.SchemeGroupVersionKind("Project"))
	if ref == nil {
		return errors.New("environment has no project owner reference")
	}

	crs := []*rbac.ClusterRole{
		// Project space viewer.
		{
			ObjectMeta: meta.ObjectMeta{Name: getProjectSpaceViewerClusterRoleName(ref.Name)},
			Rules:      []rbac.PolicyRule{getNamespaceViewerPolicyRuleFor(env.Name)},
		},
		// Project space editor.
		{
			ObjectMeta: meta.ObjectMeta{Name: getProjectSpaceEditorClusterRoleName(ref.Name)},
			Rules:      []rbac.PolicyRule{getNamespaceEditorPolicyRuleFor(env.Name)},
		},
	}
	createAlignFn := func(eCr *rbac.ClusterRole) kubeclientset.AlignWithFn[*rbac.ClusterRole] {
		return func(aCr *rbac.ClusterRole) (_ *rbac.ClusterRole, skip bool, err error) {
			included := slices.ContainsFunc(eCr.Rules, func(er rbac.PolicyRule) bool {
				return slices.ContainsFunc(aCr.Rules, func(ar rbac.PolicyRule) bool {
					return reflect.DeepEqual(er, ar)
				})
			})
			if included {
				return nil, true, nil
			}
			// Append the existing rules.
			aCr.Rules = append(aCr.Rules, eCr.Rules...)
			return aCr, false, nil
		}
	}

	for i := range crs {
		alignFn := createAlignFn(crs[i])

		// Update.
		_, err := kubeclientset.UpdateWithCtrlClient(ctx, cli, crs[i],
			kubeclientset.WithUpdateAlign(alignFn))
		if err != nil {
			return fmt.Errorf("update project cluster role: %s: %w", crs[i].Name, err)
		}
	}

	return nil
}

// DeleteEnvironmentSpace removes the given environment from the project space cluster roles.
func DeleteEnvironmentSpace(ctx context.Context, cli ctrlcli.Client, env *walrus.Environment) error {
	ref := kubemeta.GetOwnerRefOfNoCopy(env, walrus.SchemeGroupVersionKind("Project"))
	if ref == nil {
		return errors.New("environment has no project owner reference")
	}

	crs := []*rbac.ClusterRole{
		// Project space viewer.
		{
			ObjectMeta: meta.ObjectMeta{Name: getProjectSpaceViewerClusterRoleName(ref.Name)},
			Rules:      []rbac.PolicyRule{getNamespaceViewerPolicyRuleFor(env.Name)},
		},
		// Project space editor.
		{
			ObjectMeta: meta.ObjectMeta{Name: getProjectSpaceEditorClusterRoleName(ref.Name)},
			Rules:      []rbac.PolicyRule{getNamespaceEditorPolicyRuleFor(env.Name)},
		},
	}
	createAlignFn := func(eCr *rbac.ClusterRole) kubeclientset.AlignWithFn[*rbac.ClusterRole] {
		return func(aCr *rbac.ClusterRole) (_ *rbac.ClusterRole, skip bool, err error) {
			aCrNewRules := slices.DeleteFunc(aCr.Rules, func(ar rbac.PolicyRule) bool {
				return slices.ContainsFunc(eCr.Rules, func(er rbac.PolicyRule) bool {
					return reflect.DeepEqual(er, ar)
				})
			})
			exclude := len(aCr.Rules) == len(aCrNewRules)
			if exclude {
				return nil, true, nil
			}

			// Update the rules.
			aCr.Rules = aCrNewRules
			return aCr, false, nil
		}
	}

	for i := range crs {
		alignFn := createAlignFn(crs[i])

		// Update.
		_, err := kubeclientset.UpdateWithCtrlClient(ctx, cli, crs[i],
			kubeclientset.WithUpdateAlign(alignFn))
		if err != nil {
			return fmt.Errorf("update project cluster role: %s: %w", crs[i].Name, err)
		}
	}

	return nil
}
