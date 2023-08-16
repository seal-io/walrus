package server

import (
	"context"
	"net/http"

	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/role"
	"github.com/seal-io/walrus/pkg/dao/model/subject"
	"github.com/seal-io/walrus/pkg/dao/types"
)

// createBuiltinRbac creates the built-in RBAC resources.
func (r *Server) createBuiltinRbac(ctx context.Context, opts initOptions) (err error) {
	err = createRoles(ctx, opts.ModelClient)
	if err != nil {
		return
	}

	err = createSubjects(ctx, opts.ModelClient)
	if err != nil {
		return
	}

	return
}

func createRoles(ctx context.Context, mc model.ClientSet) error {
	builtin := []*model.Role{
		// System anonymity.
		{
			ID:          types.SystemRoleAnonymity,
			Kind:        types.RoleKindSystem,
			Description: "The role who cannot pass system authenticating.",
			Policies: types.RolePolicies{
				{
					Actions: types.RolePolicyFields(http.MethodPost),
					Paths: types.RolePolicyFields(
						"/account/login"),
				},
				{
					Actions: types.RolePolicyFields(http.MethodGet),
					Resources: types.RolePolicyFields(
						"settings"),
					ResourceRefers: types.RolePolicyFields(
						"BootPwdGainSource",
						"ServeUrl"),
				},
			},
			Session: true,
			Builtin: true,
		},

		// System user.
		{
			ID:          types.SystemRoleUser,
			Kind:        types.RoleKindSystem,
			Description: "The role who can manage its own projects.",
			Policies: types.RolePolicies{
				{
					Actions: types.RolePolicyFields("*"),
					Resources: types.RolePolicyFields(
						"projects"),
				},
				{
					Actions: types.RolePolicyFields(http.MethodGet),
					Resources: types.RolePolicyFields(
						"catalogs",
						"connectors",
						"environments",
						"perspectives",
						"roles",
						"settings",
						"templates",
						"templateCompletions",
						"variables"),
				},
				{
					Actions: types.RolePolicyFields(http.MethodPost, http.MethodGet),
					Resources: types.RolePolicyFields(
						"costs",
						"dashboards"), // POST for larger query body.
				},
				{
					Actions: types.RolePolicyFields(http.MethodGet, http.MethodPut),
					Paths: types.RolePolicyFields(
						"/account/info"),
				},
				{
					Actions: types.RolePolicyFields(http.MethodDelete),
					Paths: types.RolePolicyFields(
						"/account/logout"),
				},
				{
					Actions: types.RolePolicyFields(http.MethodGet, http.MethodPost),
					Paths: types.RolePolicyFields(
						"/account/tokens"),
				},
				{
					Actions: types.RolePolicyFields(http.MethodGet, http.MethodDelete),
					Paths: types.RolePolicyFields(
						"/account/tokens/:token"),
				},
			},
			Session: true,
			Builtin: true,
		},

		// System platform engineer.
		{
			ID:          types.SystemRolePlatformEngineer,
			Kind:        types.RoleKindSystem,
			Description: "The role who can manage the system resources.",
			Policies: types.RolePolicies{
				{
					Actions: types.RolePolicyFields("*"),
					Resources: types.RolePolicyFields(
						"catalogs",
						"connectors",
						"settings",
						"templates",
						"templateCompletions",
						"variables"),
				},
			},
			Session: false,
			Builtin: true,
		},

		// System admin.
		{
			ID:          types.SystemRoleAdmin,
			Kind:        types.RoleKindSystem,
			Description: "The role who can manage all resources, including system and project.",
			Policies: types.RolePolicies{
				{
					Actions: types.RolePolicyFields("*"),
					Resources: types.RolePolicyFields(
						"catalogs",
						"connectors",
						"environments",
						"perspectives",
						"projects",
						"projectSubjects",
						"roles",
						"settings",
						"subjects",
						"templates",
						"templateCompletions",
						"variables"),
				},
			},
			Session: false,
			Builtin: true,
		},

		// Project viewer.
		{
			ID:          types.ProjectRoleViewer,
			Kind:        types.RoleKindProject,
			Description: "The role who can only read the resources below the project.",
			Policies: types.RolePolicies{
				{
					Actions: types.RolePolicyFields(http.MethodGet),
					Resources: types.RolePolicyFields(
						"connectors",
						"projects",
						"environments",
						"services",
						"serviceResources",
						"serviceRevisions",
						"variables"),
				},
			},
			Session: false,
			Builtin: true,
		},

		// Project member.
		{
			ID:          types.ProjectRoleMember,
			Kind:        types.RoleKindProject,
			Description: "The role who can manage the resources below the project, excluding rbac.",
			Policies: types.RolePolicies{
				{
					Actions: types.RolePolicyFields(http.MethodGet),
					Resources: types.RolePolicyFields(
						"connectors",
						"projects"),
				},
				{
					Actions: types.RolePolicyFields("*"),
					Resources: types.RolePolicyFields(
						"environments",
						"services",
						"serviceResources",
						"serviceRevisions",
						"variables"),
				},
			},
			Session: false,
			Builtin: true,
		},

		// Project owner.
		{
			ID:          types.ProjectRoleOwner,
			Kind:        types.RoleKindProject,
			Description: "The role who can manage the whole project.",
			Policies: types.RolePolicies{
				{
					Actions: types.RolePolicyFields(http.MethodGet),
					Resources: types.RolePolicyFields(
						"subjects"),
				},
				{
					Actions: types.RolePolicyFields("*"),
					Resources: types.RolePolicyFields(
						"connectors",
						"projects",
						"projectSubjects",
						"environments",
						"services",
						"serviceResources",
						"serviceRevisions",
						"variables"),
				},
			},
			Session: false,
			Builtin: true,
		},
	}

	return mc.Roles().CreateBulk().
		Set(builtin...).
		OnConflictColumns(role.FieldID).
		Update(func(upsert *model.RoleUpsert) {
			upsert.UpdatePolicies()
			upsert.UpdateUpdateTime()
		}).
		Exec(ctx)
}

func createSubjects(ctx context.Context, mc model.ClientSet) error {
	builtin := []*model.Subject{
		// System admin.
		{
			Kind:        types.SubjectKindUser,
			Domain:      types.SubjectDomainBuiltin,
			Name:        "admin",
			Description: "The administrator user.",
			Edges: model.SubjectEdges{
				Roles: []*model.SubjectRoleRelationship{
					{
						RoleID: types.SystemRoleAdmin,
					},
				},
			},
			Builtin: true,
		},
	}

	return mc.WithTx(ctx, func(tx *model.Tx) error {
		return tx.Subjects().CreateBulk().
			Set(builtin...).
			OnConflictColumns(
				subject.FieldKind,
				subject.FieldDomain,
				subject.FieldName,
			).
			DoNothing().
			ExecE(ctx, dao.SubjectRolesEdgeSave)
	})
}
