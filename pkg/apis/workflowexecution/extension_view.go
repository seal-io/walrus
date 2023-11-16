package workflowexecution

import (
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/workflowexecution"
	"github.com/seal-io/walrus/pkg/dao/types/status"
)

type RouteRerunRequest struct {
	_ struct{} `route:"PUT=/rerun"`

	model.WorkflowExecutionQueryInput `path:",inline"`
}

func (r *RouteRerunRequest) Validate() error {
	return r.WorkflowExecutionQueryInput.Validate()
}

type RouteStopRequest struct {
	_ struct{} `route:"PUT=/stop"`

	model.WorkflowExecutionQueryInput `path:",inline"`
}

func (r *RouteStopRequest) Validate() error {
	if err := r.WorkflowExecutionQueryInput.Validate(); err != nil {
		return err
	}

	entity, err := r.Client.WorkflowExecutions().Query().
		Where(workflowexecution.ID(r.ID)).
		Only(r.Context)
	if err != nil {
		return err
	}

	if status.WorkflowExecutionStatusPending.IsUnknown(entity) ||
		status.WorkflowExecutionStatusRunning.IsTrue(entity) {
		return fmt.Errorf("workflow execution is pending or finished")
	}

	return nil
}
