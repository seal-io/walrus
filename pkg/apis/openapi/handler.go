package openapi

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/ogen-go/ogen"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/casdoor"
	"github.com/seal-io/seal/utils/bytespool"
	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/version"
)

func Index(authing bool, apiPrefix string) runtime.Handle {
	var spec = runtime.OpenAPI(ogen.NewInfo().
		SetTitle("Seal APIs").
		SetDescription("API to manage resources of Seal").
		SetVersion(version.Version))

	if apiPrefix != "" {
		spec.AddServers(ogen.NewServer().SetURL(apiPrefix))
	}

	if authing {
		var securities = ogen.SecurityRequirements{
			{"bearerAuth": {}},
			{"cookieAuth": {}},
		}
		spec.Security = securities
		var securitySchemes = make(map[string]*ogen.SecurityScheme, 2)
		securitySchemes["bearerAuth"] = &ogen.SecurityScheme{
			Type:         "http",
			In:           "header",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		}
		securitySchemes["cookieAuth"] = &ogen.SecurityScheme{
			Type: "apiKey",
			In:   "cookie",
			Name: casdoor.ExternalSessionCookieKey,
		}
		spec.Components.SecuritySchemes = securitySchemes
		for i := range spec.Paths {
			var ops = []*ogen.Operation{
				spec.Paths[i].Post,
				spec.Paths[i].Delete,
				spec.Paths[i].Put,
				spec.Paths[i].Get,
			}
			for j := range ops {
				if ops[j] == nil {
					continue
				}
				ops[j].Security = securities
			}
		}
	}

	// add named schemas.
	spec.AddNamedSchemas()

	var specBytes, err = json.Marshal(spec)
	if err != nil {
		panic(fmt.Errorf("error marshalling openapi spec: %w", err))
	}

	return func(c *gin.Context) {
		var buff = bytespool.GetBytes(0)
		defer func() { bytespool.Put(buff) }()
		_, _ = io.CopyBuffer(c.Writer, bytes.NewBuffer(specBytes), buff)
	}
}
