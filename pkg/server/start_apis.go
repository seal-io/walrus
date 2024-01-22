package server

import (
	"context"
	"errors"

	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/apis"
	"github.com/seal-io/walrus/pkg/dao/model"
)

type startApisOptions struct {
	K8sConfig   *rest.Config
	ModelClient *model.Client
}

func (r *Server) startApis(ctx context.Context, opts startApisOptions) error {
	srv, err := apis.NewServer()
	if err != nil {
		return err
	}

	serveOpts := apis.ServeOptions{
		SetupOptions: apis.SetupOptions{
			EnableAuthn:           r.EnableAuthn,
			ConnQPS:               r.ConnQPS,
			ConnBurst:             r.ConnBurst,
			WebsocketConnMaxPerIP: r.WebsocketConnMaxPerIP,
			K8sConfig:             opts.K8sConfig,
			ModelClient:           opts.ModelClient,
		},
		BindAddress:       r.BindAddress,
		BindWithDualStack: r.BindWithDualStack,
	}

	switch {
	default:
		serveOpts.TlsMode = apis.TlsModeSelfGenerated
		serveOpts.TlsCertDir = r.TlsCertDir
	case !r.EnableTls:
		serveOpts.TlsMode = apis.TlsModeDisabled
	case r.TlsCertFile != "" && r.TlsPrivateKeyFile != "":
		serveOpts.TlsMode = apis.TlsModeCustomized
		serveOpts.TlsCertFile = r.TlsCertFile
		serveOpts.TlsPrivateKeyFile = r.TlsPrivateKeyFile
	case len(r.TlsAutoCertDomains) != 0:
		serveOpts.TlsMode = apis.TlsModeAutoGenerated
		serveOpts.TlsCertDir = r.TlsCertDir
		serveOpts.TlsAutoCertDomains = r.TlsAutoCertDomains
	}

	err = srv.Serve(ctx, serveOpts)
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
