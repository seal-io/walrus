package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/auth/cache"
	"github.com/seal-io/seal/pkg/apis/auth/session"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/schema"
)

func authz(c *gin.Context, modelClient model.ClientSet) error {
	var s = session.LoadSubject(c)

	var permission, cached = cache.LoadSubjectPermission(s.Key())
	if !cached {
		permission = &cache.SubjectPermission{}
		var err error
		var roles = schema.SubjectRoles{
			{
				Domain: "system",
				Name:   "anonymity",
			},
		}
		if !s.IsAnonymous() {
			permission.Roles, err = getRoles(c, modelClient, s.Group, s.Name)
			if err != nil {
				return err
			}
			roles = append(roles, schema.SubjectRole{
				Domain: "system",
				Name:   "user",
			})
			roles = append(roles, permission.Roles...)
		}
		permission.Policies, err = getPolicies(c, modelClient, roles)
		if err != nil {
			return err
		}
		// cache
		cache.StoreSubjectPermission(s.Key(), *permission)
	}

	// validate
	session.StoreSubjectAuthzInfo(c, permission.Roles, permission.Policies)
	return nil
}

func getRoles(ctx context.Context, modelClient model.ClientSet, group, name string) (schema.SubjectRoles, error) {
	var s, err = modelClient.Subjects().Query().
		Where(subject.And(
			subject.Kind("user"),
			subject.Group(group),
			subject.Name(name),
		)).
		Select(subject.FieldRoles).
		Only(ctx)
	if err != nil {
		return nil, runtime.ErrorfP(http.StatusInternalServerError, "failed to get roles: %w", err)
	}
	return s.Roles, nil
}

func getPolicies(ctx context.Context, modelClient model.ClientSet, roles schema.SubjectRoles) (schema.RolePolicies, error) {
	var predicates []predicate.Role
	for i := 0; i < len(roles); i++ {
		predicates = append(predicates,
			role.And(
				role.Domain(roles[i].Domain),
				role.Name(roles[i].Name),
			))
	}
	var entities, err = modelClient.Roles().Query().
		Where(role.Or(predicates...)).
		Select(role.FieldPolicies).
		All(ctx)
	if err != nil {
		return nil, runtime.ErrorfP(http.StatusInternalServerError, "failed to get policies: %w", err)
	}

	var policies schema.RolePolicies
	for i := 0; i < len(entities); i++ {
		policies = append(policies, entities[i].Policies...)
	}
	policies = policies.Deduplicate().Sort()
	return policies, nil
}
