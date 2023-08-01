package bind

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const _32MiB = 32 << 20

// WithForm binds the passed struct pointer using the request form params,
// fails with false returning and aborts the request with 400 status code.
func WithForm(c *gin.Context, r any) bool {
	if c.Request == nil {
		return abortWithError(c, errors.New("invalid request"))
	}

	if err := parseRequestForm(c.Request); err != nil {
		return abortWithError(c, err)
	}

	err := MapFormWithTag(r, c.Request.Form, "form")

	return abortWithError(c, err)
}

func parseRequestForm(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	if err := req.ParseMultipartForm(_32MiB); err != nil && !errors.Is(err, http.ErrNotMultipart) {
		return err
	}

	return nil
}

// WithJSON binds the passed struct pointer using the request json params,
// fails with false returning and aborts the request with 400 status code.
func WithJSON(c *gin.Context, r any) bool {
	if c.Request == nil || c.Request.Body == nil {
		return abortWithError(c, errors.New("invalid request"))
	}

	err := json.NewDecoder(c.Request.Body).Decode(r)

	return abortWithError(c, err)
}

// WithHeader binds the passed struct pointer using the request header params,
// fails with false returning and aborts the request with 400 status code.
func WithHeader(c *gin.Context, r any) bool {
	m := c.Request.Header

	err := MapFormWithTag(r, m, "header")

	return abortWithError(c, err)
}

// WithQuery binds the passed struct pointer using the request query params,
// fails with false returning and aborts the request with 400 status code.
func WithQuery(c *gin.Context, r any) bool {
	m := c.Request.URL.Query()

	err := MapFormWithTag(r, m, "query")

	return abortWithError(c, err)
}

// WithPath binds the passed struct pointer using the request path params,
// fails with false returning and aborts the request with 400 status code.
func WithPath(c *gin.Context, r any) bool {
	m := make(map[string][]string)
	for _, v := range c.Params {
		m[v.Key] = []string{v.Value}
	}

	err := MapFormWithTag(r, m, "path")

	return abortWithError(c, err)
}

// abortWithError breaks the call chain if found error,
// and returns false.
func abortWithError(c *gin.Context, err error) bool {
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err).
			SetType(gin.ErrorTypeBind)
		return false
	}

	return true
}
