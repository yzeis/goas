# go-auto-openapi (openapigo)

Auto-generate OpenAPI 3.0 spec from your Go route registrations.

## Features (current)
- `openapi.Register()` exposes:
  - `GET /openapi.json`
  - `GET /swagger`
- Path param inference for OpenAPI style routes: `/users/{id}`
- Simple schema inference from Go struct tags (request/response samples)
- JWT (and other) security schemes via `openapi.Config.SecuritySchemes`
- Default router (net/http) via `openapi.NewRouter()`
- Optional adapters (build tags): gin / echo / fiber

## Default (net/http) usage
See `example/httprouter/main.go`.

## Adapters
Adapters are behind build tags so the repo stays usable without extra deps.

### Gin
```bash
go run -tags gin ./example/gin
```

### Echo
```bash
go run -tags echo ./example/echo
```

### Fiber
```bash
go run -tags fiber ./example/fiber
```

## Swagger UI
Open in browser:
- `http://localhost:8080/swagger`

## Notes / Roadmap
- Improve schema inference (omitempty, embedded structs, validation tags)
- Add global `servers[]` and richer config (service name, host, base path)
- Add query param inference
- Expand per-route responses (status codes, error schema)

## Full auto schema (Bind/JSON)

Go doesn't allow generic methods on types, so the fully automatic schema mode is exposed as top-level generic helpers:

```go
// No need for WithRequestSchema / WithResponseSchema
openapi.POSTT[CreateUser, User](r, "/users", func(w http.ResponseWriter, req *http.Request, in CreateUser) (User, int, error) {
	return User{ID: "1", Name: in.Name}, http.StatusCreated, nil
})
```

Notes:
- Use `struct{}` as TReq if the endpoint has no request body.
- Use `struct{}` as TRes if the endpoint has no JSON response body.
