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
	// WorkflowStepRejectedUsers is the key of reject users in spec.
	WorkflowStepRejectedUsers = "rejectedUsers"
)

type WorkflowStepApprovalSpec struct {
	ApprovalUsers []object.ID `json:"approvalUsers"`
	ApprovedUsers []object.ID `json:"approvedUsers"`
	RejectedUsers []object.ID `json:"rejectedUsers"`
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

	userIndexes := []string{
		WorkflowStepApprovalUsers,
		WorkflowStepApprovedUsers,
		WorkflowStepRejectedUsers,
	}

	for i := range userIndexes {
		if _, ok := spec[userIndexes[i]]; !ok {
			continue
		}

		if v, ok := spec[userIndexes[i]].([]any); ok {
			users, err := toObjectIDs(v)
			if err != nil {
				return nil, err
			}

			switch userIndexes[i] {
			case WorkflowStepApprovalUsers:
				s.ApprovalUsers = users
			case WorkflowStepApprovedUsers:
				s.ApprovedUsers = users
			case WorkflowStepRejectedUsers:
				s.RejectedUsers = users
			}
		}
	}

	return s, nil
}

func (s *WorkflowStepApprovalSpec) IsRejected() bool {
	return len(s.RejectedUsers) > 0
}

func (s *WorkflowStepApprovalSpec) IsApproved() bool {
	if s.IsRejected() {
		return false
	}

	if s.Type == WorkflowStepApprovalTypeOr {
		return len(s.ApprovedUsers) > 0
	}

	if s.Type == WorkflowStepApprovalTypeAnd {
		return len(s.ApprovedUsers) == len(s.ApprovalUsers) &&
			len(s.ApprovedUsers) > 0
	}

	return false
}

func (s *WorkflowStepApprovalSpec) IsApprovalUser(user object.ID) bool {
	isApprovalUser := false

	for i := range s.ApprovalUsers {
		if s.ApprovalUsers[i] != user {
			continue
		}

		isApprovalUser = true
	}

	return isApprovalUser
}

// SetUserApproval sets the user approval status.
func (s *WorkflowStepApprovalSpec) SetUserApproval(user object.ID, approved bool) error {
	if approved {
		return s.SetApprovedUser(user)
	}

	return s.SetRejectedUser(user)
}

func (s *WorkflowStepApprovalSpec) SetApprovedUser(user object.ID) error {
	if !s.IsApprovalUser(user) {
		return errors.New("user is not an approval users")
	}

	if isExist(user, s.ApprovedUsers) {
		return nil
	}

	s.ApprovedUsers = append(s.ApprovedUsers, user)

	return nil
}

func (s *WorkflowStepApprovalSpec) SetRejectedUser(user object.ID) error {
	if !s.IsApprovalUser(user) {
		return errors.New("user is not an approval users")
	}

	if isExist(user, s.ApprovedUsers) {
		return nil
	}

	s.RejectedUsers = append(s.RejectedUsers, user)

	return nil
}

func (s *WorkflowStepApprovalSpec) ToAttributes() map[string]any {
	return map[string]any{
		WorkflowStepApprovalType:  s.Type,
		WorkflowStepApprovalUsers: s.ApprovalUsers,
		WorkflowStepApprovedUsers: s.ApprovedUsers,
		WorkflowStepRejectedUsers: s.RejectedUsers,
	}
}

func toObjectIDs(users []any) ([]object.ID, error) {
	ids := make([]object.ID, 0, len(users))
	um := make(map[object.ID]struct{}, len(users))

	for i := range users {
		if v, ok := users[i].(string); ok {
			id := object.ID(v)
			if !id.Valid() {
				return nil, fmt.Errorf("invalid user id: %s", v)
			}

			if _, ok := um[id]; ok {
				continue
			}

			ids = append(ids, id)
			um[id] = struct{}{}
		}
	}

	return ids, nil
}

// isExist checks if user is in users.
func isExist(user object.ID, users []object.ID) bool {
	for i := range users {
		if users[i] == user {
			return true
		}
	}

	return false
}
