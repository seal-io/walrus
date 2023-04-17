package platformk8s

import (
	"context"
	"fmt"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformk8s/intercept"
	"github.com/seal-io/seal/pkg/platformk8s/kube"
)

// GetComponents implements operator.Operator.
func (op Operator) GetComponents(ctx context.Context, res *model.ApplicationResource) ([]*model.ApplicationResource, error) {
	if res == nil {
		return nil, nil
	}

	// parse composite resources.
	var rs, err = parseResources(ctx, op, res, intercept.Composite())
	if err != nil {
		if !isResourceParsingError(err) {
			return nil, err
		}
		// warn out if got above errors.
		op.Logger.Warn(err)
		return nil, nil
	}

	// get components of resources.
	var comps []*model.ApplicationResource
	for _, r := range rs {
		switch r.Resource {
		case "persistentvolumeclaims":
			var component, err = op.getComponentOfPersistentVolumeClaim(ctx, r.Namespace, r.Name)
			if err != nil {
				return nil, err
			}
			if component == nil {
				continue
			}
			comps = append(comps, component)
		case "cronjobs":
			var components, err = op.getComponentsOfCronJob(ctx, r.Namespace, r.Name)
			if err != nil {
				return nil, err
			}
			comps = append(comps, components...)
		default:
			var components, err = op.getComponentsOfAny(ctx, r.GroupVersionResource, r.Namespace, r.Name)
			if err != nil {
				return nil, err
			}
			comps = append(comps, components...)
		}
	}
	for i := range comps {
		// copy required field from composition resource.
		comps[i].CompositionID = res.ID
		comps[i].InstanceID = res.InstanceID
		comps[i].ConnectorID = res.ConnectorID
		comps[i].Module = res.Module
		comps[i].Mode = types.ApplicationResourceModeDiscovered
		comps[i].DeployerType = res.DeployerType
	}
	return comps, nil
}

func (op Operator) getComponentOfPersistentVolumeClaim(ctx context.Context, ns, n string) (*model.ApplicationResource, error) {
	// fetch controlled persistent volume claim.
	var coreCli, err = coreclient.NewForConfig(op.RestConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes core client: %w", err)
	}
	pvc, err := coreCli.PersistentVolumeClaims(ns).
		Get(ctx, n, meta.GetOptions{ResourceVersion: "0"}) // non quorum read
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return nil, fmt.Errorf("error getting kubernetes %s/%s persistent volume claim: %w",
				ns, n, err)
		}
		return nil, nil
	}

	// get persistent volume.
	if pvc.Spec.VolumeName == "" {
		return nil, nil
	}
	return &model.ApplicationResource{
		Type: "kubernetes_persistent_volume_v1",
		Name: kube.NamespacedName("", pvc.Spec.VolumeName),
	}, nil
}

func (op Operator) getComponentsOfCronJob(ctx context.Context, ns, n string) ([]*model.ApplicationResource, error) {
	var psp, err = op.getPodsOfCronJob(ctx, ns, n)
	if err != nil {
		return nil, err
	}
	if psp == nil {
		return nil, nil
	}

	// convert pod to application resource.
	var ps = *psp
	var rs []*model.ApplicationResource
	for i := 0; i < len(ps); i++ {
		rs = append(rs, &model.ApplicationResource{
			Type: "kubernetes_pod_v1",
			Name: kube.NamespacedName(ps[i].Namespace, ps[i].Name),
		})
	}
	return rs, nil
}

func (op Operator) getComponentsOfAny(ctx context.Context, gvr schema.GroupVersionResource, ns, n string) ([]*model.ApplicationResource, error) {
	var psp, err = op.getPodsOfAny(ctx, gvr, ns, n)
	if err != nil {
		return nil, err
	}
	if psp == nil {
		return nil, nil
	}

	// convert pod to application resource.
	var ps = *psp
	var rs []*model.ApplicationResource
	for i := 0; i < len(ps); i++ {
		rs = append(rs, &model.ApplicationResource{
			Type: "kubernetes_pod_v1",
			Name: kube.NamespacedName(ps[i].Namespace, ps[i].Name),
		})
	}
	return rs, nil
}
