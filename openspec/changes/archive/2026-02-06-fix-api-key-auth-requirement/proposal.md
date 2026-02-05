## Why

Users cannot create API keys through the web UI because the endpoint is protected by `GoogleAuthMiddleware`, which requires an Authorization header with a Bearer token. However, users authenticate via Google OAuth which establishes a session, not a Bearer token. Additionally, the API key creation form includes an unnecessary "admin" option that should be removed.

## What Changes

- Fix the API key creation authentication flow to work with session-based auth instead of requiring Bearer tokens
- Remove the "admin" option from the API key creation form in the web UI
- Ensure RND and AVP+ members can successfully create API keys through the web interface

## Capabilities

### New Capabilities
- `api-key-auth`: Authentication mechanism for API key management endpoints that supports both session-based and Bearer token auth

### Modified Capabilities
- `api-key-management`: Remove the admin option from API key creation UI and update endpoint to accept session-based authentication

## Impact

- **Backend**: Modify `internal/middlewares/google_auth.go` or create a new middleware that supports session-based auth for API key endpoints
- **Backend**: Update API key creation handler to extract user from session instead of Bearer token
- **Frontend**: Remove admin checkbox/option from API key creation form
