package cli

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/utils/files"
)

const (
	cliPath = "/var/lib/walrus/cli"
)

func Index() runtime.Handle {
	return func(c *gin.Context) {
		var (
			os   = c.Query("os")
			arch = c.Query("arch")
		)

		if os == "" || arch == "" {
			http.NotFound(c.Writer, c.Request)
			return
		}

		var suffix string

		if os == "windows" {
			suffix = ".exe"
		}

		var (
			fileName = fmt.Sprintf("cli-%s-%s%s", os, arch, suffix)
			filePath = filepath.Join(cliPath, fileName)
		)

		if !files.ExistsFile(filePath) {
			http.NotFound(c.Writer, c.Request)
			return
		}

		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="walrus%s"`, suffix))
		c.Header("Content-Type", "application/octet-stream")

		http.ServeFile(c.Writer, c.Request, filePath)
	}
}
