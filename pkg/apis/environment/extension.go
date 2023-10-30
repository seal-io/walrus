package environment

import (
	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerelationship"
	"github.com/seal-io/walrus/pkg/dao/model/template"
	"github.com/seal-io/walrus/pkg/dao/model/variable"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	pkgcomponent "github.com/seal-io/walrus/pkg/resourcecomponents"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/log"
)

var getResourceFields = resource.WithoutFields(
	resource.FieldUpdateTime)

func (h Handler) RouteGetGraph(req RouteGetGraphRequest) (*RouteGetGraphResponse, error) {
	// Fetch resource entities.
	entities, err := h.modelClient.Resources().Query().
		Where(resource.EnvironmentID(req.ID)).
		Order(model.Desc(resource.FieldCreateTime)).
		Select(getResourceFields...).
		// Must extract dependency.
		WithDependencies(func(dq *model.ResourceRelationshipQuery) {
			dq.Select(resourcerelationship.FieldDependencyID).
				Where(func(s *sql.Selector) {
					s.Where(sql.ColumnsNEQ(resourcerelationship.FieldResourceID, resourcerelationship.FieldDependencyID))
				})
		}).
		// Must extract resource.
		WithComponents(func(rq *model.ResourceComponentQuery) {
			dao.ResourceComponentShapeClassQuery(rq)
		}).
		WithTemplate(func(tq *model.TemplateVersionQuery) {
			tq.Select(template.FieldID).
				WithTemplate(func(tq *model.TemplateQuery) {
					tq.Select(template.FieldIcon)
				})
		}).
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, err
	}

	verticesCap, edgesCap := getCaps(entities)

	// Construct response.
	var (
		vertices  = make([]GraphVertex, 0, verticesCap)
		edges     = make([]GraphEdge, 0, edgesCap)
		operators = make(map[object.ID]optypes.Operator)
	)

	for i := 0; i < len(entities); i++ {
		entity := entities[i]

		// Append Resource to vertices.
		vertices = append(vertices, GraphVertex{
			GraphVertexID: GraphVertexID{
				Kind: types.VertexKindResource,
				ID:   entity.ID,
			},
			Name:        entity.Name,
			Description: entity.Description,
			Icon:        entity.Edges.Template.Edges.Template.Icon,
			Labels:      entity.Labels,
			CreateTime:  entity.CreateTime,
			UpdateTime:  entity.UpdateTime,
			Status:      entity.Status.Summary,
			Extensions: map[string]any{
				"projectID":     entity.ProjectID,
				"environmentID": entity.EnvironmentID,
			},
		})

		// Append the link of related Resources to edges.
		for j := 0; j < len(entity.Edges.Dependencies); j++ {
			edges = append(edges, GraphEdge{
				Type: types.EdgeTypeDependency,
				Start: GraphVertexID{
					Kind: types.VertexKindResource,
					ID:   entity.ID,
				},
				End: GraphVertexID{
					Kind: types.VertexKindResource,
					ID:   entity.Edges.Dependencies[j].DependencyID,
				},
			})
		}

		// Set keys for next operations, e.g. Log, Exec and so on.
		if !req.WithoutKeys {
			pkgcomponent.SetKeys(
				req.Context,
				log.WithName("api").WithName("environment"),
				h.modelClient,
				entity.Edges.Components,
				operators)
		}

		// Append ResourceComponent to vertices,
		// and append the link of related ResourceComponents to edges.
		vertices, edges = pkgcomponent.GetVerticesAndEdges(
			entity.Edges.Components, vertices, edges)

		for j := 0; j < len(entity.Edges.Components); j++ {
			// Append the link from Resource to ResourceComponent into edges.
			edges = append(edges, GraphEdge{
				Type: types.EdgeTypeComposition,
				Start: GraphVertexID{
					Kind: types.VertexKindResource,
					ID:   entity.ID,
				},
				End: GraphVertexID{
					Kind: types.VertexKindResourceComponentGroup,
					ID:   entity.Edges.Components[j].ID,
				},
			})
		}
	}

	return &RouteGetGraphResponse{
		Vertices: vertices,
		Edges:    edges,
	}, nil
}

func getCaps(entities model.Resources) (int, int) {
	// Calculate capacity for allocation.
	var verticesCap, edgesCap int

	// Count the number of Resource.
	verticesCap = len(entities)
	for i := 0; i < len(entities); i++ {
		// Count the vertex size of ResourceComponent,
		// and the edge size from Resource to ResourceComponent.
		verticesCap += len(entities[i].Edges.Components)
		edgesCap += len(entities[i].Edges.Dependencies)

		for j := 0; j < len(entities[i].Edges.Components); j++ {
			// Count the vertex size of instances,
			// and the edge size from ResourceComponentGroup to instance ResourceComponent.
			verticesCap += len(entities[i].Edges.Components[j].Edges.Instances)
			edgesCap += len(entities[i].Edges.Components[j].Edges.Instances) +
				len(entities[i].Edges.Components[j].Edges.Dependencies)

			for k := 0; k < len(entities[i].Edges.Components[j].Edges.Instances); k++ {
				verticesCap += len(entities[i].Edges.Components[j].Edges.Components)
				edgesCap += len(entities[i].Edges.Components[j].Edges.Components)
			}
		}
	}

	return verticesCap, edgesCap
}

func (h Handler) RouteCloneEnvironment(req RouteCloneEnvironmentRequest) (*RouteCloneEnvironmentResponse, error) {
	entity := req.Model()

	var variableNames []string

	variables := entity.Edges.Variables
	for i := range variables {
		v := variables[i]
		if v == nil {
			return nil, errorx.New("invalid input: nil variable")
		}

		if v.Sensitive && v.Value == "" {
			variableNames = append(variableNames, v.Name)
		}
	}

	// Fetch and fill the value with sensitive variables from the cloned environment.
	if len(variableNames) > 0 {
		vs, err := h.modelClient.Variables().Query().
			Where(
				variable.EnvironmentID(req.CloneEnvironmentId),
				variable.NameIn(variableNames...),
			).
			All(req.Context)
		if err != nil {
			return nil, err
		}

		variableValueMap := make(map[string]crypto.String)
		for _, vv := range vs {
			variableValueMap[vv.Name] = vv.Value
		}

		for _, v := range variables {
			if v.Sensitive && v.Value == "" {
				v.Value = variableValueMap[v.Name]
			}
		}
	}

	return createEnvironment(req.Context, h.modelClient, entity)
}
