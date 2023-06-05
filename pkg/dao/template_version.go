package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
)

func TemplateVersionCreates(
	mc model.ClientSet,
	input ...*model.TemplateVersion,
) ([]*model.TemplateVersionCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.TemplateVersionCreate, len(input))

	for i := range input {
		r := input[i]
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.TemplateVersions().Create().
			SetTemplateID(r.TemplateID).
			SetSource(r.Source).
			SetVersion(r.Version)

		// Optional.
		if r.Schema != nil {
			c.SetSchema(r.Schema)
		}
		rrs[i] = c
	}

	return rrs, nil
}
