# configurable-http-timeouts Specification

## Purpose

TBD - created by archiving change configurable-server-timeouts. Update Purpose after archive.

## Requirements

### Requirement: Server idle timeout is configurable

The system SHALL allow operators to configure the HTTP server's idle timeout via the SERVER_IDLE_TIMEOUT environment variable.

#### Scenario: Idle timeout uses default value

- **WHEN** the SERVER_IDLE_TIMEOUT environment variable is not set
- **THEN** the server SHALL use a default idle timeout of 1 minute

#### Scenario: Idle timeout uses custom value

- **WHEN** the SERVER_IDLE_TIMEOUT environment variable is set to "2m"
- **THEN** the server SHALL use an idle timeout of 2 minutes

### Requirement: Server read timeout is configurable

The system SHALL allow operators to configure the HTTP server's read timeout via the SERVER_READ_TIMEOUT environment variable.

#### Scenario: Read timeout uses default value

- **WHEN** the SERVER_READ_TIMEOUT environment variable is not set
- **THEN** the server SHALL use a default read timeout of 10 seconds

#### Scenario: Read timeout uses custom value

- **WHEN** the SERVER_READ_TIMEOUT environment variable is set to "30s"
- **THEN** the server SHALL use a read timeout of 30 seconds

### Requirement: Server write timeout is configurable

The system SHALL allow operators to configure the HTTP server's write timeout via the SERVER_WRITE_TIMEOUT environment variable.

#### Scenario: Write timeout uses default value

- **WHEN** the SERVER_WRITE_TIMEOUT environment variable is not set
- **THEN** the server SHALL use a default write timeout of 30 seconds

#### Scenario: Write timeout uses custom value

- **WHEN** the SERVER_WRITE_TIMEOUT environment variable is set to "60s"
- **THEN** the server SHALL use a write timeout of 60 seconds

### Requirement: Timeout values use Go duration format

The system SHALL accept timeout values in Go's duration string format.

#### Scenario: Valid duration formats are accepted

- **WHEN** the operator sets SERVER_READ_TIMEOUT to "1h30m"
- **THEN** the server SHALL parse this as 1 hour and 30 minutes

#### Scenario: Invalid duration formats are rejected

- **WHEN** the operator sets SERVER_READ_TIMEOUT to "invalid"
- **THEN** the system SHALL return an error during configuration loading
- **AND** the error message SHALL indicate the invalid duration format
