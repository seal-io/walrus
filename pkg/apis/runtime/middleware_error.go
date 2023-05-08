package runtime

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/log"
)

// Erroring is a gin middleware,
// which converts the chain calling error into response.
func Erroring() Handle {
	var logger = log.WithName("api")
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			if c.Writer.Status() >= 400 && c.Writer.Written() && c.Writer.Size() == 0 {
				_ = c.Error(Errorc(c.Writer.Status())).
					SetType(gin.ErrorTypePublic)
			} else {
				return
			}
		}

		// log private errors
		if me := c.Errors.ByType(gin.ErrorTypePrivate); len(me) != 0 {
			var reqPath = c.Request.URL.Path
			if raw := c.Request.URL.RawQuery; raw != "" {
				reqPath = reqPath + "?" + raw
			}
			logger.Errorf("error requesting %s: %v", reqPath, me[len(me)-1])
		}

		// get last error from chain and parse into response
		var he = getHttpError(c)
		c.AbortWithStatusJSON(he.code, he) // TODO negotiate
	}
}

func getHttpError(c *gin.Context) (he httpError) {
	var ge = c.Errors.Last()

	if ge == nil || ge.Err == nil {
		he.code = http.StatusInternalServerError
	} else {
		if !errors.As(ge.Err, &he) {
			he.code, he.brief = diagnoseError(ge)
			he.cause = ge.Err
		}
		if ge.Type == gin.ErrorTypePrivate {
			var we wrapError
			if !errors.As(he.cause, &we) {
				he.cause = nil // mute
			} else {
				he.cause = we.external
			}
		}
	}

	// correct the code if already write within context.
	if c.Writer.Written() {
		he.code = c.Writer.Status()
	}

	return
}

func diagnoseError(ge *gin.Error) (int, string) {
	var c = http.StatusInternalServerError
	if ge.Type == gin.ErrorTypeBind {
		c = http.StatusBadRequest
	}

	var b strings.Builder
	if ge.Meta != nil {
		var m, ok = ge.Meta.(RouteResourceHandleErrorMetadata)
		if ok {
			b.WriteString(m.String())
		}
	}

	var err = ge.Err
	if ue := errors.Unwrap(err); ue != nil {
		err = ue
	}
	for i := range diagnosis {
		var s = diagnosis[i].probe(err)
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
