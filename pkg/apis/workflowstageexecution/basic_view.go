package workflowstageexecution

import "github.com/seal-io/walrus/pkg/dao/model"

type UpdateRequest struct {
	model.WorkflowStageExecutionUpdateInput `path:",inline" json:",inline"`

	Status string `json:"status"`
}

func (r *UpdateRequest) Validate() error {
	if err := r.WorkflowStageExecutionUpdateInput.Validate(); err != nil {
		return err
	}

	return nil
}
