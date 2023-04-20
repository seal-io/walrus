package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/strs"
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
		c.SetStatusMessage(strs.NormalizeSpecialChars(r.StatusMessage))
		c.SetVariables(r.Variables)

		rrs[i] = c
	}
	return rrs, nil
}

func ApplicationInstanceUpdate(mc model.ClientSet, input *model.ApplicationInstance) (*model.ApplicationInstanceUpdateOne, error) {
	if input == nil {
		return nil, errors.New("invalid input: nil entity")
	}

	if input.ID == "" {
		return nil, errors.New("invalid input: illegal predicates")
	}

	var c = mc.ApplicationInstances().UpdateOne(input).
		SetStatusMessage(strs.NormalizeSpecialChars(input.StatusMessage)).
		SetVariables(input.Variables)
	if input.Status != "" {
		c.SetStatus(input.Status)
	}

	return c, nil
}
