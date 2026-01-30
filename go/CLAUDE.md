# Claude Context

This is a Go project following the standard Go project layout.

## Directory Structure

```
go/
├── cmd/                    # Entrypoints (main packages)
│   ├── container/          # Docker container server (port 8888)
│   └── lambda/             # AWS Lambda handler
├── internal/               # Private application code
│   ├── server/             # HTTP server, routing, CORS, error mapping
│   ├── handler/            # HTTP endpoint handlers (by feature)
│   ├── domain/             # Core domain models and interfaces
│   ├── service/            # Business logic implementations
│   ├── integration/        # External service adapters (Brevo, Amazon)
│   └── pkg/                # Shared utilities (errors, DynamoDB, AI)
└── docs/                   # Generated Swagger docs (swag init)
```

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
