package runtime

import (
	"database/sql"
	"errors"
	"net/http"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/log"
)

// Erroring is a gin middleware,
// which converts the chain calling error into response.
func Erroring() Handle {
	var logger = log.WithName("apis")
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
			var cause = errors.Unwrap(ge.Err)
			if cause == nil {
				cause = ge.Err
			}
			var code = http.StatusInternalServerError
			if ge.Type == gin.ErrorTypeBind {
				code = http.StatusBadRequest
			}
			switch {
			case isBadRequestError(cause):
				code = http.StatusBadRequest
			case isNotFoundError(cause):
				code = http.StatusNotFound
			case isConflictError(cause):
				code = http.StatusConflict
			case isUnprocessableEntityError(cause):
				code = http.StatusUnprocessableEntity
			case isInternalServerError(cause):
				code = http.StatusInternalServerError
			}
			he.code = code
			he.cause = ge.Err
		}
		if ge.Type == gin.ErrorTypePrivate {
			he.cause = nil // mute
		}
	}

	// correct the code if already write within context.
	if c.Writer.Written() {
		he.code = c.Writer.Status()
	}

	return
}

func isBadRequestError(err error) bool {
	return model.IsValidationError(err)
}

func isNotFoundError(err error) bool {
	return model.IsNotFound(err) || model.IsNotSingular(err)
}

func isConflictError(err error) bool {
	return model.IsConstraintError(err) || sqlgraph.IsConstraintError(err)
}

func isUnprocessableEntityError(err error) bool {
	return model.IsNotLoaded(err) || errors.Is(err, sql.ErrNoRows)
}

func isInternalServerError(err error) bool {
	return errors.Is(err, sql.ErrConnDone) || errors.Is(err, sql.ErrTxDone)
}
