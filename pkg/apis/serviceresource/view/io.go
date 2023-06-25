package view

import (
	"context"
	"errors"
	"time"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/status"
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

	ProjectID       oid.ID `query:"projectID"`
	EnvironmentID   oid.ID `query:"environmentID,omitempty"`
	EnvironmentName string `query:"environmentName,omitempty"`
	ServiceID       oid.ID `query:"serviceID,omitempty"`
	ServiceName     string `query:"serviceName,omitempty"`
	WithoutKeys     bool   `query:"withoutKeys,omitempty"`
}

func (r *CollectionGetRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	modelClient := input.(model.ClientSet)

	switch {
	case r.ServiceID.Valid(0):
		_, err := modelClient.Services().Query().
			Where(service.ID(r.ServiceID)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get service")
		}
	case r.ServiceName != "":
		switch {
		case r.EnvironmentID.Valid(0):
			id, err := modelClient.Services().Query().
				Where(
					service.ProjectID(r.ProjectID),
					service.EnvironmentID(r.EnvironmentID),
					service.Name(r.ServiceName),
				).
				OnlyID(ctx)
			if err != nil {
				return runtime.Errorw(err, "failed to get service by name")
			}

			r.ServiceID = id
		case r.EnvironmentName != "":
			id, err := modelClient.Services().Query().
				Where(
					service.ProjectID(r.ProjectID),
					service.HasEnvironmentWith(environment.Name(r.EnvironmentName)),
					service.Name(r.ServiceName),
				).
				OnlyID(ctx)
			if err != nil {
				return runtime.Errorw(err, "failed to get service by name")
			}

			r.ServiceID = id
		default:
			return errors.New("both environment id and environment name are blank, " +
				"one of them is required while query by service name")
		}
	default:
		return errors.New("both service id and service name are blank")
	}

	return nil
}

type ServiceResource struct {
	*model.ServiceResourceOutput `json:",inline"`
	Keys                         *optypes.Keys `json:"keys"`
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

type CollectionGetGraphRequest struct {
	ProjectID       oid.ID `query:"projectID"`
	EnvironmentID   oid.ID `query:"environmentID,omitempty"`
	EnvironmentName string `query:"environmentName,omitempty"`
	ServiceID       oid.ID `query:"serviceID,omitempty"`
	ServiceName     string `query:"serviceName,omitempty"`
}

func (r *CollectionGetGraphRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	modelClient := input.(model.ClientSet)

	switch {
	case r.ServiceID.Valid(0):
		_, err := modelClient.Services().Query().
			Where(service.ID(r.ServiceID)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get service")
		}
	case r.ServiceName != "":
		var svcID oid.ID

		switch {
		case r.EnvironmentID.Valid(0):
			var err error

			svcID, err = modelClient.Services().Query().
				Where(
					service.ProjectID(r.ProjectID),
					service.EnvironmentID(r.EnvironmentID),
					service.Name(r.ServiceName)).
				OnlyID(ctx)
			if err != nil {
				return runtime.Errorw(err, "failed to get service")
			}
		case r.EnvironmentName != "":
			var (
				envID oid.ID
				err   error
			)

			envID, err = modelClient.Environments().Query().
				Where(
					environment.ProjectID(r.ProjectID),
					environment.Name(r.EnvironmentName),
				).
				OnlyID(ctx)
			if err != nil {
				return runtime.Errorw(err, "failed to get environment")
			}

			svcID, err = modelClient.Services().Query().
				Where(
					service.ProjectID(r.ProjectID),
					service.EnvironmentID(envID),
					service.Name(r.ServiceName)).
				OnlyID(ctx)
			if err != nil {
				return runtime.Errorw(err, "failed to get service")
			}
		default:
			return errors.New("both environment id and environment name are blank")
		}

		r.ServiceID = svcID
	default:
		return errors.New("both service id and service name are blank")
	}

	return nil
}

type CollectionGetGraphResponse struct {
	Vertices []GraphVertex `json:"vertices"`
	Edges    []GraphEdge   `json:"edges"`
}

// GraphVertexID defines the identifier of the vertex,
// which uniquely represents an API resource.
type GraphVertexID struct {
	// Kind indicates the kind of the resource,
	// which should be the same as the API handler's Kind returning.
	Kind string `json:"kind"`
	// ID indicates the identifier of the resource.
	ID any `json:"id,omitempty"`
}

// GraphVertex defines the vertex of graph.
type GraphVertex struct {
	GraphVertexID `json:",inline"`

	// Name indicates a human-readable string of the vertex.
	Name string `json:"name,omitempty"`
	// Description indicates the detail of the vertex.
	Description string `json:"description,omitempty"`
	// Labels indicates the labels of the vertex.
	Labels map[string]string `json:"labels,omitempty"`
	// CreateTime indicates the time to create the vertex.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// UpdateTime indicates the time to update the vertex.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// Status observes the status of the vertex.
	Status status.Summary `json:"status,omitempty"`
	// Extensions records the other information of the vertex,
	// e.g. the physical type, logical category or the operational keys of the resource.
	Extensions map[string]any `json:"extensions,omitempty"`
}

// GraphEdge defines the edge of graph.
type GraphEdge struct {
	// Type indicates the type of the edge, like: Dependency or Composition.
	Type string `json:"type"`
	// Start indicates the beginning of the edge.
	Start GraphVertexID `json:"start"`
	// End indicates the ending of the edge.
	End GraphVertexID `json:"end"`
}
