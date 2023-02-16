package view

import (
	"context"
	"errors"
	"net/http"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Basic APIs

type CreateRequest struct {
	Group       string `json:"group"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

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

	Name string `json:"-"`
}

func (r *DeleteRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(2) {
		return errors.New("invalid id: blank")
	}

	var confirmGroup = []predicate.Subject{
		subject.Kind("group"),
		subject.Builtin(false),
	}
	if r.ID.IsNaive() {
		confirmGroup = append(confirmGroup, subject.ID(r.ID))
	} else {
		var keys = r.ID.Split()
		confirmGroup = append(confirmGroup, subject.Group(keys[0]))
		confirmGroup = append(confirmGroup, subject.Name(keys[1]))
	}
	var groupEntity, err = modelClient.Subjects().Query().
		Where(confirmGroup...).
		Select(subject.FieldID, subject.FieldGroup, subject.FieldName).
		Only(ctx)
	if err != nil {
		return err
	}
	r.ID = groupEntity.ID
	r.Name = groupEntity.Name

	return nil
}

type UpdateRequest struct {
	ID          types.ID `uri:"id"`
	Description string   `json:"description,omitempty"`

	Name string `json:"-"`
}

func (r *UpdateRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(2) {
		return errors.New("invalid id: blank")
	}

	var confirmGroup = []predicate.Subject{
		subject.Kind("group"),
	}
	if r.ID.IsNaive() {
		confirmGroup = append(confirmGroup, subject.ID(r.ID))
	} else {
		var keys = r.ID.Split()
		confirmGroup = append(confirmGroup, subject.Group(keys[0]))
		confirmGroup = append(confirmGroup, subject.Name(keys[1]))
	}
	var groupEntity, err = modelClient.Subjects().Query().
		Where(confirmGroup...).
		Select(subject.FieldID, subject.FieldGroup, subject.FieldName).
		Only(ctx)
	if err != nil {
		return err
	}
	r.ID = groupEntity.ID
	r.Name = groupEntity.Name

	var needUpdate bool
	if r.Description != "" {
		if r.Description == groupEntity.Description {
			r.Description = ""
		} else {
			needUpdate = true
		}
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
