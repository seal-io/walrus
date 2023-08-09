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
	"github.com/seal-io/seal/pkg/operator/k8s/polymorphic"
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
	var s string
	if msg == "" {
		s = "Server Error"
	} else {
		s = "Server Error: " + msg
	}

	return &typestatus.Status{
		Summary: typestatus.Summary{
			SummaryStatus:        typestatus.GeneralStatusError,
			Error:                true,
			SummaryStatusMessage: s,
		},
	}
}

// Get returns status of the given k8s resource.
func Get(
	ctx context.Context,
	dynamicCli *dynamic.DynamicClient,
	o *unstructured.Unstructured,
) (*typestatus.Status, error) {
	switch o.GetKind() {
	case "Service":
		return getService(o)
	case "Pod":
		return getPod(o)
	case "PersistentVolume", "PersistentVolumeClaim":
		return getPersistentVolume(o)
	case "APIService":
		return getAPIService(o)
	case "Deployment":
		return getDeployment(o)
	case "DaemonSet":
		return getDaemonSet(o)
	case "StatefulSet":
		return getStatefulSet(o)
	case "ReplicationController", "ReplicaSet":
		return getReplicas(ctx, dynamicCli, o)
	case "Job":
		return getJob(o)
	case "HorizontalPodAutoscaler":
		return getHorizontalPodAutoscaler(o)
	case "CertificateSigningRequest":
		return getCertificateSigningRequest(o)
	case "Ingress":
		return getIngress(o)
	case "NetworkPolicy":
		return getNetworkPolicy(o)
	case "PodDisruptionBudget":
		return getPodDisruptionBudget(o)
	case "ValidatingWebhookConfiguration", "MutatingWebhookConfiguration":
		return getWebhookConfiguration(ctx, dynamicCli, o)
	}

	return &GeneralStatusReady, nil
}

// getService returns the status of kubernetes service resource.
func getService(o *unstructured.Unstructured) (*typestatus.Status, error) {
	spec, exist, _ := unstructured.NestedMap(o.Object, "spec")
	if !exist {
		return nil, errors.New("not found 'service' spec")
	}

	// If .spec.type == ExternalName, then ready.
	specType, _, _ := unstructured.NestedString(spec, "type")
	if core.ServiceType(specType) == core.ServiceTypeExternalName {
		return &GeneralStatusReady, nil
	}

	// If .spec.clusterIP == "", then not ready.
	specClusterIP, _, _ := unstructured.NestedString(spec, "clusterIP")
	if specClusterIP == "" {
		return &GeneralStatusReadyTransitioning, nil
	}

	// If .spec.type == LoadBalancer && len(.spec.externalIPs) > 0, then ready.
	// If .spec.type == LoadBalancer && len(.status.loadBalancer.ingress) > 0, then ready.
	if core.ServiceType(specType) == core.ServiceTypeLoadBalancer {
		specExternalIPs, _, _ := unstructured.NestedStringSlice(spec, "externalIPs")
		if len(specExternalIPs) > 0 {
			return &GeneralStatusReady, nil
		}

		statusLBIngresses, _, _ := unstructured.NestedSlice(o.Object, "status", "loadBalancer", "ingress")
		if len(statusLBIngresses) > 0 {
			return &GeneralStatusReady, nil
		}

		return &GeneralStatusReadyTransitioning, nil
	}

	// Otherwise, ready.
	return &GeneralStatusReady, nil
}

// getPod returns the status of kubernetes pod resource.
func getPod(o *unstructured.Unstructured) (*typestatus.Status, error) {
	status, exist, _ := unstructured.NestedMap(o.Object, "status")
	if !exist {
		return nil, errors.New("not found 'pod' status")
	}

	var (
		statusConds, _, _ = unstructured.NestedSlice(status, "conditions")
		st                = &typestatus.Status{}
	)

	st.SetConditions(toConditions(statusConds))
	st.SetSummary(podStatusPaths.Walk(st))

	// Dig clearer error message from status.
	if st.Error {
		st.SummaryStatusMessage = digPodErrorReason(status)
	}

	return st, nil
}

// getPersistentVolume returns the status of kubernetes persistent volume(claim) resource,.
func getPersistentVolume(o *unstructured.Unstructured) (*typestatus.Status, error) {
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

// getReplicas returns the status of kubernetes replica set resource.
func getReplicas(
	ctx context.Context,
	dynamicCli *dynamic.DynamicClient,
	o *unstructured.Unstructured,
) (*typestatus.Status, error) {
	st := &typestatus.Status{}
	// Use conditions to generate status while it existed.
	statusConditions, exist, _ := unstructured.NestedSlice(o.Object, "status", "conditions")
	if exist {
		st.SetConditions(toConditions(statusConditions))
		st.SetSummary(replicaSetStatusPaths.Walk(st))

		return st, nil
	}

	// If getPod(all pod) == Running|Succeeded, then ready.
	ns, s, err := polymorphic.SelectorsForObject(o)
	if err != nil {
		return nil, fmt.Errorf("error gettting selector of kubernetes %s %s/%s: %w",
			o.GroupVersionKind(), o.GetNamespace(), o.GetName(), err)
	}
	ss := s.String()

	pl, err := dynamicCli.Resource(core.SchemeGroupVersion.WithResource("pods")).
		Namespace(ns).
		List(ctx, meta.ListOptions{ResourceVersion: "0", LabelSelector: ss})
	if err != nil {
		return nil, fmt.Errorf("error listing kubernetes %s pods with %s: %w",
			ns, ss, err)
	}

	for i := range pl.Items {
		p := pl.Items[i]

		ps, err := getPod(&p)
		if err != nil {
			return nil, fmt.Errorf("error stating kubernetes pod %s/%s: %w",
				p.GetNamespace(), p.GetName(), err)
		}

		switch ps.Summary.SummaryStatus {
		default:
			return ps, nil
		case string(core.PodReady):
			// Expected pod phase.
		}
	}

	// Otherwise, not ready.
	return &GeneralStatusReady, nil
}

// getAPIService returns the status of kubernetes api service resource,.
func getAPIService(o *unstructured.Unstructured) (*typestatus.Status, error) {
	statusConditions, exist, _ := unstructured.NestedSlice(o.Object, "status", "conditions")
	if !exist {
		return nil, errors.New("not found 'api service' status conditions")
	}

	st := &typestatus.Status{}
	st.SetConditions(toConditions(statusConditions))
	st.SetSummary(apiServiceStatusPaths.Walk(st))

	return st, nil
}

// getDeployment returns the status of kubernetes deployment resource.
func getDeployment(o *unstructured.Unstructured) (*typestatus.Status, error) {
	statusConditions, exist, _ := unstructured.NestedSlice(o.Object, "status", "conditions")
	if !exist {
		return nil, errors.New("not found 'deployment' status conditions")
	}

	st := &typestatus.Status{}
	st.SetConditions(toConditions(statusConditions))
	st.SetSummary(deploymentStatusPaths.Walk(st))

	return st, nil
}

// getDaemonSet returns the status of kubernetes daemon set resource,
// daemonSet status condition is empty, judge the summary based on other fields.
func getDaemonSet(o *unstructured.Unstructured) (*typestatus.Status, error) {
	status, exist, _ := unstructured.NestedMap(o.Object, "status")
	if !exist {
		return nil, errors.New("not found 'daemonSet' status")
	}

	// If .status.observedGeneration < .metadata.generation, then not ready.
	statusObservedGeneration, _, _ := unstructured.NestedInt64(status, "observedGeneration")
	if statusObservedGeneration < o.GetGeneration() {
		return &GeneralStatusReadyTransitioning, nil
	}

	spec, exist, _ := unstructured.NestedMap(o.Object, "spec")
	if !exist {
		return nil, errors.New("not found 'daemonSet' spec")
	}

	// If .spec.strategy.type != RollingUpdate, then ready.
	specStrategyType, _, _ := unstructured.NestedString(spec, "strategy", "type")
	if apps.DaemonSetUpdateStrategyType(specStrategyType) !=
		apps.RollingUpdateDaemonSetStrategyType {
		return &GeneralStatusReady, nil
	}

	// If .status.desiredNumberScheduled != .status.updatedNumberScheduled, then not ready.
	statusDesiredNumberScheduled, _, _ := unstructured.NestedInt64(status, "desiredNumberScheduled")
	statusUpdatedNumberScheduled, _, _ := unstructured.NestedInt64(status, "updatedNumberScheduled")

	if statusDesiredNumberScheduled != statusUpdatedNumberScheduled {
		return &GeneralStatusReadyTransitioning, nil
	}

	// Expected replicas =
	// .status.desiredNumberScheduled - min(.spec.strategy.rollingUpdate.maxUnavailable, .status.desiredNumberScheduled)
	// if .status.numberReady < expected replicas, then not ready.
	expectedReplicas := statusDesiredNumberScheduled

	if statusDesiredNumberScheduled > 0 {
		var maxUnavailable int64

		specRollingUpdate, exist, _ := unstructured.NestedMap(
			spec,
			"strategy",
			"rollingUpdate",
		)
		if exist {
			maxUnavailableIntStr := intstr.Parse(fmt.Sprint(specRollingUpdate["maxUnavailable"]))
			maxUnavailableInt, _ := intstr.GetScaledValueFromIntOrPercent(
				&maxUnavailableIntStr,
				int(statusDesiredNumberScheduled),
				true,
			)
			maxUnavailable = int64(maxUnavailableInt)
		}

		if maxUnavailable > statusDesiredNumberScheduled {
			maxUnavailable = statusDesiredNumberScheduled
		}
		expectedReplicas = statusDesiredNumberScheduled - maxUnavailable
	}

	statusNumberReady, _, _ := unstructured.NestedInt64(status, "numberReady")
	if statusNumberReady < expectedReplicas {
		return &GeneralStatusReadyTransitioning, nil
	}

	// Otherwise, ready.
	return &GeneralStatusReady, nil
}

// getStatefulSet returns the status of kubernetes stateful set resource,
// daemonSet status condition is empty, judge the summary based on other fields.
func getStatefulSet(o *unstructured.Unstructured) (*typestatus.Status, error) {
	status, exist, _ := unstructured.NestedMap(o.Object, "status")
	if !exist {
		return nil, errors.New("not found 'statefulSet' status")
	}

	// If .status.observedGeneration < .metadata.generation, then not ready.
	statusObservedGeneration, _, _ := unstructured.NestedInt64(status, "observedGeneration")
	if statusObservedGeneration < o.GetGeneration() {
		return &GeneralStatusReadyTransitioning, nil
	}

	spec, exist, _ := unstructured.NestedMap(o.Object, "spec")
	if !exist {
		return nil, errors.New("not found 'statefulSet' spec")
	}

	// If .status.strategy.type != RollingUpdate, then ready.
	specStrategyType, _, _ := unstructured.NestedString(spec, "strategy", "type")
	if apps.StatefulSetUpdateStrategyType(specStrategyType) !=
		apps.RollingUpdateStatefulSetStrategyType {
		return &GeneralStatusReady, nil
	}

	// Expected replicas = .spec.replicas - .spec.strategy.rollingUpdate.partition.
	// If .status.updateReplicas < expected replicas, then not ready.
	specReplicas, _, _ := unstructured.NestedInt64(spec, "replicas")
	specPartition, _, _ := unstructured.NestedInt64(spec, "strategy", "rollingUpdate", "partition")
	expectedReplicas := specReplicas - specPartition
	statusUpdatedReplicas, _, _ := unstructured.NestedInt64(status, "updateReplicas")

	if statusUpdatedReplicas < expectedReplicas {
		return &GeneralStatusReadyTransitioning, nil
	}

	// If .status.readyReplicas != .spec.replicas, then not ready.
	statusReadyReplicas, _, _ := unstructured.NestedInt64(status, "readyReplicas")
	if statusReadyReplicas != specReplicas {
		return &GeneralStatusReadyTransitioning, nil
	}

	// If .status.currentRevision != .status.updateRevision, then not ready.
	statusCurrentRevision, _, _ := unstructured.NestedString(status, "currentRevision")
	statusUpdateRevision, _, _ := unstructured.NestedString(status, "updateRevision")

	if statusCurrentRevision != statusUpdateRevision {
		return &GeneralStatusReadyTransitioning, nil
	}

	// Otherwise, ready.
	return &GeneralStatusReady, nil
}

// getJob returns the status of kubernetes job resource.
func getJob(o *unstructured.Unstructured) (*typestatus.Status, error) {
	statusConditions, exist, _ := unstructured.NestedSlice(o.Object, "status", "conditions")
	if !exist {
		return nil, errors.New("not found 'job' status conditions")
	}

	st := &typestatus.Status{}
	st.SetConditions(toConditions(statusConditions))
	st.SetSummary(jobStatusPaths.Walk(st))

	return st, nil
}

// getHorizontalPodAutoscaler returns the status of kubernetes hpa resource.
func getHorizontalPodAutoscaler(o *unstructured.Unstructured) (*typestatus.Status, error) {
	statusConditions, exist, _ := unstructured.NestedSlice(o.Object, "status", "conditions")
	if !exist {
		return nil, errors.New("not found 'horizontal pod autoscaler' status conditions")
	}

	st := &typestatus.Status{}
	st.SetConditions(toConditions(statusConditions))
	st.SetSummary(hpaStatusPaths.Walk(st))

	return st, nil
}

// getCertificateSigningRequest returns the status of kubernetes csr resource.
func getCertificateSigningRequest(o *unstructured.Unstructured) (*typestatus.Status, error) {
	statusConditions, exist, _ := unstructured.NestedSlice(o.Object, "status", "conditions")
	if !exist {
		return nil, errors.New("not found 'certificate signing request' status conditions")
	}

	st := &typestatus.Status{}
	st.SetConditions(toConditions(statusConditions))
	st.SetSummary(csrStatusPaths.Walk(st))

	return st, nil
}

// getIngress returns the status of kubernetes ingress resource,
// ingress status isn't contain conditions, judge the summary based on other fields.
func getIngress(o *unstructured.Unstructured) (*typestatus.Status, error) {
	// If len(.status.loadBalancer.ingress) != 0, then ready.
	statusLBIngresses, _, _ := unstructured.NestedSlice(o.Object, "status", "loadBalancer", "ingress")
	if len(statusLBIngresses) > 0 {
		return &GeneralStatusReady, nil
	}

	// Otherwise, not ready.
	return &GeneralStatusReadyTransitioning, nil
}

// getNetworkPolicy returns the status of kubernetes network policy resource.
func getNetworkPolicy(o *unstructured.Unstructured) (*typestatus.Status, error) {
	statusConditions, exist, _ := unstructured.NestedSlice(o.Object, "status", "conditions")
	if !exist {
		return nil, errors.New("not found 'network policy' status conditions")
	}

	st := &typestatus.Status{}
	st.SetConditions(toConditions(statusConditions))
	st.SetSummary(networkPolicyStatusPaths.Walk(st))

	return st, nil
}

// getPodDisruptionBudget returns the status of kubernetes pdb resource.
func getPodDisruptionBudget(o *unstructured.Unstructured) (*typestatus.Status, error) {
	statusConditions, exist, _ := unstructured.NestedSlice(o.Object, "status", "conditions")
	if !exist {
		return nil, errors.New("not found 'pod disruption budget' status conditions")
	}

	st := &typestatus.Status{}
	st.SetConditions(toConditions(statusConditions))
	st.SetSummary(pdbStatusPaths.Walk(st))

	return st, nil
}

// getWebhookConfiguration returns the status of kubernetes webhook configuration resource.
func getWebhookConfiguration(
	ctx context.Context,
	dynamicCli *dynamic.DynamicClient,
	o *unstructured.Unstructured,
) (*typestatus.Status, error) {
	// If getService(.spec.webhooks[.clientConfig.service?]) == NotReady, then not ready.
	specWebhooks, _, _ := unstructured.NestedSlice(o.Object, "spec", "webhooks")
	for i := range specWebhooks {
		webhook, ok := specWebhooks[i].(map[string]any)
		if !ok {
			continue
		}

		webhookService, exist, _ := unstructured.NestedMap(webhook, "clientConfig", "service")
		if !exist {
			continue
		}

		svcName := webhookService["name"].(string)
		if svcName == "" {
			continue
		}
		svcNamespace := webhookService["namespace"].(string)

		svc, err := dynamicCli.Resource(core.SchemeGroupVersion.WithResource("services")).
			Namespace(svcNamespace).
			Get(ctx, svcName, meta.GetOptions{ResourceVersion: "0"})
		if err != nil {
			if kerrors.IsNotFound(err) {
				return &GeneralStatusReadyTransitioning, nil
			}

			return nil, err
		}

		svcState, err := getService(svc)
		if err != nil {
			return nil, fmt.Errorf("error stating kubernetes service: %w", err)
		}

		if svcState.SummaryStatus != GeneralStatusReady.SummaryStatus {
			return svcState, nil
		}
	}

	// Otherwise, ready.
	return &GeneralStatusReady, nil
}

func toConditions(statusConds []any) (conds []typestatus.Condition) {
	for i := range statusConds {
		condition, ok := statusConds[i].(map[string]any)
		if !ok {
			continue
		}

		condType, ok, _ := unstructured.NestedString(condition, "type")
		if !ok {
			continue
		}

		condTypeStatus, ok, _ := unstructured.NestedString(condition, "status")
		if !ok {
			continue
		}
		msg, _, _ := unstructured.NestedString(condition, "message")
		reason, _, _ := unstructured.NestedString(condition, "reason")
		cond := typestatus.Condition{
			Type:    typestatus.ConditionType(condType),
			Status:  typestatus.ConditionStatus(condTypeStatus),
			Message: msg,
			Reason:  reason,
		}

		ts, ok, _ := unstructured.NestedString(condition, "lastTransitionTime")
		if ok {
			lastUpdateTime, err := time.Parse(time.RFC3339, ts)
			if err == nil {
				cond.LastUpdateTime = lastUpdateTime
			}
		}

		conds = append(conds, cond)
	}

	return conds
}
