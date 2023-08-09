package kubelabel

import (
	"context"
	"fmt"
	"reflect"

	"go.uber.org/multierr"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	typesk8s "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"

	"github.com/seal-io/seal/pkg/operator/k8s/polymorphic"
	"github.com/seal-io/seal/utils/json"
)

// Apply applies the labels to kubernetes resource.
func Apply(
	ctx context.Context,
	dynamicCli *dynamic.DynamicClient,
	o *unstructured.Unstructured,
	labels map[string]string,
) error {
	p := patcher{
		dynamicCli: dynamicCli,
	}

	switch o.GetKind() {
	case
		"Service",
		"Ingress",
		"Pod",
		"PersistentVolume", "PersistentVolumeClaim":
		return p.applyLabels(ctx, o, labels)
	case
		"Deployment",
		"DaemonSet",
		"StatefulSet",
		"Job",
		"ReplicationController", "ReplicaSet":
		return p.applyLabelsToWorkloads(ctx, o, labels)
	}

	return nil
}

type patcher struct {
	dynamicCli *dynamic.DynamicClient
}

func (p *patcher) applyLabelsToWorkloads(
	ctx context.Context,
	o *unstructured.Unstructured,
	labels map[string]string,
) error {
	pods, err := p.selectPods(ctx, o)
	if err != nil {
		return err
	}

	var berr error

	for i := range pods {
		err = p.applyLabels(ctx, &pods[i], labels)
		multierr.AppendInto(&berr, err)
	}

	return berr
}

func (p *patcher) applyLabels(ctx context.Context, o *unstructured.Unstructured, labels map[string]string) error {
	var (
		ns     = o.GetNamespace()
		name   = o.GetName()
		gvk    = o.GetObjectKind().GroupVersionKind()
		gvr, _ = meta.UnsafeGuessKindToResource(gvk)
	)

	metadata, err := meta.Accessor(o)
	if err != nil {
		return fmt.Errorf("error get metadata for %s %s/%s: %w", gvr.Resource, ns, name, err)
	}

	origin := metadata.GetLabels()

	update := metadata.GetLabels()
	if update == nil {
		update = make(map[string]string, len(labels))
	}

	for k, v := range labels {
		update[k] = v
	}

	// Unchanged.
	if reflect.DeepEqual(origin, update) {
		return nil
	}

	// Change.
	patchBytes, err := json.Marshal(map[string]any{
		"metadata": map[string]any{
			"labels": update,
		},
	})
	if err != nil {
		return fmt.Errorf("error create labels patch: %w", err)
	}

	_, err = p.dynamicCli.Resource(gvr).Namespace(ns).Patch(
		ctx,
		name,
		typesk8s.StrategicMergePatchType,
		patchBytes,
		metav1.PatchOptions{},
	)
	if err != nil {
		return fmt.Errorf("error patch labels to %s %s/%s: %w", gvr.Resource, ns, name, err)
	}

	return nil
}

func (p *patcher) selectPods(ctx context.Context, o *unstructured.Unstructured) ([]unstructured.Unstructured, error) {
	var (
		name   = o.GetName()
		gvk    = o.GetObjectKind().GroupVersionKind()
		gvr, _ = meta.UnsafeGuessKindToResource(gvk)
	)

	ns, s, err := polymorphic.SelectorsForObject(o)
	if err != nil {
		return nil, fmt.Errorf("error gettting selector of kubernetes %s %s/%s: %w", gvr.Resource, ns, name, err)
	}

	ss := s.String()

	pl, err := p.dynamicCli.
		Resource(core.SchemeGroupVersion.WithResource("pods")).
		Namespace(ns).
		List(ctx, metav1.ListOptions{ResourceVersion: "0", LabelSelector: ss})
	if err != nil {
		return nil, fmt.Errorf("error listing kubernetes %s pods with %s: %w", ns, ss, err)
	}

	return pl.Items, nil
}
