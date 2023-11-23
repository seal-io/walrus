package workflow

import (
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/workflow"
)

type RouteGetLatestExecutionRequest struct {
	_ struct{} `route:"GET=/latest-execution"`

	model.WorkflowQueryInput `path:",inline"`
}

type RouteGetLatestExecutionResponse = *model.WorkflowExecutionOutput

type RouteRunRequest struct {
	_ struct{} `route:"POST=/run"`

	model.WorkflowQueryInput `path:",inline"`

	Variables   map[string]string `json:"variables,omitempty"`
	Description string            `json:"description,omitempty"`
}

type RouteRunResponse = *model.WorkflowExecutionOutput

func (r *RouteRunRequest) Validate() error {
	err := r.WorkflowQueryInput.Validate()
	if err != nil {
		return err
	}

	wf, err := r.Client.Workflows().Query().
		Where(workflow.ID(r.ID)).
		Only(r.Context)
	if err != nil {
		return err
	}

	_, err = dao.OverwriteWorkflowVariables(r.Variables, wf.Variables)

	return err
}
