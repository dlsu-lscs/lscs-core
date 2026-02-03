## 1. Database - Add Lightweight Authorization Query

- [x] 1.1 Add `GetMemberAuthInfo` query to `query.sql` that fetches only `id`, `position_id`, `committee_id` by email
- [x] 1.2 Run `sqlc generate` to regenerate repository code

## 2. RBAC Service - Add Email-Based API Access Check

- [x] 2.1 Add `CanAccessAPIByEmail(ctx, email string) bool` method to `RBACService` in `internal/auth/rbac.go`
- [x] 2.2 Implement the method to: check if RND member OR AVP+ position (using the new lightweight query)

## 3. Middleware - Create RequireAPIAccess Middleware

- [x] 3.1 Add `RequireAPIAccess(rbacService *auth.RBACService)` middleware to `internal/middlewares/authorization.go`
- [x] 3.2 Middleware should extract email from context (`user_email`) and call `RBACService.CanAccessAPIByEmail()`
- [x] 3.3 Return HTTP 403 with `{"error": "Insufficient privileges"}` on failure

## 4. Routes - Migrate JWT Routes to RBAC Middleware

- [x] 4.1 Update `internal/server/routes.go` to replace `AuthorizationMiddleware(s.db)` with `RequireAPIAccess(s.rbacService)` on JWT-protected routes
- [x] 4.2 Verify all JWT routes (`/members`, `/committees`, `/member`, `/member-id`, `/check-email`, `/check-id`) use the new middleware

## 5. Cleanup - Remove Deprecated Code

- [x] 5.1 Remove `AuthorizationMiddleware` function from `internal/middlewares/authorization.go`
- [x] 5.2 Remove `AuthorizeIfRNDAndAVP` function from `internal/helpers/authorization.go`
- [x] 5.3 Remove unused imports from both files

## 6. Testing

- [x] 6.1 Build the project to verify no compilation errors (`make build`)
- [x] 6.2 Test JWT-protected route with valid API key (should work for RND/AVP+)
- [x] 6.3 Verify the `image_url` error no longer occurs
