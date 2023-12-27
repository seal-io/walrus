package resource

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinition"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinitionmatchingrule"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
	"github.com/seal-io/walrus/pkg/resourcedefinitions"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/strs"
)

type AccessEndpoint struct {
	// Name is identifier for the endpoint.
	Name string `json:"name,omitempty"`
	// Endpoint is access endpoint.
	Endpoints []string `json:"endpoints,omitempty"`
}

type RouteUpgradeRequest struct {
	_ struct{} `route:"PUT=/upgrade"`

	model.ResourceUpdateInput `path:",inline" json:",inline"`

	Draft           bool `json:"draft,default=false"`
	ReuseAttributes bool `json:"reuseAttributes,default=false"`
}

func (r *RouteUpgradeRequest) Validate() error {
	if err := r.ResourceUpdateInput.Validate(); err != nil {
		return err
	}

	entity, err := r.Client.Resources().Query().
		Where(resource.ID(r.ID)).
		Select(
			resource.FieldTemplateID,
			resource.FieldResourceDefinitionID,
			resource.FieldStatus,
			resource.FieldAttributes,
		).
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

	if r.Draft && !pkgresource.IsInactive(entity) {
		return errorx.HttpErrorf(http.StatusBadRequest,
			"cannot update resource draft in %q status", entity.Status.SummaryStatus)
	}

	if r.ReuseAttributes {
		r.Attributes = entity.Attributes
	}

	if entity.ResourceDefinitionID != nil {
		r.ResourceDefinition = &model.ResourceDefinitionQueryInput{
			Type: entity.Edges.ResourceDefinition.Type,
			ID:   *entity.ResourceDefinitionID,
		}
	}

	if entity.TemplateID != nil && r.ReuseAttributes {
		r.Template = &model.TemplateVersionQueryInput{
			ID:   *entity.TemplateID,
			Name: entity.Edges.Template.Name,
		}
	}

	switch {
	case r.Template != nil:
		if r.Template.Name != entity.Edges.Template.Name {
			return errors.New("invalid template name: immutable")
		}

		tv, err := r.Client.TemplateVersions().Query().
			Where(templateversion.ID(r.Template.ID)).
			Select(
				templateversion.FieldSchema,
				templateversion.FieldUiSchema,
			).
			Only(r.Context)
		if err != nil {
			return fmt.Errorf("failed to get template version: %w", err)
		}

		// Verify attributes with schema.
		// TODO(thxCode): migrate schema to ui schema, then reduce if-else.
		if s := tv.UiSchema; !s.IsEmpty() {
			err = r.Attributes.ValidateWith(s.VariableSchema())
			if err != nil {
				return fmt.Errorf("invalid variables: violate ui schema: %w", err)
			}
		} else if s := tv.Schema; !s.IsEmpty() {
			err = r.Attributes.ValidateWith(s.VariableSchema())
			if err != nil {
				return fmt.Errorf("invalid variables: %w", err)
			}
		}
	case r.ResourceDefinition != nil:
		rd, err := r.Client.ResourceDefinitions().Query().
			Where(resourcedefinition.Type(r.ResourceDefinition.Type)).
			WithMatchingRules(func(rq *model.ResourceDefinitionMatchingRuleQuery) {
				rq.Order(model.Asc(resourcedefinitionmatchingrule.FieldOrder)).
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
			Only(r.Context)
		if err != nil {
			return fmt.Errorf("failed to get resource definition: %w", err)
		}

		env, err := r.Client.Environments().Query().
			Where(environment.ID(r.Environment.ID)).
			Select(
				environment.FieldID,
				environment.FieldName,
				environment.FieldLabels,
			).
			WithProject(func(pq *model.ProjectQuery) {
				pq.Select(project.FieldName)
			}).
			Only(r.Context)
		if err != nil {
			return fmt.Errorf("failed to get environment: %w", err)
		}

		rule := resourcedefinitions.Match(
			rd.Edges.MatchingRules,
			env.Edges.Project.Name,
			env.Name,
			env.Type,
			env.Labels,
			r.Labels,
		)
		if rule == nil {
			return fmt.Errorf("resource definition %s does not match environment %s", rd.Name, env.Name)
		}
	default:
		return errors.New("template or resource definition is required")
	}

	// Verify that variables in attributes are valid.
	if err = validateVariable(r.Context, r.Client, r.Attributes, r.Name, r.Project.ID, r.Environment.ID); err != nil {
		return err
	}

	if err = ValidateRevisionsStatus(r.Context, r.Client, r.ID); err != nil {
		return err
	}

	return nil
}

type RouteRollbackRequest struct {
	_ struct{} `route:"POST=/rollback"`

	model.ResourceQueryInput `path:",inline"`

	RevisionID    object.ID `query:"revisionID"`
	ChangeComment string    `query:"changeComment"`
}

func (r *RouteRollbackRequest) Validate() error {
	if err := r.ResourceQueryInput.Validate(); err != nil {
		return err
	}

	latestRevision, err := r.Client.ResourceRevisions().Query().
		Where(resourcerevision.ResourceID(r.ID)).
		Order(model.Desc(resourcerevision.FieldCreateTime)).
		Select(resourcerevision.FieldStatus).
		First(r.Context)
	if err != nil && !model.IsNotFound(err) {
		return fmt.Errorf("failed to get the latest revision: %w", err)
	}

	if status.ResourceRevisionStatusReady.IsUnknown(latestRevision) {
		return errors.New("latest revision is running")
	}

	return nil
}

type RouteStopRequest struct {
	_ struct{} `route:"POST=/stop"`

	model.ResourceDeleteInput `path:",inline"`

	ChangeComment string `json:"changeComment"`
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

			return errorx.HttpErrorf(
				http.StatusConflict,
				"resource about to be stopped is the dependency of: %v",
				strs.Join(", ", names...),
			)
		}
	}

	return nil
}

type RouteStartRequest struct {
	_ struct{} `route:"POST=/start"`

	model.ResourceQueryInput `path:",inline"`

	ChangeComment string `json:"changeComment"`

	resource *model.Resource `json:"-"`
}

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

	if !pkgresource.IsInactive(res) {
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
	RouteGetAccessEndpointsRequest struct {
		_ struct{} `route:"GET=/access-endpoints"`

		model.ResourceQueryInput `path:",inline"`

		Stream *runtime.RequestUnidiStream
	}

	RouteGetAccessEndpointsResponse = []AccessEndpoint
)

func (r *RouteGetAccessEndpointsRequest) SetStream(stream runtime.RequestUnidiStream) {
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
		if !pkgresource.IsInactive(entity) {
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
}

func (r *CollectionRouteUpgradeRequest) Validate() error {
	if err := r.ResourceUpdateInputs.Validate(); err != nil {
		return err
	}

	entities, err := r.Client.Resources().Query().
		Where(resource.IDIn(r.IDs()...)).
		Select(
			resource.FieldTemplateID,
			resource.FieldResourceDefinitionID,
			resource.FieldStatus,
			resource.FieldAttributes,
		).
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
		All(r.Context)
	if err != nil {
		return fmt.Errorf("failed to get resources: %w", err)
	}

	entityMap := make(map[object.ID]*model.Resource)
	for i := range entities {
		entityMap[entities[i].ID] = entities[i]

		if r.Draft && !pkgresource.IsInactive(entities[i]) {
			return errorx.HttpErrorf(http.StatusBadRequest,
				"cannot update resource draft in %q status", entities[i].Status.SummaryStatus)
		}
	}

	for i := range r.Items {
		entity, ok := entityMap[r.Items[i].ID]
		if !ok {
			return fmt.Errorf("resource %s not found", r.Items[i].Name)
		}

		input := r.Items[i]

		if r.ReuseAttributes {
			input.Attributes = entity.Attributes
		}

		if entity.ResourceDefinitionID != nil {
			input.ResourceDefinition = &model.ResourceDefinitionQueryInput{
				Type: entity.Edges.ResourceDefinition.Type,
				ID:   *entity.ResourceDefinitionID,
			}
		}

		if r.ReuseAttributes && entity.TemplateID != nil {
			input.Template = &model.TemplateVersionQueryInput{
				ID:   *entity.TemplateID,
				Name: entity.Edges.Template.Name,
			}
		}

		switch {
		case input.Template != nil:
			if input.Template.Name != entity.Edges.Template.Name {
				return errors.New("invalid template name: immutable")
			}

			tv, err := r.Client.TemplateVersions().Query().
				Where(templateversion.ID(input.Template.ID)).
				Select(
					templateversion.FieldSchema,
					templateversion.FieldUiSchema,
				).
				Only(r.Context)
			if err != nil {
				return fmt.Errorf("failed to get template version: %w", err)
			}

			// Verify attributes with schema.
			// TODO(thxCode): migrate schema to ui schema, then reduce if-else.
			if s := tv.UiSchema; !s.IsEmpty() {
				err = input.Attributes.ValidateWith(s.VariableSchema())
				if err != nil {
					return fmt.Errorf("invalid variables: violate ui schema: %w", err)
				}
			} else if s := tv.Schema; !s.IsEmpty() {
				err = input.Attributes.ValidateWith(s.VariableSchema())
				if err != nil {
					return fmt.Errorf("invalid variables: %w", err)
				}
			}
		case input.ResourceDefinition != nil:
			rd, err := r.Client.ResourceDefinitions().Query().
				Where(resourcedefinition.Type(input.ResourceDefinition.Type)).
				WithMatchingRules(func(rq *model.ResourceDefinitionMatchingRuleQuery) {
					rq.Order(model.Asc(resourcedefinitionmatchingrule.FieldOrder)).
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
				Only(r.Context)
			if err != nil {
				return fmt.Errorf("failed to get resource definition: %w", err)
			}

			env, err := r.Client.Environments().Query().
				Where(environment.ID(r.Environment.ID)).
				Select(
					environment.FieldID,
					environment.FieldName,
					environment.FieldLabels,
				).
				WithProject(func(pq *model.ProjectQuery) {
					pq.Select(project.FieldName)
				}).
				Only(r.Context)
			if err != nil {
				return fmt.Errorf("failed to get environment: %w", err)
			}

			rule := resourcedefinitions.Match(
				rd.Edges.MatchingRules,
				env.Edges.Project.Name,
				env.Name,
				env.Type,
				env.Labels,
				input.Labels,
			)
			if rule == nil {
				return fmt.Errorf("resource definition %s does not match environment %s", rd.Name, env.Name)
			}
		default:
			return errors.New("template or resource definition is required")
		}

		// Verify that variables in attributes are valid.
		if err := validateVariable(r.Context, r.Client,
			input.Attributes, input.Name, r.Project.ID, r.Environment.ID); err != nil {
			return err
		}

		if err := ValidateRevisionsStatus(r.Context, r.Client, input.ID); err != nil {
			return err
		}
	}

	return nil
}
