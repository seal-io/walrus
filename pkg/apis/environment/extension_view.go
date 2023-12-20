package environment

import (
	"fmt"
	"net/http"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/strs"
)

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

		model.EnvironmentQueryInput `path:",inline"`

		WithoutKeys bool `query:"withoutKeys,omitempty"`
	}

	RouteGetGraphResponse struct {
		Vertices []GraphVertex `json:"vertices"`
		Edges    []GraphEdge   `json:"edges"`
	}
)

type (
	RouteCloneEnvironmentRequest struct {
		_ struct{} `route:"POST=/clone"`

		model.EnvironmentCreateInput `path:",inline" json:",inline"`

		CloneEnvironmentId object.ID `path:"environment"`

		// When draft is true, clone given resources as undeployed draft in target environment.
		Draft bool `json:"draft,default=false"`
	}

	RouteCloneEnvironmentResponse = model.EnvironmentOutput
)

func (r *RouteCloneEnvironmentRequest) Validate() error {
	return validateEnvironmentCreateInput(r.EnvironmentCreateInput)
}

type (
	RouteGetResourceDefinitionsRequest struct {
		_ struct{} `route:"GET=/resource-definitions"`

		model.EnvironmentQueryInput `path:",inline"`
	}

	RouteGetResourceDefinitionsResponse = []*model.ResourceDefinitionOutput
)

type RouteStopRequest struct {
	_ struct{} `route:"POST=/stop"`

	model.EnvironmentQueryInput `path:",inline"`

	stoppableResources model.Resources `json:"-"`
}

func (r *RouteStopRequest) Validate() error {
	if err := r.EnvironmentQueryInput.Validate(); err != nil {
		return err
	}

	resources, err := r.Client.Resources().Query().
		Where(resource.EnvironmentID(r.ID)).
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

	stoppableResources := make(model.Resources, 0, len(resources))
	stoppableResourceIDs := make([]object.ID, 0, len(resources))

	for _, r := range resources {
		if pkgresource.IsStoppable(r) {
			stoppableResources = append(stoppableResources, r)
			stoppableResourceIDs = append(stoppableResourceIDs, r.ID)
		}
	}

	r.stoppableResources = stoppableResources

	dependantIDs, err := dao.GetNonStoppedResourceDependantIDs(r.Context, r.Client, stoppableResourceIDs...)
	if err != nil {
		return fmt.Errorf("failed to get resource dependencies: %w", err)
	}

	dependantIDSet := sets.New[object.ID](dependantIDs...)
	toStopIDSet := sets.New[object.ID](stoppableResourceIDs...)

	// Validate if a resource is about to stop but the resources depending on it is not.
	diffIDSet := dependantIDSet.Difference(toStopIDSet)
	if diffIDSet.Len() > 0 {
		names, err := dao.GetResourceNamesByIDs(r.Context, r.Client, diffIDSet.UnsortedList()...)
		if err != nil {
			return fmt.Errorf("failed to get resources: %w", err)
		}

		return errorx.HttpErrorf(
			http.StatusConflict,
			"resource about to be stopped is the dependency of: %v",
			strs.Join(", ", names...),
		)
	}

	return nil
}

type RouteStartRequest struct {
	_ struct{} `route:"POST=/start"`

	model.EnvironmentQueryInput `path:",inline"`
}

func (r *RouteStartRequest) Validate() error {
	if err := r.EnvironmentQueryInput.Validate(); err != nil {
		return err
	}

	return nil
}

type RouteApplyRequest struct {
	_ struct{} `route:"POST=/apply"`

	model.EnvironmentQueryInput `path:",inline"`

	YAML string `json:"yaml"`
}

func (r *RouteApplyRequest) Validate() error {
	if err := r.EnvironmentQueryInput.Validate(); err != nil {
		return err
	}

	return nil
}

func serverContext(project, env, token string) *config.Config {
	return &config.Config{
		ServerContext: config.ServerContext{
			ScopeContext: config.ScopeContext{
				Project:     project,
				Environment: env,
			},
			Server:   "https://localhost",
			Insecure: true,
			Token:    token,
		},
	}
}

type RouteExportRequest struct {
	_ struct{} `route:"GET=/export"`

	model.EnvironmentQueryInput `path:",inline"`

	Names []string `query:"name,omitempty"`

	IDs []object.ID `query:"id,omitempty"`
}

func (r *RouteExportRequest) Validate() error {
	if err := r.EnvironmentQueryInput.Validate(); err != nil {
		return err
	}

	if len(r.Names) == 0 && len(r.IDs) == 0 {
		return fmt.Errorf("resource identidy is empty")
	}

	// While id existed.
	switch {
	case len(r.IDs) > 0:
		for i := range r.IDs {
			rq := model.ResourceQueryInput{
				ID: r.IDs[i],
				Environment: &model.EnvironmentQueryInput{
					ID: r.EnvironmentQueryInput.ID,
				},
				Project: &model.ProjectQueryInput{
					ID: r.Project.ID,
				},
			}

			if err := rq.ValidateWith(r.Context, r.Client, nil); err != nil {
				return err
			}
		}
	case len(r.Names) > 0:
		// While name existed.
		r.IDs = make([]object.ID, len(r.Names))
		for i := range r.Names {
			rq := model.ResourceQueryInput{
				Name: r.Names[i],
				Environment: &model.EnvironmentQueryInput{
					ID: r.EnvironmentQueryInput.ID,
				},
				Project: &model.ProjectQueryInput{
					ID: r.Project.ID,
				},
			}

			if err := rq.ValidateWith(r.Context, r.Client, nil); err != nil {
				return err
			}
			r.IDs[i] = rq.ID
		}
	}

	return nil
}

type WalrusFile struct {
	model.ResourceCreateInput
}
