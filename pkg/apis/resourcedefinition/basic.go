package resourcedefinition

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinition"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinitionmatchingrule"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	"github.com/seal-io/walrus/pkg/templates/openapi"
	"github.com/seal-io/walrus/utils/topic"
)

func (h Handler) Create(req CreateRequest) (CreateResponse, error) {
	entity := req.Model()

	if err := generateSchema(req.Context, req.Client, entity); err != nil {
		return nil, fmt.Errorf("failed to generate schema: %w", err)
	}

	err := h.modelClient.WithTx(req.Context, func(tx *model.Tx) (err error) {
		entity, err = tx.ResourceDefinitions().Create().
			Set(entity).
			SaveE(req.Context, dao.ResourceDefinitionMatchingRulesEdgeSave)

		return err
	})
	if err != nil {
		return nil, err
	}

	return model.ExposeResourceDefinition(entity), nil
}

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.ResourceDefinitions().Query().
		Where(resourcedefinition.ID(req.ID)).
		WithMatchingRules(func(rq *model.ResourceDefinitionMatchingRuleQuery) {
			rq.Order(model.Desc(resourcedefinitionmatchingrule.FieldCreateTime)).
				Unique(false).
				Select(
					resourcedefinitionmatchingrule.FieldName,
					resourcedefinitionmatchingrule.FieldTemplateID,
					resourcedefinitionmatchingrule.FieldSelector,
					resourcedefinitionmatchingrule.FieldAttributes,
				).
				WithTemplate(func(tvq *model.TemplateVersionQuery) {
					tvq.Select(
						templateversion.FieldID,
						templateversion.FieldVersion,
						templateversion.FieldName,
					).WithTemplate(func(tq *model.TemplateQuery) {
						tq.Select(templateversion.FieldID)
					})
				})
		}).
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	return dao.ExposeResourceDefinition(entity), nil
}

func (h Handler) Update(req UpdateRequest) error {
	entity := req.Model()

	if err := generateSchema(req.Context, req.Client, entity); err != nil {
		return fmt.Errorf("failed to generate schema: %w", err)
	}

	err := h.modelClient.WithTx(req.Context, func(tx *model.Tx) (err error) {
		entity, err = tx.ResourceDefinitions().UpdateOne(entity).
			Set(entity).
			SaveE(req.Context, dao.ResourceDefinitionMatchingRulesEdgeSave)

		return err
	})

	return err
}

func (h Handler) Delete(req DeleteRequest) error {
	return h.modelClient.ResourceDefinitions().DeleteOneID(req.ID).
		Exec(req.Context)
}

var (
	queryFields = []string{
		resourcedefinition.FieldName,
	}
	getFields  = resourcedefinition.WithoutFields()
	sortFields = []string{
		resourcedefinition.FieldName,
		resourcedefinition.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.ResourceDefinitions().Query()

	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	if stream := req.Stream; stream != nil {
		// Handle stream request.
		if fields, ok := req.Extracting(getFields, getFields...); ok {
			query.Select(fields...)
		}

		if orders, ok := req.Sorting(sortFields, model.Desc(resourcedefinition.FieldCreateTime)); ok {
			query.Order(orders...)
		}

		t, err := topic.Subscribe(modelchange.ResourceDefinition)
		if err != nil {
			return nil, 0, err
		}

		defer func() { t.Unsubscribe() }()

		for {
			var event topic.Event

			event, err = t.Receive(stream)
			if err != nil {
				return nil, 0, err
			}

			dm, ok := event.Data.(modelchange.Event)
			if !ok {
				continue
			}

			var items []*model.ResourceDefinitionOutput

			switch dm.Type {
			case modelchange.EventTypeCreate, modelchange.EventTypeUpdate:
				entities, err := query.Clone().
					Where(resourcedefinition.IDIn(dm.IDs...)).
					Unique(false).
					All(stream)
				if err != nil {
					return nil, 0, err
				}

				items = dao.ExposeResourceDefinitions(entities)
			case modelchange.EventTypeDelete:
				items = make([]*model.ResourceDefinitionOutput, len(dm.IDs))
				for i := range dm.IDs {
					items[i] = &model.ResourceDefinitionOutput{
						ID: dm.IDs[i],
					}
				}
			}

			if len(items) == 0 {
				continue
			}

			resp := runtime.TypedResponse(dm.Type.String(), items)
			if err = stream.SendJSON(resp); err != nil {
				return nil, 0, err
			}
		}
	}

	// Handle normal request.

	// Get count.
	cnt, err := query.Clone().Count(req.Context)
	if err != nil {
		return nil, 0, err
	}

	// Get entities.
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}

	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}

	if orders, ok := req.Sorting(sortFields, model.Desc(resourcedefinition.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return dao.ExposeResourceDefinitions(entities), cnt, nil
}

func (h Handler) CollectionDelete(req CollectionDeleteRequest) error {
	ids := req.IDs()

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		_, err := tx.ResourceDefinitions().Delete().
			Where(resourcedefinition.IDIn(ids...)).
			Exec(req.Context)

		return err
	})
}

// generateSchema generates definition schema with inputs/outputs intersection of matching template versions.
func generateSchema(ctx context.Context, mc model.ClientSet, definition *model.ResourceDefinition) error {
	tvIDs := make([]object.ID, len(definition.Edges.MatchingRules))
	for i, rule := range definition.Edges.MatchingRules {
		tvIDs[i] = rule.TemplateID
	}

	templateVersions, err := mc.TemplateVersions().Query().
		Where(templateversion.IDIn(tvIDs...)).
		All(ctx)
	if err != nil {
		return err
	}

	if len(templateVersions) == 0 {
		return nil
	}

	definition.Schema = types.Schema{
		OpenAPISchema: &openapi3.T{
			OpenAPI: openapi.OpenAPIVersion,
			Info: &openapi3.Info{
				Title:   fmt.Sprintf("OpenAPI schema for resource definition %s", definition.Name),
				Version: "v0.0.0",
			},
			Components: templateVersions[0].Schema.OpenAPISchema.Components,
		},
	}
	for i := 1; i < len(templateVersions); i++ {
		definition.Schema.Intersect(&types.Schema{
			OpenAPISchema: templateVersions[i].Schema.OpenAPISchema,
		})
	}

	return nil
}
