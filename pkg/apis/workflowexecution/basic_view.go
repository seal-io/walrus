package workflowexecution

import (
	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/workflowexecution"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/utils/errorx"
)

type (
	GetRequest struct {
		model.WorkflowExecutionQueryInput `path:",inline"`
	}

	GetResponse = *model.WorkflowExecutionOutput
)

type UpdateRequest struct {
	model.WorkflowExecutionUpdateInput `path:",inline" json:",inline"`

	Status string `json:"status"`
}

func (r *UpdateRequest) Validate() error {
	if err := r.WorkflowExecutionUpdateInput.Validate(); err != nil {
		return err
	}

	return nil
}

type (
	CollectionGetRequest struct {
		model.WorkflowExecutionQueryInputs `path:",inline"`

		runtime.RequestCollection[
			predicate.WorkflowExecution, workflowexecution.OrderOption,
		] `query:",inline"`

		Stream *runtime.RequestUnidiStream

		ID object.ID `query:"id"`
	}

	CollectionGetResponse []*model.WorkflowExecutionOutput
)

func (r *CollectionGetRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type (
	DeleteRequest struct {
		model.WorkflowExecutionQueryInput `path:",inline"`
	}

	DeleteResponse = *model.WorkflowExecutionDeleteInput
)

func (r *DeleteRequest) Validate() error {
	err := r.WorkflowExecutionQueryInput.Validate()
	if err != nil {
		return err
	}

	workflowExecution, err := r.Client.WorkflowExecutions().Query().
		Where(workflowexecution.ID(r.ID)).
		Only(r.Context)
	if err != nil {
		return err
	}

	if status.WorkflowExecutionStatusRunning.IsUnknown(workflowExecution) {
		return errorx.Errorf("workflow execution %s is running", workflowExecution.ID.String())
	}

	return nil
}

type CollectionDeleteRequest struct {
	model.WorkflowExecutionDeleteInputs
}

func (r *CollectionDeleteRequest) Validate() error {
	ids := r.IDs()

	executions, err := r.Client.WorkflowExecutions().Query().
		Where(workflowexecution.IDIn(ids...)).
		Select(workflowexecution.FieldID).
		All(r.Context)
	if err != nil {
		return err
	}

	for i := range executions {
		execution := executions[i]
		if status.WorkflowExecutionStatusRunning.IsUnknown(execution) {
			return errorx.Errorf("workflow execution %s is running", execution.ID.String())
		}
	}

	return nil
}
