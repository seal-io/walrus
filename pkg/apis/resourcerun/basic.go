package resourcerun

import (
	"context"
	"time"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerun"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/topic"
)

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.ResourceRuns().Query().
		Where(resourcerun.ID(req.ID)).
		WithResource().
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	return model.ExposeResourceRun(entity), nil
}

var (
	getFields = resourcerun.WithoutFields(
		resourcerun.FieldRecord,
		resourcerun.FieldInputConfigs,
		resourcerun.FieldTemplateName,
		resourcerun.FieldTemplateVersion,
		resourcerun.FieldAttributes,
		resourcerun.FieldVariables,
	)
	sortFields = []string{
		resourcerun.FieldCreateTime,
	}
)

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.ResourceRuns().Query()

	if req.Resource != nil && req.Resource.ID != "" {
		query.Where(resourcerun.ResourceID(req.Resource.ID))
	}

	if fieldQueries, ok := req.Queries(); ok {
		query.Where(fieldQueries...)
	}

	if req.Rollback {
		query.Where(resourcerun.TypeIn(
			types.RunTypeCreate.String(),
			types.RunTypeUpgrade.String(),
		))
	}

	if stream := req.Stream; stream != nil {
		// Handle stream request.
		if fields, ok := req.Extracting(getFields, getFields...); ok {
			query.Select(fields...)
		}

		if orders, ok := req.Sorting(sortFields, model.Desc(resourcerun.FieldCreateTime)); ok {
			query.Order(orders...)
		}

		t, err := topic.Subscribe(modelchange.ResourceRun)
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

			var items []*model.ResourceRunOutput

			ids := dm.IDs()

			switch dm.Type {
			case modelchange.EventTypeCreate, modelchange.EventTypeUpdate:
				runs, err := query.Clone().
					Where(resourcerun.IDIn(ids...)).
					WithResource().
					// Must append service ID.
					Select(resourcerun.FieldResourceID).
					Unique(false).
					All(stream)
				if err != nil {
					return nil, 0, err
				}

				items = model.ExposeResourceRuns(runs)
			case modelchange.EventTypeDelete:
				items = make([]*model.ResourceRunOutput, len(ids))
				for i := range ids {
					items[i] = &model.ResourceRunOutput{
						ID: ids[i],
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

	// Handler normal request.

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

	if orders, ok := req.Sorting(sortFields, model.Desc(resourcerun.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		// Must append service ID.
		Select(resourcerun.FieldResourceID).
		WithResource().
		Unique(false).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeResourceRuns(entities), cnt, nil
}

func (h Handler) CollectionDelete(req CollectionDeleteRequest) error {
	ids := req.IDs()

	err := h.cleanPlanFiles(req.Context, ids...)
	if err != nil {
		return err
	}

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		_, err := tx.ResourceRuns().Delete().
			Where(resourcerun.IDIn(ids...)).
			Exec(req.Context)

		return err
	})
}

func (h Handler) cleanPlanFiles(ctx context.Context, ids ...object.ID) error {
	logger := log.WithName("apis").WithName("resourcerun")

	runs, err := h.modelClient.ResourceRuns().Query().
		Where(resourcerun.IDIn(ids...)).
		All(ctx)
	if err != nil {
		return err
	}

	gopool.Go(func() {
		subCtx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		for i := range runs {
			if rerr := h.storageManager.DeleteRunPlan(subCtx, runs[i]); err != nil {
				logger.Error(rerr, "failed to delete run plan", "run", runs[i].ID)
			}
		}
	})

	return nil
}
