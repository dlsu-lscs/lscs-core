## Context

The LSCS Core API project is a monorepo containing:
- **Backend**: Go/Echo API with MySQL database
- **Frontend**: Next.js web application

Currently, the CI/CD pipeline builds a single Docker image that contains both services. This approach:
- Forces full rebuilds on any code change
- Prevents independent scaling of API and Web services
- Is incompatible with Dokploy's application model
- Couples deployment schedules between services

## Goals / Non-Goals

**Goals:**
- Enable independent Docker builds for API and Web services
- Implement path-based filtering to only rebuild changed services
- Support Dokploy deployment with separate applications
- Maintain or improve security scanning coverage
- Preserve existing test workflows

**Non-Goals:**
- Modify application code or business logic
- Change the database schema or migration strategy
- Implement blue/green or canary deployments (future work)
- Migrate from GitHub Actions to another CI/CD provider

## Decisions

### Decision 1: Separate Dockerfiles vs Multi-stage Build

**Chosen**: Separate Dockerfiles (`Dockerfile.api` and `web/Dockerfile`)

**Rationale**:
- Clear separation of concerns for builds
- Independent caching - changes in Web don't invalidate API build cache
- Dokploy can point directly to each Dockerfile
- Simpler build context for each service

**Alternatives Considered**:
- Single multi-stage build: Would still require full rebuild on any change
- Docker Compose build: Requires Dokploy to support compose natively

### Decision 2: GitHub Actions Workflow Structure

**Chosen**: Independent workflows per service with shared test workflow

**Rationale**:
- Clear trigger conditions based on path changes
- Independent parallel execution possible
- Easier maintenance and debugging
- Each workflow can be disabled independently

**Workflow Structure**:
1. `001-test.yml` - Runs on any relevant change, no deployment
2. `002-build-push-api.yml` - Triggers on Go changes
3. `003-build-push-web.yml` - Triggers on Web changes
4. `004-deploy-api.yml` - Triggers on API image push
5. `005-deploy-web.yml` - Triggers on Web image push

### Decision 3: Image Naming Convention

**Chosen**: `ghcr.io/<org>/<repo>-<service>:<tag>`

**Examples**:
- `ghcr.io/dlsu-lscs/lscs-core-api-api:main`
- `ghcr.io/dlsu-lscs/lscs-core-api-web:latest`

**Rationale**:
- Clear identification of service
- Prevents image name collisions
- Each service has independent tag history

### Decision 4: Deployment Trigger Strategy

**Chosen**: Separate webhook URLs for each service

**Rationale**:
- Dokploy supports webhooks per application
- Enables independent deployment scheduling
- Failure in one service doesn't block the other
- Can implement different deployment strategies per service

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|------|--------|------------|
| Increased GitHub Actions minutes | Higher CI/CD costs | Path-based filtering minimizes unnecessary runs |
| More workflow files to maintain | Operational complexity | Clear naming and documentation |
| Network issues with webhook | Deployment failure | Retry mechanism, manual deploy option |
| Image versioning drift | Incompatible deployments | Use consistent tagging strategy |

## Migration Plan

### Phase 1: Infrastructure Preparation
1. Create new Dockerfiles for API and Web
2. Set up two Dokploy applications with webhook URLs
3. Create GitHub repository secrets for webhook URLs

### Phase 2: CI/CD Migration
1. Refactor existing workflows to use new structure
2. Add path-based filtering to build workflows
3. Add webhook deployment steps
4. Test builds independently

### Phase 3: Cutover
1. Update Dokploy to use new image names
2. Disable old monolithic workflow
3. Deploy first changes via new pipeline
4. Monitor for issues, rollback plan ready

### Rollback Plan
1. Re-enable old monolithic workflow
2. Update Dokploy to use old image
3. Revert image references in workflows if needed

## Open Questions

1. **Should we use Docker layer caching across services?** Currently each Dockerfile is independent. Could share a base image layer.
2. **Image tag strategy for PRs?** Need to decide if PR deployments should push images and how to tag them.
3. **VPS resource allocation?** Dokploy on self-hosted VPS - ensure adequate resources for two separate containers.
