package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/seal-io/utils/version"
	"github.com/seal-io/utils/waitx"
	admreg "k8s.io/api/admissionregistration/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilversion "k8s.io/apimachinery/pkg/version"
	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/utils/ptr"

	"github.com/seal-io/walrus/pkg/apis"
	"github.com/seal-io/walrus/pkg/manager"
	"github.com/seal-io/walrus/pkg/servers/serverset/scheme"
	"github.com/seal-io/walrus/pkg/system"
	"github.com/seal-io/walrus/pkg/systemkuberes"
	"github.com/seal-io/walrus/pkg/webhooks"
)

type Config struct {
	ManagerConfig   *manager.Config
	APIServerConfig *genericapiserver.Config
	Serve           *genericoptions.SecureServingOptions
	Authn           *genericoptions.DelegatingAuthenticationOptions
	Authz           *genericoptions.DelegatingAuthorizationOptions
	Audit           *genericoptions.AuditOptions
}

func (c *Config) Apply(ctx context.Context) (*Server, error) {
	mgr, err := c.ManagerConfig.Apply(ctx)
	if err != nil {
		return nil, err
	}

	apiSrvCfg := c.APIServerConfig

	// Apply the server configuration.
	err = c.Serve.ApplyTo(&apiSrvCfg.SecureServing)
	if err != nil {
		return nil, fmt.Errorf("apply server config: %w", err)
	}

	if c.Authn != nil {
		// Apply the authentication configuration.
		err = c.Authn.ApplyTo(&apiSrvCfg.Authentication, apiSrvCfg.SecureServing, nil)
		if err != nil {
			return nil, fmt.Errorf("apply authentication config: %w", err)
		}
	}

	if c.Authz != nil {
		// Apply the authorization configuration.
		err = c.Authz.ApplyTo(&apiSrvCfg.Authorization)
		if err != nil {
			return nil, fmt.Errorf("apply authorization config: %w", err)
		}
	}

	// Apply the audit configuration.
	err = c.Audit.ApplyTo(apiSrvCfg)
	if err != nil {
		return nil, fmt.Errorf("apply audit config: %w", err)
	}

	// Apply OpenAPI configuration.
	var (
		title        = "Walrus - Open Source XaC Platform"
		fullVersion  = version.Get()
		majorVersion = version.Major()
		minorVersion = strings.TrimPrefix(version.MajorMinor(), majorVersion+".")
	)
	apiSrvCfg.Version = &utilversion.Info{
		Major:      majorVersion,
		Minor:      minorVersion,
		GitVersion: fullVersion,
		GitCommit:  version.GitCommit,
	}
	apiSrvCfg.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(
		apis.GetOpenAPIDefinitions, openapinamer.NewDefinitionNamer(scheme.Scheme))
	apiSrvCfg.OpenAPIConfig.Info.Title = title
	apiSrvCfg.OpenAPIConfig.Info.Version = fullVersion
	apiSrvCfg.OpenAPIV3Config = genericapiserver.DefaultOpenAPIV3Config(
		apis.GetOpenAPIDefinitions, openapinamer.NewDefinitionNamer(scheme.Scheme))
	apiSrvCfg.OpenAPIV3Config.Info.Title = title
	apiSrvCfg.OpenAPIV3Config.Info.Version = fullVersion

	apiSrvCompletedCfg := apiSrvCfg.Complete(nil)
	apiSrv, err := apiSrvCompletedCfg.New("walrus", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, fmt.Errorf("create APIServer: %w", err)
	}

	// By default, we suggest to deploy with HA mode(by all-in-one YAML or Helm Chart),
	// which means we stay inside the loopback Kubernetes cluster.
	// So the system Kubernetes Service is created before webhook server start.
	cc := admreg.WebhookClientConfig{
		Service: &admreg.ServiceReference{
			Namespace: systemkuberes.SystemNamespaceName,
			Name:      systemkuberes.SystemRoutingServiceName,
			Port:      ptr.To(int32(c.Serve.BindPort)),
		},
	}
	// However, if we stand closed to loopback Kubernetes cluster but not inside,
	// we can use the primary IP address to access the webhook server.
	if !system.LoopbackKubeInside.Get() && system.LoopbackKubeNearby.Get() {
		// NB(thxCode): launch multiple instances, only one takes working.
		ep := fmt.Sprintf("https://%s:%d", system.PrimaryIP.Get(), c.Serve.BindPort)
		cc = admreg.WebhookClientConfig{
			URL: ptr.To(ep),
		}
	}
	// When no cert provided, we will use Kubernetes CertificateSigningRequest to generate the server cert,
	// we should load the root CA bundle for webhook configuration.
	if c.Serve.ServerCert.CertDirectory == "" {
		// TODO(thxCode): The root CA might be expired,
		//   we can refresh the CA bundle in the future with restarting.
		//   A restarting-less way is needed in the future.
		err = waitx.PollUntilContextTimeout(ctx, time.Second, 30*time.Second, true,
			func(ctx context.Context) error {
				cm, err := system.LoopbackKubeClient.Get().CoreV1().ConfigMaps(meta.NamespaceSystem).
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
			return nil, fmt.Errorf("get kube-root-ca.crt: %w", err)
		}
	}
	// Install webhook configurations.
	err = webhooks.InstallWebhookConfigurations(ctx, system.LoopbackKubeClient.Get(), cc)
	if err != nil {
		return nil, err
	}

	return &Server{
		Manager:   mgr,
		APIServer: apiSrv,
	}, nil
}
