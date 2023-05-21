package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
)

func SubjectRoleRelationshipCreates(
	mc model.ClientSet,
	input ...*model.SubjectRoleRelationship,
) ([]*model.SubjectRoleRelationshipCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.SubjectRoleRelationshipCreate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.SubjectRoleRelationships().Create().
			SetSubjectID(r.SubjectID).
			SetRoleID(r.RoleID)

		// Optional.
		if r.ProjectID.IsNaive() {
			c.SetProjectID(r.ProjectID)
		}

		rrs[i] = c
	}

	return rrs, nil
}
