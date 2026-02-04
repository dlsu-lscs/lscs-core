## ADDED Requirements

### Requirement: Auth store properly initializes loading state
The auth store SHALL transition `isLoading` to `false` after the persist middleware has hydrated from storage.

#### Scenario: Store hydration completes with no stored user
- **WHEN** the application loads and no user is persisted in storage
- **THEN** `isLoading` is set to `false` and `isAuthenticated` is `false`

#### Scenario: Store hydration completes with stored user
- **WHEN** the application loads and a user is persisted in storage
- **THEN** `isLoading` is set to `false`, `isAuthenticated` is `true`, and `user` contains the persisted user data

### Requirement: Root page redirects based on auth state
The root page SHALL redirect to `/dashboard` if authenticated or `/login` if not authenticated, once the loading state completes.

#### Scenario: Authenticated user is redirected to dashboard
- **WHEN** `isLoading` is `false` and `isAuthenticated` is `true`
- **THEN** the page redirects to `/dashboard`

#### Scenario: Unauthenticated user is redirected to login
- **WHEN** `isLoading` is `false` and `isAuthenticated` is `false`
- **THEN** the page redirects to `/login`

#### Scenario: Loading state prevents premature redirect
- **WHEN** `isLoading` is `true`
- **THEN** the page shows a loading spinner and does not redirect
