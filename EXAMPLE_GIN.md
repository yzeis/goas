# Gin example (OpenAPIGO)

This repo uses a **config-first** style (SpringBoot-like): keep routes/handlers clean, and declare OpenAPI metadata in one place using `openapi/simple`.

## Run

Non-security:

```bash
go run -tags gin ./example/gin
```

Security (Bearer + X-API-Key):

```bash
go run -tags "gin,security" ./example/gin
```

Open Swagger UI:

- http://localhost:8080/swagger-ui/index.html#/

OpenAPI JSON:

- http://localhost:8080/openapi.json

## What to look at

- Routing (clean):
  - `example/gin/routes.go`
  - `example/gin/routes_security.go`
- Handlers (separated):
  - `example/gin/handlers_*.go`

## Upload file (multipart/form-data)

Endpoint:
- `POST /users/upload` (non-security)
- `POST /secure/users/upload` (security)

In the spec config it’s defined with:

```go
s.POST("/users/upload").MultipartUpload(
  "file",
  openapi.MultipartField{Name: "note", Type: openapi.ParamString},
).Res(map[string]string{}).OK()
```

In Swagger UI you should see:
- `file` input as file upload
- `note` as text input
- requestBody content type: `multipart/form-data`

## Security

Security demo uses two schemes:
- Bearer token (`Authorization: Bearer <token>`)
- API key (`X-API-Key: <key>`)

See:
- `example/gin/main_security.go`
- `example/gin/routes_security.go`

## Notes

- The examples intentionally do not use the removed typed mode.
- Response error schemas are auto-included via default error responses in `openapi.Config` (unless disabled).
