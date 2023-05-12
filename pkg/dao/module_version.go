package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
)

func ModuleVersionCreates(mc model.ClientSet, input ...*model.ModuleVersion) ([]*model.ModuleVersionCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*model.ModuleVersionCreate, len(input))
	for i := range input {
		r := input[i]
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		var c = mc.ModuleVersions().Create().
			SetModuleID(r.ModuleID).
			SetSource(r.Source).
			SetVersion(r.Version)

		// Optional.
		if r.Schema != nil {
			c.SetSchema(r.Schema)
		}
		rrs[i] = c
	}
	return rrs, nil
}
