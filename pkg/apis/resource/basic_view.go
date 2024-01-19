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
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/deployer/terraform"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
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

		Draft bool `json:"draft,default=false"`
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

		return fmt.Errorf(
			"resource about to be deleted is the dependency of: %v",
			strs.Join(", ", names...),
		)
	}

	if !r.WithoutCleanup {
		if err = ValidateRevisionsStatus(r.Context, r.Client, r.ID); err != nil {
			return err
		}
	}

	return nil
}

type PatchRequest struct {
	model.ResourcePatchInput `path:",inline" json:",inline"`

	Draft bool `json:"draft,default=false"`
}

func (r *PatchRequest) Validate() error {
	if err := r.ResourcePatchInput.Validate(); err != nil {
		return err
	}

	entity, err := r.Client.Resources().Query().
		Where(resource.ID(r.ID)).
		Select(
			resource.FieldTemplateID,
			resource.FieldResourceDefinitionID,
			resource.FieldType,
			resource.FieldStatus,
		).
		WithTemplate(func(tvq *model.TemplateVersionQuery) {
			tvq.Select(
				templateversion.FieldName,
				templateversion.FieldUiSchema,
				templateversion.FieldSchema,
				templateversion.FieldVersion)
		}).
		Only(r.Context)
	if err != nil {
		return fmt.Errorf("failed to get resource: %w", err)
	}

	if r.Draft && !pkgresource.IsInactive(entity) {
		return errorx.HttpErrorf(http.StatusBadRequest,
			"cannot update resource draft in %q status", entity.Status.SummaryStatus)
	}

	patched := r.Model()
	patched.Edges = entity.Edges

	switch {
	case patched.TemplateID != nil:
		if patched.TemplateID.String() != entity.TemplateID.String() {
			return errors.New("invalid template: immutable")
		}

		err = validateAttributesWithTemplate(r.Attributes, patched.Edges.Template)
		if err != nil {
			return err
		}
	case entity.Type != "":
		env, err := r.Client.Environments().Query().
			Where(environment.ID(r.Environment.ID)).
			WithProject(func(pq *model.ProjectQuery) {
				pq.Select(project.FieldName, project.FieldLabels)
			}).
			Only(r.Context)
		if err != nil {
			return fmt.Errorf("failed to get environment: %w", err)
		}

		resourceDefinitions, err := r.Client.ResourceDefinitions().Query().
			Where(resourcedefinition.Type(entity.Type)).
			Select(
				resourcedefinition.FieldID,
				resourcedefinition.FieldName,
				resourcedefinition.FieldType,
				resourcedefinition.FieldSchema,
				resourcedefinition.FieldUiSchema,
			).
			WithMatchingRules(func(rq *model.ResourceDefinitionMatchingRuleQuery) {
				rq.Order(model.Asc(resourcedefinitionmatchingrule.FieldOrder)).
					Select(
						resourcedefinitionmatchingrule.FieldName,
						resourcedefinitionmatchingrule.FieldSelector,
					)
			}).
			All(r.Context)
		if err != nil {
			return fmt.Errorf("failed to get resource definitions: %w", err)
		}

		def, rule := resourcedefinitions.MatchResourceDefinition(resourceDefinitions, types.MatchResourceMetadata{
			ProjectName:       env.Edges.Project.Name,
			EnvironmentName:   env.Name,
			EnvironmentType:   env.Type,
			ProjectLabels:     env.Edges.Project.Labels,
			EnvironmentLabels: env.Labels,
			ResourceLabels:    patched.Labels,
		})

		if def == nil {
			return errors.New("find no matching resource definition")
		}

		// Mutate definition edge according to matching resource definition.
		// The matching definition/rule may change on update.
		r.ResourceDefinition = &model.ResourceDefinitionQueryInput{
			ID: def.ID,
		}

		r.ResourceDefinitionMatchingRule = &model.ResourceDefinitionMatchingRuleQueryInput{
			ID: rule.ID,
		}

		if err = validateAttributesWithResourceDefinition(r.Attributes, def); err != nil {
			return err
		}
	default:
		return errors.New("invalid resource: missing type or template")
	}

	// Verify that variables in attributes are valid.
	if err = validateVariable(
		r.Context,
		r.Client,
		patched.Attributes,
		patched.Name,
		patched.ProjectID,
		patched.EnvironmentID); err != nil {
		return err
	}

	if err = ValidateRevisionsStatus(r.Context, r.Client, patched.ID); err != nil {
		return err
	}

	// Set computedAttributes.
	patched.ComputedAttributes, err = pkgresource.GenComputedAttributes(r.Context, r.Client, patched)
	if err != nil {
		return err
	}

	return nil
}

type (
	CollectionCreateRequest struct {
		model.ResourceCreateInputs `path:",inline" json:",inline"`

		Draft bool `json:"draft,default=false"`
	}

	CollectionCreateResponse = []*model.ResourceOutput
)

func (r *CollectionCreateRequest) Validate() error {
	if err := r.ResourceCreateInputs.Validate(); err != nil {
		return err
	}

	cache := make(map[string][]*model.ResourceDefinition)

	env, err := r.Client.Environments().Query().
		Where(environment.ID(r.Environment.ID)).
		WithConnectors(func(rq *model.EnvironmentConnectorRelationshipQuery) {
			// Includes connectors.
			rq.WithConnector()
		}).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(
				project.FieldName,
				project.FieldLabels,
			)
		}).
		Only(r.Context)
	if err != nil {
		return fmt.Errorf("failed to get environment: %w", err)
	}

	for _, rci := range r.Items {
		// Mutate definition edge according to matching resource definition.
		if rci.Type != "" {
			definitions, ok := cache[rci.Type]
			if !ok {
				definitions, err = r.Client.ResourceDefinitions().Query().
					Where(resourcedefinition.Type(rci.Type)).
					Select(
						resourcedefinition.FieldID,
						resourcedefinition.FieldName,
						resourcedefinition.FieldType,
					).
					WithMatchingRules(func(rq *model.ResourceDefinitionMatchingRuleQuery) {
						rq.Order(model.Asc(resourcedefinitionmatchingrule.FieldOrder)).
							Select(
								resourcedefinitionmatchingrule.FieldName,
								resourcedefinitionmatchingrule.FieldSelector,
							)
					}).
					All(r.Context)
				if err != nil {
					return fmt.Errorf("failed to get resource definitions: %w", err)
				}
				cache[rci.Type] = definitions
			}

			def, rule := resourcedefinitions.MatchResourceDefinition(definitions, types.MatchResourceMetadata{
				ProjectName:       env.Edges.Project.Name,
				EnvironmentName:   env.Name,
				EnvironmentType:   env.Type,
				ProjectLabels:     env.Edges.Project.Labels,
				EnvironmentLabels: env.Labels,
				ResourceLabels:    rci.Labels,
			})

			if def == nil {
				return fmt.Errorf("find no mathcing resource definition for resource %s", rci.Name)
			}

			rci.ResourceDefinition = &model.ResourceDefinitionQueryInput{
				ID: def.ID,
			}

			rci.ResourceDefinitionMatchingRule = &model.ResourceDefinitionMatchingRuleQueryInput{
				ID: rule.ID,
			}
		}
	}

	// Verify resources.
	for i := range r.Items {
		if r.Items[i] == nil {
			return errors.New("empty resource")
		}

		if err := validation.IsValidName(r.Items[i].Name); err != nil {
			return fmt.Errorf("invalid resource name: %w", err)
		}
	}

	tvIDs := make([]object.ID, 0, len(r.Items))
	// Get template versions.
	for i := range r.Items {
		if r.Items[i].Template == nil {
			continue
		}

		tvIDs = append(tvIDs, r.Items[i].Template.ID)
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

	for _, res := range r.Items {
		if res.Template != nil {
			err = validateAttributesWithTemplate(res.Attributes, tvm[res.Template.ID])
			if err != nil {
				return err
			}
		}

		// Verify that variables in attributes are valid.
		err = validateVariable(r.Context, r.Client, res.Attributes, res.Name, r.Project.ID, r.Environment.ID)
		if err != nil {
			return err
		}

		// Set computedAttributes.
		var err error
		en := createInputsItemToResource(res, r.Project, r.Environment)

		res.ComputedAttributes, err = pkgresource.GenComputedAttributes(r.Context, r.Client, en)
		if err != nil {
			return err
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

		return fmt.Errorf(
			"resource about to be deleted is the dependency of: %v",
			strs.Join(", ", names...),
		)
	}

	if r.WithoutCleanup {
		if err = ValidateRevisionsStatus(r.Context, r.Client, ids...); err != nil {
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

// ValidateRevisionsStatus validates revision status of given resource IDs.
func ValidateRevisionsStatus(ctx context.Context, mc model.ClientSet, ids ...object.ID) error {
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
	_, _, _, err := terraform.ParseModuleAttributes(ctx, mc, attrs, true, opts)

	return err
}

func ValidateCreateInput(rci *model.ResourceCreateInput) error {
	if err := rci.Validate(); err != nil {
		return err
	}

	if err := validation.IsValidName(rci.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	// Get environment.
	env, err := rci.Client.Environments().Query().
		Where(environment.ID(rci.Environment.ID)).
		WithConnectors(func(rq *model.EnvironmentConnectorRelationshipQuery) {
			// Includes connectors.
			rq.WithConnector()
		}).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(project.FieldName, project.FieldLabels)
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

		err = validateAttributesWithTemplate(rci.Attributes, tv)
		if err != nil {
			return err
		}

	case rci.Type != "":
		resourceDefinitions, err := rci.Client.ResourceDefinitions().Query().
			Where(resourcedefinition.Type(rci.Type)).
			Select(
				resourcedefinition.FieldID,
				resourcedefinition.FieldName,
				resourcedefinition.FieldType,
				resourcedefinition.FieldSchema,
				resourcedefinition.FieldUiSchema,
			).
			WithMatchingRules(func(rq *model.ResourceDefinitionMatchingRuleQuery) {
				rq.Order(model.Asc(resourcedefinitionmatchingrule.FieldOrder)).
					Select(
						resourcedefinitionmatchingrule.FieldID,
						resourcedefinitionmatchingrule.FieldName,
						resourcedefinitionmatchingrule.FieldSelector,
					)
			}).
			All(rci.Context)
		if err != nil {
			return fmt.Errorf("failed to get resource definitions: %w", err)
		}

		def, rule := resourcedefinitions.MatchResourceDefinition(resourceDefinitions, types.MatchResourceMetadata{
			ProjectName:       env.Edges.Project.Name,
			EnvironmentName:   env.Name,
			EnvironmentType:   env.Type,
			ProjectLabels:     env.Edges.Project.Labels,
			EnvironmentLabels: env.Labels,
			ResourceLabels:    rci.Labels,
		})
		if def == nil {
			return errors.New("find no matching resource definition")
		}

		if err = validateAttributesWithResourceDefinition(rci.Attributes, def); err != nil {
			return err
		}

		rci.ResourceDefinition = &model.ResourceDefinitionQueryInput{
			ID: def.ID,
		}

		rci.ResourceDefinitionMatchingRule = &model.ResourceDefinitionMatchingRuleQueryInput{
			ID: rule.ID,
		}
	default:
		return errors.New("invalid resource: missing type or template")
	}

	// Verify that variables in attributes are valid.
	err = validateVariable(rci.Context, rci.Client, rci.Attributes, rci.Name, rci.Project.ID, rci.Environment.ID)
	if err != nil {
		return err
	}

	// Set computedAttributes.
	en := rci.Model()

	rci.ComputedAttributes, err = pkgresource.GenComputedAttributes(rci.Context, rci.Client, en)
	if err != nil {
		return err
	}

	return nil
}

func validateAttributesWithTemplate(attrs property.Values, tv *model.TemplateVersion) error {
	if s := tv.UiSchema; !s.IsEmpty() {
		err := attrs.ValidateWith(s.VariableSchema())
		if err != nil {
			return fmt.Errorf("invalid variables: violate ui schema: %w", err)
		}
	}

	return nil
}

func validateAttributesWithResourceDefinition(attrs property.Values, rd *model.ResourceDefinition) error {
	rdo := dao.ExposeResourceDefinition(rd)
	if s := rdo.UiSchema; !s.IsEmpty() {
		err := attrs.ValidateWith(s.VariableSchema())
		if err != nil {
			return fmt.Errorf("invalid variables: violate ui schema: %w", err)
		}
	}

	return nil
}
