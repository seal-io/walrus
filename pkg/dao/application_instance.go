package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/status"
)

func ApplicationInstanceCreates(
	mc model.ClientSet,
	input ...*model.ApplicationInstance,
) ([]*model.ApplicationInstanceCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.ApplicationInstanceCreate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.ApplicationInstances().Create().
			SetProjectID(r.ProjectID).
			SetApplicationID(r.ApplicationID).
			SetEnvironmentID(r.EnvironmentID).
			SetName(r.Name)

		// Optional.
		c.SetVariables(r.Variables)
		status.ApplicationInstanceStatusDeployed.Unknown(r, "Deploying instance")
		r.Status.SetSummary(status.WalkApplicationInstance(&r.Status))
		c.SetStatus(r.Status)

		rrs[i] = c
	}

	return rrs, nil
}

func ApplicationInstanceUpdate(
	mc model.ClientSet,
	input *model.ApplicationInstance,
) (*model.ApplicationInstanceUpdateOne, error) {
	if input == nil {
		return nil, errors.New("invalid input: nil entity")
	}

	if input.ID == "" {
		return nil, errors.New("invalid input: illegal predicates")
	}

	c := mc.ApplicationInstances().UpdateOne(input).
		SetVariables(input.Variables)

	input.Status.SetSummary(status.WalkApplicationInstance(&input.Status))

	if input.Status.Changed() {
		c.SetStatus(input.Status)
	}

	return c, nil
}

// ApplicationInstanceStatusUpdate updates the status of the given application instance.
// TODO (alex): unify the status update logic for all application instances.
func ApplicationInstanceStatusUpdate(
	mc model.ClientSet,
	input *model.ApplicationInstance,
) (*model.ApplicationInstanceUpdateOne, error) {
	if input == nil {
		return nil, errors.New("invalid input: nil entity")
	}

	if input.ID == "" {
		return nil, errors.New("invalid input: illegal predicates")
	}

	c := mc.ApplicationInstances().UpdateOne(input)
	input.Status.SetSummary(status.WalkApplicationInstance(&input.Status))

	if input.Status.Changed() {
		c.SetStatus(input.Status)
	}

	return c, nil
}
