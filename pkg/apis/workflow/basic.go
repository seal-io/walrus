package workflow

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/workflow"
	"github.com/seal-io/walrus/pkg/dao/model/workflowexecution"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstage"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstageexecution"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstep"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/datalisten/modelchange"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/topic"
)

func (h Handler) Create(req CreateRequest) (CreateResponse, error) {
	entity := req.Model()

	var err error

	err = h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		entity, err = tx.Workflows().Create().
			Set(entity).
			SaveE(req.Context, dao.WorkflowStagesEdgeSave)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return model.ExposeWorkflow(entity), nil
}

func (h Handler) Update(req UpdateRequest) error {
	entity := req.Model()

	var err error

	old, err := h.modelClient.Workflows().Query().
		Select(workflow.FieldVersion).
		Where(workflow.ID(entity.ID)).
		Only(req.Context)
	if err != nil {
		return err
	}

	entity.Version = old.Version

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		entity, err = tx.Workflows().UpdateOne(entity).
			Set(entity).
			SaveE(req.Context, dao.WorkflowStagesEdgeSave)
		if err != nil {
			return err
		}

		return nil
	})
}

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.Workflows().Query().
		Where(workflow.ID(req.ID)).
		WithStages(func(sgq *model.WorkflowStageQuery) {
			sgq.WithSteps(func(wsq *model.WorkflowStepQuery) {
				wsq.Order(model.Asc(workflowstep.FieldOrder))
			}).
				Order(model.Asc(workflowstage.FieldOrder))
		}).
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	return model.ExposeWorkflow(entity), nil
}

func (h Handler) Delete(req DeleteRequest) (err error) {
	wf, err := h.modelClient.Workflows().Query().
		Select(
			workflow.FieldID,
			workflow.FieldName,
		).
		Where(workflow.ID(req.ID)).
		WithExecutions(func(weq *model.WorkflowExecutionQuery) {
			weq.Select(
				workflowexecution.FieldID,
				workflowexecution.FieldStatus,
				workflowexecution.FieldWorkflowID,
			).Where(func(s *sql.Selector) {
				s.Where(sqljson.ValueEQ(
					workflowexecution.FieldStatus,
					status.WorkflowExecutionStatusRunning,
					sqljson.Path("summaryStatus"),
				))
			})
		}).
		Only(req.Context)
	if err != nil {
		return err
	}

	if len(wf.Edges.Executions) > 0 {
		return errorx.Errorf("workflow %s has running executions", wf.Name)
	}

	return h.modelClient.Workflows().DeleteOneID(req.ID).
		Exec(req.Context)
}

var (
	queryFields = []string{
		workflow.FieldID,
		workflow.FieldName,
	}
	getFields  = workflow.WithoutFields()
	sortFields = []string{
		workflow.FieldID,
		workflow.FieldName,
	}

	workflowExecutionLatestQuery = func(weq *model.WorkflowExecutionQuery) {
		weq.Where(func(s *sql.Selector) {
			sq := s.Clone().
				AppendSelectExprAs(
					sql.RowNumber().
						PartitionBy(workflowexecution.FieldWorkflowID).
						OrderBy(sql.Desc(workflowexecution.FieldCreateTime)),
					"row_number",
				).
				Where(s.P()).
				From(s.Table()).
				As(workflowexecution.Table)

			s.Where(sql.EQ(s.C("row_number"), 1)).
				From(sq)
		}).WithStages(func(sgq *model.WorkflowStageExecutionQuery) {
			sgq.Select(
				workflowstageexecution.FieldID,
				workflowstageexecution.FieldName,
				workflowstageexecution.FieldStatus,
				workflowexecution.FieldDuration,
				workflowexecution.FieldCreateTime,
			).Order(model.Asc(workflowstageexecution.FieldOrder))
		})
	}
)

func (h Handler) CollectionGet(req CollectionGetRequest) (CollectionGetResponse, int, error) {
	query := h.modelClient.Workflows().Query().
		Where(workflow.ProjectID(req.Project.ID))

	if queries, ok := req.Querying(queryFields); ok {
		query = query.Where(queries)
	}

	if stream := req.Stream; stream != nil {
		// Handle stream request.
		if fields, ok := req.Extracting(getFields, getFields...); ok {
			query.Select(fields...)
		}

		if orders, ok := req.Sorting(sortFields, model.Desc(workflow.FieldCreateTime)); ok {
			query.Order(orders...)
		}

		t, err := topic.Subscribe(modelchange.Workflow)
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

			var items []*model.WorkflowOutput

			switch dm.Type {
			case modelchange.EventTypeCreate, modelchange.EventTypeUpdate:
				entities, err := query.Clone().
					Where(workflow.IDIn(dm.IDs...)).
					Unique(false).
					WithExecutions(workflowExecutionLatestQuery).
					All(stream)
				if err != nil {
					return nil, 0, err
				}

				items = model.ExposeWorkflows(entities)
			case modelchange.EventTypeDelete:
				items = make([]*model.WorkflowOutput, len(dm.IDs))
				for i := range dm.IDs {
					items[i] = &model.WorkflowOutput{
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

	if orders, ok := req.Sorting(sortFields, model.Desc(workflow.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		Unique(false).
		WithExecutions(workflowExecutionLatestQuery).
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeWorkflows(entities), count, nil
}

func (h Handler) CollectionDelete(req CollectionDeleteRequest) error {
	ids := req.IDs()

	return h.modelClient.WithTx(req.Context, func(tx *model.Tx) error {
		workflows, err := tx.Workflows().Query().
			Select(
				workflow.FieldID,
				workflow.FieldName,
			).
			Where(workflow.IDIn(ids...)).
			WithExecutions(func(weq *model.WorkflowExecutionQuery) {
				weq.Select(
					workflowexecution.FieldID,
					workflowexecution.FieldStatus,
					workflowexecution.FieldWorkflowID,
				).Where(func(s *sql.Selector) {
					s.Where(sqljson.ValueEQ(
						workflowexecution.FieldStatus,
						status.WorkflowExecutionStatusRunning,
						sqljson.Path("summaryStatus"),
					))
				})
			}).
			All(req.Context)
		if err != nil {
			return err
		}

		for i := range workflows {
			wf := workflows[i]
			if len(wf.Edges.Executions) > 0 {
				return errorx.Errorf("workflow %s has running executions", wf.Name)
			}
		}

		_, err = tx.Workflows().Delete().
			Where(workflow.IDIn(ids...)).
			Exec(req.Context)
		if err != nil {
			return err
		}

		return nil
	})
}
