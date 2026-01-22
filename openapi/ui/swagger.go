package ui

import (
	_ "embed"
	"net/http"
	"strings"
	"text/template"
)

//go:embed templates/swagger-ui.html
var swaggerUITemplate string

//go:embed templates/openapi-stack-logo.png
var openAPIStackLogoPNG []byte

var swaggerUITpl = template.Must(template.New("swagger-ui.html").Parse(swaggerUITemplate))

type SwaggerUIConfig struct {
	MountPath   string // default: /swagger-ui
	SpecURLPath string // default: /openapi.json
}

func RegisterSwaggerUI(mux interface {
	Get(string, http.HandlerFunc)
}, cfg SwaggerUIConfig) {
	mount := cfg.MountPath
	if mount == "" {
		mount = "/swagger-ui"
	}
	if !strings.HasPrefix(mount, "/") {
		mount = "/" + mount
	}
	mount = strings.TrimSuffix(mount, "/")

	// Serve favicon/logo asset
	mux.Get(mount+"/openapi-stack-logo.png", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(openAPIStackLogoPNG)
	})

	spec := cfg.SpecURLPath
	if spec == "" {
		spec = "/openapi.json"
	}

	indexPath := mount + "/index.html"
	redirectHTML := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, indexPath+"#/", http.StatusFound)
	}

	// New canonical paths
	mux.Get(mount, redirectHTML)
	mux.Get(mount+"/", redirectHTML)
	mux.Get(indexPath, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_ = swaggerUITpl.Execute(w, map[string]any{"SpecURL": spec, "MountPath": mount})
	})

	// Legacy: /swagger should redirect to new canonical UI.
	if mount != "/swagger" {
		mux.Get("/swagger", redirectHTML)
		mux.Get("/swagger/", redirectHTML)
	}
}
