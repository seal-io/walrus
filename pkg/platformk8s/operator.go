package platformk8s

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/utils/pointer"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/k8s"
	"github.com/seal-io/seal/pkg/platform/operator"
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
	if !isContainerExisted(p, ct, cn) {
		return fmt.Errorf("given %s %s is not ownered by %s/%s pod", ct, cn, ns, pn)
	}

	// stream.
	var stmOpts = &core.PodLogOptions{
		Container:    cn,
		Follow:       isContainerRunning(p, ct, cn),
		Previous:     opts.Previous,
		SinceSeconds: opts.SinceSeconds,
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
	if !isContainerExisted(p, ct, cn) {
		return fmt.Errorf("given %s %s is not ownered by %s/%s pod", ct, cn, ns, pn)
	}
	if p.Status.Phase != core.PodRunning {
		return fmt.Errorf("given %s/%s pod is not running", ns, pn)
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

// collection of container types.
const (
	initContainer      = "initContainer"
	ephemeralContainer = "ephemeralContainer"
	container          = "container"
)

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

// isContainerExisted returns true if the given type container belongs to the pod.
func isContainerExisted(pod *core.Pod, containerType, containerName string) bool {
	switch containerType {
	case initContainer:
		for i := range pod.Spec.InitContainers {
			if pod.Spec.InitContainers[i].Name == containerName {
				return true
			}
		}
	case ephemeralContainer:
		for i := range pod.Spec.EphemeralContainers {
			if pod.Spec.EphemeralContainers[i].Name == containerName {
				return true
			}
		}
	case container:
		for i := range pod.Spec.Containers {
			if pod.Spec.Containers[i].Name == containerName {
				return true
			}
		}
	}
	return false
}

// isContainerRunning returns true if the given type container is running.
func isContainerRunning(pod *core.Pod, containerType, containerName string) bool {
	switch containerType {
	case initContainer:
		for i := range pod.Status.InitContainerStatuses {
			if pod.Status.InitContainerStatuses[i].Name == containerName {
				return pod.Status.InitContainerStatuses[i].State.Running != nil
			}
		}
	case ephemeralContainer:
		for i := range pod.Status.EphemeralContainerStatuses {
			if pod.Status.EphemeralContainerStatuses[i].Name == containerName {
				return pod.Status.EphemeralContainerStatuses[i].State.Running != nil
			}
		}
	case container:
		for i := range pod.Status.ContainerStatuses {
			if pod.Status.ContainerStatuses[i].Name == containerName {
				return pod.Status.ContainerStatuses[i].State.Running != nil
			}
		}
	}
	return false
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
