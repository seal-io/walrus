package server

import (
	"context"
	"errors"

	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/apis"
	"github.com/seal-io/seal/pkg/dao/model"
)

type setupApisOptions struct {
	K8sConfig   *rest.Config
	ModelClient *model.Client
}

func (r *Server) setupApis(ctx context.Context, opts setupApisOptions) error {
	var srv, err = apis.NewServer()
	if err != nil {
		return err
	}
	var serveOpts = apis.ServeOptions{
		SetupOptions: apis.SetupOptions{
			EnableAuthn: r.EnableAuthn,
			K8sConfig:   opts.K8sConfig,
			ModelClient: opts.ModelClient,
		},
		BindAddress: r.BindAddress,
	}
	switch {
	default:
		serveOpts.TlsMode = apis.TlsModeSelfGenerated
		serveOpts.TlsCertDir = r.TlsCertDir
	case r.TlsCertFile != "" && r.TlsPrivateKeyFile != "":
		serveOpts.TlsMode = apis.TlsModeCustomized
		serveOpts.TlsCertFile = r.TlsCertFile
		serveOpts.TlsPrivateKeyFile = r.TlsPrivateKeyFile
	case len(r.TlsAutoCertDomains) != 0:
		serveOpts.TlsMode = apis.TlsModeAutoGenerated
		serveOpts.TlsCertified = true
		serveOpts.TlsCertDir = r.TlsCertDir
		serveOpts.TlsAutoCertDomains = r.TlsAutoCertDomains
	}
	err = srv.Serve(ctx, serveOpts)
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}
