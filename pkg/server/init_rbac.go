package server

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/types"
)

func (r *Server) initRbac(ctx context.Context, opts initOptions) (err error) {
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
					Paths:   types.RolePolicyFields("/account/login"),
				},
				{
					Actions:   types.RolePolicyFields(http.MethodGet),
					Resources: types.RolePolicyFields("settings"),
					ObjectIDs: types.RolePolicyFields("FirstLogin"),
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
					Actions:   types.RolePolicyFields("*"),
					Resources: types.RolePolicyFields("projects", "tokens"),
				},
				{
					Actions: types.RolePolicyFields(http.MethodGet),
					Resources: types.RolePolicyFields(
						"secrets", "connectors", "environments",
						"templates", "templateVersions", "templateCompletions",
						"perspectives", "settings", "roles"),
				},
				{
					Actions:   types.RolePolicyFields(http.MethodPost, http.MethodGet),
					Resources: types.RolePolicyFields("costs", "dashboards"), // POST for larger query body.
				},
				{
					Actions: types.RolePolicyFields(http.MethodGet, http.MethodPost),
					Paths:   types.RolePolicyFields("/account/info"),
				},
				{
					Actions: types.RolePolicyFields(http.MethodPost),
					Paths:   types.RolePolicyFields("/account/logout"),
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
						"secrets", "connectors", "templates",
						"templateVersions", "templateCompletions",
						"settings"),
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
						"secrets", "connectors", "templates",
						"templateVersions", "templateCompletions",
						"settings", "roles", "subjects",
						"subjectRoles", "perspectives"),
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
						"services", "serviceRevisions",
						"serviceResources", "environments", "connectors", "secrets"),
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
					Actions: types.RolePolicyFields("*"),
					Resources: types.RolePolicyFields(
						"services", "serviceRevisions",
						"serviceResources", "environments", "connectors", "secrets"),
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
					Actions: types.RolePolicyFields("*"),
					Resources: types.RolePolicyFields(
						"projects", "services",
						"serviceRevisions", "serviceResources",
						"environments", "connectors", "secrets",
						"subjectRoles"),
				},
				{
					Actions:   types.RolePolicyFields(http.MethodGet),
					Resources: types.RolePolicyFields("subjects"),
				},
			},
			Session: false,
			Builtin: true,
		},
	}

	creates, err := dao.RoleCreates(mc, builtin...)
	if err != nil {
		return err
	}

	for i := range creates {
		err = creates[i].
			OnConflictColumns(role.FieldID).
			Update(func(upsert *model.RoleUpsert) {
				upsert.UpdatePolicies()
				upsert.UpdateUpdateTime()
			}).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	return nil
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

	creates, err := dao.SubjectCreates(mc, builtin...)
	if err != nil {
		return err
	}

	for i := range creates {
		err = creates[i].
			OnConflictColumns(
				subject.FieldKind,
				subject.FieldDomain,
				subject.FieldName,
			).
			DoNothing().
			Exec(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// No rows error is reasonable for nothing updating.
				continue
			}

			return err
		}
	}

	return nil
}
