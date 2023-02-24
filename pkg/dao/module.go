package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/status"
)

func ModuleCreates(mc model.ClientSet, input ...*model.Module) ([]*model.ModuleCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*model.ModuleCreate, len(input))
	for i := range input {
		r := input[i]
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// required.
		var c = mc.Modules().Create().
			SetID(r.ID).
			SetSource(r.Source).
			SetStatus(status.Initializing)

		// optional.
		if r.Description != "" {
			c.SetDescription(r.Description)
		}
		if r.Labels != nil {
			c.SetLabels(r.Labels)
		}
		if r.Version != "" {
			c.SetVersion(r.Version)
		}
		rrs[i] = c
	}
	return rrs, nil
}

func ModuleUpdate(mc model.ClientSet, m *model.Module) (*model.ModuleUpdateOne, error) {
	if m == nil {
		return nil, errors.New("invalid input: nil entity")
	}

	if m.Source == "" {
		return nil, errors.New("invalid input: blank source")
	}

	var c = mc.Modules().UpdateOne(m).
		SetSource(m.Source).
		SetDescription(m.Description).
		SetVersion(m.Version).
		SetSchema(m.Schema).
		SetStatus(m.Status).
		SetStatusMessage(m.StatusMessage)

	return c, nil
}
