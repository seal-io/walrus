package status

import "github.com/seal-io/seal/utils/strs"

type (
	Converter interface {
		// Convert base on input status and message and gives a converted seal status.
		Convert(string, string) *Status
	}
)

// NewConverter creates a converter by the given status and status message,
// status not in the normalStatus and errorStatus will set transitioning to true.
//   - `normalStatus` specifies the normal status list, won't change the error and transition,
//     while summary status is in the normal status list.
//   - `errorStatus` specifies the error status list, set the error to true,
//     while the summary status is in the error status list.
func NewConverter[T ~string](normalStatus, errorStatus []T) Converter {
	var (
		ns = make(map[T]struct{})
		es = make(map[T]struct{})
	)

	for _, v := range normalStatus {
		ns[v] = struct{}{}
	}

	for _, v := range errorStatus {
		es[v] = struct{}{}
	}

	return &converter[T]{
		normalStatus: ns,
		errorStatus:  es,
	}
}

type converter[T ~string] struct {
	normalStatus map[T]struct{}
	errorStatus  map[T]struct{}
}

func (w *converter[T]) Convert(sm, msg string) *Status {
	st := &Status{}

	if sm == "" {
		return st
	}

	_, isErr := w.errorStatus[any(sm).(T)]
	_, isNormal := w.normalStatus[any(sm).(T)]

	switch {
	case isErr:
		st.Error = true
		st.Transitioning = false
	case isNormal:
		st.Error = false
		st.Transitioning = false
	default:
		st.Error = false
		st.Transitioning = true
	}

	// Format status.
	st.SummaryStatus = strs.Camelize(sm)
	st.SummaryStatusMessage = msg

	return st
}
