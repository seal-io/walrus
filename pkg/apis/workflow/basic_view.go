package workflow

import (
	"fmt"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/workflow"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/utils/validation"
)

type (
	CreateRequest struct {
		model.WorkflowCreateInput `path:",inline" json:",inline"`
	}

	CreateResponse = *model.WorkflowOutput
)

func (r *CreateRequest) Validate() error {
	if err := r.WorkflowCreateInput.Validate(); err != nil {
		return err
	}

	if err := r.WorkflowCreateInput.Variables.Validate(); err != nil {
		return fmt.Errorf("invalid variables configs: %w", err)
	}

	if err := validation.IsValidName(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if err := validateType(r.Type); err != nil {
		return fmt.Errorf("invalid type: %w", err)
	}

	if err := validateStages(r.Context, r.Client, r.Stages); err != nil {
		return fmt.Errorf("invalid stages: %w", err)
	}

	return nil
}

type (
	GetRequest struct {
		model.WorkflowQueryInput `path:",inline"`
	}

	GetResponse = *model.WorkflowOutput
)

type UpdateRequest struct {
	model.WorkflowUpdateInput `path:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if err := r.WorkflowUpdateInput.Validate(); err != nil {
		return err
	}

	if err := validation.IsValidName(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if err := validateStages(r.Context, r.Client, r.Stages); err != nil {
		return fmt.Errorf("invalid stages: %w", err)
	}

	return nil
}

type (
	DeleteRequest struct {
		model.WorkflowQueryInput `path:",inline"`
	}

	DeleteResponse = *model.WorkflowDeleteInput
)

func (r *DeleteRequest) Validate() error {
	if err := r.WorkflowQueryInput.Validate(); err != nil {
		return err
	}

	return nil
}

type (
	CollectionGetRequest struct {
		model.WorkflowQueryInputs `path:",inline"`

		runtime.RequestCollection[
			predicate.Workflow, workflow.OrderOption,
		] `query:",inline"`

		Stream *runtime.RequestUnidiStream
	}

	CollectionGetResponse = []*model.WorkflowOutput
)

func (r *CollectionGetRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type CollectionDeleteRequest = model.WorkflowDeleteInputs

func validateType(workflowType string) error {
	switch workflowType {
	case types.WorkflowTypeDefault:
		return nil
	// Add more workflow types here.
	default:
		return fmt.Errorf("invalid workflow type: %s", workflowType)
	}
}
