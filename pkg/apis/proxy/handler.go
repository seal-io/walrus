package proxy

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

type proxy interface {
	handle(path, rawQuery string)
}

func Proxy(mc *model.Client) runtime.Handle {
	return func(context *gin.Context) {
		const key = "connectorID"

		var conn *model.Connector

		var err error

		query := context.Request.URL.Query()

		connectorID := query.Get(key)
		if connectorID != "" {
			conn, err = mc.Connectors().Get(context, object.ID(connectorID))
			if err != nil {
				_ = context.AbortWithError(http.StatusBadRequest, err)
				return
			}
			// Remove connectorID from query.
			query.Del(key)
		}

		if !authorize(context, conn) {
			return
		}

		// Get proxy request path.
		p := strings.TrimPrefix(context.Param("path"), "/")

		proxy := newProxy(context, conn)
		proxy.handle(p, query.Encode())
	}
}

// authorize checks if the subject is authorized to access the connector.
// If the connector is nil, it will check if the subject is authenticated.
func authorize(context *gin.Context, conn *model.Connector) bool {
	sj, err := session.GetSubject(context)
	if err != nil {
		_ = context.AbortWithError(http.StatusBadRequest, err)
		return false
	}

	if sj.IsAnonymous() {
		context.AbortWithStatus(http.StatusUnauthorized)
		return false
	}

	if sj.IsAdmin() || conn == nil {
		return true
	}

	// Check connector permission.
	resources := []session.ActionResource{
		{Name: "projects", Refer: conn.ProjectID.String()},
		{Name: "connectors", Refer: conn.ID.String()},
	}
	if !sj.Enforce(http.MethodGet, resources, "") {
		context.AbortWithStatus(http.StatusForbidden)
		return false
	}

	return true
}

// newProxy creates a proxy instance based on the connector type.
func newProxy(context *gin.Context, conn *model.Connector) proxy {
	switch conn.Type {
	case types.ConnectorTypeKubernetes:
		return &k8sProxy{
			conn:    conn,
			context: context,
		}

	case types.ConnectorTypeAWS:
		return &awsProxy{
			conn:    conn,
			context: context,
			whiteListDomains: []string{
				"*.amazonaws.com",
				"*.amazonaws.com.cn",
			},
		}
	}

	return nil
}
