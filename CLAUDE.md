# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Run Commands

```bash
# Build the binary
go build -o ddl_guard ./cmd/ddl_guard/

# Initialize configuration and database
./ddl_guard init -d ./data

# Run the HTTP server
./ddl_guard run -d ./data

# Upgrade/migrate database
./ddl_guard upgrade -d ./data

# Run with custom config path
./ddl_guard run -d ./data -c /path/to/config.yaml
```

## Code Generation

```bash
# Regenerate Wire dependency injection code (required after adding new providers)
go generate ./cmd/...

# Regenerate Swagger documentation (required after changing API annotations)
swag init -g cmd/ddl_guard/main.go
```

## Architecture Overview

This is a Go HTTP API server using clean architecture with dependency injection via Google Wire.

### Layer Structure

```
cmd/           → CLI commands (Cobra) and Wire setup
internal/
  ├── base/           → Infrastructure (config, database, server, auth, email)
  ├── controller/     → HTTP handlers (presentation layer)
  ├── entity/         → Database models (XORM)
  ├── middleware/     → HTTP middleware (JWT auth)
  ├── migrations/     → Database migrations
  ├── repo/           → Repository layer (data access)
  ├── router/         → API route definitions
  ├── schema/         → Request/Response DTOs
  └── service/        → Business logic layer
pkg/            → Reusable utilities (jwt, uuid, time, etc.)
```

### Dependency Injection

Wire is used for compile-time DI. Provider sets are defined in `provider.go` files:
- `internal/service/provider.go` - Service layer providers
- `internal/repo/provider.go` - Repository layer providers
- `internal/base/server/provider.go` - Server infrastructure providers
- `internal/router/provider.go` - Router providers

After adding new services or repos, add them to the appropriate `ProviderSet` and run `go generate ./cmd/...`.

### Configuration

Configuration is YAML-based with global singleton access via `conf.Global()`. Default config is embedded at `configs/config.yaml`. At runtime, config is loaded from `<data-dir>/conf/config.yaml`.

Key config sections:
- `server.http.addr` - HTTP server address
- `data.database` - PostgreSQL connection
- `VISUAL_AI` - AI provider for image analysis (GLM)
- `EMAIL_OTP` - SMTP settings for email verification
- `jwt` - JWT token settings

### Error Handling

Use `internal/base/handler` for HTTP responses:
- `handler.NewError(code, message, err)` - Create AppError with HTTP status
- `handler.BadRequest/NotFound/Conflict/Internal()` - Common error constructors
- `handler.HandleResponse(ctx, err, data)` - Unified response handling
- `handler.BindAndCheck(ctx, &req)` - JSON binding with error response

Service errors are defined in respective service packages and mapped to HTTP status in `handler.NormalizeError()`.

### Database Migrations

Migrations are defined in `internal/migrations/migrations.go`. Add new migrations using `NewMigration(version, description, migrateFunc)`. Run migrations with `./ddl_guard upgrade`.

### API Documentation

Swagger annotations are in controller files. Access Swagger UI at `/swagger/index.html` when running the server.

## Key Dependencies

- **Gin** - HTTP framework
- **Cobra** - CLI framework
- **Wire** - Dependency injection
- **XORM** - ORM with PostgreSQL
- **jwt/v5** - JWT authentication
- **swaggo** - Swagger/OpenAPI generation
