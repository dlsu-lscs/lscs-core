## ADDED Requirements

### Requirement: API deployment webhook trigger
The system SHALL trigger the Dokploy API application deployment when the API Docker image is successfully built and pushed.

#### Scenario: Trigger API deployment after successful build
- **WHEN** the workflow `002-build-push-api.yml` completes successfully
- **AND** the image is pushed to `ghcr.io/<org>/<repo>-api:<tag>`
- **THEN** the workflow `004-deploy-api.yml` SHALL execute
- **AND** SHALL send a GET request to the Dokploy webhook URL for the API application
- **AND** the Authorization header SHALL contain the Bearer token from secrets

#### Scenario: Skip API deployment on failure
- **WHEN** the workflow `002-build-push-api.yml` fails (build or push failure)
- **THEN** the workflow `004-deploy-api.yml` SHALL NOT execute
- **AND** no deployment webhook SHALL be sent

### Requirement: Web deployment webhook trigger
The system SHALL trigger the Dokploy Web application deployment when the Web Docker image is successfully built and pushed.

#### Scenario: Trigger Web deployment after successful build
- **WHEN** the workflow `003-build-push-web.yml` completes successfully
- **AND** the image is pushed to `ghcr.io/<org>/<repo>-web:<tag>`
- **THEN** the workflow `005-deploy-web.yml` SHALL execute
- **AND** SHALL send a GET request to the Dokploy webhook URL for the Web application
- **AND** the Authorization header SHALL contain the Bearer token from secrets

#### Scenario: Skip Web deployment on failure
- **WHEN** the workflow `003-build-push-web.yml` fails (build or push failure)
- **THEN** the workflow `005-deploy-web.yml` SHALL NOT execute
- **AND** no deployment webhook SHALL be sent

### Requirement: Independent deployment isolation
The system SHALL ensure that API and Web deployments remain independent.

#### Scenario: API deployment does not affect Web
- **WHEN** an API deployment webhook is triggered
- **THEN** only the API Dokploy application SHALL redeploy
- **AND** the Web Dokploy application SHALL remain at its current version

#### Scenario: Web deployment does not affect API
- **WHEN** a Web deployment webhook is triggered
- **THEN** only the Web Dokploy application SHALL redeploy
- **AND** the API Dokploy application SHALL remain at its current version

### Requirement: Secret management for webhooks
The system SHALL securely store Dokploy webhook URLs and authentication tokens.

#### Scenario: Webhook secrets are GitHub secrets
- **WHEN** deployment workflows need to send webhook requests
- **THEN** the Dokploy webhook URL SHALL be stored in `DOKPLOY_API_WEBHOOK_URL` secret
- **AND** the Dokploy webhook token SHALL be stored in `DOKPLOY_API_TOKEN` secret
- **AND** the Web deployment webhook URL SHALL be stored in `DOKPLOY_WEB_WEBHOOK_URL` secret
- **AND** the Web deployment webhook token SHALL be stored in `DOKPLOY_WEB_TOKEN` secret

#### Scenario: Secrets are not exposed in logs
- **WHEN** deployment workflows execute
- **THEN** webhook URLs and tokens SHALL NOT appear in workflow logs
- **AND** the Authorization header value SHALL be masked in logs
