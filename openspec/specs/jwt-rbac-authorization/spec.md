### Requirement: Lightweight authorization query
The system SHALL provide a SQL query `GetMemberAuthInfo` that fetches only `id`, `position_id`, and `committee_id` for a member by email, without depending on optional columns like `image_url`.

#### Scenario: Query returns auth info for existing member
- **WHEN** `GetMemberAuthInfo` is called with a valid member email
- **THEN** the system returns the member's `id`, `position_id`, and `committee_id`

#### Scenario: Query returns error for non-existent member
- **WHEN** `GetMemberAuthInfo` is called with an email that doesn't exist in the database
- **THEN** the system returns `sql.ErrNoRows`

### Requirement: RBAC service supports email-based API access check
The system SHALL provide a method `CanAccessAPIByEmail(ctx, email)` on `RBACService` that checks if a member (identified by email) can access JWT-protected API routes.

#### Scenario: RND member can access API
- **WHEN** `CanAccessAPIByEmail` is called for a member in the RND committee (any position)
- **THEN** the method returns `true`

#### Scenario: AVP+ member can access API
- **WHEN** `CanAccessAPIByEmail` is called for a member with AVP, VP, EVP, or PRES position (any committee)
- **THEN** the method returns `true`

#### Scenario: Regular member cannot access API
- **WHEN** `CanAccessAPIByEmail` is called for a member with MEM, JO, or CT position in a non-RND committee
- **THEN** the method returns `false`

#### Scenario: Non-member cannot access API
- **WHEN** `CanAccessAPIByEmail` is called for an email that is not in the members table
- **THEN** the method returns `false`

### Requirement: JWT routes use RBAC middleware
The system SHALL protect JWT-authenticated routes with `RequireAPIAccess` middleware that uses `RBACService.CanAccessAPIByEmail()`.

#### Scenario: Authorized request passes middleware
- **WHEN** a JWT-authenticated request is made by an RND member or AVP+ position holder
- **THEN** the middleware allows the request to proceed to the handler

#### Scenario: Unauthorized request is rejected
- **WHEN** a JWT-authenticated request is made by a member without RND membership or AVP+ position
- **THEN** the middleware returns HTTP 403 Forbidden with error message "Insufficient privileges"

#### Scenario: Invalid JWT is rejected
- **WHEN** a request is made with an invalid or expired JWT token
- **THEN** the JWT middleware returns HTTP 401 Unauthorized (before reaching RBAC middleware)

### Requirement: Deprecated middleware removed
The system SHALL NOT use the deprecated `AuthorizationMiddleware` or `AuthorizeIfRNDAndAVP` helper for authorization.

#### Scenario: Deprecated code is removed
- **WHEN** the codebase is inspected after migration
- **THEN** `AuthorizationMiddleware` in `internal/middlewares/authorization.go` is removed
- **AND** `AuthorizeIfRNDAndAVP` in `internal/helpers/authorization.go` is removed
