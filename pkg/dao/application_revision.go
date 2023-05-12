package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/strs"
)

func ApplicationRevisionCreates(
	mc model.ClientSet,
	input ...*model.ApplicationRevision,
) ([]*model.ApplicationRevisionCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.ApplicationRevisionCreate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.ApplicationRevisions().Create().
			SetInstanceID(r.InstanceID).
			SetEnvironmentID(r.EnvironmentID).
			SetInputPlan(r.InputPlan).
			SetOutput(r.Output)

		// Optional.
		c.SetStatus(r.Status)
		c.SetStatusMessage(strs.NormalizeSpecialChars(r.StatusMessage))

		if r.Modules != nil {
			c.SetModules(r.Modules)
		}

		if r.Secrets != nil {
			c.SetSecrets(r.Secrets)
		}

		if r.Variables != nil {
			c.SetVariables(r.Variables)
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

func ApplicationRevisionUpdate(
	mc model.ClientSet,
	input *model.ApplicationRevision,
) (*model.ApplicationRevisionUpdateOne, error) {
	if input == nil {
		return nil, errors.New("invalid input: nil entity")
	}

	if input.ID == "" {
		return nil, errors.New("invalid input: illegal predicates")
	}

	c := mc.ApplicationRevisions().UpdateOne(input).
		SetStatusMessage(strs.NormalizeSpecialChars(input.StatusMessage))
	if input.Status != "" {
		c.SetStatus(input.Status)
	}

	if input.InputPlan != "" {
		c.SetInputPlan(input.InputPlan)
	}

	if input.Output != "" {
		c.SetOutput(input.Output)
	}

	if input.Duration != 0 {
		c.SetDuration(input.Duration)
	}

	if input.Secrets != nil {
		c.SetSecrets(input.Secrets)
	}

	return c, nil
}
