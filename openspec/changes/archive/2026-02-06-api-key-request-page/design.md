## Context

The API key request page allows RND committee members with AVP+ positions to generate API keys for external projects. The backend endpoint uses Google OAuth for authentication (not session cookies).

## Goals / Non-Goals

**Goals:**
- Create `/request-key` page with form for API key request
- Support three key types: Development, Production, Admin
- Handle Google OAuth authentication flow
- Display generated API key securely after creation
- Use consistent authenticated layout

**Non-Goals:**
- No backend changes (already implemented)
- No API key management (revoke, list) - existing `/api-keys` page handles this

## Decisions

### Page Structure

The page will use a card-based layout with:
- Project name input (optional)
- Key type selector (radio buttons or select)
- Allowed origin input (required for prod, optional for dev)
- Submit button
- Success state showing the generated API key

### Authentication Flow

Since `/request-key` uses Google Auth (not session):
1. User visits `/request-key`
2. If not authenticated via Google OAuth, redirect to `/auth/google/login?redirect=/request-key`
3. After OAuth callback, user returns to `/request-key`
4. Submit form to create API key

### Form Validation

| Field | Dev Key | Prod Key | Admin Key |
|-------|---------|----------|-----------|
| Project | Optional | Optional | Optional |
| Allowed Origin | localhost only* | Required, no localhost | Not required |

*Dev keys: must start with `http://localhost`

## API Integration

**Endpoint:** POST `/request-key`
**Auth:** Google OAuth (via redirect)

**Request:**
```typescript
{
  project?: string,
  allowed_origin?: string,
  is_dev: boolean,
  is_admin: boolean
}
```

**Response:**
```typescript
{
  email: string,
  api_key: string,
  expires_at?: string
}
```

## UI States

1. **Loading**: Show spinner while checking auth
2. **Unauthenticated**: Redirect to Google login
3. **Unauthorized**: Show "Access Denied" message (not RND AVP+)
4. **Form**: Empty form ready for input
5. **Submitting**: Disable form, show loading
6. **Success**: Show generated API key in a copyable format
7. **Error**: Show error message with retry option

## Risks / Trade-offs

- **Risk**: Google OAuth redirect might lose form state
  - Mitigation: Store form data in sessionStorage before redirect
