package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/service"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/deployer/terraform"
	"github.com/seal-io/walrus/pkg/terraform/convertor"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/strs"
	"github.com/seal-io/walrus/utils/validation"
)

type (
	CreateRequest struct {
		model.ServiceCreateInput `path:",inline" json:",inline"`
	}

	CreateResponse = *model.ServiceOutput
)

func (r *CreateRequest) Validate() error {
	return ValidateCreateInput(r.ServiceCreateInput)
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
		return fmt.Errorf("failed to get service relationships: %w", err)
	}

	if len(ids) > 0 {
		names, err := dao.GetServiceNamesByIDs(r.Context, r.Client, ids...)
		if err != nil {
			return fmt.Errorf("failed to get services: %w", err)
		}

		return errorx.HttpErrorf(
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
			templateversion.FieldSchema,
			templateversion.FieldUiSchema,
		).
		All(r.Context)
	if err != nil {
		return fmt.Errorf("failed to get template version: %w", err)
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
		return fmt.Errorf("failed to get environment: %w", err)
	}

	tvm := make(map[object.ID]*model.TemplateVersion, len(tvs))

	// Validate template version whether match the target environment.
	for i := range tvs {
		if err = validateEnvironment(tvs[i], env); err != nil {
			return errorx.HttpErrorf(
				http.StatusBadRequest, "environment %s missing required connectors", env.Name)
		}

		// Map template version by ID for service validation.
		tvm[tvs[i].ID] = tvs[i]
	}

	for _, svc := range r.Items {
		// Verify service's variables with variables schema that defined on the template version.
		if !tvm[svc.Template.ID].Schema.IsEmpty() {
			err = svc.Attributes.ValidateWith(tvm[svc.Template.ID].Schema.VariableSchemas())
			if err != nil {
				return fmt.Errorf("invalid variables: %w", err)
			}
		}

		// Verify that variables in attributes are valid.
		err = validateVariable(r.Context, r.Client, svc.Attributes, svc.Name, r.Project.ID, r.Environment.ID)
		if err != nil {
			return err
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
		return fmt.Errorf("failed to get service dependencies: %w", err)
	}

	dependantIDSet := sets.New[object.ID](dependantIDs...)
	toDeleteIDSet := sets.New[object.ID](ids...)

	diffIDSet := dependantIDSet.Difference(toDeleteIDSet)
	if diffIDSet.Len() > 0 {
		names, err := dao.GetServiceNamesByIDs(r.Context, r.Client, diffIDSet.UnsortedList()...)
		if err != nil {
			return fmt.Errorf("failed to get services: %w", err)
		}

		return errorx.HttpErrorf(
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
		return errorx.NewHttpError(http.StatusBadRequest, "no connectors")
	}

	providers := make([]string, 0)

	if len(tv.Schema.RequiredProviders) != 0 {
		for _, provider := range tv.Schema.RequiredProviders {
			providers = append(providers, provider.Name)
		}
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
		return fmt.Errorf("failed to get service revisions: %w", err)
	}

	for _, r := range revisions {
		switch r.Status.SummaryStatus {
		case status.ServiceRevisionSummaryStatusSucceed:
		case status.ServiceRevisionSummaryStatusFailed:
		case status.ServiceRevisionSummaryStatusRunning:
			return errorx.HttpErrorf(
				http.StatusBadRequest,
				"deployment of service %q is running, please wait for it to finish",
				r.Edges.Service.Name,
			)
		default:
			return errorx.HttpErrorf(
				http.StatusBadRequest,
				"invalid deployment status of service %q: %s",
				r.Edges.Service.Name,
				r.Status.SummaryStatus,
			)
		}
	}

	return nil
}

func validateVariable(
	ctx context.Context,
	mc model.ClientSet,
	attributes property.Values,
	serviceName string,
	projectID object.ID,
	environmentID object.ID,
) error {
	attrs := make(map[string]any, len(attributes))
	for k, v := range attributes {
		attrs[k] = string(json.ShouldMarshal(v))
	}

	opts := terraform.ServiceOpts{
		ServiceName:   serviceName,
		ProjectID:     projectID,
		EnvironmentID: environmentID,
	}
	_, _, err := terraform.ParseModuleAttributes(ctx, mc, attrs, true, opts)

	return err
}

func ValidateCreateInput(sci model.ServiceCreateInput) error {
	if err := sci.Validate(); err != nil {
		return err
	}

	if err := validation.IsDNSLabel(sci.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	// Get template version.
	tv, err := sci.Client.TemplateVersions().Query().
		Where(templateversion.ID(sci.Template.ID)).
		Select(
			templateversion.FieldID,
			templateversion.FieldName,
			templateversion.FieldSchema,
			templateversion.FieldUiSchema).
		Only(sci.Context)
	if err != nil {
		return fmt.Errorf("failed to get template version: %w", err)
	}

	// Get environment.
	env, err := sci.Client.Environments().Query().
		Where(environment.ID(sci.Environment.ID)).
		Select(
			environment.FieldID,
			environment.FieldName).
		WithConnectors(func(rq *model.EnvironmentConnectorRelationshipQuery) {
			// Includes connectors.
			rq.WithConnector()
		}).
		Only(sci.Context)
	if err != nil {
		return fmt.Errorf("failed to get environment: %w", err)
	}

	// Validate template version whether match the target environment.
	if err = validateEnvironment(tv, env); err != nil {
		return err
	}

	// Verify variables with variables schema that defined on the template version.
	if !tv.Schema.IsEmpty() {
		err = sci.Attributes.ValidateWith(tv.Schema.VariableSchemas())
		if err != nil {
			return fmt.Errorf("invalid variables: %w", err)
		}
	}

	// Verify that variables in attributes are valid.
	err = validateVariable(sci.Context, sci.Client, sci.Attributes, sci.Name, sci.Project.ID, sci.Environment.ID)
	if err != nil {
		return err
	}

	return nil
}
