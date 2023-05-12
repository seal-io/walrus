package view

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Basic APIs.

type CreateRequest struct {
	Group       string `json:"group"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	Paths []string `json:"-"`
}

func (r *CreateRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if r.Group == "" {
		return errors.New("invalid group: blank")
	}

	if r.Name == "" {
		return errors.New("invalid name: blank")
	}

	group, err := modelClient.Subjects().Query().
		Where(
			subject.Kind("group"),
			subject.Name(r.Group),
		).
		Only(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get group")
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
	modelClient := input.(model.ClientSet)

	if !r.ID.Valid(2) {
		return errors.New("invalid id: blank")
	}

	confirmGroup := []predicate.Subject{
		subject.Kind("group"),
		subject.Builtin(false),
	}
	if r.ID.IsNaive() {
		confirmGroup = append(confirmGroup, subject.ID(r.ID))
	} else {
		keys := r.ID.Split()
		confirmGroup = append(confirmGroup, subject.Group(keys[0]))
		confirmGroup = append(confirmGroup, subject.Name(keys[1]))
	}

	groupEntity, err := modelClient.Subjects().Query().
		Where(confirmGroup...).
		Select(subject.FieldID, subject.FieldGroup, subject.FieldName).
		Only(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get group")
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
	modelClient := input.(model.ClientSet)

	if !r.ID.Valid(2) {
		return errors.New("invalid id: blank")
	}

	confirmGroup := []predicate.Subject{
		subject.Kind("group"),
	}
	if r.ID.IsNaive() {
		confirmGroup = append(confirmGroup, subject.ID(r.ID))
	} else {
		keys := r.ID.Split()
		confirmGroup = append(confirmGroup, subject.Group(keys[0]))
		confirmGroup = append(confirmGroup, subject.Name(keys[1]))
	}

	groupEntity, err := modelClient.Subjects().Query().
		Where(confirmGroup...).
		Select(subject.FieldID, subject.FieldGroup, subject.FieldName).
		Only(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get group")
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

// Batch APIs.

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.Subject, subject.OrderOption] `query:",inline"`

	Group string `query:"_group,omitempty"`
}

type CollectionGetResponse = []*model.SubjectOutput

// Extensional APIs.
