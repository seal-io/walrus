package platformk8s

import (
	"context"
	"fmt"
	"sort"

	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicclient "k8s.io/client-go/dynamic"
	batchclient "k8s.io/client-go/kubernetes/typed/batch/v1"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/utils/pointer"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/platform/operator"
	"github.com/seal-io/seal/pkg/platformk8s/intercept"
	"github.com/seal-io/seal/pkg/platformk8s/kube"
	"github.com/seal-io/seal/pkg/platformk8s/polymorphic"
)

// GetKeys implements operator.Operator.
func (op Operator) GetKeys(ctx context.Context, res *model.ApplicationResource) (*operator.Keys, error) {
	var psp, err = op.getPods(ctx, res)
	if err != nil {
		return nil, err
	}

	// {
	//      "labels": ["Pod", "Container"],
	//      "keys":   [
	//          {
	//              "name": "<pod name>",
	//              "keys": [
	//                  {
	//                      "name":  "<container name>",
	//                      "value": "<key>"
	//                  }
	//              ]
	//          }
	//      ]
	// }
	var ks = operator.Keys{
		Labels: []string{"Pod", "Container"},
		Keys:   make([]operator.Key, 0),
	}
	if psp == nil {
		return &ks, nil
	}
	var ps = *psp
	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].CreationTimestamp.Time.After(ps[j].CreationTimestamp.Time)
	})
	for i := 0; i < len(ps); i++ {
		var running = kube.IsPodRunning(&ps[i])
		var states = kube.GetContainerStates(&ps[i])

		var k = operator.Key{
			Name: ps[i].Name, // pod name
			Keys: make([]operator.Key, 0, len(states)),
		}
		for j := 0; j < len(states); j++ {
			k.Keys = append(k.Keys, operator.Key{
				Name:       states[j].Name,     // container name
				Value:      states[j].String(), // key
				Loggable:   pointer.Bool(states[j].State > kube.ContainerStateUnknown),
				Executable: pointer.Bool(running && states[j].State == kube.ContainerStateRunning),
			})
		}
		ks.Keys = append(ks.Keys, k)
	}
	return &ks, nil
}

func (op Operator) getPods(ctx context.Context, res *model.ApplicationResource) (*[]core.Pod, error) {
	if res == nil {
		return nil, nil
	}

	// parse operable resources.
	var rs, err = parseResources(ctx, op, res, intercept.Operable())
	if err != nil {
		if !isResourceParsingError(err) {
			return nil, err
		}
		// warn out if got above errors.
		op.Logger.Warn(err)
		return nil, nil
	}

	// get pods of resources.
	var ps []core.Pod
	for _, r := range rs {
		switch r.Resource {
		case "pods":
			var p, err = op.getPod(ctx, r.Namespace, r.Name)
			if err != nil {
				return nil, err
			}
			if p == nil {
				continue
			}
			ps = append(ps, *p)
		case "cronjobs":
			var psp, err = op.getPodsOfCronJob(ctx, r.Namespace, r.Name)
			if err != nil {
				return nil, err
			}
			if psp == nil {
				continue
			}
			ps = append(ps, *psp...)
		default:
			var psp, err = op.getPodsOfAny(ctx, r.GroupVersionResource, r.Namespace, r.Name)
			if err != nil {
				return nil, err
			}
			if psp == nil {
				continue
			}
			ps = append(ps, *psp...)
		}
	}
	if len(ps) == 0 {
		return nil, nil
	}
	return &ps, nil
}

func (op Operator) getPod(ctx context.Context, ns, n string) (*core.Pod, error) {
	// fetch pod with name.
	coreCli, err := coreclient.NewForConfig(op.RestConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes core client: %w", err)
	}
	p, err := coreCli.Pods(ns).
		Get(ctx, n, meta.GetOptions{ResourceVersion: "0"}) // non quorum read
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return nil, fmt.Errorf("error getting kubernetes %s/%s pod: %w",
				ns, n, err)
		}
		return nil, nil
	}
	return p, nil
}

func (op Operator) getPodsOfCronJob(ctx context.Context, ns, n string) (*[]core.Pod, error) {
	// fetch controlled cronjob with name.
	batchCli, err := batchclient.NewForConfig(op.RestConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kuberentes batch client: %w", err)
	}
	cj, err := batchCli.CronJobs(ns).
		Get(ctx, n, meta.GetOptions{ResourceVersion: "0"}) // non quorum read
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return nil, fmt.Errorf("error getting kubernetes %s/%s cronjob: %w",
				ns, n, err)
		}
		return nil, nil
	}

	// fetch jobs in pagination and filter out.
	var js []*batch.Job
	var jlo = meta.ListOptions{Limit: 100}
	for {
		var jl *batch.JobList
		jl, err = batchCli.Jobs(ns).List(ctx, jlo)
		if err != nil {
			return nil, fmt.Errorf("error listing kubernetes %s jobs: %w",
				ns, err)
		}
		for i := 0; i < len(jl.Items); i++ {
			if !meta.IsControlledBy(&jl.Items[i], cj) {
				continue
			}
			js = append(js, &jl.Items[i])
		}
		jlo.Continue = jl.Continue
		if jlo.Continue == "" {
			break
		}
	}
	if len(js) == 0 {
		return nil, nil
	}

	// fetch pods with job label.
	var ps []core.Pod
	coreCli, err := coreclient.NewForConfig(op.RestConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes core client: %w", err)
	}
	for i := 0; i < len(js); i++ {
		var ss = labels.SelectorFromSet(labels.Set{
			"controller-uid": string(js[i].UID),
			"job-name":       js[i].Name,
		}).String()
		var pl *core.PodList
		pl, err = coreCli.Pods(ns).
			List(ctx, meta.ListOptions{ResourceVersion: "0", LabelSelector: ss}) // non quorum read
		if err != nil {
			if !kerrors.IsNotFound(err) {
				return nil, fmt.Errorf("error listing kubernetes %s pods with %s: %w",
					ns, ss, err)
			}
			continue
		}
		for j := 0; j < len(pl.Items); j++ {
			ps = append(ps, pl.Items[j])
		}
	}
	if len(ps) == 0 {
		return nil, nil
	}
	return &ps, nil
}

func (op Operator) getPodsOfAny(ctx context.Context, gvr schema.GroupVersionResource, ns, n string) (*[]core.Pod, error) {
	// fetch label selector with dynamic client.
	dynamicCli, err := dynamicclient.NewForConfig(op.RestConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes dynamic client: %w", err)
	}
	o, err := dynamicCli.Resource(gvr).Namespace(ns).
		Get(ctx, n, meta.GetOptions{ResourceVersion: "0"}) // non quorum read
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return nil, fmt.Errorf("error getting kubernetes %s %s/%s: %w",
				gvr.Resource, ns, n, err)
		}
		return nil, nil
	}
	_, s, err := polymorphic.SelectorsForObject(o)
	if err != nil {
		return nil, fmt.Errorf("error gettting selector of kubernetes %s %s/%s: %w",
			gvr.Resource, ns, n, err)
	}

	// fetch pods with label selector.
	var ss = s.String()
	coreCli, err := coreclient.NewForConfig(op.RestConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes core client: %w", err)
	}
	pl, err := coreCli.Pods(ns).
		List(ctx, meta.ListOptions{ResourceVersion: "0", LabelSelector: ss}) // non quorum read
	if err != nil {
		return nil, fmt.Errorf("error listing kubernetes %s pods with %s: %w",
			ns, ss, err)
	}
	return &pl.Items, nil
}
