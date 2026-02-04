## Why

The root page (`/`) is stuck in a constant loading state because the `useAuthStore` initializes with `isLoading: true` but never transitions to `false` on initial load. This creates an infinite loading spinner that prevents users from being redirected to the appropriate page (dashboard or login).

## What Changes

- Fix `useAuthStore` to properly handle initial hydration state
- Add session initialization mechanism that sets `isLoading: false` after store hydration
- Ensure persisted auth state is correctly restored on page load

## Capabilities

### New Capabilities
- `auth-store-initialization`: Proper initialization and hydration handling for the auth store

## Impact

- **Modified**: `web/src/lib/auth-store.ts` - Fix initial loading state handling
- **Modified**: `web/src/hooks/use-auth.ts` - May need updates to work with fixed store
- **Modified**: `web/src/app/page.tsx` - Root page will work correctly after store fix
