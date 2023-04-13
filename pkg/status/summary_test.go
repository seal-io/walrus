package status

import (
	"testing"

	core "k8s.io/api/core/v1"

	"github.com/seal-io/seal/pkg/dao/types/status"
)

var testingStatusCases = []struct {
	name  string
	input status.Status
}{
	{
		name:  "No Conditions",
		input: status.Status{},
	},
	{
		name: "One Conditions",
		input: status.Status{
			Conditions: []status.Condition{
				{
					Type:   status.ConditionType(core.PodInitialized),
					Status: status.ConditionStatusUnknown,
				},
			},
		},
	},
	{
		name: "Half Conditions",
		input: status.Status{
			Conditions: []status.Condition{
				{
					Type:   status.ConditionType(core.PodInitialized),
					Status: status.ConditionStatusTrue,
				},
				{
					Type:   status.ConditionType(core.PodScheduled),
					Status: status.ConditionStatusTrue,
				},
				{
					Type:   status.ConditionType(core.ContainersReady),
					Status: status.ConditionStatusUnknown,
				},
			},
		},
	},
	{
		name: "All Conditions",
		input: status.Status{
			Conditions: []status.Condition{
				{
					Type:   status.ConditionType(core.PodInitialized),
					Status: status.ConditionStatusTrue,
				},
				{
					Type:   status.ConditionType(core.PodScheduled),
					Status: status.ConditionStatusTrue,
				},
				{
					Type:   status.ConditionType(core.ContainersReady),
					Status: status.ConditionStatusTrue,
				},
				{
					Type:   status.ConditionType(core.ContainersReady),
					Status: status.ConditionStatusTrue,
				},
			},
		},
	},
	{
		name: "Last Failed Step Conditions",
		input: status.Status{
			Conditions: []status.Condition{
				{
					Type:   status.ConditionType(core.PodInitialized),
					Status: status.ConditionStatusTrue,
				},
				{
					Type:   status.ConditionType(core.PodScheduled),
					Status: status.ConditionStatusTrue,
				},
				{
					Type:   status.ConditionType(core.ContainersReady),
					Status: status.ConditionStatusTrue,
				},
				{
					Type:   status.ConditionType(core.ContainersReady),
					Status: status.ConditionStatusFalse,
				},
			},
		},
	},
	{
		name: "Out-of-Sync Conditions",
		input: status.Status{
			Conditions: []status.Condition{
				{
					Type:   status.ConditionType(core.PodInitialized),
					Status: status.ConditionStatusTrue,
				},
				{
					Type:   status.ConditionType(core.PodScheduled),
					Status: status.ConditionStatusTrue,
				},
				{
					Type:   status.ConditionType(core.ContainersReady),
					Status: status.ConditionStatusFalse,
				},
				{
					Type:   status.ConditionType(core.PodReady),
					Status: status.ConditionStatusFalse,
				},
				{
					Type:   status.ConditionType(core.DisruptionTarget),
					Status: status.ConditionStatusFalse,
				},
			},
		},
	},
}

func newSummarizer() *Summarizer {
	var s = NewSummarizer(status.ConditionType(core.PodReady))
	s.AddErrorFalseTransitioningUnknown(status.ConditionType(core.PodInitialized), string(core.PodInitialized))
	s.AddErrorFalseTransitioningUnknown(status.ConditionType(core.PodScheduled), string(core.PodScheduled))
	s.AddErrorFalseTransitioningUnknown(status.ConditionType(core.ContainersReady), string(core.ContainersReady))
	s.AddErrorFalseTransitioningUnknown(status.ConditionType(core.PodReady), string(core.PodReady))
	s.AddErrorTrueTransitioningUnknown(status.ConditionType(core.DisruptionTarget), string(core.DisruptionTarget))
	return s
}

func newWalker() status.Walker {
	return status.NewWalker(
		[][]core.PodConditionType{
			{
				core.PodInitialized,
				core.PodScheduled,
				core.ContainersReady,
				core.PodReady,
				core.DisruptionTarget,
			},
		},
		func(d status.Decision[core.PodConditionType]) {
			d.Make(core.DisruptionTarget,
				func(st status.ConditionStatus, reason string) (display string, isError bool, isTransitioning bool) {
					switch st {
					case status.ConditionStatusTrue:
						return "Evicted", true, false
					case status.ConditionStatusFalse:
						return "Preparing", false, false
					}
					return "Evicting", false, true
				})
		})
}

func BenchmarkSummarizer_NewAndSummarize(b *testing.B) {
	for _, c := range testingStatusCases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				newSummarizer().Summarize(&c.input)
			}
			b.ReportAllocs()
		})
	}
}

func BenchmarkWalker_NewAndWalk(b *testing.B) {
	for _, c := range testingStatusCases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				newWalker().Walk(&c.input)
			}
			b.ReportAllocs()
		})
	}
}

func BenchmarkSummarizer_Summarize(b *testing.B) {
	var s = newSummarizer()
	for _, c := range testingStatusCases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s.Summarize(&c.input)
			}
			b.ReportAllocs()
		})
	}
}

func BenchmarkWalker_Walk(b *testing.B) {
	var w = newWalker()
	for _, c := range testingStatusCases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				w.Walk(&c.input)
			}
			b.ReportAllocs()
		})
	}
}
