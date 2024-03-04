package clis

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"path/filepath"

	"github.com/seal-io/utils/httpx"
	"github.com/seal-io/utils/osx"

	"github.com/seal-io/walrus/pkg/system"
)

var Dir = system.SubDataDir("clis")

func Index(_ context.Context) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		reqFile := path.Base(r.URL.Path)
		locFile := filepath.Join(Dir, reqFile)

		if !osx.ExistsFile(locFile) {
			httpx.Error(rw, http.StatusNotFound)
			return
		}

		rwh := rw.Header()
		rwh.Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, reqFile))
		rwh.Add("Content-Type", "application/octet-stream")

		http.ServeFile(rw, r, locFile)
	})
}
