package apis

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/apis/cli"
	"github.com/seal-io/seal/pkg/apis/connector"
	"github.com/seal-io/seal/pkg/apis/cost"
	"github.com/seal-io/seal/pkg/apis/dashboard"
	"github.com/seal-io/seal/pkg/apis/debug"
	"github.com/seal-io/seal/pkg/apis/environment"
	"github.com/seal-io/seal/pkg/apis/measure"
	"github.com/seal-io/seal/pkg/apis/openapi"
	"github.com/seal-io/seal/pkg/apis/perspective"
	"github.com/seal-io/seal/pkg/apis/project"
	"github.com/seal-io/seal/pkg/apis/role"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/apis/secret"
	"github.com/seal-io/seal/pkg/apis/service"
	"github.com/seal-io/seal/pkg/apis/serviceresource"
	"github.com/seal-io/seal/pkg/apis/servicerevision"
	"github.com/seal-io/seal/pkg/apis/setting"
	"github.com/seal-io/seal/pkg/apis/subject"
	"github.com/seal-io/seal/pkg/apis/subjectrole"
	"github.com/seal-io/seal/pkg/apis/swagger"
	"github.com/seal-io/seal/pkg/apis/template"
	"github.com/seal-io/seal/pkg/apis/templatecompletion"
	"github.com/seal-io/seal/pkg/apis/templateversion"
	"github.com/seal-io/seal/pkg/apis/token"
	"github.com/seal-io/seal/pkg/apis/ui"
	"github.com/seal-io/seal/pkg/apis/variable"
	"github.com/seal-io/seal/pkg/auths"
	"github.com/seal-io/seal/pkg/dao/model"
)

type SetupOptions struct {
	// Configure from launching.
	EnableAuthn bool
	ConnQPS     int
	ConnBurst   int
	// Derived from configuration.
	K8sConfig    *rest.Config
	ModelClient  *model.Client
	TlsCertified bool
}

func (s *Server) Setup(ctx context.Context, opts SetupOptions) (http.Handler, error) {
	gin.DefaultWriter = s.logger
	gin.DefaultErrorWriter = s.logger
	apis := gin.New()
	account := auths.RequestAccount(opts.ModelClient, opts.EnableAuthn)
	throttler := runtime.RequestThrottling(opts.ConnQPS, opts.ConnBurst)
	rectifier := runtime.RequestShaping(opts.ConnQPS, opts.ConnQPS, 5*time.Second)
	wsCounter := runtime.If(
		// Validate websocket connection.
		runtime.IsBidiStreamRequest,
		// Maximum 10 connection per ip.
		runtime.PerIP(func() runtime.Handle {
			return runtime.RequestCounting(10, 5*time.Second)
		}),
	)

	apis.NoMethod(runtime.NoMethod())
	apis.NoRoute(ui.Index(ctx, opts.ModelClient), runtime.NotFound())
	apis.Use(
		runtime.Observing(
			"/",
			"/assets/*any",
			"/verify-auth",
			"/livez",
			"/metrics",
			"/openapi",
			"/swagger/*any",
			"/debug/version"),
		runtime.Recovering(),
		runtime.Erroring(),
		runtime.I18n(),
	)

	runtime.MustRouteGet(apis, "/cli", cli.Index())

	measureApis := apis.Group("",
		throttler)
	{
		r := measureApis
		runtime.MustRouteGet(r, "/readyz", measure.Readyz())
		runtime.MustRouteGet(r, "/livez", measure.Livez())
		runtime.MustRouteGet(r, "/metrics", measure.Metrics())
	}

	accountApis := apis.Group("/account",
		rectifier,
		account.Filter)
	{
		r := accountApis
		runtime.MustRoutePost(r, "/login", account.Login)
		runtime.MustRoutePost(r, "/logout", account.Logout)
		runtime.MustRoutePost(r, "/info", account.UpdateInfo)
		runtime.MustRouteGet(r, "/info", account.GetInfo)
	}

	resourceApis := apis.Group("/v1",
		throttler,
		account.Filter)
	{
		r := resourceApis
		runtime.MustRouteResource(r, service.Handle(opts.ModelClient, opts.K8sConfig, opts.TlsCertified))
		runtime.MustRouteResource(r.Group("", rectifier, wsCounter),
			serviceresource.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, servicerevision.Handle(opts.ModelClient, opts.K8sConfig, opts.TlsCertified))
		runtime.MustRouteResource(r, connector.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, cost.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, dashboard.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, environment.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, template.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, templatecompletion.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, templateversion.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, perspective.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, project.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, role.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, secret.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, setting.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, subject.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, subjectrole.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, token.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, variable.Handle(opts.ModelClient))
	}
	runtime.MustRouteGet(apis, "/openapi", openapi.Index(opts.EnableAuthn, resourceApis.BasePath()))
	runtime.MustRouteStatic(apis, "/swagger/*any", swagger.Index("/openapi"))

	debugApis := apis.Group("/debug")
	{
		r := debugApis
		runtime.MustRouteGet(r, "/version", debug.Version())
		runtime.MustRouteGet(r.Group("", runtime.OnlyLocalIP()), "/pprof/*any", debug.PProf())
		runtime.MustRoutePut(r.Group("", runtime.OnlyLocalIP()), "/flags", debug.SetFlags())
		runtime.MustRouteGet(r, "/flags", debug.GetFlags())
	}

	return apis, nil
}
