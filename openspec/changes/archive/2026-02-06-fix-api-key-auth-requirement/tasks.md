## 1. Backend - Fix Authentication for API Key Creation

- [x] 1.1 Change `/request-key` endpoint to use `SessionMiddleware` instead of `GoogleAuthMiddleware` in `internal/server/routes.go`
- [x] 1.2 Verify the handler correctly extracts user email from session context

## 2. Frontend - Remove Admin Option

- [x] 2.1 Remove "Admin" radio button option from `web/src/app/request-key/page.tsx`
- [x] 2.2 Remove `is_admin` field from `RequestKeyRequest` type and form state
- [x] 2.3 Update validation logic to remove admin-specific checks
- [x] 2.4 Remove admin-related UI code (Admin key type handling)

## 3. Testing

- [x] 3.1 Test API key creation from web UI (should succeed with session auth)
- [x] 3.2 Verify admin option is no longer visible in the form
- [x] 3.3 Test that existing Bearer token auth still works for API clients
