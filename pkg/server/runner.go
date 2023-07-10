package server

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	stdlog "log"
	"net"
	"path/filepath"
	"strings"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"           // Db = postgres.
	_ "github.com/mattn/go-sqlite3" // Db = sqlite3.
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	klogv2 "k8s.io/klog/v2"

	"github.com/seal-io/seal/pkg/casdoor"
	"github.com/seal-io/seal/pkg/consts"
	"github.com/seal-io/seal/pkg/dao/model"
	_ "github.com/seal-io/seal/pkg/dao/model/runtime" // Default = ent.
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/k8s"
	"github.com/seal-io/seal/pkg/rds"
	"github.com/seal-io/seal/utils/clis"
	"github.com/seal-io/seal/utils/cryptox"
	"github.com/seal-io/seal/utils/files"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
	"github.com/seal-io/seal/utils/version"
)

type Server struct {
	Logger clis.Logger

	BindAddress        string
	EnableTls          bool
	TlsCertFile        string
	TlsPrivateKeyFile  string
	TlsCertDir         string
	TlsAutoCertDomains []string
	BootstrapPassword  string
	ConnQPS            int
	ConnBurst          int
	GopoolWorkerFactor int

	KubeConfig      string
	KubeConnTimeout time.Duration
	KubeConnQPS     float64
	KubeConnBurst   int

	DataSourceAddress        string
	DataSourceConnMaxOpen    int
	DataSourceConnMaxIdle    int
	DataSourceConnMaxLife    time.Duration
	DataSourceDataEncryptAlg string
	DataSourceDataEncryptKey []byte

	EnableAuthn         bool
	AuthnSessionMaxIdle time.Duration
	CasdoorServer       string
}

func New() *Server {
	return &Server{
		BindAddress:           "0.0.0.0",
		EnableTls:             true,
		TlsCertDir:            filepath.FromSlash(filepath.Join(consts.DataDir, "tls")),
		ConnQPS:               10,
		ConnBurst:             20,
		KubeConnTimeout:       5 * time.Minute,
		KubeConnQPS:           16,
		KubeConnBurst:         64,
		DataSourceConnMaxOpen: 15,
		DataSourceConnMaxIdle: 5,
		DataSourceConnMaxLife: 10 * time.Minute,
		EnableAuthn:           true,
		AuthnSessionMaxIdle:   30 * time.Minute,
		GopoolWorkerFactor:    100,
	}
}

func (r *Server) Flags(cmd *cli.Command) {
	flags := [...]cli.Flag{
		&cli.StringFlag{
			Name:        "bind-address",
			Usage:       "The IP address on which to listen.",
			Destination: &r.BindAddress,
			Value:       r.BindAddress,
			Action: func(c *cli.Context, s string) error {
				if s != "" && net.ParseIP(s) == nil {
					return errors.New("--bind-address: invalid IP address")
				}
				return nil
			},
		},
		&cli.BoolFlag{
			Name:        "enable-tls",
			Usage:       "Enable HTTPs.",
			Destination: &r.EnableTls,
			Value:       r.EnableTls,
		},
		&cli.StringFlag{
			Name: "tls-cert-file",
			Usage: "The file containing the default x509 certificate for HTTPS. " +
				"If any CA certs, concatenated after server cert file. ",
			Destination: &r.TlsCertFile,
			Value:       r.TlsCertFile,
			Action: func(c *cli.Context, s string) error {
				if s != "" &&
					!files.Exists(s) {
					return errors.New("--tls-cert-file: file is not existed")
				}
				return nil
			},
		},
		&cli.StringFlag{
			Name:        "tls-private-key-file",
			Usage:       "The file containing the default x509 private key matching --tls-cert-file.",
			Destination: &r.TlsPrivateKeyFile,
			Value:       r.TlsPrivateKeyFile,
			Action: func(c *cli.Context, s string) error {
				if s != "" &&
					!files.Exists(s) {
					return errors.New("--tls-private-key-file: file is not existed")
				}
				return nil
			},
		},
		&cli.StringFlag{
			Name: "tls-cert-dir",
			Usage: "The directory where the TLS certs are located. " +
				"If --tls-cert-file and --tls-private-key-file are provided, this flag will be ignored. " +
				"If --tls-cert-file and --tls-private-key-file are not provided, " +
				"the self-signed certificate and key are generated and saved to the directory specified by this flag.",
			Destination: &r.TlsCertDir,
			Value:       r.TlsCertDir,
			Action: func(c *cli.Context, s string) error {
				if s == "" &&
					(c.String("tls-cert-file") == "" || c.String("tls-private-key-file") == "") {
					return errors.New(
						"--tls-cert-dir: must be filled if --tls-cert-file and --tls-private-key-file are not provided",
					)
				}
				return nil
			},
		},
		&cli.StringSliceFlag{
			Name: "tls-auto-cert-domains",
			Usage: "The domains to accept ACME HTTP-01 or TLS-ALPN-01 challenge to " +
				"generate HTTPS x509 certificate and private key, " +
				"and saved to the directory specified by --tls-cert-dir. " +
				"If --tls-cert-file and --tls-key-file are provided, this flag will be ignored.",
			Action: func(c *cli.Context, v []string) error {
				f := field.NewPath("--tls-auto-cert-domains")
				for i := range v {
					if err := validation.IsFullyQualifiedDomainName(f, v[i]).ToAggregate(); err != nil {
						return err
					}
				}
				if len(v) != 0 &&
					(c.String("tls-cert-dir") == "" &&
						(c.String("tls-cert-file") == "" || c.String("tls-private-key-file") == "")) {
					return errors.New("--tls-cert-dir: must be filled")
				}
				r.TlsAutoCertDomains = v
				return nil
			},
			Value: cli.NewStringSlice(r.TlsAutoCertDomains...),
		},
		&cli.StringFlag{
			Name:  "bootstrap-password",
			Usage: "The password to bootstrap instead of random generating.",
			Action: func(c *cli.Context, s string) error {
				if strs.StringWidth(s) < 6 {
					return errors.New("invalid bootstrap-password: too short")
				}
				return nil
			},
			Destination: &r.BootstrapPassword,
			Value:       r.BootstrapPassword,
		},
		&cli.IntFlag{
			Name:        "conn-qps",
			Usage:       "The qps(maximum average number per second) when dialing the server.",
			Destination: &r.ConnQPS,
			Value:       r.ConnQPS,
		},
		&cli.IntFlag{
			Name:        "conn-burst",
			Usage:       "The burst(maximum number at the same moment) when dialing the server.",
			Destination: &r.ConnBurst,
			Value:       r.ConnBurst,
		},
		&cli.StringFlag{
			Name:        "kubeconfig",
			Usage:       "The configuration path of the worker kubernetes cluster.",
			Destination: &r.KubeConfig,
			Value:       r.KubeConfig,
		},
		&cli.DurationFlag{
			Name:        "kube-conn-timeout",
			Usage:       "The timeout for dialing the worker kubernetes cluster.",
			Destination: &r.KubeConnTimeout,
			Value:       r.KubeConnTimeout,
		},
		&cli.Float64Flag{
			Name:        "kube-conn-qps",
			Usage:       "The qps(maximum average number per second) when dialing the worker kubernetes cluster.",
			Destination: &r.KubeConnQPS,
			Value:       r.KubeConnQPS,
		},
		&cli.IntFlag{
			Name:        "kube-conn-burst",
			Usage:       "The burst(maximum number at the same moment) when dialing the worker kubernetes cluster.",
			Destination: &r.KubeConnBurst,
			Value:       r.KubeConnBurst,
		},
		&cli.StringFlag{
			Name: "data-source-address",
			Usage: "The addresses for connecting data source, e.g. " +
				"Postgres(postgres://[username[:password]@][protocol[(address)]]/dbname" +
				"[?param1=value1&...&paramN=valueN]), " +
				"MySQL(mysql://[username[:password]@][protocol[(address)]]/dbname" +
				"[?param1=value1&...&paramN=valueN]), " +
				"SQLite3(file:dbpath[?param1=value1&...&paramN=valueN]).",
			Destination: &r.DataSourceAddress,
			Value:       r.DataSourceAddress,
		},
		&cli.IntFlag{
			Name:        "data-source-conn-max-open",
			Usage:       "The maximum opening connections for connecting data source.",
			Destination: &r.DataSourceConnMaxOpen,
			Value:       r.DataSourceConnMaxOpen,
		},
		&cli.IntFlag{
			Name:        "data-source-conn-max-idle",
			Usage:       "The maximum idling connections for connecting data source.",
			Destination: &r.DataSourceConnMaxIdle,
			Value:       r.DataSourceConnMaxIdle,
		},
		&cli.DurationFlag{
			Name:        "data-source-conn-max-life",
			Usage:       "The maximum lifetime for connecting data source.",
			Destination: &r.DataSourceConnMaxLife,
			Value:       r.DataSourceConnMaxLife,
		},
		&cli.StringFlag{
			Name: "data-source-data-encryption",
			Usage: "The algorithm and key(in-hex string) for encrypting the user credentials storing in data source, " +
				"e.g. aesgcm:3a9b4000d0ad8fbcd01eb922231d395d, " +
				"aesgcm:b4d1c09dcf62214a05d85548b9217b34da63224d2605938abb6bf384050d2222",
			Action: func(c *cli.Context, s string) error {
				ss := strings.SplitN(s, ":", 2)
				if len(ss) != 2 {
					return errors.New("invalid data-source-data-encryption: illegal format")
				}
				alg := ss[0]
				key, err := hex.DecodeString(ss[1])
				if err != nil {
					return errors.New("invalid data-source-data-encryption: must in hex string")
				}
				switch alg {
				default:
					return errors.New(
						"invalid data-source-data-encryption: unknown algorithm " + alg,
					)
				case "aesgcm":
					if ks := len(key); ks != 16 && ks != 32 {
						return errors.New(
							"invalid data-source-data-encryption: must in 16 bytes or 32 bytes",
						)
					}
				}
				r.DataSourceDataEncryptAlg, r.DataSourceDataEncryptKey = alg, key
				return nil
			},
		},
		&cli.BoolFlag{
			Name:        "enable-authn",
			Usage:       "Enable authentication.",
			Destination: &r.EnableAuthn,
			Value:       r.EnableAuthn,
		},
		&cli.DurationFlag{
			Name: "authn-session-max-idle",
			Usage: "The maximum idling duration for keeping authenticated session, " +
				"it represents the max-age of authenticated cookie.",
			Action: func(c *cli.Context, d time.Duration) error {
				if d < 0 {
					return errors.New("invalid authn-session-max-idle: negative")
				}
				return nil
			},
			Destination: &r.AuthnSessionMaxIdle,
			Value:       r.AuthnSessionMaxIdle,
		},
		&cli.StringFlag{
			Name:        "casdoor-server",
			Usage:       "The URL for connecting external casdoor server.",
			Destination: &r.CasdoorServer,
			Value:       r.CasdoorServer,
		},
		&cli.IntFlag{
			Name: "gopool-worker-factor",
			Usage: "The gopool worker factor determines the number of tasks of the goroutine worker pool," +
				"it is calculated by the number of CPU cores multiplied by this factor.",
			Destination: &r.GopoolWorkerFactor,
			Value:       r.GopoolWorkerFactor,
		},
	}
	for i := range flags {
		cmd.Flags = append(cmd.Flags, flags[i])
	}

	r.Logger.Flags(cmd)
}

func (r *Server) Before(cmd *cli.Command) {
	r.Logger.Before(cmd)
	// Compatible with other loggers.
	logger := log.GetLogger()
	stdlog.SetOutput(logger)
	logrus.SetOutput(logger)
	klog.SetOutput(logger)
	klogv2.SetLogger(log.AsLogr(logger))
}

func (r *Server) Action(cmd *cli.Command) {
	cmd.Action = func(c *cli.Context) error {
		return r.Run(c.Context)
	}
}

func (r *Server) Run(c context.Context) error {
	if err := r.configure(c); err != nil {
		log.Errorf("error configuring: %v", err)
		return fmt.Errorf("error configuring: %w", err)
	}

	g, ctx := gopool.GroupWithContext(c)

	// Get kubernetes config.
	k8sCfg, err := k8s.GetConfig(r.KubeConfig)
	if err != nil {
		// If not found, launch embedded kubernetes.
		var e k8s.Embedded

		g.Go(func() error {
			log.Info("running embedded kubernetes")

			err := e.Run(ctx)
			if err != nil {
				log.Errorf("error running embedded kubernetes: %v", err)
			}

			return err
		})
		// And get embedded kubernetes config.
		r.KubeConfig, k8sCfg, err = e.GetConfig(ctx)
		if err != nil {
			return fmt.Errorf("error getting embedded kubernetes config: %w", err)
		}
	}
	// Wait kubernetes to be ready.
	if err = k8s.Wait(ctx, k8sCfg); err != nil {
		return fmt.Errorf("error waiting kubernetes cluster ready: %w", err)
	}

	r.setKubernetesConfig(k8sCfg)

	// Load database driver.
	rdsDrvDialect, rdsDrv, err := rds.LoadDriver(r.DataSourceAddress)
	if err != nil {
		// If not found, launch embedded database.
		var e rds.Embedded

		g.Go(func() error {
			log.Info("running embedded database")

			err := e.Run(ctx)
			if err != nil {
				log.Errorf("error running embedded database: %v", err)
			}

			return err
		})
		// And get embedded database driver.
		r.DataSourceAddress, rdsDrvDialect, rdsDrv, err = e.GetDriver(ctx)
		if err != nil {
			return fmt.Errorf("error getting embedded database driver: %w", err)
		}
	}
	// Wait database to be ready.
	if err = rds.Wait(ctx, rdsDrv); err != nil {
		return fmt.Errorf("error waiting database ready: %w", err)
	}

	r.setDataSourceDriver(rdsDrv)

	if r.EnableAuthn {
		// Enable authentication.
		if r.CasdoorServer == "" {
			// If not specified, launch embedded casdoor,.
			var e casdoor.Embedded

			g.Go(func() error {
				log.Info("running embedded casdoor")

				err := e.Run(ctx, r.DataSourceAddress)
				if err != nil {
					log.Errorf("error running embedded casdoor: %v", err)
				}

				return err
			})
			// And get embedded casdoor address.
			r.CasdoorServer, err = e.GetAddress(ctx)
			if err != nil {
				return fmt.Errorf("error getting embedded casdor: %w", err)
			}
		}
		// Wait casdoor to be ready.
		if err = casdoor.Wait(ctx, r.CasdoorServer); err != nil {
			return fmt.Errorf("error waiting casdoor ready: %w", err)
		}
	}

	// Initialize some resources.
	log.Info("initializing")
	modelClient := getModelClient(rdsDrvDialect, rdsDrv)

	initOpts := initOptions{
		K8sConfig:     k8sCfg,
		K8sCacheReady: make(chan struct{}),
		ModelClient:   modelClient,
		SkipTLSVerify: len(r.TlsAutoCertDomains) != 0,
		RdsDialect:    rdsDrvDialect,
		RdsDriver:     rdsDrv,
	}
	if err = r.init(ctx, initOpts); err != nil {
		log.Errorf("error initializing: %v", err)
		return fmt.Errorf("error initializing: %w", err)
	}

	// Setup k8s controllers.
	setupK8sCtrlsOpts := setupK8sCtrlsOptions{
		K8sConfig:     k8sCfg,
		K8sCacheReady: initOpts.K8sCacheReady,
		ModelClient:   modelClient,
	}

	g.Go(func() error {
		log.Info("setting up kubernetes controller")

		err := r.setupK8sCtrls(ctx, setupK8sCtrlsOpts)
		if err != nil {
			log.Errorf("error setting up kubernetes controller: %v", err)
		}

		return err
	})

	// Setup apis.
	setupApisOpts := setupApisOptions{
		ModelClient: modelClient,
		K8sConfig:   k8sCfg,
	}

	g.Go(func() error {
		log.Info("setting up apis")

		err := r.setupApis(ctx, setupApisOpts)
		if err != nil {
			log.Errorf("error setting up apis: %v", err)
		}

		return err
	})

	return g.Wait()
}

// configure performs necessary configuration to support the whole server running.
func (r *Server) configure(_ context.Context) error {
	// Configure gopool.
	gopool.ResetPool(r.GopoolWorkerFactor)

	// Configure data encryption.
	if r.DataSourceDataEncryptKey != nil {
		var (
			alg = r.DataSourceDataEncryptAlg
			key = r.DataSourceDataEncryptKey

			enc cryptox.Encryptor
			err error
		)

		switch alg {
		default:
			return fmt.Errorf("unknown data encryptor algorithm: %s", alg)
		case "aesgcm":
			enc, err = cryptox.AesGcm(key)
		}

		if err != nil {
			return fmt.Errorf("error gaining data encryptor: %w", err)
		}

		crypto.EncryptorConfig.Set(enc)
	}

	return nil
}

func (r *Server) setKubernetesConfig(cfg *rest.Config) {
	cfg.Timeout = r.KubeConnTimeout
	cfg.QPS = float32(r.KubeConnQPS)
	cfg.Burst = r.KubeConnBurst
	cfg.UserAgent = version.GetUserAgent()
}

func (r *Server) setDataSourceDriver(drv *sql.DB) {
	drv.SetConnMaxLifetime(r.DataSourceConnMaxLife)
	drv.SetMaxIdleConns(r.DataSourceConnMaxIdle)
	drv.SetMaxOpenConns(r.DataSourceConnMaxOpen)
}

func getModelClient(drvDialect string, drv *sql.DB) *model.Client {
	logger := log.WithName("model")

	return model.NewClient(
		model.Log(func(args ...any) { logger.Debug(args...) }),
		model.Driver(entsql.NewDriver(drvDialect, entsql.Conn{ExecQuerier: drv})),
	)
}
