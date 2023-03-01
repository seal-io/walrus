package view

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/platform/operator"
)

// Basic APIs

// GetRequest loads model.ApplicationResource with the request ID in validating.
type GetRequest struct {
	*model.ApplicationResourceQueryInput `uri:",inline"`

	Entity *model.ApplicationResource `json:"-"`
}

func (r *GetRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	var entity, err = modelClient.ApplicationResources().Query().
		Where(applicationresource.ID(r.ID)).
		Select(applicationresource.WithoutFields(
			applicationresource.FieldApplicationID,
			applicationresource.FieldUpdateTime)...).
		WithConnector(func(cq *model.ConnectorQuery) {
			cq.Select(
				connector.FieldName,
				connector.FieldType,
				connector.FieldConfigVersion,
				connector.FieldConfigData)
		}).
		Only(ctx)
	if err != nil {
		return err
	}
	r.Entity = entity
	return nil
}

// Batch APIs

// Extensional APIs

type GetKeysRequest = GetRequest

type GetKeysResponse = *operator.Keys

type StreamLogRequest struct {
	GetRequest `uri:",inline"`

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
	return r.GetRequest.ValidateWith(ctx, input)
}

type StreamExecRequest struct {
	GetRequest `uri:",inline"`

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
	return r.GetRequest.ValidateWith(ctx, input)
}
