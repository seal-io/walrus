package runtime

import (
	"errors"
	"io"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"

	"github.com/seal-io/seal/utils/strs"
)

func init() {
	// disable gin default binding
	binding.Validator = nil
}

// MustRoutePost likes RoutePost,
// but panics if error occurs.
func MustRoutePost(r gin.IRoutes, p string, h any) {
	var err = RoutePost(r, p, h)
	if err != nil {
		panic(err)
	}
}

// RoutePost registers the given handler as POST router.
func RoutePost(r gin.IRoutes, p string, h any) error {
	switch t := h.(type) {
	default:
		return errors.New("invalid handler")
	case ErrorHandler:
		r.POST(p, WrapErrorHandler(t))
	case Handler:
		r.POST(p, WrapHandler(t))
	case ErrorHandle:
		r.POST(p, WrapErrorHandle(t))
	case Handle:
		r.POST(p, t)
	case HTTPHandle:
		r.POST(p, WrapHTTPHandle(t))
	case HTTPHandler:
		r.POST(p, WrapHTTPHandler(t))
	}
	return nil
}

// MustRouteDelete likes RouteDelete,
// but panics if error occurs.
func MustRouteDelete(r gin.IRoutes, p string, h any) {
	var err = RouteDelete(r, p, h)
	if err != nil {
		panic(err)
	}
}

// RouteDelete registers the given handler as DELETE router.
func RouteDelete(r gin.IRoutes, p string, h any) error {
	switch t := h.(type) {
	default:
		return errors.New("invalid handler")
	case ErrorHandler:
		r.DELETE(p, WrapErrorHandler(t))
	case Handler:
		r.DELETE(p, WrapHandler(t))
	case ErrorHandle:
		r.DELETE(p, WrapErrorHandle(t))
	case Handle:
		r.DELETE(p, t)
	case HTTPHandle:
		r.DELETE(p, WrapHTTPHandle(t))
	case HTTPHandler:
		r.DELETE(p, WrapHTTPHandler(t))
	}
	return nil
}

// MustRoutePut likes RoutePut,
// but panics if error occurs.
func MustRoutePut(r gin.IRoutes, p string, h any) {
	var err = RoutePut(r, p, h)
	if err != nil {
		panic(err)
	}
}

// RoutePut registers the given handler as PUT router.
func RoutePut(r gin.IRoutes, p string, h any) error {
	switch t := h.(type) {
	default:
		return errors.New("invalid handler")
	case ErrorHandler:
		r.PUT(p, WrapErrorHandler(t))
	case Handler:
		r.PUT(p, WrapHandler(t))
	case ErrorHandle:
		r.PUT(p, WrapErrorHandle(t))
	case Handle:
		r.PUT(p, t)
	case HTTPHandle:
		r.PUT(p, WrapHTTPHandle(t))
	case HTTPHandler:
		r.PUT(p, WrapHTTPHandler(t))
	}
	return nil
}

// MustRouteGet likes RouteGet,
// but panics if error occurs.
func MustRouteGet(r gin.IRoutes, p string, h any) {
	var err = RouteGet(r, p, h)
	if err != nil {
		panic(err)
	}
}

// RouteGet registers the given handler as GET router.
func RouteGet(r gin.IRoutes, p string, h any) error {
	switch t := h.(type) {
	default:
		return errors.New("invalid handler")
	case ErrorHandler:
		r.GET(p, WrapErrorHandler(t))
	case Handler:
		r.GET(p, WrapHandler(t))
	case ErrorHandle:
		r.GET(p, WrapErrorHandle(t))
	case Handle:
		r.GET(p, t)
	case HTTPHandle:
		r.GET(p, WrapHTTPHandle(t))
	case HTTPHandler:
		r.GET(p, WrapHTTPHandler(t))
	}
	return nil
}

// MustRouteHead likes RouteHead,
// but panics if error occurs.
func MustRouteHead(r gin.IRoutes, p string, h any) {
	var err = RouteHead(r, p, h)
	if err != nil {
		panic(err)
	}
}

// RouteHead registers the given handler as HEAD router.
func RouteHead(r gin.IRoutes, p string, h any) error {
	switch t := h.(type) {
	default:
		return errors.New("invalid handler")
	case ErrorHandler:
		r.HEAD(p, WrapErrorHandler(t))
	case Handler:
		r.HEAD(p, WrapHandler(t))
	case ErrorHandle:
		r.HEAD(p, WrapErrorHandle(t))
	case Handle:
		r.HEAD(p, t)
	case HTTPHandle:
		r.HEAD(p, WrapHTTPHandle(t))
	case HTTPHandler:
		r.HEAD(p, WrapHTTPHandler(t))
	}
	return nil
}

// MustRouteStatic likes RouteStatic,
// but panics if error occurs.
func MustRouteStatic(r gin.IRoutes, p string, h any) {
	var err = RouteStatic(r, p, h)
	if err != nil {
		panic(err)
	}
}

// RouteStatic registers the given handler as HEAD and GET router.
func RouteStatic(r gin.IRoutes, p string, h any) error {
	var err = RouteHead(r, p, h)
	if err != nil {
		return err
	}
	err = RouteGet(r, p, h)
	if err != nil {
		return err
	}
	return nil
}

// MustRouteResource likes RouteResource,
// but panics if error occurs.
func MustRouteResource(r gin.IRoutes, h Resource) {
	var err = RouteResource(r, h)
	if err != nil {
		panic(err)
	}
}

// RouteResource reflects the function descriptors of the given handler,
// and registers as router if satisfy the below rules.
//
//	Input : struct type.
//	Output: any types.
//
//	* Basic APIs
//
//	func Create(*gin.Context, <Input>) ((<Output>,) error)
//	 ->   POST /<plural>
//	func Delete(*gin.Context, <Input>) ((<Output>,) error)
//	 -> DELETE /<plural>/:id
//	func Update(*gin.Context, <Input>) ((<Output>,) error)
//	 ->    PUT /<plural>/:id
//	func Get(*gin.Context, <Input>) (<Output>, error)
//	 ->    GET /<plural>/:id
//
//	* Batch APIs
//
//	func CollectionCreate(*gin.Context, <Input>) ((<Output>,) error)
//	 ->   POST /<plural>/_/batch
//	func CollectionDelete(*gin.Context, <Input>) ((<Output>,) error)
//	 -> DELETE /<plural>
//	func CollectionUpdate(*gin.Context, <Input>) ((<Output>,) error)
//	 ->    PUT /<plural>
//	func CollectionGet(*gin.Context, <Input>) (<Output>, (int,) error)
//	 ->    GET /<plural>
//
//	* Extensional APIs
//
//	func Create<Something>(*gin.Context, <Input>) ((<Output>,) error)
//	 ->   POST /<plural>/:id/<something>
//	func CollectionCreate<Something>(*gin.Context, <Input>) ((<Output>,) error)
//	 ->   POST /<plural>/_/<something>
//	func Delete<Something>(*gin.Context, <Input>) ((<Output>,) error)
//	 -> DELETE /<plural>/:id/<something>
//	func CollectionDelete<Something>(*gin.Context, <Input>) ((<Output>,) error)
//	 -> DELETE /<plural>/_/<something>
//	func Update<Something>(*gin.Context, <Input>) ((<Output>,) error)
//	 ->    PUT /<plural>/:id/<something>
//	func CollectionUpdate<Something>(*gin.Context, <Input>) ((<Output>,) error)
//	 ->    PUT /<plural>/_/<something>
//	func Get<Something>(*gin.Context, <Input>) (<Output>, (int,) error)
//	 ->    GET /<plural>/:id/<something>
//	func CollectionGet<Something>(*gin.Context, <Input>) (<Output>, (int,) error)
//	 ->    GET /<plural>/_/<something>
//	func Route<Something>(*gin.Context, <Input(route:method=subpath)>) ((<Output>), (int,) error)
//	 -> method /<plural>/:id/<subpath>
//	func CollectionRoute<Something>(*gin.Context, <Input(route:method=subpath)>) ((<Output>), (int,) error)
//	 -> method /<plural>/_/<subpath>
//
//	* Stream APIs
//
//	func Stream(RequestStream) error
//	  ->   websocket /<plural>
func RouteResource(r gin.IRoutes, h Resource) error {
	if adv, ok := r.(AdviceBeforeResourceRegistering); ok {
		var err = adv.BeforeAdvice(adviceResource{Resource: h})
		if err != nil {
			return err
		}
	}

	if adv, ok := r.(AdviceAfterResourceRegistering); ok {
		var err = adv.AfterAdvice(adviceResource{Resource: h})
		if err != nil {
			return err
		}
	}

	return nil
}

type rendCloser interface {
	io.Closer
	render.Render
}

func isImplementationOfError(r reflect.Type) bool {
	var expected = reflect.TypeOf((*error)(nil)).Elem()
	switch r.Kind() {
	default:
		return false
	case reflect.Interface:
		return r.Implements(expected)
	}
}

func isTypeOfRequestStream(r reflect.Type) bool {
	var expected = reflect.TypeOf(RequestStream{}).String()
	switch r.Kind() {
	default:
		return false
	case reflect.Struct:
		var given = r.String()
		return given == expected
	}
}

func isTypeOfGinContext(r reflect.Type) bool {
	var expected = reflect.TypeOf(gin.Context{}).String()
	switch r.Kind() {
	default:
		return false
	case reflect.Pointer:
		var given = r.Elem().String()
		return given == expected
	}
}

func isTypeOfInt(r reflect.Type) bool {
	return r.Kind() == reflect.Int
}

func getResourceAndResourcePath(k string) (rk, rp string) {
	rk = strs.CamelizeDownFirst(strs.Pluralize(k))
	rp = strings.ToLower(strs.Pluralize(strs.Dasherize(rk)))
	return
}

type adviceResource struct {
	Resource
}

func (r adviceResource) ResourceAndResourcePath() (resource, resourcePath string) {
	return getResourceAndResourcePath(r.Kind())
}

func (r adviceResource) Unwrap() Resource {
	return r.Resource
}
