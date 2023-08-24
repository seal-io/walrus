package k8s

import (
	"bufio"
	"context"
	"fmt"
	"io"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/seal-io/walrus/pkg/operator/k8s/key"
	"github.com/seal-io/walrus/pkg/operator/k8s/kube"

	optypes "github.com/seal-io/walrus/pkg/operator/types"
)

// Log implements operator.Operator.
func (op Operator) Log(ctx context.Context, k string, opts optypes.LogOptions) error {
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

	// Stream.
	stmOpts := &core.PodLogOptions{
		Container:    cn,
		Follow:       !opts.WithoutFollow && kube.IsContainerRunning(p, kube.Container{Type: ct, Name: cn}),
		Previous:     opts.Previous,
		SinceSeconds: opts.SinceSeconds,
		Timestamps:   opts.Timestamps,
		TailLines:    opts.TailLines,
	}

	return GetPodLogs(ctx, op.CoreCli, ns, pn, opts.Out, stmOpts)
}

// GetPodLogs get the logs of a pod.
func GetPodLogs(
	ctx context.Context,
	cli *coreclient.CoreV1Client,
	namespace string,
	name string,
	w io.Writer,
	opts *core.PodLogOptions,
) error {
	stm, err := cli.Pods(namespace).
		GetLogs(name, opts).
		Stream(ctx)
	if err != nil {
		return fmt.Errorf("failed to create log stream: %w", err)
	}

	defer func() { _ = stm.Close() }()
	r := bufio.NewReader(stm)

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
