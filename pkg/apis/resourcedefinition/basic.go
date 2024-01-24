package resourcedefinition

import (
	"fmt"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinition"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinitionmatchingrule"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	"github.com/seal-io/walrus/pkg/resourcedefinitions"
	"github.com/seal-io/walrus/pkg/templates"
	"github.com/seal-io/walrus/utils/topic"
)

func (h Handler) Create(req CreateRequest) (CreateResponse, error) {
	entity := req.Model()

	err := templates.SetResourceDefinitionSchemaDefault(req.Context, req.Client, entity)
	if err != nil {
		return nil, err
	}

	err = resourcedefinitions.GenSchema(req.Context, req.Client, entity)
	if err != nil {
		return nil, fmt.Errorf("failed to generate schema: %w", err)
	}

	err = h.modelClient.WithTx(req.Context, func(tx *model.Tx) (err error) {
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
			rq.Order(model.Asc(resourcedefinitionmatchingrule.FieldOrder)).
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

	err := templates.SetResourceDefinitionSchemaDefault(req.Context, req.Client, entity)
	if err != nil {
		return err
	}

	err = resourcedefinitions.GenSchema(req.Context, req.Client, entity)
	if err != nil {
		return fmt.Errorf("failed to generate schema: %w", err)
	}

	err = h.modelClient.WithTx(req.Context, func(tx *model.Tx) (err error) {
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

			ids := dm.IDs()

			switch dm.Type {
			case modelchange.EventTypeCreate, modelchange.EventTypeUpdate:
				entities, err := query.Clone().
					Where(resourcedefinition.IDIn(ids...)).
					Unique(false).
					All(stream)
				if err != nil {
					return nil, 0, err
				}

				items = dao.ExposeResourceDefinitions(entities)
			case modelchange.EventTypeDelete:
				items = make([]*model.ResourceDefinitionOutput, len(ids))
				for i := range ids {
					items[i] = &model.ResourceDefinitionOutput{
						ID:   ids[i],
						Name: dm.Data[i].Name,
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
