package runtime

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gin-gonic/gin/render"

	"github.com/seal-io/seal/pkg/apis/auth/session"
	"github.com/seal-io/seal/utils/log"
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
	var logger = log.WithName("restful")

	if adv, ok := r.(AdviceBeforeResourceRegistering); ok {
		var err = adv.BeforeAdvice(adviceResource{Resource: h})
		if err != nil {
			return err
		}
	}

	var k = h.Kind()
	var hr = reflect.ValueOf(h)
	var ht = hr.Type()
	for i := 0; i < ht.NumMethod(); i++ {
		var mn = ht.Method(i).Name
		switch mn {
		case "Create", "Delete", "Update", "Get":
			var err = registerBasicHandler(r, hr, mn, k)
			if err != nil {
				logger.Warnf("error registering basic handler of '%s': %v", hr.Type().String(), err)
				continue
			}
		case "CollectionCreate", "CollectionDelete", "CollectionUpdate", "CollectionGet":
			var err = registerCollectionHandler(r, hr, mn, k)
			if err != nil {
				logger.Warnf("error registering collection handler of '%s': %v", hr.Type().String(), err)
				continue
			}
		default:
			var err = registerExtensionalHandler(r, hr, mn, k)
			if err != nil {
				logger.Warnf("error registering extensional handler of '%s': %v", hr.Type().String(), err)
				continue
			}
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

func registerBasicHandler(r gin.IRoutes, hr reflect.Value, mn string, k string) error {
	var isCreateMethod = mn == "Create"
	var isDeleteMethod = mn == "Delete"
	var isUpdateMethod = mn == "Update"
	var isGetMethod = mn == "Get"
	if !(isCreateMethod || isDeleteMethod || isUpdateMethod || isGetMethod) {
		return nil
	}

	// validate
	var mr = hr.MethodByName(mn)
	var mrt = mr.Type()
	switch mrt.NumIn() {
	default:
		return fmt.Errorf("invalid input number of '%s': got %d but expected 2", mn, mrt.NumIn())
	case 2:
		if !isTypeOfGinContext(mrt.In(0)) {
			return fmt.Errorf("illegal first input type of '%s': expected *gin.Context", mn)
		}
		switch mrt.In(1).Kind() {
		default:
			return fmt.Errorf("illegal last input type of '%s': expected struct or pointer type", mn)
		case reflect.Struct, reflect.Pointer:
		}
	}
	switch mrt.NumOut() {
	default:
		return fmt.Errorf("invalid output number of'%s': got %d", mn, mrt.NumOut())
	case 2:
		if !isImplementationOfError(mrt.Out(1)) {
			return fmt.Errorf("illegal last output type of '%s': expected error", mn)
		}
	case 1:
		if isGetMethod {
			return fmt.Errorf("invalid output number of '%s': got %d but expected at least 2", mn, mrt.NumOut())
		}
		if !isImplementationOfError(mrt.Out(0)) {
			return fmt.Errorf("illegal last output type of '%s': expected error", mn)
		}
	}

	// construct virtual handler
	var resource, resourcePath = getResourceAndResourcePath(k)
	var it = mrt.In(1)
	var ip = GetInputProfile(it)
	var ipState = ip.State()
	var vh = func(c *gin.Context) {
		var s = session.LoadSubject(c)
		if !s.Enforce(c, resource) {
			if s.IsAnonymous() {
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
			return
		}
		session.StoreSubjectCurrentOperation(c, s.Give(resource).If(c.Request.Method))

		// bind input
		var ri reflect.Value
		if it.Kind() == reflect.Pointer {
			ri = reflect.New(it.Elem())
		} else {
			ri = reflect.New(it)
		}
		if ipState[ProfileCategoryHeader] {
			if err := c.BindHeader(ri.Interface()); err != nil {
				return
			}
		}
		if ipState[ProfileCategoryUri] {
			if err := c.BindUri(ri.Interface()); err != nil {
				return
			}
		}
		if ipState[ProfileCategoryQuery] {
			if err := binding.MapFormWithTag(ri.Interface(), c.Request.URL.Query(), "query"); err != nil {
				return
			}
		}
		if c.Request.ContentLength != 0 {
			switch c.ContentType() {
			case binding.MIMEPOSTForm:
				if isCreateMethod {
					if err := c.MustBindWith(ri.Interface(), binding.Form); err != nil {
						return
					}
				}
			case binding.MIMEMultipartPOSTForm:
				if isCreateMethod {
					if err := c.MustBindWith(ri.Interface(), binding.FormMultipart); err != nil {
						return
					}
				}
			default:
				if !isGetMethod {
					if err := c.BindJSON(ri.Interface()); err != nil {
						return
					}
				}
			}
		}
		if rv, ok := ri.Interface().(Validator); ok {
			if err := rv.Validate(); err != nil {
				_ = c.Error(Errorw(http.StatusBadRequest, err)).
					SetType(gin.ErrorTypeBind)
				return
			}
		}
		if rv, ok := ri.Interface().(ValidatorWithInput); ok {
			hv, ok := hr.Interface().(ValidatingInput)
			if ok {
				if err := rv.ValidateWith(c, hv.Validating()); err != nil {
					_ = c.Error(Errorw(http.StatusBadRequest, err)).
						SetType(gin.ErrorTypeBind)
					return
				}
			}
		}
		if it.Kind() != reflect.Pointer {
			ri = ri.Elem()
		}

		// process
		var inputs = make([]reflect.Value, 0, 2)
		inputs = append(inputs, reflect.ValueOf(c))
		inputs = append(inputs, ri)
		var outputs = mr.Call(inputs)

		// render response
		if c.Request.Context().Err() != nil ||
			c.Writer.Size() >= 0 ||
			len(c.Errors) != 0 {
			// already render inside the above processing
			return
		}
		var errInterface = outputs[len(outputs)-1].Interface()
		if errInterface != nil {
			var err = errInterface.(error)
			var ge = c.Error(err)
			if !isGinError(err) {
				_ = ge.SetType(gin.ErrorTypePublic)
			}
			return
		}
		var code = http.StatusOK
		switch len(outputs) {
		case 1:
			if isCreateMethod {
				code = http.StatusNoContent
			}
			c.Writer.WriteHeader(code)
		case 2:
			if outputs[0].IsZero() {
				if isCreateMethod {
					code = http.StatusNoContent
				}
				c.Writer.WriteHeader(code)
				return
			}
			if isCreateMethod {
				code = http.StatusCreated
			}
			var obj = outputs[0].Interface()
			switch ot := obj.(type) {
			case rendCloser:
				if ot == nil {
					return
				}
				defer func() { _ = ot.Close() }()
				c.Render(code, ot)
			case render.Render:
				if ot == nil {
					return
				}
				c.Render(code, ot)
			default:
				c.JSON(code, obj) // TODO negotiate
			}
		}
	}

	// register route
	var rp = path.Join("/", resourcePath)
	var xrp = path.Join(rp, ":id")
	switch {
	case isCreateMethod:
		r.POST(rp, vh)
	case isDeleteMethod:
		r.DELETE(xrp, vh)
	case isUpdateMethod:
		r.PUT(xrp, vh)
	case isGetMethod:
		r.GET(xrp, vh)
	}

	// register schema
	var handle = hr.Type().String() + "." + mn
	var ots []reflect.Type
	switch mrt.NumOut() {
	case 2:
		ots = []reflect.Type{mrt.Out(0)}
	case 3:
		ots = []reflect.Type{mrt.Out(0), mrt.Out(1)}
	}
	var op = GetOutputProfile(ots...)
	switch {
	case isCreateMethod:
		schemeRoute(resource, handle, http.MethodPost, rp, ip, op)
	case isDeleteMethod:
		schemeRoute(resource, handle, http.MethodDelete, xrp, ip, op)
	case isUpdateMethod:
		schemeRoute(resource, handle, http.MethodPut, xrp, ip, op)
	case isGetMethod:
		schemeRoute(resource, handle, http.MethodGet, xrp, ip, op)
	}
	return nil
}

func registerCollectionHandler(r gin.IRoutes, hr reflect.Value, mn string, k string) error {
	var isCreateMethod = mn == "CollectionCreate"
	var isDeleteMethod = mn == "CollectionDelete"
	var isUpdateMethod = mn == "CollectionUpdate"
	var isGetMethod = mn == "CollectionGet"
	if !(isCreateMethod || isDeleteMethod || isUpdateMethod || isGetMethod) {
		return nil
	}

	// validate
	var mr = hr.MethodByName(mn)
	var mrt = mr.Type()
	switch mrt.NumIn() {
	default:
		return fmt.Errorf("invalid input number of '%s': got %d but expected 2", mn, mrt.NumIn())
	case 2:
		if !isTypeOfGinContext(mrt.In(0)) {
			return fmt.Errorf("illegal first input type of '%s': expected *gin.Context", mn)
		}
		switch mrt.In(1).Kind() {
		default:
			return fmt.Errorf("illegal last input type of '%s': expected struct, pointer or slice type", mn)
		case reflect.Struct, reflect.Pointer, reflect.Slice:
		}
	}
	switch mrt.NumOut() {
	default:
		return fmt.Errorf("invalid output number of '%s': got %d but expected not more than 3", mn, mrt.NumOut())
	case 3:
		if !isGetMethod {
			return fmt.Errorf("invalid output number of '%s': got %d but exepcted not more than 2", mn, mrt.NumOut())
		}
		if !isTypeOfInt(mrt.Out(1)) {
			return fmt.Errorf("illegal second output type of '%s': expected int", mn)
		}
		if !isImplementationOfError(mrt.Out(2)) {
			return fmt.Errorf("illegal last output type of '%s': expected error", mn)
		}
	case 2:
		if !isImplementationOfError(mrt.Out(1)) {
			return fmt.Errorf("illegal last output type of '%s': expected error", mn)
		}
	case 1:
		if isGetMethod {
			return fmt.Errorf("invalid output number of '%s': got %d but expected at least 2", mn, mrt.NumOut())
		}
		if !isImplementationOfError(mrt.Out(0)) {
			return fmt.Errorf("illegal last output type of '%s': expected error", mn)
		}
	}

	// get stream handler
	var smr = func() reflect.Value {
		if !isGetMethod {
			return reflect.Value{}
		}
		var smr = hr.MethodByName("Stream")
		if !smr.IsValid() {
			return reflect.Value{}
		}
		var smrt = smr.Type()
		switch smrt.NumIn() {
		default:
			return reflect.Value{}
		case 1:
			if !isTypeOfRequestStream(smrt.In(0)) {
				return reflect.Value{}
			}
		}
		switch smrt.NumOut() {
		default:
			return reflect.Value{}
		case 1:
			if !isImplementationOfError(smrt.Out(0)) {
				return reflect.Value{}
			}
		}
		return smr
	}()

	// construct virtual handler
	var resource, resourcePath = getResourceAndResourcePath(k)
	var it = mrt.In(1)
	var ip = GetInputProfile(it)
	var ipState = ip.State()
	var vh = func(c *gin.Context) {
		var s = session.LoadSubject(c)
		if !s.Enforce(c, resource) {
			if s.IsAnonymous() {
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
			return
		}
		session.StoreSubjectCurrentOperation(c, s.Give(resource).If(c.Request.Method))

		// upgrade to stream
		if isGetMethod && smr.IsValid() && isUpgradeStreamRequest(c) {
			doUpgradeStreamRequest(c, smr)
			return
		}

		// bind input
		var ri = reflect.New(it)
		if it.Kind() == reflect.Pointer {
			ri = reflect.New(it.Elem())
		} else {
			ri = reflect.New(it)
		}
		if ipState[ProfileCategoryHeader] {
			if err := c.BindHeader(ri.Interface()); err != nil {
				return
			}
		}
		if ipState[ProfileCategoryUri] {
			if err := c.BindUri(ri.Interface()); err != nil {
				return
			}
		}
		if ipState[ProfileCategoryQuery] {
			if err := binding.MapFormWithTag(ri.Interface(), c.Request.URL.Query(), "query"); err != nil {
				return
			}
		}
		if c.Request.ContentLength != 0 {
			switch c.ContentType() {
			case binding.MIMEPOSTForm:
				if isCreateMethod {
					if err := c.MustBindWith(ri.Interface(), binding.Form); err != nil {
						return
					}
				}
			case binding.MIMEMultipartPOSTForm:
				if isCreateMethod {
					if err := c.MustBindWith(ri.Interface(), binding.FormMultipart); err != nil {
						return
					}
				}
			default:
				if !isGetMethod {
					if err := c.BindJSON(ri.Interface()); err != nil {
						return
					}
				}
			}
		}
		if rv, ok := ri.Interface().(Validator); ok {
			if err := rv.Validate(); err != nil {
				_ = c.Error(Errorw(http.StatusBadRequest, err)).
					SetType(gin.ErrorTypeBind)
				return
			}
		}
		if rv, ok := ri.Interface().(ValidatorWithInput); ok {
			hv, ok := hr.Interface().(ValidatingInput)
			if ok {
				if err := rv.ValidateWith(c, hv.Validating()); err != nil {
					_ = c.Error(Errorw(http.StatusBadRequest, err)).
						SetType(gin.ErrorTypeBind)
					return
				}
			}
		}
		if it.Kind() != reflect.Pointer {
			ri = ri.Elem()
		}

		// process
		var inputs = make([]reflect.Value, 0, 2)
		inputs = append(inputs, reflect.ValueOf(c))
		inputs = append(inputs, ri)
		var outputs = mr.Call(inputs)

		// render response
		if c.Request.Context().Err() != nil ||
			c.Writer.Size() >= 0 ||
			len(c.Errors) != 0 {
			// already render inside the above processing
			return
		}
		var errInterface = outputs[len(outputs)-1].Interface()
		if errInterface != nil {
			var err = errInterface.(error)
			var ge = c.Error(err)
			if !isGinError(err) {
				_ = ge.SetType(gin.ErrorTypePublic)
			}
			return
		}
		var code = http.StatusOK
		switch len(outputs) {
		case 1:
			if isCreateMethod {
				code = http.StatusNoContent
			}
			c.Writer.WriteHeader(code)
		case 2:
			if outputs[0].IsZero() {
				if isCreateMethod {
					code = http.StatusNoContent
				}
				c.Writer.WriteHeader(code)
				return
			}
			if isCreateMethod {
				code = http.StatusCreated
			}
			var obj = outputs[0].Interface()
			switch ot := obj.(type) {
			case rendCloser:
				if ot == nil {
					return
				}
				defer func() { _ = ot.Close() }()
				c.Render(code, ot)
			case render.Render:
				if ot == nil {
					return
				}
				c.Render(code, ot)
			default:
				c.JSON(code, obj) // TODO negotiate
			}
		case 3:
			var obj = GetResponseCollection(c, outputs[0].Interface(), int(outputs[1].Int()))
			c.JSON(code, obj) // TODO negotiate
		}
	}

	// register route
	var rp = path.Join("/", resourcePath)
	var xrp = path.Join(rp, "_", "batch")
	switch {
	case isCreateMethod:
		r.POST(xrp, vh)
	case isDeleteMethod:
		r.DELETE(rp, vh)
	case isUpdateMethod:
		r.PUT(rp, vh)
	case isGetMethod:
		r.GET(rp, vh)
	}

	// register schema
	var handle = hr.Type().String() + "." + mn
	var ots []reflect.Type
	switch mrt.NumOut() {
	case 2:
		ots = []reflect.Type{mrt.Out(0)}
	case 3:
		ots = []reflect.Type{mrt.Out(0), mrt.Out(1)}
	}
	var op = GetOutputProfile(ots...)
	switch {
	case isCreateMethod:
		schemeRoute(resource, handle, http.MethodPost, xrp, ip, op)
	case isDeleteMethod:
		schemeRoute(resource, handle, http.MethodDelete, rp, ip, op)
	case isUpdateMethod:
		schemeRoute(resource, handle, http.MethodPut, rp, ip, op)
	case isGetMethod:
		schemeRoute(resource, handle, http.MethodGet, rp, ip, op)
	}
	return nil
}

func registerExtensionalHandler(r gin.IRoutes, hr reflect.Value, mn string, k string) error {
	var isPluralMethod = strings.HasPrefix(mn, "Collection")
	var isCreateMethod = strings.HasPrefix(mn, "Create") || strings.HasPrefix(mn, "CollectionCreate")
	var isDeleteMethod = strings.HasPrefix(mn, "Delete") || strings.HasPrefix(mn, "CollectionDelete")
	var isUpdateMethod = strings.HasPrefix(mn, "Update") || strings.HasPrefix(mn, "CollectionUpdate")
	var isGetMethod = strings.HasPrefix(mn, "Get") || strings.HasPrefix(mn, "CollectionGet")
	var isRouteMethod = strings.HasPrefix(mn, "Route") || strings.HasPrefix(mn, "CollectionRoute")
	if !(isCreateMethod || isDeleteMethod || isUpdateMethod || isGetMethod || isRouteMethod) {
		return nil
	}
	switch mn {
	case "CollectionCreateBatch":
		// avoid conflicting with batch apis.
		return fmt.Errorf("invalid name of '%s': use '%s' for batch operation", mn, mn[:len(mn)-5])
	default:
	}

	// validate
	var mr = hr.MethodByName(mn)
	var mrt = mr.Type()
	switch mrt.NumIn() {
	default:
		return fmt.Errorf("invalid input number of '%s': got %d but expected 2", mn, mrt.NumIn())
	case 2:
		if !isTypeOfGinContext(mrt.In(0)) {
			return fmt.Errorf("illegal first input type of '%s': expected *gin.Context", mn)
		}
		switch mrt.In(1).Kind() {
		default:
			return fmt.Errorf("illegal last input type of '%s': expected struct, pointer or slice type", mn)
		case reflect.Struct, reflect.Pointer, reflect.Slice:
		}
	}
	switch mrt.NumOut() {
	default:
		return fmt.Errorf("invalid output number of '%s': got %d but expected not more than 3", mn, mrt.NumOut())
	case 3:
		if !isGetMethod {
			return fmt.Errorf("invalid output number of '%s': got %d but expected not more than 2", mn, mrt.NumOut())
		}
		if !isTypeOfInt(mrt.Out(1)) {
			return fmt.Errorf("illegal second output type of '%s': expected int", mn)
		}
		if !isImplementationOfError(mrt.Out(2)) {
			return fmt.Errorf("illegal last output type of '%s': expected error", mn)
		}
	case 2:
		if !isImplementationOfError(mrt.Out(1)) {
			return fmt.Errorf("illegal last output type of '%s': expected error", mn)
		}
	case 1:
		if isGetMethod {
			return fmt.Errorf("invalid output number of '%s': got %d but expected at least 2", mn, mrt.NumOut())
		}
		if !isImplementationOfError(mrt.Out(0)) {
			return fmt.Errorf("illegal last output type of '%s': expected error", mn)
		}
	}

	// construct virtual handler
	var resource, resourcePath = getResourceAndResourcePath(k)
	var it = mrt.In(1)
	var ip = GetInputProfile(it)
	var ipState = ip.State()
	if ip.Router != nil {
		isCreateMethod = ip.Router.Method == http.MethodPost
		isDeleteMethod = ip.Router.Method == http.MethodDelete
		isUpdateMethod = ip.Router.Method == http.MethodPut
		isGetMethod = ip.Router.Method == http.MethodGet
	}
	var vh = func(c *gin.Context) {
		var s = session.LoadSubject(c)
		if !s.Enforce(c, resource) {
			if s.IsAnonymous() {
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
			return
		}
		session.StoreSubjectCurrentOperation(c, s.Give(resource).If(c.Request.Method))

		// validate
		if !isPluralMethod {
			var id = c.Param("id")
			if id == "" || id == "_" {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
		}
		// bind input
		var ri = reflect.New(it)
		if it.Kind() == reflect.Pointer {
			ri = reflect.New(it.Elem())
		} else {
			ri = reflect.New(it)
		}
		if ipState[ProfileCategoryHeader] {
			if err := c.BindHeader(ri.Interface()); err != nil {
				return
			}
		}
		if ipState[ProfileCategoryUri] {
			if err := c.BindUri(ri.Interface()); err != nil {
				return
			}
		}
		if ipState[ProfileCategoryQuery] {
			if err := binding.MapFormWithTag(ri.Interface(), c.Request.URL.Query(), "query"); err != nil {
				return
			}
		}
		if c.Request.ContentLength != 0 {
			switch c.ContentType() {
			case binding.MIMEPOSTForm:
				if isCreateMethod {
					if err := c.MustBindWith(ri.Interface(), binding.Form); err != nil {
						return
					}
				}
			case binding.MIMEMultipartPOSTForm:
				if isCreateMethod {
					if err := c.MustBindWith(ri.Interface(), binding.FormMultipart); err != nil {
						return
					}
				}
			default:
				if !isGetMethod {
					if err := c.BindJSON(ri.Interface()); err != nil {
						return
					}
				}
			}
		}
		if rv, ok := ri.Interface().(Validator); ok {
			if err := rv.Validate(); err != nil {
				_ = c.Error(Errorw(http.StatusBadRequest, err)).
					SetType(gin.ErrorTypeBind)
				return
			}
		}
		if rv, ok := ri.Interface().(ValidatorWithInput); ok {
			hv, ok := hr.Interface().(ValidatingInput)
			if ok {
				if err := rv.ValidateWith(c, hv.Validating()); err != nil {
					_ = c.Error(Errorw(http.StatusBadRequest, err)).
						SetType(gin.ErrorTypeBind)
					return
				}
			}
		}
		if it.Kind() != reflect.Pointer {
			ri = ri.Elem()
		}

		// process
		var inputs = make([]reflect.Value, 0, 2)
		inputs = append(inputs, reflect.ValueOf(c))
		inputs = append(inputs, ri)
		var outputs = mr.Call(inputs)

		// render response
		if c.Request.Context().Err() != nil ||
			c.Writer.Size() >= 0 ||
			len(c.Errors) != 0 {
			// already render inside the above processing
			return
		}
		var errInterface = outputs[len(outputs)-1].Interface()
		if errInterface != nil {
			var err = errInterface.(error)
			var ge = c.Error(err)
			if !isGinError(err) {
				_ = ge.SetType(gin.ErrorTypePublic)
			}
			return
		}
		var code = http.StatusOK
		switch len(outputs) {
		case 1:
			if isCreateMethod {
				code = http.StatusNoContent
			}
			c.Writer.WriteHeader(code)
		case 2:
			if outputs[0].IsZero() {
				if isCreateMethod {
					code = http.StatusNoContent
				}
				c.Writer.WriteHeader(code)
				return
			}
			if isCreateMethod {
				code = http.StatusCreated
			}
			var obj = outputs[0].Interface()
			switch ot := obj.(type) {
			case rendCloser:
				if ot == nil {
					return
				}
				defer func() { _ = ot.Close() }()
				c.Render(code, ot)
			case render.Render:
				if ot == nil {
					return
				}
				c.Render(code, ot)
			default:
				c.JSON(code, obj) // TODO negotiate
			}
		case 3:
			var obj = GetResponseCollection(c, outputs[0].Interface(), int(outputs[1].Int()))
			c.JSON(code, obj) // TODO negotiate
		}
	}

	// register route
	var rps = []string{"/", resourcePath}
	var mnt = strings.TrimPrefix(mn, "Collection")
	if mnt != mn { // general items
		if ip.Router != nil && isCreateMethod && ip.Router.SubPath == "/batch" {
			// avoid conflicting with batch apis.
			return fmt.Errorf("invalid subpath definition of '%s': confilict with batch operation", mn)
		}
		rps = append(rps, "_")
		mn = mnt
	} else { // specific item
		rps = append(rps, ":id")
	}
	if ip.Router != nil {
		rps = append(rps, ip.Router.SubPath)
	} else {
		if !isGetMethod {
			rps = append(rps, strs.Dasherize(mn[6:])) // Create/Update/Delete
		} else {
			rps = append(rps, strs.Dasherize(mn[3:])) // Get
		}
	}
	var rp = path.Join(rps...)
	switch {
	case isCreateMethod:
		r.POST(rp, vh)
	case isDeleteMethod:
		r.DELETE(rp, vh)
	case isUpdateMethod:
		r.PUT(rp, vh)
	case isGetMethod:
		r.GET(rp, vh)
	}

	// register schema
	var handle = hr.Type().String() + "." + mn
	var ots []reflect.Type
	switch mrt.NumOut() {
	case 2:
		ots = []reflect.Type{mrt.Out(0)}
	case 3:
		ots = []reflect.Type{mrt.Out(0), mrt.Out(1)}
	}
	var op = GetOutputProfile(ots...)
	switch {
	case isCreateMethod:
		schemeRoute(resource, handle, http.MethodPost, rp, ip, op)
	case isDeleteMethod:
		schemeRoute(resource, handle, http.MethodDelete, rp, ip, op)
	case isUpdateMethod:
		schemeRoute(resource, handle, http.MethodPut, rp, ip, op)
	case isGetMethod:
		schemeRoute(resource, handle, http.MethodGet, rp, ip, op)
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
