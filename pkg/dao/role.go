package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/types"
)

func RoleCreates(mc model.ClientSet, input ...*model.Role) ([]*model.RoleCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.RoleCreate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.Roles().Create().
			SetID(r.ID)

		// Optional.
		if r.Kind != "" {
			c.SetKind(r.Kind)
		}

		c.SetDescription(r.Description)

		if len(r.Policies) != 0 {
			c.SetPolicies(r.Policies.Normalize().Deduplicate().Sort())
		}

		c.SetBuiltin(r.Builtin)
		c.SetSession(r.Session)

		rrs[i] = c
	}

	return rrs, nil
}

func RoleUpdates(mc model.ClientSet, input ...*model.Role) ([]*model.RoleUpdate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.RoleUpdate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Predicated.
		if r.ID == "" {
			return nil, errors.New("invalid input: illegal predicates")
		}

		// Conditional.
		c := mc.Roles().Update().
			Where(role.ID(r.ID)).
			SetDescription(r.Description)
		if len(r.Policies) != 0 {
			c.SetPolicies(r.Policies.Normalize().Deduplicate().Sort())
		} else {
			c.SetPolicies(types.DefaultRolePolicies())
		}

		rrs[i] = c
	}

	return rrs, nil
}
