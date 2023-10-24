package resource

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
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinition"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinitionmatchingrule"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/deployer/terraform"
	"github.com/seal-io/walrus/pkg/resourcedefinitions"
	"github.com/seal-io/walrus/pkg/terraform/convertor"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/strs"
	"github.com/seal-io/walrus/utils/validation"
)

type (
	CreateRequest struct {
		model.ResourceCreateInput `path:",inline" json:",inline"`
	}

	CreateResponse = *model.ResourceOutput
)

func (r *CreateRequest) Validate() error {
	return ValidateCreateInput(&r.ResourceCreateInput)
}

type (
	GetRequest = model.ResourceQueryInput

	GetResponse = *model.ResourceOutput
)

type DeleteRequest struct {
	model.ResourceDeleteInput `path:",inline"`

	WithoutCleanup bool `query:"withoutCleanup,omitempty"`
}

func (r *DeleteRequest) Validate() error {
	if err := r.ResourceDeleteInput.Validate(); err != nil {
		return err
	}

	ids, err := dao.GetResourceDependantIDs(r.Context, r.Client, r.ID)
	if err != nil {
		return fmt.Errorf("failed to get resource relationships: %w", err)
	}

	if len(ids) > 0 {
		names, err := dao.GetResourceNamesByIDs(r.Context, r.Client, ids...)
		if err != nil {
			return fmt.Errorf("failed to get resources: %w", err)
		}

		return errorx.HttpErrorf(
			http.StatusConflict,
			"resource about to be deleted is the dependency of: %v",
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
		model.ResourceCreateInputs `path:",inline" json:",inline"`
	}

	CollectionCreateResponse = []*model.ResourceOutput
)

func (r *CollectionCreateRequest) Validate() error {
	// Resource type maps to type in definition edge.
	for _, item := range r.Items {
		if item.Type != "" {
			item.ResourceDefinition = &model.ResourceDefinitionQueryInput{
				Type: item.Type,
			}
		}
	}

	if err := r.ResourceCreateInputs.Validate(); err != nil {
		return err
	}

	// Verify resources.
	for i := range r.Items {
		if r.Items[i] == nil {
			return errors.New("empty resource")
		}

		if err := validation.IsDNSLabel(r.Items[i].Name); err != nil {
			return fmt.Errorf("invalid resource name: %w", err)
		}
	}

	var tvIDs []object.ID
	// Get template versions.
	for i := range r.Items {
		if r.Items[i].Template == nil {
			continue
		}
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
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(project.FieldName)
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

		// Map template version by ID for resource validation.
		tvm[tvs[i].ID] = tvs[i]
	}

	// Get resource definitions.
	var rdTypes []string

	for i := range r.Items {
		if r.Items[i].Type != "" {
			rdTypes = append(rdTypes, r.Items[i].Type)
		}
	}

	rds, err := r.Client.ResourceDefinitions().Query().
		Where(resourcedefinition.TypeIn(rdTypes...)).
		Select(
			resourcedefinition.FieldID,
			resourcedefinition.FieldName,
		).
		WithMatchingRules(func(rq *model.ResourceDefinitionMatchingRuleQuery) {
			rq.Order(model.Desc(resourcedefinitionmatchingrule.FieldCreateTime)).
				Select(resourcedefinitionmatchingrule.FieldResourceDefinitionID).
				Unique(false).
				Select(resourcedefinitionmatchingrule.FieldTemplateID).
				WithTemplate(func(tq *model.TemplateVersionQuery) {
					tq.Select(
						templateversion.FieldID,
						templateversion.FieldVersion,
						templateversion.FieldName,
					)
				})
		}).
		All(r.Context)
	if err != nil {
		return fmt.Errorf("failed to get resource definition: %w", err)
	}

	rdm := make(map[string]*model.ResourceDefinition, len(rds))
	for _, rd := range rds {
		rdm[rd.Type] = rd
	}

	for _, res := range r.Items {
		if res.Template != nil {
			// Verify resource's variables with variables schema that defined on the template version.
			if !tvm[res.Template.ID].Schema.IsEmpty() {
				err = res.Attributes.ValidateWith(tvm[res.Template.ID].Schema.VariableSchemas())
				if err != nil {
					return fmt.Errorf("invalid variables: %w", err)
				}
			} else if res.Type != "" {
				rule := resourcedefinitions.Match(
					rdm[res.Type].Edges.MatchingRules,
					env.Edges.Project.Name,
					env.Name,
					env.Type,
					env.Labels,
					res.Labels,
				)
				if rule == nil {
					return fmt.Errorf("no matching resource definition for %q", res.Name)
				}
			}
			// Verify that variables in attributes are valid.
			err = validateVariable(r.Context, r.Client, res.Attributes, res.Name, r.Project.ID, r.Environment.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type (
	CollectionGetRequest struct {
		model.ResourceQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.Resource, resource.OrderOption,
		] `query:",inline"`

		IsService *bool `query:"isService,omitempty"`

		WithSchema bool `query:"withSchema,omitempty"`

		Stream *runtime.RequestUnidiStream
	}

	CollectionGetResponse = []*model.ResourceOutput
)

func (r *CollectionGetRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type CollectionDeleteRequest struct {
	model.ResourceDeleteInputs `path:",inline" json:",inline"`

	WithoutCleanup bool `query:"withoutCleanup,omitempty"`
}

func (r *CollectionDeleteRequest) Validate() error {
	if err := r.ResourceDeleteInputs.Validate(); err != nil {
		return err
	}

	ids := r.IDs()

	dependantIDs, err := dao.GetResourceDependantIDs(r.Context, r.Client, ids...)
	if err != nil {
		return fmt.Errorf("failed to get resource dependencies: %w", err)
	}

	dependantIDSet := sets.New[object.ID](dependantIDs...)
	toDeleteIDSet := sets.New[object.ID](ids...)

	diffIDSet := dependantIDSet.Difference(toDeleteIDSet)
	if diffIDSet.Len() > 0 {
		names, err := dao.GetResourceNamesByIDs(r.Context, r.Client, diffIDSet.UnsortedList()...)
		if err != nil {
			return fmt.Errorf("failed to get resources: %w", err)
		}

		return errorx.HttpErrorf(
			http.StatusConflict,
			"resource about to be deleted is the dependency of: %v",
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
		return fmt.Errorf("failed to get resource revisions: %w", err)
	}

	for _, r := range revisions {
		switch r.Status.SummaryStatus {
		case status.ResourceRevisionSummaryStatusSucceed:
		case status.ResourceRevisionSummaryStatusFailed:
		case status.ResourceRevisionSummaryStatusRunning:
			return errorx.HttpErrorf(
				http.StatusBadRequest,
				"deployment of resource %q is running, please wait for it to finish",
				r.Edges.Resource.Name,
			)
		default:
			return errorx.HttpErrorf(
				http.StatusBadRequest,
				"invalid deployment status of resource %q: %s",
				r.Edges.Resource.Name,
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
	resourceName string,
	projectID object.ID,
	environmentID object.ID,
) error {
	attrs := make(map[string]any, len(attributes))
	for k, v := range attributes {
		attrs[k] = string(json.ShouldMarshal(v))
	}

	opts := terraform.RevisionOpts{
		ResourceName:  resourceName,
		ProjectID:     projectID,
		EnvironmentID: environmentID,
	}
	_, _, err := terraform.ParseModuleAttributes(ctx, mc, attrs, true, opts)

	return err
}

func ValidateCreateInput(rci *model.ResourceCreateInput) error {
	// Resource type maps to type in definition edge.
	if rci.Type != "" {
		rci.ResourceDefinition = &model.ResourceDefinitionQueryInput{
			Type: rci.Type,
		}
	}

	if err := rci.Validate(); err != nil {
		return err
	}

	if err := validation.IsDNSLabel(rci.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	// Get environment.
	env, err := rci.Client.Environments().Query().
		Where(environment.ID(rci.Environment.ID)).
		Select(
			environment.FieldID,
			environment.FieldName).
		WithConnectors(func(rq *model.EnvironmentConnectorRelationshipQuery) {
			// Includes connectors.
			rq.WithConnector()
		}).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(project.FieldName)
		}).
		Only(rci.Context)
	if err != nil {
		return fmt.Errorf("failed to get environment: %w", err)
	}

	switch {
	case rci.Template != nil:
		// Get template version.
		tv, err := rci.Client.TemplateVersions().Query().
			Where(templateversion.ID(rci.Template.ID)).
			Select(
				templateversion.FieldID,
				templateversion.FieldName,
				templateversion.FieldSchema,
				templateversion.FieldUiSchema).
			Only(rci.Context)
		if err != nil {
			return fmt.Errorf("failed to get template version: %w", err)
		}

		// Validate template version whether match the target environment.
		if err = validateEnvironment(tv, env); err != nil {
			return err
		}

		// Verify variables with variables schema that defined on the template version.
		if !tv.Schema.IsEmpty() {
			err = rci.Attributes.ValidateWith(tv.Schema.VariableSchemas())
			if err != nil {
				return fmt.Errorf("invalid variables: %w", err)
			}
		}

	case rci.Type != "":
		rd, err := rci.Client.ResourceDefinitions().Query().
			Where(resourcedefinition.Type(rci.Type)).
			WithMatchingRules(func(rq *model.ResourceDefinitionMatchingRuleQuery) {
				rq.Order(model.Desc(resourcedefinitionmatchingrule.FieldCreateTime)).
					Select(resourcedefinitionmatchingrule.FieldResourceDefinitionID).
					Unique(false).
					Select(resourcedefinitionmatchingrule.FieldTemplateID).
					WithTemplate(func(tq *model.TemplateVersionQuery) {
						tq.Select(
							templateversion.FieldID,
							templateversion.FieldVersion,
							templateversion.FieldName,
						)
					})
			}).
			Select(resourcedefinition.FieldID, resourcedefinition.FieldName).
			Only(rci.Context)
		if err != nil {
			return fmt.Errorf("failed to get resource definition: %w", err)
		}

		rule := resourcedefinitions.Match(
			rd.Edges.MatchingRules,
			env.Edges.Project.Name,
			env.Name,
			env.Type,
			env.Labels,
			rci.Labels,
		)
		if rule == nil {
			return fmt.Errorf("resource definition %s does not match environment %s", rd.Name, env.Name)
		}
	default:
		return errors.New("template or resource definition is required")
	}

	// Verify that variables in attributes are valid.
	err = validateVariable(rci.Context, rci.Client, rci.Attributes, rci.Name, rci.Project.ID, rci.Environment.ID)
	if err != nil {
		return err
	}

	return nil
}
