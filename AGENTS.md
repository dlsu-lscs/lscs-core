# AGENTS.md - LSCS Core (Monorepo)

Guidelines for AI agents working in this Go/Echo + Next.js monorepo.

## Project Overview

This is a **monorepo** with two main components:

### Backend (Go/Echo)
- **Language**: Go 1.24
- **Framework**: Echo v4
- **Database**: MySQL with sqlc for type-safe SQL
- **Module**: `github.com/dlsu-lscs/lscs-core-api`

### Frontend (Next.js)
- **Framework**: Next.js 14+ (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **UI Components**: shadcn/ui, MagicUI
- **State**: TanStack Query + Zustand

**See `PLAN.md`** for project roadmap, phases, and current status.

## Monorepo Structure

```
lscs-core-api/
├── cmd/api/main.go           # Go API entry point
├── internal/                 # Go backend code
│   auth/                     # Authentication handlers
│   committee/                # Committee handlers
│   database/                 # Database connection
│   helpers/                  # Shared utilities
│   member/                   # Member handlers, DTOs
│   middlewares/              # Echo middlewares
│   repository/               # sqlc-generated code
│   server/                   # Server setup
│   shared/                   # Shared types
├── query.sql                 # SQL queries for sqlc
├── schema.sql                # Database schema
├── migrations/               # Goose migrations
├── logs/                     # Change logs
└── web/                      # Next.js frontend
    ├── src/
    │   ├── app/              # Next.js App Router pages
    │   ├── components/       # React components
    │   ├── lib/              # Utilities, API client
    │   └── hooks/            # React hooks
    ├── public/               # Static assets
    └── ...config files
```

## Build/Run/Test Commands

### Backend (Go)

```bash
cd /home/eisen/work/dlsu-y3/lscs/lscs-core-api

# Build and run
make build              # Build to ./main
make run                # Run directly
make watch              # Live reload with Air
make test               # Run all tests

# Single test
go test -v ./internal/member

# Code generation
sqlc generate           # Regenerate repository code

# Database migrations
make migrate-up
make migrate-down
make migrate-status
make migrate-create
```

### Frontend (Next.js)

```bash
cd /home/eisen/work/dlsu-y3/lscs/lscs-core-api/web

# Install dependencies (pnpm is required)
pnpm install

# Development
pnpm dev

# Build for production
pnpm build

# Start production server
pnpm start

# Lint
pnpm lint
```

## Backend Code Style Guidelines

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

## Frontend Code Style Guidelines

### Project Setup

- **Package Manager**: pnpm (required)
- **UI Components**: shadcn/ui (install via CLI)
- **Animations**: MagicUI
- **State Management**: TanStack Query (server), Zustand (client)

### Component Pattern

```tsx
// src/components/ui/button.tsx pattern
import { Button } from "@/components/ui/button"

// src/app/page.tsx pattern
import { Button } from "@/components/ui/button"

export default function Page() {
  return <Button>Click me</Button>
}
```

### API Client Pattern

```tsx
// src/lib/api.ts
const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"

export async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${endpoint}`, {
    headers: {
      "Content-Type": "application/json",
      ...options?.headers,
    },
    ...options,
  })
  if (!res.ok) throw new Error("API request failed")
  return res.json()
}
```

### React Query Pattern

```tsx
// src/hooks/use-member.ts
import { useQuery } from "@tanstack/react-query"
import { fetchAPI } from "@/lib/api"

export function useMember(id: number) {
  return useQuery({
    queryKey: ["member", id],
    queryFn: () => fetchAPI(`/member/${id}`),
  })
}
```

## Important Notes

### Backend
1. **DO NOT edit files in `internal/repository/`** - auto-generated by sqlc
2. **Run `sqlc generate`** after modifying `query.sql` or `schema.sql`
3. **Environment variables** are loaded via `godotenv`
4. **Required env vars**: `DB_*`, `JWT_SECRET`, `GOOGLE_CLIENT_ID`, etc.

### Frontend
1. **Use pnpm** for all package management
2. **shadcn/ui components** go in `src/components/ui/`
3. **Custom components** go in `src/components/`
4. **API routes** are on the Go backend, not Next.js API routes
5. **Environment variables** prefixed with `NEXT_PUBLIC_` are exposed to client

## Database Migrations

See backend section for migration commands.

## Change Logs

Major changes are logged in `logs/` directory:

```
logs/<timestamp>-phase<X>-<descriptive-title>.md
```

Examples:
- `logs/20260127-2200-phase1-fix-security-add-validation.md`
- `logs/20260202-0123-phase6-image-upload.md`
