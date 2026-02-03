## 1. Add timeout fields to Config struct

- [x] 1.1 Add `ServerIdleTimeout`, `ServerReadTimeout`, and `ServerWriteTimeout` fields to the `Config` struct in `internal/config/config.go`
- [x] 1.2 Add helper function `getEnvDuration()` for parsing duration strings from environment variables
- [x] 1.3 Initialize the three timeout fields in `Load()` using `getEnvDuration()` with appropriate defaults

## 2. Update server initialization to use config timeouts

- [x] 2.1 Replace hardcoded `time.Minute` with `cfg.ServerIdleTimeout` in `internal/server/server.go`
- [x] 2.2 Replace hardcoded `10 * time.Second` with `cfg.ServerReadTimeout`
- [x] 2.3 Replace hardcoded `30 * time.Second` with `cfg.ServerWriteTimeout`

## 3. Verify implementation

- [x] 3.1 Run the application with default values to confirm it starts correctly
- [x] 3.2 Test with custom timeout values (e.g., SERVER_READ_TIMEOUT=5s) to verify parsing works
- [x] 3.3 Test with invalid duration format to verify error handling
