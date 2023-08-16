package runtime

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/walrus/utils/log"
)

// recovering is a gin middleware,
// which is the same as gin.Recovery,
// but friendly message information can be provided according to the request header.
func recovering(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("panic observing: %v", r)
			}

			log.WithName("api").
				Errorf("panic observing: %v, callstack: \n%s",
					err, getPanicCallstack(3))

			if isStreamRequest(c) {
				// Stream request always send header at first,
				// so we don't need to rewrite.
				return
			}

			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()

	c.Next()
}

func getPanicCallstack(skip int) []byte {
	var buf bytes.Buffer

	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := "???"
		f := runtime.FuncForPC(pc)

		if f != nil {
			fn = f.Name()
		}
		_, _ = fmt.Fprintf(&buf, "%s\n\t%s:%d (0x%x)\n", fn, file, line, pc)
	}

	return buf.Bytes()
}
