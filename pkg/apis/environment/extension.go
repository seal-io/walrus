package environment

import (
	"entgo.io/ent/dialect/sql"

	"github.com/seal-io/walrus/pkg/auths"
	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/cli/config"
	"github.com/seal-io/walrus/pkg/cli/manifest"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinitionmatchingrule"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerelationship"
	"github.com/seal-io/walrus/pkg/dao/model/template"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/model/variable"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/deployer"
	deployertf "github.com/seal-io/walrus/pkg/deployer/terraform"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
	pkgcomponent "github.com/seal-io/walrus/pkg/resourcecomponents"
	"github.com/seal-io/walrus/pkg/resourcedefinitions"
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
					s.Where(
						sql.ColumnsNEQ(resourcerelationship.FieldResourceID, resourcerelationship.FieldDependencyID),
					)
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

		vertex := GraphVertex{
			GraphVertexID: GraphVertexID{
				Kind: types.VertexKindResource,
				ID:   entity.ID,
			},
			Name:        entity.Name,
			Description: entity.Description,
			Labels:      entity.Labels,
			CreateTime:  entity.CreateTime,
			UpdateTime:  entity.UpdateTime,
			Status:      entity.Status.Summary,
			Extensions: map[string]any{
				"projectID":     entity.ProjectID,
				"environmentID": entity.EnvironmentID,
				"labels":        entity.Labels,
				"isService":     pkgresource.IsService(entity),
			},
		}

		// TODO resource definition icon.
		if pkgresource.IsService(entity) {
			vertex.Icon = entity.Edges.Template.Edges.Template.Icon
		}

		// Append Resource to vertices.
		vertices = append(vertices, vertex)

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

	dp, err := h.getDeployer(req.Context)
	if err != nil {
		return nil, err
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

	return createEnvironment(req.Context, h.modelClient, dp, entity, req.Draft)
}

func (h Handler) RouteGetResourceDefinitions(
	req RouteGetResourceDefinitionsRequest,
) (RouteGetResourceDefinitionsResponse, error) {
	rds, err := h.modelClient.ResourceDefinitions().Query().
		WithMatchingRules(func(rq *model.ResourceDefinitionMatchingRuleQuery) {
			rq.Order(model.Asc(resourcedefinitionmatchingrule.FieldOrder)).
				Unique(false).
				Select(
					resourcedefinitionmatchingrule.FieldResourceDefinitionID,
					resourcedefinitionmatchingrule.FieldTemplateID,
					resourcedefinitionmatchingrule.FieldSelector,
				).
				WithTemplate(func(tq *model.TemplateVersionQuery) {
					tq.Select(
						templateversion.FieldID,
						templateversion.FieldVersion,
						templateversion.FieldName,
					)
				})
		}).
		All(req.Context)
	if err != nil {
		return nil, err
	}

	env, err := h.modelClient.Environments().Query().
		Where(environment.ID(req.ID)).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(project.FieldName)
		}).
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	var availableRds []*model.ResourceDefinition

	for _, rd := range rds {
		m := resourcedefinitions.MatchEnvironment(
			rd.Edges.MatchingRules,
			env.Edges.Project.Name,
			env.Name,
			env.Type,
			env.Labels,
		)
		if m != nil {
			availableRds = append(availableRds, rd)
		}
	}

	return dao.ExposeResourceDefinitions(availableRds), nil
}

func (h Handler) RouteStart(req RouteStartRequest) error {
	dp, err := h.getDeployer(req.Context)
	if err != nil {
		return err
	}

	resources, err := req.Client.Resources().Query().
		Where(resource.EnvironmentID(req.ID)).
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
		All(req.Context)
	if err != nil {
		return err
	}

	toStartResources := make(model.Resources, 0, len(resources))

	for _, r := range resources {
		if pkgresource.IsInactive(r) {
			toStartResources = append(toStartResources, r)
		}
	}

	err = h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		if err := pkgresource.SetSubjectID(req.Context, toStartResources...); err != nil {
			return err
		}

		return pkgresource.SetResourceStatusScheduled(req.Context, tx, dp, toStartResources...)
	})
	if err != nil {
		return errorx.Wrap(err, "failed to start environment")
	}

	return nil
}

func (h Handler) RouteStop(req RouteStopRequest) error {
	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		deployerOpts := deptypes.CreateOptions{
			Type:       deployertf.DeployerType,
			KubeConfig: h.kubeConfig,
		}

		dp, err := deployer.Get(req.Context, deployerOpts)
		if err != nil {
			return err
		}

		destroyOpts := pkgresource.Options{
			Deployer: dp,
		}

		for _, s := range req.stoppableResources {
			if !pkgresource.CanBeStopped(s) {
				continue
			}

			err = pkgresource.Stop(req.Context, tx, s, destroyOpts)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (h Handler) RouteApply(req RouteApplyRequest) error {
	sj := session.MustGetSubject(req.Context)

	token, err := auths.GetAccessToken(
		req.Context, h.modelClient,
		sj.ID, types.TokenKindAPI, types.WalrusOperationTokenName)
	if err != nil {
		if !model.IsNotFound(err) {
			return errorx.New(err.Error())
		}

		token, err = auths.CreateAccessToken(req.Context, h.modelClient,
			sj.ID, types.TokenKindAPI, types.WalrusOperationTokenName, nil)
		if err != nil {
			return errorx.New(err.Error())
		}
	}

	sc := serverContext(req.Project.Name, req.Name, token.AccessToken)

	loader := manifest.DefaultLoader(sc, true)

	set, err := loader.LoadFromByte([]byte(req.YAML))
	if err != nil {
		return errorx.Wrap(err, "failed to load walrus file")
	}

	operator := manifest.DefaultApplyOperator(sc, false)

	r, err := operator.Operate(set)
	if err != nil {
		return errorx.Wrap(err, "failed to apply walrus file")
	}

	if r.Failed.Len() != 0 {
		var keys []string
		for _, v := range r.Failed.All() {
			keys = append(keys, v.Key())
		}

		return errorx.Errorf("failed to apply: %v", keys)
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
