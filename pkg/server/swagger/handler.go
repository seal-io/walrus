package swagger

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/seal-io/utils/httpx"
)

func Index() http.Handler {
	srv := http.FileServer(httpx.FS(ui, httpx.FSOptions().WithEmbedded()))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "", "index", "index.html":
			_, _ = fmt.Fprint(w, index)
		default:
			srv.ServeHTTP(w, r)
		}
	})
}

// downloaded form https://github.com/swagger-api/swagger-ui/releases.
//
//go:embed ui/*
var ui embed.FS

const index = `
<!DOCTYPE html>
<html lang="en">
	<head>
	    <meta charset="utf-8" />
	    <meta name="viewport" content="width=device-width, initial-scale=1" />
	    <meta name="description" content="Walrus OpenAPI Swagger UI"/>
	    <title>Walrus OpenAPI Swagger UI</title>
	    <link rel="icon" type="image/png" href="./ui/favicon-32x32.png" sizes="32x32" />
	    <link rel="icon" type="image/png" href="./ui/favicon-16x16.png" sizes="16x16" />
	    <link rel="stylesheet" type="text/css" href="./ui/swagger-ui.css" />
	    <style>
	      html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
	      *, *:before, *:after { box-sizing: inherit; }
	      body { margin: 0; background: #fafafa; }
	    </style>
	</head>
	<body>
		<div id="swagger-ui"></div>
		<script src="./ui/swagger-ui-bundle.js" charset="UTF-8"></script>
		<script>
		  window.onload = () => {
		    window.ui = SwaggerUIBundle({
		      url: '/openapi/v3/apis/walrus.seal.io/v1',
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
