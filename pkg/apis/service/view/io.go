package view

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/pointer"

	"github.com/seal-io/seal/pkg/apis/runtime"
	serviceresourceview "github.com/seal-io/seal/pkg/apis/serviceresource/view"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/servicedependency"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/model/servicerevision"
	"github.com/seal-io/seal/pkg/dao/model/templateversion"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/strs"
	"github.com/seal-io/seal/utils/validation"
)

// Basic APIs.

type CreateRequest struct {
	model.ServiceCreateInput `json:",inline"`

	ProjectID   oid.ID   `query:"projectID"`
	ProjectName string   `query:"projectName"`
	RemarkTags  []string `json:"remarkTags,omitempty"`
}

func (r *CreateRequest) ValidateWith(ctx context.Context, input any) error {
	modelClient := input.(model.ClientSet)

	switch {
	case r.ProjectID != "":
		if !r.ProjectID.Valid(0) {
			return errors.New("invalid project id: blank")
		}
	case r.ProjectName != "":
		projectID, err := modelClient.Projects().Query().
			Where(project.Name(r.ProjectName)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get project")
		}

		r.ProjectID = projectID
	}

	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.Environment.ID.Valid(0) {
		return errors.New("invalid environment id: blank")
	}

	if err := validation.IsDNSSubdomainName(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	// Verify module version.

	templateVersion, err := modelClient.TemplateVersions().Query().
		Where(templateversion.TemplateID(r.Template.ID), templateversion.Version(r.Template.Version)).
		Only(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get template version")
	}

	// Verify environment if it has no connectors.
	_, err = modelClient.Environments().Query().
		Where(environment.ID(r.Environment.ID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get environment")
	}

	count, _ := modelClient.EnvironmentConnectorRelationships().Query().
		Where(environmentconnectorrelationship.EnvironmentID(r.Environment.ID)).
		Count(ctx)
	if count == 0 {
		return runtime.Error(http.StatusNotFound, "invalid environment: no connectors")
	}

	// Verify variables with variables schema that defined on the template version.
	err = r.Attributes.ValidateWith(templateVersion.Schema.Variables)
	if err != nil {
		return fmt.Errorf("invalid variables: %w", err)
	}

	return nil
}

type CreateResponse = *model.ServiceOutput

type DeleteRequest struct {
	model.ServiceQueryInput `uri:",inline"`

	ProjectID oid.ID `query:"projectID"`
	Force     *bool  `query:"force,default=true"`
}

func (r *DeleteRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	modelClient := input.(model.ClientSet)

	ids, err := modelClient.ServiceDependencies().Query().
		Where(servicedependency.DependentID(r.ID)).
		IDs(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get service dependencies")
	}

	if len(ids) > 0 {
		return runtime.Error(http.StatusConflict, "service has dependencies")
	}

	if r.Force == nil {
		// By default, clean deployed native resources too.
		r.Force = pointer.Bool(true)
	}

	if *r.Force {
		err := validateRevisionStatus(ctx, modelClient, r.ID, "delete")
		if err != nil {
			return err
		}
	}

	return nil
}

type GetRequest struct {
	model.ServiceQueryInput `uri:",inline"`

	ProjectID oid.ID `query:"projectID"`
}

func (r *GetRequest) Validate() error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetResponse = *model.ServiceOutput

type StreamRequest struct {
	ID oid.ID `uri:"id"`

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

	exist, err := modelClient.Services().Query().
		Where(service.ID(r.ID)).
		Exist(ctx)
	if err != nil || !exist {
		return runtime.Errorw(err, "invalid id: not found")
	}

	return nil
}

type StreamResponse struct {
	Type       datamessage.EventType  `json:"type"`
	IDs        []oid.ID               `json:"ids,omitempty"`
	Collection []*model.ServiceOutput `json:"collection,omitempty"`
}

// Batch APIs.

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.Service, service.OrderOption] `query:",inline"`

	ProjectID       oid.ID `query:"projectID"`
	EnvironmentID   oid.ID `query:"environmentID,omitempty"`
	EnvironmentName string `query:"environmentName,omitempty"`
}

func (r *CollectionGetRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	modelClient := input.(model.ClientSet)

	switch {
	case r.EnvironmentID.Valid(0):
		_, err := modelClient.Environments().Query().
			Where(environment.ID(r.EnvironmentID)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get environment")
		}
	case r.EnvironmentName != "":
		envID, err := modelClient.Environments().Query().
			Where(
				environment.ProjectID(r.ProjectID),
				environment.Name(r.EnvironmentName),
			).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get environment")
		}

		r.EnvironmentID = envID
	default:
		return errors.New("both environment id and environment name are blank")
	}

	return nil
}

type CollectionCreateRequest struct {
	EnvironmentIDs []oid.ID                   `json:"environmentIDs"`
	Services       []model.ServiceCreateInput `json:"services"`
}

func (r *CollectionCreateRequest) ValidateWith(ctx context.Context, input any) error {
	if len(r.EnvironmentIDs) == 0 {
		return errors.New("invalid environment ids: blank")
	}

	if len(r.Services) == 0 {
		return errors.New("invalid services: blank")
	}

	for _, envID := range r.EnvironmentIDs {
		if !envID.Valid(0) {
			return errors.New("invalid environment id: blank")
		}
	}

	modelClient := input.(model.ClientSet)

	environments, err := modelClient.Environments().Query().
		Select(environment.FieldID).
		Where(environment.IDIn(r.EnvironmentIDs...)).
		WithConnectors().
		All(ctx)
	if err != nil {
		return err
	}

	if len(environments) != len(r.EnvironmentIDs) {
		return errors.New("invalid environment IDs")
	}

	for _, env := range environments {
		if len(env.Edges.Connectors) == 0 {
			return fmt.Errorf("environment %s has no connectors", env.ID)
		}
	}

	// Get template versions.
	templateVersionKeys := sets.NewString()
	templateVersionPredicates := make([]predicate.TemplateVersion, 0)

	for _, s := range r.Services {
		key := strs.Join("/", s.Template.ID, s.Template.Version)
		if templateVersionKeys.Has(key) {
			continue
		}

		templateVersionKeys.Insert(key)

		templateVersionPredicates = append(templateVersionPredicates, templateversion.And(
			templateversion.TemplateID(s.Template.ID),
			templateversion.Version(s.Template.Version),
		))
	}

	templateVersions, err := modelClient.TemplateVersions().Query().
		Select(
			templateversion.FieldTemplateID,
			templateversion.FieldVersion,
			templateversion.FieldSchema,
		).
		Where(templateversion.Or(
			templateVersionPredicates...,
		)).
		All(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get template version")
	}
	templateVersionMap := make(map[string]*model.TemplateVersion, len(templateVersions))

	for _, tv := range templateVersions {
		key := strs.Join("/", tv.TemplateID, tv.Version)
		if _, ok := templateVersionMap[key]; !ok {
			templateVersionMap[key] = tv
		}
	}

	for _, s := range r.Services {
		if s.Name == "" {
			return errors.New("invalid service name: blank")
		}

		if err := validation.IsDNSSubdomainName(s.Name); err != nil {
			return fmt.Errorf("invalid name: %w", err)
		}

		// Verify template version.
		key := strs.Join("/", s.Template.ID, s.Template.Version)

		templateVersion, ok := templateVersionMap[key]
		if !ok {
			return runtime.Errorw(err, "failed to get template version")
		}

		// Verify variables with variables schema that defined on the template version.
		err = s.Attributes.ValidateWith(templateVersion.Schema.Variables)
		if err != nil {
			return fmt.Errorf("invalid variables: %w", err)
		}
	}

	return nil
}

type CollectionGetResponse = []*model.ServiceOutput

type CollectionStreamRequest struct {
	runtime.RequestExtracting `query:",inline"`

	ProjectID     oid.ID `query:"projectID"`
	EnvironmentID oid.ID `query:"environmentID,omitempty"`
}

func (r *CollectionStreamRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if r.EnvironmentID.Valid(0) {
		modelClient := input.(model.ClientSet)

		if !r.EnvironmentID.Valid(0) {
			return errors.New("invalid environment id: blank")
		}

		_, err := modelClient.Environments().Query().
			Where(environment.ID(r.EnvironmentID)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get environment")
		}
	}

	return nil
}

// Extensional APIs.

type RouteUpgradeRequest struct {
	_ struct{} `route:"PUT=/upgrade"`

	model.ServiceUpdateInput `uri:",inline" json:",inline"`

	ProjectID  oid.ID   `query:"projectID"`
	RemarkTags []string `json:"remarkTags,omitempty"`
}

func (r *RouteUpgradeRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	modelClient := input.(model.ClientSet)

	svc, err := modelClient.Services().Query().
		Select(
			service.FieldTemplate,
		).
		Where(service.ID(r.ID)).
		Only(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get service")
	}

	if r.Template.ID != svc.Template.ID {
		return errors.New("invalid template id: immutable")
	}

	tv, err := modelClient.TemplateVersions().Query().
		Select(templateversion.FieldSchema).
		Where(
			templateversion.TemplateID(r.Template.ID),
			templateversion.Version(r.Template.Version),
		).
		Only(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get template version")
	}

	// Verify attributes with variables schema of the template version.
	err = r.Attributes.ValidateWith(tv.Schema.Variables)
	if err != nil {
		return fmt.Errorf("invalid variables: %w", err)
	}

	err = validateRevisionStatus(ctx, modelClient, r.ID, "upgrade")
	if err != nil {
		return err
	}

	return nil
}

type RouteRollbackRequest struct {
	_ struct{} `route:"POST=/rollback"`

	model.ServiceQueryInput `uri:",inline" json:",inline"`

	ProjectID  oid.ID `query:"projectID"`
	RevisionID oid.ID `query:"revisionID"`
}

func (r *RouteRollbackRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.RevisionID.Valid(0) {
		return errors.New("invalid revision id: blank")
	}

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	// Check if the latest revision is running.
	modelClient := input.(model.ClientSet)

	latestRevision, err := modelClient.ServiceRevisions().Query().
		Select(servicerevision.FieldStatus).
		Where(servicerevision.ServiceID(r.ID)).
		Order(model.Desc(servicerevision.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return runtime.Errorw(err, "failed to get the latest revision")
	}

	if latestRevision.Status == status.ServiceRevisionStatusRunning {
		return errors.New("latest revision is running")
	}

	return nil
}

func IsEndpointOutput(outputName string) bool {
	return strings.HasPrefix(outputName, "endpoint")
}

type RouteAccessEndpointRequest struct {
	_ struct{} `route:"GET=/access-endpoints"`

	model.ServiceQueryInput `uri:",inline"`

	ProjectID oid.ID `query:"projectID"`
}

func (r *RouteAccessEndpointRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	modelClient := input.(model.ClientSet)

	_, err := modelClient.Services().Query().
		Where(service.ID(r.ID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get service")
	}

	return nil
}

type RouteAccessEndpointResponse = []AccessEndpoint

type AccessEndpoint struct {
	// Name is identifier for the endpoint.
	Name string `json:"name,omitempty"`
	// Endpoint is access endpoint.
	Endpoints []string `json:"endpoints,omitempty"`
}

type RouteOutputRequest struct {
	_ struct{} `route:"GET=/outputs"`

	model.ServiceQueryInput `uri:",inline"`

	ProjectID oid.ID `query:"projectID"`
}

func (r *RouteOutputRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	modelClient := input.(model.ClientSet)

	_, err := modelClient.Services().Query().
		Where(service.ID(r.ID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get service")
	}

	return nil
}

type RouteOutputResponse = []types.OutputValue

type CreateCloneRequest struct {
	ID            oid.ID   `uri:"id"`
	EnvironmentID oid.ID   `json:"environmentID"`
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

	if r.EnvironmentID != "" {
		if !r.EnvironmentID.IsNaive() {
			return fmt.Errorf("invalid environment id: %s", r.EnvironmentID)
		}
		modelClient := input.(model.ClientSet)

		_, err := modelClient.Environments().Query().
			Where(environment.ID(r.EnvironmentID)).
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
	revision, err := modelClient.ServiceRevisions().Query().
		Where(servicerevision.ServiceID(id)).
		Order(model.Desc(servicerevision.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return runtime.Errorw(err, "failed to get service revision")
	}

	if revision != nil {
		switch revision.Status {
		case status.ServiceRevisionStatusSucceeded:
		case status.ServiceRevisionStatusRunning:
			return runtime.Error(http.StatusBadRequest,
				"deployment is running, please wait for it to finish before deleting the service")
		case status.ServiceRevisionStatusFailed:
			if action != "delete" {
				return nil
			}

			resourceExist, err := modelClient.ServiceResources().Query().
				Where(serviceresource.ServiceID(id)).
				Exist(ctx)
			if err != nil {
				return err
			}

			if resourceExist {
				return runtime.Error(
					http.StatusBadRequest,
					"latest deployment is not succeeded,"+
						" please fix the service configuration or rollback before deleting it",
				)
			}
		default:
			return runtime.Error(http.StatusBadRequest, "invalid deployment status")
		}
	}

	return nil
}

type StreamAccessEndpointResponse struct {
	Type       datamessage.EventType       `json:"type"`
	Collection RouteAccessEndpointResponse `json:"collection,omitempty"`
}

type StreamOutputResponse struct {
	Type       datamessage.EventType `json:"type"`
	Collection RouteOutputResponse   `json:"collection,omitempty"`
}

type CollectionGetGraphRequest struct {
	ProjectID       oid.ID `query:"projectID"`
	EnvironmentID   oid.ID `query:"environmentID,omitempty"`
	EnvironmentName string `query:"environmentName,omitempty"`
}

func (r *CollectionGetGraphRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ProjectID.Valid(0) {
		return errors.New("invalid project id: blank")
	}

	modelClient := input.(model.ClientSet)

	switch {
	case r.EnvironmentID.Valid(0):
		_, err := modelClient.Environments().Query().
			Where(environment.ID(r.EnvironmentID)).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get environment")
		}
	case r.EnvironmentName != "":
		envID, err := modelClient.Environments().Query().
			Where(
				environment.ProjectID(r.ProjectID),
				environment.Name(r.EnvironmentName),
			).
			OnlyID(ctx)
		if err != nil {
			return runtime.Errorw(err, "failed to get environment")
		}

		r.EnvironmentID = envID
	default:
		return errors.New("both environment id and environment name are blank")
	}

	return nil
}

type (
	CollectionGetGraphResponse = serviceresourceview.CollectionGetGraphResponse

	// GraphVertexID defines the identifier of the vertex,
	// which uniquely represents an API resource.
	GraphVertexID = serviceresourceview.GraphVertexID
	// GraphVertex defines the vertex of graph.
	GraphVertex = serviceresourceview.GraphVertex
	// GraphEdge defines the edge of graph.
	GraphEdge = serviceresourceview.GraphEdge
)
