package k8s

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/seal-io/walrus/pkg/operator/k8s/key"
	"github.com/seal-io/walrus/pkg/operator/k8s/kube"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
)

// Exec implements operator.Operator.
func (op Operator) Exec(ctx context.Context, k string, opts optypes.ExecOptions) error {
	// Parse key.
	ns, pn, ct, cn, ok := key.Decode(k)
	if !ok {
		return fmt.Errorf("failed to parse given key: %q", k)
	}

	// Confirm.
	p, err := op.CoreCli.Pods(ns).
		Get(ctx, pn, meta.GetOptions{ResourceVersion: "0"}) // Non quorum read.
	if err != nil {
		return fmt.Errorf("error getting kubernetes pod %s/%s: %w", ns, pn, err)
	}

	if !kube.IsContainerExisted(p, kube.Container{Type: ct, Name: cn}) {
		return fmt.Errorf("given %s container %s is not ownered by %s/%s pod", ct, cn, ns, pn)
	}

	if !kube.IsContainerRunning(p, kube.Container{Type: ct, Name: cn}) {
		return fmt.Errorf("given %s container %s is not running in %s/%s pod", ct, cn, ns, pn)
	}

	// Stream.
	stmURL := op.CoreCli.RESTClient().Post().
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

	stmOpts := remotecommand.StreamOptions{
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

	return func() (w, h uint16, ok bool) {
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

	w, h, ok := t()
	if !ok {
		return nil
	}

	return &remotecommand.TerminalSize{
		Width:  w,
		Height: h,
	}
}
