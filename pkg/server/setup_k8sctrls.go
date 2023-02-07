package server

import (
	"context"
	"errors"

	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/k8sctrls"
)

type setupK8sCtrlsOptions struct {
	K8sConfig   *rest.Config
	ModelClient *model.Client
}

func (r *Server) setupK8sCtrls(ctx context.Context, opts setupK8sCtrlsOptions) error {
	var mgr, err = k8sctrls.NewManager(opts.K8sConfig)
	if err != nil {
		return err
	}
	var startOpts = k8sctrls.StartOptions{
		SetupOptions: k8sctrls.SetupOptions{
			ModelClient: opts.ModelClient,
		},
	}
	err = mgr.Start(ctx, startOpts)
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}
