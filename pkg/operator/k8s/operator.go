package k8s

import (
	"context"
	"fmt"
	"time"

	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicclient "k8s.io/client-go/dynamic"
	batchclient "k8s.io/client-go/kubernetes/typed/batch/v1"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/k8s"
	"github.com/seal-io/seal/pkg/operator/k8s/polymorphic"
	optypes "github.com/seal-io/seal/pkg/operator/types"
	"github.com/seal-io/seal/utils/log"
)

const OperatorType = types.ConnectorTypeK8s

// New returns operator.Operator with the given options.
func New(ctx context.Context, opts optypes.CreateOptions) (optypes.Operator, error) {
	// NB(thxCode): disable timeout as we don't know the maximum time-cost of once operation,
	// and rely on the session context timeout control of each operation.
	restConfig, err := GetConfig(opts.Connector, WithoutTimeout())
	if err != nil {
		return nil, err
	}
	op := Operator{
		Logger:     log.WithName("operator").WithName("k8s"),
		RestConfig: restConfig,
	}

	return op, nil
}

type Operator struct {
	Logger     log.Logger
	RestConfig *rest.Config
}

// Type implements operator.Operator.
func (Operator) Type() optypes.Type {
	return OperatorType
}

// IsConnected implements operator.Operator.
func (op Operator) IsConnected(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return k8s.Wait(ctx, op.RestConfig)
}

func (op Operator) getPod(ctx context.Context, ns, n string) (*core.Pod, error) {
	// Fetch pod with name.
	coreCli, err := coreclient.NewForConfig(op.RestConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes core client: %w", err)
	}

	p, err := coreCli.Pods(ns).
		Get(ctx, n, meta.GetOptions{ResourceVersion: "0"}) // Non quorum read.
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
	// Fetch controlled cronjob with name.
	batchCli, err := batchclient.NewForConfig(op.RestConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kuberentes batch client: %w", err)
	}

	cj, err := batchCli.CronJobs(ns).
		Get(ctx, n, meta.GetOptions{ResourceVersion: "0"}) // Non quorum read.
	if err != nil {
		if !kerrors.IsNotFound(err) {
			return nil, fmt.Errorf("error getting kubernetes %s/%s cronjob: %w",
				ns, n, err)
		}

		return nil, nil
	}

	// Fetch jobs in pagination and filter out.
	var js []*batch.Job
	jlo := meta.ListOptions{Limit: 100}

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

	// Fetch pods with job label.
	var ps []core.Pod

	coreCli, err := coreclient.NewForConfig(op.RestConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes core client: %w", err)
	}

	for i := 0; i < len(js); i++ {
		ss := labels.SelectorFromSet(labels.Set{
			"controller-uid": string(js[i].UID),
			"job-name":       js[i].Name,
		}).String()

		var pl *core.PodList

		pl, err = coreCli.Pods(ns).
			List(ctx, meta.ListOptions{ResourceVersion: "0", LabelSelector: ss}) // Non quorum read.
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

func (op Operator) getPodsOfAny(
	ctx context.Context,
	gvr schema.GroupVersionResource,
	ns, n string,
) (*[]core.Pod, error) {
	// Fetch label selector with dynamic client.
	dynamicCli, err := dynamicclient.NewForConfig(op.RestConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes dynamic client: %w", err)
	}

	o, err := dynamicCli.Resource(gvr).Namespace(ns).
		Get(ctx, n, meta.GetOptions{ResourceVersion: "0"}) // Non quorum read.
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

	// Fetch pods with label selector.
	ss := s.String()

	coreCli, err := coreclient.NewForConfig(op.RestConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes core client: %w", err)
	}

	pl, err := coreCli.Pods(ns).
		List(ctx, meta.ListOptions{ResourceVersion: "0", LabelSelector: ss}) // Non quorum read.
	if err != nil {
		return nil, fmt.Errorf("error listing kubernetes %s pods with %s: %w",
			ns, ss, err)
	}

	return &pl.Items, nil
}
