package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/token"
)

func TokenCreates(mc model.ClientSet, input ...*model.Token) ([]*model.TokenCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*model.TokenCreate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// required.
		if r.CasdoorTokenName == "" {
			return nil, errors.New("invalid input: blank casddor token name")
		}
		if r.CasdoorTokenOwner == "" {
			return nil, errors.New("invalid input: blank casddor token owner")
		}
		if r.Name == "" {
			return nil, errors.New("invalid input: blank name")
		}
		var c = mc.Tokens().Create().
			SetCasdoorTokenName(r.CasdoorTokenName).
			SetCasdoorTokenOwner(r.CasdoorTokenOwner).
			SetName(r.Name)

		// optional.
		c.SetNillableExpiration(r.Expiration)
		rrs[i] = c
	}
	return rrs, nil
}

func TokenDeletes(mc model.ClientSet, input ...*model.Token) ([]*model.TokenDelete, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*model.TokenDelete, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// predicated.
		var ps []predicate.Token
		switch {
		case r.ID.IsNaive():
			ps = append(ps, token.ID(r.ID))
		case r.CasdoorTokenName != "":
			ps = append(ps, token.CasdoorTokenName(r.CasdoorTokenName))
		}
		if len(ps) == 0 {
			return nil, errors.New("invalid input: illegal predicates")
		}
		rrs[i] = mc.Tokens().Delete().
			Where(ps...)
	}
	return rrs, nil
}
