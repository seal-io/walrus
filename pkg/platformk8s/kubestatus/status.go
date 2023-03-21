package kubestatus

import (
	"context"
	"errors"
	"fmt"
	"time"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"

	typestatus "github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platformk8s/polymorphic"
)

var (
	GeneralStatusReady              typestatus.Status
	GeneralStatusReadyTransitioning typestatus.Status
)

func init() {
	GeneralStatusReady.SetSummary(
		&typestatus.Summary{
			SummaryStatus: string(typestatus.GeneralStatusReady),
		},
	)

	GeneralStatusReadyTransitioning.SetSummary(
		&typestatus.Summary{
			SummaryStatus: typestatus.GeneralStatusReadyTransitioning,
			Transitioning: true,
		},
	)
}

// StatusError with error message to generate the status error.
func StatusError(msg string) *typestatus.Status {
	st := &typestatus.Status{}
	st.SetSummary(
		&typestatus.Summary{
			SummaryStatus:        typestatus.GeneralStatusError,
			Error:                true,
			SummaryStatusMessage: fmt.Sprintf("Server Error: %s", msg),
		},
	)
	return st
}

type GetOptions struct {
	WithJobs     bool
	IgnorePaused bool
}

// Get returns status of the given k8s resource.
func Get(ctx context.Context, dynamicCli *dynamic.DynamicClient, o *unstructured.Unstructured, opts GetOptions) (*typestatus.Status, error) {
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
			return getJob(ctx, o)
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

	return &GeneralStatusReady, nil
}

// getService returns the status of kubernetes service resource.
func getService(ctx context.Context, o *unstructured.Unstructured) (*typestatus.Status, error) {
	spec, exist, _ := unstructured.NestedMap(o.Object, "spec")
	if !exist {
		return nil, errors.New("not found 'service' spec")
	}

	// if .spec.type == ExternalName, then ready.
	var specType, _, _ = unstructured.NestedString(spec, "type")
	if core.ServiceType(specType) == core.ServiceTypeExternalName {
		return &GeneralStatusReady, nil
	}

	// if .spec.clusterIP == "", then unready.
	var specClusterIP, _, _ = unstructured.NestedString(spec, "clusterIP")
	if specClusterIP == "" {
		return &GeneralStatusReadyTransitioning, nil
	}

	// if .spec.type == LoadBalancer && len(.spec.externalIPs) > 0, then ready.
	// if .spec.type == LoadBalancer && len(.status.loadBalancer.ingress) > 0, then ready.
	if core.ServiceType(specType) == core.ServiceTypeLoadBalancer {
		var specExternalIPs, _, _ = unstructured.NestedStringSlice(spec, "externalIPs")
		if len(specExternalIPs) > 0 {
			return &GeneralStatusReady, nil
		}
		var statusLBIngresses, _, _ = unstructured.NestedSlice(o.Object, "status", "loadBalancer", "ingress")
		if len(statusLBIngresses) > 0 {
			return &GeneralStatusReady, nil
		}
		return &GeneralStatusReadyTransitioning, nil
	}

	// otherwise, ready.
	return &GeneralStatusReady, nil
}

// getPod returns the status of kubernetes pod resource.
func getPod(ctx context.Context, o *unstructured.Unstructured) (*typestatus.Status, error) {
	status, exist, _ := unstructured.NestedMap(o.Object, "status")
	if !exist {
		return nil, errors.New("not found 'pod' status")
	}

	var (
		statusConds, _, _ = unstructured.NestedSlice(status, "conditions")
		st                = &typestatus.Status{}
	)
	st.SetConditions(toConditions(statusConds))
	st.SetSummary(podSummarizer().Summarize(st))
	return st, nil
}

// getPersistentVolume returns the status of kubernetes persistent volume(claim) resource,
func getPersistentVolume(ctx context.Context, o *unstructured.Unstructured) (*typestatus.Status, error) {
	status, exist, _ := unstructured.NestedMap(o.Object, "status")
	if !exist {
		return nil, errors.New("not found 'persistent volume(claim)' status")
	}

	var (
		isErr, isTransition bool
		phase, _, _         = unstructured.NestedString(status, "phase")
		message, _, _       = unstructured.NestedString(status, "message")
	)

	switch phase {
	case string(core.VolumePending):
		isTransition = true
	case string(core.VolumeFailed):
		isErr = true
	}
	return &typestatus.Status{
		Summary: typestatus.Summary{
			SummaryStatus:        phase,
			SummaryStatusMessage: message,
			Error:                isErr,
			Transitioning:        isTransition,
		},
	}, nil
}

// getReplicas returns the status of kubernetes replica set resource
func getReplicas(ctx context.Context, dynamicCli *dynamic.DynamicClient, o *unstructured.Unstructured) (*typestatus.Status, error) {
	st := &typestatus.Status{}
	// use conditions to generate status while it existed
	statusConditions, exist, _ := unstructured.NestedSlice(o.Object, "status", "conditions")
	if exist {
		st.SetConditions(toConditions(statusConditions))
		st.SetSummary(replicaSetSummarizer().Summarize(st))
		return st, nil
	}

	// if getPod(all pod) == Running|Succeeded, then ready.
	var ns, s, err = polymorphic.SelectorsForObject(o)
	if err != nil {
		return nil, fmt.Errorf("error gettting selector of kubernetes %s %s/%s: %w",
			o.GroupVersionKind(), o.GetNamespace(), o.GetName(), err)
	}
	var ss = s.String()
	pl, err := dynamicCli.Resource(core.SchemeGroupVersion.WithResource("pods")).
		Namespace(ns).
		List(ctx, meta.ListOptions{ResourceVersion: "0", LabelSelector: ss})
	if err != nil {
		return nil, fmt.Errorf("error listing kubernetes %s pods with %s: %w",
			ns, ss, err)
	}
	for i := range pl.Items {
		var p = pl.Items[i]
		var ps, err = getPod(ctx, &p)
		if err != nil {
			return nil, fmt.Errorf("error stating kubernetes pod %s/%s: %w",
				p.GetNamespace(), p.GetName(), err)
		}
		switch ps.Summary.SummaryStatus {
		default:
			return ps, nil
		case string(core.PodReady):
			// expected pod phase.
		}
	}

	// otherwise, unready.
	return &GeneralStatusReady, nil
}

// getAPIService returns the status of kubernetes api service resource,
func getAPIService(ctx context.Context, o *unstructured.Unstructured) (*typestatus.Status, error) {
	statusConditions, exist, _ := unstructured.NestedSlice(o.Object, "status", "conditions")
	if !exist {
		return nil, errors.New("not found 'api service' status conditions")
	}

	st := &typestatus.Status{}
	st.SetConditions(toConditions(statusConditions))
	st.SetSummary(apiServiceSummarizer().Summarize(st))
	return st, nil
}

// getDeployment returns the status of kubernetes deployment resource.
func getDeployment(ctx context.Context, o *unstructured.Unstructured, pauseAsReady bool) (*typestatus.Status, error) {
	// paused processing.
	spec, exist, _ := unstructured.NestedMap(o.Object, "spec")
	if !exist {
		return nil, errors.New("not found 'deployment' spec")
	}
	var specPaused, _, _ = unstructured.NestedBool(spec, "paused")
	if specPaused {
		if pauseAsReady {
			return &GeneralStatusReady, nil
		}
		return &GeneralStatusReadyTransitioning, nil
	}

	statusConditions, exist, _ := unstructured.NestedSlice(o.Object, "status", "conditions")
	if !exist {
		return nil, errors.New("not found 'deployment' status conditions")
	}

	st := &typestatus.Status{}
	st.SetConditions(toConditions(statusConditions))
	st.SetSummary(deploymentSummarizer().Summarize(st))
	return st, nil
}

// getDaemonSet returns the status of kubernetes daemon set resource,
// daemonSet status condition is empty, judge the summary based on other fields
func getDaemonSet(ctx context.Context, o *unstructured.Unstructured) (*typestatus.Status, error) {
	status, exist, _ := unstructured.NestedMap(o.Object, "status")
	if !exist {
		return nil, errors.New("not found 'daemonSet' status")
	}

	// if .status.observedGeneration < .metadata.generation, then unready.
	var statusObservedGeneration, _, _ = unstructured.NestedInt64(status, "observedGeneration")
	if statusObservedGeneration < o.GetGeneration() {
		return &GeneralStatusReadyTransitioning, nil
	}

	spec, exist, _ := unstructured.NestedMap(o.Object, "spec")
	if !exist {
		return nil, errors.New("not found 'daemonSet' spec")
	}

	// if .spec.strategy.type != RollingUpdate, then ready.
	var specStrategyType, _, _ = unstructured.NestedString(spec, "strategy", "type")
	if apps.DaemonSetUpdateStrategyType(specStrategyType) != apps.RollingUpdateDaemonSetStrategyType {
		return &GeneralStatusReady, nil
	}

	// if .status.desiredNumberScheduled != .status.updatedNumberScheduled, then unready.
	var statusDesiredNumberScheduled, _, _ = unstructured.NestedInt64(status, "desiredNumberScheduled")
	var statusUpdatedNumberScheduled, _, _ = unstructured.NestedInt64(status, "updatedNumberScheduled")
	if statusDesiredNumberScheduled != statusUpdatedNumberScheduled {
		return &GeneralStatusReadyTransitioning, nil
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
		return &GeneralStatusReadyTransitioning, nil
	}

	// otherwise, ready.
	return &GeneralStatusReady, nil
}

// getStatefulSet returns the status of kubernetes stateful set resource,
// daemonSet status condition is empty, judge the summary based on other fields
func getStatefulSet(ctx context.Context, o *unstructured.Unstructured) (*typestatus.Status, error) {
	status, exist, _ := unstructured.NestedMap(o.Object, "status")
	if !exist {
		return nil, errors.New("not found 'statefulSet' status")
	}

	// if .status.observedGeneration < .metadata.generation, then unready.
	var statusObservedGeneration, _, _ = unstructured.NestedInt64(status, "observedGeneration")
	if statusObservedGeneration < o.GetGeneration() {
		return &GeneralStatusReadyTransitioning, nil
	}

	spec, exist, _ := unstructured.NestedMap(o.Object, "spec")
	if !exist {
		return nil, errors.New("not found 'statefulSet' spec")
	}

	// if .status.strategy.type != RollingUpdate, then ready.
	var specStrategyType, _, _ = unstructured.NestedString(spec, "strategy", "type")
	if apps.StatefulSetUpdateStrategyType(specStrategyType) != apps.RollingUpdateStatefulSetStrategyType {
		return &GeneralStatusReady, nil
	}

	// expected replicas = .spec.replicas - .spec.strategy.rollingUpdate.partition.
	// if .status.updateReplicas < expected replicas, then unready.
	var specReplicas, _, _ = unstructured.NestedInt64(spec, "replicas")
	var specPartition, _, _ = unstructured.NestedInt64(spec, "strategy", "rollingUpdate", "partition")
	var expectedReplicas = specReplicas - specPartition
	var statusUpdatedReplicas, _, _ = unstructured.NestedInt64(status, "updateReplicas")
	if statusUpdatedReplicas < expectedReplicas {
		return &GeneralStatusReadyTransitioning, nil
	}

	// if .status.readyReplicas != .spec.replicas, then unready.
	var statusReadyReplicas, _, _ = unstructured.NestedInt64(status, "readyReplicas")
	if statusReadyReplicas != specReplicas {
		return &GeneralStatusReadyTransitioning, nil
	}

	// if .status.currentRevision != .status.updateRevision, then unready.
	var statusCurrentRevision, _, _ = unstructured.NestedString(status, "currentRevision")
	var statusUpdateRevision, _, _ = unstructured.NestedString(status, "updateRevision")
	if statusCurrentRevision != statusUpdateRevision {
		return &GeneralStatusReadyTransitioning, nil
	}

	// otherwise, ready.
	return &GeneralStatusReady, nil
}

// getJob returns the status of kubernetes job resource.
func getJob(ctx context.Context, o *unstructured.Unstructured) (*typestatus.Status, error) {
	statusConditions, exist, _ := unstructured.NestedSlice(o.Object, "status", "conditions")
	if !exist {
		return nil, errors.New("not found 'job' status conditions")
	}

	st := &typestatus.Status{}
	st.SetConditions(toConditions(statusConditions))
	st.SetSummary(jobSummarizer().Summarize(st))
	return st, nil
}

// getHorizontalPodAutoscaler returns the status of kubernetes hpa resource.
func getHorizontalPodAutoscaler(ctx context.Context, o *unstructured.Unstructured) (*typestatus.Status, error) {
	statusConditions, exist, _ := unstructured.NestedSlice(o.Object, "status", "conditions")
	if !exist {
		return nil, errors.New("not found 'horizontal pod autoscaler' status conditions")
	}

	st := &typestatus.Status{}
	st.SetConditions(toConditions(statusConditions))
	st.SetSummary(hpaSummarizer().Summarize(st))
	return st, nil
}

// getCertificateSigningRequest returns the status of kubernetes csr resource.
func getCertificateSigningRequest(ctx context.Context, o *unstructured.Unstructured) (*typestatus.Status, error) {
	statusConditions, exist, _ := unstructured.NestedSlice(o.Object, "status", "conditions")
	if !exist {
		return nil, errors.New("not found 'certificate signing request' status conditions")
	}

	st := &typestatus.Status{}
	st.SetConditions(toConditions(statusConditions))
	st.SetSummary(csrSummarizer().Summarize(st))
	return st, nil
}

// getIngress returns the status of kubernetes ingress resource,
// ingress status isn't contain conditions, judge the summary based on other fields.
func getIngress(ctx context.Context, o *unstructured.Unstructured) (*typestatus.Status, error) {
	// if len(.status.loadBalancer.ingress) != 0, then ready.
	var statusLBIngresses, _, _ = unstructured.NestedSlice(o.Object, "status", "loadBalancer", "ingress")
	if len(statusLBIngresses) > 0 {
		return &GeneralStatusReady, nil
	}

	// otherwise, unready.
	return &GeneralStatusReadyTransitioning, nil
}

// getNetworkPolicy returns the status of kubernetes network policy resource.
func getNetworkPolicy(ctx context.Context, o *unstructured.Unstructured) (*typestatus.Status, error) {
	statusConditions, exist, _ := unstructured.NestedSlice(o.Object, "status", "conditions")
	if !exist {
		return nil, errors.New("not found 'network policy' status conditions")
	}

	st := &typestatus.Status{}
	st.SetConditions(toConditions(statusConditions))
	st.SetSummary(networkPolicySummarizer().Summarize(st))
	return st, nil
}

// getPodDisruptionBudget returns the status of kubernetes pdb resource.
func getPodDisruptionBudget(ctx context.Context, o *unstructured.Unstructured) (*typestatus.Status, error) {
	statusConditions, exist, _ := unstructured.NestedSlice(o.Object, "status", "conditions")
	if !exist {
		return nil, errors.New("not found 'pod disruption budget' status conditions")
	}

	st := &typestatus.Status{}
	st.SetConditions(toConditions(statusConditions))
	st.SetSummary(pdbSummarizer().Summarize(st))
	return st, nil
}

// getWebhookConfiguration returns the status of kubernetes webhook configuration resource.
func getWebhookConfiguration(ctx context.Context, dynamicCli *dynamic.DynamicClient, o *unstructured.Unstructured) (*typestatus.Status, error) {
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
				return &GeneralStatusReadyTransitioning, nil
			}
			return nil, err
		}
		svcState, err := getService(ctx, svc)
		if err != nil {
			return nil, fmt.Errorf("error stating kubernetes service: %w", err)
		}
		if svcState.SummaryStatus != GeneralStatusReady.SummaryStatus {
			return svcState, nil
		}
	}

	// otherwise, ready.
	return &GeneralStatusReady, nil
}

func toConditions(statusConds []interface{}) []typestatus.Condition {
	var conds []typestatus.Condition
	for i := range statusConds {
		var condition, ok = statusConds[i].(map[string]interface{})
		if !ok {
			continue
		}

		condType, ok := condition["type"].(string)
		if !ok {
			continue
		}
		condTypeStatus, ok := condition["status"].(string)
		if !ok {
			continue
		}
		msg, _ := condition["message"].(string)
		reason, _ := condition["reason"].(string)
		cond := typestatus.Condition{
			Type:    typestatus.ConditionType(condType),
			Status:  typestatus.ConditionStatus(condTypeStatus),
			Message: msg,
			Reason:  reason,
		}

		ts := condition["lastTransitionTime"]
		if ts != nil {
			if lastTransitionTime, ok := ts.(string); ok {
				lastUpdateTime, err := time.Parse(time.RFC3339, lastTransitionTime)
				if err == nil {
					cond.LastUpdateTime = lastUpdateTime
				}
			}
		}

		conds = append(conds, cond)
	}
	return conds
}
