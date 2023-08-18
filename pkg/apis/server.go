package apis

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	stdlog "log"
	"net"
	"net/http"
	"time"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"

	"github.com/seal-io/walrus/utils/dynacert"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
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

func (s *Server) Serve(c context.Context, opts ServeOptions) error {
	s.logger.Info("starting")

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

		ls, err := newTcpListener(ctx, opts.BindAddress, 443)
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
			s.logger.InfoS("serving in HTTPs with self-generated keypair",
				"cache", opts.TlsCertDir)

			mgr := &dynacert.Manager{
				Cache: dynacert.DirCache(opts.TlsCertDir),
			}
			tlsConfig.GetCertificate = mgr.GetCertificate
			ls = tls.NewListener(ls, tlsConfig)
			httpHandler <- http.HandlerFunc(redirectHandler)
		case TlsModeAutoGenerated:
			s.logger.InfoS("serving in HTTPs with auto-generated keypair",
				"domains", opts.TlsAutoCertDomains,
				"cache", opts.TlsCertDir)

			mgr := &autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				Cache:      autocert.DirCache(opts.TlsCertDir),
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

		s.logger.Info("serving https")

		return serve(ctx, h, lg, ls)
	})

	// Serve http.
	g.Go(func(ctx context.Context) error {
		h := <-httpHandler
		lg := newStdErrorLogger(s.logger.WithName("http"))

		ls, err := newTcpListener(ctx, opts.BindAddress, 80)
		if err != nil {
			return err
		}

		defer func() { _ = ls.Close() }()
		s.logger.Info("serving http")

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

func newTcpListener(ctx context.Context, ip string, port int) (net.Listener, error) {
	address := fmt.Sprintf("%s:%d", ip, port)
	lc := net.ListenConfig{
		KeepAlive: 3 * time.Minute,
	}

	ls, err := lc.Listen(ctx, "tcp", address)
	if err != nil {
		return nil, fmt.Errorf("error creating tcp listener for %s: %w", address, err)
	}

	return ls, nil
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Use HTTPS", http.StatusBadRequest)
		return
	}

	host := r.Host

	rawHost, _, err := net.SplitHostPort(host)
	if err == nil {
		host = net.JoinHostPort(rawHost, "443")
	}

	http.Redirect(w, r, "https://"+host+r.URL.RequestURI(), http.StatusFound)
}
