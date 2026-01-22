//go:build security

package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/aizacoders/openapigo/openapi"
	"github.com/aizacoders/openapigo/openapi/simple"
)

type SecUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	base := openapi.NewRouter()

	cfg := openapi.Config{
		Title:   "User API (Security)",
		Version: "1.0.0",
		Tags: openapi3.Tags{
			{Name: "Secure Users", Description: "Secured endpoints (Bearer / X-API-Key)"},
		},
		SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{
			"bearerAuth": {Value: &openapi3.SecurityScheme{Type: "http", Scheme: "bearer", BearerFormat: "JWT"}},
			"apiKeyAuth": {Value: &openapi3.SecurityScheme{Type: "apiKey", In: "header", Name: "X-API-Key"}},
		},
	}

	bearer := openapi3.NewSecurityRequirement().Authenticate("bearerAuth")
	apiKey := openapi3.NewSecurityRequirement().Authenticate("apiKeyAuth")

	spec := simple.Spec{
		"GET /secure/users":  {Tags: []string{"Secure Users"}, Security: &bearer, ResSchema: []SecUser{}, Status: http.StatusOK},
		"POST /secure/users": {Tags: []string{"Secure Users"}, Security: &apiKey, ResSchema: struct{}{}, Status: http.StatusCreated},
	}

	r := simple.New(base, spec)
	secure := r.Group("", openapi.WithTags("Secure Users"))

	// Bearer-protected endpoint
	secure.GET("/secure/users", func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_ = json.NewEncoder(w).Encode([]SecUser{{ID: "1", Name: "Alice"}})
	})

	// API-key-protected endpoint
	secure.POST("/secure/users", func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-API-Key") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	openapi.Register(base, cfg)
	_ = http.ListenAndServe(":8080", r)
}
