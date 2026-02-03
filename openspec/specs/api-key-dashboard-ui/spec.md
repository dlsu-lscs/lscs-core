# api-key-dashboard-ui Specification

## Purpose
TBD - created by archiving change api-key-dashboard. Update Purpose after archive.
## Requirements
### Requirement: API Key Dashboard page accessible to RND members
The system SHALL provide a frontend page at `/api-keys` that displays all API keys for the authenticated RND member.

#### Scenario: RND member views dashboard
- **WHEN** an authenticated RND member navigates to `/api-keys`
- **THEN** the system SHALL display a page showing all their API keys
- **AND** each key SHALL display: project name, allowed origin, type (dev/prod), creation date, expiration date
- **AND** the page SHALL include a "Create API Key" button

#### Scenario: Non-RND member attempts to access dashboard
- **WHEN** an authenticated user who is not in RND attempts to navigate to `/api-keys`
- **THEN** the system SHALL redirect them to the dashboard or show an access denied message

#### Scenario: Dashboard shows empty state
- **WHEN** an authenticated RND member with no API keys navigates to `/api-keys`
- **THEN** the system SHALL display an empty state message
- **AND** the page SHALL still show the "Create API Key" button

### Requirement: Revoke button on dashboard
The system SHALL provide a "Revoke" button for each API key on the dashboard with confirmation.

#### Scenario: User revokes key with confirmation
- **WHEN** an RND member clicks the "Revoke" button next to an API key
- **THEN** the system SHALL display a confirmation dialog
- **AND** **WHEN** the user confirms
- **THEN** the system SHALL call the DELETE endpoint
- **AND** the system SHALL refresh the list to show the key is removed

#### Scenario: User cancels revocation
- **WHEN** an RND member clicks the "Revoke" button
- **THEN** the system SHALL display a confirmation dialog
- **AND** **WHEN** the user cancels
- **THEN** the system SHALL close the dialog without deleting the key

