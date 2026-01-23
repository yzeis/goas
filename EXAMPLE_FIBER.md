# Fiber example (OpenAPIGO)

Fiber example uses the same config-first style with `openapi/simple`.

## Run

Non-security:

```bash
go run ./example/fiber
```

Security:

```bash
go run -tags "security" ./example/fiber
```

Swagger UI:
- http://localhost:8080/swagger-ui/index.html#/

OpenAPI JSON:
- http://localhost:8080/openapi.json

## Upload file

Endpoint:
- `POST /users/upload`
- `POST /secure/users/upload`

In Swagger UI it must show a file chooser for `file`.

## Notes

- Examples now use the adapter `NewFromApp` pattern so you can initialize your Fiber app as usual and wrap it with the adapter before registering OpenAPI.
- Request/response JSON binding uses adapter helpers.
- Spec is declared once and injected into the router:
  - `simple.NewFiber(base, spec)`
