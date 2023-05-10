package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/setting"
)

func SettingCreates(mc model.ClientSet, input ...*model.Setting) ([]*model.SettingCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.SettingCreate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.Settings().Create().
			SetName(r.Name)

		// Optional.
		c.SetValue(r.Value)
		c.SetNillableHidden(r.Hidden)
		c.SetNillableEditable(r.Editable)
		c.SetNillablePrivate(r.Private)
		rrs[i] = c
	}
	return rrs, nil
}

func SettingUpdates(mc model.ClientSet, input ...*model.Setting) ([]*model.SettingUpdate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.SettingUpdate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Predicated.
		var ps []predicate.Setting
		switch {
		case r.ID.IsNaive():
			ps = append(ps, setting.ID(r.ID))
		case r.Name != "":
			ps = append(ps, setting.Name(r.Name))
		}
		if len(ps) == 0 {
			return nil, errors.New("invalid input: illegal predicates")
		}

		c := mc.Settings().Update().
			Where(ps...).
			SetValue(r.Value).
			SetNillableHidden(r.Hidden).
			SetNillableEditable(r.Editable).
			SetNillablePrivate(r.Private)
		rrs[i] = c
	}
	return rrs, nil
}
