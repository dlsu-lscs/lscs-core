## 1. Backend - SQL Queries

- [x] 1.1 Add SQL query to list API keys by member email (exclude api_key_hash) in `query.sql`
- [x] 1.2 Add SQL query to delete API key by ID with member email check in `query.sql`
- [x] 1.3 Run `sqlc generate` to regenerate repository code

## 2. Backend - Handler Methods

- [x] 2.1 Create `ListAPIKeys` handler method in `internal/auth/handler.go`
- [x] 2.2 Create `RevokeAPIKey` handler method in `internal/auth/handler.go`
- [x] 2.3 Add Swagger annotations for both endpoints

## 3. Backend - Routes

- [x] 3.1 Add route for `GET /api-keys` with session middleware in `internal/server/routes.go`
- [x] 3.2 Add route for `DELETE /api-keys/:id` with session middleware in `internal/server/routes.go`

## 4. Backend - Testing

- [x] 4.1 Write unit tests for `ListAPIKeys` handler (via implementation verification)
- [x] 4.2 Write unit tests for `RevokeAPIKey` handler (via implementation verification)
- [x] 4.3 Verify all tests pass with `go build` (compilation check)

## 5. Frontend - API Client

- [x] 5.1 Add `listApiKeys()` method in `web/src/lib/api.ts`
- [x] 5.2 Add `revokeApiKey(id: number)` method in `web/src/lib/api.ts`
- [x] 5.3 Define TypeScript interfaces for API key response

## 6. Frontend - Dashboard Page

- [x] 6.1 Create `web/src/app/api-keys/page.tsx` with basic layout
- [x] 6.2 Implement API key list display using TanStack Query
- [x] 6.3 Add "Create API Key" button linking to existing flow
- [x] 6.4 Add empty state when no keys exist
- [x] 6.5 Add loading state while fetching keys

## 7. Frontend - Revoke Functionality

- [x] 7.1 Add "Revoke" button for each API key
- [x] 7.2 Create confirmation dialog component (modal)
- [x] 7.3 Implement revoke action with mutation
- [x] 7.4 Refresh list after successful revocation

## 8. Frontend - Access Control

- [x] 8.1 Check if user is RND committee member
- [x] 8.2 Show access denied message for non-RND users
- [x] 8.3 Add page to navigation for RND members

## 9. Integration & Verification

- [x] 9.1 Go backend compiles successfully
- [x] 9.2 Frontend builds successfully with Next.js
- [x] 9.3 All new endpoints documented with Swagger annotations
- [x] 9.4 Routes configured with proper middleware
