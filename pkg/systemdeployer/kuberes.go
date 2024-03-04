package systemdeployer

import (
	"context"
	"fmt"
	"reflect"
	"slices"

	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/seal-io/walrus/pkg/clients/clientset"
	"github.com/seal-io/walrus/pkg/kubeclientset"
	"github.com/seal-io/walrus/pkg/kubeclientset/review"
)

// ClusterRoleName is the Kubernetes ClusterRole name for running deployer.
const ClusterRoleName = "walrus-deployer"

// Initialize initializes Kubernetes resources for deployer.
//
// Initialize creates a Kubernetes ClusterRole for running deployer.
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
	eCr := &rbac.ClusterRole{
		ObjectMeta: meta.ObjectMeta{
			Name: ClusterRoleName,
		},
		Rules: []rbac.PolicyRule{
			// The below rules are used for kaniko to build images.
			{
				APIGroups: []string{batch.GroupName},
				Resources: []string{"jobs"},
				Verbs:     []string{rbac.VerbAll},
			},
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"secrets", "pods", "pods/log"},
				Verbs:     []string{rbac.VerbAll},
			},
		},
	}
	alignFn := func(aCr *rbac.ClusterRole) (*rbac.ClusterRole, bool, error) {
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

	_, err = kubeclientset.Create(ctx, crCli, eCr,
		kubeclientset.WithUpdateIfExisted(alignFn))
	if err != nil {
		return fmt.Errorf("install deployer cluster role %q: %w", eCr.GetName(), err)
	}

	return nil
}
