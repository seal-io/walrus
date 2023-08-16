package cli

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/walrus/pkg/apis/runtime"
)

const (
	cliPath = "/var/lib/seal/cli"
	cli     = "seal"
)

func Index() runtime.Handle {
	return func(c *gin.Context) {
		osArch := fmt.Sprintf("%s-%s", c.Query("os"), c.Query("arch"))
		filePath := path.Join(cliPath, fmt.Sprintf("cli-%s", osArch))

		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, cli))
		c.Header("Content-Type", "application/octet-stream")

		http.ServeFile(c.Writer, c.Request, filePath)
	}
}
