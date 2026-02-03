## ADDED Requirements

### Requirement: Revoke API key endpoint deletes key by ID
The system SHALL provide a `DELETE /api-keys/:id` endpoint that allows authenticated RND members to revoke their own API keys.

#### Scenario: Successfully revoke own API key
- **WHEN** an authenticated RND member sends a DELETE request to `/api-keys/123` where key 123 belongs to them
- **THEN** the system SHALL delete the API key from the database
- **AND** the system SHALL return a 200 OK response with a success message

#### Scenario: Attempt to revoke another user's key
- **WHEN** an authenticated RND member sends a DELETE request to `/api-keys/123` where key 123 belongs to a different user
- **THEN** the system SHALL return a 403 Forbidden response
- **AND** the system SHALL NOT delete the key

#### Scenario: Attempt to revoke non-existent key
- **WHEN** an authenticated RND member sends a DELETE request to `/api-keys/99999` where key 99999 does not exist
- **THEN** the system SHALL return a 404 Not Found response

#### Scenario: Non-RND member attempts to revoke
- **WHEN** an authenticated user who is not in the RND committee sends a DELETE request to `/api-keys/123`
- **THEN** the system SHALL return a 403 Forbidden response

#### Scenario: Unauthenticated request
- **WHEN** an unauthenticated user sends a DELETE request to `/api-keys/123`
- **THEN** the system SHALL return a 401 Unauthorized response
