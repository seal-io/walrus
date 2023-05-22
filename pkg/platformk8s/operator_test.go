package platformk8s

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/utils/pointer"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/k8s"
	optypes "github.com/seal-io/seal/pkg/operator/types"
	"github.com/seal-io/seal/utils/log"
)

// TestOperator tests all actions of Operator if found a valid local kubeconfig.
func TestOperator(t *testing.T) {
	k8sCfg, err := k8s.GetConfig("")
	if err != nil {
		t.Logf("error getting kubeconfig: %v", err)
		t.Skip("cannot get kubeconfig")

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	op := Operator{
		Logger:     log.WithName("operator").WithName("k8s"),
		RestConfig: k8sCfg,
	}

	t.Run("IsConnected", func(t *testing.T) {
		err := op.IsConnected(ctx)
		if err != nil {
			t.Fatalf("error connecting kubernetes cluster: %v", err)
		}
	})

	// Start testing pod.
	cli, err := coreclient.NewForConfig(k8sCfg)
	if err != nil {
		t.Fatalf("error createing kubernetes client: %v", err)
	}
	p := &core.Pod{
		ObjectMeta: meta.ObjectMeta{
			Namespace:    core.NamespaceDefault,
			GenerateName: "nginx-",
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:  "nginx",
					Image: "nginx",
				},
			},
		},
	}

	p, err = cli.Pods(core.NamespaceDefault).Create(ctx, p, meta.CreateOptions{})
	if err != nil {
		t.Fatalf("error applying kubernetes pod: %v", err)
	}

	pw, err := cli.Pods(p.Namespace).Watch(ctx, meta.ListOptions{Watch: true, ResourceVersion: "0"})
	if err != nil {
		return
	}

	for evt := range pw.ResultChan() {
		wp, ok := evt.Object.(*core.Pod)
		if !ok {
			continue
		}

		if wp.Name != p.Name {
			continue
		}

		if wp.Status.Phase == core.PodRunning {
			pw.Stop()
			break
		}
	}

	defer func() {
		// Clean testing pod.
		_ = cli.Pods(p.Namespace).Delete(ctx, p.Name, meta.DeleteOptions{})
	}()

	// Mock application resource.
	res := &model.ApplicationResource{
		Type:         "kubernetes_pod",
		Name:         p.Namespace + "/" + p.Name,
		DeployerType: types.DeployerTypeTF,
	}

	t.Run("GetKeys", func(t *testing.T) {
		keys, err := op.GetKeys(ctx, res)
		if err != nil {
			t.Errorf("error getting keys: %v", err)
		}

		assert.Equalf(t, keys, &optypes.Keys{
			Labels: []string{"Pod", "Container"},
			Keys: []optypes.Key{
				{
					Name: p.Name,
					Keys: []optypes.Key{
						{
							Name:       "nginx",
							Value:      p.Namespace + "/" + p.Name + "/run/nginx",
							Loggable:   pointer.Bool(true),
							Executable: pointer.Bool(true),
						},
					},
				},
			},
		}, "invaild keys")
	})

	t.Run("Log", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		err := op.Log(ctx, p.Namespace+"/"+p.Name+"/run/nginx", optypes.LogOptions{
			Out:  testLogWriter(t.Logf),
			Tail: true,
		})
		if err != nil {
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Errorf("error logging: %v", err)
			}

			t.Log("finished")
		}
	})

	t.Run("Exec", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		r, w := io.Pipe()

		go func() {
			tk := time.NewTicker(2 * time.Second)
			defer tk.Stop()

			for {
				select {
				case <-ctx.Done():
					_ = r.Close()
					_ = w.Close()

					return
				case <-tk.C:
					_, _ = w.Write([]byte(fmt.Sprintf("echo %q \n", rand.String(16))))
				}
			}
		}()

		err = op.Exec(ctx, p.Namespace+"/"+p.Name+"/run/nginx", optypes.ExecOptions{
			Out:   testLogWriter(t.Logf),
			In:    r,
			Shell: "bash",
		})
		if err != nil {
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Errorf("error execing: %v", err)
			}

			t.Log("finished")
		}
	})
}

type testLogWriter func(format string, args ...any)

func (f testLogWriter) Write(p []byte) (n int, err error) {
	f(string(p))
	return len(p), nil
}
