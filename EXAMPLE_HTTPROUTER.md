# net/http (chi) router example (OpenAPIGO)

The “default” router in this repo is `openapi.Router` (built on top of `chi`).

## Run

Non-security:

```bash
go run ./example/httprouter
```

Security:

```bash
go run -tags security ./example/httprouter
```

Swagger UI:
- http://localhost:8080/swagger-ui/index.html#/

OpenAPI JSON:
- http://localhost:8080/openapi.json

## Upload file

Endpoint:
- `POST /users/upload`
- `POST /secure/users/upload`

Uses the same `MultipartUpload(...)` helper in the spec builder.

## Why "httprouter" directory name?

Historical naming: originally this example used another router.
Now it uses the built-in `openapi.Router` (chi-based) but we keep the folder name for compatibility.

## Notes

- Examples follow the pattern: build a router/engine, wrap with adapter (when applicable), build spec via `simple.NewSpec()` and then use `simple.New*` wrappers.
