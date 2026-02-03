## Context

Currently, HTTP server timeouts in `internal/server/server.go` are hardcoded:
- `IdleTimeout: time.Minute`
- `ReadTimeout: 10 * time.Second`
- `WriteTimeout: 30 * time.Second`

The Go `http.Server` struct accepts `time.Duration` values for these fields. We need to parse duration strings from environment variables and use them when creating the server.

## Goals / Non-Goals

**Goals:**
- Make all three HTTP server timeouts configurable via environment variables
- Maintain backward compatibility with current default values
- Support standard Go duration format (e.g., "30s", "1m", "2h")
- Validate duration values at startup and fail fast on invalid input

**Non-Goals:**
- Dynamic timeout changes at runtime (requires server restart)
- Per-route or per-handler timeout configuration
- Integration with external configuration management systems (Consul, etcd, etc.)

## Decisions

### Decision 1: Use Go's time.ParseDuration for parsing

**Approach:** Use the standard library's `time.ParseDuration()` function to parse duration strings.

**Rationale:**
- Zero external dependencies
- Well-documented format (e.g., "30s", "1m", "2h30m")
- Provides validation out of the box
- Consistent with Go ecosystem conventions

**Alternative considered:** Custom parsing with regex or manual parsing. Rejected because it adds complexity and maintenance burden for no benefit.

### Decision 2: Add three separate timeout fields to Config struct

**Approach:** Add `ServerIdleTimeout`, `ServerReadTimeout`, and `ServerWriteTimeout` fields to the `Config` struct.

**Rationale:**
- Clear, explicit configuration
- Type-safe (all are `time.Duration`)
- Easy to access throughout the application if needed

**Alternative considered:** A single `ServerTimeouts` nested struct. Rejected because it adds unnecessary nesting for just three fields.

### Decision 3: Validate durations during config loading

**Approach:** Parse and validate duration strings in `config.Load()` and return an error if any are invalid.

**Rationale:**
- Fail fast at startup rather than during request handling
- Clear error messages for operators
- Prevents silent misconfiguration

### Decision 4: Keep existing hardcoded values as defaults

**Approach:** If environment variables are not set, use the current hardcoded values as defaults.

**Rationale:**
- Maintains backward compatibility
- Current values are reasonable defaults
- No breaking changes for existing deployments

## Risks / Trade-offs

**Risk:** Invalid duration format causes application startup failure.
- **Mitigation:** This is the desired behavior—operators should fix the configuration. Error message will clearly indicate which variable has the invalid format.

**Risk:** Very short timeouts in production could cause request failures.
- **Mitigation:** Not addressed in this change—operators are responsible for choosing appropriate values. Current defaults are safe.

**Risk:** Very long timeouts could allow slowloris-style attacks.
- **Mitigation:** This is an operational concern, not a code issue. Operators should configure timeouts based on their security requirements.
