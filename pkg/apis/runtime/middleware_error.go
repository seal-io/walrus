package runtime

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/gin-gonic/gin"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/log"
)

// erroring is a gin middleware,
// which converts the chain calling error into response.
func erroring(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		if c.Writer.Status() >= http.StatusBadRequest && c.Writer.Size() == 0 {
			// Detail the error status message.
			_ = c.Error(errorx.NewHttpError(c.Writer.Status(), ""))
		} else {
			// No errors.
			return
		}
	}

	// Get errors from chain and parse into response.
	he := getHttpError(c)

	// Log errors.
	if len(he.errs) != 0 && withinStacktraceStatus(he.Status) {
		reqMethod := c.Request.Method

		reqPath := c.Request.URL.Path
		if raw := c.Request.URL.RawQuery; raw != "" {
			reqPath = reqPath + "?" + raw
		}

		log.WithName("api").
			Errorf("error requesting %s %s: %v", reqMethod, reqPath, errorx.Format(he.errs))
	}

	c.AbortWithStatusJSON(he.Status, he)
}

func getHttpError(c *gin.Context) (he ErrorResponse) {
	var errs []error

	for i := range c.Errors {
		if c.Errors[i].Err != nil {
			errs = append(errs, c.Errors[i].Err)
		}
	}
	he.errs = errs

	if len(errs) == 0 {
		he.Status = http.StatusInternalServerError
	} else {
		// Get the public error.
		he.Status, he.Message = errorx.Public(errs)

		// Get the last error.
		if he.Status == 0 {
			st, msg := diagnoseError(c.Errors.Last())
			he.Status = st

			if he.Message == "" {
				he.Message = msg
			}
		}
	}

	// Correct the code if already write within context.
	if c.Writer.Written() {
		he.Status = c.Writer.Status()
	}

	he.StatusText = http.StatusText(he.Status)

	return
}

type ErrorResponse struct {
	Message    string `json:"message"`
	Status     int    `json:"status"`
	StatusText string `json:"statusText"`

	// Errs is the all errors from gin context errors.
	errs []error
}

func diagnoseError(ge *gin.Error) (int, string) {
	c := http.StatusInternalServerError
	if ge.Type == gin.ErrorTypeBind {
		c = http.StatusBadRequest
	}

	var b strings.Builder

	if ge.Meta != nil {
		m, ok := ge.Meta.(string)
		if ok {
			b.WriteString("failed to ")
			b.WriteString(m)
		}
	}

	err := ge.Err
	if ue := errors.Unwrap(err); ue != nil {
		err = ue
	}

	for i := range diagnosis {
		s := diagnosis[i].probe(err)
		if s == "" {
			continue
		}
		c = diagnosis[i].code

		if b.Len() != 0 {
			b.WriteString(": ")
		}

		b.WriteString(s)

		break
	}

	return c, b.String()
}

var diagnosis = []struct {
	code  int
	probe func(error) string
}{
	{
		code:  http.StatusBadRequest,
		probe: isBadRequestError,
	},
	{
		code:  http.StatusNotFound,
		probe: isNotFoundError,
	},
	{
		code:  http.StatusConflict,
		probe: isConflictError,
	},
	{
		code:  http.StatusUnprocessableEntity,
		probe: isUnprocessableEntityError,
	},
	{
		code:  http.StatusInternalServerError,
		probe: isInternalServerError,
	},
}

func isBadRequestError(err error) string {
	if model.IsValidationError(err) {
		return "datasource: violates validity detecting"
	}

	if strings.Contains(err.Error(), "sql: converting argument") {
		return "datasource: invalid field parsing"
	}

	return ""
}

func isNotFoundError(err error) string {
	switch {
	case model.IsNotFound(err):
		return "datasource: not found"
	case model.IsNotSingular(err):
		return "datasource: found more than one"
	}

	return ""
}

func isConflictError(err error) string {
	switch {
	case sqlgraph.IsUniqueConstraintError(err):
		return "datasource: duplicated"
	case sqlgraph.IsForeignKeyConstraintError(err):
		return "datasource: be depended on other resources"
	}

	return ""
}

func isUnprocessableEntityError(err error) string {
	switch {
	case model.IsNotLoaded(err):
		return "datasource: no dependencies"
	case errors.Is(err, sql.ErrNoRows):
		return "datasource: no changed"
	}

	return ""
}

func isInternalServerError(err error) string {
	switch {
	case errors.Is(err, sql.ErrConnDone):
		return "datasource: closed"
	case errors.Is(err, sql.ErrTxDone):
		return "datasource: transaction closed"
	}

	return ""
}

func withinStacktraceStatus(status int) bool {
	return (status < http.StatusOK || status >= http.StatusInternalServerError) &&
		status != http.StatusSwitchingProtocols
}
