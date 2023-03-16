package status

import (
	"time"
)

const (
	ConnectorStatusProvisioned   ConditionType = "Provisioned"
	ConnectorStatusToolsDeployed ConditionType = "Deployed"
	ConnectorStatusCostSynced    ConditionType = "CostSynced"
	ConnectorStatusReady         ConditionType = "Ready"

	ConnectorStatusToolsDeployedTransitioning string = "Deploying"    // transitioning status of ConnectorStatusToolsDeployed
	ConnectorStatusCostSyncedTransitioning    string = "CostSyncing"  // transitioning status of ConnectorFinOpsSyncStatusSynced
	ConnectorStatusProvisionedTransitioning   string = "Provisioning" // transitioning status of ConnectorStatusProvisioned

)

const (
	ModuleStatusInitializing = "Initializing"
	ModuleStatusReady        = "Ready"
	ModuleStatusError        = "Error"
)

const (
	ApplicationInstanceStatusDeploying    = "Deploying"
	ApplicationInstanceStatusDeployed     = "Deployed"
	ApplicationInstanceStatusDeployFailed = "DeployFailed"
	ApplicationInstanceStatusDeleting     = "Deleting"
	ApplicationInstanceStatusDeleteFailed = "DeleteFailed"
)

const (
	ApplicationRevisionStatusRunning   = "Running"
	ApplicationRevisionStatusSucceeded = "Succeeded"
	ApplicationRevisionStatusFailed    = "Failed"
)

// Status wrap the summary of conditions and condition details
type Status struct {
	Summary    `json:",inline"`
	Conditions []Condition `json:"conditions,omitempty"`

	// used for
	conditionChanged bool
}

func (s Status) ConditionChanged() bool {
	return s.conditionChanged
}

// Condition is the condition details
type Condition struct {
	// type of condition in CamelCase.
	Type ConditionType `json:"type,omitempty"`
	// status of the condition, one of True, False, Unknown.
	Status ConditionStatus `json:"status,omitempty"`
	// This should be when the underlying condition changed.
	LastUpdateTime time.Time `json:"lastUpdateTime,omitempty"`
	// message is a human-readable message indicating details about the status.
	Message string `json:"message,omitempty"`
}

// Summary is the summary of conditions
type Summary struct {
	Status        string `json:"status,omitempty"`
	StatusMessage string `json:"statusMessage,omitempty"`
	Error         bool   `json:"error,omitempty"`
	Transitioning bool   `json:"transitioning,omitempty"`
}
