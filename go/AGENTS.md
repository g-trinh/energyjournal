# Global architecture

- Go HTTP service with a small router dispatching to handlers via http.ServeMux.
- Handlers validate JSON and call domain services; map domain errors to HTTP codes and simple JSON bodies.
- Domain/service layer enforces basic group and gifters rules and delegates persistence.
- Storage layer is abstracted;
- Lightweight domain structs (group, participant, rule) with UUID IDs; helper packages for error types

# Open API specifications

Open API specifications are generated thanks to the Swagger annotations in http handlers.

Everytime you create, delete or make changes to an API endpoint, review the Swagger annotations attached to it.
Before committing changes, re-generate the swagger documentation using `swag init`

Always update code before updating Open API specifications.
Never edit Open API specifications without using `swag` binary.

# Templating

## HTML
Use go/template when you need to generate HTML. HTML templates can be found in ../front/tpl.
Inside this folder, you will find templates for all purposes.

If you cannot find a perfect matching template, ask me what you should do.
Do not try to guess.

# Application architecture

This project follows Go's standard project layout conventions.

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

> Infrastructure (AWS CDK) is managed in a separate module at the project root.

# Integrations

Default external service integrations are :
- logrus for logging
