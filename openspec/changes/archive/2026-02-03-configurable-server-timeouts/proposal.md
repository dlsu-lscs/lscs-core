## Why

HTTP server timeouts (IdleTimeout, ReadTimeout, WriteTimeout) are currently hardcoded in `internal/server/server.go`. This prevents operators from tuning server behavior for different deployment scenarios (development vs production, high-latency vs low-latency networks, etc.). Making these configurable via environment variables improves operational flexibility.

## What Changes

- Add three new configuration options in `internal/config/config.go`:
  - `SERVER_IDLE_TIMEOUT` (default: "1m")
  - `SERVER_READ_TIMEOUT` (default: "10s")
  - `SERVER_WRITE_TIMEOUT` (default: "30s")
- Update `internal/server/server.go` to use config values instead of hardcoded constants
- Support Go duration string format (e.g., "30s", "1m", "2h")

## Capabilities

### New Capabilities
- `configurable-http-timeouts`: Allow operators to configure HTTP server timeout values via environment variables

### Modified Capabilities
<!-- No existing capabilities are being modified, only implementation -->

## Impact

- `internal/config/config.go`: Add three new timeout fields and helper function for duration parsing
- `internal/server/server.go`: Replace hardcoded timeouts with config values
- No breaking changes—all values have sensible defaults matching current behavior
