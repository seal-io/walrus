package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/templateversion"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/terraform/convertor"
	"github.com/seal-io/seal/utils/strs"
	"github.com/seal-io/seal/utils/validation"
)

type (
	CreateRequest struct {
		model.ServiceCreateInput `path:",inline" json:",inline"`

		RemarkTags []string `json:"remarkTags,omitempty"`
	}

	CreateResponse = *model.ServiceOutput
)

func (r *CreateRequest) Validate() error {
	if err := r.ServiceCreateInput.Validate(); err != nil {
		return err
	}

	if err := validation.IsDNSLabel(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	// Get template version.
	tv, err := r.Client.TemplateVersions().Query().
		Where(templateversion.ID(r.Template.ID)).
		Select(
			templateversion.FieldID,
			templateversion.FieldName,
			templateversion.FieldSchema).
		Only(r.Context)
	if err != nil {
		return runtime.Errorw(err, "failed to get template version")
	}

	// Get environment.
	env, err := r.Client.Environments().Query().
		Where(environment.ID(r.Environment.ID)).
		Select(
			environment.FieldID,
			environment.FieldName).
		WithConnectors(func(rq *model.EnvironmentConnectorRelationshipQuery) {
			// Includes connectors.
			rq.WithConnector()
		}).
		Only(r.Context)
	if err != nil {
		return runtime.Errorw(err, "failed to get environment")
	}

	// Validate template version whether match the target environment.
	if err = validateEnvironment(tv, env); err != nil {
		return err
	}

	// Verify variables with variables schema that defined on the template version.
	err = r.Attributes.ValidateWith(tv.Schema.Variables)
	if err != nil {
		return fmt.Errorf("invalid variables: %w", err)
	}

	return nil
}

type (
	GetRequest = model.ServiceQueryInput

	GetResponse = *model.ServiceOutput
)

type DeleteRequest struct {
	model.ServiceDeleteInput `path:",inline"`

	WithoutCleanup bool `query:"withoutCleanup,omitempty"`
}

func (r *DeleteRequest) Validate() error {
	if err := r.ServiceDeleteInput.Validate(); err != nil {
		return err
	}

	ids, err := dao.GetServiceDependantIDs(r.Context, r.Client, r.ID)
	if err != nil {
		return runtime.Errorw(err, "failed to get service relationships")
	}

	if len(ids) > 0 {
		names, err := dao.GetServiceNamesByIDs(r.Context, r.Client, ids...)
		if err != nil {
			return runtime.Errorw(err, "failed to get services")
		}

		return runtime.Errorf(
			http.StatusConflict,
			"service about to be deleted is the dependency of: %v",
			strs.Join(", ", names...),
		)
	}

	if !r.WithoutCleanup {
		if err = validateRevisionsStatus(r.Context, r.Client, r.ID); err != nil {
			return err
		}
	}

	return nil
}

type (
	CollectionCreateRequest struct {
		model.ServiceCreateInputs `path:",inline" json:",inline"`
	}

	CollectionCreateResponse = []*model.ServiceOutput
)

func (r *CollectionCreateRequest) Validate() error {
	if err := r.ServiceCreateInputs.Validate(); err != nil {
		return err
	}

	// Verify services.
	for i := range r.Items {
		if r.Items[i] == nil {
			return errors.New("empty service")
		}

		if err := validation.IsDNSLabel(r.Items[i].Name); err != nil {
			return fmt.Errorf("invalid service name: %w", err)
		}
	}

	// Get template versions.
	tvIDs := make([]object.ID, len(r.Items))
	for i := range r.Items {
		tvIDs[i] = r.Items[i].Template.ID
	}

	tvs, err := r.Client.TemplateVersions().Query().
		Where(templateversion.IDIn(tvIDs...)).
		Select(
			templateversion.FieldID,
			templateversion.FieldName,
			templateversion.FieldVersion,
			templateversion.FieldSchema).
		All(r.Context)
	if err != nil {
		return runtime.Errorw(err, "failed to get template version")
	}

	// Get environment.
	env, err := r.Client.Environments().Query().
		Where(environment.ID(r.Environment.ID)).
		Select(
			environment.FieldID,
			environment.FieldName).
		WithConnectors(func(rq *model.EnvironmentConnectorRelationshipQuery) {
			// Includes connectors.
			rq.WithConnector()
		}).
		Only(r.Context)
	if err != nil {
		return runtime.Errorw(err, "failed to get environment")
	}

	tvm := make(map[object.ID]*model.TemplateVersion, len(tvs))

	// Validate template version whether match the target environment.
	for i := range tvs {
		if err = validateEnvironment(tvs[i], env); err != nil {
			return runtime.Errorf(
				http.StatusBadRequest, "environment %s missing required connectors", env.Name)
		}

		// Map template version by ID for service validation.
		tvm[tvs[i].ID] = tvs[i]
	}

	// Verify service's variables with variables schema that defined on the template version.
	for _, svc := range r.Items {
		err = svc.Attributes.ValidateWith(tvm[svc.Template.ID].Schema.Variables)
		if err != nil {
			return fmt.Errorf("invalid variables: %w", err)
		}
	}

	return nil
}

type (
	CollectionGetRequest struct {
		model.ServiceQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.Service, service.OrderOption,
		] `query:",inline"`

		WithSchema bool `query:"withSchema,omitempty"`

		Stream *runtime.RequestUnidiStream
	}

	CollectionGetResponse = []*model.ServiceOutput
)

func (r *CollectionGetRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type CollectionDeleteRequest struct {
	model.ServiceDeleteInputs `path:",inline" json:",inline"`

	WithoutCleanup bool `query:"withoutCleanup,omitempty"`
}

func (r *CollectionDeleteRequest) Validate() error {
	if err := r.ServiceDeleteInputs.Validate(); err != nil {
		return err
	}

	ids := r.IDs()

	dependantIDs, err := dao.GetServiceDependantIDs(r.Context, r.Client, ids...)
	if err != nil {
		return runtime.Errorw(err, "failed to get service dependencies")
	}

	dependantIDSet := sets.New[object.ID](dependantIDs...)
	toDeleteIDSet := sets.New[object.ID](ids...)

	diffIDSet := dependantIDSet.Difference(toDeleteIDSet)
	if diffIDSet.Len() > 0 {
		names, err := dao.GetServiceNamesByIDs(r.Context, r.Client, diffIDSet.UnsortedList()...)
		if err != nil {
			return runtime.Errorw(err, "failed to get services")
		}

		return runtime.Errorf(
			http.StatusConflict,
			"service about to be deleted is the dependency of: %v",
			strs.Join(", ", names...),
		)
	}

	if r.WithoutCleanup {
		if err = validateRevisionsStatus(r.Context, r.Client, ids...); err != nil {
			return err
		}
	}

	return nil
}

func validateEnvironment(tv *model.TemplateVersion, env *model.Environment) error {
	if len(env.Edges.Connectors) == 0 {
		return runtime.Error(http.StatusBadRequest, errors.New("no connectors"))
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

func validateRevisionsStatus(ctx context.Context, mc model.ClientSet, ids ...object.ID) error {
	revisions, err := dao.GetLatestRevisions(ctx, mc, ids...)
	if err != nil {
		return runtime.Errorw(err, "failed to get service revisions")
	}

	for _, r := range revisions {
		switch r.Status {
		case status.ServiceRevisionStatusSucceeded:
		case status.ServiceRevisionStatusFailed:
		case status.ServiceRevisionStatusRunning:
			return runtime.Errorf(
				http.StatusBadRequest,
				"deployment of service %q is running, please wait for it to finish",
				r.Edges.Service.Name,
			)
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
