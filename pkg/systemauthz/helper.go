package systemauthz

import (
	"fmt"

	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"

	walrus "github.com/seal-io/walrus/pkg/apis/walrus/v1"
)

const (
	_ProjectSpaceViewerClusterRoleNameFormat = "walrus-project-%s-space-viewer"
	_ProjectSpaceEditorClusterRoleNameFormat = "walrus-project-%s-space-editor"
)

// GetProjectSpaceViewerClusterRoleName returns the project space viewer cluster role name of the given project.
func GetProjectSpaceViewerClusterRoleName(proj *walrus.Project) string {
	return getProjectSpaceViewerClusterRoleName(proj.Name)
}

// GetProjectSpaceEditorClusterRoleName returns the project space editor cluster role name of the given project.
func GetProjectSpaceEditorClusterRoleName(proj *walrus.Project) string {
	return getProjectSpaceEditorClusterRoleName(proj.Name)
}

func getProjectSpaceViewerClusterRoleName(
	projName string,
) string {
	return fmt.Sprintf(_ProjectSpaceViewerClusterRoleNameFormat, projName)
}

func getProjectSpaceEditorClusterRoleName(
	projName string,
) string {
	return fmt.Sprintf(_ProjectSpaceEditorClusterRoleNameFormat, projName)
}

func getNamespaceViewerPolicyRuleFor(
	namespaceName string,
) rbac.PolicyRule {
	return rbac.PolicyRule{
		APIGroups: []string{
			core.GroupName,
		},
		Resources: []string{
			"namespaces",
		},
		ResourceNames: []string{
			namespaceName,
		},
		Verbs: []string{
			"get",
			"list",
			"watch",
		},
	}
}

func getNamespaceEditorPolicyRuleFor(
	namespaceName string,
) rbac.PolicyRule {
	return rbac.PolicyRule{
		APIGroups: []string{
			core.GroupName,
		},
		Resources: []string{
			"namespaces",
		},
		ResourceNames: []string{
			namespaceName,
		},
		Verbs: []string{
			rbac.ResourceAll,
		},
	}
}
