package kubestatus

import (
	apps "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v1"
	batch "k8s.io/api/batch/v1"
	certificates "k8s.io/api/certificates/v1"
	core "k8s.io/api/core/v1"
	policy "k8s.io/api/policy/v1"
	apiregistration "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	"github.com/seal-io/walrus/pkg/dao/types/status"
)

const (
	displayFailed      = "Failed"
	displaySigning     = "Signing"
	displayProgressing = "Progressing"
)

// podStatusPaths makes the following decision.
//
//	| Condition Type   |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| ---------------- | ----------------------- | --------------------- | --------------------- |
//	| Initialized      | Unknown                 | Initializing          | Transitioning         |
//	| Initialized      | False                   | InitializeFailed      | Error                 |
//	| Initialized      | True                    | Initialized           |                       |
//	| PodScheduled     | Unknown                 | PodScheduling         | Transitioning         |
//	| PodScheduled     | False                   | PodScheduleFailed     | Error                 |
//	| PodScheduled     | True                    | PodScheduled          |                       |
//	| ContainersReady  | Unknown                 | ContainersPreparing   | Transitioning         |
//	| ContainersReady  | False                   | ContainersNotReady    | Error                 |
//	| ContainersReady  | True                    | ContainersReady       |                       |
//	| DisruptionTarget | Unknown                 | Evicting              | Transitioning         |
//	| DisruptionTarget | False                   | Preparing             |                       |
//	| DisruptionTarget | True                    | Evicted               | Error                 |
//	| Ready            | Unknown                 | Preparing             | Transitioning         |
//	| Ready            | False                   | NotReady              | Error                 |
//	| Ready            | True                    | Ready                 |                       |
var podStatusPaths = status.NewWalker(
	[][]core.PodConditionType{
		{
			core.PodInitialized,
			core.PodScheduled,
			core.ContainersReady,
			core.DisruptionTarget,
			core.PodReady,
		},
	},
	func(d status.Decision[core.PodConditionType]) {
		const (
			reasonContainersNotInitialized = "ContainersNotInitialized"
			reasonPodCompleted             = "PodCompleted"
		)

		d.Make(core.PodInitialized,
			func(st status.ConditionStatus, reason string) *status.Summary {
				switch st {
				case status.ConditionStatusTrue:
					return &status.Summary{SummaryStatus: "Initialized"}
				case status.ConditionStatusFalse:
					if reason == reasonContainersNotInitialized {
						return &status.Summary{
							SummaryStatus: "Initializing",
							Transitioning: true,
						}
					}
					return &status.Summary{
						SummaryStatus: "InitializeFailed",
						Error:         true,
					}
				}
				return &status.Summary{
					SummaryStatus: "Initializing",
					Transitioning: true,
				}
			})

		d.Make(core.ContainersReady,
			func(st status.ConditionStatus, reason string) *status.Summary {
				switch st {
				case status.ConditionStatusTrue:
					return &status.Summary{SummaryStatus: "ContainersReady"}
				case status.ConditionStatusFalse:
					if reason == reasonPodCompleted {
						// Completed job.
						return &status.Summary{
							SummaryStatus: "ContainersCompleted",
							Inactive:      true,
						}
					}
					return &status.Summary{
						SummaryStatus: "ContainersNotReady",
						Error:         true,
					}
				}
				return &status.Summary{
					SummaryStatus: "ContainersPreparing",
					Transitioning: true,
				}
			})

		d.Make(core.PodReady,
			func(st status.ConditionStatus, reason string) *status.Summary {
				switch st {
				case status.ConditionStatusTrue:
					return &status.Summary{SummaryStatus: "Ready"}
				case status.ConditionStatusFalse:
					if reason == reasonPodCompleted {
						// Completed job.
						return &status.Summary{
							SummaryStatus: "Completed",
							Inactive:      true,
						}
					}
					return &status.Summary{
						SummaryStatus: "NotReady",
						Error:         true,
					}
				}
				return &status.Summary{
					SummaryStatus: "Preparing",
					Transitioning: true,
				}
			})

		d.Make(core.DisruptionTarget,
			func(st status.ConditionStatus, reason string) *status.Summary {
				switch st {
				case status.ConditionStatusTrue:
					return &status.Summary{
						SummaryStatus: "Evicted",
						Error:         true,
					}
				case status.ConditionStatusFalse:
					return &status.Summary{SummaryStatus: "Preparing"}
				}
				return &status.Summary{
					SummaryStatus: "Evicting",
					Transitioning: true,
				}
			})
	},
)

// replicaSetStatusPaths makes the following decision.
//
//	| Condition Type |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| -------------- | ----------------------- | --------------------- | --------------------- |
//	| ReplicaFailure | Unknown                 | ReplicaDeploying      | Transitioning         |
//	| ReplicaFailure | False                   | ReplicaDeployed       |                       |
//	| ReplicaFailure | True                    | ReplicaDeployFailed   | Error                 |
var replicaSetStatusPaths = status.NewWalker(
	[][]apps.ReplicaSetConditionType{
		{
			apps.ReplicaSetReplicaFailure,
		},
	},
	func(d status.Decision[apps.ReplicaSetConditionType]) {
		d.Make(
			apps.ReplicaSetReplicaFailure,
			func(st status.ConditionStatus, reason string) *status.Summary {
				switch st {
				case status.ConditionStatusFalse:
					return &status.Summary{SummaryStatus: "ReplicaDeployed"}
				case status.ConditionStatusTrue:
					return &status.Summary{
						SummaryStatus: "ReplicaDeployFailed",
						Error:         true,
					}
				}
				return &status.Summary{
					SummaryStatus: "ReplicaDeploying",
					Transitioning: true,
				}
			},
		)
	},
)

// apiServiceStatusPaths makes the following decision.
//
//	| Condition Type |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| -------------- | ----------------------- | --------------------- | --------------------- |
//	| Available      | Unknown                 | Preparing             | Transitioning         |
//	| Available      | False                   | Unavailable           | Error                 |
//	| Available      | True                    | Available             |                       |
var apiServiceStatusPaths = status.NewWalker(
	[][]apiregistration.APIServiceConditionType{
		{
			apiregistration.Available,
		},
	},
)

// deploymentStatusPaths makes the following decision.
//
//	| Condition Type |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| -------------- | ----------------------- | --------------------- | --------------------- |
//	| Progressing    | Unknown                 | Progressing           | Transitioning         |
//	| Progressing    | False                   | Progressing           | Error                 |
//	| Progressing    | True(ReplicaSetUpdated) | Progressing           | Transitioning         |
//	| Progressing    | True(DeploymentPaused)  | Pausing               | Transitioning         |
//	| Progressing    | True                    | Progressed            |                       |
//	| ReplicaFailure | Unknown                 | ReplicaDeploying      | Transitioning         |
//	| ReplicaFailure | False                   | ReplicaDeployed       |                       |
//	| ReplicaFailure | True                    | ReplicaDeployFailed   | Error                 |
//	| Available      | Unknown                 | Preparing             | Transitioning         |
//	| Available      | False                   | Unavailable           | Error                 |
//	| Available      | True                    | Available             |                       |
var deploymentStatusPaths = status.NewWalker(
	[][]apps.DeploymentConditionType{
		{
			apps.DeploymentProgressing,
			apps.DeploymentReplicaFailure,
			apps.DeploymentAvailable,
		},
	},
	func(d status.Decision[apps.DeploymentConditionType]) {
		d.Make(apps.DeploymentProgressing,
			func(st status.ConditionStatus, reason string) *status.Summary {
				if st == status.ConditionStatusTrue && reason != "ReplicaSetUpdated" {
					return &status.Summary{SummaryStatus: "Progressed"}
				}
				if st == status.ConditionStatusUnknown && reason == "DeploymentPaused" {
					return &status.Summary{
						SummaryStatus: "Pausing",
						Transitioning: true,
					}
				}
				return &status.Summary{
					SummaryStatus: displayProgressing,
					Error:         st == status.ConditionStatusFalse,
					Transitioning: st != status.ConditionStatusFalse,
				}
			})

		d.Make(apps.DeploymentReplicaFailure,
			func(st status.ConditionStatus, reason string) *status.Summary {
				switch st {
				case status.ConditionStatusFalse:
					return &status.Summary{SummaryStatus: "ReplicaDeployed"}
				case status.ConditionStatusTrue:
					return &status.Summary{
						SummaryStatus: "ReplicaDeployFailed",
						Error:         true,
					}
				}
				return &status.Summary{
					SummaryStatus: "ReplicaDeploying",
					Transitioning: true,
				}
			})
	},
)

// jobStatusPaths makes the following decision.
//
//	| Condition Type |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| -------------- | ----------------------- | --------------------- | --------------------- |
//	| Suspended      | Unknown                 | Progressing           | Transitioning         |
//	| Suspended      | False(JobResumed)       | Progressing           |                       |
//	| Suspended      | False                   | Progressing           | Error                 |
//	| Suspended      | True                    | Suspending            | Transitioning         |
//	| FailureTarget  | Unknown                 | Progressing           | Transitioning         |
//	| FailureTarget  | False                   | Progressing           |                       |
//	| FailureTarget  | True                    | Failed                | Error                 |
//	| Failed         | Unknown                 | Progressing           | Transitioning         |
//	| Failed         | False                   | Progressing           |                       |
//	| Failed         | True                    | Failed                | Error                 |
//	| Complete       | Unknown                 | Progressing           | Transitioning         |
//	| Complete       | False                   | Progressing           | Error                 |
//	| Complete       | True                    | Completed             |                       |
var jobStatusPaths = status.NewWalker(
	[][]batch.JobConditionType{
		{
			batch.JobSuspended,
			batch.JobFailureTarget,
			batch.JobFailed,
			batch.JobComplete,
		},
	},
	func(d status.Decision[batch.JobConditionType]) {
		d.Make(batch.JobSuspended,
			func(st status.ConditionStatus, reason string) *status.Summary {
				switch st {
				case status.ConditionStatusTrue:
					return &status.Summary{
						SummaryStatus: "Suspended",
						Transitioning: true,
					}
				case status.ConditionStatusFalse:
					if reason != "JobResumed" {
						return &status.Summary{
							SummaryStatus: displayProgressing,
							Error:         true,
						}
					}
				}
				return &status.Summary{
					SummaryStatus: displayProgressing,
					Transitioning: st == status.ConditionStatusUnknown,
				}
			})

		d.Make(batch.JobFailureTarget,
			func(st status.ConditionStatus, reason string) *status.Summary {
				if st == status.ConditionStatusTrue {
					return &status.Summary{
						SummaryStatus: "Failed",
						Error:         true,
					}
				}
				return &status.Summary{
					SummaryStatus: displayProgressing,
					Error:         st == status.ConditionStatusFalse,
					Transitioning: st == status.ConditionStatusUnknown,
				}
			})

		d.Make(batch.JobFailed,
			func(st status.ConditionStatus, reason string) *status.Summary {
				if st == status.ConditionStatusTrue {
					return &status.Summary{
						SummaryStatus: displayFailed,
						Error:         true,
					}
				}
				return &status.Summary{
					SummaryStatus: displayProgressing,
					Error:         st == status.ConditionStatusFalse,
					Transitioning: st == status.ConditionStatusUnknown,
				}
			})

		d.Make(batch.JobComplete,
			func(st status.ConditionStatus, reason string) *status.Summary {
				if st == status.ConditionStatusTrue {
					return &status.Summary{
						SummaryStatus: "Completed",
						Inactive:      true,
					}
				}
				return &status.Summary{
					SummaryStatus: displayProgressing,
					Error:         st == status.ConditionStatusFalse,
					Transitioning: st == status.ConditionStatusUnknown,
				}
			})
	},
)

// hpaStatusPaths makes the following decision.
//
//	| Condition Type |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| -------------- | ----------------------- | --------------------- | --------------------- |
//	| ScalingLimited | Unknown                 | Scaling               | Transitioning         |
//	| ScalingLimited | False                   | Scaling               |                       |
//	| ScalingLimited | True                    | ScalingLimited        | Error                 |
//	| AbleToScale    | Unknown                 | Scaling               | Transitioning         |
//	| AbleToScale    | False                   | Scaling               | Error                 |
//	| AbleToScale    | True(SucceededRescale)  | Scaled                |                       |
//	| AbleToScale    | True                    | Scaling               | Transitioning         |
//	| ScalingActive  | Unknown                 | ScalingPreparing      | Transitioning         |
//	| ScalingActive  | False                   | ScalingInactive       | Error                 |
//	| ScalingActive  | True                    | ScalingActive         |                       |
var hpaStatusPaths = status.NewWalker(
	[][]autoscaling.HorizontalPodAutoscalerConditionType{
		{
			autoscaling.ScalingLimited,
			autoscaling.AbleToScale,
			autoscaling.ScalingActive,
		},
	},
	func(d status.Decision[autoscaling.HorizontalPodAutoscalerConditionType]) {
		d.Make(autoscaling.ScalingLimited,
			func(st status.ConditionStatus, reason string) *status.Summary {
				if st == status.ConditionStatusTrue {
					return &status.Summary{
						SummaryStatus: "ScalingLimited",
						Error:         true,
					}
				}
				return &status.Summary{
					SummaryStatus: "Scaling",
					Transitioning: st == status.ConditionStatusUnknown,
				}
			})

		d.Make(autoscaling.AbleToScale,
			func(st status.ConditionStatus, reason string) *status.Summary {
				if st == status.ConditionStatusTrue && reason == "SucceededRescale" {
					return &status.Summary{SummaryStatus: "Scaled"}
				}
				return &status.Summary{
					SummaryStatus: "Scaling",
					Error:         st == status.ConditionStatusFalse,
					Transitioning: true,
				}
			})
	},
)

// csrStatusPaths makes the following decision.
//
//	| Condition Type |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| -------------- | ----------------------- | --------------------- | --------------------- |
//	| Failed         | Unknown                 | Signing               | Transitioning         |
//	| Failed         | False                   | Signing               |                       |
//	| Failed         | True                    | Failed                | Error                 |
//	| Denied         | Unknown                 | Signing               | Transitioning         |
//	| Denied         | False                   | Signing               |                       |
//	| Denied         | True                    | Denied                | Error                 |
//	| Approved       | Unknown                 | Signing               | Transitioning         |
//	| Approved       | False                   | Signing               | Error                 |
//	| Approved       | True                    | Approved              |                       |
var csrStatusPaths = status.NewWalker(
	[][]certificates.RequestConditionType{
		{
			certificates.CertificateFailed,
			certificates.CertificateDenied,
			certificates.CertificateApproved,
		},
	},
	func(d status.Decision[certificates.RequestConditionType]) {
		d.Make(certificates.CertificateFailed,
			func(st status.ConditionStatus, reason string) *status.Summary {
				if st == status.ConditionStatusTrue {
					return &status.Summary{
						SummaryStatus: displayFailed,
						Error:         true,
					}
				}
				return &status.Summary{
					SummaryStatus: displaySigning,
					Transitioning: st == status.ConditionStatusUnknown,
				}
			})

		d.Make(certificates.CertificateDenied,
			func(st status.ConditionStatus, reason string) *status.Summary {
				if st == status.ConditionStatusTrue {
					return &status.Summary{
						SummaryStatus: "Denied",
						Error:         true,
					}
				}
				return &status.Summary{
					SummaryStatus: displaySigning,
					Transitioning: st == status.ConditionStatusUnknown,
				}
			})

		d.Make(certificates.CertificateApproved,
			func(st status.ConditionStatus, reason string) *status.Summary {
				if st == status.ConditionStatusTrue {
					return &status.Summary{SummaryStatus: "Approved"}
				}
				return &status.Summary{
					SummaryStatus: displaySigning,
					Error:         st == status.ConditionStatusFalse,
					Transitioning: st == status.ConditionStatusUnknown,
				}
			})
	},
)

// pdbStatusPaths makes the following decision.
//
//	|  Condition Type   |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| ----------------- | ----------------------- | --------------------- | --------------------- |
//	| DisruptionAllowed | Unknown                 | Preparing             | Transitioning         |
//	| DisruptionAllowed | False(InsufficientPods) | Active                |                       |
//	| DisruptionAllowed | False                   | Inactive              | Error                 |
//	| DisruptionAllowed | True                    | Active                |                       |
var pdbStatusPaths = status.NewWalker(
	[][]string{
		{
			policy.DisruptionAllowedCondition,
		},
	},
	func(d status.Decision[string]) {
		d.Make(policy.DisruptionAllowedCondition,
			func(st status.ConditionStatus, reason string) *status.Summary {
				switch st {
				case status.ConditionStatusTrue:
					return &status.Summary{SummaryStatus: "Active"}
				case status.ConditionStatusFalse:
					if reason == "InsufficientPods" {
						return &status.Summary{SummaryStatus: "Active"}
					}
					return &status.Summary{
						SummaryStatus: "Inactive",
						Error:         true,
					}
				}
				return &status.Summary{
					SummaryStatus: "Preparing",
					Transitioning: true,
				}
			})
	},
)
