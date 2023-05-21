package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Basic APIs.

type CreateRequest struct {
	*model.SubjectCreateInput `json:",inline"`

	Password string `json:"password"`
}

func (r *CreateRequest) Validate() error {
	if !types.IsSubjectKind(r.Kind) {
		return errors.New("invalid kind: unknown")
	}

	if r.Name == "" {
		return errors.New("invalid name: blank")
	}

	if r.Password == "" {
		return errors.New("invalid password: blank")
	}

	r.Builtin = false

	return nil
}

type CreateResponse = *model.SubjectOutput

type DeleteRequest struct {
	*model.SubjectQueryInput `uri:",inline"`
}

func (r *DeleteRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type UpdateRequest struct {
	*model.SubjectUpdateInput `uri:",inline" json:",inline"`

	Password string `json:"password,omitempty"` // Allow to reset password.
}

func (r *UpdateRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

// Batch APIs.

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.Subject, subject.OrderOption] `query:",inline"`

	Kind string `query:"kind,omitempty"`
}

func (r *CollectionGetRequest) Validate() error {
	if r.Kind != "" && !types.IsSubjectKind(r.Kind) {
		return errors.New("invalid kind: unknown")
	}

	return nil
}

type CollectionGetResponse = []*model.SubjectOutput

// Extensional APIs.
