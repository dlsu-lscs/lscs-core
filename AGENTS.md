# AGENTS.md - LSCS Core API

Guidelines for AI agents working in this Go/Echo codebase.

## Project Overview

- **Language**: Go 1.24
- **Framework**: Echo v4 (web framework)
- **Database**: MySQL with sqlc for type-safe SQL generation
- **Module**: `github.com/dlsu-lscs/lscs-core-api`

**See `PLAN.md`** for project roadmap, phases, and current status.

## Build/Run/Test Commands

```bash
make build              # Build to ./main
make run                # Run directly
make watch              # Live reload with Air
make test               # Run all tests with verbose output
make itest              # Run integration tests (./internal/database)
make docker-run         # Start containers
make docker-down        # Stop containers

# Run single test
go test -v ./internal/member                           # All tests in package
go test -v -run TestGetMemberInfo ./internal/member    # Specific test function
go test -v -run TestGetMemberInfo/success ./internal/member  # Specific subtest

# Code generation
sqlc generate           # Regenerate repository code from query.sql

# Database migrations (requires goose CLI)
make migrate-up         # Run all pending migrations
make migrate-down       # Rollback last migration
make migrate-status     # Show migration status
make migrate-create     # Create new migration file
make migrate-baseline   # Mark baseline as applied (for existing DBs)
```

## Project Structure

```
cmd/api/main.go           # Application entry point
internal/
  auth/                   # Authentication handlers & JWT service
  committee/              # Committee handlers
  database/               # Database connection & Service interface
  helpers/                # Shared utilities (NullableString, authorization)
  member/                 # Member handlers, DTOs, types
  middlewares/            # Echo middlewares (auth, JWT)
  repository/             # sqlc-generated code (DO NOT EDIT MANUALLY)
  server/                 # Server setup & route registration
  shared/                 # Shared types
query.sql                 # SQL queries for sqlc
schema.sql                # Database schema
```

## Code Style Guidelines

### Import Organization

Group imports: standard library, external packages, internal packages (blank lines between):

```go
import (
    "database/sql"
    "net/http"

    "github.com/labstack/echo/v4"

    "github.com/dlsu-lscs/lscs-core-api/internal/database"
)
```

### Naming Conventions

- **Packages**: lowercase single words (`member`, `auth`, `helpers`)
- **Exported types**: PascalCase (`Handler`, `EmailRequest`, `Service`)
- **Unexported types**: camelCase (`service`, `mockDBService`)
- **Handler methods**: `VerbNounHandler` or `VerbNoun` (`GetMemberInfo`, `CheckEmailHandler`)
- **JSON fields**: snake_case in struct tags (`json:"full_name"`)
- **Interfaces**: Name by behavior, often ending in `-er` or describe capability

### Handler Pattern

Handlers are structs with a database service dependency:

```go
type Handler struct {
    dbService database.Service
}

func NewHandler(dbService database.Service) *Handler {
    return &Handler{dbService: dbService}
}

func (h *Handler) GetMemberInfo(c echo.Context) error {
    ctx := c.Request().Context()
    q := repository.New(h.dbService.GetConnection())
    // ... handler logic
}
```

### Error Handling

- Use `log/slog` for structured logging
- Return JSON error responses with consistent format
- Check `sql.ErrNoRows` specifically for not-found cases

```go
if err != nil {
    if err == sql.ErrNoRows {
        return c.JSON(http.StatusNotFound, map[string]string{"error": "Not found", "state": "absent"})
    }
    slog.Error("operation failed", "error", err)
    return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
}
```

### DTOs and Responses

- Define request/response types in `dto.go` or `types.go`
- Use conversion functions (`toMemberResponse`, `toFullInfoMemberResponse`)
- Use `helpers.NullableString` for nullable database fields

### Interface-based Design

Define interfaces for dependencies to enable mocking:

```go
type Service interface {
    Health() map[string]string
    Close() error
    GetConnection() *sql.DB
}
```

## Testing

- Use `github.com/stretchr/testify/assert` for assertions
- Use `github.com/DATA-DOG/go-sqlmock` for database mocking
- Use `testcontainers-go` for integration tests
- Write table-driven tests with `t.Run()`

```go
func TestHandler(t *testing.T) {
    t.Run("success", func(t *testing.T) {
        db, mock, _ := sqlmock.New()
        defer db.Close()
        mock.ExpectQuery("SELECT").WillReturnRows(...)
        h := NewHandler(&mockDBService{db: db})
        // Create request and test
    })
}

type mockDBService struct{ db *sql.DB }
func (m *mockDBService) Health() map[string]string { return nil }
func (m *mockDBService) Close() error { return nil }
func (m *mockDBService) GetConnection() *sql.DB { return m.db }
```

## Important Notes

1. **DO NOT edit files in `internal/repository/`** - these are auto-generated by sqlc
2. **Run `sqlc generate`** after modifying `query.sql` or `schema.sql`
3. **Environment variables** are loaded via `github.com/joho/godotenv/autoload`
4. **Required env vars**: `DB_DATABASE`, `DB_PASSWORD`, `DB_USERNAME`, `DB_PORT`, `DB_HOST`, `JWT_SECRET`
5. **Comments**: Start with lowercase, explain "why" not "what"

## Database Migrations

This project uses [Goose](https://github.com/pressly/goose) for database migrations.

### Setup

Install the goose CLI:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Migration Workflow

**For existing databases** (already has schema):

```bash
# mark the baseline migration as applied without running it
make migrate-baseline
```

**For new databases**:

```bash
# run all pending migrations
make migrate-up
```

### Creating New Migrations

```bash
make migrate-create
# enter migration name when prompted (e.g., "add_sessions_table")
```

This creates a new SQL file in `migrations/` with `-- +goose Up` and `-- +goose Down` sections.

### Migration Commands

| Command | Description |
|---------|-------------|
| `make migrate-up` | Run all pending migrations |
| `make migrate-down` | Rollback the last migration |
| `make migrate-status` | Show current migration status |
| `make migrate-create` | Create a new migration file |
| `make migrate-baseline` | Mark baseline as applied (for existing DBs) |

### Migration Best Practices

1. **Never modify applied migrations** - create a new migration instead
2. **Always include a Down migration** - enables rollback
3. **Test migrations locally** before committing
4. **Use descriptive names** (e.g., `add_sessions_table`, `add_member_image_url`)
5. **Keep migrations small** - one logical change per migration

## Change Logs

Major changes are logged in `logs/` directory. Use this naming format:

```
logs/<timestamp>-phase<X>-<descriptive-title>.md
```

- **Timestamp format**: `YYYYMMDD-HHMM`
- **Phase**: Current phase number from PLAN.md
- **Title**: kebab-case summary of what changed

Examples:
- `logs/20260127-2200-phase1-fix-security-add-validation.md`
- `logs/20260127-2330-phase2-add-swagger-docs-monorepo-structure.md`

