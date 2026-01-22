package ui

import (
	"io"
)

// WriteSwaggerUIHTML renders the embedded Swagger UI template into the provided writer.
//
// It intentionally does not register any routes; use RegisterSwaggerUI for that.
func WriteSwaggerUIHTML(w io.Writer, cfg SwaggerUIConfig) {
	spec := cfg.SpecURLPath
	if spec == "" {
		spec = "/openapi.json"
	}

	_ = swaggerUITpl.Execute(w, map[string]any{"SpecURL": spec})
}
