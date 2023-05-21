package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
)

func TokenCreates(mc model.ClientSet, input ...*model.Token) ([]*model.TokenCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.TokenCreate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.Tokens().Create().
			SetSubjectID(r.SubjectID).
			SetName(r.Name).
			SetValue(r.Value)

		// Optional.
		if r.Kind != "" {
			c.SetKind(r.Kind)
		}

		c.SetNillableExpiration(r.Expiration)

		rrs[i] = c
	}

	return rrs, nil
}
