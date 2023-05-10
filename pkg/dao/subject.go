package dao

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/subject"
)

func SubjectCreates(mc model.ClientSet, input ...*model.Subject) ([]*model.SubjectCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.SubjectCreate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		if len(r.Paths) == 0 {
			return nil, errors.New("invalid input: empty paths")
		}
		c := mc.Subjects().Create().
			SetName(r.Name).
			SetPaths(r.Paths)

		// Optional.
		c.SetDescription(r.Description)
		c.SetNillableMountTo(r.MountTo)
		c.SetNillableLoginTo(r.LoginTo)
		if r.Kind != "" {
			c.SetKind(r.Kind)
		}
		if r.Group != "" {
			c.SetGroup(r.Group)
		}
		if len(r.Roles) != 0 {
			c.SetRoles(r.Roles.Deduplicate().Sort())
		}
		rrs[i] = c
	}
	return rrs, nil
}

func SubjectUpdates(mc model.ClientSet, input ...*model.Subject) ([]*model.SubjectUpdate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*model.SubjectUpdate, len(input))
	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Predicated.
		var ps []predicate.Subject
		switch {
		case r.ID.IsNaive():
			ps = append(ps, subject.ID(r.ID))
		case r.Kind != "" && r.Group != "" && r.Name != "":
			ps = append(ps, subject.And(
				subject.Kind(r.Kind),
				subject.Group(r.Group),
				subject.Name(r.Name),
			))
		}
		if len(ps) == 0 {
			return nil, errors.New("invalid input: illegal predicates")
		}

		// Conditional.
		c := mc.Subjects().Update().
			Where(ps...).
			SetDescription(r.Description)
		if r.Group != "" {
			c.SetGroup(r.Group)
		}
		if r.LoginTo != nil {
			c.SetLoginTo(*r.LoginTo)
		}
		if len(r.Roles) != 0 {
			c.SetRoles(r.Roles.Deduplicate().Sort())
		}
		if len(r.Paths) != 0 {
			c.SetPaths(r.Paths)
		}
		rrs[i] = c
	}
	return rrs, nil
}
