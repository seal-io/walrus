package view

import (
	"context"
	"errors"
	"net/http"

	"k8s.io/utils/pointer"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/topic/datamessage"
)

// Basic APIs

type CreateRequest struct {
	*model.ApplicationInstanceCreateInput `json:",inline"`
}

func (r *CreateRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.Application.ID.Valid(0) {
		return errors.New("invalid application id: blank")
	}
	if !r.Environment.ID.Valid(0) {
		return errors.New("invalid environment id: blank")
	}
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}

	_, err := modelClient.Applications().Query().
		Where(application.ID(r.Application.ID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get application")
	}
	count, _ := modelClient.ApplicationModuleRelationships().Query().
		Where(applicationmodulerelationship.ApplicationID(r.Application.ID)).
		Count(ctx)
	if count == 0 {
		return runtime.Error(http.StatusNotFound, "invalid application: no modules")
	}
	_, err = modelClient.Environments().Query().
		Where(environment.ID(r.Environment.ID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get environment")
	}
	count, _ = modelClient.EnvironmentConnectorRelationships().Query().
		Where(environmentconnectorrelationship.EnvironmentID(r.Environment.ID)).
		Count(ctx)
	if count == 0 {
		return runtime.Error(http.StatusNotFound, "invalid environment: no connectors")
	}
	return nil
}

type CreateResponse = *model.ApplicationInstanceOutput

type DeleteRequest struct {
	*model.ApplicationInstanceQueryInput `uri:",inline"`

	Force *bool `query:"force,default=true"`
}

func (r *DeleteRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	if r.Force == nil {
		// by default, clean deployed native resources too.
		r.Force = pointer.Bool(true)
	}
	return nil
}

type GetRequest struct {
	*model.ApplicationInstanceQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetResponse = *model.ApplicationInstanceOutput

type StreamResponse struct {
	Type       datamessage.EventType              `json:"type"`
	IDs        []types.ID                         `json:"ids"`
	Collection []*model.ApplicationInstanceOutput `json:"collection"`
}

type StreamRequest struct {
	ID types.ID `uri:"id"`
}

func (r *StreamRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	var modelClient = input.(model.ClientSet)
	exist, err := modelClient.ApplicationInstances().Query().
		Where(applicationinstance.ID(r.ID)).
		Exist(ctx)
	if err != nil || !exist {
		return runtime.Errorw(err, "invalid id: not found")
	}
	return nil
}

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.ApplicationInstance] `query:",inline"`

	ApplicationID types.ID `query:"applicationID"`
}

func (r *CollectionGetRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ApplicationID.Valid(0) {
		return errors.New("invalid application id: blank")
	}
	_, err := modelClient.Applications().Query().
		Where(application.ID(r.ApplicationID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get application")
	}
	return nil
}

type CollectionGetResponse = []*model.ApplicationInstanceOutput

type CollectionStreamRequest struct {
	runtime.RequestExtracting `query:",inline"`
}

// Extensional APIs

type RouteUpgradeRequest struct {
	_ struct{} `route:"PUT=/upgrade"`

	*model.ApplicationInstanceUpdateInput `uri:",inline" json:",inline"`
}

func (r *RouteUpgradeRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	_, err := modelClient.ApplicationInstances().Query().
		Where(applicationinstance.ID(r.ID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get application instance")
	}

	return nil
}

type AccessEndpointRequest struct {
	_ struct{} `route:"GET=/access-endpoints"`

	*model.ApplicationInstanceQueryInput `uri:",inline"`
}

func (r *AccessEndpointRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	_, err := modelClient.ApplicationInstances().Query().
		Where(applicationinstance.ID(r.ID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get application instance")
	}
	return nil
}

type AccessEndpointResponse struct {
	Endpoints []ResourceEndpoint `json:"endpoints"`
}

type ResourceEndpoint struct {
	// ResourceID is the namespaced name.
	ResourceID string `json:"resourceID,omitempty"`
	// ResourceKind be Ingress or Service.
	ResourceKind string `json:"resourceKind,omitempty"`
	// ResourceSubKind is the sub kind for endpoint, like nodePort, loadBalance.
	ResourceSubKind string `json:"resourceSubKind,omitempty"`
	// Endpoints is access endpoints.
	Endpoints []string `json:"endpoints,omitempty"`
}
