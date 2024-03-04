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

const (
	// SystemUserClusterRoleName is the Kubernetes ClusterRole name for system user.
	SystemUserClusterRoleName = "walrus-system-user"
	// SystemManagerClusterRoleName is the Kubernetes ClusterRole name for system manager.
	SystemManagerClusterRoleName = "walrus-system-manager"
	// SystemAdminClusterRoleName is the Kubernetes ClusterRole name for system admin.
	SystemAdminClusterRoleName = "walrus-system-admin"

	// ProjectViewerClusterRoleName is the Kubernetes ClusterRole name for project viewer.
	ProjectViewerClusterRoleName = "walrus-project-viewer"
	// ProjectMemberClusterRoleName is the Kubernetes ClusterRole name for project member.
	ProjectMemberClusterRoleName = "walrus-project-member"
	// ProjectOwnerClusterRoleName is the Kubernetes ClusterRole name for project owner.
	ProjectOwnerClusterRoleName = "walrus-project-owner"
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
		review.WithUpdate(),
	)
	if err != nil {
		return err
	}

	crCli := cli.RbacV1().ClusterRoles()
	eCrs := []*rbac.ClusterRole{
		// Project viewer.
		{
			ObjectMeta: meta.ObjectMeta{
				Name: ProjectViewerClusterRoleName,
			},
			Rules: []rbac.PolicyRule{
				{
					APIGroups: []string{
						walrus.SchemeGroupVersion.String(),
						walruscore.SchemeGroupVersion.String(),
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
				Name: ProjectMemberClusterRoleName,
			},
			Rules: []rbac.PolicyRule{
				{
					APIGroups: []string{
						walrus.SchemeGroupVersion.String(),
						walruscore.SchemeGroupVersion.String(),
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
						walrus.SchemeGroupVersion.String(),
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
						walruscore.SchemeGroupVersion.String(),
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
				Name: ProjectOwnerClusterRoleName,
			},
			Rules: []rbac.PolicyRule{
				{
					APIGroups: []string{
						walrus.SchemeGroupVersion.String(),
						walruscore.SchemeGroupVersion.String(),
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
		_, err = kubeclientset.Create(ctx, crCli, eCrs[i],
			kubeclientset.CreateOrUpdate(createAlignFn(eCrs[i])))
		if err != nil {
			return fmt.Errorf("install project cluster role %q: %w", eCrs[i].Name, err)
		}
	}

	return nil
}
