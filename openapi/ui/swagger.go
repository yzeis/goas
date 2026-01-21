package ui

import "net/http"

type SwaggerUIConfig struct {
	MountPath   string // default: /swagger
	SpecURLPath string // default: /openapi.json
}

func RegisterSwaggerUI(mux interface {
	Get(string, http.HandlerFunc)
}, cfg SwaggerUIConfig) {
	mount := cfg.MountPath
	if mount == "" {
		mount = "/swagger"
	}
	spec := cfg.SpecURLPath
	if spec == "" {
		spec = "/openapi.json"
	}

	mux.Get(mount, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
<script>
SwaggerUIBundle({
  url: '` + spec + `',
  dom_id: '#swagger-ui'
});
</script>
</body>
</html>
`))
	})
}
