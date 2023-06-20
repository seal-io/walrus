package cli

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/runtime"
)

var (
	cliPath = "/var/lib/seal/cli/seal-cli"
	cli     = "seal-cli"
)

func Index() runtime.Handle {
	return func(c *gin.Context) {
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, cli))
		c.Header("Content-Type", "application/octet-stream")
		http.ServeFile(c.Writer, c.Request, cliPath)
	}
}
