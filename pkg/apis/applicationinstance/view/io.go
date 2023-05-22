package view

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"k8s.io/utils/pointer"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/validation"
)

// Basic APIs.

type CreateRequest struct {
	*model.ApplicationInstanceCreateInput `json:",inline"`

	RemarkTags []string `json:"remarkTags,omitempty"`
}

func (r *CreateRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if !r.Application.ID.Valid(0) {
		return errors.New("invalid application id: blank")
	}

	if !r.Environment.ID.Valid(0) {
		return errors.New("invalid environment id: blank")
	}

	if err := validation.IsDNSSubdomainName(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	// Verify application if it has no modules.
	app, err := modelClient.Applications().Query().
		Select(
			application.FieldID,
			application.FieldVariables).
		Where(application.ID(r.Application.ID)).
		Only(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get application")
	}

	count, _ := modelClient.ApplicationModuleRelationships().Query().
		Where(applicationmodulerelationship.ApplicationID(r.Application.ID)).
		Count(ctx)
	if count == 0 {
		return runtime.Error(http.StatusNotFound, "invalid application: no modules")
	}

	// Verify environment if it has no connectors.
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

	// Verify variables with variables schema that defined on application.
	err = r.Variables.ValidateWith(app.Variables)
	if err != nil {
		return fmt.Errorf("invalid variables: %w", err)
	}

	return nil
}

type CreateResponse = *model.ApplicationInstanceOutput

type DeleteRequest struct {
	*model.ApplicationInstanceQueryInput `uri:",inline"`

	Force *bool `query:"force,default=true"`
}

func (r *DeleteRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	if r.Force == nil {
		// By default, clean deployed native resources too.
		r.Force = pointer.Bool(true)
	}

	if *r.Force {
		modelClient := input.(model.ClientSet)

		err := validateRevisionStatus(ctx, modelClient, r.ID, "delete")
		if err != nil {
			return err
		}
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
	IDs        []oid.ID                           `json:"ids,omitempty"`
	Collection []*model.ApplicationInstanceOutput `json:"collection,omitempty"`
}

type StreamRequest struct {
	ID oid.ID `uri:"id"`
}

func (r *StreamRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	modelClient := input.(model.ClientSet)

	exist, err := modelClient.ApplicationInstances().Query().
		Where(applicationinstance.ID(r.ID)).
		Exist(ctx)
	if err != nil || !exist {
		return runtime.Errorw(err, "invalid id: not found")
	}

	return nil
}

// Batch APIs.

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.ApplicationInstance, applicationinstance.OrderOption] `query:",inline"`

	ApplicationID oid.ID `query:"applicationID"`
}

func (r *CollectionGetRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

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

	ApplicationID oid.ID `query:"applicationID,omitempty"`
}

func (r *CollectionStreamRequest) ValidateWith(ctx context.Context, input any) error {
	if r.ApplicationID != "" {
		modelClient := input.(model.ClientSet)

		if !r.ApplicationID.Valid(0) {
			return errors.New("invalid application id: blank")
		}

		_, err := modelClient.Applications().Query().
			Where(application.ID(r.ApplicationID)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get application")
		}
	}

	return nil
}

// Extensional APIs.

type RouteUpgradeRequest struct {
	_ struct{} `route:"PUT=/upgrade"`

	*model.ApplicationInstanceUpdateInput `uri:",inline" json:",inline"`

	RemarkTags []string `json:"remarkTags,omitempty"`
}

func (r *RouteUpgradeRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	ai, err := modelClient.ApplicationInstances().Query().
		Select(
			applicationinstance.FieldID,
			applicationinstance.FieldApplicationID).
		Where(applicationinstance.ID(r.ID)).
		WithApplication(func(aq *model.ApplicationQuery) {
			aq.Select(application.FieldVariables)
		}).
		Only(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get application instance")
	}

	// Verify variables with variables schema that defined on application.
	err = r.Variables.ValidateWith(ai.Edges.Application.Variables)
	if err != nil {
		return fmt.Errorf("invalid variables: %w", err)
	}

	err = validateRevisionStatus(ctx, modelClient, r.ID, "upgrade")
	if err != nil {
		return err
	}

	return nil
}

func IsEndpointOuput(outputName string) bool {
	return strings.HasPrefix(outputName, "endpoint")
}

type AccessEndpointRequest struct {
	_ struct{} `route:"GET=/access-endpoints"`

	*model.ApplicationInstanceQueryInput `uri:",inline"`
}

func (r *AccessEndpointRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

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

type AccessEndpointResponse = []Endpoint

type Endpoint struct {
	// ModuleName is the name of module.
	ModuleName string `json:"moduleName,omitempty"`
	// Name is identifier for the endpoint.
	Name string `json:"name,omitempty"`
	// Endpoint is access endpoint.
	Endpoints []string `json:"endpoints,omitempty"`
}

type OutputRequest struct {
	_ struct{} `route:"GET=/outputs"`

	*model.ApplicationInstanceQueryInput `uri:",inline"`
}

func (r *OutputRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

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

type OutputResponse = []types.OutputValue

type CreateCloneRequest struct {
	_ struct{} `route:"POST=/clone"`

	ID            oid.ID   `uri:"id"`
	EnviornmentID oid.ID   `json:"enviornmentID"`
	Name          string   `json:"name"`
	RemarkTags    []string `json:"remarkTags,omitempty"`
}

func (r *CreateCloneRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	if r.Name == "" {
		return errors.New("invalid name: blank")
	}

	if r.EnviornmentID != "" {
		if !r.EnviornmentID.IsNaive() {
			return fmt.Errorf("invalid environment id: %s", r.EnviornmentID)
		}
		modelClient := input.(model.ClientSet)

		_, err := modelClient.Environments().Query().
			Where(environment.ID(r.EnviornmentID)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get environment")
		}
	}

	return nil
}

func validateRevisionStatus(
	ctx context.Context,
	modelClient model.ClientSet,
	id oid.ID,
	action string,
) error {
	revision, err := modelClient.ApplicationRevisions().Query().
		Where(applicationrevision.InstanceID(id)).
		Order(model.Desc(applicationrevision.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return runtime.Errorw(err, "failed to get application deployment")
	}

	if revision != nil {
		switch revision.Status {
		case status.ApplicationRevisionStatusSucceeded:
		case status.ApplicationRevisionStatusRunning:
			return runtime.Error(http.StatusBadRequest,
				"deployment is running, please wait for it to finish before deleting the instance")
		case status.ApplicationRevisionStatusFailed:
			if action != "delete" {
				return nil
			}

			resourceExist, err := modelClient.ApplicationResources().Query().
				Where(applicationresource.InstanceID(id)).
				Exist(ctx)
			if err != nil {
				return err
			}

			if resourceExist {
				return runtime.Error(
					http.StatusBadRequest,
					"latest deployment is not succeeded,"+
						" please fix the app configuration or rollback the instance before deleting it",
				)
			}
		default:
			return runtime.Error(http.StatusBadRequest, "invalid deployment status")
		}
	}

	return nil
}

type StreamAccessEndpointResponse struct {
	Type       datamessage.EventType  `json:"type"`
	Collection AccessEndpointResponse `json:"collection,omitempty"`
}

type StreamOutputResponse struct {
	Type       datamessage.EventType `json:"type"`
	Collection OutputResponse        `json:"collection,omitempty"`
}

type Diff struct {
	Variables property.Schemas          `json:"variables"`
	Modules   []types.ApplicationModule `json:"modules"`
}

type DiffLatestRequest struct {
	ID oid.ID `uri:"id"`
}

func (r *DiffLatestRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type DiffLatestResponse struct {
	Old *Diff `json:"old,omitempty"`
	New *Diff `json:"new,omitempty"`
}
