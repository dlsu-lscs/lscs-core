## Why

The current JWT-protected API routes use a deprecated `AuthorizationMiddleware` that:
1. Has overly restrictive logic (RND members AND AVP+, instead of OR)
2. Uses `GetMemberInfo` query which includes `image_url` column - causing failures when the database migration hasn't been applied
3. Doesn't leverage the existing RBAC infrastructure (`RBACService`, RBAC middlewares)

Additionally, there's a critical database query mismatch: the authorization helper calls `GetMemberInfo` which selects `m.image_url`, but production databases may not have this column if the migration wasn't applied.

## What Changes

- **Remove deprecated `AuthorizationMiddleware`** from JWT-protected routes
- **Create lightweight authorization query** that only fetches `position_id` and `committee_id` (no `image_url`)
- **Migrate JWT routes to use RBAC-based middlewares** (`RequireAPIKeyAccess`, `RequirePosition`, etc.)
- **Remove `AuthorizeIfRNDAndAVP` helper** after migration (no longer needed)
- **Ensure consistent auth flow** for both session-based (Web UI) and JWT (API keys) authentication

## Capabilities

### New Capabilities
- `jwt-rbac-authorization`: RBAC-based authorization for JWT-protected API routes, replacing the deprecated `AuthorizationMiddleware`

### Modified Capabilities

## Impact

- `internal/server/routes.go` - Replace `AuthorizationMiddleware` with RBAC middlewares on JWT routes
- `internal/helpers/authorization.go` - Remove or deprecate `AuthorizeIfRNDAndAVP`
- `internal/middlewares/authorization.go` - Remove deprecated `AuthorizationMiddleware`
- `query.sql` - Add lightweight query for authorization checks (position + committee only)
- `internal/repository/` - Regenerate after query.sql changes
