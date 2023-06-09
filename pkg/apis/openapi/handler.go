package openapi

import (
	"bytes"
	"fmt"
	"io"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/auths"
	"github.com/seal-io/seal/utils/bytespool"
)

func Index(authing bool, apiPrefix string) runtime.Handle {
	spec := runtime.OpenAPI(&openapi3.Info{
		Title:       "Seal APIs",
		Description: "API to manage resources of Seal",
		Version:     "version.Version",
	})

	if apiPrefix != "" {
		spec.AddServer(&openapi3.Server{
			URL: apiPrefix,
		})
	}

	if authing {
		securities := openapi3.SecurityRequirements{
			{"bearerAuth": {}},
			{"cookieAuth": {}},
		}
		spec.Security = securities
		securitySchemes := make(map[string]*openapi3.SecuritySchemeRef, 2)
		securitySchemes["bearerAuth"] = &openapi3.SecuritySchemeRef{
			Value: &openapi3.SecurityScheme{
				Type:         "http",
				In:           "header",
				Scheme:       "bearer",
				BearerFormat: "JWT",
			},
		}
		securitySchemes["cookieAuth"] = &openapi3.SecuritySchemeRef{
			Value: &openapi3.SecurityScheme{
				Type: "apiKey",
				In:   "cookie",
				Name: auths.SessionCookieName,
			},
		}

		spec.Components.SecuritySchemes = securitySchemes
		for i := range spec.Paths {
			ops := []*openapi3.Operation{
				spec.Paths[i].Post,
				spec.Paths[i].Delete,
				spec.Paths[i].Put,
				spec.Paths[i].Get,
			}
			for j := range ops {
				if ops[j] == nil {
					continue
				}
				ops[j].Security = &securities
			}
		}
	}

	specBytes, err := spec.MarshalJSON()
	if err != nil {
		panic(fmt.Errorf("error marshaling openapi spec: %w", err))
	}

	return func(c *gin.Context) {
		buff := bytespool.GetBytes(0)
		defer func() { bytespool.Put(buff) }()
		_, _ = io.CopyBuffer(c.Writer, bytes.NewBuffer(specBytes), buff)
	}
}
