package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/strs"
)

func ServiceRevisionCreates(
	mc model.ClientSet,
	input ...*model.ServiceRevision,
) ([]*model.ServiceRevisionCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.ServiceRevisionCreate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.ServiceRevisions().Create().
			SetProjectID(r.ProjectID).
			SetServiceID(r.ServiceID).
			SetEnvironmentID(r.EnvironmentID).
			SetTemplateID(r.TemplateID).
			SetTemplateVersion(r.TemplateVersion).
			SetAttributes(r.Attributes).
			SetInputPlan(r.InputPlan).
			SetOutput(r.Output)

		// Optional.
		c.SetStatus(r.Status)
		c.SetStatusMessage(strs.NormalizeSpecialChars(r.StatusMessage))

		if r.Tags != nil {
			c.SetTags(r.Tags)
		}

		if r.Variables != nil {
			c.SetVariables(r.Variables)
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

func ServiceRevisionUpdate(
	mc model.ClientSet,
	input *model.ServiceRevision,
) (*model.ServiceRevisionUpdateOne, error) {
	if input == nil {
		return nil, errors.New("invalid input: nil entity")
	}

	if input.ID == "" {
		return nil, errors.New("invalid input: illegal predicates")
	}

	c := mc.ServiceRevisions().UpdateOne(input).
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

	if input.Variables != nil {
		c.SetVariables(input.Variables)
	}

	return c, nil
}
