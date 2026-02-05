## MODIFIED Requirements

### Requirement: API key creation form does not include admin option

The API key creation form in the web UI SHALL NOT include an "admin" checkbox or option. API keys are created with standard member permissions by default.

#### Scenario: Creating API key without admin option

- **WHEN** a user accesses the API key creation form in the web UI
- **THEN** the form SHALL NOT display an "admin" or "administrator" checkbox
- **AND** API keys are created with standard member permissions

#### Scenario: API key creation via API

- **WHEN** an API request is made to create an API key
- **THEN** the system SHALL ignore any "admin" parameter in the request body
- **AND** the API key SHALL be created with standard member permissions

## REMOVED Requirements

### Requirement: API key creation includes admin option

**Reason**: The admin option is not needed for the RND and AVP+ use case. API keys should be created with standard member permissions by default.

**Migration**: Remove the admin checkbox from the frontend form and update any API endpoints to ignore admin parameters in request bodies.
