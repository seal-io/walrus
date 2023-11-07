package workflowstepexecution

import (
	"github.com/seal-io/walrus/pkg/dao/model"
)

type (
	GetRequest struct {
		model.WorkflowStepExecutionQueryInput `path:",inline"`
	}

	GetResponse = *model.WorkflowStepExecutionOutput
)

type UpdateRequest struct {
	model.WorkflowStepExecutionUpdateInput `path:",inline" json:",inline"`

	Status string `json:"status"`
}

func (r *UpdateRequest) Validate() error {
	if err := r.WorkflowStepExecutionUpdateInput.Validate(); err != nil {
		return err
	}

	return nil
}
