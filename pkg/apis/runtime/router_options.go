package runtime

import (
	"io"

	"github.com/gin-gonic/gin"
)

// ginGlobalOption is the function type of RouterOption,
// which is used to set global options for gin.
type ginGlobalOption func()

func (ginGlobalOption) isOption() {}

// ginEngineOption is the function type of RouterOption,
// which is used to set engine for gin.
type ginEngineOption func(*gin.Engine)

func (ginEngineOption) isOption() {}

// routerOption is the function type of RouterOption,
// which is used to set Router.
type routerOption func(*Router)

func (routerOption) isOption() {}

// ginRouteOption is the function type of RouterOption,
// which is used to register the routes within raw gin method.
type ginRouteOption func(gin.IRouter)

func (ginRouteOption) isOption() {}

// WithDefaultWriter is a RouterOption to configure the default writer for gin.
func WithDefaultWriter(w io.Writer) RouterOption {
	return ginGlobalOption(func() {
		gin.DefaultWriter = w
		gin.DefaultErrorWriter = w
	})
}

// WithDefaultHandler is a RouterOption to configure the default handler for gin.
func WithDefaultHandler(handler IHandler) RouterOption {
	return ginEngineOption(func(eng *gin.Engine) {
		eng.NoRoute(
			asHandle(handler),
			noRoute)
	})
}
