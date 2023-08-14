package runtime

import (
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func init() {
	// Disable gin default binding.
	binding.Validator = nil
}

type (
	IHandler any

	IResourceHandler interface {
		IHandler

		Kind() string
	}

	IRouter interface {
		http.Handler

		// Use attaches a global middleware to the router.
		Use(...IHandler) IRouter
		// Group creates a new router group with the given string,
		// and returns the new router group.
		Group(string) IRouter
		// GroupIn creates a new router group with the given string,
		// but returns the original router group,
		// the new router group is passed to the given function.
		GroupIn(string, func(groupRouter IRouter)) IRouter
		// GroupRelativePath returns the relative path of router group.
		GroupRelativePath() string

		// Static registers GET/HEAD routers to serve the handler.
		Static(string, http.FileSystem) IRouter
		// Get registers GET router to serve the handler.
		Get(string, IHandler) IRouter
		// Post registers POST router to serve the handler.
		Post(string, IHandler) IRouter
		// Delete registers DELETE router to serve the handler.
		Delete(string, IHandler) IRouter
		// Patch registers PATCH router to serve the handler.
		Patch(string, IHandler) IRouter
		// Put registers PUT router to serve the handler.
		Put(string, IHandler) IRouter

		// Routes registers the reflected routes of a IHandler.
		//
		// Routes reflects the function descriptors as the below rules,
		// if the handler implements IResourceHandler as well.
		//
		//	Input : struct type.
		//	Output: any types.
		//
		//	* Basic APIs
		//
		//	func Create(<Input>) (<Output>, error)
		//	 ->   POST /<plural>
		//	func Get(<Input>) (<Output>, error)
		//	 ->    GET /<plural>/:id(?watch=true)
		//	func Update(<Input>) error
		//	 ->    PUT /<plural>/:id
		//	func Delete(<Input>) error
		//	 -> DELETE /<plural>/:id
		//	func CollectionCreate(<Input>) (<Output>, error)
		//	 ->   POST /<plural>/_/batch
		//	func CollectionGet(<Input>) (<Output>, (int,) error)
		//	 ->    GET /<plural>(?watch=true)
		//	func CollectionUpdate(<Input>) error
		//	 ->    PUT /<plural>
		//	func CollectionDelete(<Input>) error
		//	 -> DELETE /<plural>
		//
		//	* Extensional APIs
		//
		//	func Route<Something>(<Input(route:method=subpath)>) ((<Output>), (int,) error)
		//	 -> method /<plural>/:id/<subpath>(?watch=true)
		//	func CollectionRoute<Something>(<Input(route:method=subpath)>) ((<Output>), (int,) error)
		//	 -> method /<plural>/_/<subpath>(?watch=true)
		//
		// Otherwise, Routes tries to reflect the function descriptors as the below rules.
		//
		//	Input : struct type.
		//	Output: any types.
		//
		//	func <Anything>(<Input(route:method=path)>) ((<Output>), (int,) error)
		//	 -> method /<path>(?watch=true)
		//
		Routes(IHandler) IRouter
	}
)

type (
	RouterOption interface {
		isOption()
	}

	RouterOptions []RouterOption

	Router struct {
		options RouterOptions
		engine  *gin.Engine
		router  gin.IRouter

		adviceProviders []RouteAdviceProvider
		authorizer      RouteAuthorizer
	}
)

// Apply applies the options one by one,
// and returns the residual options.
func (opts RouterOptions) Apply(fn func(o RouterOption) bool) (rOpts RouterOptions) {
	for i := range opts {
		if opts[i] == nil {
			continue
		}

		if !fn(opts[i]) {
			rOpts = append(rOpts, opts[i])
		}
	}

	return
}

func NewRouter(options ...RouterOption) IRouter {
	opts := RouterOptions(options)

	// Apply global options.
	opts = opts.Apply(func(o RouterOption) bool {
		op, ok := o.(ginGlobalOption)
		if ok {
			op()
		}

		return ok
	})

	e := gin.New()
	e.NoMethod(noMethod)
	e.NoRoute(noRoute)

	// Apply engine options.
	opts = opts.Apply(func(o RouterOption) bool {
		op, ok := o.(ginEngineOption)
		if ok {
			op(e)
		}

		return ok
	})

	rt := &Router{
		options: opts,
		engine:  e,
		router:  e,
	}

	// Apply router options.
	rt.options = opts.Apply(func(o RouterOption) bool {
		op, ok := o.(routerOption)
		if ok {
			op(rt)
		}

		return ok
	})

	e.Use(observing, recovering, erroring)

	// Apply route options.
	rt.options = rt.options.Apply(func(o RouterOption) bool {
		op, ok := o.(ginRouteOption)
		if ok {
			op(rt.engine)
		}

		return ok
	})

	return rt
}

func (rt *Router) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	rt.engine.ServeHTTP(resp, req)
}

func (rt *Router) Use(handlers ...IHandler) IRouter {
	hs := make([]gin.HandlerFunc, 0, len(handlers))
	for i := range handlers {
		hs = append(hs, asHandle(handlers[i]))
	}

	rt.router.Use(hs...)

	return rt
}

func (rt *Router) Group(relativePath string) IRouter {
	grt := *rt
	grt.router = rt.router.Group(relativePath)

	return &grt
}

func (rt *Router) GroupIn(relativePath string, doGroupRoute func(IRouter)) IRouter {
	grt := *rt
	grt.router = rt.router.Group(relativePath)

	if doGroupRoute != nil {
		doGroupRoute(&grt)
	}

	return rt
}

func (rt *Router) GroupRelativePath() string {
	if t, ok := rt.router.(interface{ BasePath() string }); ok {
		return t.BasePath()
	}

	return "/"
}

func (rt *Router) Static(p string, fs http.FileSystem) IRouter {
	skipLoggingPath(path.Join(p, "/*filepath"))
	rt.engine.StaticFS(p, fs)

	return rt
}

func (rt *Router) Get(path string, handler IHandler) IRouter {
	rt.router.GET(path, asHandle(handler))
	return rt
}

func (rt *Router) Post(path string, handler IHandler) IRouter {
	rt.router.POST(path, asHandle(handler))
	return rt
}

func (rt *Router) Delete(path string, handler IHandler) IRouter {
	rt.router.DELETE(path, asHandle(handler))
	return rt
}

func (rt *Router) Patch(path string, handler IHandler) IRouter {
	rt.router.PATCH(path, asHandle(handler))
	return rt
}

func (rt *Router) Put(path string, handler IHandler) IRouter {
	rt.router.PUT(path, asHandle(handler))
	return rt
}
