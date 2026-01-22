//go:build fiber && !typed && !security

package main

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	fiberlib "github.com/gofiber/fiber/v2"

	"github.com/aizacoders/openapigo/adapters/fiber"
	"github.com/aizacoders/openapigo/openapi"
	"github.com/aizacoders/openapigo/openapi/simple"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UpdateUser struct {
	Name string `json:"name"`
}

type CreateUser struct {
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	base := fiber.New()

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

	r := simple.NewFiber(base, spec)
	users := r.Group("", fiber.WithTags("Users"))

	users.GET("/users", func(c *fiberlib.Ctx) error {
		return fiber.JSON(c, http.StatusOK, []User{{ID: "1", Name: "Alice"}})
	})

	users.GET("/search", func(c *fiberlib.Ctx) error {
		_ = c.Query("q")
		return c.SendStatus(http.StatusOK)
	})

	users.POST("/users", func(c *fiberlib.Ctx) error {
		var in CreateUser
		if err := fiber.Bind(c, &in); err != nil || in.Name == "" {
			_ = fiber.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		return c.SendStatus(http.StatusCreated)
	})

	users.GET("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		if id == "404" {
			return fiber.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
		}
		return fiber.JSON(c, http.StatusOK, User{ID: id, Name: "Alice"})
	})

	users.PUT("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		var in UpdateUser
		if err := fiber.Bind(c, &in); err != nil {
			_ = fiber.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		if id == "404" {
			_ = fiber.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return fiber.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	})

	users.PATCH("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		var in UpdateUser
		if err := fiber.Bind(c, &in); err != nil {
			_ = fiber.JSON(c, http.StatusBadRequest, ErrorResponse{Error: "invalid body"})
			return nil
		}
		if id == "404" {
			_ = fiber.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return fiber.JSON(c, http.StatusOK, User{ID: id, Name: in.Name})
	})

	users.DELETE("/users/:id", func(c *fiberlib.Ctx) error {
		id := c.Params("id")
		if id == "404" {
			_ = fiber.JSON(c, http.StatusNotFound, ErrorResponse{Error: "user not found"})
			return nil
		}
		return c.SendStatus(http.StatusNoContent)
	})

	fiber.Register(base, openapi.Config{
		Title:   "User API",
		Version: "1.0.0",
		Tags:    openapi3.Tags{{Name: "Users", Description: "User management endpoints"}},
	})
	_ = base.App.Listen(":8080")
}
