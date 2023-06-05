package view

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	optypes "github.com/seal-io/seal/pkg/operator/types"
	"github.com/seal-io/seal/pkg/topic/datamessage"
)

// ServiceResourceQuery loads model.ServiceResource with the request ID in validating.
type ServiceResourceQuery struct {
	*model.ServiceResourceQueryInput `uri:",inline"`

	Entity *model.ServiceResource `uri:"-" json:"-"`
}

func (r *ServiceResourceQuery) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	entity, err := modelClient.ServiceResources().Query().
		Where(serviceresource.ID(r.ID)).
		Select(serviceresource.WithoutFields(
			serviceresource.FieldServiceID,
			serviceresource.FieldUpdateTime)...).
		WithConnector(func(cq *model.ConnectorQuery) {
			cq.Select(
				connector.FieldName,
				connector.FieldType,
				connector.FieldCategory,
				connector.FieldConfigVersion,
				connector.FieldConfigData)
		}).
		Only(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get service resource")
	}
	r.Entity = entity

	return nil
}

type StreamRequest struct {
	ID        oid.ID `uri:"id"`
	ProjectID oid.ID `query:"projectID"`
}

func (r *StreamRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	modelClient := input.(model.ClientSet)

	exist, err := modelClient.ServiceResources().Query().
		Where(serviceresource.ID(r.ID)).
		Exist(ctx)
	if err != nil || !exist {
		return runtime.Errorw(err, "invalid id: not found")
	}

	return nil
}

type StreamResponse struct {
	Type       datamessage.EventType `json:"type"`
	IDs        []oid.ID              `json:"ids,omitempty"`
	Collection []ServiceResource     `json:"collection,omitempty"`
}

// Batch APIs.

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.ServiceResource, serviceresource.OrderOption] `query:",inline"`

	ProjectID   oid.ID `query:"projectID"`
	ServiceID   oid.ID `query:"serviceID"`
	WithoutKeys bool   `query:"withoutKeys,omitempty"`
}

func (r *CollectionGetRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.ServiceID.Valid(0) {
		return errors.New("invalid service id: blank")
	}

	modelClient := input.(model.ClientSet)

	_, err := modelClient.Services().Query().
		Where(service.ID(r.ServiceID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get service")
	}

	return nil
}

type ServiceResource struct {
	Resource *model.ServiceResourceOutput `json:",inline"`
	Keys     *optypes.Keys                `json:"keys"`
}

type CollectionGetResponse = []ServiceResource

type CollectionStreamRequest struct {
	runtime.RequestExtracting `query:",inline"`

	ProjectID   oid.ID `query:"projectID"`
	ServiceID   oid.ID `query:"serviceID,omitempty"`
	WithoutKeys bool   `query:"withoutKeys,omitempty"`
}

func (r *CollectionStreamRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if r.ServiceID != "" {
		modelClient := input.(model.ClientSet)

		if !r.ServiceID.Valid(0) {
			return errors.New("invalid service id: blank")
		}

		_, err := modelClient.Services().Query().
			Where(service.ID(r.ServiceID)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get service")
		}
	}

	return nil
}

// Extensional APIs.

type GetKeysRequest struct {
	ServiceResourceQuery `query:"-" uri:",inline"`

	ProjectID oid.ID `query:"projectID"`
}

func (r *GetKeysRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	return r.ServiceResourceQuery.ValidateWith(ctx, input)
}

type GetKeysResponse = *optypes.Keys

type StreamLogRequest struct {
	ServiceResourceQuery `query:"-" uri:",inline"`

	ProjectID    oid.ID `query:"projectID"`
	Key          string `query:"key"`
	Previous     bool   `query:"previous,omitempty"`
	Tail         bool   `query:"tail,omitempty"`
	SinceSeconds *int64 `query:"sinceSeconds,omitempty"`
	Timestamps   bool   `query:"timestamps,omitempty"`
}

func (r *StreamLogRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if r.Key == "" {
		return errors.New("invalid key: blank")
	}

	if r.SinceSeconds != nil {
		if *r.SinceSeconds <= 0 {
			return errors.New("invalid since seconds: illegal")
		}
	}

	return r.ServiceResourceQuery.ValidateWith(ctx, input)
}

type StreamExecRequest struct {
	ServiceResourceQuery `query:"-" uri:",inline"`

	ProjectID oid.ID `query:"projectID"`
	Key       string `query:"key"`
	Shell     string `query:"shell,omitempty"`
	Width     int32  `query:"width,omitempty"`
	Height    int32  `query:"height,omitempty"`
}

func (r *StreamExecRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if r.Key == "" {
		return errors.New("invalid key: blank")
	}

	if r.Shell == "" {
		r.Shell = "sh"
	}

	if r.Width < 0 {
		return errors.New("invalid width: negative")
	} else if r.Width == 0 {
		r.Width = 100
	}

	if r.Height < 0 {
		return errors.New("invalid height: negative")
	} else if r.Height == 0 {
		r.Height = 100
	}

	return r.ServiceResourceQuery.ValidateWith(ctx, input)
}
