## Why

The API Keys page (`/api-keys`) has a "Create API Key" button that links to `/request-key`, but this page doesn't exist yet. Users need a form to request new API keys with project name, allowed origin, and key type (dev/prod/admin).

## What Changes

- Create new page at `web/src/app/request-key/page.tsx`
- Add form with fields: Project name, Allowed Origin, Key Type (dev/prod/admin)
- Submit to backend via POST `/request-key` (Google Auth protected)
- Display generated API key after successful submission
- Use authenticated layout with navbar for consistency

## Capabilities

### New Capabilities
- `api-key-request-ui`: Frontend page for requesting API keys

## Impact

- **New file**: `web/src/app/request-key/page.tsx`
- **Modified**: `web/src/lib/api.ts` - add `requestKey` API function
- **Modified**: `web/src/components/authenticated-layout.tsx` - may need to add navigation
