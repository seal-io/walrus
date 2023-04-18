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

	"github.com/seal-io/seal/utils/dynamicert"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
)

func NewServer() (*Server, error) {
	var logger = log.WithName("apis")
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

	var handler, err = s.Setup(c, opts.SetupOptions)
	if err != nil {
		return fmt.Errorf("error setting up apis server: %w", err)
	}
	var httpHandler = make(chan http.Handler)

	var g = gopool.GroupWithContextIn(c)

	// serve https.
	g.Go(func(ctx context.Context) error {
		if opts.TlsMode == TlsModeDisabled {
			httpHandler <- handler
			return nil
		}

		var h = handler
		var lg = newStdLogger(s.logger.WithName("https"))
		var ls, err = newTcpListener(ctx, opts.BindAddress, 443)
		if err != nil {
			return err
		}
		defer func() { _ = ls.Close() }()
		var tlsConfig = &tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
			MinVersion: tls.VersionTLS12,
		}
		switch opts.TlsMode {
		default: // TlsModeSelfGenerated
			var mgr = &dynamicert.Manager{
				Cache: dynamicert.DirCache(opts.TlsCertDir),
			}
			ls, h, err = mgr.Handle(ls, handler)
			if err != nil {
				return err
			}
			httpHandler <- http.HandlerFunc(redirectHandler)
		case TlsModeAutoGenerated:
			var mgr = &autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(opts.TlsAutoCertDomains...),
				Cache:      autocert.DirCache(opts.TlsCertDir),
			}
			tlsConfig.NextProtos = append(tlsConfig.NextProtos, acme.ALPNProto)
			tlsConfig.GetCertificate = func(i *tls.ClientHelloInfo) (*tls.Certificate, error) {
				if i.ServerName == "localhost" || i.ServerName == "" {
					var ni = *i
					ni.ServerName = opts.TlsAutoCertDomains[0]
					return mgr.GetCertificate(&ni)
				}
				return mgr.GetCertificate(i)
			}
			ls = tls.NewListener(ls, tlsConfig)
			httpHandler <- mgr.HTTPHandler(nil)
		case TlsModeCustomized:
			var cert, err = tls.LoadX509KeyPair(opts.TlsCertFile, opts.TlsPrivateKeyFile)
			if err != nil {
				return err
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
			httpHandler <- http.HandlerFunc(redirectHandler)
		}
		s.logger.Info("serving https")
		return serve(ctx, h, lg, ls)
	})

	// serve http.
	g.Go(func(ctx context.Context) error {
		var h = <-httpHandler
		var lg = newStdLogger(s.logger.WithName("http"))
		var ls, err = newTcpListener(ctx, opts.BindAddress, 80)
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
	var s = http.Server{
		Handler:     handler,
		ErrorLog:    errorLog,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}
	defer func() {
		var sCtx, sCancel = context.WithTimeout(context.Background(), 15*time.Second)
		defer sCancel()
		_ = s.Shutdown(sCtx)
	}()
	var err = s.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func newTcpListener(ctx context.Context, ip string, port int) (net.Listener, error) {
	var address = fmt.Sprintf("%s:%d", ip, port)
	var lc = net.ListenConfig{
		KeepAlive: 3 * time.Minute,
	}
	var ls, err = lc.Listen(ctx, "tcp", address)
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
	var host = r.Host
	var rawHost, _, err = net.SplitHostPort(host)
	if err == nil {
		host = net.JoinHostPort(rawHost, "443")
	}
	http.Redirect(w, r, "https://"+host+r.URL.RequestURI(), http.StatusFound)
}
