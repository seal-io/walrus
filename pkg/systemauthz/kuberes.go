package systemauthz

import (
	"context"
	"fmt"
	"reflect"
	"slices"

	rbac "k8s.io/api/rbac/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	walrus "github.com/seal-io/walrus/pkg/apis/walrus/v1"
	walruscore "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
	"github.com/seal-io/walrus/pkg/clients/clientset"
	"github.com/seal-io/walrus/pkg/kubeclientset"
	"github.com/seal-io/walrus/pkg/kubeclientset/review"
)

// Initialize initializes Kubernetes resources for authorization.
//
// Initialize creates Kubernetes ClusterRoles for system and project.
func Initialize(ctx context.Context, cli clientset.Interface) error {
	err := review.CanDoCreate(ctx,
		cli.AuthorizationV1().SelfSubjectAccessReviews(),
		review.Simples{
			{
				Group:    rbac.SchemeGroupVersion.Group,
				Version:  rbac.SchemeGroupVersion.Version,
				Resource: "clusterroles",
			},
		},
		review.WithUpdateIfExisted(),
	)
	if err != nil {
		return err
	}

	crCli := cli.RbacV1().ClusterRoles()
	eCrs := []*rbac.ClusterRole{
		// Project viewer.
		{
			ObjectMeta: meta.ObjectMeta{
				Name: walrus.ProjectSubjectRoleViewer.String(),
			},
			Rules: []rbac.PolicyRule{
				{
					APIGroups: []string{
						walrus.GroupName,
						walruscore.GroupName,
					},
					Resources: []string{
						rbac.ResourceAll,
					},
					Verbs: []string{
						"get",
						"list",
						"watch",
					},
				},
			},
		},
		// Project member.
		{
			ObjectMeta: meta.ObjectMeta{
				Name: walrus.ProjectSubjectRoleMember.String(),
			},
			Rules: []rbac.PolicyRule{
				{
					APIGroups: []string{
						walrus.GroupName,
						walruscore.GroupName,
					},
					Resources: []string{
						rbac.ResourceAll,
					},
					Verbs: []string{
						"get",
						"list",
						"watch",
					},
				},
				{
					APIGroups: []string{
						walrus.GroupName,
					},
					Resources: []string{
						"environments",
						"resourcecomponents",
					},
					Verbs: []string{
						rbac.VerbAll,
					},
				},
				{
					APIGroups: []string{
						walruscore.GroupName,
					},
					Resources: []string{
						"resources",
						"resources/components",
					},
					Verbs: []string{
						rbac.VerbAll,
					},
				},
			},
		},
		// Project owner.
		{
			ObjectMeta: meta.ObjectMeta{
				Name: walrus.ProjectSubjectRoleOwner.String(),
			},
			Rules: []rbac.PolicyRule{
				{
					APIGroups: []string{
						walrus.GroupName,
						walruscore.GroupName,
					},
					Resources: []string{
						rbac.ResourceAll,
					},
					Verbs: []string{
						rbac.VerbAll,
					},
				},
			},
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

	for i := range eCrs {
		alignFn := createAlignFn(eCrs[i])

		// Create.
		_, err = kubeclientset.Create(ctx, crCli, eCrs[i],
			kubeclientset.WithUpdateIfExisted(alignFn))
		if err != nil {
			return fmt.Errorf("install project cluster role %q: %w", eCrs[i].Name, err)
		}
	}

	return nil
}
