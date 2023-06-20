package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/status"
)

func ServiceCreates(
	mc model.ClientSet,
	input ...*model.Service,
) ([]*model.ServiceCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.ServiceCreate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.Services().Create().
			SetProjectID(r.ProjectID).
			SetEnvironmentID(r.EnvironmentID).
			SetTemplate(r.Template).
			SetName(r.Name)

		// Optional.
		c.SetDescription(r.Description)

		if r.Labels != nil {
			c.SetLabels(r.Labels)
		}

		if r.Attributes != nil {
			c.SetAttributes(r.Attributes)
		}

		status.ServiceStatusDeployed.Unknown(r, "Deploying service")
		r.Status.SetSummary(status.WalkService(&r.Status))
		c.SetStatus(r.Status)

		rrs[i] = c
	}

	return rrs, nil
}

func ServiceUpdate(
	mc model.ClientSet,
	input *model.Service,
) (*model.ServiceUpdateOne, error) {
	if input == nil {
		return nil, errors.New("invalid input: nil entity")
	}

	if input.ID == "" {
		return nil, errors.New("invalid input: illegal predicates")
	}

	c := mc.Services().UpdateOne(input).
		SetTemplate(input.Template)

	c.SetDescription(input.Description)

	if input.Labels != nil {
		c.SetLabels(input.Labels)
	}

	if input.Attributes != nil {
		c.SetAttributes(input.Attributes)
	}

	input.Status.SetSummary(status.WalkService(&input.Status))

	if input.Status.Changed() {
		c.SetStatus(input.Status)
	}

	return c, nil
}

// ServiceStatusUpdate updates the status of the given service.
// TODO (alex): unify the status update logic for all services.
func ServiceStatusUpdate(
	mc model.ClientSet,
	input *model.Service,
) (*model.ServiceUpdateOne, error) {
	if input == nil {
		return nil, errors.New("invalid input: nil entity")
	}

	if input.ID == "" {
		return nil, errors.New("invalid input: illegal predicates")
	}

	c := mc.Services().UpdateOne(input)
	input.Status.SetSummary(status.WalkService(&input.Status))

	if input.Status.Changed() {
		c.SetStatus(input.Status)
	}

	return c, nil
}
