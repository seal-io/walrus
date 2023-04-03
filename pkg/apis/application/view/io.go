package view

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/topic/datamessage"
)

// Basic APIs

type CreateRequest struct {
	*model.ApplicationCreateInput `json:",inline"`
}

func (r *CreateRequest) Validate() error {
	if r.Project.ID == "" {
		return errors.New("invalid project id: blank")
	}
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}
	if len(r.Modules) != 0 {
		for i := 0; i < len(r.Modules); i++ {
			if r.Modules[i].Module.ID == "" {
				return errors.New("invalid module id: blank")
			}
			if r.Modules[i].Name == "" {
				return errors.New("invalid module name: blank")
			}
		}
	}
	return nil
}

type CreateResponse = *model.ApplicationOutput

type DeleteRequest = GetRequest

type UpdateRequest struct {
	*model.ApplicationUpdateInput `uri:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	if len(r.Modules) != 0 {
		for i := 0; i < len(r.Modules); i++ {
			if r.Modules[i].Module.ID == "" {
				return errors.New("invalid module id: blank")
			}
			if r.Modules[i].Name == "" {
				return errors.New("invalid module name: blank")
			}
		}
	}
	return nil
}

type GetRequest struct {
	*model.ApplicationQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetResponse = *model.ApplicationOutput

type StreamResponse struct {
	Type       datamessage.EventType      `json:"type"`
	IDs        []types.ID                 `json:"ids"`
	Collection []*model.ApplicationOutput `json:"collection"`
}

type StreamRequest struct {
	ID types.ID `uri:"id"`
}

func (r *StreamRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	var client = input.(model.ClientSet)
	exist, err := client.Applications().Query().
		Where(application.ID(r.ID)).
		Exist(ctx)
	if err != nil || !exist {
		return runtime.Errorw(err, "invalid id: not found")
	}

	return nil
}

// Batch APIs

type CollectionDeleteRequest []*model.ApplicationQueryInput

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
	runtime.RequestCollection[predicate.Application] `query:",inline"`

	ProjectIDs []types.ID `query:"projectID"`
}

func (r *CollectionGetRequest) Validate() error {
	if len(r.ProjectIDs) == 0 {
		return errors.New("invalid input: missing project id")
	}

	for i := range r.ProjectIDs {
		if !r.ProjectIDs[i].Valid(0) {
			return errors.New("invalid project id: blank")
		}
	}
	return nil
}

type CollectionGetResponse = []*model.ApplicationOutput

type CollectionStreamRequest struct {
	runtime.RequestExtracting `query:",inline"`

	ProjectIDs []types.ID `query:"projectID,omitempty"`
}

func (r *CollectionStreamRequest) Validate() error {
	for i := range r.ProjectIDs {
		if !r.ProjectIDs[i].Valid(0) {
			return errors.New("invalid project id: blank")
		}
	}
	return nil
}

// Extensional APIs
