## ADDED Requirements

### Requirement: Path-based workflow triggering
The CI/CD system SHALL trigger builds only for services affected by the changed files.

#### Scenario: API-only changes trigger API workflow
- **WHEN** changes are pushed that modify only Go source files (`**/*.go`, `go.mod`, `go.sum`)
- **AND** no changes are made to files in the `web/` directory
- **THEN** the workflow `002-build-push-api.yml` SHALL trigger
- **AND** the workflow `003-build-push-web.yml` SHALL NOT trigger

#### Scenario: Web-only changes trigger Web workflow
- **WHEN** changes are pushed that modify only files in the `web/` directory
- **AND** no changes are made to Go source files
- **THEN** the workflow `003-build-push-web.yml` SHALL trigger
- **AND** the workflow `002-build-push-api.yml` SHALL NOT trigger

#### Scenario: Both services changed trigger both workflows
- **WHEN** changes are pushed that modify both Go source files and files in the `web/` directory
- **THEN** both workflow `002-build-push-api.yml` and `003-build-push-web.yml` SHALL trigger
- **AND** both workflows SHALL execute independently

### Requirement: Configuration file changes trigger all builds
The CI/CD system SHALL trigger builds for all services when configuration files are modified.

#### Scenario: Root configuration changes
- **WHEN** changes are pushed to `Dockerfile`, `docker-compose.yml`, or other root-level configuration files
- **THEN** both API and Web workflows SHALL trigger
- **AND** both services SHALL be rebuilt with the updated configuration

### Requirement: Test workflow execution
The test workflow SHALL run on any relevant changes to validate code quality.

#### Scenario: Test on API changes
- **WHEN** changes are pushed to Go source files or `go.mod`/`go.sum`
- **THEN** the workflow `001-test.yml` SHALL execute
- **AND** SHALL build the Go application using `make build`
- **AND** SHALL NOT deploy any images

#### Scenario: Test on Web changes
- **WHEN** changes are pushed to files in the `web/` directory
- **THEN** the workflow `001-test.yml` SHALL execute
- **AND** SHOULD run Web-specific tests (linting, type checking)
- **AND** SHALL NOT deploy any images

#### Scenario: Test on any PR to main
- **WHEN** a pull request is opened or updated targeting the main branch
- **THEN** the workflow `001-test.yml` SHALL execute for all relevant changes
- **AND** the test results SHALL be reported on the PR
