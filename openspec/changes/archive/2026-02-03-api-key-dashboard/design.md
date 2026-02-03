## Context

The LSCS Core API already has an API key management system through the `/request-key` endpoint. API keys are stored in the `api_keys` table with fields: `api_key_id`, `member_email`, `api_key_hash`, `project`, `allowed_origin`, `is_dev`, `is_admin`, `created_at`, `expires_at`.

Current authentication flow uses session-based auth (httpOnly cookies) for web UI users. The system has RBAC middleware (`RequireRNDAndAVP`) that restricts access to certain endpoints based on committee membership and position level.

The frontend uses Next.js 16 with TanStack Query for server state and Zustand for client state. The project follows the existing pattern of organizing pages under `web/src/app/`.

## Goals / Non-Goals

**Goals:**
- Provide RND committee members visibility into their existing API keys
- Allow RND members to revoke unused or compromised keys
- Integrate seamlessly with existing auth system and UI patterns
- Maintain security by never returning actual key values (only metadata)

**Non-Goals:**
- Edit/modify existing API keys (only create new or revoke)
- Track API key usage statistics or analytics
- Admin view of all users' keys (only own keys visible)
- Bulk operations on multiple keys

## Decisions

### Decision 1: Reuse existing auth pattern (session + RND check)

**Approach:** Use the same `RequireRNDAndAVP` middleware pattern as `/request-key` endpoint.

**Rationale:**
- Consistent with existing API key management security model
- Only RND AVP+ members should manage API keys (business requirement)
- Less code duplication, follows established patterns

**Alternative considered:** Create a new `RequireRND` middleware that allows any RND member. Rejected because we want to maintain the same authorization level as key creation.

### Decision 2: Return full metadata but exclude `api_key_hash`

**Approach:** Query returns all fields from `api_keys` table except the sensitive `api_key_hash`.

**Rationale:**
- Security best practice - actual key values should never be retrievable after creation
- Users already received the key at creation time (`/request-key` returns it once)
- Metadata is useful for identifying which key is which

### Decision 3: Implement confirmation dialog on frontend only

**Approach:** Show a confirmation modal before calling DELETE endpoint, but no backend confirmation flow.

**Rationale:**
- Prevents accidental deletions through misclicks
- Keeps backend simple and stateless
- Follows existing UI patterns in the codebase

### Decision 4: No soft delete - hard delete from database

**Approach:** Use SQL DELETE to permanently remove the API key record.

**Rationale:**
- API keys can be recreated if needed (new key generation)
- No audit trail requirements specified
- Simpler implementation, less database clutter

**Risk:** Accidental deletion is permanent. Mitigated by confirmation dialog and clear UI messaging.

## Risks / Trade-offs

**Risk:** User accidentally revokes a key that's still in use
- **Mitigation:** Confirmation dialog with clear messaging. User can create a new key immediately if needed.

**Risk:** Key enumeration by ID
- **Mitigation:** Only allows deletion of own keys (checked via `member_email` match). No way to list other users' keys.

**Trade-off:** No "undo" for revocation
- **Acceptance:** Recreation is trivial, new key will have new expiration date.

## Migration Plan

**Deployment:**
1. Backend endpoints can be deployed independently of frontend
2. Database: No migrations needed (existing schema sufficient)
3. Frontend page can be deployed after backend is live

**Rollback:**
- Remove routes from `routes.go` to disable endpoints
- Remove frontend page directory

## Open Questions

None - requirements are clear from PLAN.md Phase 7.5 section.
