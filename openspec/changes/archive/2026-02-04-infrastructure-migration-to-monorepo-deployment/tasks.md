## 1. Docker Infrastructure Setup

- [x] 1.1 Rename existing `Dockerfile` to `Dockerfile.api`
- [x] 1.2 Create `web/Dockerfile` for Next.js frontend
- [x] 1.3 Test both Dockerfiles build independently
- [x] 1.4 Verify API Dockerfile runs correctly locally
- [x] 1.5 Verify Web Dockerfile runs correctly locally

## 2. Dokploy Application Setup

- [x] 2.1 Create API application in Dokploy pointing to `Dockerfile.api`
- [x] 2.2 Create Web application in Dokploy pointing to `web/Dockerfile`
- [x] 2.3 Obtain API webhook URL and token from Dokploy
- [x] 2.4 Obtain Web webhook URL and token from Dokploy
- [x] 2.5 Configure environment variables in Dokploy for each application
- [x] 2.6 Test manual deployment via Dokploy UI

## 3. GitHub Secrets Configuration

- [x] 3.1 Add `DOKPLOY_API_WEBHOOK_URL` to GitHub secrets
- [x] 3.3 Add `DOKPLOY_WEB_WEBHOOK_URL` to GitHub secrets
- [x] 3.5 Verify secrets are accessible in workflows (test run)

## 5. CI/CD Workflow Refactoring - API Build/Push

- [x] 5.1 Create `002-build-push-api.yml` workflow
- [x] 5.2 Add path filters for Go files only
- [x] 5.3 Configure image naming to `ghcr.io/<org>/<repo>-api`
- [x] 5.4 Add Docker Buildx setup
- [x] 5.5 Add GitHub Container Registry login
- [x] 5.6 Add metadata extraction for tagging
- [x] 5.7 Add Trivy security scanning
- [x] 5.8 Add GitHub Actions cache for Docker layers
- [x] 5.9 Test API build workflow

## 6. CI/CD Workflow Refactoring - Web Build/Push

- [x] 6.1 Create `003-build-push-web.yml` workflow
- [x] 6.2 Add path filters for `web/` directory only
- [x] 6.3 Configure image naming to `ghcr.io/<org>/<repo>-web`
- [x] 6.4 Set build context to `web/` directory
- [x] 6.5 Add Docker Buildx setup
- [x] 6.6 Add GitHub Container Registry login
- [x] 6.7 Add metadata extraction for tagging
- [x] 6.8 Add Trivy security scanning
- [x] 6.9 Add GitHub Actions cache for Docker layers
- [x] 6.10 Test Web build workflow

## 7. CI/CD Workflow Refactoring - Deployment Workflows

- [x] 7.1 Create `004-deploy-api.yml` workflow
- [x] 7.2 Trigger on `002-build-push-api.yml` completion
- [x] 7.3 Add Dokploy webhook call with authorization
- [x] 7.4 Create `005-deploy-web.yml` workflow
- [x] 7.5 Trigger on `003-build-push-web.yml` completion
- [x] 7.6 Add Dokploy webhook call with authorization
- [x] 7.7 Test deployment workflows

## 8. Cleanup - Remove Old Workflows

- [x] 8.1 Archive or delete `002-build-push-image.yml` (old monolithic build)
- [x] 8.2 Archive or delete `003-cd.yml` (old Coolify deployment)
- [x] 8.3 Keep `release.yml` (goreleaser - still needed for releases)
- [x] 8.4 Update `README.md` with new deployment documentation

## 9. Integration Testing

- [x] 9.1 Push a test Go change and verify API pipeline triggers
- [x] 9.2 Push a test Web change and verify Web pipeline triggers
- [x] 9.3 Verify API deployment completes in Dokploy
- [x] 9.4 Verify Web deployment completes in Dokploy
- [x] 9.5 Test rollback procedure for both services
- [x] 9.6 Verify security scans run and report results

## 10. Documentation

- [x] 10.1 Update README.md with new deployment documentation (covers build/deploy commands)
- [x] 10.2 Create deployment runbook for Dokploy (covered in README.md)
- [x] 10.3 Document webhook secret management procedures (covered in README.md)
