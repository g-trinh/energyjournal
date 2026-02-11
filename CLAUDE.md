# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Energy Journal is a full-stack application for tracking users' energy levels across three dimensions (Physical, Mental, Emotional) on a 0-10 scale. It visualizes energy variations, marking events, and time spending summaries.

**Architecture**: Monorepo with Go backend (`go/`) and React TypeScript frontend (`front/` - git submodule).

## Commands

### Development
```bash
# Backend (from go/)
go run ./cmd/container/main.go    # Start server on :8888

# Frontend (from front/)
npm run dev                        # Start Vite dev server on :8080 (proxies /api to :8888)
```

### Build & Deploy
```bash
# Root level
make build-front                   # Build frontend (npm ci && npm run build)
make deploy                        # Rebase front, build, then CDK deploy
make deploy-no-rebase              # Build and deploy without rebasing front

# Backend
go build ./...
go test ./...

# Frontend
npm run build                      # TypeScript check + Vite build
npm run lint                       # ESLint
```

### Swagger
```bash
swag init                          # Regenerate docs (run after API changes)
```

## Architecture

### Backend (Go)

Layered architecture: `handler → service → domain → storage`

- **cmd/container/**: HTTP server entrypoint (port 8888)
- **cmd/lambda/**: AWS Lambda adapter
- **internal/server/**: Routing (http.ServeMux), CORS, error-to-HTTP mapping
- **internal/handler/**: HTTP handlers by feature (validate input, call services)
- **internal/domain/**: Models and service interfaces
- **internal/service/**: Business logic implementations
- **internal/integration/**: External adapters (Brevo, Amazon)
- **internal/pkg/error/**: Custom error types (`InputValidationError`, `NotFoundError`, `RateLimitError`)

Error handling pattern: domain errors map to HTTP codes in `server/errors.go`.

### Frontend (React)

- React 19 + TypeScript + Vite
- Recharts for visualization
- Path alias: `@/*` → `./src/*`
- Vite proxies `/api/*` to backend at `:8888`

## Key Conventions

- **Go module**: `energyjournal` — imports use `energyjournal/internal/...`
- **Swagger**: Update annotations in handlers, then run `swag init` before committing