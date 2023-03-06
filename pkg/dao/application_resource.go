package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
)

func ApplicationResourceCreates(mc model.ClientSet, input ...*model.ApplicationResource) ([]*model.ApplicationResourceCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var resources = make([]*model.ApplicationResourceCreate, len(input))
	for i := range input {
		r := input[i]
		// required.
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		c := mc.ApplicationResources().Create().
			SetApplicationID(r.ApplicationID).
			SetConnectorID(r.ConnectorID).
			SetName(r.Name).
			SetType(r.Type).
			SetModule(r.Module).
			SetMode(r.Mode).
			SetDeployerType(r.DeployerType)

		resources[i] = c
	}

	return resources, nil
}
