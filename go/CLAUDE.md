# Claude Context

This is a Go project following the standard Go project layout.

> Infrastructure (CDK) lives in a separate module at the project root.

## Key Conventions

- **Module**: `energyjournal`
- **Imports**: Use `energyjournal/internal/...` for internal packages
- **Handlers**: Add to `internal/handler/<feature>/`
- **Domain models**: Add to `internal/domain/<entity>/`
- **External adapters**: Add to `internal/integration/<service>/`

## Commands

```bash
go build ./...              # Build all packages
go test ./...               # Run all tests
swag init                   # Regenerate Swagger docs
```

## Stack

- HTTP routing: `net/http` with `http.ServeMux`
- Logging: `logrus`
- Deployment: Container

## See Also

- `AGENTS.md` - Detailed architecture and integration docs
