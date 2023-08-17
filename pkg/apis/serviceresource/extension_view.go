package serviceresource

import (
	"errors"
	"fmt"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/model/serviceresource"
	"github.com/seal-io/walrus/pkg/dao/types"
)

type Request struct {
	model.ServiceResourceQueryInput `path:",inline"`

	Entity *model.ServiceResource
}

func (r *Request) Validate() error {
	if err := r.ServiceResourceQueryInput.Validate(); err != nil {
		return err
	}

	entity, err := r.Client.ServiceResources().Query().
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
		Only(r.Context)
	if err != nil {
		return fmt.Errorf("failed to get service resource: %w", err)
	}

	r.Entity = entity

	return nil
}

type (
	RouteGetKeysRequest struct {
		_ struct{} `route:"GET=/keys"`

		Request `path:",inline"`
	}

	RouteGetKeysResponse = *types.ServiceResourceOperationKeys
)

type RouteLogRequest struct {
	_ struct{} `route:"GET=/log"`

	Request ` path:",inline"`

	Key          string `query:"key"`
	Previous     bool   `query:"previous,omitempty"`
	SinceSeconds *int64 `query:"sinceSeconds,omitempty"`
	TailLines    *int64 `query:"tailLines,omitempty"`
	Timestamps   bool   `query:"timestamps,omitempty"`

	Stream *runtime.RequestUnidiStream
}

func (r *RouteLogRequest) Validate() error {
	if err := r.Request.Validate(); err != nil {
		return err
	}

	if r.Key == "" {
		return errors.New("invalid key: blank")
	}

	if r.SinceSeconds != nil {
		if *r.SinceSeconds <= 0 {
			return errors.New("invalid since seconds: illegal")
		}
	}

	if r.TailLines != nil {
		if *r.TailLines <= 0 {
			return errors.New("invalid tail lines: illegal")
		}
	}

	return nil
}

func (r *RouteLogRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type RouteExecRequest struct {
	_ struct{} `route:"GET=/exec"`

	Request ` path:",inline"`

	Key    string `query:"key"`
	Shell  string `query:"shell,omitempty"`
	Width  int32  `query:"width,omitempty"`
	Height int32  `query:"height,omitempty"`

	Stream *runtime.RequestBidiStream
}

func (r *RouteExecRequest) Validate() error {
	if err := r.Request.Validate(); err != nil {
		return err
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

	return nil
}

func (r *RouteExecRequest) SetStream(stream runtime.RequestBidiStream) {
	r.Stream = &stream
}
