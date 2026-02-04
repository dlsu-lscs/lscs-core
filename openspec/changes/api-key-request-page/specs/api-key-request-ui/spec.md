## ADDED Requirements

### Requirement: Request key page renders with authenticated layout
The request key page SHALL render with the authenticated layout containing navigation and user info.

#### Scenario: Page loads with authenticated layout
- **WHEN** a user navigates to `/request-key`
- **THEN** the page displays the navbar with Dashboard, Profile, and API Keys navigation links

### Requirement: Page handles unauthenticated users via Google OAuth redirect
Unauthenticated users SHALL be redirected to Google OAuth login and returned to the request key page after successful authentication.

#### Scenario: Unauthenticated user is redirected to login
- **WHEN** a user visits `/request-key` without valid authentication
- **THEN** the page redirects to `/auth/google/login?redirect=/request-key`

#### Scenario: User returns from OAuth flow
- **WHEN** Google OAuth completes successfully and redirects back to `/request-key`
- **THEN** the form is displayed for the authenticated user

### Requirement: Page handles unauthorized users
Users without RND AVP+ position SHALL see an access denied message.

#### Scenario: Non-RND member sees access denied
- **WHEN** an authenticated user without RND committee membership or AVP+ position visits `/request-key`
- **THEN** the page displays "Access Denied" message

#### Scenario: RND member without AVP+ sees access denied
- **WHEN** an authenticated RND member with MEM, JO, or CT position visits `/request-key`
- **THEN** the page displays "Access Denied" message

### Requirement: Form accepts project name
Users SHALL be able to enter an optional project name.

#### Scenario: User enters project name
- **WHEN** the user types a project name in the project field
- **THEN** the project name is stored and submitted with the request

#### Scenario: User leaves project name empty
- **WHEN** the user submits without a project name
- **THEN** an empty project name is submitted (backend uses NULL)

### Requirement: Form allows key type selection
Users SHALL be able to select between Development, Production, and Admin key types.

#### Scenario: User selects development key
- **WHEN** the user selects "Development" as the key type
- **THEN** the form shows the allowed origin field with localhost validation hint

#### Scenario: User selects production key
- **WHEN** the user selects "Production" as the key type
- **THEN** the form shows the allowed origin field with production requirements

#### Scenario: User selects admin key
- **WHEN** the user selects "Admin" as the key type
- **THEN** the allowed origin field is hidden (not required for admin keys)

### Requirement: Form validates allowed origin based on key type
The allowed origin field SHALL have different validation rules based on key type.

#### Scenario: Development key with localhost origin
- **WHEN** the user requests a dev key with origin starting with `http://localhost`
- **THEN** the form is valid and can be submitted

#### Scenario: Development key with non-localhost origin
- **WHEN** the user requests a dev key with origin not starting with `http://localhost`
- **THEN** the form shows an error "For dev keys, allowed_origin must start with http://localhost"

#### Scenario: Production key with valid origin
- **WHEN** the user requests a prod key with a valid https:// origin (not localhost)
- **THEN** the form is valid and can be submitted

#### Scenario: Production key with localhost origin
- **WHEN** the user requests a prod key with origin containing localhost
- **THEN** the form shows an error "localhost is not a valid origin for production keys"

#### Scenario: Production key with empty origin
- **WHEN** the user requests a prod key without an allowed origin
- **THEN** the form shows an error "allowed_origin is required for production keys"

### Requirement: Form submits to backend and displays API key
Successfully submitted requests SHALL display the generated API key.

#### Scenario: Form submits successfully
- **WHEN** the user submits a valid form
- **THEN** the API key is displayed in a copyable format

#### Scenario: User copies API key
- **WHEN** the user clicks the copy button
- **THEN** the API key is copied to the clipboard

#### Scenario: API returns error
- **WHEN** the backend returns an error (e.g., origin already exists)
- **THEN** the error message is displayed to the user

### Requirement: Page shows loading states
The page SHALL display appropriate loading indicators during async operations.

#### Scenario: Loading authentication state
- **WHEN** the page is loading
- **THEN** a spinner is displayed

#### Scenario: Submitting form
- **WHEN** the form is being submitted
- **THEN** the submit button is disabled with loading spinner

### Requirement: Page includes back navigation
Users SHALL be able to navigate back to the API keys list.

#### Scenario: User clicks back button
- **WHEN** the user clicks "Back to API Keys"
- **THEN** the page navigates to `/api-keys`
