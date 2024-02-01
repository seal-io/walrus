package proxy

import (
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/dao/model"
	opk8s "github.com/seal-io/walrus/pkg/operator/k8s"
)

type k8sProxy struct {
	conn    *model.Connector
	context *gin.Context
}

func (p *k8sProxy) handle(requestPath, rawQuery string) {
	restCfg, err := opk8s.GetConfig(*p.conn)
	if err != nil {
		_ = p.context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	clusterClient, err := rest.HTTPClientFor(restCfg)
	if err != nil {
		_ = p.context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	u, err := getProxyURL(restCfg.Host, requestPath, rawQuery)
	if err != nil {
		_ = p.context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	response, err := clusterClient.Get(u)
	if err != nil {
		_ = p.context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, err = io.Copy(p.context.Writer, response.Body)
	if err != nil {
		_ = p.context.AbortWithError(http.StatusBadRequest, err)
		return
	}
}

func getProxyURL(host, requestPath, query string) (string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, requestPath)
	u.RawQuery = query

	return u.String(), nil
}
