package types

import (
	"errors"
	"fmt"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"

	"github.com/seal-io/walrus/pkg/dao/types/object"
)

const (
	WorkflowTypeDefault = "default"
	WorkflowTypeCron    = "cron"
)

// RetryStrategy is the retry strategy of a workflow step.
// See https://raw.githubusercontent.com/argoproj/argo-workflows/master/examples/retry-conditional.yaml
type RetryStrategy struct {
	Limit       int                  `json:"limit"`
	RetryPolicy v1alpha1.RetryPolicy `json:"retryPolicy"`
	Backoff     *v1alpha1.Backoff    `json:"backoff"`
}

const (
	WorkflowStepTypeService  = "service"
	WorkflowStepTypeApproval = "approval"
)

const (
	WorkflowExecutionTriggerTypeManual = "manual"
)

// WorkflowExecutionTrigger is the trigger of a workflow execution.
type WorkflowExecutionTrigger struct {
	Type string `json:"type"`
	User string `json:"user"`
}

const (
	ExecutionStatusRunning   = "Running"
	ExecutionStatusSucceeded = "Succeeded"
	ExecutionStatusFailed    = "Failed"
	ExecutionStatusError     = "Error"
)

const (
	// WorkflowStepApprovalTypeOr means step is approved
	// if any of the approval users approves it.
	WorkflowStepApprovalTypeOr = "or"
	// WorkflowStepApprovalTypeAnd means step is approved
	// only all of the approval users approve it.
	WorkflowStepApprovalTypeAnd = "and"

	// WorkflowStepApprovalType is the key of type in spec.
	WorkflowStepApprovalType = "approvalType"
	// WorkflowStepApprovalUsers is the key of approval users in spec.
	WorkflowStepApprovalUsers = "approvalUsers"
	// WorkflowStepApprovedUsers is the key of approved users in spec.
	WorkflowStepApprovedUsers = "approvedUsers"
)

type WorkflowStepApprovalSpec struct {
	ApprovalUsers []object.ID `json:"approvalUsers"`
	ApprovedUsers []object.ID `json:"approvedUsers"`
	Type          string      `json:"approvalType"`
}

func NewWorkflowStepApprovalSpec(spec map[string]any) (*WorkflowStepApprovalSpec, error) {
	if spec == nil {
		return nil, errors.New("invalid input: nil spec")
	}

	s := &WorkflowStepApprovalSpec{}

	if v, ok := spec[WorkflowStepApprovalType]; ok {
		s.Type = v.(string)
	}

	switch s.Type {
	case WorkflowStepApprovalTypeOr, WorkflowStepApprovalTypeAnd:
	default:
		return nil, errors.New("invalid input: invalid approval type")
	}

	if v, ok := spec[WorkflowStepApprovalUsers].([]any); ok {
		users, err := toObjectIDs(v)
		if err != nil {
			return nil, err
		}
		s.ApprovalUsers = removeDuplicatedUsers(users)
	}

	if v, ok := spec[WorkflowStepApprovedUsers].([]any); ok {
		users, err := toObjectIDs(v)
		if err != nil {
			return nil, err
		}
		s.ApprovedUsers = removeDuplicatedUsers(users)
	}

	return s, nil
}

func (s *WorkflowStepApprovalSpec) IsApproved() bool {
	if s.Type == WorkflowStepApprovalTypeOr {
		return len(s.ApprovedUsers) > 0
	}

	if s.Type == WorkflowStepApprovalTypeAnd {
		return len(s.ApprovedUsers) == len(s.ApprovalUsers) &&
			len(s.ApprovedUsers) > 0
	}

	return false
}

func (s *WorkflowStepApprovalSpec) SetApprovedUser(user object.ID) error {
	isApprovalUser := false

	for i := range s.ApprovalUsers {
		if s.ApprovalUsers[i] != user {
			continue
		}

		isApprovalUser = true
	}

	if !isApprovalUser {
		return errors.New("user is not an approval users")
	}

	s.ApprovedUsers = append(s.ApprovedUsers, user)

	// Remove duplicated users.
	s.ApprovedUsers = removeDuplicatedUsers(s.ApprovedUsers)

	return nil
}

func (s *WorkflowStepApprovalSpec) ToAttributes() map[string]any {
	return map[string]any{
		WorkflowStepApprovalType:  s.Type,
		WorkflowStepApprovalUsers: s.ApprovalUsers,
		WorkflowStepApprovedUsers: s.ApprovedUsers,
	}
}

func removeDuplicatedUsers(users []object.ID) []object.ID {
	m := make(map[object.ID]struct{}, len(users))

	for i := range users {
		m[users[i]] = struct{}{}
	}

	newUsers := make([]object.ID, 0, len(m))

	for k := range m {
		newUsers = append(newUsers, k)
	}

	return newUsers
}

func toObjectIDs(users []any) ([]object.ID, error) {
	ids := make([]object.ID, 0, len(users))

	for i := range users {
		if v, ok := users[i].(string); ok {
			id := object.ID(v)
			if !id.Valid() {
				return nil, fmt.Errorf("invalid user id: %s", v)
			}

			ids = append(ids, id)
		}
	}

	return ids, nil
}
