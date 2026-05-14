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
  ├── cli/            → CLI subcommand implementations (install, upgrade)
  ├── controller/     → HTTP handlers (presentation layer)
  ├── entity/         → Database models (XORM)
  ├── errors/         → Centralized error types and codes
  ├── middleware/     → HTTP middleware (JWT auth)
  ├── migrations/     → Database migrations
  ├── repo/           → Repository layer (data access)
  ├── router/         → API route definitions
  ├── schema/         → Request/Response DTOs
  ├── service/        → Business logic layer
  └── worker/         → Background workers (DDL expiration, email publishing)
pkg/            → Reusable utilities (jwt, uuid, time, pic, etc.)
specs/          → Feature specification documents
```

### Dependency Injection

Wire is used for compile-time DI. Provider sets are defined in `provider.go` files:
- `internal/service/provider.go` - Service layer providers
- `internal/repo/provider.go` - Repository layer providers
- `internal/base/server/provider.go` - Server infrastructure providers
- `internal/router/provider.go` - Router providers
- `internal/controller/provider.go` - Controller providers

After adding new services or repos, add them to the appropriate `ProviderSet` and run `go generate ./cmd/...`.

### Background Workers

The application runs background workers alongside the HTTP server:
- **PublishWorker** - Sends email notifications for upcoming DDLs
- **ExpirationWorker** - Marks DDLs as expired when their deadline passes

Workers are wired in `cmd/wire.go` and started via `app.StartWorkers()`.

### Configuration

Configuration is YAML-based with global singleton access via `conf.Global()`. Default config is embedded at `configs/config.yaml`. At runtime, config is loaded from `<data-dir>/conf/config.yaml`.

All config fields support **environment variable overrides** via `env` struct tags and the `applyEnvVars()` mechanism.

Key config sections:
- `server.http.addr` - HTTP server address
- `data.database` - PostgreSQL connection
- `redis` - Redis connection (`addr`, `password`, `db`)
- `VISUAL_AI` - AI provider for image analysis (GLM)
- `EMAIL_OTP` - SMTP settings for email verification
- `WECHAT` - WeChat Mini Program config (`app_id`, `app_secret`)
- `publish` - Email notification settings (`email.enabled`, `email.smtp`)
- `jwt` - JWT token settings

### Error Handling

The centralized error system lives in `internal/errors/`:

- `ErrorCode` - String constants organized by domain (DDL, User, WeChat, Auth, Infrastructure, HTTP, AI)
- `AppError` - Structured error with `HTTPStatus`, `Code` (machine-readable), `Message` (user-friendly), `Err` (wrapped), `RequestID`, `UserID`, `Operation`, and stack trace
- Sentinel errors for each domain (e.g., `ErrWechatCodeInvalid`, `ErrInvalidRefreshToken`, `ErrEmailAlreadyExists`)
- Factory functions: `DatabaseError()`, `RedisError()`, `ValidationError()`, `NotFoundError()`, `UnauthorizedError()`, `ForbiddenError()`, `AIRequestFailed()`, `AIResponseInvalid()`
- Fluent methods: `.WithRequestID()`, `.WithUserID()`, `.WithOperation()`, `.Wrap()`
- Compatible with `errors.Is`/`errors.As` via `Unwrap()`

The `internal/base/handler` package delegates to `internal/errors`:
- `handler.NewError/BadRequest/NotFound/Conflict/Internal()` - Create AppError via `apperrors.WrapError()`
- `handler.HandleResponse(ctx, err, data)` - Unified response handling
- `handler.BindAndCheck(ctx, &req)` - JSON binding with error response
- `handler.NormalizeError()` - Uses `errors.As` to detect `*AppError`

Response format uses string `code` (e.g., "SUCCESS", "DRAFT_NOT_FOUND") and includes `request_id`.

### Database Migrations

Migrations are defined in `internal/migrations/migrations.go`. Add new migrations using `NewMigration(version, description, migrateFunc)`. Run migrations with `./ddl_guard upgrade`.

### API Documentation

Swagger annotations are in controller files. Access Swagger UI at `/swagger/index.html` when running the server.

### API Route Groups

| Router | Prefix | Description |
|--------|--------|-------------|
| SwaggerRouter | `/swagger/` | API docs |
| AuthApiRouter | `/auth/` | Token refresh |
| UserApiRouter | `/users/` | User auth (email, phone, WeChat login), registration, password |
| DDLApiRouter | `/ddl/` | DDL CRUD and listing |
| ExamApiRouter | `/exams/` | Exam CRUD |
| FinalGradeApiRouter | `/final-grades/` | Final grade CRUD |
| QuizScoreApiRouter | `/final-grades/:uuid/quiz-scores/` | Quiz score CRUD |
| HomeworkScoreApiRouter | `/final-grades/:uuid/homework-scores/` | Homework score CRUD |

## Key Dependencies

- **Gin** - HTTP framework
- **Cobra** - CLI framework
- **Wire** - Dependency injection
- **XORM** - ORM with PostgreSQL
- **go-redis/v9** - Redis client (OTP storage, session management)
- **jwt/v5** - JWT authentication
- **swaggo** - Swagger/OpenAPI generation
