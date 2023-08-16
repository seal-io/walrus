package projectsubject

import (
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/subject"
	"github.com/seal-io/walrus/pkg/dao/model/subjectrolerelationship"
)

func (h Handler) Delete(req DeleteRequest) error {
	return h.modelClient.SubjectRoleRelationships().DeleteOneID(req.ID).
		Exec(req.Context)
}

var (
	queryFields = []string{
		subject.FieldName,
	}
	getFields  = subjectrolerelationship.WithoutFields()
	sortFields = []string{
		subjectrolerelationship.FieldCreateTime,
	}
)

func (h Handler) CollectionCreate(req CollectionCreateRequest) (CollectionCreateResponse, error) {
	entities := req.Model()

	entities, err := h.modelClient.SubjectRoleRelationships().CreateBulk().
		Set(entities...).
		Save(req.Context)
	if err != nil {
		return nil, err
	}

	return model.ExposeSubjectRoleRelationships(entities), nil
}

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.SubjectRoleRelationships().Query()

	if req.Project != nil {
		// Project scope only.
		query.Where(subjectrolerelationship.ProjectID(req.Project.ID))
	} else {
		// Global scope.
		query.Where(subjectrolerelationship.ProjectIDIsNil())
	}

	if queries, ok := req.Querying(queryFields); ok {
		query.Where(subjectrolerelationship.HasSubjectWith(predicate.Subject(queries)))
	}

	// Get count.
	cnt, err := query.Clone().Count(req.Context)
	if err != nil {
		return nil, 0, err
	}

	// Get entities.
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}

	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}

	if orders, ok := req.Sorting(sortFields, model.Desc(subjectrolerelationship.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		// Must append project ID.
		Select(subjectrolerelationship.FieldProjectID).
		// Must extract subject.
		Select(subjectrolerelationship.FieldSubjectID).
		WithSubject(func(sq *model.SubjectQuery) {
			sq.Select(
				subject.FieldID,
				subject.FieldKind,
				subject.FieldDomain,
				subject.FieldName)
		}).
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeSubjectRoleRelationships(entities), cnt, nil
}

func (h Handler) CollectionDelete(req CollectionDeleteRequest) error {
	ids := req.IDs()

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		_, err := tx.SubjectRoleRelationships().Delete().
			Where(subjectrolerelationship.IDIn(ids...)).
			Exec(req.Context)

		return err
	})
}
