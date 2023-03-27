package user

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"k8s.io/utils/pointer"

	"github.com/seal-io/seal/pkg/apis/user/view"
	"github.com/seal-io/seal/pkg/casdoor"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/settings"
)

func Handle(mc model.ClientSet) Handler {
	return Handler{
		modelClient: mc,
	}
}

type Handler struct {
	modelClient model.ClientSet
}

func (h Handler) Kind() string {
	return "User"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) error {
	var input = &model.Subject{
		Kind:        "user",
		Group:       req.Group,
		Name:        req.Name,
		Description: req.Description,
		MountTo:     pointer.Bool(false),
		LoginTo:     pointer.Bool(true),
		Roles:       req.Roles,
		Paths:       req.Paths,
		Builtin:     false,
	}

	return h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		var creates, err = dao.SubjectCreates(tx, input)
		if err != nil {
			return err
		}
		_, err = creates[0].Save(ctx)
		if err != nil {
			return err
		}
		// create user from casdoor.
		var cred casdoor.ApplicationCredential
		err = settings.CasdoorCred.ValueJSONUnmarshal(ctx, tx, &cred)
		if err != nil {
			return err
		}
		err = casdoor.CreateUser(ctx, cred.ClientID, cred.ClientSecret,
			casdoor.BuiltinApp, casdoor.BuiltinOrg, req.Name, req.Password)
		if err != nil {
			return fmt.Errorf("failed to create user to casdoor: %w", err)
		}
		return nil
	})
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	var input = []predicate.Subject{
		subject.Kind("user"),
		subject.Group(req.Group),
		subject.Name(req.Name),
	}
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		// TODO cascade delete token.
		var _, err = tx.Subjects().Delete().
			Where(input...).
			Exec(ctx)
		if err != nil {
			return err
		}
		switch {
		case !req.MountTo: // created user.
			var cred casdoor.ApplicationCredential
			err = settings.CasdoorCred.ValueJSONUnmarshal(ctx, tx, &cred)
			if err != nil {
				return err
			}
			err = casdoor.DeleteUser(ctx, cred.ClientID, cred.ClientSecret,
				casdoor.BuiltinOrg, req.Name)
			if err != nil {
				if !strings.HasSuffix(err.Error(), "not found") {
					return fmt.Errorf("failed to delete user from casdoor: %w", err)
				}
			}
			return nil
		case req.LoginTo: // mounted user but login on.
			return tx.Subjects().Update().
				SetLoginTo(true).
				Where(
					subject.Kind("user"),
					subject.Name(req.Name),
					subject.MountTo(false),
				).
				Exec(ctx)
		}
		return nil
	})
	// TODO clean cache
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	var input = &model.Subject{
		Kind:        "user",
		ID:          req.ID,
		Description: req.Description,
		Roles:       req.Roles,
	}

	return h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		var updates, err = dao.SubjectUpdates(tx, input)
		if err != nil {
			return err
		}
		err = updates[0].Exec(ctx)
		if err != nil {
			return err
		}
		// update password.
		if req.Password == "" {
			return nil
		}
		var cred casdoor.ApplicationCredential
		err = settings.CasdoorCred.ValueJSONUnmarshal(ctx, tx, &cred)
		if err != nil {
			return err
		}
		err = casdoor.UpdateUserPassword(ctx, cred.ClientID, cred.ClientSecret,
			casdoor.BuiltinOrg, req.Name, "", req.Password)
		if err != nil {
			return fmt.Errorf("failed to update user password to casdoor: %w", err)
		}
		return nil
	})
	// TODO clean cache
}

// Batch APIs

var (
	queryFields = []string{
		subject.FieldName,
	}
	getFields = subject.WithoutFields(
		subject.FieldCreateTime,
		subject.FieldLoginTo)
	sortFields = []string{
		subject.FieldCreateTime,
		subject.FieldUpdateTime}
)

func (h Handler) CollectionGet(ctx *gin.Context, req view.CollectionGetRequest) (view.CollectionGetResponse, int, error) {
	var input = []predicate.Subject{
		subject.Kind("user"),
	}
	if req.Group != "" {
		input = append(input, subject.Group(req.Group)) // include mounted user.
	} else {
		input = append(input, subject.MountTo(false)) // created user.
	}

	var query = h.modelClient.Subjects().Query().
		Where(input...)
	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	// get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// get entities.
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}
	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}
	if orders, ok := req.Sorting(sortFields); ok {
		query.Order(orders...)
	}
	entities, err := query.
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeSubjects(entities), cnt, nil
}

// Extensional APIs

func (h Handler) RouteMount(ctx *gin.Context, req view.RouteMountRequest) error {
	var input = &model.Subject{
		Kind:    "user",
		Group:   req.Group,
		Name:    req.Name,
		MountTo: pointer.Bool(true),
		LoginTo: pointer.Bool(false),
		Roles:   req.Roles,
		Paths:   req.Paths,
		Builtin: false,
	}

	var creates, err = dao.SubjectCreates(h.modelClient, input)
	if err != nil {
		return err
	}
	_, err = creates[0].Save(ctx)
	return err
}
