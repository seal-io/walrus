package platformk8s

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicclient "k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	batchclient "k8s.io/client-go/kubernetes/typed/batch/v1"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/utils/pointer"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/k8s"
	"github.com/seal-io/seal/pkg/platform/operator"
	"github.com/seal-io/seal/pkg/platformk8s/intercept"
	"github.com/seal-io/seal/pkg/platformk8s/pods"
	"github.com/seal-io/seal/pkg/platformk8s/polymorphic"
)

const OperatorType = types.ConnectorTypeK8s

// NewOperator returns operator.Operator with the given options.
func NewOperator(ctx context.Context, opts operator.CreateOptions) (operator.Operator, error) {
	var restConfig, err = GetConfig(opts.Connector)
	if err != nil {
		return nil, err
	}
	var op = Operator{
		RestConfig: restConfig,
	}
	return op, nil
}

type Operator struct {
	RestConfig *rest.Config
}

// Type implements operator.Operator.
func (Operator) Type() operator.Type {
	return OperatorType
}

// IsConnected implements operator.Operator.
func (op Operator) IsConnected(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var err = k8s.Wait(ctx, op.RestConfig)
	return err == nil, err
}

// GetKeys implements operator.Operator.
func (op Operator) GetKeys(ctx context.Context, res model.ApplicationResource) (*operator.Keys, error) {
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
	var k = operator.Keys{
		Labels: []string{"Pod", "Container"},
	}
	if psp == nil {
		return &k, nil
	}
	var ps = *psp
	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].CreationTimestamp.Time.After(ps[j].CreationTimestamp.Time)
	})
	for i := 0; i < len(ps); i++ {
		var running = pods.IsPodRunning(&ps[i])
		var states = pods.GetContainerStates(&ps[i])

		var key = operator.Key{
			Name: ps[i].Name, // pod name
			Keys: make([]operator.Key, 0, len(states)),
		}
		for j := 0; j < len(states); j++ {
			key.Keys = append(key.Keys, operator.Key{
				Name:       states[j].Name,     // container name
				Value:      states[j].String(), // key
				Loggable:   pointer.Bool(states[j].State > pods.ContainerStateUnknown),
				Executable: pointer.Bool(running && states[j].State == pods.ContainerStateRunning),
			})
		}
		k.Keys = append(k.Keys, key)
	}
	return &k, nil
}

func (op Operator) getPods(ctx context.Context, res model.ApplicationResource) (*[]core.Pod, error) {
	// get group version resources.
	var gvr, ok = intercept.Terraform().GetGVR(res.Type)
	if !ok {
		// no error.
		return nil, nil
	}

	// parse resource name.
	ns, n, ok := parseNamespacedName(res.Name)
	if !ok {
		return nil, fmt.Errorf("failed to parse given resource name: %q", res.Name)
	}

	switch gvr.Resource {
	case "pods":
		var p, err = op.getPod(ctx, ns, n)
		if err != nil {
			return nil, err
		}
		return &[]core.Pod{*p}, nil
	case "cronjobs":
		return op.getPodsOfCronJob(ctx, ns, n)
	default:
	}
	return op.getPodsOfAny(ctx, gvr, ns, n)
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
		return nil, fmt.Errorf("error getting kubernetes %s/%s pod: %w",
			ns, n, err)
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
		return nil, fmt.Errorf("error getting kubernetes %s/%s cronjob: %w",
			ns, n, err)
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
			return nil, fmt.Errorf("error listing kubernetes %s pods with %s: %w",
				ns, ss, err)
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
		return nil, fmt.Errorf("error getting kubernetes %s %s/%s: %w",
			gvr.Resource, ns, n, err)
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

// Log implements operator.Operator.
func (op Operator) Log(ctx context.Context, res model.ApplicationResource, opts operator.LogOptions) error {
	// parse resource name.
	var ns, _, ok = parseNamespacedName(res.Name)
	if !ok {
		return fmt.Errorf("failed to parse given resource name: %q", res.Name)
	}

	// parse key.
	pn, ct, cn, ok := parseKey(opts.Key)
	if !ok {
		return fmt.Errorf("failed to parse given key: %q", opts.Key)
	}

	// confirm.
	var cli, err = coreclient.NewForConfig(op.RestConfig)
	if err != nil {
		return fmt.Errorf("error creating kubernetes client: %w", err)
	}
	p, err := cli.Pods(ns).
		Get(ctx, pn, meta.GetOptions{ResourceVersion: "0"}) // non quorum read
	if err != nil {
		return fmt.Errorf("error getting kubernetes pod %s/%s: %w", ns, pn, err)
	}
	if !pods.IsContainerExisted(p, pods.Container{Type: ct, Name: cn}) {
		return fmt.Errorf("given %s container %s is not ownered by %s/%s pod", ct, cn, ns, pn)
	}

	// stream.
	var stmOpts = &core.PodLogOptions{
		Container:    cn,
		Follow:       pods.IsContainerRunning(p, pods.Container{Type: ct, Name: cn}),
		Previous:     opts.Previous,
		SinceSeconds: opts.SinceSeconds,
		Timestamps:   opts.Timestamps,
	}
	if opts.Tail {
		stmOpts.TailLines = pointer.Int64(10)
	}
	stm, err := cli.Pods(ns).
		GetLogs(pn, stmOpts).
		Stream(ctx)
	if err != nil {
		return fmt.Errorf("failed to create log stream: %w", err)
	}
	defer func() { _ = stm.Close() }()
	var r = bufio.NewReader(stm)
	var w = opts.Out
	for {
		var bs []byte
		bs, err = r.ReadBytes('\n')
		if err != nil {
			if isTrivialError(err) {
				err = nil
			}
			break
		}
		_, err = w.Write(bs)
		if err != nil {
			if isTrivialError(err) {
				err = nil
			}
			break
		}
	}
	if err != nil {
		return fmt.Errorf("error streaming log: %w", err)
	}
	return nil
}

// Exec implements operator.Operator.
func (op Operator) Exec(ctx context.Context, res model.ApplicationResource, opts operator.ExecOptions) error {
	// parse resource name.
	var ns, _, ok = parseNamespacedName(res.Name)
	if !ok {
		return fmt.Errorf("failed to parse given resource name: %q", res.Name)
	}

	// parse key.
	pn, ct, cn, ok := parseKey(opts.Key)
	if !ok {
		return fmt.Errorf("failed to parse given key: %q", opts.Key)
	}

	// confirm.
	var cli, err = coreclient.NewForConfig(op.RestConfig)
	if err != nil {
		return fmt.Errorf("error creating kubernetes client: %w", err)
	}
	p, err := cli.Pods(ns).
		Get(ctx, pn, meta.GetOptions{ResourceVersion: "0"}) // non quorum read
	if err != nil {
		return fmt.Errorf("error getting kubernetes pod %s/%s: %w", ns, pn, err)
	}
	if !pods.IsContainerExisted(p, pods.Container{Type: ct, Name: cn}) {
		return fmt.Errorf("given %s container %s is not ownered by %s/%s pod", ct, cn, ns, pn)
	}
	if !pods.IsContainerRunning(p, pods.Container{Type: ct, Name: cn}) {
		return fmt.Errorf("given %s container %s is not running in %s/%s pod", ct, cn, ns, pn)
	}

	// stream.
	var stmURL = cli.RESTClient().Post().
		Resource("pods").
		Name(pn).
		Namespace(ns).
		SubResource("exec").
		VersionedParams(
			&core.PodExecOptions{
				Container: cn,
				Command:   strings.Split(opts.Shell, " "),
				Stdin:     true,
				Stdout:    true,
				TTY:       true,
			},
			scheme.ParameterCodec,
		).
		URL()
	var stmOpts = remotecommand.StreamOptions{
		Stdin:  opts.In,
		Stdout: opts.Out,
		Tty:    true,
	}
	if opts.Resizer != nil {
		stmOpts.TerminalSizeQueue = terminalResizer(opts.Resizer.Next)
	} else {
		stmOpts.TerminalSizeQueue = terminalSize(100, 100)
	}
	stm, err := remotecommand.NewSPDYExecutor(op.RestConfig, http.MethodPost, stmURL)
	if err != nil {
		return fmt.Errorf("failed to create exec stream: %w", err)
	}
	err = stm.StreamWithContext(ctx, stmOpts)
	if err != nil {
		if !isTrivialError(err) {
			return fmt.Errorf("error streaming exec: %w", err)
		}
	}
	return nil
}

// parseNamespacedName parses the given string into {namespace, name},
// returns false if not a valid namespaced name, e.g. kube-system/coredns.
func parseNamespacedName(s string) (namespace, name string, ok bool) {
	var ss = strings.SplitN(s, "/", 2)
	ok = len(ss) == 2
	if !ok {
		return
	}
	namespace = ss[0]
	name = ss[1]
	return
}

// parseKey parses the given string into {pod name, container type, container name},
// returns false if not a valid token, e.g. coredns-64897985d-6x2jm/container/coredns.
// valid container types have `initContainer`, `ephemeralContainer`, `container`.
func parseKey(s string) (podName, containerType, containerName string, ok bool) {
	var ss = strings.SplitN(s, "/", 3)
	ok = len(ss) == 3
	if !ok {
		return
	}
	podName = ss[0]
	containerType = ss[1]
	containerName = ss[2]
	return
}

// isTrivialError returns true if the given error can be ignored.
func isTrivialError(e error) bool {
	for _, t := range []error{
		io.EOF,
		context.Canceled,
	} {
		if errors.Is(e, t) {
			return true
		}
	}
	return false
}

// terminalSize returns terminalResizer with fixed width and height.
func terminalSize(width, height uint16) terminalResizer {
	var o sync.Once
	return func() (w uint16, h uint16, ok bool) {
		o.Do(func() {
			w = width
			h = height
		})
		return
	}
}

// terminalResizer implements remotecommand.TerminalSizeQueue.
type terminalResizer func() (width, height uint16, ok bool)

func (t terminalResizer) Next() *remotecommand.TerminalSize {
	if t == nil {
		return nil
	}
	var w, h, ok = t()
	if !ok {
		return nil
	}
	return &remotecommand.TerminalSize{
		Width:  w,
		Height: h,
	}
}
