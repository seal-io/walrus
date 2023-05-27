package dao

import (
	"context"
	"errors"

	"entgo.io/ent/dialect/sql"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/model/subjectrolerelationship"
	"github.com/seal-io/seal/utils/strs"
)

// WrappedSubjectUpsertOne is a wrapper for model.SubjectUpsertOne
// to process the relationship with model.Role.
// TODO(thxCode): generate this with entc.
type WrappedSubjectUpsertOne struct {
	*WrappedSubjectCreate

	upsertOne *model.SubjectUpsertOne
}

func (suo *WrappedSubjectUpsertOne) DoNothing() *WrappedSubjectUpsertOne {
	suo.upsertOne.DoNothing()
	return suo
}

func (suo *WrappedSubjectUpsertOne) Ignore() *WrappedSubjectUpsertOne {
	suo.upsertOne.Ignore()
	return suo
}

func (suo *WrappedSubjectUpsertOne) UpdateNewValues() *WrappedSubjectUpsertOne {
	suo.upsertOne.UpdateNewValues()
	return suo
}

// WrappedSubjectCreate is a wrapper for model.SubjectCreate
// to process the relationship with model.Role.
// TODO(thxCode): generate this with entc.
type WrappedSubjectCreate struct {
	*model.SubjectCreate

	entity *model.Subject
}

func (sc *WrappedSubjectCreate) OnConflictColumns(cols ...string) *WrappedSubjectUpsertOne {
	return &WrappedSubjectUpsertOne{
		WrappedSubjectCreate: sc,
		upsertOne:            sc.SubjectCreate.OnConflictColumns(cols...),
	}
}

func (sc *WrappedSubjectCreate) Save(ctx context.Context) (created *model.Subject, err error) {
	mc := sc.SubjectCreate.Mutation().Client()

	// Save entity.
	created, err = sc.SubjectCreate.Save(ctx)
	if err != nil {
		return
	}

	// Construct relationships.
	newRss := sc.entity.Edges.Roles
	createRss := make([]*model.SubjectRoleRelationshipCreate, len(newRss))

	for i, rs := range newRss {
		if rs == nil {
			return nil, errors.New("invalid input: nil relationship")
		}

		// Required.
		c := mc.SubjectRoleRelationships().Create().
			SetSubjectID(created.ID).
			SetRoleID(rs.RoleID)

		createRss[i] = c
	}

	// Save relationships.
	newRss, err = mc.SubjectRoleRelationships().CreateBulk(createRss...).
		Save(ctx)
	if err != nil {
		return
	}
	created.Edges.Roles = newRss

	return created, nil
}

func (sc *WrappedSubjectCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

func SubjectCreates(mc model.ClientSet, input ...*model.Subject) ([]*WrappedSubjectCreate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*WrappedSubjectCreate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Required.
		c := mc.Subjects().Create().
			SetName(r.Name)

		// Optional.
		if r.Kind != "" {
			c.SetKind(r.Kind)
		}

		if r.Domain != "" {
			c.SetDomain(r.Domain)
		}

		c.SetDescription(r.Description)
		c.SetBuiltin(r.Builtin)

		rrs[i] = &WrappedSubjectCreate{
			SubjectCreate: c,
			entity:        input[i],
		}
	}

	return rrs, nil
}

// WrappedSubjectUpdate is a wrapper for model.SubjectUpdate
// to process the relationship with model.Role.
// TODO(thxCode): generate this with entc.
type WrappedSubjectUpdate struct {
	*model.SubjectUpdate

	entity           *model.Subject
	entityPredicates []predicate.Subject
}

func (su *WrappedSubjectUpdate) Save(ctx context.Context) (updated int, err error) {
	mc := su.SubjectUpdate.Mutation().Client()

	if len(su.SubjectUpdate.Mutation().Fields()) != 0 {
		// Update entity.
		updated, err = su.SubjectUpdate.Save(ctx)
		if err != nil {
			return
		}
	}

	// Get old relationships.
	oldEntity, err := mc.Subjects().Query().
		Where(su.entityPredicates...).
		Select(subject.FieldID).
		WithRoles(func(rq *model.SubjectRoleRelationshipQuery) {
			rq.Where(subjectrolerelationship.ProjectIDIsNil()).
				Select(
					subjectrolerelationship.FieldSubjectID,
					subjectrolerelationship.FieldRoleID,
				)
		}).
		Only(ctx)
	if err != nil {
		return
	}

	// Create new relationship or update relationship.
	subjectID := oldEntity.ID
	newRsKeys := sets.New[string]()
	newRss := su.entity.Edges.Roles

	for _, rs := range newRss {
		newRsKeys.Insert(strs.Join("/", string(subjectID), rs.RoleID))

		// Required.
		c := mc.SubjectRoleRelationships().Create().
			SetSubjectID(subjectID).
			SetRoleID(rs.RoleID)

		err = c.
			OnConflict(
				sql.ConflictColumns(
					subjectrolerelationship.FieldSubjectID,
					subjectrolerelationship.FieldRoleID,
				),
				sql.ConflictWhere(sql.P().
					IsNull(subjectrolerelationship.FieldProjectID)),
			).
			DoNothing().
			Exec(ctx)
		if err != nil {
			return
		}
	}

	// Delete stale relationship.
	oldRss := oldEntity.Edges.Roles
	for _, rs := range oldRss {
		if newRsKeys.Has(strs.Join("/", string(rs.SubjectID), rs.RoleID)) {
			continue
		}

		_, err = mc.SubjectRoleRelationships().Delete().
			Where(
				subjectrolerelationship.ProjectIDIsNil(),
				subjectrolerelationship.SubjectID(rs.SubjectID),
				subjectrolerelationship.RoleID(rs.RoleID),
			).
			Exec(ctx)
		if err != nil {
			return
		}
	}

	return updated, nil
}

func (su *WrappedSubjectUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

func SubjectUpdates(mc model.ClientSet, input ...*model.Subject) ([]*WrappedSubjectUpdate, error) {
	if len(input) == 0 {
		return nil, errors.New("invalid input: empty list")
	}

	rrs := make([]*WrappedSubjectUpdate, len(input))

	for i, r := range input {
		if r == nil {
			return nil, errors.New("invalid input: nil entity")
		}

		// Predicated.
		var ps []predicate.Subject

		switch {
		case r.ID.IsNaive():
			ps = append(ps, subject.ID(r.ID))
		case r.Kind != "" && r.Domain != "" && r.Name != "":
			ps = append(ps, subject.And(
				subject.Kind(r.Kind),
				subject.Domain(r.Domain),
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

		if r.Domain != "" {
			c.SetDomain(r.Domain)
		}

		rrs[i] = &WrappedSubjectUpdate{
			SubjectUpdate:    c,
			entity:           input[i],
			entityPredicates: ps,
		}
	}

	return rrs, nil
}
