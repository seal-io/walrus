package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/status"
)

func TemplateCreates(mc model.ClientSet, input ...*model.Template) ([]*model.TemplateCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.TemplateCreate, len(input))

	for i := range input {
		r := input[i]
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.Templates().Create().
			SetID(r.ID).
			SetSource(r.Source).
			SetStatus(status.TemplateStatusInitializing)

		// Optional.
		c.SetDescription(r.Description)
		c.SetIcon(r.Icon)

		if r.Labels != nil {
			c.SetLabels(r.Labels)
		}
		rrs[i] = c
	}

	return rrs, nil
}

func TemplateUpdate(mc model.ClientSet, input *model.Template) (*model.TemplateUpdateOne, error) {
	if input == nil {
		return nil, errors.New("invalid input: nil entity")
	}

	// Predicated.
	if input.ID == "" {
		return nil, errors.New("invalid input: illegal predicates")
	}

	// Conditional.
	c := mc.Templates().UpdateOne(input).
		SetDescription(input.Description).
		SetIcon(input.Icon).
		SetStatus(input.Status).
		SetStatusMessage(input.StatusMessage)
	if input.Labels != nil {
		c.SetLabels(input.Labels)
	}

	if input.Source != "" {
		c.SetSource(input.Source)
	}

	return c, nil
}
