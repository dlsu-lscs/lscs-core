## ADDED Requirements

### Requirement: List API keys endpoint returns user's keys
The system SHALL provide a `GET /api-keys` endpoint that returns all API keys belonging to the authenticated user.

#### Scenario: Successfully list API keys
- **WHEN** an authenticated RND member sends a GET request to `/api-keys`
- **THEN** the system SHALL return a list of all API keys created by that user
- **AND** each key SHALL include: api_key_id, project name, allowed_origin, is_dev, is_admin, created_at, expires_at
- **AND** the actual api_key_hash SHALL NOT be returned (only metadata)
- **AND** the response SHALL be in JSON format

#### Scenario: Non-RND member attempts to access
- **WHEN** an authenticated user who is not in the RND committee sends a GET request to `/api-keys`
- **THEN** the system SHALL return a 403 Forbidden response
- **AND** the response SHALL include an error message indicating insufficient permissions

#### Scenario: Unauthenticated request
- **WHEN** an unauthenticated user sends a GET request to `/api-keys`
- **THEN** the system SHALL return a 401 Unauthorized response

#### Scenario: Empty list when no keys exist
- **WHEN** an authenticated RND member with no API keys sends a GET request to `/api-keys`
- **THEN** the system SHALL return an empty array
