//go:build !security

package main

import (
	"net/http"

	"github.com/aizacoders/openapigo/adapters/httprouter"
	"github.com/aizacoders/openapigo/openapi"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateUser struct {
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UpdateUser struct {
	Name string `json:"name"`
}

func main() {
	r := httprouter.New()

	// Full auto schema: request schema inferred from CreateUser, response schema from User.
	httprouter.POSTT[CreateUser, User](r, "/users", func(w http.ResponseWriter, req *http.Request, in CreateUser) (User, int, error) {
		_ = req
		return User{ID: "1", Name: in.Name}, http.StatusCreated, nil
	})

	// No request body: use struct{}
	httprouter.GETT[struct{}, []User](r, "/users", func(w http.ResponseWriter, req *http.Request, _ struct{}) ([]User, int, error) {
		_ = w
		_ = req
		return []User{{ID: "1", Name: "Alice"}}, http.StatusOK, nil
	})

	httprouter.GETT[struct{}, User](r, "/users/{id}", func(w http.ResponseWriter, req *http.Request, _ struct{}) (User, int, error) {
		id := openapi.PathValue(req, "id")
		if id == "404" {
			return User{}, http.StatusNotFound, nil
		}
		return User{ID: id, Name: "Alice"}, http.StatusOK, nil
	}, httprouter.WithTags("Users"), httprouter.WithResponses(
		openapi.ResponseSpec{Status: http.StatusOK, Schema: User{}, Description: "OK"},
		openapi.ResponseSpec{Status: http.StatusNotFound, Schema: ErrorResponse{}, Description: "Not Found"},
	))

	httprouter.PUTT[UpdateUser, User](r, "/users/{id}", func(w http.ResponseWriter, req *http.Request, in UpdateUser) (User, int, error) {
		id := openapi.PathValue(req, "id")
		if in.Name == "" {
			return User{}, http.StatusBadRequest, nil
		}
		if id == "404" {
			return User{}, http.StatusNotFound, nil
		}
		return User{ID: id, Name: in.Name}, http.StatusOK, nil
	}, httprouter.WithTags("Users"), httprouter.WithResponses(
		openapi.ResponseSpec{Status: http.StatusOK, Schema: User{}, Description: "OK"},
		openapi.ResponseSpec{Status: http.StatusBadRequest, Schema: ErrorResponse{}, Description: "Bad Request"},
		openapi.ResponseSpec{Status: http.StatusNotFound, Schema: ErrorResponse{}, Description: "Not Found"},
	))

	httprouter.PATCHT[UpdateUser, User](r, "/users/{id}", func(w http.ResponseWriter, req *http.Request, in UpdateUser) (User, int, error) {
		id := openapi.PathValue(req, "id")
		if in.Name == "" {
			return User{}, http.StatusBadRequest, nil
		}
		if id == "404" {
			return User{}, http.StatusNotFound, nil
		}
		return User{ID: id, Name: in.Name}, http.StatusOK, nil
	}, httprouter.WithTags("Users"), httprouter.WithResponses(
		openapi.ResponseSpec{Status: http.StatusOK, Schema: User{}, Description: "OK"},
		openapi.ResponseSpec{Status: http.StatusBadRequest, Schema: ErrorResponse{}, Description: "Bad Request"},
		openapi.ResponseSpec{Status: http.StatusNotFound, Schema: ErrorResponse{}, Description: "Not Found"},
	))

	httprouter.DELETET[struct{}, struct{}](r, "/users/{id}", func(w http.ResponseWriter, req *http.Request, _ struct{}) (struct{}, int, error) {
		id := openapi.PathValue(req, "id")
		if id == "404" {
			return struct{}{}, http.StatusNotFound, nil
		}
		return struct{}{}, http.StatusNoContent, nil
	}, httprouter.WithTags("Users"), httprouter.WithResponses(
		openapi.ResponseSpec{Status: http.StatusNoContent, Schema: struct{}{}, Description: "No Content"},
		openapi.ResponseSpec{Status: http.StatusNotFound, Schema: ErrorResponse{}, Description: "Not Found"},
	))

	openapi.Register(r, openapi.Config{Title: "User API", Version: "1.0.0"})
	_ = http.ListenAndServe(":8080", r)
}
