package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
)

func ApplicationRevisionCreates(mc model.ClientSet, input ...*model.ApplicationRevision) ([]*model.ApplicationRevisionCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	var rrs = make([]*model.ApplicationRevisionCreate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// required.
		var c = mc.ApplicationRevisions().Create().
			SetInstanceID(r.InstanceID).
			SetEnvironmentID(r.EnvironmentID).
			SetInputPlan(r.InputPlan).
			SetOutput(r.Output)

		// optional.
		c.SetStatus(r.Status)
		c.SetStatusMessage(r.StatusMessage)
		if r.Modules != nil {
			c.SetModules(r.Modules)
		}
		if r.InputVariables != nil {
			c.SetInputVariables(r.InputVariables)
		}
		if r.DeployerType != "" {
			c.SetDeployerType(r.DeployerType)
		}
		if r.Duration != 0 {
			c.SetDuration(r.Duration)
		}
		if len(r.PreviousRequiredProviders) != 0 {
			c.SetPreviousRequiredProviders(r.PreviousRequiredProviders)
		}

		rrs[i] = c
	}
	return rrs, nil
}
