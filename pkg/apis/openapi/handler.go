package openapi

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/ogen-go/ogen"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/auths"
	"github.com/seal-io/seal/utils/bytespool"
	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/version"
)

func Index(authing bool, apiPrefix string) runtime.Handle {
	spec := runtime.OpenAPI(ogen.NewInfo().
		SetTitle("Seal APIs").
		SetDescription("API to manage resources of Seal").
		SetVersion(version.Version))

	if apiPrefix != "" {
		spec.AddServers(ogen.NewServer().SetURL(apiPrefix))
	}

	if authing {
		securities := ogen.SecurityRequirements{
			{"bearerAuth": {}},
			{"cookieAuth": {}},
		}
		spec.Security = securities
		securitySchemes := make(map[string]*ogen.SecurityScheme, 2)
		securitySchemes["bearerAuth"] = &ogen.SecurityScheme{
			Type:         "http",
			In:           "header",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		}
		securitySchemes["cookieAuth"] = &ogen.SecurityScheme{
			Type: "apiKey",
			In:   "cookie",
			Name: auths.SessionCookieName,
		}

		spec.Components.SecuritySchemes = securitySchemes
		for i := range spec.Paths {
			ops := []*ogen.Operation{
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

	// Add named schemas.
	spec.AddNamedSchemas()

	specBytes, err := json.Marshal(spec)
	if err != nil {
		panic(fmt.Errorf("error marshaling openapi spec: %w", err))
	}

	return func(c *gin.Context) {
		buff := bytespool.GetBytes(0)
		defer func() { bytespool.Put(buff) }()
		_, _ = io.CopyBuffer(c.Writer, bytes.NewBuffer(specBytes), buff)
	}
}
