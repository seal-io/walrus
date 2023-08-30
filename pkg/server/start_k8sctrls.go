package server

import (
	"context"
	"sync/atomic"

	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/k8sctrls"
)

type startK8sCtrlsOptions struct {
	MgrIsReady  *atomic.Bool
	RestConfig  *rest.Config
	ModelClient *model.Client
}

func (r *Server) startK8sCtrls(ctx context.Context, opts startK8sCtrlsOptions) error {
	mgr, err := k8sctrls.NewManager(k8sctrls.ManagerOptions{
		IsReady:            opts.MgrIsReady,
		LeaderElection:     r.KubeLeaderElection,
		LeaderLease:        r.KubeLeaderLease,
		LeaderRenewTimeout: r.KubeLeaderRenewTimeout,
	})
	if err != nil {
		return err
	}

	startOpts := k8sctrls.StartOptions{
		RestConfig: opts.RestConfig,
		SetupOptions: k8sctrls.SetupOptions{
			ModelClient: opts.ModelClient,
		},
	}

	return mgr.Start(ctx, startOpts)
}
