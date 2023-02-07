package runtime

import (
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
			logger.Error("runtime errors: ", me.Errors())
		}

		// get last error from chain and parse into response
		var he = getHttpError(c)
		c.AbortWithStatusJSON(he.code, he) // TODO negotiate
	}
}

func getHttpError(c *gin.Context) (he httpError) {
	var ge = c.Errors.Last()

	if ge == nil || ge.Err == nil {
		he.code = http.StatusServiceUnavailable
	} else {
		if !errors.As(ge.Err, &he) {
			var cause = errors.Unwrap(ge.Err)
			if cause == nil {
				cause = ge.Err
			}
			var code = http.StatusServiceUnavailable
			switch {
			case model.IsNotFound(cause) || model.IsNotSingular(cause):
				code = http.StatusNotFound
			case model.IsConstraintError(cause) || sqlgraph.IsConstraintError(cause):
				code = http.StatusConflict
			case model.IsNotLoaded(cause):
				code = http.StatusUnprocessableEntity
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
