package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
)

func ConnectorCreates(mc model.ClientSet, input ...*model.Connector) ([]*model.ConnectorCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*model.ConnectorCreate, len(input))
	for i := range input {
		r := input[i]
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// required.
		if r.Name == "" {
			return nil, errors.New("invalid input: blank name")
		}

		if r.ConfigVersion == "" {
			return nil, errors.New("invalid input: blank configVersion")
		}

		if r.ConfigData == nil {
			return nil, errors.New("invalid input: blank configData")
		}

		var c = mc.Connectors().Create().
			SetName(r.Name).
			SetType(r.Type).
			SetConfigVersion(r.ConfigVersion).
			SetConfigData(r.ConfigData).
			SetEnableFinOps(r.EnableFinOps)

		// optional.
		if r.Description != "" {
			c.SetDescription(r.Description)
		}

		rrs[i] = c
	}
	return rrs, nil
}

func ConnectorUpdate(mc model.ClientSet, r *model.Connector) (*model.ConnectorUpdateOne, error) {
	if r == nil {
		return nil, errors.New("invalid input: nil entity")
	}

	if r.ConfigVersion == "" {
		return nil, errors.New("invalid input: blank configVersion")
	}

	if r.ConfigData == nil {
		return nil, errors.New("invalid input: blank configData")
	}

	var c = mc.Connectors().UpdateOne(r).
		SetConfigVersion(r.ConfigVersion).
		SetConfigData(r.ConfigData).
		SetEnableFinOps(r.EnableFinOps)

	return c, nil
}
