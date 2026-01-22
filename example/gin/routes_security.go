//go:build gin && security && !typed

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi"
	"github.com/aizacoders/openapigo/openapi/simple"
)

func openAPICfgSecurity() (openapi.Config, *openapi3.SecurityRequirement, *openapi3.SecurityRequirement) {
	cfg := openapi.Config{
		Title:   "User API (Gin + Security)",
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
	return cfg, &bearer, &apiKey
}

func registerSecureRoutes(r *simple.GinRouter, bearer, apiKey *openapi3.SecurityRequirement) {
	spec := r.Spec
	// enrich existing spec
	spec["GET /secure/healthz"] = simple.RouteDef{
		Tags:      []string{"System"},
		Security:  bearer,
		ResSchema: map[string]string{},
		Status:    http.StatusOK,
	}
	spec["GET /secure/users"] = simple.RouteDef{
		Tags:      []string{"Secure Users"},
		Security:  bearer,
		ResSchema: []SecUser{},
		Status:    http.StatusOK,
	}
	spec["POST /secure/users"] = simple.RouteDef{
		Tags:      []string{"Secure Users"},
		Security:  apiKey,
		ResSchema: struct{}{},
		Status:    http.StatusCreated,
	}
	r.Spec = spec

	r.GET("/secure/healthz", handleSecureHealthz)

	secure := r.Group("", gin.WithTags("Secure Users"))
	secure.GET("/secure/users", handleSecureListUsers)
	secure.POST("/secure/users", handleSecureCreateUser)
}
