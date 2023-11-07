package workflow

import (
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
)

type RouteGetLatestExecutionRequest struct {
	_ struct{} `route:"GET=/latest-execution"`

	model.WorkflowQueryInput `path:",inline"`
}

type RouteGetLatestExecutionResponse = *model.WorkflowExecutionOutput

type RouteRunRequest struct {
	_ struct{} `route:"POST=/run"`

	model.WorkflowQueryInput `path:",inline" json:",inline"`

	Params      map[string]string `json:"params"`
	Description string            `json:"description,omitempty"`
}

type RouteRunResponse = *model.WorkflowExecutionOutput

func (r *RouteRunRequest) Validate() error {
	err := r.WorkflowQueryInput.Validate()
	if err != nil {
		return err
	}

	return dao.CheckParams(r.Params)
}
