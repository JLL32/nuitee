package main

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed docs/*
var docsFS embed.FS

// serveOpenAPISpec serves the OpenAPI specification file
func (app *application) serveOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	data, err := docsFS.ReadFile("docs/openapi.yaml")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// serveSwaggerUI serves the Swagger UI for interactive API documentation
func (app *application) serveSwaggerUI(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Nuitee API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui.css" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin:0;
            background: #fafafa;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '/docs/openapi.yaml',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                tryItOutEnabled: true,
                supportedSubmitMethods: ['get', 'post', 'put', 'delete', 'patch'],
                docExpansion: 'list',
                filter: true,
                showRequestHeaders: true
            });
        };
    </script>
</body>
</html>`

	t, err := template.New("swagger").Parse(tmpl)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err = t.Execute(w, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// serveReDocUI serves the ReDoc UI as an alternative documentation interface
func (app *application) serveReDocUI(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>Nuitee API Documentation - ReDoc</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
    <style>
        body {
            margin: 0;
            padding: 0;
        }
    </style>
</head>
<body>
    <redoc spec-url='/docs/openapi.yaml'></redoc>
    <script src="https://cdn.jsdelivr.net/npm/redoc@2.1.3/bundles/redoc.standalone.js"></script>
</body>
</html>`

	t, err := template.New("redoc").Parse(tmpl)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err = t.Execute(w, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// serveDocsIndex serves a simple index page with links to different documentation formats
func (app *application) serveDocsIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Nuitee API Documentation</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 40px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            background: white;
            padding: 40px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            margin-bottom: 30px;
        }
        .docs-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin-top: 30px;
        }
        .docs-card {
            border: 1px solid #ddd;
            border-radius: 6px;
            padding: 20px;
            text-decoration: none;
            color: #333;
            transition: box-shadow 0.2s;
        }
        .docs-card:hover {
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
            text-decoration: none;
        }
        .docs-card h3 {
            margin-top: 0;
            color: #0066cc;
        }
        .docs-card p {
            margin-bottom: 0;
            color: #666;
        }
        .version {
            background: #e9ecef;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 0.85em;
            color: #495057;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Nuitee API Documentation <span class="version">v1.0.0</span></h1>
        <p>Welcome to the Nuitee API documentation. Choose your preferred documentation format below:</p>
        
        <div class="docs-grid">
            <a href="/docs/swagger" class="docs-card">
                <h3>Swagger UI</h3>
                <p>Interactive API documentation with the ability to test endpoints directly in your browser.</p>
            </a>
            
            <a href="/docs/redoc" class="docs-card">
                <h3>ReDoc</h3>
                <p>Clean, responsive API documentation with a focus on readability and navigation.</p>
            </a>
            
            <a href="/docs/openapi.yaml" class="docs-card">
                <h3>OpenAPI Spec</h3>
                <p>Raw OpenAPI 3.0 specification in YAML format for integration with other tools.</p>
            </a>
            
            <a href="/docs/simple" class="docs-card">
                <h3>Simple HTML</h3>
                <p>Offline-friendly HTML documentation that works without external dependencies.</p>
            </a>
        </div>
        
        <div style="margin-top: 40px; padding-top: 20px; border-top: 1px solid #eee; color: #666; font-size: 0.9em;">
            <p><strong>API Base URL:</strong> <code>/v1</code></p>
            <p><strong>Health Check:</strong> <a href="/v1/healthcheck">/v1/healthcheck</a></p>
        </div>
    </div>
</body>
</html>`

	t, err := template.New("index").Parse(tmpl)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err = t.Execute(w, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// serveSimpleHTML serves a simple offline HTML documentation page
func (app *application) serveSimpleHTML(w http.ResponseWriter, r *http.Request) {
	data, err := docsFS.ReadFile("docs/simple.html")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
