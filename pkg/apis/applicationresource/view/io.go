package view

import (
	"context"
	"errors"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platform/operator"
)

// Basic APIs

// Batch APIs

// Extensional APIs

type GetKeysRequest struct {
	EntityLoader `uri:",inline"`
}

type GetKeysResponse = operator.Keys

type StreamLogRequest struct {
	EntityLoader `uri:",inline"`

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
	return r.EntityLoader.ValidateWith(ctx, input)
}

type StreamExecRequest struct {
	EntityLoader `uri:",inline"`

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
	return r.EntityLoader.ValidateWith(ctx, input)
}

// EntityLoader loads model.ApplicationResource with the request ID in validating.
type EntityLoader struct {
	ID types.ID `uri:"id"`

	Entity *model.ApplicationResource `json:"-"`
}

func (r *EntityLoader) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	var entity, err = modelClient.ApplicationResources().Query().
		Where(applicationresource.ID(r.ID)).
		Select(
			applicationresource.FieldID,
			applicationresource.FieldConnectorID,
			applicationresource.FieldType,
			applicationresource.FieldName).
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
