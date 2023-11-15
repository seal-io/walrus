package workflow

import (
	"encoding/json"
	"fmt"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/sets"

	apiresource "github.com/seal-io/walrus/pkg/apis/resource"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/subject"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/utils/validation"
)

func init() {
	stepValidatorCreators = map[string]func(*model.WorkflowStepCreateInput) WorkflowStepValidator{
		types.WorkflowStepTypeService: func(step *model.WorkflowStepCreateInput) WorkflowStepValidator {
			return &WorkflowStepServiceValidator{step}
		},
		types.WorkflowStepTypeApproval: func(step *model.WorkflowStepCreateInput) WorkflowStepValidator {
			return &WorkflowStepApprovalValidator{step}
		},
	}
}

func validateStages(ctx *gin.Context, client *model.Client, stages []*model.WorkflowStageCreateInput) error {
	for _, stage := range stages {
		if err := validateStage(ctx, client, stage); err != nil {
			return err
		}
	}

	return nil
}

func validateStage(ctx *gin.Context, client *model.Client, stage *model.WorkflowStageCreateInput) error {
	if err := validateSteps(ctx, client, stage.Steps); err != nil {
		return fmt.Errorf("invalid steps: %w", err)
	}

	return nil
}

func validateSteps(ctx *gin.Context, client *model.Client, steps []*model.WorkflowStepCreateInput) error {
	for _, step := range steps {
		if err := validateStep(ctx, client, step); err != nil {
			return err
		}
	}

	return nil
}

func validateStep(ctx *gin.Context, client *model.Client, step *model.WorkflowStepCreateInput) error {
	creator, ok := stepValidatorCreators[step.Type]
	if !ok {
		return fmt.Errorf("unknown step type: %s", step.Type)
	}

	err := validation.MapStringNoMustache(step.Attributes)
	if err != nil {
		return fmt.Errorf("invalid attributes: %w", err)
	}

	stepValidator := creator(step)

	if step.RetryStrategy != nil {
		switch step.RetryStrategy.RetryPolicy {
		case wfv1.RetryPolicyAlways,
			wfv1.RetryPolicyOnFailure,
			wfv1.RetryPolicyOnError,
			wfv1.RetryPolicyOnTransientError:
		default:
			return fmt.Errorf("invalid retry policy: %s", step.RetryStrategy.RetryPolicy)
		}
	}

	return stepValidator.Validate(ctx, client)
}

type WorkflowStepValidator interface {
	Set(*model.WorkflowStepCreateInput)
	// Validate validates the attributes of the workflow step.
	Validate(*gin.Context, *model.Client) error
}

// WorkflowStepServiceValidator validates the attributes of a service workflow step.
type WorkflowStepServiceValidator struct {
	*model.WorkflowStepCreateInput `path:",inline" json:",inline"`
}

type ServiceCreateMeta struct {
	Project     *model.ProjectQueryInput     `json:"project"`
	Environment *model.EnvironmentQueryInput `json:"environment"`
}

func (s *WorkflowStepServiceValidator) Set(input *model.WorkflowStepCreateInput) {
	s.WorkflowStepCreateInput = input
}

func (s *WorkflowStepServiceValidator) Validate(ctx *gin.Context, client *model.Client) error {
	rci := &model.ResourceCreateInput{}
	scm := &ServiceCreateMeta{}

	rci.SetGinContext(ctx)
	rci.SetModelClient(client)

	jsonData, err := json.Marshal(s.Attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal service attributes: %w", err)
	}

	if err := json.Unmarshal(jsonData, rci); err != nil {
		return fmt.Errorf("failed to unmarshal service input: %w", err)
	}

	if err := json.Unmarshal(jsonData, scm); err != nil {
		return fmt.Errorf("failed to unmarshal service meta: %w", err)
	}

	rci.Project = scm.Project
	rci.Environment = scm.Environment

	if err := apiresource.ValidateCreateInput(rci); err != nil {
		return err
	}

	return nil
}

// WorkflowStepApprovalValidator validates the attributes of an approval workflow step.
type WorkflowStepApprovalValidator struct {
	*model.WorkflowStepCreateInput
}

func (s *WorkflowStepApprovalValidator) Set(input *model.WorkflowStepCreateInput) {
	s.WorkflowStepCreateInput = input
}

func (s *WorkflowStepApprovalValidator) Validate(ctx *gin.Context, client *model.Client) error {
	approvalType, ok := s.Attributes[types.WorkflowStepApprovalType].(string)
	if !ok {
		return fmt.Errorf("invalid approval type")
	}

	switch approvalType {
	case types.WorkflowStepApprovalTypeOr, types.WorkflowStepApprovalTypeAnd:
	default:
		return fmt.Errorf("invalid approval type: %s", approvalType)
	}

	approvalUsers, ok := s.Attributes[types.WorkflowStepApprovalUsers].([]any)
	if !ok {
		return fmt.Errorf("invalid approval users")
	}

	for _, user := range approvalUsers {
		uid := object.ID(user.(string))
		if !uid.Valid() {
			return fmt.Errorf("invalid user id: %s", user)
		}

		// Check user roles.
		s, err := client.Subjects().Query().
			Where(subject.ID(uid)).
			WithRoles().
			Only(ctx)
		if err != nil {
			return err
		}

		// Only project owner and member can approve steps.
		approvalRoles := sets.NewString(
			types.SystemRoleUser,
			types.ProjectRoleMember,
			types.ProjectRoleOwner,
		)

		hasApprovalRole := false

		for _, role := range s.Edges.Roles {
			rid := role.RoleID
			if approvalRoles.Has(rid) {
				hasApprovalRole = true
				break
			}
		}

		if !hasApprovalRole {
			return fmt.Errorf("user %s has no approval role", uid)
		}
	}

	s.Attributes = map[string]any{
		types.WorkflowStepApprovalType:  approvalType,
		types.WorkflowStepApprovalUsers: approvalUsers,
	}

	return nil
}

var stepValidatorCreators = map[string]func(*model.WorkflowStepCreateInput) WorkflowStepValidator{}
