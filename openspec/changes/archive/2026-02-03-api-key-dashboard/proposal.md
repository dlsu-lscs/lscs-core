## Why

RND committee members currently have no way to view or manage their API keys after creation. They can only request new keys via `/request-key` but cannot see existing keys, revoke old ones, or track usage. This creates operational friction—members lose track of which keys they've created and have no way to clean up unused or compromised keys.

## What Changes

- Add `GET /api-keys` endpoint to list all API keys for the authenticated user (returns metadata only, not the actual key)
- Add `DELETE /api-keys/:id` endpoint to revoke a specific API key
- Create new frontend page `/api-keys` accessible only to RND committee members
- Add list view showing all user's API keys with project name, allowed origin, creation date, and expiration
- Add "Create API Key" button that navigates to existing key request flow
- Add "Revoke" button for each key with confirmation dialog
- Show key details (created date, expiry, origin, type) but not the actual key value (shown only once at creation)

## Capabilities

### New Capabilities
- `api-key-list`: List all API keys belonging to the authenticated user
- `api-key-revocation`: Revoke/delete a specific API key by ID
- `api-key-dashboard-ui`: Frontend page for RND members to view and manage their API keys

### Modified Capabilities
<!-- No existing capabilities are being modified, only new endpoints added -->

## Impact

- `internal/auth/handler.go`: Add two new handler methods (ListAPIKeys, RevokeAPIKey)
- `internal/server/routes.go`: Register new routes with appropriate middleware
- `query.sql`: Add new SQL queries for listing and deleting API keys by member email
- `web/src/app/api-keys/`: New Next.js page for API key dashboard
- `web/src/lib/api.ts`: Add API client methods for new endpoints
- Database: No schema changes required (using existing `api_keys` table)
