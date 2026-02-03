## Why

The current CI/CD pipeline is designed for a single monolithic Docker image that bundles both the Go API and Next.js frontend together. This setup is incompatible with the planned Dokploy deployment strategy, which requires separate Docker images for independent application management. Additionally, the existing workflows lack proper isolation, meaning any code change triggers a full rebuild and deployment of both services, leading to slower deploys and unnecessary resource usage.

## What Changes

- **Create separate Dockerfiles for API and Web services** for independent builds
- **Refactor GitHub Actions workflows** to support monorepo deployment with path-based triggers
- **Add separate Docker image builds** for API and Web with independent versioning
- **Implement Dokploy webhook triggers** for each service to enable selective deployments
- **Add security scanning** for both container images
- **Maintain test coverage** for both backend and frontend code

## Capabilities

### New Capabilities
- **docker-images**: Multi-image build and push pipeline supporting separate API and Web Docker images
- **selective-deployment**: Path-based deployment triggers that only rebuild changed services
- **dokploy-webhooks**: Automated deployment triggers via Dokploy webhooks for each service

### Modified Capabilities
- (None - this is a new infrastructure capability)

## Impact

- **`.github/workflows/`**: Complete refactoring of existing workflows
- **`Dockerfile`**: Rename to `Dockerfile.api`, create new `web/Dockerfile`
- **Build system**: Independent Docker builds for API and Web
- **Deployment**: Separate Dokploy applications with dedicated webhooks
- **CI/CD pipeline**: Path-based filtering for build triggers
