## Context

The LSCS Core API has two authentication flows:
1. **Session-based** (Web UI) - Uses Google OAuth, stores session in cookies, protects `/auth/*` and `/upload/*` routes
2. **JWT-based** (API keys) - Uses JWT tokens, protects `/members`, `/committees`, `/member`, `/member-id`, `/check-email`, `/check-id` routes

Currently, JWT routes use a deprecated `AuthorizationMiddleware` that:
- Calls `helpers.AuthorizeIfRNDAndAVP()` which queries `GetMemberInfo`
- `GetMemberInfo` includes `m.image_url` column that may not exist in production databases
- Results in `Error 1054: Unknown column 'm.image_url'` errors

The RBAC infrastructure already exists (`RBACService`, RBAC middlewares like `RequireAPIKeyAccess`) but isn't being used for JWT routes.

## Goals / Non-Goals

**Goals:**
- Fix the `image_url` column error by using a lightweight authorization query
- Migrate JWT routes from deprecated `AuthorizationMiddleware` to RBAC-based middlewares
- Maintain backward compatibility - API behavior should remain the same (RND OR AVP+ can access)
- Clean up deprecated code after migration

**Non-Goals:**
- Changing session-based authentication (already working correctly)
- Adding new authorization rules or roles
- Modifying the RBAC service itself (it's already well-designed)

## Decisions

### 1. Add lightweight SQL query for authorization

**Decision:** Create a new `GetMemberAuthInfo` query that only fetches `id`, `position_id`, and `committee_id`.

**Rationale:** 
- The current `GetMemberInfo` query fetches 15+ columns including `image_url`
- Authorization only needs position and committee info
- Avoids dependency on columns that may not exist in older databases
- Better performance (fewer columns to fetch)

**Alternatives considered:**
- Modify `GetMemberInfo` to remove `image_url` - rejected because other code depends on it
- Fix the database migration - helps but doesn't solve the root cause of overfetching

### 2. Extend RBACService with email-based lookup

**Decision:** Add `CanAccessAPIByEmail(ctx, email string)` method to RBACService.

**Rationale:**
- JWT tokens contain email, not member ID
- Current `CanAccessAPIKeyManagement` takes member ID
- Need a method that works with email directly for JWT middleware
- Keeps RBAC logic centralized in RBACService

**Alternatives considered:**
- Lookup member ID from email in middleware then call existing method - adds extra DB call
- Keep authorization in helper function - doesn't leverage existing RBAC infrastructure

### 3. Replace deprecated middleware with RequireAPIAccess

**Decision:** Create `RequireAPIAccess` middleware that uses `RBACService.CanAccessAPIByEmail()`.

**Rationale:**
- JWT routes need email-based authorization (token contains email)
- Follows existing RBAC middleware patterns (`RequireAdmin`, `RequireAPIKeyAccess`, etc.)
- Single middleware replaces `AuthorizationMiddleware`

### 4. Remove deprecated code after migration

**Decision:** Remove `AuthorizationMiddleware` and `AuthorizeIfRNDAndAVP` after migration is complete.

**Rationale:**
- Dead code creates confusion
- Deprecated code may be accidentally used in future
- Cleaner codebase

## Risks / Trade-offs

**[Risk] Database query change could affect existing behavior** → Mitigation: New query returns same authorization fields, just fewer columns. Write tests to verify behavior.

**[Risk] Email lookup in RBAC adds complexity** → Mitigation: Uses existing sqlc-generated queries. Single DB call, same as before.

**[Trade-off] Adding new SQL query** → Acceptable: Small query, necessary for decoupling from `image_url` column.
