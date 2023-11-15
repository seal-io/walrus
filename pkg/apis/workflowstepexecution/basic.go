package workflowstepexecution

import (
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstepexecution"
)

func (h Handler) Get(req GetRequest) (GetResponse, error) {
	entity, err := h.modelClient.WorkflowStepExecutions().Query().
		Where(workflowstepexecution.ID(req.ID)).
		Only(req.Context)
	if err != nil {
		return nil, err
	}

	return model.ExposeWorkflowStepExecution(entity), nil
}
