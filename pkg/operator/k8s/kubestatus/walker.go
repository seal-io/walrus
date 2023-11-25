package kubestatus

import (
	apps "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v1"
	batch "k8s.io/api/batch/v1"
	certificates "k8s.io/api/certificates/v1"
	core "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
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
//	| Ready            | Unknown                 | Preparing             | Transitioning         |
//	| Ready            | False                   | NotReady              | Error                 |
//	| Ready            | True                    | Ready                 |                       |
//	| DisruptionTarget | Unknown                 | Evicting              | Transitioning         |
//	| DisruptionTarget | False                   | Preparing             |                       |
//	| DisruptionTarget | True                    | Evicted               | Error                 |
var podStatusPaths = status.NewWalker(
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
		const reasonPodCompleted = "PodCompleted"

		d.Make(core.ContainersReady,
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				switch st {
				case status.ConditionStatusTrue:
					return "ContainersReady", false, false
				case status.ConditionStatusFalse:
					if reason == reasonPodCompleted {
						// Completed job.
						return "ContainersCompleted", false, false
					}
					return "ContainersNotReady", true, false
				}
				return "ContainersPreparing", false, true
			})

		d.Make(core.PodReady,
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				switch st {
				case status.ConditionStatusTrue:
					return "Ready", false, false
				case status.ConditionStatusFalse:
					if reason == reasonPodCompleted {
						// Completed job.
						return "Completed", false, false
					}
					return "NotReady", true, false
				}
				return "Preparing", false, true //nolint: goconst
			})

		d.Make(core.DisruptionTarget,
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				switch st {
				case status.ConditionStatusTrue:
					return "Evicted", true, false
				case status.ConditionStatusFalse:
					return "Preparing", false, false
				}
				return "Evicting", false, true
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
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				switch st {
				case status.ConditionStatusFalse:
					return "ReplicaDeployed", false, false
				case status.ConditionStatusTrue:
					return "ReplicaDeployFailed", true, false
				}
				return "ReplicaDeploying", false, true
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
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				if st == status.ConditionStatusTrue && reason != "ReplicaSetUpdated" {
					return "Progressed", false, false
				}
				if st == status.ConditionStatusUnknown && reason == "DeploymentPaused" {
					return "Pausing", false, true
				}
				return displayProgressing, st == status.ConditionStatusFalse, st != status.ConditionStatusFalse
			})

		d.Make(apps.DeploymentReplicaFailure,
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				switch st {
				case status.ConditionStatusFalse:
					return "ReplicaDeployed", false, false
				case status.ConditionStatusTrue:
					return "ReplicaDeployFailed", true, false
				}
				return "ReplicaDeploying", false, true
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
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				switch st {
				case status.ConditionStatusTrue:
					return "Suspending", false, true
				case status.ConditionStatusFalse:
					if reason != "JobResumed" {
						return displayProgressing, true, false
					}
				}
				return displayProgressing, false, st == status.ConditionStatusUnknown
			})

		d.Make(batch.JobFailureTarget,
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				if st == status.ConditionStatusTrue {
					return displayFailed, true, false
				}
				return displayProgressing, st == status.ConditionStatusFalse, st == status.ConditionStatusUnknown
			})

		d.Make(batch.JobFailed,
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				if st == status.ConditionStatusTrue {
					return displayFailed, true, false
				}
				return displayProgressing, st == status.ConditionStatusFalse, st == status.ConditionStatusUnknown
			})

		d.Make(batch.JobComplete,
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				if st == status.ConditionStatusTrue {
					return "Completed", false, false
				}
				return displayProgressing, st == status.ConditionStatusFalse, st == status.ConditionStatusUnknown
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
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				if st == status.ConditionStatusTrue {
					return "ScalingLimited", true, false
				}
				return "Scaling", false, st == status.ConditionStatusUnknown
			})

		d.Make(autoscaling.AbleToScale,
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				if st == status.ConditionStatusTrue && reason == "SucceededRescale" {
					return "Scaled", false, false
				}
				return "Scaling", st == status.ConditionStatusFalse, true
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
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				if st == status.ConditionStatusTrue {
					return displayFailed, true, false
				}
				return displaySigning, false, st == status.ConditionStatusUnknown
			})

		d.Make(certificates.CertificateDenied,
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				if st == status.ConditionStatusTrue {
					return "Denied", true, false
				}
				return displaySigning, false, st == status.ConditionStatusUnknown
			})

		d.Make(certificates.CertificateApproved,
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				if st == status.ConditionStatusTrue {
					return "Approved", false, false
				}
				return displaySigning, st == status.ConditionStatusFalse, st == status.ConditionStatusUnknown
			})
	},
)

// networkPolicyStatusPaths makes the following decision.
//
//	| Condition Type |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| -------------- | ----------------------- | --------------------- | --------------------- |
//	| PartialFailure | Unknown                 | Accepting             | Transitioning         |
//	| PartialFailure | False                   | Accepting             |                       |
//	| PartialFailure | True                    | PartialFailed         | Error                 |
//	| Failure        | Unknown                 | Accepting             | Transitioning         |
//	| Failure        | False                   | Accepting             |                       |
//	| Failure        | True                    | Failed                | Error                 |
//	| Accepted       | Unknown                 | Accepting             | Transitioning         |
//	| Accepted       | False                   | NotAccepted           | Error                 |
//	| Accepted       | True                    | Accepted              |                       |
var networkPolicyStatusPaths = status.NewWalker(
	[][]networking.NetworkPolicyConditionType{
		{
			networking.NetworkPolicyConditionStatusPartialFailure,
			networking.NetworkPolicyConditionStatusFailure,
			networking.NetworkPolicyConditionStatusAccepted,
		},
	},
	func(d status.Decision[networking.NetworkPolicyConditionType]) {
		d.Make(networking.NetworkPolicyConditionStatusPartialFailure,
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				if st == status.ConditionStatusTrue {
					return "PartialFailed", true, false
				}
				return "Accepting", false, st == status.ConditionStatusUnknown
			})

		d.Make(networking.NetworkPolicyConditionStatusFailure,
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				if st == status.ConditionStatusTrue {
					return displayFailed, true, false
				}
				return "Accepting", false, st == status.ConditionStatusUnknown
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
			func(st status.ConditionStatus, reason string) (display string, isError, isTransitioning bool) {
				switch st {
				case status.ConditionStatusTrue:
					return "Active", false, false
				case status.ConditionStatusFalse:
					if reason == "InsufficientPods" {
						return "Active", false, false
					}
					return "Inactive", true, false
				}
				return "Preparing", false, true
			})
	},
)
