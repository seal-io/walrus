package apis

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	stdlog "log"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"

	"github.com/seal-io/walrus/pkg/apis/config"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/kms"
	"github.com/seal-io/walrus/utils/dynacert"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

func NewServer() (*Server, error) {
	logger := log.WithName("api")

	return &Server{
		logger: logger,
	}, nil
}

type Server struct {
	logger log.Logger
}

type ServeOptions struct {
	SetupOptions

	BindAddress        string
	BindWithDualStack  bool
	TlsMode            TlsMode
	TlsCertFile        string
	TlsPrivateKeyFile  string
	TlsCertDir         string
	TlsAutoCertDomains []string
}

type TlsMode uint64

const (
	TlsModeDisabled TlsMode = iota
	TlsModeSelfGenerated
	TlsModeAutoGenerated
	TlsModeCustomized
)

type TlsCertDirMode = string

const (
	TlsCertDirByK8sSecrets TlsCertDirMode = "k8s:///secrets"
)

func (s *Server) Serve(c context.Context, opts ServeOptions) error {
	s.logger.Info("starting")

	config.TlsCertified.Set(opts.TlsCertified)

	handler, err := s.Setup(c, opts.SetupOptions)
	if err != nil {
		return fmt.Errorf("error setting up apis server: %w", err)
	}
	httpHandler := make(chan http.Handler)

	g := gopool.GroupWithContextIn(c)

	// Serve https.
	g.Go(func(ctx context.Context) error {
		if opts.TlsMode == TlsModeDisabled {
			s.logger.Info("serving in HTTP")

			httpHandler <- handler

			return nil
		}

		h := handler
		lg := newStdErrorLogger(s.logger.WithName("https"))

		nw, addr, err := parseBindAddress(opts.BindAddress, 443, opts.BindWithDualStack)
		if err != nil {
			return err
		}

		ls, err := newTcpListener(ctx, nw, addr)
		if err != nil {
			return err
		}

		defer func() { _ = ls.Close() }()

		tlsConfig := &tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
			MinVersion: tls.VersionTLS12,
		}

		switch opts.TlsMode {
		default: // TlsModeSelfGenerated.
			var cache dynacert.Cache
			if opts.TlsCertDir == TlsCertDirByK8sSecrets {
				cache, err = kms.NewKubernetes(ctx, kms.KubernetesOptions{
					Namespace: types.WalrusSystemNamespace,
					Config:    opts.K8sConfig,
					Group:     "dynacert",
					Logger:    s.logger.WithName("https"),
					RaiseNotFound: func(_ string) error {
						return dynacert.ErrCacheMiss
					},
				})
				if err != nil {
					return fmt.Errorf("failed to create HTTPs certificate cache: %w", err)
				}
			} else {
				cache = dynacert.DirCache(opts.TlsCertDir)
			}

			s.logger.InfoS("serving in HTTPs with self-generated keypair",
				"cache", opts.TlsCertDir)

			mgr := &dynacert.Manager{
				Cache: cache,
			}
			tlsConfig.GetCertificate = mgr.GetCertificate
			ls = tls.NewListener(ls, tlsConfig)
			httpHandler <- http.HandlerFunc(redirectHandler)
		case TlsModeAutoGenerated:
			var cache autocert.Cache
			if opts.TlsCertDir == TlsCertDirByK8sSecrets {
				cache, err = kms.NewKubernetes(ctx, kms.KubernetesOptions{
					Namespace: types.WalrusSystemNamespace,
					Config:    opts.K8sConfig,
					Group:     "autocert:" + strs.Join(",", opts.TlsAutoCertDomains...),
					Logger:    s.logger.WithName("https"),
					RaiseNotFound: func(_ string) error {
						return autocert.ErrCacheMiss
					},
				})
				if err != nil {
					return fmt.Errorf("failed to create HTTPs certificate cache: %w", err)
				}
			} else {
				cache = dynacert.DirCache(opts.TlsCertDir)
			}

			s.logger.InfoS("serving in HTTPs with auto-generated keypair",
				"domains", opts.TlsAutoCertDomains,
				"cache", opts.TlsCertDir)

			mgr := &autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				Cache:      cache,
				HostPolicy: autocert.HostWhitelist(opts.TlsAutoCertDomains...),
			}

			tlsConfig.NextProtos = append(tlsConfig.NextProtos, acme.ALPNProto)
			tlsConfig.GetCertificate = func(i *tls.ClientHelloInfo) (*tls.Certificate, error) {
				if i.ServerName == "localhost" || i.ServerName == "" {
					ni := *i
					ni.ServerName = opts.TlsAutoCertDomains[0]

					return mgr.GetCertificate(&ni)
				}

				return mgr.GetCertificate(i)
			}
			ls = tls.NewListener(ls, tlsConfig)
			httpHandler <- mgr.HTTPHandler(http.HandlerFunc(redirectHandler))
		case TlsModeCustomized:
			s.logger.Info("serving in HTTPs with custom keypair")

			cert, err := tls.LoadX509KeyPair(opts.TlsCertFile, opts.TlsPrivateKeyFile)
			if err != nil {
				return err
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
			ls = tls.NewListener(ls, tlsConfig)
			httpHandler <- http.HandlerFunc(redirectHandler)
		}

		s.logger.Infof("serving https on %q by %q", addr, nw)

		return serve(ctx, h, lg, ls)
	})

	// Serve http.
	g.Go(func(ctx context.Context) error {
		h := <-httpHandler
		lg := newStdErrorLogger(s.logger.WithName("http"))

		nw, addr, err := parseBindAddress(opts.BindAddress, 80, opts.BindWithDualStack)
		if err != nil {
			return err
		}

		ls, err := newTcpListener(ctx, nw, addr)
		if err != nil {
			return err
		}

		defer func() { _ = ls.Close() }()

		s.logger.Infof("serving http on %q by %q", addr, nw)

		return serve(ctx, h, lg, ls)
	})

	return g.Wait()
}

func serve(ctx context.Context, handler http.Handler, errorLog *stdlog.Logger, listener net.Listener) error {
	s := http.Server{
		Handler:     handler,
		ErrorLog:    errorLog,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}
	defer func() {
		sCtx, sCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer sCancel()
		_ = s.Shutdown(sCtx)
	}()

	err := s.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func parseBindAddress(ip string, port int, dual bool) (network, address string, err error) {
	p := net.ParseIP(ip)
	if p == nil {
		return "", "", fmt.Errorf("invalid IP address: %s", ip)
	}

	nw := "tcp"

	p = p.To4()
	if p != nil {
		if !dual {
			nw = "tcp4"
		}

		return nw, fmt.Sprintf("%s:%d", p.String(), port), nil
	}

	if !dual {
		nw = "tcp6"
	}

	return nw, fmt.Sprintf("[%s]:%d", ip, port), nil
}

func newTcpListener(ctx context.Context, network, address string) (net.Listener, error) {
	lc := net.ListenConfig{
		KeepAlive: 3 * time.Minute,
	}

	ls, err := lc.Listen(ctx, network, address)
	if err != nil {
		return nil, fmt.Errorf("error creating %s listener for %s: %w",
			network, address, err)
	}

	return ls, nil
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Use HTTPS", http.StatusBadRequest)
		return
	}

	// From Kubernetes probes guide,
	// https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes,
	// the `User-Agent: kube-probe/<version>` header identify the incoming request is for Kubelet health check,
	// in order to avoid stuck in the readiness check, we don't redirect the probes request.
	if ua := r.Header.Get("User-Agent"); strings.HasPrefix(ua, "kube-probe/") {
		w.WriteHeader(http.StatusOK)
		return
	}

	host := r.Host
	if rawHost, _, err := net.SplitHostPort(host); err == nil {
		host = net.JoinHostPort(rawHost, "443")
	}

	http.Redirect(w, r, "https://"+host+r.URL.RequestURI(), http.StatusFound)
}
