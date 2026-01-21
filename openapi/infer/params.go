package infer

import "strings"

// NormalizePath ensures the stored path matches OpenAPI style.
func NormalizePath(path string) string {
	return path
}

// PathParams extracts OpenAPI style path params from /users/{id}.
func PathParams(path string) []string {
	var params []string
	for _, part := range strings.Split(path, "/") {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			name := strings.Trim(part, "{}")
			if name != "" {
				params = append(params, name)
			}
		}
	}
	return params
}
