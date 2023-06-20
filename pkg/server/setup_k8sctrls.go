package server

import (
	"context"
	"errors"

	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/k8sctrls"
	"github.com/seal-io/seal/utils/gopool"
)

type setupK8sCtrlsOptions struct {
	K8sConfig     *rest.Config
	K8sCacheReady chan struct{}
	ModelClient   *model.Client
}

func (r *Server) setupK8sCtrls(ctx context.Context, opts setupK8sCtrlsOptions) error {
	mgr, err := k8sctrls.NewManager(opts.K8sConfig)
	if err != nil {
		return err
	}
	startOpts := k8sctrls.StartOptions{
		SetupOptions: k8sctrls.SetupOptions{
			ModelClient: opts.ModelClient,
		},
	}

	gopool.Go(func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if !mgr.IsReady(ctx) {
				continue
			}

			// Close the channel to notify the cache is ready.
			close(opts.K8sCacheReady)

			break
		}
	})

	err = mgr.Start(ctx, startOpts)
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
