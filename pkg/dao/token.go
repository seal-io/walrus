package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
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
