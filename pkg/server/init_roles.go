package server

import (
	"context"
	"net/http"

	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/schema"
)

func (r *Server) initRoles(ctx context.Context, opts initOptions) error {
	var builtin = []*model.Role{
		// system anonymity
		{
			Domain:      "system",
			Name:        "anonymity",
			Description: "system/anonymity",
			Policies: schema.RolePolicies{
				{
					Actions: schema.RolePolicyFields(http.MethodPost),
					Paths:   schema.RolePolicyFields("/account/login"),
				},
				{
					Actions:   schema.RolePolicyFields(http.MethodGet),
					Scope:     schema.RolePolicyResourceScopeGlobal,
					Resources: schema.RolePolicyFields("settings"),
					ObjectIDs: schema.RolePolicyFields("FirstLogin"),
				},
			},
			Builtin: true,
			Session: true,
		},

		// system user
		{
			Domain:      "system",
			Name:        "user",
			Description: "system/user",
			Policies: schema.RolePolicies{
				{
					Actions:          schema.RolePolicyFields("*"),
					Scope:            schema.RolePolicyResourceScopeInherit,
					Resources:        schema.RolePolicyFields("*"),
					ResourceExcludes: schema.RolePolicyFields("groups", "users", "roles", "settings", "tokens"),
				},
				{
					Actions:   schema.RolePolicyFields("*"),
					Scope:     schema.RolePolicyResourceScopePrivate,
					Resources: schema.RolePolicyFields("tokens"),
				},
				{
					Actions: schema.RolePolicyFields("*"),
					Paths:   schema.RolePolicyFields("/account/info"),
				},
				{
					Actions: schema.RolePolicyFields(http.MethodPost),
					Paths:   schema.RolePolicyFields("/account/logout"),
				},
			},
			Builtin: true,
			Session: true,
		},

		// system admin
		{
			Domain:      "system",
			Name:        "admin",
			Description: "system/admin",
			Policies: schema.RolePolicies{
				schema.RolePolicyResourceAdminFor("*"),
			},
			Builtin: true,
		},
	}

	var creates, err = dao.RoleCreates(opts.ModelClient, builtin...)
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
