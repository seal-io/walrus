package workflowstepexecution

import (
	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/workflowstepexecution"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/utils/errorx"
)

type RouteLogRequest struct {
	_ struct{} `route:"GET=/log"`

	model.WorkflowStepExecutionQueryInput `path:",inline"`

	Previous     bool   `query:"previous,omitempty"`
	SinceSeconds *int64 `query:"sinceSeconds,omitempty"`
	TailLines    *int64 `query:"tailLines,omitempty"`
	Timestamps   bool   `query:"timestamps,omitempty"`

	Stream *runtime.RequestUnidiStream
}

func (r *RouteLogRequest) Validate() error {
	if err := r.WorkflowStepExecutionQueryInput.Validate(); err != nil {
		return err
	}

	return nil
}

func (r *RouteLogRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type RouteApproveRequest struct {
	_ struct{} `route:"POST=/approve"`

	model.WorkflowStepExecutionQueryInput `path:",inline"`

	Approve bool `json:"approve"`
}

func (r *RouteApproveRequest) Validate() error {
	if err := r.WorkflowStepExecutionQueryInput.Validate(); err != nil {
		return err
	}

	step, err := r.Client.WorkflowStepExecutions().Query().
		Where(workflowstepexecution.ID(r.ID)).
		Only(r.Context)
	if err != nil {
		return err
	}

	wasa, err := types.NewWorkflowStepApprovalSpec(step.Attributes)
	if err != nil {
		return err
	}

	if wasa.IsApproved() {
		return errorx.New("step is already approved")
	}

	if wasa.IsRejected() {
		return errorx.New("step is rejected")
	}

	return nil
}
