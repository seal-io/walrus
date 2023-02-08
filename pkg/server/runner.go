package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"           // db = postgres
	_ "github.com/mattn/go-sqlite3" // db = sqlite3
	"github.com/urfave/cli/v2"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/casdoor"
	"github.com/seal-io/seal/pkg/dao/model"
	_ "github.com/seal-io/seal/pkg/dao/model/runtime" // default = ent
	"github.com/seal-io/seal/pkg/k8s"
	"github.com/seal-io/seal/pkg/rds"
	"github.com/seal-io/seal/utils/clis"
	"github.com/seal-io/seal/utils/files"
	"github.com/seal-io/seal/utils/gopool"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
)

type Server struct {
	Logger clis.Logger

	BindAddress        string
	TlsCertFile        string
	TlsPrivateKeyFile  string
	TlsCertDir         string
	TlsAutoCertDomains []string
	BootstrapPassword  string

	KubeConfig      string
	KubeConnTimeout time.Duration
	KubeConnQPS     float64
	KubeConnBurst   int

	DataSourceAddress     string
	DataSourceConnMaxOpen int
	DataSourceConnMaxIdle int
	DataSourceConnMaxLife time.Duration

	EnableAuthn   bool
	CasdoorServer string
}

func New() *Server {
	return &Server{
		BindAddress:           "0.0.0.0",
		TlsCertDir:            filepath.FromSlash("/var/run/seal"),
		KubeConnTimeout:       5 * time.Minute,
		KubeConnQPS:           16,
		KubeConnBurst:         64,
		DataSourceConnMaxOpen: 15,
		DataSourceConnMaxIdle: 5,
		DataSourceConnMaxLife: 10 * time.Minute,
		EnableAuthn:           true,
	}
}

func (r *Server) Flags(cmd *cli.Command) {
	var flags = [...]cli.Flag{
		&cli.StringFlag{
			Name:        "bind-address",
			Usage:       "The IP address on which to listen.",
			Destination: &r.BindAddress,
			Value:       r.BindAddress,
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
				"If --tls-cert-file and --tls-private-key-file are not provided, the self-signed certificate and key are generated and saved to the directory specified by this flag.",
			Destination: &r.TlsCertDir,
			Value:       r.TlsCertDir,
			Action: func(c *cli.Context, s string) error {
				if s == "" &&
					(c.String("tls-cert-file") == "" || c.String("tls-private-key-file") == "") {
					return errors.New("--tls-cert-dir: must be filled if --tls-cert-file and --tls-private-key-file are not provided")
				}
				return nil
			},
		},
		&cli.StringSliceFlag{
			Name: "tls-auto-cert-domains",
			Usage: "The domains to accept ACME HTTP-01 challenge to generate HTTPS x509 certificate and private key, " +
				"and saved to the directory specified by --tls-cert-dir. " +
				"If --tls-cert-file and --tls-key-file are provided, this flag will be ignored.",
			Action: func(c *cli.Context, v []string) error {
				var f = field.NewPath("--tls-auto-cert-domains")
				for i := range v {
					if err := validation.IsFullyQualifiedDomainName(f, v[i]).ToAggregate(); err != nil {
						return err
					}
				}
				if len(v) != 0 &&
					(c.String("tls-cert-dir") == "" && c.String("tls-cert-file") == "" || c.String("tls-private-key-file") == "") {
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
			Usage:       "The qps when dialing the worker kubernetes cluster.",
			Destination: &r.KubeConnQPS,
			Value:       r.KubeConnQPS,
		},
		&cli.IntFlag{
			Name:        "kube-conn-burst",
			Usage:       "The burst when dialing the worker kubernetes cluster.",
			Destination: &r.KubeConnBurst,
			Value:       r.KubeConnBurst,
		},
		&cli.StringFlag{
			Name: "data-source-address",
			Usage: "The addresses for connecting data source, e.g. " +
				"Postgres(postgres://[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]), " +
				"MySQL(mysql://[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]), " +
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
		&cli.BoolFlag{
			Name:        "enable-authn",
			Usage:       "Enable authentication",
			Destination: &r.EnableAuthn,
			Value:       r.EnableAuthn,
		},
		&cli.StringFlag{
			Name:        "casdoor-server",
			Usage:       "The URL for connecting external casdoor server.",
			Destination: &r.CasdoorServer,
			Value:       r.CasdoorServer,
		},
	}
	for i := range flags {
		cmd.Flags = append(cmd.Flags, flags[i])
	}

	r.Logger.Flags(cmd)
}

func (r *Server) Before(cmd *cli.Command) {
	r.Logger.Before(cmd)
}

func (r *Server) Action(cmd *cli.Command) {
	cmd.Action = func(c *cli.Context) error {
		return r.Run(c.Context)
	}
}

func (r *Server) Run(c context.Context) error {
	var g, ctx = gopool.WithContext(c)

	// get kubernetes config.
	var k8sCfg, err = k8s.GetConfig(r.KubeConfig)
	if err != nil {
		// if not found, launch embedded kubernetes
		var e k8s.Embedded
		g.Go(func() error {
			log.Info("running embedded kubernetes")
			var err = e.Run(ctx)
			if err != nil {
				log.Errorf("error running embedded kubernetes: %v", err)
			}
			return err
		})
		// and get embedded kubernetes config.
		r.KubeConfig, k8sCfg, err = e.GetConfig(ctx)
		if err != nil {
			return fmt.Errorf("error getting embedded kubernetes config: %w", err)
		}
	}
	// wait kubernetes to be ready.
	if err = k8s.Wait(ctx, k8sCfg); err != nil {
		return fmt.Errorf("error waiting kubernetes cluster ready: %w", err)
	}
	r.setKubernetesConfig(k8sCfg)

	// load database driver.
	rdsDrvDialect, rdsDrv, err := rds.LoadDriver(r.DataSourceAddress)
	if err != nil {
		// if not found, launch embedded database
		var e rds.Embedded
		g.Go(func() error {
			log.Info("running embedded database")
			var err = e.Run(ctx)
			if err != nil {
				log.Errorf("error running embedded database: %v", err)
			}
			return err
		})
		// and get embedded database driver.
		r.DataSourceAddress, rdsDrvDialect, rdsDrv, err = e.GetDriver(ctx)
		if err != nil {
			return fmt.Errorf("error getting embedded database driver: %w", err)
		}
	}
	// wait database to be ready.
	if err = rds.Wait(ctx, rdsDrv); err != nil {
		return fmt.Errorf("error waiting database ready: %w", err)
	}
	r.setDataSourceDriver(rdsDrv)

	if r.EnableAuthn {
		// enable authentication.
		if r.CasdoorServer == "" {
			// if not specified, launch embedded casdoor,
			var e casdoor.Embedded
			g.Go(func() error {
				log.Info("running embedded casdoor")
				var err = e.Run(ctx, r.DataSourceAddress)
				if err != nil {
					log.Errorf("error running embedded casdoor: %v", err)
				}
				return err
			})
			// and get embedded casdoor address.
			r.CasdoorServer, err = e.GetAddress(ctx)
			if err != nil {
				return fmt.Errorf("error getting embedded casdor: %w", err)
			}
		}
		// wait casdoor to be ready.
		if err = casdoor.Wait(ctx, r.CasdoorServer); err != nil {
			return fmt.Errorf("error waiting casdoor ready: %w", err)
		}
	}

	// initialize some resources.
	log.Info("initializing")
	var modelClient = getModelClient(rdsDrvDialect, rdsDrv)
	var initOpts = initOptions{
		ModelClient: modelClient,
	}
	if err = r.init(ctx, initOpts); err != nil {
		log.Errorf("error initializing: %v", err)
		return fmt.Errorf("error initializing: %w", err)
	}

	// setup k8s controllers.
	var setupK8sCtrlsOpts = setupK8sCtrlsOptions{
		K8sConfig:   k8sCfg,
		ModelClient: modelClient,
	}
	g.Go(func() error {
		log.Info("setting up kubernetes controller")
		var err = r.setupK8sCtrls(ctx, setupK8sCtrlsOpts)
		if err != nil {
			log.Errorf("error setting up kubernetes controller: %v", err)
		}
		return err
	})

	// setup apis.
	var setupApisOpts = setupApisOptions{
		ModelClient: modelClient,
	}
	g.Go(func() error {
		log.Info("setting up apis")
		var err = r.setupApis(ctx, setupApisOpts)
		if err != nil {
			log.Errorf("error setting up apis: %v", err)
		}
		return err
	})

	return g.Wait()
}

func (r *Server) setKubernetesConfig(cfg *rest.Config) {
	cfg.Timeout = r.KubeConnTimeout
	cfg.QPS = float32(r.KubeConnQPS)
	cfg.Burst = r.KubeConnBurst
}

func (r *Server) setDataSourceDriver(drv *sql.DB) {
	drv.SetConnMaxLifetime(r.DataSourceConnMaxLife)
	drv.SetMaxIdleConns(r.DataSourceConnMaxIdle)
	drv.SetMaxOpenConns(r.DataSourceConnMaxOpen)
}

func getModelClient(drvDialect string, drv *sql.DB) *model.Client {
	var logger = log.WithName("model")
	return model.NewClient(
		model.Log(func(args ...interface{}) { logger.Debug(args...) }),
		model.Driver(entsql.NewDriver(drvDialect, entsql.Conn{ExecQuerier: drv})),
	)
}
