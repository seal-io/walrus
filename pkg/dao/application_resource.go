package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
)

func ApplicationResourceCreates(
	mc model.ClientSet,
	input ...*model.ApplicationResource,
) ([]*model.ApplicationResourceCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.ApplicationResourceCreate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.ApplicationResources().Create().
			SetInstanceID(r.InstanceID).
			SetConnectorID(r.ConnectorID).
			SetName(r.Name).
			SetType(r.Type).
			SetModule(r.Module).
			SetMode(r.Mode).
			SetDeployerType(r.DeployerType)

		// Optional.
		if r.CompositionID.Valid(0) {
			c.SetCompositionID(r.CompositionID)
		}

		rrs[i] = c
	}

	return rrs, nil
}
