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
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinition"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinitionmatchingrule"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerun"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/resourcedefinitions"
	runstatus "github.com/seal-io/walrus/pkg/resourceruns/status"
	pkgresource "github.com/seal-io/walrus/pkg/resources"
	"github.com/seal-io/walrus/pkg/resources/status"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/strs"
)

type AccessEndpoint struct {
	// Name is identifier for the endpoint.
	Name string `json:"name,omitempty"`
	// Endpoint is access endpoint.
	Endpoints []string `json:"endpoints,omitempty"`
}

type (
	RouteUpgradeRequest struct {
		_ struct{} `route:"PUT=/upgrade"`

		model.ResourceUpdateInput `path:",inline" json:",inline"`

		Draft           bool   `json:"draft,default=false"`
		ChangeComment   string `json:"changeComment,omitempty"`
		ReuseAttributes bool   `json:"reuseAttributes,default=false"`
		Preview         bool   `json:"preview,default=false"`
	}
)

func (r *RouteUpgradeRequest) Validate() error {
	if err := r.ResourceUpdateInput.Validate(); err != nil {
		return err
	}

	entity, err := r.Client.Resources().Query().
		Where(resource.ID(r.ID)).
		WithTemplate(func(tvq *model.TemplateVersionQuery) {
			tvq.Select(
				templateversion.FieldName,
				templateversion.FieldVersion)
		}).
		WithResourceDefinition(func(rdq *model.ResourceDefinitionQuery) {
			rdq.Select(
				resourcedefinition.FieldType,
			)
		}).
		Only(r.Context)
	if err != nil {
		return fmt.Errorf("failed to get resource: %w", err)
	}

	if r.Draft && !status.IsInactive(entity) {
		return errorx.HttpErrorf(http.StatusBadRequest,
			"cannot update resource draft in %q status", entity.Status.SummaryStatus)
	}

	if entity.TemplateID != nil && r.ReuseAttributes {
		r.Template = &model.TemplateVersionQueryInput{
			ID:   *entity.TemplateID,
			Name: entity.Edges.Template.Name,
		}
	}

	en := r.Model()
	// Set environment ID since Model will not set the environment ID.
	en.EnvironmentID = r.Environment.ID

	switch {
	case r.Template != nil:
		if r.Template.Name != entity.Edges.Template.Name {
			return errors.New("invalid template name: immutable")
		}

		tv, err := r.Client.TemplateVersions().Query().
			Where(templateversion.ID(r.Template.ID)).
			Select(
				templateversion.FieldSchema,
				templateversion.FieldUISchema,
			).
			Only(r.Context)
		if err != nil {
			return fmt.Errorf("failed to get template version: %w", err)
		}

		if err = validateAttributesWithTemplate(
			r.Context, r.Client, r.Project.ID, r.Environment.ID, r.Attributes, tv); err != nil {
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
				resourcedefinition.FieldUISchema,
			).
			WithMatchingRules(func(rq *model.ResourceDefinitionMatchingRuleQuery) {
				rq.Order(model.Asc(resourcedefinitionmatchingrule.FieldOrder)).
					Select(
						resourcedefinitionmatchingrule.FieldID,
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
			ProjectLabels:     env.Edges.Project.Labels,
			EnvironmentName:   env.Name,
			EnvironmentType:   env.Type,
			EnvironmentLabels: env.Labels,
			ResourceLabels:    r.Labels,
		})

		if def == nil {
			return errors.New("no matching resource definition found")
		}

		if err = validateAttributesWithResourceDefinition(
			r.Context, r.Client, r.Project.ID, r.Environment.ID, r.Attributes, def); err != nil {
			return err
		}

		r.ResourceDefinition = &model.ResourceDefinitionQueryInput{
			ID: def.ID,
		}

		r.ResourceDefinitionMatchingRule = &model.ResourceDefinitionMatchingRuleQueryInput{
			ID: rule.ID,
		}

		en.ResourceDefinitionMatchingRuleID = &rule.ID
	default:
		return errors.New("template or resource definition is required")
	}

	if r.ReuseAttributes {
		r.Attributes = entity.Attributes
		r.ComputedAttributes = entity.ComputedAttributes

		if r.Labels == nil {
			r.Labels = entity.Labels
		}
	} else {
		computedAttr, err := pkgresource.GenComputedAttributes(r.Context, r.Client, en)
		if err != nil {
			return err
		}

		r.ComputedAttributes = computedAttr
	}

	// Verify that variables in attributes are valid.
	if err = validateVariable(r.Context, r.Client, r.Attributes, r.Name, r.Project.ID, r.Environment.ID); err != nil {
		return err
	}

	if err = validateRunsStatus(r.Context, r.Client, r.ID); err != nil {
		return err
	}

	return nil
}

type RouteRollbackRequest struct {
	_ struct{} `route:"POST=/rollback"`

	model.ResourceQueryInput `path:",inline"`

	RunID object.ID `query:"runID"`

	ChangeComment string `json:"changeComment"`
	Preview       bool   `json:"preview,default=false"`
}

func (r *RouteRollbackRequest) Validate() error {
	if err := r.ResourceQueryInput.Validate(); err != nil {
		return err
	}

	if !r.RunID.Valid() {
		return errorx.New("run ID is required")
	}

	latestRun, err := r.Client.ResourceRuns().Query().
		Where(resourcerun.ResourceID(r.ID)).
		Order(model.Desc(resourcerun.FieldCreateTime)).
		Select(resourcerun.FieldStatus).
		First(r.Context)
	if err != nil && !model.IsNotFound(err) {
		return fmt.Errorf("failed to get the latest run: %w", err)
	}

	if runstatus.IsStatusRunning(latestRun) {
		return errors.New("latest run is running")
	}

	return nil
}

type RouteStopRequest struct {
	_ struct{} `route:"POST=/stop"`

	model.ResourceDeleteInput `path:",inline"`

	ChangeComment string `json:"changeComment"`
	Preview       bool   `json:"preview,default=false"`
}

func (r *RouteStopRequest) Validate() error {
	if err := r.ResourceDeleteInput.Validate(); err != nil {
		return err
	}

	res, err := r.Client.Resources().Get(r.Context, r.ID)
	if err != nil {
		return err
	}

	return validateStop(r.Context, r.Client, res)
}

func validateStop(ctx context.Context, mc model.ClientSet, resources ...*model.Resource) error {
	resourceNames := sets.NewString()
	for i := range resources {
		resourceNames.Insert(resources[i].Name)
	}

	for i := range resources {
		res := resources[i]
		if !pkgresource.IsStoppable(res) {
			return errorx.HttpErrorf(
				http.StatusBadRequest,
				"resource %s is non-stoppable",
				res.Name,
			)
		}

		if !pkgresource.CanBeStopped(res) {
			return errorx.HttpErrorf(
				http.StatusBadRequest,
				"cannot stop resource %q: in %q status",
				res.Name, res.Status.SummaryStatus,
			)
		}

		ids, err := dao.GetNonStoppedResourceDependantIDs(ctx, mc, res.ID)
		if err != nil {
			return fmt.Errorf("failed to get resource relationships: %w", err)
		}

		if len(ids) > 0 {
			names, err := dao.GetResourceNamesByIDs(ctx, mc, ids...)
			if err != nil {
				return fmt.Errorf("failed to get resources: %w", err)
			}

			// If the resource is the dependency of other non-stopped resources
			// and not exist in the input, then return error.
			if !resourceNames.HasAll(names...) {
				return errorx.HttpErrorf(
					http.StatusConflict,
					"resource about to be stopped is the dependency of: %v",
					strs.Join(", ", names...),
				)
			}
		}
	}

	return nil
}

type (
	RouteStartRequest struct {
		_ struct{} `route:"POST=/start"`

		model.ResourceQueryInput `path:",inline"`

		ChangeComment string `json:"changeComment"`
		Preview       bool   `json:"preview,default=false"`

		resource *model.Resource `json:"-"`
	}
)

func (r *RouteStartRequest) Validate() error {
	if err := r.ResourceQueryInput.Validate(); err != nil {
		return err
	}

	res, err := r.Client.Resources().Query().
		Where(resource.ID(r.ID)).
		WithTemplate(func(tvq *model.TemplateVersionQuery) {
			tvq.Select(
				templateversion.FieldID,
				templateversion.FieldTemplateID,
				templateversion.FieldName,
				templateversion.FieldVersion,
				templateversion.FieldProjectID)
		}).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(project.FieldName)
		}).
		WithEnvironment(func(eq *model.EnvironmentQuery) {
			eq.Select(environment.FieldName)
		}).
		Only(r.Context)
	if err != nil {
		return err
	}

	if !status.IsInactive(res) {
		return errorx.HttpErrorf(
			http.StatusBadRequest,
			"cannot start resource %q: in %q status",
			res.Name, res.Status.SummaryStatus,
		)
	}

	r.resource = res

	return nil
}

type (
	RouteGetEndpointsRequest struct {
		_ struct{} `route:"GET=/endpoints"`

		model.ResourceQueryInput `path:",inline"`

		Stream *runtime.RequestUnidiStream
	}

	RouteGetEndpointsResponse = types.ResourceEndpoints
)

func (r *RouteGetEndpointsRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type (
	RouteGetOutputsRequest struct {
		_ struct{} `route:"GET=/outputs"`

		model.ResourceQueryInput `path:",inline"`

		Stream *runtime.RequestUnidiStream
	}

	RouteGetOutputsResponse = []types.OutputValue
)

func (r *RouteGetOutputsRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type (
	// GraphVertexID defines the identifier of the vertex,
	// which uniquely represents an API resource.
	GraphVertexID = types.GraphVertexID
	// GraphVertex defines the vertex of graph.
	GraphVertex = types.GraphVertex
	// GraphEdge defines the edge of graph.
	GraphEdge = types.GraphEdge

	RouteGetGraphRequest struct {
		_ struct{} `route:"GET=/graph"`

		model.ResourceQueryInput `path:",inline"`

		WithoutKeys bool `query:"withoutKeys,omitempty"`
	}

	RouteGetGraphResponse struct {
		Vertices []GraphVertex `json:"vertices"`
		Edges    []GraphEdge   `json:"edges"`
	}
)

type StartInputs = model.ResourceDeleteInputs

type CollectionRouteStartRequest struct {
	_ struct{} `route:"POST=/start"`

	StartInputs `path:",inline" json:",inline"`

	ChangeComment string `json:"changeComment"`
	Preview       bool   `json:"preview,default=false"`

	Resources []*model.Resource `json:"-"`
}

func (r *CollectionRouteStartRequest) Validate() error {
	if err := r.StartInputs.Validate(); err != nil {
		return err
	}

	entities, err := r.Client.Resources().Query().
		Where(resource.IDIn(r.IDs()...)).
		WithTemplate(func(tvq *model.TemplateVersionQuery) {
			tvq.Select(
				templateversion.FieldID,
				templateversion.FieldTemplateID,
				templateversion.FieldName,
				templateversion.FieldVersion,
				templateversion.FieldProjectID)
		}).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(project.FieldName)
		}).
		WithEnvironment(func(eq *model.EnvironmentQuery) {
			eq.Select(environment.FieldName)
		}).
		All(r.Context)
	if err != nil {
		return err
	}

	for i := range entities {
		entity := entities[i]
		if !status.IsInactive(entity) {
			return errorx.HttpErrorf(
				http.StatusBadRequest,
				"cannot start resource %q: in %q status",
				entity.Name, entity.Status.SummaryStatus,
			)
		}
	}

	r.Resources = entities

	return nil
}

type CollectionRouteStopRequest struct {
	_ struct{} `route:"POST=/stop"`

	model.ResourceDeleteInputs `path:",inline" json:",inline"`

	ChangeComment string `json:"changeComment"`
	Preview       bool   `json:"preview,default=false"`

	Resources []*model.Resource `json:"-"`
}

func (r *CollectionRouteStopRequest) Validate() error {
	if err := r.ResourceDeleteInputs.Validate(); err != nil {
		return err
	}

	entities, err := r.Client.Resources().Query().
		Where(resource.IDIn(r.IDs()...)).
		All(r.Context)
	if err != nil {
		return err
	}

	r.Resources = entities

	return validateStop(r.Context, r.Client, entities...)
}

type CollectionRouteUpgradeRequest struct {
	_ struct{} `route:"POST=/upgrade"`

	model.ResourceUpdateInputs `path:",inline" json:",inline"`

	ChangeComment   string `json:"changeComment"`
	Draft           bool   `json:"draft,default=false"`
	ReuseAttributes bool   `json:"reuseAttributes,default=false"`
	Preview         bool   `json:"preview,default=false"`
}

func (r *CollectionRouteUpgradeRequest) Validate() error {
	if err := r.ResourceUpdateInputs.Validate(); err != nil {
		return err
	}

	var err error

	// Fetch database entities to refill input items.
	var (
		entities    []*model.Resource
		entitiesMap = make(map[object.ID]*model.Resource)
	)
	{
		entities, err = r.Client.Resources().Query().
			Where(resource.IDIn(r.IDs()...)).
			WithProject(func(pq *model.ProjectQuery) {
				pq.Select(
					project.FieldID,
					project.FieldName,
					project.FieldLabels)
			}).
			WithEnvironment(func(eq *model.EnvironmentQuery) {
				eq.Select(
					environment.FieldID,
					environment.FieldName,
					environment.FieldLabels,
					environment.FieldType)
			}).
			WithTemplate(func(tvq *model.TemplateVersionQuery) {
				tvq.Select(
					templateversion.FieldID,
					templateversion.FieldName,
					templateversion.FieldVersion)
			}).
			Select(
				resource.WithoutFields(
					resource.FieldUpdateTime,
					resource.FieldCreateTime,
					resource.FieldComputedAttributes)...).
			All(r.Context)
		if err != nil {
			return fmt.Errorf("failed to get resources: %w", err)
		}

		for i := range entities {
			entitiesMap[entities[i].ID] = entities[i]

			if r.Draft && !status.IsInactive(entities[i]) {
				return errorx.HttpErrorf(http.StatusBadRequest,
					"cannot update resource draft in %q status", entities[i].Status.SummaryStatus)
			}
		}
	}

	// Refill input items if reusing.
	if r.ReuseAttributes {
		for i := range r.Items {
			input := r.Items[i]

			entity, ok := entitiesMap[input.ID]
			if !ok {
				return fmt.Errorf("resource %s not found", input.Name)
			}

			// Refill attributes.
			input.Attributes = entity.Attributes

			// Reuse template if exists.
			if entity.TemplateID != nil {
				input.Template = &model.TemplateVersionQueryInput{
					ID:   *entity.TemplateID,
					Name: entity.Edges.Template.Name,
				}

				continue
			}

			// Otherwise, reuse resource definition type.
			input.Type = entity.Type

			// NB(thxCode): Refill labels for resource definition matching.
			input.Labels = entity.Labels
		}
	}

	// Validate input items.
	defsCache := make(map[string][]*model.ResourceDefinition)

	for i := range r.Items {
		input := r.Items[i]

		entity, ok := entitiesMap[input.ID]
		if !ok {
			return fmt.Errorf("resource %s not found", input.Name)
		}

		switch {
		default:
			return errors.New("template or resource definition is required")
		case input.Template != nil:
			// Validate if the template has changed.
			if input.Template.Name != entity.Edges.Template.Name {
				return errors.New("invalid template name: immutable")
			}

			tv, err := r.Client.TemplateVersions().Query().
				Where(templateversion.ID(input.Template.ID)).
				Select(
					templateversion.FieldSchema,
					templateversion.FieldUISchema).
				Only(r.Context)
			if err != nil {
				return fmt.Errorf("failed to get template version: %w", err)
			}

			// Validate attributes with template schema.
			err = validateAttributesWithTemplate(
				r.Context, r.Client,
				entity.ProjectID, entity.EnvironmentID,
				input.Attributes, tv)
			if err != nil {
				return err
			}
		case input.Type != "":
			// Get resource definitions by type.
			rfs, ok := defsCache[input.Type]
			if !ok {
				rfs, err = r.Client.ResourceDefinitions().Query().
					Where(resourcedefinition.Type(input.Type)).
					Select(
						resourcedefinition.FieldID,
						resourcedefinition.FieldName,
						resourcedefinition.FieldType,
						resourcedefinition.FieldSchema,
						resourcedefinition.FieldUISchema,
					).
					WithMatchingRules(func(rq *model.ResourceDefinitionMatchingRuleQuery) {
						rq.Order(model.Asc(resourcedefinitionmatchingrule.FieldOrder)).
							Select(
								resourcedefinitionmatchingrule.FieldName,
								resourcedefinitionmatchingrule.FieldSelector)
					}).
					All(r.Context)
				if err != nil {
					return fmt.Errorf("failed to get resource definitions: %w", err)
				}
				defsCache[input.Type] = rfs // Reuse.
			}

			// Pick rules.
			rf, rule := resourcedefinitions.MatchResourceDefinition(
				rfs,
				types.MatchResourceMetadata{
					ProjectName:       entity.Edges.Project.Name,
					EnvironmentName:   entity.Edges.Environment.Name,
					EnvironmentType:   entity.Edges.Environment.Type,
					ProjectLabels:     entity.Edges.Project.Labels,
					EnvironmentLabels: entity.Edges.Environment.Labels,
					ResourceLabels:    input.Labels,
				})
			if rf == nil {
				return fmt.Errorf("no matching resource definition found for resource %s", input.Name)
			}

			// Validate attributes with resource definition schema.
			err = validateAttributesWithResourceDefinition(
				r.Context, r.Client,
				entity.ProjectID, entity.EnvironmentID,
				input.Attributes, rf)
			if err != nil {
				return err
			}

			// Get matched result for calculating computed attributes.
			input.ResourceDefinitionMatchingRule = &model.ResourceDefinitionMatchingRuleQueryInput{
				ID: rule.ID,
			}
		}

		// Validate referring variables.
		err = validateVariable(
			r.Context, r.Client,
			input.Attributes, input.Name,
			entity.ProjectID, entity.EnvironmentID)
		if err != nil {
			return err
		}

		// Validate whether the resource is deploying.
		err = validateRunsStatus(
			r.Context, r.Client,
			input.ID)
		if err != nil {
			return err
		}
	}

	// Get computed attributes for deploying.
	for i := range r.Items {
		input := r.Items[i]
		entity := entitiesMap[input.ID] // Guaranteed to exist.

		switch {
		case input.Template != nil:
			entity.TemplateID = &input.Template.ID
		case input.Type != "":
			entity.ResourceDefinitionMatchingRuleID = &input.ResourceDefinitionMatchingRule.ID
		}

		attrs, err := pkgresource.GenComputedAttributes(r.Context, r.Client, entity)
		if err != nil {
			return err
		}

		input.ComputedAttributes = attrs
	}

	return nil
}
