package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/secret"
)

func SecretCreates(mc model.ClientSet, input ...*model.Secret) ([]*model.SecretCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*model.SecretCreate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		var c = mc.Secrets().Create().
			SetName(r.Name).
			SetValue(r.Value)

		// Optional.
		if r.ProjectID != "" {
			c.SetProjectID(r.ProjectID)
		}
		rrs[i] = c
	}
	return rrs, nil
}

func SecretUpdates(mc model.ClientSet, input ...*model.Secret) ([]*model.SecretUpdate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*model.SecretUpdate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Predicated.
		var ps []predicate.Secret
		switch {
		case r.ID.IsNaive():
			ps = append(ps, secret.ID(r.ID))
		case r.Name != "":
			ps = append(ps, secret.Name(r.Name))
			if r.ProjectID != "" {
				ps = append(ps, secret.ProjectID(r.ProjectID))
			} else {
				ps = append(ps, secret.ProjectIDIsNil())
			}
		}
		if len(ps) == 0 {
			return nil, errors.New("invalid input: illegal predicates")
		}

		var c = mc.Secrets().Update().
			Where(ps...).
			SetValue(r.Value)
		rrs[i] = c
	}
	return rrs, nil
}
