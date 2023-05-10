package auth

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/auth/cache"
	"github.com/seal-io/seal/pkg/apis/auth/session"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/types"
)

func authz(c *gin.Context, modelClient model.ClientSet) error {
	s := session.LoadSubject(c)

	permission, cached := cache.LoadSubjectPermission(s.Key())
	if !cached {
		permission = &cache.SubjectPermission{}
		var err error
		roles := types.SubjectRoles{
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
			roles = append(roles, types.SubjectRole{
				Domain: "system",
				Name:   "user",
			})
			roles = append(roles, permission.Roles...)
		}
		permission.Policies, err = getPolicies(c, modelClient, roles)
		if err != nil {
			return err
		}
		// Cache.
		cache.StoreSubjectPermission(s.Key(), *permission)
	}

	// Validate.
	session.StoreSubjectAuthzInfo(c, permission.Roles, permission.Policies)
	return nil
}

func getRoles(
	ctx context.Context,
	modelClient model.ClientSet,
	group, name string,
) (types.SubjectRoles, error) {
	s, err := modelClient.Subjects().Query().
		Where(subject.And(
			subject.Kind("user"),
			subject.Group(group),
			subject.Name(name),
		)).
		Select(subject.FieldRoles).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	return s.Roles, nil
}

func getPolicies(
	ctx context.Context,
	modelClient model.ClientSet,
	roles types.SubjectRoles,
) (types.RolePolicies, error) {
	var predicates []predicate.Role
	for i := 0; i < len(roles); i++ {
		predicates = append(predicates,
			role.And(
				role.Domain(roles[i].Domain),
				role.Name(roles[i].Name),
			))
	}
	entities, err := modelClient.Roles().Query().
		Where(role.Or(predicates...)).
		Select(role.FieldPolicies).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get policies: %w", err)
	}

	var policies types.RolePolicies
	for i := 0; i < len(entities); i++ {
		policies = append(policies, entities[i].Policies...)
	}
	policies = policies.Deduplicate().Sort()
	return policies, nil
}
