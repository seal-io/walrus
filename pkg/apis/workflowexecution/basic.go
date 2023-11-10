package workflowexecution

import (
	"time"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/workflowexecution"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstageexecution"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstepexecution"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	"github.com/seal-io/walrus/utils/topic"
)

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.WorkflowExecutions().Query().
		Where(workflowexecution.ID(req.ID)).
		WithStages(func(wsgq *model.WorkflowStageExecutionQuery) {
			wsgq.WithSteps(func(wseq *model.WorkflowStepExecutionQuery) {
				wseq.
					Select(workflowstepexecution.WithoutFields(workflowstepexecution.FieldRecord)...).
					Order(model.Asc(workflowstepexecution.FieldOrder))
			}).
				Order(model.Asc(workflowstageexecution.FieldOrder))
		}).
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	return model.ExposeWorkflowExecution(entity), nil
}

func (h Handler) Update(req UpdateRequest) error {
	entity, err := h.modelClient.WorkflowExecutions().Query().
		Where(workflowexecution.ID(req.ID)).
		Only(req.Context)
	if err != nil {
		return err
	}

	update := h.modelClient.WorkflowExecutions().UpdateOne(entity)

	switch req.Status {
	case types.ExecutionStatusSucceeded:
		status.WorkflowExecutionStatusRunning.True(entity, "")
	case types.ExecutionStatusFailed, types.ExecutionStatusError:
		status.WorkflowExecutionStatusRunning.False(entity, "")
	case types.ExecutionStatusRunning:
		status.WorkflowExecutionStatusPending.True(entity, "")
		status.WorkflowExecutionStatusRunning.Unknown(entity, "")

		update.SetExecuteTime(time.Now())
	default:
		return nil
	}

	entity.Status.SetSummary(status.WalkWorkflowExecution(&entity.Status))
	update.SetStatus(entity.Status)

	// If workflow execution is not running, set duration.
	if req.Status != types.ExecutionStatusRunning {
		update.SetDuration(int(time.Since(entity.ExecuteTime).Seconds()))
	}

	entity, err = update.Save(req.Context)
	if err != nil {
		return err
	}

	// Publish workflow topic.
	// Execution update will trigger workflow update of the workflow list.
	return topic.Publish(req.Context, modelchange.Workflow, modelchange.Event{
		Type: modelchange.EventTypeUpdate,
		IDs:  []object.ID{entity.WorkflowID},
	})
}

var (
	queryFields = []string{
		workflowexecution.FieldID,
		workflowexecution.FieldName,
		workflowexecution.FieldWorkflowID,
	}
	getFields  = workflowexecution.WithoutFields()
	sortFields = []string{
		workflowexecution.FieldID,
		workflowexecution.FieldName,
	}
)

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.WorkflowExecutions().Query().
		Where(workflowexecution.WorkflowID(req.Workflow.ID))
	if req.ID.Valid() {
		query.Where(workflowexecution.ID(req.ID))
	}

	if queries, ok := req.Querying(queryFields); ok {
		query = query.Where(queries)
	}

	if stream := req.Stream; stream != nil {
		// Handle stream request.
		if fields, ok := req.Extracting(getFields, getFields...); ok {
			query.Select(fields...)
		}

		if orders, ok := req.Sorting(sortFields, model.Desc(workflowexecution.FieldCreateTime)); ok {
			query.Order(orders...)
		}

		t, err := topic.Subscribe(modelchange.WorkflowExecution)
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

			var items []*model.WorkflowExecutionOutput

			switch dm.Type {
			case modelchange.EventTypeCreate, modelchange.EventTypeUpdate:
				entities, err := query.Clone().
					Where(workflowexecution.IDIn(dm.IDs...)).
					WithStages(func(wsgq *model.WorkflowStageExecutionQuery) {
						wsgq.WithSteps(func(wseq *model.WorkflowStepExecutionQuery) {
							wseq.Select(workflowstepexecution.WithoutFields(workflowstepexecution.FieldRecord)...).
								Order(model.Asc(workflowstepexecution.FieldOrder))
						}).Order(model.Asc(workflowstageexecution.FieldOrder))
					}).
					Unique(false).
					All(stream)
				if err != nil {
					return nil, 0, err
				}

				items = model.ExposeWorkflowExecutions(entities)
			case modelchange.EventTypeDelete:
				items = make([]*model.WorkflowExecutionOutput, len(dm.IDs))
				for i := range dm.IDs {
					items[i] = &model.WorkflowExecutionOutput{
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
	count, err := query.Clone().Count(req.Context)
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

	if orders, ok := req.Sorting(sortFields, model.Desc(workflowexecution.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		Unique(false).
		WithStages(func(wsgq *model.WorkflowStageExecutionQuery) {
			wsgq.WithSteps(func(wseq *model.WorkflowStepExecutionQuery) {
				wseq.Select(workflowstepexecution.WithoutFields(workflowstepexecution.FieldRecord)...).
					Order(model.Asc(workflowstepexecution.FieldOrder))
			}).Order(model.Asc(workflowstageexecution.FieldOrder))
		}).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeWorkflowExecutions(entities), count, nil
}

func (h Handler) Delete(req DeleteRequest) error {
	return h.modelClient.WorkflowExecutions().DeleteOneID(req.ID).
		Exec(req.Context)
}

func (h Handler) CollectionDelete(req CollectionDeleteRequest) error {
	ids := req.IDs()

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		_, err := tx.WorkflowStepExecutions().Delete().
			Where(workflowstepexecution.IDIn(ids...)).
			Exec(req.Context)
		if err != nil {
			return err
		}

		return nil
	})
}
