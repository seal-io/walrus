package step

import (
	"context"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/workflow/step/types"
)

type ApprovalStepManager struct {
	mc model.ClientSet
}

func NewApprovalStepManager(opts types.CreateOptions) types.StepManager {
	return &ApprovalStepManager{
		mc: opts.ModelClient,
	}
}

func (m *ApprovalStepManager) GenerateTemplates(
	ctx context.Context,
	stepExecution *model.WorkflowStepExecution,
) (main *wfv1.Template, subTemplates []*wfv1.Template, err error) {
	main = &wfv1.Template{
		Name:    StepTemplateName(stepExecution),
		Suspend: &wfv1.SuspendTemplate{},
		Metadata: wfv1.Metadata{
			Labels: map[string]string{
				"step-execution-id": stepExecution.ID.String(),
			},
		},
	}

	return
}
