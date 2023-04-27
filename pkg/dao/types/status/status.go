package status

import (
	"reflect"
	"sort"
	"time"
)

const (
	GeneralStatusReady              ConditionType = "Ready"
	GeneralStatusReadyTransitioning string        = "Provisioning" // transitioning status of Ready
	GeneralStatusError              string        = "Error"        // error status of Ready
)

const (
	ModuleStatusInitializing = "Initializing"
	ModuleStatusReady        = "Ready"
	ModuleStatusError        = "Error"
)

const (
	ApplicationRevisionStatusRunning   = "Running"
	ApplicationRevisionStatusSucceeded = "Succeeded"
	ApplicationRevisionStatusFailed    = "Failed"
)

// Status wrap the summary of conditions and condition details.
type Status struct {
	Summary    `json:",inline"`
	Conditions []Condition `json:"conditions,omitempty"`

	// used for check whether status changed.
	conditionChanged bool
	summaryChanged   bool
}

func (s *Status) Changed() bool {
	return s.conditionChanged || s.summaryChanged
}

func (s *Status) SetConditions(conds []Condition) {
	// sort conditions
	sortConditions := func(conds []Condition) {
		sort.Slice(conds, func(i, j int) bool {
			return conds[i].Type < conds[j].Type
		})
	}
	sortConditions(s.Conditions)
	sortConditions(conds)

	// unchanged
	if reflect.DeepEqual(s.Conditions, conds) {
		s.conditionChanged = false
	}

	// change
	s.Conditions = conds
	s.conditionChanged = true
}

func (s *Status) SetSummary(summary *Summary) {
	// unchanged
	if reflect.DeepEqual(s.Summary, *summary) {
		s.summaryChanged = false
	}

	// changed
	s.summaryChanged = true
	s.Summary = *summary
}

func (s Status) Equal(newStatue Status) bool {
	if !reflect.DeepEqual(s.Summary, newStatue.Summary) {
		return false
	}

	if len(s.Conditions) != len(newStatue.Conditions) {
		return false
	}

	sortConditions := func(conds []Condition) {
		sort.Slice(conds, func(i, j int) bool {
			return conds[i].Type < conds[j].Type
		})
	}

	sortConditions(s.Conditions)
	sortConditions(newStatue.Conditions)
	return reflect.DeepEqual(s.Conditions, newStatue.Conditions)
}

// Condition is the condition details.
type Condition struct {
	// type of condition in CamelCase.
	Type ConditionType `json:"type,omitempty"`
	// status of the condition, one of True, False, Unknown.
	Status ConditionStatus `json:"status,omitempty"`
	// This should be when the underlying condition changed.
	LastUpdateTime time.Time `json:"lastUpdateTime,omitempty"`
	// message is a human-readable message indicating details about the status.
	Message string `json:"message,omitempty"`
	// reason contains a programmatic identifier indicating the reason for the condition's last transition.
	Reason string `json:"reason"`
}

// Summary is the summary of conditions.
type Summary struct {
	SummaryStatus        string `json:"summaryStatus,omitempty"`
	SummaryStatusMessage string `json:"summaryStatusMessage,omitempty"`
	Error                bool   `json:"error,omitempty"`
	Transitioning        bool   `json:"transitioning,omitempty"`
}
