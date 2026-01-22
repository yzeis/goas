//go:build gin && !typed && !security

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/aizacoders/openapigo/adapters/gin"
	"github.com/aizacoders/openapigo/openapi"
	"github.com/aizacoders/openapigo/openapi/simple"
)

// registerRoutes wires the endpoints in a readable and grouped way.
// (Non-typed, non-security variant.)

func registerSystemRoutes(r *simple.GinRouter) {
	r.GET("/healthz", handleHealthz)
}

func registerUserRoutes(r *simple.GinRouter) {
	users := r.Group("", gin.WithTags("Users"))

	users.GET("/users", handleListUsers)
	users.GET("/search", handleSearchUsers)
	users.POST("/users", handleCreateUser)
	users.GET("/users/:id", handleGetUser)
	users.PUT("/users/:id", handlePutUser)
	users.PATCH("/users/:id", handlePatchUser)
	users.DELETE("/users/:id", handleDeleteUser)
}

func openAPICfg() openapi.Config {
	return openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags: openapi3.Tags{
			{Name: "Users", Description: "User management endpoints"},
		},
	}
}

func springSpec() simple.Spec {
	return simple.Spec{
		"GET /healthz": {Tags: []string{"System"}, ResSchema: map[string]string{}, Status: http.StatusOK},
		"GET /users":   {Tags: []string{"Users"}, ResSchema: []User{}, Status: http.StatusOK},
		"GET /search": {
			Tags: []string{"Users"},
			QueryParams: []openapi.QueryParam{
				{Name: "q", Type: openapi.ParamString, Required: true, Description: "Search term"},
				{Name: "limit", Type: openapi.ParamInteger, Required: false, Description: "Max results"},
			},
			ResSchema: struct{}{},
			Status:    http.StatusOK,
		},
		"POST /users":       {Tags: []string{"Users"}, ReqSchema: CreateUser{}, ResSchema: struct{}{}, Status: http.StatusCreated},
		"GET /users/:id":    {Tags: []string{"Users"}, ResSchema: User{}, Status: http.StatusOK},
		"PUT /users/:id":    {Tags: []string{"Users"}, ReqSchema: UpdateUser{}, ResSchema: User{}, Status: http.StatusOK},
		"PATCH /users/:id":  {Tags: []string{"Users"}, ReqSchema: UpdateUser{}, ResSchema: User{}, Status: http.StatusOK},
		"DELETE /users/:id": {Tags: []string{"Users"}, ResSchema: struct{}{}, Status: http.StatusNoContent},
	}
}
