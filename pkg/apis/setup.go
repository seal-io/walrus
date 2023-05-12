package apis

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/apis/account"
	"github.com/seal-io/seal/pkg/apis/application"
	"github.com/seal-io/seal/pkg/apis/applicationinstance"
	"github.com/seal-io/seal/pkg/apis/applicationresource"
	"github.com/seal-io/seal/pkg/apis/applicationrevision"
	"github.com/seal-io/seal/pkg/apis/auth"
	"github.com/seal-io/seal/pkg/apis/connector"
	"github.com/seal-io/seal/pkg/apis/cost"
	"github.com/seal-io/seal/pkg/apis/dashboard"
	"github.com/seal-io/seal/pkg/apis/debug"
	"github.com/seal-io/seal/pkg/apis/environment"
	"github.com/seal-io/seal/pkg/apis/group"
	"github.com/seal-io/seal/pkg/apis/health"
	"github.com/seal-io/seal/pkg/apis/module"
	"github.com/seal-io/seal/pkg/apis/modulecompletion"
	"github.com/seal-io/seal/pkg/apis/moduleversion"
	"github.com/seal-io/seal/pkg/apis/openapi"
	"github.com/seal-io/seal/pkg/apis/perspective"
	"github.com/seal-io/seal/pkg/apis/project"
	"github.com/seal-io/seal/pkg/apis/role"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/apis/secret"
	"github.com/seal-io/seal/pkg/apis/setting"
	"github.com/seal-io/seal/pkg/apis/swagger"
	"github.com/seal-io/seal/pkg/apis/token"
	"github.com/seal-io/seal/pkg/apis/ui"
	"github.com/seal-io/seal/pkg/apis/user"
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
	var apis = gin.New()
	var auths = auth.Auth(opts.EnableAuthn, opts.ModelClient)
	var throttler = runtime.RequestThrottling(opts.ConnQPS, opts.ConnBurst)
	var rectifier = runtime.RequestShaping(opts.ConnQPS, opts.ConnQPS, 5*time.Second)
	var wsCounter = runtime.If(
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
		runtime.Logging(
			"/",
			"/assets/*any",
			"/verify-auth",
			"/livez",
			"/openapi",
			"/swagger/*any",
			"/debug/version"),
		runtime.Recovering(),
		runtime.Erroring(),
		runtime.I18n(),
	)

	runtime.MustRouteGet(apis, "/livez", health.Livez())

	var accountApis = apis.Group("/account",
		rectifier,
		auths)
	{
		var r = accountApis
		runtime.MustRoutePost(r, "/login", account.Login())
		runtime.MustRoutePost(r, "/logout", account.Logout())
		runtime.MustRoutePost(r, "/info", account.Info(opts.ModelClient))
		runtime.MustRouteGet(r, "/info", account.Info(opts.ModelClient))
	}

	var resourceApis = apis.Group("/v1",
		throttler,
		auths)
	{
		var r = auth.WithResourceRoleGenerator(ctx, resourceApis, opts.ModelClient)
		runtime.MustRouteResource(r, application.Handle(opts.ModelClient, opts.K8sConfig, opts.TlsCertified))
		runtime.MustRouteResource(r, applicationinstance.Handle(opts.ModelClient, opts.K8sConfig, opts.TlsCertified))
		runtime.MustRouteResource(r.Group("", rectifier, wsCounter), applicationresource.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, applicationrevision.Handle(opts.ModelClient, opts.K8sConfig, opts.TlsCertified))
		runtime.MustRouteResource(r, connector.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, cost.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, group.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, project.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, role.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, secret.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, setting.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, token.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, user.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, module.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, moduleversion.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, perspective.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, environment.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, dashboard.Handle(opts.ModelClient))
		runtime.MustRouteResource(r, modulecompletion.Handle(opts.ModelClient))
	}
	runtime.MustRouteGet(apis, "/openapi", openapi.Index(opts.EnableAuthn, resourceApis.BasePath()))
	runtime.MustRouteStatic(apis, "/swagger/*any", swagger.Index("/openapi"))

	var debugApis = apis.Group("/debug")
	{
		var r = debugApis
		runtime.MustRouteGet(r, "/version", debug.Version())
		runtime.MustRouteGet(r.Group("", runtime.OnlyLocalIP()), "/pprof/*any", debug.PProf())
		runtime.MustRoutePut(r.Group("", runtime.OnlyLocalIP()), "/flags", debug.SetFlags())
		runtime.MustRouteGet(r, "/flags", debug.GetFlags())
	}

	return apis, nil
}
