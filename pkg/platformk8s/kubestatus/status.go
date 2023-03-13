package kubestatus

import (
	"context"
	"errors"
	"fmt"

	apps "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v2"
	core "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	policy "k8s.io/api/policy/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	"github.com/seal-io/seal/pkg/platformk8s/polymorphic"
)

type GeneralStatus = string

const (
	GeneralStatusUnknown GeneralStatus = "Unknown"
	GeneralStatusReady   GeneralStatus = "Ready"
	GeneralStatusUnready GeneralStatus = "Unready"
)

type GetOptions struct {
	WithJobs     bool
	IgnorePaused bool
}

// Get returns status of the given k8s resource.
func Get(ctx context.Context, dynamicCli *dynamic.DynamicClient, o *unstructured.Unstructured, opts GetOptions) (string, error) {
	switch o.GetKind() {
	case "Service":
		return getService(ctx, o)
	case "Pod":
		return getPod(ctx, o)
	case "PersistentVolume", "PersistentVolumeClaim":
		return getPersistentVolume(ctx, o)
	case "APIService":
		return getAPIService(ctx, o)
	case "Deployment":
		return getDeployment(ctx, o, opts.IgnorePaused)
	case "DaemonSet":
		return getDaemonSet(ctx, o)
	case "StatefulSet":
		return getStatefulSet(ctx, o)
	case "ReplicationController", "ReplicaSet":
		return getReplicas(ctx, dynamicCli, o)
	case "Job":
		if opts.WithJobs {
			return getJob(ctx, o, opts.IgnorePaused)
		}
	case "HorizontalPodAutoscaler":
		return getHorizontalPodAutoscaler(ctx, o)
	case "CertificateSigningRequest":
		return getCertificateSigningRequest(ctx, o)
	case "Ingress":
		return getIngress(ctx, o)
	case "NetworkPolicy":
		return getNetworkPolicy(ctx, o)
	case "PodDisruptionBudget":
		return getPodDisruptionBudget(ctx, o)
	case "ValidatingWebhookConfiguration", "MutatingWebhookConfiguration":
		return getWebhookConfiguration(ctx, dynamicCli, o)
	}

	return GeneralStatusReady, nil
}

// getService returns the status of kubernetes service resource,
// status is one of "Unknown", "Unready", "Ready".
func getService(ctx context.Context, o *unstructured.Unstructured) (string, error) {
	spec, exist, _ := unstructured.NestedMap(o.Object, "spec")
	if !exist {
		return GeneralStatusUnknown, errors.New("not found 'svc' spec")
	}

	// if .spec.type == ExternalName, then ready.
	var specType, _, _ = unstructured.NestedString(spec, "type")
	if core.ServiceType(specType) == core.ServiceTypeExternalName {
		return GeneralStatusReady, nil
	}

	// if .spec.clusterIP == "", then unready.
	var specClusterIP, _, _ = unstructured.NestedString(spec, "clusterIP")
	if specClusterIP == "" {
		return GeneralStatusUnready, nil
	}

	// if .spec.type == LoadBalancer && len(.spec.externalIPs) > 0, then ready.
	// if .spec.type == LoadBalancer && len(.status.loadBalancer.ingress) > 0, then ready.
	if core.ServiceType(specType) == core.ServiceTypeLoadBalancer {
		var specExternalIPs, _, _ = unstructured.NestedStringSlice(spec, "externalIPs")
		if len(specExternalIPs) > 0 {
			return GeneralStatusReady, nil
		}
		var statusLBIngresses, _, _ = unstructured.NestedSlice(o.Object, "status", "loadBalancer", "ingress")
		if len(statusLBIngresses) > 0 {
			return GeneralStatusReady, nil
		}
		return GeneralStatusUnready, nil
	}

	// otherwise, ready.
	return GeneralStatusReady, nil
}

// getPod returns the status of kubernetes pod resource,
// status is one of "Unknown", "Unready", "Pending", "Running", "Succeeded", "Failed".
func getPod(ctx context.Context, o *unstructured.Unstructured) (string, error) {
	status, exist, _ := unstructured.NestedMap(o.Object, "status")
	if !exist {
		return GeneralStatusUnknown, errors.New("not found 'pod' status")
	}

	// if .status.phase != Running, then .status.phase.
	var statusPhase, _, _ = unstructured.NestedString(status, "phase")
	if core.PodPhase(statusPhase) != core.PodRunning {
		return statusPhase, nil
	}

	// if .status.phase == Running && .status.conditions[type==Ready].status == True, then .status.phase.
	var statusConditions, _, _ = unstructured.NestedSlice(status, "conditions")
	for i := range statusConditions {
		var condition, ok = statusConditions[i].(map[string]interface{})
		if !ok {
			continue
		}
		var conditionType = condition["type"].(string)
		if core.PodConditionType(conditionType) != core.PodReady {
			continue
		}
		var conditionStatus = condition["status"].(string)
		if core.ConditionStatus(conditionStatus) == core.ConditionTrue {
			return statusPhase, nil
		}
		break
	}

	// otherwise, unready.
	return GeneralStatusUnready, nil
}

// getPersistentVolume returns the status of kubernetes persistent volume(claim) resource,
// status is one of "Unknown", "Pending", "Bound", "Lost".
func getPersistentVolume(ctx context.Context, o *unstructured.Unstructured) (string, error) {
	var statusPhase, exist, _ = unstructured.NestedString(o.Object, "status", "phase")
	if !exist {
		return GeneralStatusUnknown, errors.New("not found 'pv(c)' phase")
	}
	return statusPhase, nil
}

// getReplicas returns the status of kubernetes replica set resource,
// status is one of "Unknown", "Unready", "Ready".
func getReplicas(ctx context.Context, dynamicCli *dynamic.DynamicClient, o *unstructured.Unstructured) (string, error) {
	// if getPod(all pod) == Running|Succeeded, then ready.
	var ns, s, err = polymorphic.SelectorsForObject(o)
	if err != nil {
		return GeneralStatusUnknown, fmt.Errorf("error gettting selector of kubernetes %s %s/%s: %w",
			o.GroupVersionKind(), o.GetNamespace(), o.GetName(), err)
	}
	var ss = s.String()
	pl, err := dynamicCli.Resource(core.SchemeGroupVersion.WithResource("pods")).
		Namespace(ns).
		List(ctx, meta.ListOptions{ResourceVersion: "0", LabelSelector: ss})
	if err != nil {
		return GeneralStatusUnknown, fmt.Errorf("error listing kubernetes %s pods with %s: %w",
			ns, ss, err)
	}
	for i := range pl.Items {
		var p = pl.Items[i]
		var ps, err = getPod(ctx, &p)
		if err != nil {
			return GeneralStatusUnknown, fmt.Errorf("error stating kubernetes pod %s/%s: %w",
				p.GetNamespace(), p.GetName(), err)
		}
		switch core.PodPhase(ps) {
		default:
			return GeneralStatusUnready, nil
		case core.PodRunning, core.PodSucceeded:
			// expected pod phase.
		}
	}

	// otherwise, unready.
	return GeneralStatusReady, nil
}

// getAPIService returns the status of kubernetes api service resource,
// status is one of "Unknown", "Unready", "Ready".
func getAPIService(ctx context.Context, o *unstructured.Unstructured) (string, error) {
	// if .status.conditions[type==Available].status == True, then ready.
	var statusConditions, _, _ = unstructured.NestedSlice(o.Object, "status", "conditions")
	for i := range statusConditions {
		var condition, ok = statusConditions[i].(map[string]interface{})
		if !ok {
			continue
		}
		var conditionType = condition["type"].(string)
		if apiregistrationv1.APIServiceConditionType(conditionType) != apiregistrationv1.Available {
			continue
		}
		var conditionStatus = condition["status"].(string)
		if apiregistrationv1.ConditionStatus(conditionStatus) == apiregistrationv1.ConditionTrue {
			return GeneralStatusReady, nil
		}
		break
	}

	// otherwise, unready.
	return GeneralStatusUnready, nil
}

// getDeployment returns the status of kubernetes deployment resource,
// status is one of "Unknown", "Unready", "Ready".
func getDeployment(ctx context.Context, o *unstructured.Unstructured, pauseAsReady bool) (string, error) {
	status, exist, _ := unstructured.NestedMap(o.Object, "status")
	if !exist {
		return GeneralStatusUnknown, errors.New("not found 'deploy' status")
	}

	// if .status.observedGeneration < .metadata.generation, then unready.
	var statusObservedGeneration, _, _ = unstructured.NestedInt64(status, "observedGeneration")
	if statusObservedGeneration < o.GetGeneration() {
		return GeneralStatusUnready, nil
	}

	spec, exist, _ := unstructured.NestedMap(o.Object, "spec")
	if !exist {
		return GeneralStatusUnknown, errors.New("not found 'deploy' spec")
	}

	// paused processing.
	var specPaused, _, _ = unstructured.NestedBool(spec, "paused")
	if specPaused {
		if pauseAsReady {
			return GeneralStatusReady, nil
		}
		return GeneralStatusUnready, nil
	}

	// expected replicas = .spec.replicas - min(.spec.strategy.rollingUpdate.maxUnavailable, .spec.replicas)
	// if .status.readyReplicas < expected replicas, then unready.
	var specReplicas, _, _ = unstructured.NestedInt64(spec, "replicas")
	var specStrategyType, _, _ = unstructured.NestedString(spec, "strategy", "type")
	var expectedReplicas = specReplicas
	if specReplicas > 0 && apps.DeploymentStrategyType(specStrategyType) == apps.RollingUpdateDeploymentStrategyType {
		var maxUnavailable int64
		var specRollingUpdate, exist, _ = unstructured.NestedMap(spec, "strategy", "rollingUpdate")
		if exist {
			var maxUnavailableIntStr = intstr.Parse(fmt.Sprint(specRollingUpdate["maxUnavailable"]))
			var maxUnavailableInt, _ = intstr.GetScaledValueFromIntOrPercent(&maxUnavailableIntStr, int(specReplicas), false)
			maxUnavailable = int64(maxUnavailableInt)
		}
		if maxUnavailable > specReplicas {
			maxUnavailable = specReplicas
		}
		expectedReplicas = specReplicas - maxUnavailable
	}
	var statusReadyReplicas, _, _ = unstructured.NestedInt64(status, "readyReplicas")
	if statusReadyReplicas < expectedReplicas {
		return GeneralStatusUnready, nil
	}

	// otherwise, ready.
	return GeneralStatusReady, nil
}

// getDaemonSet returns the status of kubernetes daemon set resource,
// status is one of "Unknown", "Unready", "Ready".
func getDaemonSet(ctx context.Context, o *unstructured.Unstructured) (string, error) {
	status, exist, _ := unstructured.NestedMap(o.Object, "status")
	if !exist {
		return GeneralStatusUnknown, errors.New("not found 'ds' status")
	}

	// if .status.observedGeneration < .metadata.generation, then unready.
	var statusObservedGeneration, _, _ = unstructured.NestedInt64(status, "observedGeneration")
	if statusObservedGeneration < o.GetGeneration() {
		return GeneralStatusUnready, nil
	}

	spec, exist, _ := unstructured.NestedMap(o.Object, "spec")
	if !exist {
		return GeneralStatusUnknown, errors.New("not found 'ds' spec")
	}

	// if .spec.strategy.type != RollingUpdate, then ready.
	var specStrategyType, _, _ = unstructured.NestedString(spec, "strategy", "type")
	if apps.DaemonSetUpdateStrategyType(specStrategyType) != apps.RollingUpdateDaemonSetStrategyType {
		return GeneralStatusReady, nil
	}

	// if .status.desiredNumberScheduled != .status.updatedNumberScheduled, then unready.
	var statusDesiredNumberScheduled, _, _ = unstructured.NestedInt64(status, "desiredNumberScheduled")
	var statusUpdatedNumberScheduled, _, _ = unstructured.NestedInt64(status, "updatedNumberScheduled")
	if statusDesiredNumberScheduled != statusUpdatedNumberScheduled {
		return GeneralStatusUnready, nil
	}

	// expected replicas = .status.desiredNumberScheduled - min(.spec.strategy.rollingUpdate.maxUnavailable, .status.desiredNumberScheduled)
	// if .status.numberReady < expected replicas, then unready.
	var expectedReplicas = statusDesiredNumberScheduled
	if statusDesiredNumberScheduled > 0 {
		var maxUnavailable int64
		var specRollingUpdate, exist, _ = unstructured.NestedMap(spec, "strategy", "rollingUpdate")
		if exist {
			var maxUnavailableIntStr = intstr.Parse(fmt.Sprint(specRollingUpdate["maxUnavailable"]))
			var maxUnavailableInt, _ = intstr.GetScaledValueFromIntOrPercent(&maxUnavailableIntStr, int(statusDesiredNumberScheduled), true)
			maxUnavailable = int64(maxUnavailableInt)
		}
		if maxUnavailable > statusDesiredNumberScheduled {
			maxUnavailable = statusDesiredNumberScheduled
		}
		expectedReplicas = statusDesiredNumberScheduled - maxUnavailable
	}
	var statusNumberReady, _, _ = unstructured.NestedInt64(status, "numberReady")
	if statusNumberReady < expectedReplicas {
		return GeneralStatusUnready, nil
	}

	// otherwise, ready.
	return GeneralStatusReady, nil
}

// getStatefulSet returns the status of kubernetes stateful set resource,
// status is one of "Unknown", "Unready", "Ready".
func getStatefulSet(ctx context.Context, o *unstructured.Unstructured) (string, error) {
	status, exist, _ := unstructured.NestedMap(o.Object, "status")
	if !exist {
		return GeneralStatusUnknown, errors.New("not found 'sts' status")
	}

	// if .status.observedGeneration < .metadata.generation, then unready.
	var statusObservedGeneration, _, _ = unstructured.NestedInt64(status, "observedGeneration")
	if statusObservedGeneration < o.GetGeneration() {
		return GeneralStatusUnready, nil
	}

	spec, exist, _ := unstructured.NestedMap(o.Object, "spec")
	if !exist {
		return GeneralStatusUnknown, errors.New("not found 'sts' spec")
	}

	// if .status.strategy.type != RollingUpdate, then ready.
	var specStrategyType, _, _ = unstructured.NestedString(spec, "strategy", "type")
	if apps.StatefulSetUpdateStrategyType(specStrategyType) != apps.RollingUpdateStatefulSetStrategyType {
		return GeneralStatusReady, nil
	}

	// expected replicas = .spec.replicas - .spec.strategy.rollingUpdate.partition.
	// if .status.updateReplicas < expected replicas, then unready.
	var specReplicas, _, _ = unstructured.NestedInt64(spec, "replicas")
	var specPartition, _, _ = unstructured.NestedInt64(spec, "strategy", "rollingUpdate", "partition")
	var expectedReplicas = specReplicas - specPartition
	var statusUpdatedReplicas, _, _ = unstructured.NestedInt64(status, "updateReplicas")
	if statusUpdatedReplicas < expectedReplicas {
		return GeneralStatusUnready, nil
	}

	// if .status.readyReplicas != .spec.replicas, then unready.
	var statusReadyReplicas, _, _ = unstructured.NestedInt64(status, "readyReplicas")
	if statusReadyReplicas != specReplicas {
		return GeneralStatusUnready, nil
	}

	// if .status.currentRevision != .status.updateRevision, then unready.
	var statusCurrentRevision, _, _ = unstructured.NestedString(status, "currentRevision")
	var statusUpdateRevision, _, _ = unstructured.NestedString(status, "updateRevision")
	if statusCurrentRevision != statusUpdateRevision {
		return GeneralStatusUnready, nil
	}

	// otherwise, ready.
	return GeneralStatusReady, nil
}

// getJob returns the status of kubernetes job resource,
// status is one of "Unknown", "Unready", "Ready".
func getJob(ctx context.Context, o *unstructured.Unstructured, suspendAsReady bool) (string, error) {
	spec, exist, _ := unstructured.NestedMap(o.Object, "spec")
	if !exist {
		return GeneralStatusUnknown, errors.New("not found 'job' spec")
	}

	// suspend processing.
	var specSuspend, _, _ = unstructured.NestedBool(spec, "suspend")
	if specSuspend {
		if suspendAsReady {
			return GeneralStatusReady, nil
		}
		return GeneralStatusUnready, nil
	}

	status, exist, _ := unstructured.NestedMap(o.Object, "status")
	if !exist {
		return GeneralStatusUnknown, errors.New("not found 'job' status")
	}

	// if .status.failed > .spec.backoffLimit, then unready.
	var statusFailed, _, _ = unstructured.NestedInt64(status, "failed")
	var specBackoffLimit, _, _ = unstructured.NestedInt64(spec, "backoffLimit")
	if statusFailed > specBackoffLimit {
		return GeneralStatusUnready, nil
	}

	// if .status.succeeded < .spec.completions, then unready.
	var statusSucceeded, _, _ = unstructured.NestedInt64(status, "succeeded")
	var specCompletions, _, _ = unstructured.NestedInt64(spec, "completions")
	if statusSucceeded < specCompletions {
		return GeneralStatusUnready, nil
	}

	// otherwise, ready.
	return GeneralStatusReady, nil
}

// getHorizontalPodAutoscaler returns the status of kubernetes hpa resource,
// status is one of "Unknown", "Unready", "Ready".
func getHorizontalPodAutoscaler(ctx context.Context, o *unstructured.Unstructured) (string, error) {
	status, exist, _ := unstructured.NestedMap(o.Object, "status")
	if !exist {
		return GeneralStatusUnknown, errors.New("not found hpa status")
	}

	// if .status.observedGeneration < .metadata.generation, then unready.
	var statusObservedGeneration, _, _ = unstructured.NestedInt64(status, "observedGeneration")
	if exist && statusObservedGeneration < o.GetGeneration() {
		return GeneralStatusUnready, nil
	}

	// if .status.conditions[type==ScalingActive].status == True, then ready.
	var statusConditions, _, _ = unstructured.NestedSlice(status, "conditions")
	for i := range statusConditions {
		var condition, ok = statusConditions[i].(map[string]interface{})
		if !ok {
			continue
		}
		var conditionType = condition["type"].(string)
		if autoscaling.HorizontalPodAutoscalerConditionType(conditionType) != autoscaling.ScalingActive {
			continue
		}
		var conditionStatus = condition["status"].(string)
		if core.ConditionStatus(conditionStatus) == core.ConditionTrue {
			return GeneralStatusReady, nil
		}
		break
	}

	// otherwise, unready.
	return GeneralStatusUnready, nil
}

// getCertificateSigningRequest returns the status of kubernetes csr resource,
// status is one of "Unready", "Ready".
func getCertificateSigningRequest(ctx context.Context, o *unstructured.Unstructured) (string, error) {
	// if len(.status.certificate) != 0, then ready.
	var statusCertificate, exist, _ = unstructured.NestedFieldNoCopy(o.Object, "status", "certificate")
	if !exist || len(statusCertificate.([]byte)) == 0 {
		return GeneralStatusUnready, nil
	}

	// otherwise, unready.
	return GeneralStatusReady, nil
}

// getIngress returns the status of kubernetes ingress resource,
// status is one of "Unready", "Ready".
func getIngress(ctx context.Context, o *unstructured.Unstructured) (string, error) {
	// if len(.status.loadBalancer.ingress) != 0, then ready.
	var statusLBIngresses, _, _ = unstructured.NestedSlice(o.Object, "status", "loadBalancer", "ingress")
	if len(statusLBIngresses) > 0 {
		return GeneralStatusReady, nil
	}

	// otherwise, unready.
	return GeneralStatusUnready, nil
}

// getNetworkPolicy returns the status of kubernetes network policy resource,
// status is one of "Unready", "Ready".
func getNetworkPolicy(ctx context.Context, o *unstructured.Unstructured) (string, error) {
	// if .status.conditions[type==Accepted].status == True, then ready.
	var statusConditions, _, _ = unstructured.NestedSlice(o.Object, "status", "conditions")
	for i := range statusConditions {
		var condition, ok = statusConditions[i].(map[string]interface{})
		if !ok {
			continue
		}
		var conditionType = condition["type"].(string)
		if networking.NetworkPolicyConditionType(conditionType) != networking.NetworkPolicyConditionStatusAccepted {
			continue
		}
		var conditionStatus = condition["status"].(string)
		if core.ConditionStatus(conditionStatus) == core.ConditionTrue {
			return GeneralStatusReady, nil
		}
		break
	}

	// otherwise, unready.
	return GeneralStatusUnready, nil
}

// getPodDisruptionBudget returns the status of kubernetes pdb resource,
// status is one of "Unknown", "Unready", "Ready".
func getPodDisruptionBudget(ctx context.Context, o *unstructured.Unstructured) (string, error) {
	status, exist, _ := unstructured.NestedMap(o.Object, "status")
	if !exist {
		return GeneralStatusUnknown, errors.New("not found 'pdb' status")
	}

	// if .status.observedGeneration < .metadata.generation, then unready.
	var statusObservedGeneration, _, _ = unstructured.NestedInt64(status, "observedGeneration")
	if statusObservedGeneration < o.GetGeneration() {
		return GeneralStatusUnready, nil
	}

	// if .status.conditions[type==DisruptionAllowed].status == True, then ready.
	var statusConditions, _, _ = unstructured.NestedSlice(status, "conditions")
	for i := range statusConditions {
		var condition, ok = statusConditions[i].(map[string]interface{})
		if !ok {
			continue
		}
		var conditionType = condition["type"].(string)
		if conditionType != policy.DisruptionAllowedCondition {
			continue
		}
		var conditionStatus = condition["status"].(string)
		if core.ConditionStatus(conditionStatus) == core.ConditionTrue {
			return GeneralStatusReady, nil
		}
		break
	}

	// otherwise, unready.
	return GeneralStatusUnready, nil
}

// getWebhookConfiguration returns the status of kubernetes webhook configuration resource,
// status is one of "Unknown", "Unready", "Ready".
func getWebhookConfiguration(ctx context.Context, dynamicCli *dynamic.DynamicClient, o *unstructured.Unstructured) (string, error) {
	// if getService(.spec.webhooks[.clientConfig.service?]) == Unready, then unready.
	var specWebhooks, _, _ = unstructured.NestedSlice(o.Object, "spec", "webhooks")
	for i := range specWebhooks {
		var webhook, ok = specWebhooks[i].(map[string]interface{})
		if !ok {
			continue
		}
		var webhookService, exist, _ = unstructured.NestedMap(webhook, "clientConfig", "service")
		if !exist {
			continue
		}
		var svcName = webhookService["name"].(string)
		if svcName == "" {
			continue
		}
		var svcNamespace = webhookService["namespace"].(string)
		var svc, err = dynamicCli.Resource(core.SchemeGroupVersion.WithResource("services")).
			Namespace(svcNamespace).
			Get(ctx, svcName, meta.GetOptions{ResourceVersion: "0"})
		if err != nil {
			if kerrors.IsNotFound(err) {
				return GeneralStatusUnready, nil
			}
			return GeneralStatusUnknown, err
		}
		svcState, err := getService(ctx, svc)
		if err != nil {
			return GeneralStatusUnknown, fmt.Errorf("error stating kubernetes service: %w", err)
		}
		if svcState == GeneralStatusUnready {
			return GeneralStatusUnready, nil
		}
	}

	// otherwise, ready.
	return GeneralStatusReady, nil
}
