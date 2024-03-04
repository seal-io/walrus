package ui

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/seal-io/utils/funcx"
	"github.com/seal-io/utils/httpx"

	"github.com/seal-io/walrus/pkg/system"
)

var Dir = system.SubDataDir("ui")

func Index(ctx context.Context) http.Handler {
	// TODO: support settings.
	uri := funcx.MustNoError(url.Parse(
		"https://walrus-ui-1303613262.cos.ap-guangzhou.myqcloud.com/latest/index.html"))
	return serve(uri)
}

func serve(uri *url.URL) http.Handler {
	if uri.Scheme == "file" {
		return local(uri.Path)
	}
	return remote(uri.String())
}

func local(dir string) http.Handler {
	return http.FileServer(httpx.FS(os.DirFS(dir)))
}

func remote(uri string) http.HandlerFunc {
	cli := httpx.DefaultClient

	return func(w http.ResponseWriter, r *http.Request) {
		req, err := httpx.NewGetRequestWithContext(r.Context(), uri)
		if err != nil {
			httpx.Error(w, http.StatusNotFound)
			return
		}

		resp, err := cli.Do(req)
		if err != nil {
			httpx.Error(w, http.StatusGone)
			return
		}
		defer httpx.Close(resp)

		_, _ = io.Copy(w, resp.Body)
	}
}
