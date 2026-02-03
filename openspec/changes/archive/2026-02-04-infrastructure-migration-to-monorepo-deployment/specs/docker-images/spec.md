## ADDED Requirements

### Requirement: API service Docker image build
The system SHALL build a Docker image for the API service from `Dockerfile.api` when changes are pushed to Go source files or configuration.

#### Scenario: Build on Go source changes
- **WHEN** changes are pushed to files matching `**/*.go`, `go.mod`, `go.sum`, or `Dockerfile.api`
- **THEN** the GitHub Actions workflow `002-build-push-api.yml` SHALL execute
- **AND** SHALL build a Docker image using the context root directory
- **AND** SHALL push the image to `ghcr.io/<org>/<repo>-api:<tag>`

#### Scenario: API image tagging
- **WHEN** the API image build completes successfully
- **THEN** the image SHALL be tagged with:
  - `ghcr.io/<org>/<repo>-api:sha-<short-sha>` (full commit SHA)
  - `ghcr.io/<org>/<repo>-api:branch-<branch-name>` (for non-main branches)
  - `ghcr.io/<org>/<repo>-api:latest` (for main branch only)

### Requirement: Web service Docker image build
The system SHALL build a Docker image for the Web service from `web/Dockerfile` when changes are pushed to web source files.

#### Scenario: Build on Web source changes
- **WHEN** changes are pushed to files in the `web/` directory or `web/Dockerfile`
- **THEN** the GitHub Actions workflow `003-build-push-web.yml` SHALL execute
- **AND** SHALL build a Docker image using the context `web/` directory
- **AND** SHALL push the image to `ghcr.io/<org>/<repo>-web:<tag>`

#### Scenario: Web image tagging
- **WHEN** the Web image build completes successfully
- **THEN** the image SHALL be tagged with:
  - `ghcr.io/<org>/<repo>-web:sha-<short-sha>` (full commit SHA)
  - `ghcr.io/<org>/<repo>-web:branch-<branch-name>` (for non-main branches)
  - `ghcr.io/<org>/<repo>-web:latest` (for main branch only)

### Requirement: Docker build security scanning
The system SHALL scan each built Docker image for security vulnerabilities.

#### Scenario: API image vulnerability scan
- **WHEN** the API Docker image is built and pushed
- **THEN** Trivy SHALL scan the image for vulnerabilities
- **AND** the scan results SHALL be uploaded as a SARIF artifact
- **AND** the scan SHALL check for OS and library vulnerabilities
- **AND** the scan severity filter SHALL be set to CRITICAL and HIGH
- **AND** the build SHALL continue even if vulnerabilities are found (informational only)

#### Scenario: Web image vulnerability scan
- **WHEN** the Web Docker image is built and pushed
- **THEN** Trivy SHALL scan the image for vulnerabilities
- **AND** the scan results SHALL be uploaded as a SARIF artifact
- **AND** the scan severity filter SHALL be set to CRITICAL and HIGH
- **AND** the build SHALL continue even if vulnerabilities are found (informational only)

### Requirement: Docker build caching
The system SHALL use GitHub Actions cache to optimize Docker build times.

#### Scenario: API image build cache
- **WHEN** the API Docker image is built
- **THEN** the build SHALL use cache from GitHub Actions cache backend
- **AND** the cache mode SHALL be set to `max` to cache all layers

#### Scenario: Web image build cache
- **WHEN** the Web Docker image is built
- **THEN** the build SHALL use cache from GitHub Actions cache backend
- **AND** the cache mode SHALL be set to `max` to cache all layers
