## 1. Fix Auth Store Initialization

- [ ] 1.1 Add `onRehydrateStorage` callback to `useAuthStore` in `web/src/lib/auth-store.ts`
- [ ] 1.2 Verify `setLoading(false)` is called after hydration completes

## 2. Verify Root Page Works

- [ ] 2.1 Test root page redirects to `/dashboard` when authenticated
- [ ] 2.2 Test root page redirects to `/login` when not authenticated
- [ ] 2.3 Verify loading spinner shows during hydration
