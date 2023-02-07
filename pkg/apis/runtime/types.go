package runtime

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Resource interface {
	Kind() string
}

type (
	AdviceResource interface {
		Resource

		ResourceAndResourcePath() (resource, resourcePath string)
		Unwrap() Resource
	}

	AdviceBeforeResourceRegistering interface {
		BeforeAdvice(AdviceResource) error
	}

	AdviceAfterResourceRegistering interface {
		AfterAdvice(AdviceResource) error
	}
)

type (
	Validator interface {
		Validate() error
	}

	ValidatorWithInput interface {
		ValidateWith(ctx context.Context, input any) error
	}

	ValidatingInput interface {
		Validating() any
	}
)

type ErrorHandler interface {
	Handle(c *gin.Context) error
}

func WrapErrorHandler(h ErrorHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if h == nil {
			c.Next()
			return
		}
		var err = h.Handle(c)
		if err != nil {
			_ = c.Error(err).
				SetType(gin.ErrorTypePublic)
			c.Abort()
		}
	}
}

type ErrorHandle func(c *gin.Context) error

func WrapErrorHandle(f ErrorHandle) gin.HandlerFunc {
	return func(c *gin.Context) {
		if f == nil {
			c.Next()
			return
		}
		var err = f(c)
		if err != nil {
			_ = c.Error(err).
				SetType(gin.ErrorTypePublic)
			c.Abort()
		}
	}
}

type Handler interface {
	Handle(c *gin.Context)
}

func WrapHandler(h Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if h == nil {
			c.Next()
			return
		}
		h.Handle(c)
	}
}

type Handle = gin.HandlerFunc

type HTTPHandler = http.Handler

func WrapHTTPHandler(h http.Handler) gin.HandlerFunc {
	return gin.WrapH(h)
}

type HTTPHandle = http.HandlerFunc

func WrapHTTPHandle(h http.HandlerFunc) gin.HandlerFunc {
	return gin.WrapF(h)
}
