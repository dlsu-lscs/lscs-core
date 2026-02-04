## Context

The root page (`/`) shows an infinite loading spinner because:

1. `useAuthStore` initializes with `isLoading: true`
2. The Zustand persist middleware hydrates the stored user, but `isLoading` never transitions to `false`
3. `page.tsx` waits for `!isLoading` before redirecting
4. Result: infinite spinner

The store needs to properly handle the hydration complete state.

## Goals / Non-Goals

**Goals:**
- Fix the infinite loading spinner on the root page
- Ensure proper initialization of auth state on page load
- Maintain existing authentication flow (session-based cookies)

**Non-Goals:**
- No changes to the API authentication endpoints
- No changes to the backend auth logic
- No changes to login/logout flows

## Decisions

### Use Zustand's `onRehydrateStorage` callback

Decision: Use Zustand's `onRehydrateStorage` to detect when the store is hydrated and set `isLoading: false`.

```typescript
export const useAuthStore = create<AuthState>()(
    persist(
        (set) => ({
            // ... existing state and actions
        }),
        {
            name: "auth-storage",
            partialize: (state) => ({ user: state.user }),
            onRehydrateStorage: () => (state) => {
                state?.setLoading(false)
            },
        },
    ),
)
```

Rationale: This is the cleanest way to detect when Zustand's persist middleware has finished hydrating from localStorage.

Alternative considered: Using a useEffect with a mounted flag, but onRehydrateStorage is more idiomatic for Zustand.

## Risks / Trade-offs

- **Risk**: The persist middleware might fail to hydrate in some browsers (private browsing, storage disabled)
  - Mitigation: The persist middleware has built-in error handling and will still work

- **Risk**: Loading state might flicker briefly
  - Acceptable: Brief loading state is better than infinite spinner
