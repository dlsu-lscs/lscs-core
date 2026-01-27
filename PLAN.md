# LSCS Core - Project Plan

> **Project**: LSCS Core - Member Management & API Key Service  
> **Status**: Phase 3 - Authentication & Session Management  
> **Last Updated**: 2026-01-27

## Overview

LSCS Core is evolving from an "API Key Management" service into a full-featured web application for La Salle Computer Society. The system provides:

1. **API Key Management** - RND officers can request API keys for external projects
2. **Member Management** - Members and officers can manage their profiles
3. **Web UI** - Modern Next.js 16 frontend for all functionality

### Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────┐
│   Next.js 16    │────▶│    Go API       │────▶│   MySQL     │
│   (Frontend)    │     │  (Echo v4)      │     │             │
│   /web          │◀────│  /cmd/api       │◀────│             │
└─────────────────┘     └─────────────────┘     └─────────────┘
                               │
                               ▼
                        ┌─────────────┐
                        │   Garage    │
                        │   (S3)      │
                        └─────────────┘
```

### Tech Stack

**Backend (Go API):**

- Go 1.24, Echo v4
- sqlc (type-safe SQL), goose (migrations)
- zerolog (structured logging)
- go-playground/validator (input validation)
- swaggo/swag (OpenAPI/Swagger)

**Frontend (Next.js):**

- Next.js 16 (App Router)
- TanStack Query (server state)
- Zustand (client state)
- Tailwind CSS

**Infrastructure:**

- MySQL 8.0
- Garage (S3-compatible object storage)
- Dokploy (Docker-based deployment)

---

## Current State Assessment

**Date**: 2026-01-27

### Test Status: FAILING

- `auth/handler_test.go`: Panic - nil interface conversion (user_email not set)
- `member/handler_test.go`: SQL mock column mismatch (expected 8-9, got 18)
- Coverage: 74.1% (middlewares/repository/server at 0%)

### Security Issues

| Issue                  | Severity | Location          | Status  |
| ---------------------- | -------- | ----------------- | ------- |
| CORS allows any origin | HIGH     | `routes.go:19`    | Pending |
| JWT has no expiration  | HIGH     | `auth/service.go` | Pending |
| Panic in health check  | HIGH     | `database.go:78`  | Pending |
| No input validation    | MEDIUM   | All handlers      | Pending |

### Existing Endpoints

| Method | Path           | Auth         | Description                 |
| ------ | -------------- | ------------ | --------------------------- |
| GET    | `/`            | None         | Health check                |
| POST   | `/request-key` | Google OAuth | Generate API key (RND AVP+) |
| GET    | `/members`     | JWT          | List all members            |
| GET    | `/committees`  | JWT          | List all committees         |
| POST   | `/member`      | JWT          | Get member by email         |
| POST   | `/member-id`   | JWT          | Get member by ID            |
| POST   | `/check-email` | JWT          | Check if email is member    |
| POST   | `/check-id`    | JWT          | Check if ID is member       |

---

## Phase 1: Security & Stability Fixes ✅

> **Goal**: Fix all critical security issues and stabilize the test suite before adding features.
> **Status**: COMPLETED (2026-01-27)

### 1.1 Fix CORS Configuration

- [x] Replace wildcard origins with environment-based configuration
- [x] Add `ALLOWED_ORIGINS` env var (comma-separated list)
- [x] Default to restrictive origins in production

**Files**: `internal/server/routes.go`, `.env.example`

### 1.2 Add JWT Expiration

- [x] Add expiration to JWT tokens (configurable via env)
- [x] Dev API keys: 30 days expiration
- [x] Prod API keys: 1 year expiration
- [x] Admin API keys: No expiration (existing behavior)
- [x] Ensure backward compatibility with existing tokens

**Files**: `internal/auth/service.go`, `.env.example`

### 1.3 Remove Panic in Production Code

- [x] Replace `log.Fatalf` with error return in `Health()`
- [x] Add graceful error handling for database connection failures

**Files**: `internal/database/database.go`

### 1.4 Add Input Validation

- [x] Add `go-playground/validator/v10` dependency
- [x] Create validation middleware/helper
- [x] Add validation tags to all request DTOs
- [x] Standardize validation error responses

**Files**: `internal/helpers/validation.go`, all handler DTOs

### 1.5 Fix Broken Tests

- [x] Fix `auth/handler_test.go` - set `user_email` in context
- [x] Fix `member/handler_test.go` - update mock column counts to match query (18 columns)
- [x] Ensure all tests pass before proceeding

**Files**: `internal/auth/handler_test.go`, `internal/member/handler_test.go`

### 1.6 Standardize Error Responses

- [x] Create `APIError` struct for consistent error format
- [x] Create error response helper functions
- [x] Update all handlers to use standardized errors

**Files**: `internal/helpers/errors.go`, all handlers

---

## Phase 2: Foundation Setup ✅

> **Goal**: Set up infrastructure for new features (migrations, logging, docs, config).
> **Status**: COMPLETED (2026-01-27)

### 2.1 Database Migrations with Goose

- [x] Add `pressly/goose` dependency
- [x] Create `migrations/` directory
- [x] Generate initial migration from existing `schema.sql` (baseline)
- [x] Add migration commands to Makefile
- [x] Document migration workflow in AGENTS.md

**Important**: Do NOT drop existing tables. Initial migration should be marked as baseline (already applied).

**Files**: `migrations/`, `Makefile`, `AGENTS.md`

### 2.2 Database Schema Changes Overview

All schema changes across phases are listed here for reference. Each will be implemented as a separate goose migration.

#### New Tables

| Table                   | Phase   | Description                           |
| ----------------------- | ------- | ------------------------------------- |
| `sessions`              | Phase 3 | Web UI session management             |
| `roles`                 | Phase 4 | System-level roles (ADMIN, MODERATOR) |
| `member_roles`          | Phase 4 | Many-to-many: members ↔ roles         |
| `registration_requests` | Phase 8 | Self-registration approval queue      |

#### Schema Modifications

| Table     | Change                       | Phase   | Migration                                                             |
| --------- | ---------------------------- | ------- | --------------------------------------------------------------------- |
| `members` | Add `image_url VARCHAR(512)` | Phase 5 | `ALTER TABLE members ADD COLUMN image_url VARCHAR(512) DEFAULT NULL;` |

#### Complete New Table Schemas

**sessions** (Phase 3):

```sql
CREATE TABLE sessions (
    id VARCHAR(64) PRIMARY KEY,
    member_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_agent VARCHAR(512),
    ip_address VARCHAR(45),
    FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE,
    INDEX idx_sessions_member_id (member_id),
    INDEX idx_sessions_expires_at (expires_at)
);
```

**roles** (Phase 4):

```sql
CREATE TABLE roles (
    id VARCHAR(20) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT
);

-- Seed data
INSERT INTO roles (id, name, description) VALUES
    ('ADMIN', 'Administrator', 'Full system access, can manage all members and settings'),
    ('MODERATOR', 'Moderator', 'Can moderate content and manage basic member issues');
```

**member_roles** (Phase 4):

```sql
CREATE TABLE member_roles (
    member_id INT NOT NULL,
    role_id VARCHAR(20) NOT NULL,
    granted_by INT,
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (member_id, role_id),
    FOREIGN KEY (member_id) REFERENCES members(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (granted_by) REFERENCES members(id) ON DELETE SET NULL
);
```

**registration_requests** (Phase 8 - Future):

```sql
CREATE TABLE registration_requests (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    college VARCHAR(255),
    program VARCHAR(255),
    status ENUM('PENDING', 'APPROVED', 'REJECTED') DEFAULT 'PENDING',
    reviewed_by INT,
    reviewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    notes TEXT,
    FOREIGN KEY (reviewed_by) REFERENCES members(id) ON DELETE SET NULL
);
```

#### Migration Order

```
migrations/
├── 00001_baseline_schema.sql          # Phase 2: Existing schema as baseline
├── 00002_add_sessions_table.sql       # Phase 3: Sessions for web auth
├── 00003_add_roles_tables.sql         # Phase 4: RBAC tables + seed
├── 00004_add_member_image_url.sql     # Phase 5: Profile image support
└── 00005_add_registration_requests.sql # Phase 8: Self-registration (future)
```

### 2.3 Structured Logging with Zerolog

- [x] Replace `log/slog` with `rs/zerolog`
- [x] Add request ID middleware
- [x] Configure log levels via env (`LOG_LEVEL`)
- [x] Add structured context to all log calls

**Files**: `internal/logging/logging.go`, `internal/middlewares/request_logger.go`, all files using `slog`

### 2.4 Configuration Management

- [x] Create `internal/config/config.go`
- [x] Centralize all env var access
- [x] Add validation for required config
- [x] Add config documentation

**Files**: `internal/config/config.go`, `.env.example`

### 2.5 API Documentation (Swagger)

- [x] Add `swaggo/swag` and `swaggo/echo-swagger` dependencies
- [x] Add Swagger annotations to all endpoints
- [x] Serve Swagger UI at `/docs`
- [x] Generate static OpenAPI spec file
- [x] Add `make swagger` command

**Files**: `cmd/api/main.go`, all handlers, `Makefile`, `docs/`

### 2.6 Monorepo Structure

- [x] Create `web/` directory for Next.js frontend
- [x] Update `.gitignore` for Node.js

**Files**: `web/`, `.gitignore`

---

## Phase 3: Authentication & Session Management

> **Goal**: Implement proper auth for web UI while maintaining API key system.

### Authentication Architecture

The system uses **two separate auth mechanisms**:

| Aspect         | Web UI Sessions                    | API Keys (JWT)                        |
| -------------- | ---------------------------------- | ------------------------------------- |
| **For**        | LSCS members using web UI          | External projects (RND) consuming API |
| **Storage**    | `sessions` table                   | `api_keys` table (hash only)          |
| **Transport**  | httpOnly cookie                    | Bearer token header                   |
| **Lifetime**   | 24h (sliding) or 30d (remember me) | 30d (dev) / 1yr (prod)                |
| **Revocation** | Delete session record              | Delete api_key record                 |
| **Refresh**    | Sliding expiration                 | None (request new key)                |

### 3.1 Session Management

- [ ] Create `sessions` table migration (see Phase 2.2 for schema)
- [ ] Create session service (`internal/auth/session.go`)
- [ ] Implement session creation, validation, deletion
- [ ] Implement sliding expiration (extend session if < 50% time remaining)
- [ ] Add session cleanup job (delete expired sessions)
- [ ] Support multiple active sessions per user (different devices)

**Session Configuration**:

- Default duration: 24 hours
- "Remember Me" duration: 30 days
- Sliding expiration threshold: Extend if < 12h remaining (for 24h sessions)
- Cleanup job: Run every hour, delete sessions where `expires_at < NOW()`

**Sliding Expiration Logic**:

```go
// On each authenticated request:
// 1. If (expires_at - now) < (session_duration / 2):
//    - Extend expires_at by session_duration from now
// 2. Update last_activity timestamp
// This prevents constant DB writes while keeping active users logged in
```

**Files**: `migrations/00002_add_sessions_table.sql`, `internal/auth/session.go`

### 3.2 Google OAuth for Web UI

- [ ] Add `/auth/google/login` - redirect to Google (accepts `?remember=true`)
- [ ] Add `/auth/google/callback` - handle OAuth callback
- [ ] Create session on successful login (24h or 30d based on remember flag)
- [ ] Set secure httpOnly cookie with appropriate Max-Age
- [ ] Add `/auth/logout` endpoint (delete session, clear cookie)
- [ ] Add `/auth/me` endpoint (get current user from session)

**Cookie Settings**:

```go
cookie := &http.Cookie{
    Name:     "session_id",
    Value:    sessionID,
    Path:     "/",
    HttpOnly: true,                    // Not accessible via JavaScript
    Secure:   true,                    // HTTPS only (disable for localhost)
    SameSite: http.SameSiteLaxMode,    // CSRF protection
    MaxAge:   86400,                   // 24h (or 30d for remember me)
}
```

**Env vars needed**: `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `OAUTH_REDIRECT_URL`, `SESSION_SECRET`

**Files**: `internal/auth/oauth.go`, `internal/server/routes.go`

### 3.3 Session Middleware

- [ ] Create session validation middleware
- [ ] Extract session ID from cookie
- [ ] Validate session exists and not expired
- [ ] Apply sliding expiration logic
- [ ] Set `user_id` and `user_email` in request context
- [ ] Support both cookie (web UI) and Bearer token (API keys) on same endpoints

**Files**: `internal/middlewares/session.go`

---

## Phase 4: RBAC & Permissions

> **Goal**: Implement role-based access control for member management.

### 4.1 Roles System

- [ ] Create `roles` and `member_roles` tables migration (see Phase 2.2 for schema)
- [ ] Seed initial roles (ADMIN, MODERATOR)

**Files**: `migrations/00003_add_roles_tables.sql`, `internal/auth/rbac.go`

### 4.2 Permission Helpers

- [ ] Create permission checking functions
- [ ] Implement position hierarchy: `PRES > EVP > VP > AVP > CT > JO > MEM`
- [ ] Implement committee-based permissions (VP can edit their committee)
- [ ] Admin role overrides all checks

**Permission Matrix**:
| Actor | Can Edit |
|-------|----------|
| MEM, JO, CT | Own profile only |
| AVP | Own profile only |
| VP | Own profile + AVP, CT, JO, MEM in same committee |
| EVP | Own profile + VP, AVP, CT, JO, MEM (any committee) |
| PRES | Own profile + EVP, VP, AVP, CT, JO, MEM (any committee) |
| ADMIN role | Any member |

**Files**: `internal/helpers/authorization.go`, `internal/helpers/permissions.go`

### 4.3 Authorization Middleware

- [ ] Update existing authorization middleware
- [ ] Add endpoint-specific permission checks
- [ ] Add audit logging for permission checks

**Files**: `internal/middlewares/authorization.go`

---

## Phase 5: Member Management API

> **Goal**: Add CRUD endpoints for member profile management.

### 5.1 Member Profile Endpoints

- [ ] `GET /me` - Get own profile (authenticated user)
- [ ] `PUT /me` - Update own profile
- [ ] `GET /members/:id` - Get member by ID (with permission check)
- [ ] `PUT /members/:id` - Update member (with permission check)

**Editable fields** (by self):

- nickname, telegram, discord, interests, contact_number, fb_link

**Editable fields** (by authorized users):

- All above + full_name, email, position_id, committee_id, college, program, house_id

**Files**: `internal/member/handler.go`, `internal/member/dto.go`

### 5.2 Add image_url Field

- [ ] Create migration to add `image_url` to members table (see Phase 2.2 for schema)
- [ ] Update sqlc queries to include `image_url`
- [ ] Update DTOs and responses (backward compatible - new field appended)

**Files**: `migrations/00004_add_member_image_url.sql`, `query.sql`, `internal/member/dto.go`

---

## Phase 6: Image Upload

> **Goal**: Allow members to upload profile photos to S3-compatible storage (Garage).

### 6.1 S3 Client Setup

- [ ] Add AWS SDK for Go v2
- [ ] Create S3 service (`internal/storage/s3.go`)
- [ ] Configure for Garage endpoint
- [ ] Add env vars: `S3_ENDPOINT`, `S3_BUCKET`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`

**Files**: `internal/storage/s3.go`, `.env.example`

### 6.2 Pre-signed Upload URLs

- [ ] `POST /upload/profile-image` - Generate pre-signed upload URL
    - Validate: auth, file type (JPEG, PNG, WebP), max size (5MB)
    - Return: upload URL, object key, expiration
- [ ] `POST /upload/profile-image/complete` - Confirm upload & update DB
    - Validate: object exists in S3
    - Update `image_url` in members table
    - Delete old image if exists

**Files**: `internal/storage/handler.go`, `internal/server/routes.go`

### 6.3 Image Deletion

- [ ] Delete image from S3 when member updates to new image
- [ ] Delete image from S3 when member is deleted (if applicable)

**Files**: `internal/storage/s3.go`

---

## Phase 7: Next.js Frontend Setup

> **Goal**: Initialize Next.js 16 frontend with modern tooling.

### 7.1 Project Initialization

- [ ] Initialize Next.js 16 in `web/` directory
- [ ] Configure TypeScript
- [ ] Configure Tailwind CSS
- [ ] Configure path aliases

**Files**: `web/`

### 7.2 Core Dependencies

- [ ] Add TanStack Query for server state
- [ ] Add Zustand for client state
- [ ] Add axios or fetch wrapper for API calls
- [ ] Configure API base URL from env

**Files**: `web/package.json`, `web/src/lib/`

### 7.3 Authentication Flow

- [ ] Create auth context/store
- [ ] Implement Google OAuth redirect
- [ ] Handle OAuth callback
- [ ] Implement logout
- [ ] Add auth-protected routes

**Files**: `web/src/lib/auth/`, `web/src/app/`

### 7.4 Core Pages

- [ ] Login page (with "Remember Me" checkbox)
- [ ] Dashboard (landing after login)
- [ ] Profile page (view/edit own profile with image upload)
- [ ] Members list (for authorized users based on position)
- [ ] Member detail/edit (for authorized users based on RBAC)

**Files**: `web/src/app/`

### 7.5 API Key Dashboard (RND Exclusive)

> **Access**: Only visible/accessible to members with `committee_id = 'RND'`

- [ ] API Keys list page (view all user's API keys)
- [ ] Create API key form (project name, allowed origin, dev/prod toggle)
- [ ] Revoke API key functionality
- [ ] Show key details (created date, expiry, origin, last used - if tracked)
- [ ] Copy key to clipboard (on creation only - keys shown once)

**API Endpoints needed**:

- `GET /api-keys` - List user's API keys (without the actual key, just metadata)
- `DELETE /api-keys/:id` - Revoke an API key

**Files**: `web/src/app/api-keys/`, `internal/auth/handler.go`

---

## Phase 8: Registration & Onboarding (Future)

> **Goal**: Allow new LSCS members to self-register with approval workflow.
> **Priority**: Low - Not needed immediately

### 8.1 Registration Flow

- [ ] Create `registration_requests` table
- [ ] `POST /register` - Submit registration request
- [ ] Admin approval workflow
- [ ] Email notification on approval

### 8.2 Onboarding UI

- [ ] Registration form
- [ ] Pending approval status page
- [ ] Welcome flow after approval

---

## Environment Variables

### Required (Current)

```env
# Database
DB_HOST=
DB_PORT=
DB_DATABASE=
DB_USERNAME=
DB_PASSWORD=

# Auth
JWT_SECRET=
GOOGLE_CLIENT_ID=

# Server
PORT=8080
GO_ENV=development
```

### New (To Be Added)

```env
# Auth (new)
GOOGLE_CLIENT_SECRET=
OAUTH_REDIRECT_URL=
SESSION_SECRET=
SESSION_DURATION=86400
SESSION_REMEMBER_DURATION=2592000

# Security (new)
ALLOWED_ORIGINS=http://localhost:3000,https://core.lscs.org

# JWT Expiration (new)
JWT_DEV_EXPIRY_DAYS=30
JWT_PROD_EXPIRY_DAYS=365

# Logging (new)
LOG_LEVEL=info

# S3/Garage (new)
S3_ENDPOINT=
S3_BUCKET=
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_REGION=garage

# Frontend (new, for web/)
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

## Changelog Location

Major changes are logged in `logs/` directory with format:

```
logs/YYYYMMDD-HHMM-<title>.md
```

Each log entry includes:

- Timestamp
- Summary of changes
- Problem/Rationale
- Solution
- Files affected
- Notes/Follow-up

---

## Success Metrics

### Phase 1 Complete When: ✅

- [x] All tests pass
- [x] No security warnings (CORS, JWT expiration)
- [x] No panics in production code paths
- [x] Input validation on all endpoints

### Phase 2 Complete When: ✅

- [x] Goose migrations working
- [x] Zerolog integrated
- [x] Swagger docs accessible at `/docs`
- [x] Static OpenAPI spec generated

### Phase 3 Complete When:

- [ ] Web UI login with Google works
- [ ] Sessions stored in DB
- [ ] Both session (cookie) and API key (JWT) auth work

### Phase 4 Complete When:

- [ ] RBAC tables created
- [ ] Permission checks enforced on all endpoints
- [ ] Position hierarchy respected

### Phase 5 Complete When:

- [ ] Members can edit own profile
- [ ] Authorized users can edit other profiles
- [ ] `image_url` field added (backward compatible)

### Phase 6 Complete When:

- [ ] Profile image upload works
- [ ] Images stored in Garage/S3
- [ ] Old images cleaned up on update

### Phase 7 Complete When:

- [ ] Next.js app running
- [ ] Login/logout working
- [ ] Profile view/edit working
- [ ] Members list working (for authorized users)

---

## References

- [Echo Framework Docs](https://echo.labstack.com/docs)
- [sqlc Documentation](https://docs.sqlc.dev)
- [Goose Migrations](https://github.com/pressly/goose)
- [Swagger/OpenAPI](https://swagger.io/specification/)
- [Next.js 16 Docs](https://nextjs.org/docs)
- [TanStack Query](https://tanstack.com/query)
- [Zustand](https://zustand-demo.pmnd.rs/)
