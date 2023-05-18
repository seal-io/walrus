package status

import "github.com/seal-io/seal/utils/slice"

type (
	// SummaryDecide get current summary and return the updated summary, and whether it is change.
	SummaryDecide func(current Summary) (change bool, update Summary)
)

// NewSummaryWalker creates a stacking walker by the given status groups,
// and applies the customized decision to the status group.
//   - `normalStatus` specifies the normal status list, won't change the error and transition,
//     while summary status is in the normal status list.
//   - `errorStatus` specifies the error status list, set the error to true,
//     while the summary status is in the error status list.
//   - `transitioningStatus` specifies the transitioning status list, set the transitioning to true,
//     while the summary status is in the transitioningStatus list.
//   - `arrange` applies the customized decision,
//     for example, customize the summary display status or changing summary message display content by its status.
func NewSummaryWalker[T ~string](normalStatus, errorStatus, transitioningStatus []T, decides ...SummaryDecide) Walker {
	return &summaryWalker[T]{
		normalStatus:        normalStatus,
		errorStatus:         errorStatus,
		transitioningStatus: transitioningStatus,
		decides:             decides,
	}
}

type summaryWalker[T ~string] struct {
	normalStatus        []T
	errorStatus         []T
	transitioningStatus []T
	decides             []SummaryDecide
}

func (w *summaryWalker[T]) Walk(st *Status) *Summary {
	sm := &Summary{
		SummaryStatus: st.SummaryStatus,
	}

	if slice.ContainsAny[T](w.transitioningStatus, any(st.SummaryStatus).(T)) {
		sm.Transitioning = true
	}

	if slice.ContainsAny(w.errorStatus, any(st.SummaryStatus).(T)) {
		sm.Error = true
		sm.Transitioning = false
	}

	for i := range w.decides {
		change, update := w.decides[i](st.Summary)
		if change {
			sm = &update
		}
	}

	return sm
}
