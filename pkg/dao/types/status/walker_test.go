package status

import (
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestWalker_sxs(t *testing.T) {
	// 1. define resource with status.
	type ExampleResource struct {
		Status Status
	}

	// 2. define the condition types of ExampleResource,
	// condition type can be past tense or present tense.
	const (
		ExampleResourceStatusProgressing    ConditionType = "Progressing"
		ExampleResourceStatusReplicaFailure ConditionType = "ReplicaFailure"
		ExampleResourceStatusAvailable      ConditionType = "Available"
	)

	// 2.1  clarify the condition type and its status meaning as below.
	//      | Condition Type |     Condition Status    | Human Readable Status | Human Sensible Status |
	//      | -------------- | ----------------------- | --------------------- | --------------------- |
	//      | Progressing    | Unknown                 | Progressing           | Transitioning         |
	//      | Progressing    | False                   | Progressing           | Error                 |
	//      | Progressing    | True(ReplicaSetUpdated) | Progressing           | Transitioning         |
	//      | Progressing    | True(DeploymentPaused)  | Pausing               | Transitioning         |
	//      | Progressing    | True                    | Progressed            | Done                  |
	//      | ReplicaFailure | Unknown                 | ReplicaDeploying      | Transitioning         |
	//      | ReplicaFailure | False                   | ReplicaDeployed       | Done                  |
	//      | ReplicaFailure | True                    | ReplicaDeployFailed   | Error                 |
	//      | Available      | Unknown                 | Preparing             | Transitioning         |
	//      | Available      | False                   | Unavailable           | Error                 |
	//      | Available      | True                    | Available             | Done                  |

	// 3. create a flow to connect the above condition types.
	var f = NewWalker(
		// define paths.
		[][]ConditionType{
			{
				ExampleResourceStatusProgressing,
				ExampleResourceStatusReplicaFailure,
				ExampleResourceStatusAvailable,
			},
		},
		// arrange the default step decision logic.
		func(d Decision[ConditionType]) {
			d.Make(ExampleResourceStatusProgressing,
				func(st ConditionStatus, reason string) (display string, isError bool, isTransitioning bool) {
					if st == ConditionStatusTrue && reason != "ReplicaSetUpdated" {
						return "Progressed", false, false
					}
					if st == ConditionStatusUnknown && reason == "DeploymentPaused" {
						return "Pausing", false, true
					}
					return "Progressing", st == ConditionStatusFalse, st != ConditionStatusFalse
				})

			d.Make(ExampleResourceStatusReplicaFailure,
				func(st ConditionStatus, reason string) (display string, isError bool, isTransitioning bool) {
					switch st {
					case ConditionStatusFalse:
						return "ReplicaDeployed", false, false
					case ConditionStatusTrue:
						return "ReplicaDeployFailed", true, false
					}
					return "ReplicaDeploying", false, true
				})
		},
	)

	var p printer

	// 4. create an instance of ExampleResource.
	var r ExampleResource
	// 4.1  at beginning, the status is empty(we haven't configured any conditions or summary result),
	//      the path will walk to the end step and display the info of the last step,
	//      so we should get a done available summary,
	//      which can treat as Default Status.
	p.Dump("Default Available [D]", f.Walk(&r.Status))
	// 4.2  marked the "Progressing" status to Unknown, which means progressing,
	//      we should get a transitioning progressing summary.
	ExampleResourceStatusProgressing.Unknown(&r, "")
	p.Dump("Progressing [T]", f.Walk(&r.Status))
	// 4.3  marked the "Progressing" status to True with ReplicaSetUpdated reason,
	//      we should still get a transitioning progressing summary.
	r.Status.Conditions[0].Status = ConditionStatusTrue
	r.Status.Conditions[0].Reason = "ReplicaSetUpdated"
	p.Dump("Still Progressing [T]", f.Walk(&r.Status))
	// 4.4  marked the "Progressing" reason to NewReplicaSetAvailable,
	//      we should get a done progressing summary.
	//      at the same time, we haven't configured other conditions,
	//      so we only can see the progressing result.
	r.Status.Conditions[0].Reason = "NewReplicaSetAvailable"
	p.Dump("Progressed [D]", f.Walk(&r.Status))
	// 4.5  marked the "ReplicaFailure" status to Unknown, which means replica deploying,
	//      we should get a transitioning replica deploying summary.
	ExampleResourceStatusReplicaFailure.Unknown(&r, "")
	p.Dump("Replica Deploying [T]", f.Walk(&r.Status))
	// 4.6  marked the "ReplicaFailure" status to True, which means replica deploying failed,
	//      we should get a failed replica deploy summary.
	ExampleResourceStatusReplicaFailure.True(&r, "")
	p.Dump("Replica Deploy Failed [E]", f.Walk(&r.Status))
	// 4.7  marked the "Available" status to Unknown,
	//      we still get a failed replica deploy summary,
	//      as the path cannot move the next step as the "ReplicaFailure" step is not False.
	ExampleResourceStatusAvailable.Unknown(&r, "")
	p.Dump("Still Replica Deploy Failed [E]", f.Walk(&r.Status))
	// 4.8  until marked the "ReplicaFailure" status to False or remove "ReplicaFailure" condition,
	//      we will get a transitioning preparing summary.
	ExampleResourceStatusReplicaFailure.False(&r, "")
	p.Dump("Preparing [T]", f.Walk(&r.Status))
	// 4.9  marked the "Available" status to False, which means replica deploying failed,
	//      we should get an error unavailable summary.
	ExampleResourceStatusAvailable.False(&r, "")
	p.Dump("Unavailable [E]", f.Walk(&r.Status))
	// 4.10 marked the "Progressing" status to Unknown, which means progressing again,
	//      we should get a transitioning progressing summary.
	ExampleResourceStatusProgressing.Unknown(&r, "")
	p.Dump("Progressing Again [T]", f.Walk(&r.Status))

	t.Log(p.String())
}

type printer struct {
	sb strings.Builder
}

func (p *printer) Dump(title string, s *Summary) {
	p.sb.WriteString(title)
	p.sb.WriteString(": ")
	spew.Fdump(&p.sb, s)
	p.sb.WriteString("\n")
}

func (p *printer) String() string {
	return p.sb.String()
}
