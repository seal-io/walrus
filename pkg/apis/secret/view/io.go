package view

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/secret"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// Basic APIs.

type CreateRequest struct {
	model.SecretCreateInput `json:",inline"`

	ProjectID oid.ID `query:"projectID,omitempty"`
}

func (r *CreateRequest) Validate() error {
	if r.ProjectID != "" && !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if r.Name == "" {
		return errors.New("invalid name: blank")
	}

	if r.Value == "" {
		return errors.New("invalid value: blank")
	}

	return nil
}

type CreateResponse = *model.SecretOutput

type DeleteRequest struct {
	model.SecretQueryInput `uri:",inline"`

	ProjectID oid.ID `query:"projectID,omitempty"`
}

func (r *DeleteRequest) Validate() error {
	if r.ProjectID != "" && !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type UpdateRequest struct {
	model.SecretUpdateInput `uri:",inline" json:",inline"`

	ProjectID oid.ID `query:"projectID,omitempty"`
}

func (r *UpdateRequest) Validate() error {
	if r.ProjectID != "" && !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	if r.Value == "" {
		return errors.New("invalid value: blank")
	}

	return nil
}

// Batch APIs.

type CollectionDeleteRequest []*model.SecretQueryInput

func (r CollectionDeleteRequest) Validate() error {
	if len(r) == 0 {
		return errors.New("invalid input: empty")
	}

	for _, i := range r {
		if !i.ID.Valid(0) {
			return errors.New("invalid id: blank")
		}
	}

	return nil
}

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.Secret, secret.OrderOption] `query:",inline"`

	ProjectID   oid.ID `query:"projectID,omitempty"`
	ProjectName string `query:"projectName,omitempty"`
	WithGlobal  bool   `query:"withGlobal,omitempty"`
}

func (r *CollectionGetRequest) ValidateWith(ctx context.Context, input any) error {
	var (
		err         error
		modelClient = input.(model.ClientSet)
	)

	if r.ProjectID != "" && !r.ProjectID.Valid(0) {
		return errors.New("invalid project id")
	}

	if r.ProjectName != "" {
		r.ProjectID, err = modelClient.Projects().Query().
			Where(project.Name(r.ProjectName)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "fail to get project by name")
		}
	}

	return nil
}

type CollectionGetResponse = []*model.SecretOutput

// Extensional APIs.
