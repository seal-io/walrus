package subject

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/casdoor"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/role"
	"github.com/seal-io/walrus/pkg/dao/model/subject"
	"github.com/seal-io/walrus/pkg/dao/model/subjectrolerelationship"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/settings"
)

func (h Handler) Create(req CreateRequest) (CreateResponse, error) {
	entity := req.Model()

	err := h.modelClient.WithTx(req.Context, func(tx *model.Tx) (err error) {
		entity, err = tx.Subjects().Create().
			Set(entity).
			SaveE(req.Context, dao.SubjectRolesEdgeSave)
		if err != nil {
			return err
		}

		if entity.Kind != types.SubjectKindUser ||
			entity.Domain != types.SubjectDomainBuiltin {
			return nil
		}

		// Create user from casdoor.
		var cred casdoor.ApplicationCredential

		err = settings.CasdoorCred.ValueJSONUnmarshal(req.Context, tx, &cred)
		if err != nil {
			return err
		}

		err = casdoor.CreateUser(req.Context, cred.ClientID, cred.ClientSecret,
			casdoor.BuiltinApp, casdoor.BuiltinOrg, req.Name, req.Password)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return model.ExposeSubject(entity), nil
}

func (h Handler) Delete(req DeleteRequest) error {
	entity, err := h.modelClient.Subjects().Query().
		Where(subject.ID(req.ID)).
		Select(
			subject.FieldID,
			subject.FieldKind,
			subject.FieldDomain,
			subject.FieldName,
			subject.FieldBuiltin).
		Only(req.Context)
	if err != nil {
		return err
	}

	if entity.Builtin {
		return runtime.Error(http.StatusForbidden, "cannot delete builtin subject")
	}

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		err = tx.Subjects().DeleteOne(entity).
			Exec(req.Context)
		if err != nil {
			return err
		}

		if entity.Kind != types.SubjectKindUser ||
			entity.Domain != types.SubjectDomainBuiltin {
			return nil
		}

		// Delete user from casdoor.
		var cred casdoor.ApplicationCredential

		err = settings.CasdoorCred.ValueJSONUnmarshal(req.Context, tx, &cred)
		if err != nil {
			return err
		}

		err = casdoor.DeleteUser(req.Context, cred.ClientID, cred.ClientSecret,
			casdoor.BuiltinOrg, entity.Name)
		if err != nil {
			if !strings.HasSuffix(err.Error(), "not found") {
				return fmt.Errorf("failed to delete user from casdoor: %w", err)
			}
		}

		return nil
	})
}

func (h Handler) Update(req UpdateRequest) error {
	entity := req.Model()

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		err := tx.Subjects().UpdateOne(entity).
			Set(entity).
			ExecE(req.Context, dao.SubjectRolesEdgeSave)
		if err != nil {
			return err
		}

		if entity.Kind != types.SubjectKindUser ||
			entity.Domain != types.SubjectDomainBuiltin {
			return nil
		}

		// Update password.
		if req.Password == "" {
			return nil
		}

		var cred casdoor.ApplicationCredential

		err = settings.CasdoorCred.ValueJSONUnmarshal(req.Context, tx, &cred)
		if err != nil {
			return err
		}

		err = casdoor.UpdateUserPassword(req.Context, cred.ClientID, cred.ClientSecret,
			casdoor.BuiltinOrg, entity.Name, "", req.Password)
		if err != nil {
			return fmt.Errorf("failed to update user password to casdoor: %w", err)
		}

		return nil
	})
}

var (
	queryFields = []string{
		subject.FieldDomain,
		subject.FieldName,
	}
	getFields = subject.WithoutFields(
		subject.FieldUpdateTime)
	sortFields = []string{
		subject.FieldID,
		subject.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.Subjects().Query()

	if req.Kind != "" {
		query.Where(subject.Kind(req.Kind))
	}

	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
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

	if orders, ok := req.Sorting(sortFields, model.Desc(subject.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		Unique(false).
		WithRoles(func(rq *model.SubjectRoleRelationshipQuery) {
			rq.Order(model.Desc(subjectrolerelationship.FieldCreateTime)).
				Where(subjectrolerelationship.ProjectIDIsNil()).
				Select(subjectrolerelationship.FieldRoleID).
				Unique(false).
				WithRole(func(rq *model.RoleQuery) {
					rq.Select(role.FieldID)
				})
		}).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeSubjects(entities), cnt, nil
}
