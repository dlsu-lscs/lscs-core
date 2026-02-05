## ADDED Requirements

### Requirement: API key endpoints accept session-based authentication

The API key management endpoints SHALL accept authentication via session cookies in addition to Bearer tokens. When a request contains a valid session cookie with user information, the endpoint SHALL extract the user email from the session and proceed with the request.

#### Scenario: Web UI request with valid session

- **WHEN** a request is made to create an API key from the web UI with a valid session cookie
- **THEN** the system SHALL extract the user email from the session
- **AND** the system SHALL create the API key associated with that user email

#### Scenario: API request with Bearer token

- **WHEN** a request is made with a valid Authorization Bearer token header
- **THEN** the system SHALL validate the token and extract the user email
- **AND** the system SHALL proceed with the request as before

#### Scenario: Request with no authentication

- **WHEN** a request is made without a valid session cookie or Bearer token
- **THEN** the system SHALL return a 401 Unauthorized response

### Requirement: Session contains user email

The session established after Google OAuth login SHALL contain the user's email address in a retrievable format.

#### Scenario: Session after OAuth login

- **WHEN** a user successfully authenticates via Google OAuth
- **THEN** the system SHALL store the user's email in the session
- **AND** the session SHALL be retrievable by the API key creation endpoint
