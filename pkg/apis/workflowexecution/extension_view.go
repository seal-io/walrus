package workflowexecution

import "github.com/seal-io/walrus/pkg/dao/model"

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
	return r.WorkflowExecutionQueryInput.Validate()
}
