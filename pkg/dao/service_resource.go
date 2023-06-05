package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
)

func ServiceResourceCreates(
	mc model.ClientSet,
	input ...*model.ServiceResource,
) ([]*model.ServiceResourceCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.ServiceResourceCreate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.ServiceResources().Create().
			SetProjectID(r.ProjectID).
			SetServiceID(r.ServiceID).
			SetConnectorID(r.ConnectorID).
			SetName(r.Name).
			SetType(r.Type).
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
