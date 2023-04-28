package status

import "strings"

type (
	// Walker holds the steps and makes a summary of the status's conditions.
	Walker interface {
		// Walk walks all conditions of the status and gives a proper summary.
		Walk(*Status) *Summary
	}

	// Decide returns readable and sensible status by the given condition status and reason,
	// and moves to next path step if both returning `isError` and `isTransitioning` are false.
	Decide func(st ConditionStatus, reason string) (display string, isError bool, isTransitioning bool)
)

// NewWalker creates a stacking walker by the given steps group,
// and applies the customized decision to the steps group.
//   - `stepsGroup` specifies the path steps in line,
//     logically, move to next step if the current step is done.
//     By default, Walker decides to move to the next step on whether the corresponding condition is True status.
//   - `arrange` applies the customized decision,
//     for example, moving to next step on a dedicated step if its status is False,
//     or changing step's display content by its status.
func NewWalker[T ~string](stepsGroup [][]T, arranges ...func(Decision[T])) Walker {
	if len(stepsGroup) == 0 {
		panic("empty steps group")
	}

	var fs = make(paths[T], 0, len(stepsGroup))
	for i := range stepsGroup {
		fs = append(fs, newPath(stepsGroup[i], arranges...))
	}
	return fs
}

// paths stacks a collection of path,
// and picks the highest score result.
type paths[T ~string] []path[T]

func (ps paths[T]) Walk(st *Status) (r *Summary) {
	for i := range ps {
		var l = ps[i].Walk(st)
		if r == nil {
			r = l
			continue
		}

		// accept the result that has a higher score.
		var ls, rs = getSummaryScore(l), getSummaryScore(r)
		if ls <= rs {
			continue
		}
		r, rs = l, ls

		// quit soon if found one highest result.
		if rs == highestSummaryScore {
			break
		}
	}
	return
}

const (
	summaryScoreDone = iota
	summaryScoreTransitioning
	summaryScoreError

	highestSummaryScore = summaryScoreError
)

// getSummaryScore returns the score of the given status summary.
func getSummaryScore(s *Summary) int {
	switch {
	case s.Error:
		return summaryScoreError
	case s.Transitioning:
		return summaryScoreTransitioning
	}
	return summaryScoreDone
}

// newPath creates a path and initializes it.
func newPath[T ~string](steps []T, arranges ...func(Decision[T])) path[T] {
	if len(steps) == 0 {
		panic("empty steps")
	}

	var p = path[T]{
		steps:       steps,
		stepsIndex:  make(map[T]int, len(steps)),
		stepsDecide: make([]Decide, len(steps)),
	}
	for i := range steps {
		// loop check, panic if found.
		if _, exist := p.stepsIndex[steps[i]]; exist {
			panic("found loop")
		}
		p.stepsIndex[steps[i]] = i
		p.stepsDecide[i] = getGeneralDecide(steps[i])
	}

	// change the default decide logic after arranging.
	for i := range arranges {
		arranges[i](Decision[T](p))
	}

	return p
}

// path holds the steps and makes a summary of the status's conditions.
type path[T ~string] struct {
	steps       []T
	stepsIndex  map[T]int
	stepsDecide []Decide
}

func (f path[T]) Walk(st *Status) *Summary {
	var s Summary

	// walk the status if condition list is not empty.
	if len(st.Conditions) != 0 {
		// map conditions with the specified steps for quick indexing.
		var stepsConditionIndex = make([]int, len(f.steps))
		for i, c := range st.Conditions {
			// plus 1 to avoid aligning not found item.
			if idx, exist := f.stepsIndex[T(c.Type)]; exist {
				stepsConditionIndex[idx] = i + 1
			}
		}

		// walk the path to configure the summary.
		for i := range f.steps {
			if stepsConditionIndex[i] == 0 {
				// not found step in the given status's condition list.
				continue
			}
			var c = &st.Conditions[stepsConditionIndex[i]-1]

			// get summary from display result.
			s.SummaryStatus, s.Error, s.Transitioning = f.stepsDecide[i](c.Status, c.Reason)
			s.SummaryStatusMessage = c.Message

			// quit from the walk if still error or being transitioning.
			if s.Error || s.Transitioning {
				break
			}
		}
	}

	// default summary if it hasn't been configured.
	if s.SummaryStatus == "" {
		s.SummaryStatus, s.Error, s.Transitioning = f.stepsDecide[len(f.steps)-1]("", "")
		s.SummaryStatusMessage = ""
	}
	return &s
}

// Decision exposes ability to customize how to make a decision on one specified step.
type Decision[T ~string] path[T]

// Make makes a decision on the given specified step with dedicated decide logic.
func (d Decision[T]) Make(step T, with Decide) Decision[T] {
	if with != nil {
		if idx, exist := d.stepsIndex[step]; exist {
			d.stepsDecide[idx] = with
		}
	}
	return d
}

// getGeneralDecide returns a decision that adapts general scene, including,
//   - displays step pretty,
//   - marks step as error if status is False,
//   - marks step as transitioning if status is Unknown,
//   - and moves to next step if status is True.
func getGeneralDecide[T ~string](step T) Decide {
	var s = string(step)

	// pretty the display with some rules,
	// most rules are for not present tense word.
	var displays = [3]string{s, s, s} // Transitioning, Error, Done
	for m, r := range replacements {
		if !strings.HasSuffix(s, m) {
			continue
		}
		var p = s[:len(s)-len(m)]
		displays[0], displays[1], displays[2] = p+r.T, p+r.E, p+r.D
	}

	return func(st ConditionStatus, _ string) (string, bool, bool) {
		switch st {
		case ConditionStatusUnknown:
			return displays[0], false, true
		case ConditionStatusFalse:
			return displays[1], true, false
		}
		return displays[2], false, false
	}
}

// replacements collects the rules for replacing phased descriptor of the key,
// includes transitioning(T), error(E) and done(D).
var replacements = map[string]struct {
	T, E, D string
}{
	"Progressing": {"Progressing", "Progressing", "Progressed"},
	"Provisioned": {"Provisioning", "ProvisionFailed", "Provisioned"},
	"Initialized": {"Initializing", "InitializeFailed", "Initialized"},
	"Scheduled":   {"Scheduling", "ScheduleFailed", "Scheduled"},
	"Accepted":    {"Accepting", "NotAccepted", "Accepted"},
	"Deployed":    {"Deploying", "DeployFailed", "Deployed"},
	"Synced":      {"Syncing", "SyncFailed", "Synced"},
	"Available":   {"Preparing", "Unavailable", "Available"},
	"Ready":       {"Preparing", "Unready", "Ready"},
	"Active":      {"Preparing", "Inactive", "Active"},
}
