package step

import (
	"context"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"

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
) (main *v1alpha1.Template, subTemplates []*v1alpha1.Template, err error) {
	main = &v1alpha1.Template{
		Name:    "suspend-" + stepExecution.ID.String(),
		Suspend: &v1alpha1.SuspendTemplate{},
		Metadata: v1alpha1.Metadata{
			Labels: map[string]string{
				"step-execution-id": stepExecution.ID.String(),
			},
		},
	}

	return
}
