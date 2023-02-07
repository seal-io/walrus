package swagger

import (
	"embed"
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/runtime"
)

// downloaded form https://github.com/swagger-api/swagger-ui/releases.
//
//go:embed swagger-ui/*
var swaggerUI embed.FS

func Index(openapiURL string) runtime.Handle {
	const dir = "swagger-ui"
	var fs = runtime.StaticHttpFileSystem{
		FileSystem: http.FS(swaggerUI),
		Embedded:   true,
	}
	var srv = http.FileServer(fs)
	var index = fmt.Sprintf(indexTemplate, openapiURL)
	return func(c *gin.Context) {
		if len(c.Params) == 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		var p = path.Join(dir, c.Params[len(c.Params)-1].Value)
		if p == dir {
			// index
			_, _ = fmt.Fprint(c.Writer, index)
			return
		}
		// assets
		var req = c.Request.Clone(c.Request.Context())
		req.URL.Path = p
		req.URL.RawPath = p
		srv.ServeHTTP(c.Writer, req)
		c.Abort()
	}
}

const indexTemplate = `
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
    });
  };
</script>
</body>
</html>
`
