package auth

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/schema"
)

// WithResourceRoleGenerator wraps the given gin.IRoutes
// to support generating resource roles.
func WithResourceRoleGenerator(ctx context.Context, r gin.IRouter, modelClient model.ClientSet) gin.IRouter {
	return generator{
		IRouter:     r,
		ctx:         ctx,
		modelClient: modelClient,
	}
}

type generator struct {
	gin.IRouter

	ctx         context.Context
	modelClient model.ClientSet
}

func (g generator) AfterAdvice(h runtime.AdviceResource) error {
	var resource, resourcePath = h.ResourceAndResourcePath()

	// NB(thxCode): do not generate role for the following resource,
	// as we already granted the related permission to "system/user" role.
	switch resource {
	case "", "tokens":
		return nil
	}

	var builtin = []*model.Role{
		// resource admin
		{
			Domain:      resourcePath,
			Name:        "admin",
			Description: resourcePath + "/admin",
			Policies: schema.RolePolicies{
				schema.RolePolicyResourceAdminFor(resource),
			},
			Builtin: true,
		},
		// resource edit
		{
			Domain:      resourcePath,
			Name:        "edit",
			Description: resourcePath + "/edit",
			Policies: schema.RolePolicies{
				schema.RolePolicyResourceEditFor(resource),
			},
			Builtin: true,
		},
		// resource view
		{
			Domain:      resourcePath,
			Name:        "view",
			Description: resourcePath + "/view",
			Policies: schema.RolePolicies{
				schema.RolePolicyResourceViewFor(resource),
			},
			Builtin: true,
		},
	}

	var creates, err = dao.RoleCreates(g.modelClient, builtin...)
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
			Exec(g.ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
