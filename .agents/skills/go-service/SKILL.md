---
name: go-service
description: Build, test, and refactor Go runtime services following platform conventions (gRPC, SQLite, config, mTLS).
---

## go-service

### Conventions

- Service layout: `cmd/<service>/main.go`, `internal/server/`, `internal/store/`, `internal/config/`, `proto/`
- Every service is an independent Go process. Inter-service communication is gRPC over local Unix sockets with mTLS.
- SQLite via `modernc.org/sqlite` for data storage. `database/sql` with `SetMaxOpenConns(1)`.
- Config loaded from environment variables. Defaults suitable for dev.
- Error wrapping: `fmt.Errorf("doing X: %w", err)`. Context.Context is the first parameter.
- Use `errgroup` for goroutine lifecycle management.

### Testing

- Unit tests next to source. Use `testing` stdlib. `testify/assert` and `testify/require` for assertions.
- Integration tests in `_integration_test.go` with `//go:build integration` build tag.
- SQLite-dependent tests use `modernc.org/sqlite` — no external database needed.

### Adding a new service

1. Create `runtime/<service>/cmd/<service>/main.go` with signal handling, config load, server start/shutdown.
2. Create `runtime/<service>/internal/config/config.go` with env-var-based config.
3. Create `runtime/<service>/internal/server/server.go` with gRPC server, mTLS via `runtime/pkg/tls`.
4. Create `runtime/<service>/internal/store/store.go` with SQLite migrations and CRUD.
5. Define Protobuf schemas in `contracts/protobuf/<service>/v1/`, run `make protoc-gen`.
6. Add `runtime/<service>/go.mod` and wire it into the root `Makefile`.
