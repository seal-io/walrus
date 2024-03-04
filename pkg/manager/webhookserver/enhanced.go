package webhookserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/seal-io/utils/netx"
	"github.com/seal-io/utils/pools/gopool"
	"github.com/seal-io/utils/stringx"
	"github.com/seal-io/utils/waitx"
	admreg "k8s.io/api/admissionregistration/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	ctrlwebhook "sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/seal-io/walrus/pkg/clients/clientset"
	"github.com/seal-io/walrus/pkg/system"
	"github.com/seal-io/walrus/pkg/systemkuberes"
	"github.com/seal-io/walrus/pkg/webhooks"
)

func Enhance(ln net.Listener, certDir string, cli clientset.Interface) ctrlwebhook.Server {
	addr := ln.Addr().(*net.TCPAddr)
	return enhanced{
		DefaultServer: &ctrlwebhook.DefaultServer{
			Options: ctrlwebhook.Options{
				Host:       addr.IP.String(),
				Port:       addr.Port,
				WebhookMux: http.NewServeMux(),
				CertDir:    certDir,
			},
		},
		ln:  ln,
		cli: cli,
	}
}

type enhanced struct {
	*ctrlwebhook.DefaultServer

	ln  net.Listener
	cli clientset.Interface
}

func (s enhanced) Start(ctx context.Context) error {
	// By default, we suggest to deploy with HA mode(by all-in-one YAML or Helm Chart),
	// which means we stay inside the loopback Kubernetes cluster.
	// So the system Kubernetes Service is created before webhook server start.
	cc := admreg.WebhookClientConfig{
		Service: &admreg.ServiceReference{
			Namespace: systemkuberes.SystemNamespaceName,
			Name:      systemkuberes.SystemRoutingServiceName,
			Port:      ptr.To(int32(s.Options.Port)),
		},
	}
	// However, if we stand closed to loopback Kubernetes cluster but not inside,
	// we can use the primary IP address to access the webhook server.
	if !system.LoopbackKubeInside.Get() && system.LoopbackKubeNearby.Get() {
		// NB(thxCode): launch multiple instances, only one takes working.
		ep := fmt.Sprintf("https://%s:%d", system.PrimaryIP.Get(), s.Options.Port)
		cc = admreg.WebhookClientConfig{
			URL: ptr.To(ep),
		}
	}
	// When no cert provided, we will use Kubernetes CertificateSigningRequest to generate the server cert,
	// we should load the root CA bundle for webhook configuration.
	if s.Options.CertDir == "" {
		// TODO(thxCode): The root CA might be expired,
		//   we can refresh the CA bundle in the future with restarting.
		//   A restarting-less way is needed in the future.
		err := waitx.PollUntilContextTimeout(ctx, time.Second, 30*time.Second, true,
			func(ctx context.Context) error {
				cm, err := s.cli.CoreV1().
					ConfigMaps(meta.NamespaceSystem).
					Get(ctx, "kube-root-ca.crt", meta.GetOptions{ResourceVersion: "0"})
				if err != nil {
					return fmt.Errorf("get kube-root-ca.crt configmap: %w", err)
				}
				if cm.Data["ca.crt"] == "" {
					return fmt.Errorf("empty ca.crt in kube-root-ca.crt configmap")
				}
				cc.CABundle = []byte(cm.Data["ca.crt"])
				return nil
			})
		if err != nil {
			return fmt.Errorf("get kube-root-ca.crt: %w", err)
		}
	}
	// Install webhook configurations.
	err := webhooks.InstallWebhookConfigurations(ctx, s.cli, cc)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Handler:           s.WebhookMux(),
		MaxHeaderBytes:    1 << 20,
		IdleTimeout:       90 * time.Second,
		ReadHeaderTimeout: 32 * time.Second,
	}

	gopool.Go(func() {
		<-ctx.Done()

		klog.Info("shutting down the webhook server with timeout of 1 minute")
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		_ = srv.Shutdown(ctx)
	})

	err = srv.Serve(s.ln)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s enhanced) StartedChecker() healthz.Checker {
	addr := net.JoinHostPort(s.Options.Host, stringx.FromInt(s.Options.Port))

	return func(req *http.Request) error {
		return netx.IsConnected(req.Context(), "tls", addr, 10*time.Second)
	}
}

func (s enhanced) WebhookMux() *http.ServeMux {
	return s.Options.WebhookMux
}

func (enhanced) IsDummy() bool {
	return false
}
