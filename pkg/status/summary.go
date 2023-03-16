package status

import (
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

	// condition types in sequence
	conditionTypes []status.ConditionType
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

// AddErrorFalseTransitioningTrue add condition which treat condition status
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

// AddErrorTrueTransitioningUnknown add condition which treat condition status
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

// AddErrorTrueTransitioningFalse add condition which treat condition status
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

// AddErrorUnknownTransitioningFalse add condition which treat condition status
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

// AddErrorUnknownTransitioningTrue add condition which treat condition status
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

// SetSummarize summarize conditions in the status and set summary to status
func (s *Summarizer) SetSummarize(st *status.Status) {
	// reset
	st.Summary.Error = false
	st.Summary.Transitioning = false
	st.Summary.Status = ""
	st.Summary.StatusMessage = ""

	summarizers := []func(st *status.Status){
		s.summarizeErrors,
		s.summarizeTransitioning,
	}

	for _, summarizer := range summarizers {
		summarizer(st)
	}

	if st.Status == "" {
		st.Status = string(s.defaultType)
	}
}

func (s *Summarizer) summarizeErrors(st *status.Status) {
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
				st.Status = string(c.Type)
				st.Error = true
				st.StatusMessage = c.Message
				return
			}
		case status.ConditionStatusTrue:
			_, ok1 := s.errorTrueTransitioningUnknown[c.Type]
			_, ok2 := s.errorTrueTransitioningFalse[c.Type]
			if ok1 || ok2 {
				st.Status = string(c.Type)
				st.Error = true
				st.StatusMessage = c.Message
				return
			}
		case status.ConditionStatusUnknown:
			_, ok1 := s.errorUnknownTransitioningTrue[c.Type]
			_, ok2 := s.errorUnknownTransitioningFalse[c.Type]
			if ok1 || ok2 {
				st.Status = string(c.Type)
				st.Error = true
				st.StatusMessage = c.Message
				return
			}
		}
	}
}

func (s *Summarizer) summarizeTransitioning(st *status.Status) {
	// already set
	if st.Error {
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
				st.Status = tras
				st.Transitioning = true
				st.StatusMessage = c.Message
				return
			}

			if tras, ok := s.errorTrueTransitioningFalse[c.Type]; ok {
				st.Status = tras
				st.Transitioning = true
				st.StatusMessage = c.Message
				return
			}
		case status.ConditionStatusTrue:
			if tras, ok := s.errorUnknownTransitioningTrue[c.Type]; ok {
				st.Status = tras
				st.Transitioning = true
				st.StatusMessage = c.Message
				return
			}

			if tras, ok := s.errorFalseTransitioningTrue[c.Type]; ok {
				st.Status = tras
				st.Transitioning = true
				st.StatusMessage = c.Message
				return
			}
		case status.ConditionStatusUnknown:
			if tras, ok := s.errorTrueTransitioningUnknown[c.Type]; ok {
				st.Status = tras
				st.Transitioning = true
				st.StatusMessage = c.Message
				return
			}

			if tras, ok := s.errorFalseTransitioningUnknown[c.Type]; ok {
				st.Status = tras
				st.Transitioning = true
				st.StatusMessage = c.Message
				return
			}
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
