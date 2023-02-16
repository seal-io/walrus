package view

import (
	"context"
	"errors"
	"net/http"
	"reflect"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/schema"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Basic APIs

type CreateRequest struct {
	Group       string              `json:"group"`
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	Roles       schema.SubjectRoles `json:"roles,omitempty"`
	Password    string              `json:"password"`

	Paths []string `json:"-"`
}

func (r *CreateRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if r.Group == "" {
		return errors.New("invalid group: blank")
	}
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}
	r.Roles = r.Roles.Deduplicate().Sort()
	for i := range r.Roles {
		if r.Roles[i].IsZero() {
			return errors.New("invalid role")
		}
	}
	if r.Password == "" {
		return errors.New("invalid password: blank")
	}

	var group, err = modelClient.Subjects().Query().
		Where(
			subject.Kind("group"),
			subject.Name(r.Group),
		).
		Only(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return runtime.Error(http.StatusBadRequest, "invalid group: not found")
		}
		return runtime.ErrorfP(http.StatusInternalServerError, "failed to get requesting group: %w", err)
	}
	r.Paths = group.Paths
	r.Paths = append(r.Paths, r.Name)
	return nil
}

type DeleteRequest struct {
	ID types.ID `uri:"id"`

	Group   string `json:"-"`
	Name    string `json:"-"`
	MountTo bool   `json:"-"`
	LoginTo bool   `json:"-"`
}

func (r *DeleteRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(2) {
		return errors.New("invalid id: blank")
	}

	var confirmUser = []predicate.Subject{
		subject.Kind("user"),
		subject.Builtin(false),
	}
	if r.ID.IsNaive() {
		confirmUser = append(confirmUser, subject.ID(r.ID))
	} else {
		var keys = r.ID.Split()
		confirmUser = append(confirmUser, subject.Group(keys[0]))
		confirmUser = append(confirmUser, subject.Name(keys[1]))
	}
	var userEntity, err = modelClient.Subjects().Query().
		Where(confirmUser...).
		Select(subject.FieldID, subject.FieldGroup, subject.FieldName, subject.FieldMountTo, subject.FieldLoginTo).
		Only(ctx)
	if err != nil {
		return err
	}
	r.Name = userEntity.Name
	r.Group = userEntity.Group
	r.MountTo = *userEntity.MountTo
	r.LoginTo = *userEntity.LoginTo
	if !r.MountTo {
		r.Group = ""
	}
	return nil
}

type UpdateRequest struct {
	ID          types.ID            `uri:"id"`
	Description string              `json:"description,omitempty"`
	Roles       schema.SubjectRoles `json:"roles,omitempty"`
	Password    string              `json:"password,omitempty"`

	Name string `json:"-"`
}

func (r *UpdateRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(2) {
		return errors.New("invalid id: blank")
	}

	var confirmUser = []predicate.Subject{
		subject.Kind("user"),
	}
	if r.ID.IsNaive() {
		confirmUser = append(confirmUser, subject.ID(r.ID))
	} else {
		var keys = r.ID.Split()
		confirmUser = append(confirmUser, subject.Group(keys[0]))
		confirmUser = append(confirmUser, subject.Name(keys[1]))
	}
	var userEntity, err = modelClient.Subjects().Query().
		Where(confirmUser...).
		Select(subject.FieldID, subject.FieldGroup, subject.FieldName, subject.FieldMountTo).
		Only(ctx)
	if err != nil {
		return err
	}
	r.ID = userEntity.ID
	r.Name = userEntity.Name

	var needUpdate bool
	if r.Description != "" {
		if *userEntity.MountTo {
			return errors.New("invalid user: mounting user cannot update description")
		}
		if r.Description == userEntity.Description {
			r.Description = ""
		} else {
			needUpdate = true
		}
	}
	r.Roles = r.Roles.Deduplicate().Sort()
	if len(r.Roles) != 0 {
		for i := range r.Roles {
			if r.Roles[i].IsZero() {
				return errors.New("invalid role")
			}
		}
		if reflect.DeepEqual(r.Roles, userEntity.Roles) {
			r.Roles = nil
		} else {
			needUpdate = true
		}
	}
	if r.Password != "" {
		if *userEntity.MountTo {
			return errors.New("invalid user: mounting user cannot update password")
		}
		needUpdate = true
	}
	if !needUpdate {
		return errors.New("invalid input: nothing update")
	}

	return nil
}

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestExtracting `query:",inline"`
	runtime.RequestSorting    `query:",inline"`

	Group string `query:"_group,omitempty"`
}

type CollectionGetResponse = []*model.Subject

// Extensional APIs

type RouteMountRequest struct {
	_ struct{} `route:"POST=/mount"`

	ID    types.ID            `uri:"id"`
	Group string              `json:"group"`
	Roles schema.SubjectRoles `json:"roles,omitempty"`

	Name  string   `json:"-"`
	Paths []string `json:"-"`
}

func (r *RouteMountRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(2) {
		return errors.New("invalid id: blank")
	}
	if r.Group == "" {
		return errors.New("invalid group: blank")
	}
	r.Roles = r.Roles.Deduplicate().Sort()
	for i := range r.Roles {
		if r.Roles[i].IsZero() {
			return errors.New("invalid role")
		}
	}

	var confirmUser = []predicate.Subject{
		subject.Kind("user"),
		subject.Builtin(false),
	}
	if r.ID.IsNaive() {
		confirmUser = append(confirmUser, subject.ID(r.ID))
	} else {
		var keys = r.ID.Split()
		confirmUser = append(confirmUser, subject.Group(keys[0]))
		confirmUser = append(confirmUser, subject.Name(keys[1]))
	}
	var userEntity, err = modelClient.Subjects().Query().
		Where(confirmUser...).
		Select(subject.FieldGroup, subject.FieldName, subject.FieldMountTo).
		Only(ctx)
	if err != nil {
		return err
	}
	if *userEntity.MountTo {
		return runtime.Error(http.StatusBadRequest, "invalid user: already mounting")
	} else if userEntity.Group == r.Group {
		return runtime.Error(http.StatusBadRequest, "invalid group: the same")
	}
	r.Name = userEntity.Name

	groupEntity, err := modelClient.Subjects().Query().
		Where(
			subject.Kind("group"),
			subject.Name(r.Group),
		).
		Only(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return runtime.Error(http.StatusBadRequest, "invalid group: not found")
		}
		return runtime.ErrorfP(http.StatusInternalServerError, "failed to get requesting group: %w", err)
	}
	r.Paths = groupEntity.Paths
	r.Paths = append(r.Paths, r.Name)
	return nil
}
