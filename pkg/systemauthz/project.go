package systemauthz

import (
	"context"
	"errors"
	"fmt"

	"github.com/seal-io/utils/stringx"
	rbac "k8s.io/api/rbac/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	authnuser "k8s.io/apiserver/pkg/authentication/user"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/utils/ptr"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"

	walrus "github.com/seal-io/walrus/pkg/apis/walrus/v1"
	"github.com/seal-io/walrus/pkg/kubeclientset"
	"github.com/seal-io/walrus/pkg/kubemeta"
	"github.com/seal-io/walrus/pkg/systemmeta"
)

// CreateProjectSpace creates the project space cluster roles.
func CreateProjectSpace(ctx context.Context, cli ctrlcli.Client, proj *walrus.Project) error {
	crs := []*rbac.ClusterRole{
		// Project space viewer.
		{
			ObjectMeta: meta.ObjectMeta{Name: getProjectSpaceViewerClusterRoleName(proj.Name)},
			Rules:      []rbac.PolicyRule{getNamespaceViewerPolicyRuleFor(proj.Name)},
		},
		// Project space editor.
		{
			ObjectMeta: meta.ObjectMeta{Name: getProjectSpaceEditorClusterRoleName(proj.Name)},
			Rules:      []rbac.PolicyRule{getNamespaceEditorPolicyRuleFor(proj.Name)},
		},
	}

	for i := range crs {
		systemmeta.NoteResource(crs[i], "roles", map[string]string{"project": proj.Name})
		kubemeta.ControlOnWithoutBlock(crs[i], proj, walrus.SchemeGroupVersionKind("Project"))

		// Create.
		_, err := kubeclientset.CreateWithCtrlClient(ctx, cli, crs[i])
		if err != nil {
			return fmt.Errorf("create project cluster role: %s: %w", crs[i].Name, err)
		}
	}

	// Bind project subject role.
	err := BindProjectSubjectRole(ctx, cli, proj, walrus.ProjectSubjectRoleOwner)
	if err != nil {
		return fmt.Errorf("bind project subject role: %w", err)
	}

	return nil
}

// BindProjectSubjectRoleFor binds the given subject role of the given project for the given user info.
func BindProjectSubjectRoleFor(
	ctx context.Context,
	cli ctrlcli.Client,
	proj *walrus.Project,
	role walrus.ProjectSubjectRole,
	uInfo authnuser.Info,
) error {
	spaceRoleName := getProjectSpaceEditorClusterRoleName(proj.Name)
	if role != walrus.ProjectSubjectRoleOwner {
		spaceRoleName = getProjectSpaceViewerClusterRoleName(proj.Name)
	}
	resourceRoleName := role.String()

	id := stringx.SumByFNV64a(uInfo.GetName())
	crbs := []*rbac.ClusterRoleBinding{
		// Project resource role binding.
		{
			ObjectMeta: meta.ObjectMeta{
				Name: fmt.Sprintf("%s-%s", resourceRoleName, id),
			},
			Subjects: []rbac.Subject{
				{
					APIGroup: rbac.GroupName,
					Kind:     "User",
					Name:     uInfo.GetName(),
				},
			},
			RoleRef: rbac.RoleRef{
				APIGroup: rbac.GroupName,
				Kind:     "ClusterRole",
				Name:     resourceRoleName,
			},
		},
		// Project space role binding.
		{
			ObjectMeta: meta.ObjectMeta{
				Name: fmt.Sprintf("%s-%s-space", resourceRoleName, id),
			},
			Subjects: []rbac.Subject{
				{
					APIGroup: rbac.GroupName,
					Kind:     "User",
					Name:     uInfo.GetName(),
				},
			},
			RoleRef: rbac.RoleRef{
				APIGroup: rbac.GroupName,
				Kind:     "ClusterRole",
				Name:     spaceRoleName,
			},
		},
	}

	for i := range crbs {
		systemmeta.NoteResource(crbs[i], "rolebindings", map[string]string{"project": proj.Name})
		kubemeta.ControlOnWithoutBlock(crbs[i], proj, walrus.SchemeGroupVersionKind("Project"))

		// Create.
		_, err := kubeclientset.CreateWithCtrlClient(ctx, cli, crbs[i])
		if err != nil {
			return fmt.Errorf("create project cluster role binding: %s: %w", crbs[i].Name, err)
		}
	}

	return nil
}

// BindProjectSubjectRole is similar to BindProjectSubjectRoleFor,
// but it can retrieve user information from context automatically.
func BindProjectSubjectRole(ctx context.Context, cli ctrlcli.Client, proj *walrus.Project, role walrus.ProjectSubjectRole) error {
	ui, ok := genericapirequest.UserFrom(ctx)
	if !ok {
		return errors.New("cannot retrieve kubernetes request user information from context")
	}
	return BindProjectSubjectRoleFor(ctx, cli, proj, role, ui)
}

// UnbindProjectSubjectRoleFor unbinds the given subject role of the given project from the given user info.
func UnbindProjectSubjectRoleFor(
	ctx context.Context,
	cli ctrlcli.Client,
	role walrus.ProjectSubjectRole,
	uInfo authnuser.Info,
) error {
	resourceRoleName := role.String()

	id := stringx.SumByFNV64a(uInfo.GetName())
	crbs := []*rbac.ClusterRoleBinding{
		// Project resource role binding.
		{
			ObjectMeta: meta.ObjectMeta{
				Name: fmt.Sprintf("%s-%s", resourceRoleName, id),
			},
		},
		// Project space role binding.
		{
			ObjectMeta: meta.ObjectMeta{
				Name: fmt.Sprintf("%s-%s-space", resourceRoleName, id),
			},
		},
	}

	for i := range crbs {
		// Delete.
		err := kubeclientset.DeleteWithCtrlClient(ctx, cli, crbs[i],
			kubeclientset.WithDeleteMetaOptions(meta.DeleteOptions{
				PropagationPolicy: ptr.To(meta.DeletePropagationForeground),
			}))
		if err != nil {
			return err
		}
	}

	return nil
}

// UnbindProjectSubjectRole is similar to UnbindProjectSubjectRoleFor,
// but it can retrieve user information from context automatically.
func UnbindProjectSubjectRole(ctx context.Context, cli ctrlcli.Client, role walrus.ProjectSubjectRole) error {
	ui, ok := genericapirequest.UserFrom(ctx)
	if !ok {
		return errors.New("cannot retrieve kubernetes request user information from context")
	}
	return UnbindProjectSubjectRoleFor(ctx, cli, role, ui)
}
