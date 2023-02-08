package ui

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/settings"
	"github.com/seal-io/seal/utils/req"
)

func Index(ctx context.Context, modelClient model.ClientSet) runtime.Handle {
	var defaultUiIndex, _ = url.Parse("file:///var/lib/seal/ui")
	return func(c *gin.Context) {
		var path = c.Request.URL.Path

		// redirect.
		switch path {
		case "":
			http.Redirect(c.Writer, c.Request, "/", http.StatusMovedPermanently)
			c.Abort()
			return
		case "/verify-auth":
			http.Redirect(c.Writer, c.Request, "/#/integration/oauth?"+c.Request.URL.RawQuery, http.StatusFound)
			c.Abort()
			return
		}

		// ui handle.
		var uiIndex = settings.ServeUiIndex.ShouldValueURL(ctx, modelClient)
		if uiIndex == nil {
			uiIndex = defaultUiIndex
		}
		switch path {
		case "/":
			uiSrv(uiIndex).ServeHTTP(c.Writer, c.Request)
			c.Abort()
		default:
			if strings.HasPrefix(path, "/assets/") {
				// assets
				uiSrv(uiIndex).ServeHTTP(c.Writer, c.Request)
				c.Abort()
				return
			}
			c.Next()
		}
	}
}

func uiSrv(uiIndex *url.URL) http.Handler {
	switch uiIndex.Scheme {
	case "file":
		return local(uiIndex.Path)
	default:
		return remote(uiIndex.String())
	}
}

func local(dir string) http.Handler {
	var fs = runtime.StaticHttpFileSystem{
		FileSystem: http.FS(os.DirFS(dir)),
	}
	return http.FileServer(fs)
}

func remote(uri string) http.HandlerFunc {
	var httpClient = req.HTTP().
		WithMaxConnDuration(0).
		Request()
	return func(w http.ResponseWriter, r *http.Request) {
		var body, err = httpClient.GetWithContext(r.Context(), uri).
			Body()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusGone), http.StatusGone)
			return
		}
		_, _ = io.Copy(w, body)
	}
}
