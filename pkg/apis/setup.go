package apis

import (
	"context"
	"net/http"
	"time"

	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/apis/catalog"
	"github.com/seal-io/walrus/pkg/apis/cli"
	"github.com/seal-io/walrus/pkg/apis/connector"
	"github.com/seal-io/walrus/pkg/apis/cost"
	"github.com/seal-io/walrus/pkg/apis/dashboard"
	"github.com/seal-io/walrus/pkg/apis/debug"
	"github.com/seal-io/walrus/pkg/apis/measure"
	"github.com/seal-io/walrus/pkg/apis/perspective"
	"github.com/seal-io/walrus/pkg/apis/project"
	"github.com/seal-io/walrus/pkg/apis/proxy"
	"github.com/seal-io/walrus/pkg/apis/resourcedefinition"
	"github.com/seal-io/walrus/pkg/apis/role"
	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/apis/setting"
	"github.com/seal-io/walrus/pkg/apis/subject"
	"github.com/seal-io/walrus/pkg/apis/template"
	"github.com/seal-io/walrus/pkg/apis/templatecompletion"
	"github.com/seal-io/walrus/pkg/apis/templateversion"
	"github.com/seal-io/walrus/pkg/apis/ui"
	"github.com/seal-io/walrus/pkg/apis/variable"
	"github.com/seal-io/walrus/pkg/auths"
	"github.com/seal-io/walrus/pkg/dao/model"
	pkgworkflow "github.com/seal-io/walrus/pkg/workflow"
)

type SetupOptions struct {
	// Configure from launching.
	EnableAuthn           bool
	ConnQPS               int
	ConnBurst             int
	WebsocketConnMaxPerIP int
	// Derived from configuration.
	K8sConfig    *rest.Config
	ModelClient  *model.Client
	TlsCertified bool
}

func (s *Server) Setup(ctx context.Context, opts SetupOptions) (http.Handler, error) {
	// Prepare middlewares.
	account := auths.RequestAccount(opts.ModelClient, opts.EnableAuthn)
	throttler := runtime.RequestThrottling(opts.ConnQPS, opts.ConnBurst)
	rectifier := runtime.RequestShaping(opts.ConnQPS, opts.ConnQPS, 5*time.Second)
	wsCounter := runtime.If(
		// Validate websocket connection.
		runtime.IsBidiStreamRequest,
		// Maximum 10 connection per ip.
		runtime.PerIP(func() runtime.Handle {
			return runtime.RequestCounting(opts.WebsocketConnMaxPerIP, 5*time.Second)
		}),
	)
	i18n := runtime.I18n()

	wc, err := pkgworkflow.NewArgoWorkflowClient(opts.ModelClient, opts.K8sConfig)
	if err != nil {
		return nil, err
	}

	// Initial router.
	apisOpts := []runtime.RouterOption{
		runtime.WithDefaultWriter(s.logger),
		runtime.WithDefaultHandler(ui.Index(ctx, opts.ModelClient)),
		runtime.SkipLoggingPaths(
			"/",
			"/cli",
			"/assets/*filepath",
			"/readyz",
			"/livez",
			"/metrics",
			"/debug/version"),
		runtime.ExposeOpenAPI(),
		runtime.WithRouteAdviceProviders(provideModelClient(opts.ModelClient)),
		runtime.WithResourceAuthorizer(account),
	}

	apis := runtime.NewRouter(apisOpts...).
		Use(i18n)

	accountApis := apis.Group("/account").
		Use(rectifier, account.Filter)
	{
		r := accountApis
		r.Routes(account)
	}

	resourceApis := apis.Group("/v1").
		Use(throttler, wsCounter, account.Filter)
	{
		r := resourceApis
		r.Routes(catalog.Handle(opts.ModelClient))
		r.Routes(connector.Handle(opts.ModelClient))
		r.Routes(cost.Handle(opts.ModelClient))
		r.Routes(dashboard.Handle(opts.ModelClient))
		r.Routes(perspective.Handle(opts.ModelClient))
		r.Routes(project.Handle(opts.ModelClient, opts.K8sConfig, wc))
		r.Routes(resourcedefinition.Handle(opts.ModelClient))
		r.Routes(role.Handle(opts.ModelClient))
		r.Routes(setting.Handle(opts.ModelClient))
		r.Routes(subject.Handle(opts.ModelClient))
		r.Routes(template.Handle(opts.ModelClient))
		r.Routes(templateversion.Handle(opts.ModelClient))
		r.Routes(templatecompletion.Handle(opts.ModelClient))
		r.Routes(variable.Handle(opts.ModelClient))
	}

	cliApis := apis.Group("")
	{
		r := cliApis
		r.Get("/cli", cli.Index())
	}

	measureApis := apis.Group("").
		Use(throttler)
	{
		r := measureApis
		r.Get("/readyz", measure.Readyz())
		r.Get("/livez", measure.Livez())
		r.Get("/metrics", measure.Metrics())
	}

	proxyApis := apis.Group("/proxy").
		Use(throttler, account.Filter)
	{
		r := proxyApis
		r.Get("/*path", proxy.Proxy(opts.ModelClient))
	}

	debugApis := apis.Group("/debug")
	{
		r := debugApis
		r.Get("/version", debug.Version())
		r.Get("/flags", debug.GetFlags())
		r.Group("").
			Use(runtime.OnlyLocalIP()).
			Get("/pprof/*any", debug.PProf()).
			Put("/flags", debug.SetFlags())
	}

	return apis, nil
}
