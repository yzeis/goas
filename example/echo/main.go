//go:build echo && !typed && !security

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	echolib "github.com/labstack/echo/v4"

	"github.com/aizacoders/openapigo/adapters/echo"
	"github.com/aizacoders/openapigo/openapi"
	"github.com/aizacoders/openapigo/openapi/simple"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UpdateUser struct {
	Name string `json:"name"`
}

type CreateUser struct {
	Name string `json:"name"`
}

func main() {
	base := echo.New()

	spec := simple.Spec{
		"GET /users": {Tags: []string{"Users"}, ResSchema: []User{}, Status: http.StatusOK},
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

	r := simple.NewEcho(base, spec)
	users := r.Group("", echo.WithTags("Users"))

	users.GET("/users", func(c echolib.Context) error {
		return echo.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	})

	users.GET("/search", func(c echolib.Context) error {
		_ = c.QueryParam("q")
		return c.NoContent(http.StatusOK)
	})

	users.POST("/users", func(c echolib.Context) error {
		var in CreateUser
		if err := echo.Bind(c, &in); err != nil || in.Name == "" {
			_ = echo.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		return c.NoContent(http.StatusCreated)
	})

	users.GET("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		if id == "404" {
			return echo.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		return echo.JSON(c, http.StatusOK, User{ID: id, Name: "Alice"})
	})

	users.PUT("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		var in UpdateUser
		if err := echo.Bind(c, &in); err != nil {
			_ = echo.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		if id == "404" {
			_ = echo.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return echo.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	})

	users.PATCH("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		var in UpdateUser
		if err := echo.Bind(c, &in); err != nil {
			_ = echo.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		if id == "404" {
			_ = echo.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return echo.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	})

	users.DELETE("/users/:id", func(c echolib.Context) error {
		id := c.Param("id")
		if id == "404" {
			_ = echo.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return c.NoContent(http.StatusNoContent)
	})

	echo.Register(base, openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags:    openapi3.Tags{{Name: "Users", Description: "User management endpoints"}},
	})
	_ = base.Echo.Start(":8080")
}
