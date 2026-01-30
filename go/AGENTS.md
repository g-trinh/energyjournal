# Global architecture

- Go HTTP service with a small router dispatching to handlers via http.ServeMux.
- Handlers validate JSON and call domain services; map domain errors to HTTP codes and simple JSON bodies.
- Domain/service layer enforces basic group and gifters rules and delegates persistence.
- Storage layer is abstracted; current backing is DynamoDB through AWS SDK v2.
- Lightweight domain structs (group, participant, rule) with UUID IDs; helper packages for error types and a DynamoDB client builder

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

## Directories

- `/cmd`: Application entrypoints (main packages)
  - `/container`: Containerized server bootstrap
  - `/lambda`: AWS Lambda handler entrypoint
- `/internal`: Private application code (compiler-enforced privacy)
  - `/server`: HTTP server setup, routing, error mapping
  - `/handler`: HTTP handlers
    - `/groups`: group/gifter API handlers
  - `/domain`: Core domain models
    - `/group`: Core group/gifter/gift domain models, rules, AI/notifier interfaces; DynamoDB storage interfaces; includes gift tests
      - `/storage`: Domain Repository & Persistence interfaces
    - `/notification`: Notification service interface
  - `/integration`: External services adapters
    - `/amazon`: Amazon gift integration
    - `/brevo`: Brevo emailing adapter
    - `/mailtrap`: Mailtrap emailing adapter
  - `/service`: Domain service-layer implementations and adapters
    - `/groups`: Business logic for groups/gifters/AI, with tests
      - `/storage`: DynamoDB-backed repository implementations and helpers, with tests
        - `/helper`: Helpers for storage (gifter utilities)
  - `/pkg`: Infrastructure/util packages
    - `/ai`: AI client abstraction
      - `/gemini`: Gemini implementation and tests
    - `/dynamodb`: DynamoDB client builder
    - `/error`: Custom error types (validation, notfound, rate-limit, unknown)
    - `/front`: Frontend helper (template hosting/serving)
    - `/secrets`: Secrets loader
- `/docs`: Generated Swagger docs (`swagger.json`, `swagger.yaml`, `docs.go`)

> Infrastructure (AWS CDK) is managed in a separate module at the project root.

# Integrations

Default external service integrations are :
- DynamoDB for storage
- Brevo for mailing
- logrus for logging
