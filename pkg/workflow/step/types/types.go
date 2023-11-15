package types

import (
	"context"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
)

// Type is a type of step.
type Type string

func (t Type) String() string {
	return string(t)
}

const (
	StepTypeService  Type = types.WorkflowStepTypeService
	StepTypeApproval Type = types.WorkflowStepTypeApproval
)

type StepManager interface {
	// GenerateTemplates generates templates for step.
	// Returns main template and sub templates.
	// Main template is a template that will be executed.
	// A main template can have sub templates.
	GenerateTemplates(
		context.Context,
		*model.WorkflowStepExecution,
	) (mainTemplate *wfv1.Template, subTemplates []*wfv1.Template, err error)
}
