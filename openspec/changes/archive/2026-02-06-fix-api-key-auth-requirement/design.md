## Context

The web UI uses Google OAuth for authentication, establishing a session cookie upon successful login. However, the API key creation endpoint (`internal/middlewares/google_auth.go:16`) is protected by `GoogleAuthMiddleware`, which expects an `Authorization: Bearer <token>` header. This mismatch prevents API key creation from the web UI.

## Goals / Non-Goals

**Goals:**
- Allow authenticated users to create API keys through the web UI without Bearer token requirements
- Remove the unnecessary "admin" option from the API key creation form
- Maintain security by ensuring only authenticated users can create API keys

**Non-Goals:**
- Modify the Google OAuth login flow
- Change authentication for other API endpoints
- Implement a completely new authentication system

## Decisions

1. **Create a session-aware middleware for API key endpoints**

   Instead of modifying `GoogleAuthMiddleware` globally, create a new middleware or modify the API key creation endpoint to accept session-based authentication. The session cookie should contain user information that can be validated.

   Alternative considered: Modify `GoogleAuthMiddleware` to check for session first, but this could affect other endpoints unexpectedly.

2. **Remove admin option from frontend**

   The API key creation form includes an "admin" checkbox that should be removed since it's not needed for the RND/AVP+ use case.

## Risks / Trade-offs

- [Risk] Session validation might differ from Bearer token validation → Ensure session contains sufficient user info (email) for API key creation
- [Risk] Removing admin option might affect existing functionality → Verify admin flag is not used elsewhere in API key creation
