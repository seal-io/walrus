package status

import (
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/dao/types/status"
)

func NewSummarizer(defaultType status.ConditionType) *Summarizer {
	return &Summarizer{
		defaultType: defaultType,
	}
}

type Summarizer struct {
	defaultType status.ConditionType
	// True ==
	// False == error
	// Unknown == transitioning
	// e.g. Initialized
	errorFalseTransitioningUnknown map[status.ConditionType]string

	// True == transitioning
	// False == error
	// Unknown ==
	errorFalseTransitioningTrue map[status.ConditionType]string

	// True == error
	// False ==
	// Unknown == transitioning
	// e.g. DiskPressure
	errorTrueTransitioningUnknown map[status.ConditionType]string

	// True == error
	// False == transitioning
	// Unknown ==
	// e.g MemoryPressure
	errorTrueTransitioningFalse map[status.ConditionType]string

	// True ==
	// False == transitioning
	// Unknown == error
	// e.g. Completed
	errorUnknownTransitioningFalse map[status.ConditionType]string

	// True == transitioning
	// False ==
	// Unknown == error
	errorUnknownTransitioningTrue map[status.ConditionType]string

	// condition types in sequence.
	conditionTypes []status.ConditionType

	// reasons mapping keep the error identifier for condition type,
	// for example:
	// conditionType 'Processing', reason 'ReplicaSetUpdatedReason' and 'NewReplicaSetAvailable' both have True status,
	// could use reason to mark ReplicaSetUpdatedReason is in transitioning.

	// errorReasons keep programmatic identifier for error.
	errorReasons map[status.ConditionType]sets.Set[string]

	// transitioningReasons keep programmatic identifier for transitioning.
	transitioningReasons map[status.ConditionType]sets.Set[string]

	// readyReasons keep programmatic identifier for ready.
	readyReasons map[status.ConditionType]sets.Set[string]
}

// AddErrorFalseTransitioningUnknown add condition which treat condition status
// False == error
// Unknown == transitioning
// True == success
func (s *Summarizer) AddErrorFalseTransitioningUnknown(conditionType status.ConditionType, inTransitioningDisplay string) {
	s.conditionTypes = append(s.conditionTypes, conditionType)

	if len(s.errorFalseTransitioningUnknown) == 0 {
		s.errorFalseTransitioningUnknown = make(map[status.ConditionType]string)
	}
	s.errorFalseTransitioningUnknown[conditionType] = inTransitioningDisplay
}

// AddErrorFalseTransitioningTrue add condition which treat condition status.
// False == error
// Unknown == success
// True == transitioning
func (s *Summarizer) AddErrorFalseTransitioningTrue(conditionType status.ConditionType, inTransitioningDisplay string) {
	s.conditionTypes = append(s.conditionTypes, conditionType)

	if len(s.errorFalseTransitioningTrue) == 0 {
		s.errorFalseTransitioningTrue = make(map[status.ConditionType]string)
	}
	s.errorFalseTransitioningTrue[conditionType] = inTransitioningDisplay
}

// AddErrorTrueTransitioningUnknown add condition which treat condition status.
// False == success
// Unknown == transitioning
// True == error
func (s *Summarizer) AddErrorTrueTransitioningUnknown(conditionType status.ConditionType, inTransitioningDisplay string) {
	s.conditionTypes = append(s.conditionTypes, conditionType)

	if len(s.errorTrueTransitioningUnknown) == 0 {
		s.errorTrueTransitioningUnknown = make(map[status.ConditionType]string)
	}
	s.errorTrueTransitioningUnknown[conditionType] = inTransitioningDisplay
}

// AddErrorTrueTransitioningFalse add condition which treat condition status.
// False == transitioning
// Unknown == success
// True == error
func (s *Summarizer) AddErrorTrueTransitioningFalse(conditionType status.ConditionType, inTransitioningDisplay string) {
	s.conditionTypes = append(s.conditionTypes, conditionType)

	if len(s.errorTrueTransitioningFalse) == 0 {
		s.errorTrueTransitioningFalse = make(map[status.ConditionType]string)
	}
	s.errorTrueTransitioningFalse[conditionType] = inTransitioningDisplay
}

// AddErrorUnknownTransitioningFalse add condition which treat condition status.
// False == transitioning
// Unknown == error
// True == success
func (s *Summarizer) AddErrorUnknownTransitioningFalse(conditionType status.ConditionType, inTransitioningDisplay string) {
	s.conditionTypes = append(s.conditionTypes, conditionType)

	if len(s.errorUnknownTransitioningFalse) == 0 {
		s.errorUnknownTransitioningFalse = make(map[status.ConditionType]string)
	}
	s.errorUnknownTransitioningFalse[conditionType] = inTransitioningDisplay
}

// AddErrorUnknownTransitioningTrue add condition which treat condition status.
// False == success
// Unknown == error
// True == transitioning
func (s *Summarizer) AddErrorUnknownTransitioningTrue(conditionType status.ConditionType, inTransitioningDisplay string) {
	s.conditionTypes = append(s.conditionTypes, conditionType)

	if len(s.errorUnknownTransitioningTrue) == 0 {
		s.errorUnknownTransitioningTrue = make(map[status.ConditionType]string)
	}
	s.errorUnknownTransitioningTrue[conditionType] = inTransitioningDisplay
}

func (s *Summarizer) AddErrorReason(conditionType status.ConditionType, reason ...string) {
	if s.errorReasons == nil {
		s.errorReasons = make(map[status.ConditionType]sets.Set[string])
	}

	if _, ok := s.errorReasons[conditionType]; !ok {
		s.errorReasons[conditionType] = sets.Set[string]{}
	}

	s.errorReasons[conditionType].Insert(reason...)
}

// AddTransitionReason add reason for specific condition type.
func (s *Summarizer) AddTransitionReason(conditionType status.ConditionType, reason ...string) {
	if s.transitioningReasons == nil {
		s.transitioningReasons = make(map[status.ConditionType]sets.Set[string])
	}

	if _, ok := s.transitioningReasons[conditionType]; !ok {
		s.transitioningReasons[conditionType] = sets.Set[string]{}
	}
	s.transitioningReasons[conditionType].Insert(reason...)
}

func (s *Summarizer) AddReadyReason(conditionType status.ConditionType, reason ...string) {
	if s.readyReasons == nil {
		s.readyReasons = make(map[status.ConditionType]sets.Set[string])
	}

	if _, ok := s.readyReasons[conditionType]; !ok {
		s.readyReasons[conditionType] = sets.Set[string]{}
	}
	s.readyReasons[conditionType].Insert(reason...)
}

// Summarize summarize conditions from status.
func (s *Summarizer) Summarize(st *status.Status) *status.Summary {
	// init
	summary := &status.Summary{
		Error:                false,
		Transitioning:        false,
		SummaryStatus:        "",
		SummaryStatusMessage: "",
	}

	summarizers := []func(st *status.Status, summary *status.Summary){
		s.summarizeErrors,
		s.summarizeTransitioning,
		s.summarizeReason,
	}

	for _, summarizer := range summarizers {
		summarizer(st, summary)
	}

	if summary.SummaryStatus == "" {
		summary.SummaryStatus = string(s.defaultType)
	}
	return summary
}

func (s *Summarizer) summarizeErrors(st *status.Status, summary *status.Summary) {
	for _, v := range s.conditionTypes {
		c := getCondition(st.Conditions, v)
		if c == nil {
			continue
		}

		switch c.Status {
		case status.ConditionStatusFalse:
			_, ok1 := s.errorFalseTransitioningUnknown[c.Type]
			_, ok2 := s.errorFalseTransitioningTrue[c.Type]
			if ok1 || ok2 {
				summary.SummaryStatus = string(c.Type)
				summary.Error = true
				summary.SummaryStatusMessage = c.Message
				return
			}
		case status.ConditionStatusTrue:
			_, ok1 := s.errorTrueTransitioningUnknown[c.Type]
			_, ok2 := s.errorTrueTransitioningFalse[c.Type]
			if ok1 || ok2 {
				summary.SummaryStatus = string(c.Type)
				summary.Error = true
				summary.SummaryStatusMessage = c.Message
				return
			}
		case status.ConditionStatusUnknown:
			_, ok1 := s.errorUnknownTransitioningTrue[c.Type]
			_, ok2 := s.errorUnknownTransitioningFalse[c.Type]
			if ok1 || ok2 {
				summary.SummaryStatus = string(c.Type)
				summary.Error = true
				summary.SummaryStatusMessage = c.Message
				return
			}
		}
	}
}

func (s *Summarizer) summarizeTransitioning(st *status.Status, summary *status.Summary) {
	// already set
	if summary.Error {
		return
	}

	for _, v := range s.conditionTypes {
		c := getCondition(st.Conditions, v)
		if c == nil {
			continue
		}
		switch c.Status {
		case status.ConditionStatusFalse:
			if tras, ok := s.errorUnknownTransitioningFalse[c.Type]; ok {
				summary.SummaryStatus = tras
				summary.Transitioning = true
				summary.SummaryStatusMessage = c.Message
				return
			}

			if tras, ok := s.errorTrueTransitioningFalse[c.Type]; ok {
				summary.SummaryStatus = tras
				summary.Transitioning = true
				summary.SummaryStatusMessage = c.Message
				return
			}
		case status.ConditionStatusTrue:
			if tras, ok := s.errorUnknownTransitioningTrue[c.Type]; ok {
				summary.SummaryStatus = tras
				summary.Transitioning = true
				summary.SummaryStatusMessage = c.Message
				return
			}

			if tras, ok := s.errorFalseTransitioningTrue[c.Type]; ok {
				summary.SummaryStatus = tras
				summary.Transitioning = true
				summary.SummaryStatusMessage = c.Message
				return
			}
		case status.ConditionStatusUnknown:
			if tras, ok := s.errorTrueTransitioningUnknown[c.Type]; ok {
				summary.SummaryStatus = tras
				summary.Transitioning = true
				summary.SummaryStatusMessage = c.Message
				return
			}

			if tras, ok := s.errorFalseTransitioningUnknown[c.Type]; ok {
				summary.SummaryStatus = tras
				summary.Transitioning = true
				summary.SummaryStatusMessage = c.Message
				return
			}
		}
	}
}

func (s *Summarizer) summarizeReason(st *status.Status, summary *status.Summary) {
	// already set
	if summary.Error {
		return
	}

	// error reason
	for condType, v := range s.errorReasons {
		c := getCondition(st.Conditions, condType)
		if c == nil {
			continue
		}
		if v.Has(c.Reason) {
			summary.Error = true
			return
		}
	}

	// transition reason
	for condType, v := range s.transitioningReasons {
		c := getCondition(st.Conditions, condType)
		if c == nil {
			continue
		}
		if v.Has(c.Reason) {
			summary.Transitioning = true
			return
		}
	}

	// ready reason
	for condType, v := range s.readyReasons {
		c := getCondition(st.Conditions, condType)
		if c == nil {
			continue
		}
		if v.Has(c.Reason) {
			// reset to the default
			summary.SummaryStatus = string(s.defaultType)
			summary.SummaryStatusMessage = ""
			summary.Error = false
			summary.Transitioning = false
			return
		}
	}
}

func getCondition(conditions []status.Condition, conditionType status.ConditionType) *status.Condition {
	for i, v := range conditions {
		if v.Type == conditionType {
			return &conditions[i]
		}
	}
	return nil
}
