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
	"k8s.io/apimachinery/pkg/util/wait"
	dynamicclient "k8s.io/client-go/dynamic"
	batchclient "k8s.io/client-go/kubernetes/typed/batch/v1"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	networkingclient "k8s.io/client-go/kubernetes/typed/networking/v1"
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/k8s"
	"github.com/seal-io/walrus/pkg/operator/k8s/polymorphic"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/utils/hash"
	"github.com/seal-io/walrus/utils/log"
)

const OperatorType = types.ConnectorTypeKubernetes

// New returns operator.Operator with the given options.
func New(ctx context.Context, opts optypes.CreateOptions) (optypes.Operator, error) {
	// NB(thxCode): disable timeout as we don't know the maximum time-cost of once operation,
	// and rely on the session context timeout control of each operation.
	restConfig, err := GetConfig(opts.Connector, WithoutTimeout())
	if err != nil {
		return nil, err
	}

	// NB(thxCode): since we rely on fewer APIs,
	// we don't need to initialize the nanny via kubernetes.NewForConfig.
	restCli, err := rest.HTTPClientFor(restConfig)
	if err != nil {
		return nil, err
	}

	coreCli, err := coreclient.NewForConfigAndClient(restConfig, restCli)
	if err != nil {
		return nil, err
	}

	batchCli, err := batchclient.NewForConfigAndClient(restConfig, restCli)
	if err != nil {
		return nil, err
	}

	networkingCli, err := networkingclient.NewForConfigAndClient(restConfig, restCli)
	if err != nil {
		return nil, err
	}

	dynamicCli, err := dynamicclient.NewForConfigAndClient(restConfig, restCli)
	if err != nil {
		return nil, err
	}

	op := Operator{
		Logger:        log.WithName("operator").WithName("k8s"),
		Identifier:    hash.SumStrings("k8s:", restConfig.Host, restConfig.APIPath),
		ModelClient:   opts.ModelClient,
		RestConfig:    restConfig,
		CoreCli:       coreCli,
		BatchCli:      batchCli,
		NetworkingCli: networkingCli,
		DynamicCli:    dynamicCli,
	}

	return op, nil
}

type Operator struct {
	Logger        log.Logger
	Identifier    string
	ModelClient   model.ClientSet
	RestConfig    *rest.Config
	CoreCli       *coreclient.CoreV1Client
	BatchCli      *batchclient.BatchV1Client
	NetworkingCli *networkingclient.NetworkingV1Client
	DynamicCli    *dynamicclient.DynamicClient
}

// Type implements operator.Operator.
func (Operator) Type() optypes.Type {
	return OperatorType
}

// IsConnected implements operator.Operator.
func (op Operator) IsConnected(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var lastErr error

	err := wait.PollUntilContextCancel(ctx, time.Second, true,
		func(_ context.Context) (bool, error) {
			// NB(shanewxy): Keep the real error from request.
			lastErr = k8s.IsConnected(context.TODO(), op.CoreCli.RESTClient())

			return lastErr == nil, ctx.Err()
		},
	)

	if lastErr != nil {
		err = lastErr // Use last error to overwrite context error while existed.
	}

	return err
}

// Burst implements operator.Operator.
func (op Operator) Burst() int {
	if op.RestConfig.Burst == 0 {
		return rest.DefaultBurst
	}

	return op.RestConfig.Burst
}

// ID implements operator.Operator.
func (op Operator) ID() string {
	return op.Identifier
}

func (op Operator) getPod(ctx context.Context, ns, n string) (*core.Pod, error) {
	// Fetch pod with name.
	p, err := op.CoreCli.Pods(ns).
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
	cj, err := op.BatchCli.CronJobs(ns).
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

		jl, err = op.BatchCli.Jobs(ns).List(ctx, jlo)
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
	o, err := op.DynamicCli.Resource(gvr).Namespace(ns).
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

	pl, err := op.CoreCli.Pods(ns).
		List(ctx, meta.ListOptions{ResourceVersion: "0", LabelSelector: ss}) // Non quorum read.
	if err != nil {
		return nil, fmt.Errorf("error listing kubernetes %s pods with %s: %w",
			ns, ss, err)
	}

	return &pl.Items, nil
}
