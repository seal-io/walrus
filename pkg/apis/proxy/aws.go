package proxy

import (
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/gin-gonic/gin"

	"github.com/seal-io/walrus/pkg/dao/model"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
)

type awsProxy struct {
	conn             *model.Connector
	context          *gin.Context
	whiteListDomains []string
}

func (p *awsProxy) handle(requestPath, rawQuery string) {
	u, err := getURL(requestPath)
	if err != nil {
		_ = p.context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if !isWhiteListDomain(u.Hostname(), p.whiteListDomains...) {
		p.context.AbortWithStatus(http.StatusForbidden)
		return
	}

	cred, err := optypes.GetCredential(p.conn.ConfigData)
	if err != nil {
		_ = p.context.AbortWithError(http.StatusBadRequest, err)
		return
	}

	rp := &httputil.ReverseProxy{
		Director: func(request *http.Request) {
			request.URL = u
			request.Host = u.Host
			request.URL.RawQuery = rawQuery

			removeHeaders(request)

			credential := credentials.NewStaticCredentials(cred.AccessKey, cred.AccessSecret, "")
			awsSigner := v4.NewSigner(credential)

			service, region := getServiceAndRegion(u.Host)

			_, err = awsSigner.Sign(request, nil, service, region, time.Now())
			if err != nil {
				_ = p.context.AbortWithError(http.StatusBadRequest, err)
				return
			}
		},
	}

	rp.ServeHTTP(p.context.Writer, p.context.Request)
}

func getServiceAndRegion(host string) (string, string) {
	// Format : service.region.*.
	parts := strings.Split(host, ".")
	return parts[0], parts[1]
}
