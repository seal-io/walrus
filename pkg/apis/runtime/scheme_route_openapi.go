package runtime

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"net/http"
	"path"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/utils/bytespool"
	"github.com/seal-io/seal/utils/json"
)

// ExposeOpenAPI is a RouterOption to add route to serve the OpenAPI schema spec,
// and provide the SwaggerUI as well.
func ExposeOpenAPI() RouterOption {
	return ginRouteOption(func(r gin.IRouter) {
		const openAPIPath = "/openapi"

		skipLoggingPath(openAPIPath)
		openAPIIndexer := indexOpenAPI()
		r.GET(openAPIPath, openAPIIndexer)

		const swaggerUIPath = "/swagger/*filepath"

		skipLoggingPath(swaggerUIPath)
		swaggerUIIndexer := indexSwaggerUI(openAPIPath)
		r.HEAD(swaggerUIPath, swaggerUIIndexer)
		r.GET(swaggerUIPath, swaggerUIIndexer)
	})
}

func indexOpenAPI() Handle {
	var (
		once        sync.Once
		schemaBytes []byte
	)

	return func(c *gin.Context) {
		once.Do(func() {
			var err error

			schemaBytes, err = json.Marshal(openAPISchemas)
			if err != nil {
				panic(fmt.Errorf("error marshaling openapi schema spec: %w", err))
			}
		})

		buff := bytespool.GetBytes(0)
		defer func() { bytespool.Put(buff) }()
		_, _ = io.CopyBuffer(c.Writer, bytes.NewBuffer(schemaBytes), buff)
	}
}

// downloaded form https://github.com/swagger-api/swagger-ui/releases.
//
//go:embed swagger-ui/*
var swaggerUI embed.FS

func indexSwaggerUI(schemaPath string) Handle {
	const dir = "swagger-ui"
	fs := StaticHttpFileSystem{
		FileSystem: http.FS(swaggerUI),
		Embedded:   true,
	}
	srv := http.FileServer(fs)
	index := fmt.Sprintf(swaggerUIIndexTemplate, schemaPath)

	return func(c *gin.Context) {
		if len(c.Params) == 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		p := path.Join(dir, c.Params[len(c.Params)-1].Value)
		if p == dir {
			// Index.
			_, _ = fmt.Fprint(c.Writer, index)
			return
		}
		// Assets.
		req := c.Request.Clone(c.Request.Context())
		req.URL.Path = p
		req.URL.RawPath = p
		srv.ServeHTTP(c.Writer, req)
		c.Abort()
	}
}

const swaggerUIIndexTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="description" content="SwaggerUI"/>
    <title>SwaggerUI</title>
    <link rel="stylesheet" type="text/css" href="./swagger-ui.css" />
    <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
    <style>
      html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
      *, *:before, *:after { box-sizing: inherit; }
      body { margin: 0; background: #fafafa; }
    </style>
</head>
<body>
<div id="swagger-ui"></div>
<script src="./swagger-ui-bundle.js" charset="UTF-8"></script>
<script>
  window.onload = () => {
    window.ui = SwaggerUIBundle({
      url: '%s',
      validatorUrl: 'none',
      dom_id: '#swagger-ui',
      deepLinking: true,
      docExpansion: 'none',
      displayRequestDuration: true,
      persistAuthorization: true,
      requestInterceptor: (r) => {
        if (r.headers.Cookie) {
          document.cookie = r.headers.Cookie+'; path=/; domain=;'
        }
        return r
      },
      withCredentials: true,
    });
  };
</script>
</body>
</html>
`
