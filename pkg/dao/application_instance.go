package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
)

func ApplicationInstanceCreates(mc model.ClientSet, input ...*model.ApplicationInstance) ([]*model.ApplicationInstanceCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*model.ApplicationInstanceCreate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// required.
		var c = mc.ApplicationInstances().Create().
			SetApplicationID(r.ApplicationID).
			SetEnvironmentID(r.EnvironmentID).
			SetName(r.Name)

		// optional.
		c.SetStatus(r.Status)
		c.SetStatusMessage(r.StatusMessage)
		c.SetVariables(r.Variables)

		rrs[i] = c
	}
	return rrs, nil
}
