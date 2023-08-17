package runtime

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/log"
)

type Handle = func(*gin.Context)

type ErrorHandle = func(c *gin.Context) error

func wrapErrorHandle(f ErrorHandle) Handle {
	return func(c *gin.Context) {
		if f == nil {
			c.Next()
			return
		}

		err := f(c)
		if err != nil {
			_ = c.Error(errorx.Wrap(err, ""))

			c.Abort()
		}
	}
}

type Handler interface {
	Handle(c *gin.Context)
}

func wrapHandler(h Handler) Handle {
	return func(c *gin.Context) {
		if h == nil {
			c.Next()
			return
		}

		h.Handle(c)
	}
}

type ErrorHandler interface {
	Handle(c *gin.Context) error
}

func wrapErrorHandler(h ErrorHandler) Handle {
	return func(c *gin.Context) {
		if h == nil {
			c.Next()
			return
		}

		err := h.Handle(c)
		if err != nil {
			_ = c.Error(errorx.Wrap(err, ""))

			c.Abort()
		}
	}
}

type HTTPHandler = http.Handler

func wrapHTTPHandler(h http.Handler) Handle {
	return gin.WrapH(h)
}

type HTTPHandle = http.HandlerFunc

func wrapHTTPHandle(h http.HandlerFunc) Handle {
	return gin.WrapF(h)
}

func asHandle(h IHandler) Handle {
	if h != nil {
		switch t := h.(type) {
		case Handle:
			return t
		case ErrorHandle:
			return wrapErrorHandle(t)
		case Handler:
			return wrapHandler(t)
		case ErrorHandler:
			return wrapErrorHandler(t)
		case HTTPHandle:
			return wrapHTTPHandle(t)
		case HTTPHandler:
			return wrapHTTPHandler(t)
		}
	}

	log.WithName("api").
		Errorf("unknown handle type: %T", h)

	return func(c *gin.Context) {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
