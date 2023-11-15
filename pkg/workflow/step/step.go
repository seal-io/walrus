package step

import (
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/workflow/step/types"
)

const StepNamePrefix = "workflow-step-execution-"

var managerCreators map[types.Type]func(types.CreateOptions) types.StepManager

func init() {
	managerCreators = map[types.Type]func(types.CreateOptions) types.StepManager{
		types.StepTypeService:  NewServiceStepManager,
		types.StepTypeApproval: NewApprovalStepManager,
	}
}

func GetStepManager(opts types.CreateOptions) (types.StepManager, error) {
	constructor, ok := managerCreators[opts.Type]
	if !ok {
		return nil, fmt.Errorf("unknown step type: %s", opts.Type)
	}

	return constructor(opts), nil
}

func StepTemplateName(stepExecution *model.WorkflowStepExecution) string {
	return fmt.Sprintf("%s%s", StepNamePrefix, stepExecution.ID.String())
}
