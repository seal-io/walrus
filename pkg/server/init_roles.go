package server

import (
	"context"
	"net/http"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/types"
)

func (r *Server) initRoles(ctx context.Context, opts initOptions) error {
	builtin := []*model.Role{
		// System anonymity.
		{
			Domain:      "system",
			Name:        "anonymity",
			Description: "system/anonymity",
			Policies: types.RolePolicies{
				{
					Actions: types.RolePolicyFields(http.MethodPost),
					Paths:   types.RolePolicyFields("/account/login"),
				},
				{
					Actions:   types.RolePolicyFields(http.MethodGet),
					Resources: types.RolePolicyFields("settings"),
					ObjectIDs: types.RolePolicyFields("FirstLogin"),
					Scope:     types.RolePolicyResourceScopeGlobal,
				},
			},
			Builtin: true,
			Session: true,
		},

		// System user.
		{
			Domain:      "system",
			Name:        "user",
			Description: "system/user",
			Policies: types.RolePolicies{
				{
					Actions:          types.RolePolicyFields("*"),
					Resources:        types.RolePolicyFields("*"),
					ResourceExcludes: types.RolePolicyFields("groups", "users", "roles", "settings", "tokens"),
					Scope:            types.RolePolicyResourceScopeInherit,
				},
				{
					Actions:   types.RolePolicyFields("*"),
					Resources: types.RolePolicyFields("tokens"),
					Scope:     types.RolePolicyResourceScopePrivate,
				},
				{
					Actions: types.RolePolicyFields("*"),
					Paths:   types.RolePolicyFields("/account/info"),
				},
				{
					Actions: types.RolePolicyFields(http.MethodPost),
					Paths:   types.RolePolicyFields("/account/logout"),
				},
			},
			Builtin: true,
			Session: true,
		},

		// System admin.
		{
			Domain:      "system",
			Name:        "admin",
			Description: "system/admin",
			Policies: types.RolePolicies{
				{
					Actions:          types.RolePolicyFields("*"),
					Resources:        types.RolePolicyFields("*"),
					ResourceExcludes: types.RolePolicyFields("tokens"),
					Scope:            types.RolePolicyResourceScopeGlobal,
				},
			},
			Builtin: true,
		},
	}

	creates, err := dao.RoleCreates(opts.ModelClient, builtin...)
	if err != nil {
		return err
	}
	for i := range creates {
		err = creates[i].
			OnConflict(
				sql.ConflictColumns(
					role.FieldDomain,
					role.FieldName,
				),
			).
			Update(func(upsert *model.RoleUpsert) {
				upsert.UpdatePolicies()
			}).
			Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
