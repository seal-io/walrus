package view

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platform/operator"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/json"
)

// ApplicationResourceQuery loads model.ApplicationResource with the request ID in validating.
type ApplicationResourceQuery struct {
	*model.ApplicationResourceQueryInput `uri:",inline"`

	Entity *model.ApplicationResource `uri:"-" json:"-"`
}

func (r *ApplicationResourceQuery) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	var entity, err = modelClient.ApplicationResources().Query().
		Where(applicationresource.ID(r.ID)).
		Select(applicationresource.WithoutFields(
			applicationresource.FieldInstanceID,
			applicationresource.FieldUpdateTime)...).
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
		return runtime.Errorw(err, "failed to get application resource")
	}
	r.Entity = entity
	return nil
}

type StreamResponse struct {
	Type       datamessage.EventType `json:"type"`
	IDs        []types.ID            `json:"ids,omitempty"`
	Collection []ApplicationResource `json:"collection,omitempty"`
}

type StreamRequest struct {
	ID types.ID `uri:"id"`
}

func (r *StreamRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	var modelClient = input.(model.ClientSet)
	exist, err := modelClient.ApplicationResources().Query().
		Where(applicationresource.ID(r.ID)).
		Exist(ctx)
	if err != nil || !exist {
		return runtime.Errorw(err, "invalid id: not found")
	}

	return nil
}

// Basic APIs

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.ApplicationResource, applicationresource.OrderOption] `query:",inline"`

	InstanceID  types.ID `query:"instanceID"`
	WithoutKeys bool     `query:"withoutKeys,omitempty"`
}

func (r *CollectionGetRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.InstanceID.Valid(0) {
		return errors.New("invalid instance id: blank")
	}
	_, err := modelClient.ApplicationInstances().Query().
		Where(applicationinstance.ID(r.InstanceID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get application instance")
	}
	return nil
}

type ApplicationResource struct {
	Resource *model.ApplicationResourceOutput
	Keys     *operator.Keys
}

// MarshalJSON implements the json.Marshaler to avoid the impact from model.ApplicationResourceOutput's marshaller.
func (in ApplicationResource) MarshalJSON() ([]byte, error) {
	type (
		AliasResource model.ApplicationResourceOutput
		Alias         struct {
			*AliasResource `json:",inline"`
			Keys           *operator.Keys `json:"keys"`
		}
	)
	return json.Marshal(&Alias{
		AliasResource: (*AliasResource)(in.Resource.Normalize()),
		Keys:          in.Keys,
	})
}

type CollectionGetResponse = []ApplicationResource

type CollectionStreamRequest struct {
	runtime.RequestExtracting `query:",inline"`

	InstanceID  types.ID `query:"instanceID,omitempty"`
	WithoutKeys bool     `query:"withoutKeys,omitempty"`
}

func (r *CollectionStreamRequest) ValidateWith(ctx context.Context, input any) error {
	if r.InstanceID != "" {
		var modelClient = input.(model.ClientSet)
		if !r.InstanceID.Valid(0) {
			return errors.New("invalid instance id: blank")
		}
		_, err := modelClient.ApplicationInstances().Query().
			Where(applicationinstance.ID(r.InstanceID)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get application instance")
		}
	}

	return nil
}

// Extensional APIs

type GetKeysRequest = ApplicationResourceQuery

type GetKeysResponse = *operator.Keys

type StreamLogRequest struct {
	ApplicationResourceQuery `query:"-" uri:",inline"`

	Key          string `query:"key"`
	Previous     bool   `query:"previous,omitempty"`
	Tail         bool   `query:"tail,omitempty"`
	SinceSeconds *int64 `query:"sinceSeconds,omitempty"`
	Timestamps   bool   `query:"timestamps,omitempty"`
}

func (r *StreamLogRequest) ValidateWith(ctx context.Context, input any) error {
	if r.Key == "" {
		return errors.New("invalid key: blank")
	}
	if r.SinceSeconds != nil {
		if *r.SinceSeconds <= 0 {
			return errors.New("invalid since seconds: illegal")
		}
	}
	return r.ApplicationResourceQuery.ValidateWith(ctx, input)
}

type StreamExecRequest struct {
	ApplicationResourceQuery `query:"-" uri:",inline"`

	Key    string `query:"key"`
	Shell  string `query:"shell,omitempty"`
	Width  int32  `query:"width,omitempty"`
	Height int32  `query:"height,omitempty"`
}

func (r *StreamExecRequest) ValidateWith(ctx context.Context, input any) error {
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
	return r.ApplicationResourceQuery.ValidateWith(ctx, input)
}
