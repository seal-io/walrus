package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"time"

	certcache "github.com/seal-io/utils/certs/cache"
	"github.com/seal-io/utils/certs/kubecert"
	"github.com/seal-io/utils/osx"
	"github.com/spf13/pflag"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/client-go/rest"
	cliflag "k8s.io/component-base/cli/flag"

	"github.com/seal-io/walrus/pkg/manager"
	"github.com/seal-io/walrus/pkg/servers/serverset/scheme"
	"github.com/seal-io/walrus/pkg/system"
	"github.com/seal-io/walrus/pkg/systemkuberes"
)

type Options struct {
	ManagerOptions *manager.Options

	// Control.
	DisableAuths        bool
	DisableApplications []string

	// Authentication.
	AuthnTokenWebhookCacheTTL time.Duration
	AuthnTokenRequestTimeout  time.Duration

	// Authorization.
	AuthzAllowCacheTTL time.Duration
	AuthzDenyCacheTTL  time.Duration

	// Audit.
	AuditPolicyFile        string
	AuditLogFile           string
	AuditWebhookConfigFile string
}

func NewOptions() *Options {
	mgrOptions := manager.NewOptions()
	mgrOptions.Serve = false

	return &Options{
		ManagerOptions: mgrOptions,

		// Control.
		DisableAuths:        false,
		DisableApplications: []string{},

		// Authentication.
		AuthnTokenWebhookCacheTTL: 10 * time.Second,
		AuthnTokenRequestTimeout:  10 * time.Second,

		// Authorization.
		AuthzAllowCacheTTL: 10 * time.Second,
		AuthzDenyCacheTTL:  10 * time.Second,

		// Audit.
		AuditPolicyFile:        "",
		AuditLogFile:           "",
		AuditWebhookConfigFile: "",
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.ManagerOptions.AddFlags(fs)

	// Control.
	fs.BoolVar(&o.DisableAuths, "disable-auths", o.DisableAuths,
		"disable authentication and authorization.")
	fs.StringSliceVar(&o.DisableApplications, "disable-applications", o.DisableApplications,
		"disable installing applications, select from [minio, hermitcrab].")

	// Authentication.
	fs.DurationVar(&o.AuthnTokenWebhookCacheTTL, "authentication-token-webhook-cache-ttl",
		o.AuthnTokenWebhookCacheTTL,
		"the duration to cache responses from the webhook token authenticator.")
	fs.DurationVar(&o.AuthnTokenRequestTimeout, "authentication-token-request-timeout",
		o.AuthnTokenRequestTimeout,
		"the duration to wait for a response from the webhook token authenticator.")

	// Authorization.
	fs.DurationVar(&o.AuthzAllowCacheTTL, "authorization-webhook-cache-authorized-ttl",
		o.AuthzAllowCacheTTL,
		"the duration to cache 'authorized' responses from the webhook authorizer.")
	fs.DurationVar(&o.AuthzDenyCacheTTL, "authorization-webhook-cache-unauthorized-ttl",
		o.AuthzDenyCacheTTL,
		"the duration to cache 'unauthorized' responses from the webhook authorizer.")

	// Audit.
	fs.StringVar(&o.AuditPolicyFile, "audit-policy-file", o.AuditPolicyFile,
		"path to the file that defines the audit policy configuration.")
	fs.StringVar(&o.AuditLogFile, "audit-log-path", o.AuditLogFile,
		"if set, all requests coming to the server will be logged to this file. "+
			"'-' means standard out.")
	fs.StringVar(&o.AuditWebhookConfigFile, "audit-webhook-config-file", o.AuditWebhookConfigFile,
		"path to a kubeconfig formatted file that defines the audit webhook configuration.")
}

func (o *Options) Validate(ctx context.Context) error {
	if err := o.ManagerOptions.Validate(ctx); err != nil {
		return err
	}

	if !o.DisableAuths {
		// Authentication.
		if o.AuthnTokenWebhookCacheTTL < 10*time.Second {
			return errors.New("--authentication-token-webhook-cache-ttl: less than 10s")
		}
		if o.AuthnTokenRequestTimeout < 10*time.Second {
			return errors.New("--authentication-token-request-timeout: less than 10s")
		}

		// Authorization.
		if o.AuthzAllowCacheTTL < 10*time.Second {
			return errors.New("--authorization-webhook-cache-authorized-ttl: less than 10s")
		}
		if o.AuthzDenyCacheTTL < 10*time.Second {
			return errors.New("--authorization-webhook-cache-unauthorized-ttl: less than 10s")
		}
	}

	// Audit.
	if o.AuditPolicyFile != "" && !osx.ExistsFile(o.AuditPolicyFile) {
		return errors.New("--audit-policy-file: no found file")
	}
	if o.AuditLogFile != "" && o.AuditLogFile != "-" && !osx.ExistsDir(filepath.Dir(o.AuditLogFile)) {
		return errors.New("--audit-log-path: no found parent directory")
	}
	if o.AuditWebhookConfigFile != "" && !osx.ExistsFile(o.AuditWebhookConfigFile) {
		return errors.New("--audit-webhook-config-file: no found file")
	}

	return nil
}

func (o *Options) Complete(ctx context.Context) (*Config, error) {
	mgrConfig, err := o.ManagerOptions.Complete(ctx)
	if err != nil {
		return nil, err
	}

	system.ConfigureDisallowApplications(o.DisableApplications)

	serve := &genericoptions.SecureServingOptions{
		BindAddress: o.ManagerOptions.BindAddress,
		BindPort:    o.ManagerOptions.BindPort,
		ServerCert: genericoptions.GeneratableKeyCert{
			PairName:      "tls",
			CertDirectory: o.ManagerOptions.CertDir,
		},
		CipherSuites:                 cliflag.PreferredTLSCipherNames(),
		MinTLSVersion:                "VersionTLS12",
		HTTP2MaxStreamsPerConnection: 1000,
	}
	if serve.ServerCert.CertDirectory == "" {
		// Deploy in standalone mode(by Docker run) or laptop development,
		// the loopback Kubernetes cluster is nearby.
		certCache, err := certcache.NewK8sCache(ctx,
			"server", system.LoopbackKubeClient.Get().CoreV1().Secrets(systemkuberes.SystemNamespaceName))
		if err != nil {
			return nil, fmt.Errorf("create cert cache: %w", err)
		}
		certMgr := &kubecert.StaticManager{
			CertCli: system.LoopbackKubeClient.Get().CertificatesV1().CertificateSigningRequests(),
			Cache:   certCache,
			Host:    systemkuberes.SystemRoutingServiceName,
			AlternateIPs: func() []net.IP {
				if system.LoopbackKubeInside.Get() {
					return nil
				}
				return []net.IP{
					net.ParseIP("127.0.0.1"),
					net.ParseIP(system.PrimaryIP.Get()),
				}
			}(),
			AlternateDNSNames: []string{
				fmt.Sprintf("%s.%s.svc", systemkuberes.SystemRoutingServiceName, systemkuberes.SystemNamespaceName),
				fmt.Sprintf("%s.%s", systemkuberes.SystemRoutingServiceName, systemkuberes.SystemNamespaceName),
				systemkuberes.SystemRoutingServiceName,
				"localhost",
			},
		}
		serve.ServerCert.GeneratedCert = certMgr
	}

	var authn *genericoptions.DelegatingAuthenticationOptions
	if !o.DisableAuths {
		authn = &genericoptions.DelegatingAuthenticationOptions{
			CacheTTL:             o.AuthnTokenWebhookCacheTTL,
			TokenRequestTimeout:  o.AuthnTokenRequestTimeout,
			WebhookRetryBackoff:  genericoptions.DefaultAuthWebhookRetryBackoff(),
			RemoteKubeConfigFile: mgrConfig.KubeConfigPath,
			DisableAnonymous:     false,
		}
	}

	var authz *genericoptions.DelegatingAuthorizationOptions
	if !o.DisableAuths {
		authz = &genericoptions.DelegatingAuthorizationOptions{
			AllowCacheTTL:        o.AuthzAllowCacheTTL,
			DenyCacheTTL:         o.AuthzDenyCacheTTL,
			WebhookRetryBackoff:  genericoptions.DefaultAuthWebhookRetryBackoff(),
			RemoteKubeConfigFile: mgrConfig.KubeConfigPath,
			ClientTimeout:        10 * time.Second,
			AlwaysAllowGroups:    []string{"system:masters"},
			AlwaysAllowPaths: []string{
				"/", "/assets/*", "/favicon.ico", // UI assets
				"/mutate-*", "/validate-*", // Webhooks
				"/livez", "/readyz", "/metrics", "/debug/*", // Measure
				"/clis/*",    // CLI binaries
				"/openapi/*", // OpenAPI
				"/swagger/*", // Swagger
			},
		}
	}

	audit := genericoptions.NewAuditOptions()
	audit.PolicyFile = o.AuditPolicyFile
	audit.LogOptions.Path = o.AuditLogFile
	audit.WebhookOptions.ConfigFile = o.AuditWebhookConfigFile

	apiSrvCfg := genericapiserver.NewConfig(scheme.Codecs)
	{
		// Feedback Kubernetes client configuration.
		apiSrvCfg.LoopbackClientConfig = rest.CopyConfig(&mgrConfig.KubeClientConfig)
		// Disable default metrics service.
		apiSrvCfg.EnableMetrics = false
		// Disable default profiling service.
		apiSrvCfg.EnableProfiling = false
		// Disable default index service.
		apiSrvCfg.EnableIndex = false
		// Disable following post start hooks,
		// because the registered apiserver can manage them.
		apiSrvCfg.DisabledPostStartHooks.Insert(
			"priority-and-fairness-filter",
			"max-in-flight-filter",
			"storage-object-count-tracker-hook",
		)
	}

	return &Config{
		ManagerConfig:   mgrConfig,
		APIServerConfig: apiSrvCfg,
		Serve:           serve,
		Authn:           authn,
		Authz:           authz,
		Audit:           audit,
	}, nil
}
