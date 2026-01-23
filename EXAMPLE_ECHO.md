# Echo example (OpenAPIGO)

This example shows the **config-first** (SpringBoot-like) way to use OpenAPIGO with Echo.

## Run

Non-security:

```bash
go run ./example/echo
```

Security (Bearer + X-API-Key):

```bash
go run -tags "security" ./example/echo
```

Swagger UI:
- http://localhost:8080/swagger-ui/index.html#/

OpenAPI JSON:
- http://localhost:8080/openapi.json

## Code structure

- Spec definition (grouped):
  - `example/echo/main.go`
  - `example/echo/main_security.go`
- Routes + handlers: contained in the same file for Echo examples, but handlers are still separated as functions.

## Upload file

Endpoint:
- `POST /users/upload`
- `POST /secure/users/upload` (security)

Spec uses:

```go
s.POST("/users/upload").MultipartUpload(
  "file",
  openapi.MultipartField{Name: "note", Type: openapi.ParamString},
).Res(map[string]string{}).OK()
```

## Security

- Bearer: `Authorization: Bearer <token>`
- API key: `X-API-Key: <key>`

The security examples also include `/secure/demo-errors` for quick error-response preview.

## Notes

- Examples now use the adapter `NewFromEcho` pattern so you can initialize your Echo instance as usual and wrap it with the adapter before registering OpenAPI.
