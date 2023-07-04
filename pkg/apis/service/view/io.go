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
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/model/servicerevision"
	"github.com/seal-io/seal/pkg/dao/model/templateversion"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/terraform/convertor"
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

	if err := validation.IsDNSLabel(r.Name); err != nil {
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
	env, err := modelClient.Environments().Query().
		WithConnectors(func(rq *model.EnvironmentConnectorRelationshipQuery) {
			// Includes connectors.
			rq.WithConnector()
		}).
		Where(environment.ID(r.Environment.ID)).
		Only(ctx)
	if err != nil {
		return runtime.Errorw(err, "failed to get environment")
	}

	if err = validateEnvironment(templateVersion, env); err != nil {
		return err
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

	ids, err := dao.GetServiceDependantIDs(ctx, modelClient, r.ID)
	if err != nil {
		return runtime.Errorw(err, "failed to get service relationships")
	}

	if len(ids) > 0 {
		names, err := dao.GetServiceNamesByIDs(ctx, modelClient, ids...)
		if err != nil {
			return runtime.Errorw(err, "failed to get services")
		}

		return runtime.Errorf(
			http.StatusConflict,
			"service about to be deleted is the dependency of: %v",
			strs.Join(", ", names...),
		)
	}

	if r.Force == nil {
		// By default, clean deployed native resources too.
		r.Force = pointer.Bool(true)
	}

	if *r.Force {
		err := validateRevisionsStatus(ctx, modelClient, "delete", r.ID)
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
		Select(
			environment.FieldID,
			environment.FieldName,
		).
		Where(environment.IDIn(r.EnvironmentIDs...)).
		WithConnectors(func(rq *model.EnvironmentConnectorRelationshipQuery) {
			// Includes connectors.
			rq.WithConnector()
		}).
		All(ctx)
	if err != nil {
		return err
	}

	if len(environments) != len(r.EnvironmentIDs) {
		return errors.New("invalid environment IDs")
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

		for _, env := range environments {
			if err := validateEnvironment(tv, env); err != nil {
				return runtime.Errorf(
					http.StatusBadRequest, "environment %s missing required connectors", env.Name)
			}
		}
	}

	for _, s := range r.Services {
		if s.Name == "" {
			return errors.New("invalid service name: blank")
		}

		if err := validation.IsDNSLabel(s.Name); err != nil {
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

type CollectionDeleteRequest struct {
	IDs []oid.ID `json:"ids,omitempty"`

	ProjectID oid.ID `query:"projectID"`
	Force     bool   `query:"force,default=true"`
}

func (r CollectionDeleteRequest) ValidateWith(ctx context.Context, input any) error {
	if len(r.IDs) == 0 {
		return errors.New("invalid input: empty")
	}

	for _, i := range r.IDs {
		if !i.Valid(0) {
			return errors.New("invalid id: blank")
		}
	}

	modelClient := input.(model.ClientSet)

	ids, err := dao.GetServiceDependantIDs(ctx, modelClient, r.IDs...)
	if err != nil {
		return runtime.Errorw(err, "failed to get service dependencies")
	}

	dependantIDSet := sets.New[oid.ID](ids...)
	toDeleteIDSet := sets.New[oid.ID](r.IDs...)

	diffIDSet := dependantIDSet.Difference(toDeleteIDSet)
	if diffIDSet.Len() > 0 {
		names, err := dao.GetServiceNamesByIDs(ctx, modelClient, diffIDSet.UnsortedList()...)
		if err != nil {
			return runtime.Errorw(err, "failed to get services")
		}

		return runtime.Errorf(
			http.StatusConflict,
			"service about to be deleted is the dependency of: %v",
			strs.Join(", ", names...),
		)
	}

	if r.Force {
		if err = validateRevisionsStatus(ctx, modelClient, "delete", r.IDs...); err != nil {
			return err
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

	err = validateRevisionsStatus(ctx, modelClient, "upgrade", r.ID)
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

func validateRevisionsStatus(
	ctx context.Context,
	modelClient model.ClientSet,
	action string,
	serviceIDs ...oid.ID,
) error {
	revisions, err := dao.GetLatestRevisions(ctx, modelClient, serviceIDs...)
	if err != nil {
		return runtime.Errorw(err, "failed to get service revisions")
	}

	for _, r := range revisions {
		switch r.Status {
		case status.ServiceRevisionStatusSucceeded:
		case status.ServiceRevisionStatusRunning:
			return runtime.Errorf(
				http.StatusBadRequest,
				"deployment of service %q is running, please wait for it to finish before deleting it",
				r.Edges.Service.Name,
			)
		case status.ServiceRevisionStatusFailed:
			if action != "delete" {
				return nil
			}

			resourceExist, err := modelClient.ServiceResources().Query().
				Where(serviceresource.ServiceID(r.ServiceID)).
				Exist(ctx)
			if err != nil {
				return err
			}

			if resourceExist {
				return runtime.Errorf(
					http.StatusBadRequest,
					"latest deployment of %q is not succeeded,"+
						" please fix the service configuration or rollback before deleting it",
					r.Edges.Service.Name,
				)
			}
		default:
			return runtime.Errorf(
				http.StatusBadRequest,
				"invalid deployment status of service %q: %s",
				r.Edges.Service.Name,
				r.Status,
			)
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

func validateEnvironment(
	tv *model.TemplateVersion,
	env *model.Environment,
) error {
	if len(env.Edges.Connectors) == 0 {
		return runtime.Error(
			http.StatusBadRequest,
			errors.New("no connectors"),
		)
	}

	providers := make([]string, len(tv.Schema.RequiredProviders))

	for i, provider := range tv.Schema.RequiredProviders {
		providers[i] = provider.Name
	}

	var connectors model.Connectors

	for _, ecr := range env.Edges.Connectors {
		connectors = append(connectors, ecr.Edges.Connector)
	}

	_, err := convertor.ToProvidersBlocks(providers, connectors, convertor.ConvertOptions{
		Providers: providers,
	})

	return err
}
