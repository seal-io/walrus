package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
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
		c.SetNillableSensitive(r.Sensitive)
		c.SetNillablePrivate(r.Private)
		rrs[i] = c
	}

	return rrs, nil
}
