package kubestatus

import (
	apps "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v1"
	batch "k8s.io/api/batch/v1"
	certificates "k8s.io/api/certificates/v1"
	core "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	policy "k8s.io/api/policy/v1"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	"github.com/seal-io/seal/pkg/dao/types/status"
	pkgstatus "github.com/seal-io/seal/pkg/status"
)

func podSummarizer() *pkgstatus.Summarizer {
	// init
	summarizer := pkgstatus.NewSummarizer(status.ConditionType(core.PodReady))

	// status False is error, Unknown is transitioning.
	summarizer.AddErrorFalseTransitioningUnknown(status.ConditionType(core.PodInitialized), string(core.PodInitialized))
	summarizer.AddErrorFalseTransitioningUnknown(status.ConditionType(core.PodScheduled), string(core.PodScheduled))
	summarizer.AddErrorFalseTransitioningUnknown(status.ConditionType(core.ContainersReady), string(core.ContainersReady))
	summarizer.AddErrorFalseTransitioningUnknown(status.ConditionType(core.PodReady), string(core.PodReady))
	return summarizer
}

func replicaSetSummarizer() *pkgstatus.Summarizer {
	// init
	summarizer := pkgstatus.NewSummarizer(status.GeneralStatusReady)

	// status True is error, Unknown is transitioning.
	summarizer.AddErrorTrueTransitioningUnknown(status.ConditionType(apps.ReplicaSetReplicaFailure), string(apps.ReplicaSetReplicaFailure))
	return summarizer
}

func apiServiceSummarizer() *pkgstatus.Summarizer {
	// init
	summarizer := pkgstatus.NewSummarizer(status.ConditionType(apiregistrationv1.Available))

	// status False is error, Unknown is transitioning.
	summarizer.AddErrorFalseTransitioningUnknown(status.ConditionType(apiregistrationv1.Available), string(apiregistrationv1.Available))
	return summarizer
}

func deploymentSummarizer() *pkgstatus.Summarizer {
	// init
	summarizer := pkgstatus.NewSummarizer(status.ConditionType(apps.DeploymentAvailable))

	// status False is error, Unknown is transitioning.
	summarizer.AddErrorFalseTransitioningUnknown(status.ConditionType(apps.DeploymentAvailable), string(apps.DeploymentAvailable))

	// status False is error, True is transitioning.
	summarizer.AddErrorFalseTransitioningTrue(status.ConditionType(apps.DeploymentProgressing), string(apps.DeploymentProgressing))

	// status True is error, Unknown is transitioning.
	summarizer.AddErrorTrueTransitioningUnknown(status.ConditionType(apps.DeploymentReplicaFailure), string(apps.DeploymentReplicaFailure))

	summarizer.AddErrorReason(status.ConditionType(apps.DeploymentProgressing),
		"ReplicaSetCreateError",
		"ProgressDeadlineExceeded",
	)

	summarizer.AddTransitionReason(status.ConditionType(apps.DeploymentProgressing),
		"ReplicaSetUpdated",
		"DeploymentPaused",
		"DeploymentResumed",
	)

	summarizer.AddReadyReason(status.ConditionType(apps.DeploymentProgressing),
		"NewReplicaSetCreated",
		"FoundNewReplicaSet",
		"NewReplicaSetAvailable",
	)
	return summarizer
}

func jobSummarizer() *pkgstatus.Summarizer {
	// init
	summarizer := pkgstatus.NewSummarizer(status.ConditionType(batch.JobComplete))

	// status False is error, Unknown is transitioning.
	summarizer.AddErrorFalseTransitioningUnknown(status.ConditionType(batch.JobComplete), string(batch.JobComplete))

	// status True is error, Unknown is transitioning.
	summarizer.AddErrorTrueTransitioningUnknown(status.ConditionType(batch.JobFailed), string(batch.JobFailed))
	summarizer.AddErrorTrueTransitioningUnknown(status.ConditionType(batch.JobFailureTarget), string(batch.JobFailureTarget))

	// status Unknown is error, True is transitioning
	summarizer.AddErrorUnknownTransitioningTrue(status.ConditionType(batch.JobSuspended), string(batch.JobSuspended))
	return summarizer
}

func hpaSummarizer() *pkgstatus.Summarizer {
	// init
	summarizer := pkgstatus.NewSummarizer(status.ConditionType(autoscaling.ScalingActive))

	// status False is error, Unknown is transitioning.
	summarizer.AddErrorFalseTransitioningUnknown(status.ConditionType(autoscaling.ScalingActive), string(autoscaling.ScalingActive))
	summarizer.AddErrorFalseTransitioningUnknown(status.ConditionType(autoscaling.AbleToScale), string(autoscaling.AbleToScale))

	// status True is error, Unknown is transitioning.
	summarizer.AddErrorTrueTransitioningUnknown(status.ConditionType(autoscaling.ScalingLimited), string(autoscaling.ScalingLimited))
	return summarizer
}

func csrSummarizer() *pkgstatus.Summarizer {
	// init
	summarizer := pkgstatus.NewSummarizer(status.ConditionType(certificates.CertificateApproved))

	// status False is error, Unknown is transitioning.
	summarizer.AddErrorFalseTransitioningUnknown(status.ConditionType(certificates.CertificateApproved), string(certificates.CertificateApproved))

	// status True is error, Unknown is transitioning
	summarizer.AddErrorTrueTransitioningUnknown(status.ConditionType(certificates.CertificateDenied), string(certificates.CertificateDenied))
	summarizer.AddErrorTrueTransitioningUnknown(status.ConditionType(certificates.CertificateFailed), string(certificates.CertificateFailed))
	return summarizer
}

func networkPolicySummarizer() *pkgstatus.Summarizer {
	// init
	summarizer := pkgstatus.NewSummarizer(status.ConditionType(networking.NetworkPolicyConditionStatusAccepted))

	// status False is error, Unknown is transitioning.
	summarizer.AddErrorFalseTransitioningUnknown(status.ConditionType(networking.NetworkPolicyConditionStatusAccepted), string(networking.NetworkPolicyConditionStatusAccepted))

	// status True is error, Unknown is transitioning.
	summarizer.AddErrorTrueTransitioningUnknown(status.ConditionType(networking.NetworkPolicyConditionStatusPartialFailure), string(networking.NetworkPolicyConditionStatusPartialFailure))
	summarizer.AddErrorTrueTransitioningUnknown(status.ConditionType(networking.NetworkPolicyConditionStatusFailure), string(networking.NetworkPolicyConditionStatusFailure))
	return summarizer
}

func pdbSummarizer() *pkgstatus.Summarizer {
	// init
	summarizer := pkgstatus.NewSummarizer(policy.DisruptionAllowedCondition)

	// status False is error, Unknown is transitioning.
	summarizer.AddErrorFalseTransitioningUnknown(policy.DisruptionAllowedCondition, policy.DisruptionAllowedCondition)
	return summarizer
}
