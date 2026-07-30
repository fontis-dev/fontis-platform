# fontis-platform Engineering Standards

In addition to the org-wide standards in fontis-foundation `docs/STANDARDS.md`, the following platform-specific standards apply.

---

## Go Standards (Runtime Services)

### Project Layout

Each runtime service is in its own subdirectory under `runtime/<service>/` with this structure:

```
runtime/<service>/
├── cmd/
│   └── <service>/               Main binary entry point
│       └── main.go
├── internal/
│   ├── server/                  gRPC server implementation
│   ├── store/                   Data access layer (SQLite)
│   └── config/                  Configuration loading
├── api/
│   └── <service>.go             Public API types and interfaces
├── proto/                       Generated Protobuf code (committed)
│   └── <service>.pb.go
├── go.mod
├── go.sum
├── <service>_test.go            Unit tests
└── <service>_integration_test.go Integration tests
```

### Dependencies

- Standard library preferred. External dependencies require engineering review.
- Allowed without review: google.golang.org/grpc, google.golang.org/protobuf, modernc.org/sqlite.
- Any other dependency must be justified in a comment in go.mod or a README.
- Pin all dependencies to exact versions. No `latest` or range-version imports.

### Concurrency

- Use goroutines and channels for concurrency. Avoid sync primitives where channels suffice.
- Every goroutine must have a clear lifecycle: start condition, stop condition, error handling.
- Context.Context is the first parameter of every function that may block or make RPCs.
- Use errgroup for managing groups of goroutines.

### Error Handling

- Errors are values. Wrap them with context using `fmt.Errorf("doing X: %w", err)`.
- Do not use `log.Fatal` in libraries. Panic only in initialization code that cannot proceed.
- Every gRPC handler must return appropriate status codes (codes.NotFound, codes.InvalidArgument, etc.).
- Log errors at the appropriate level. Errors that require human attention are `log.Error`. Recoverable errors are `log.Warn`.

---

## Rust Standards (HAL)

### Project Layout

```
core/hal/
├── src/
│   ├── lib.rs                  Public API re-exports
│   ├── block/                  Block device HAL
│   ├── tpm/                    TPM 2.0 HAL
│   ├── power/                  Power management HAL
│   └── ffi/                    C FFI exports for Go
├── tests/                      Integration tests
├── Cargo.toml
└── build.rs
```

### Conventions

- Follow the Rust API Guidelines. All public items have doc comments.
- Use `#![deny(unsafe_code)]` in lib.rs. Unsafe code is forbidden in HAL except in the `ffi/` module.
- The `ffi/` module must be minimal: thin C-compatible wrappers around safe Rust functions.
- All HAL functions return `Result<T, HalError>`. Errors are opaque enums with Display and Debug.

### Safety

- The FFI layer validates all inputs from Go (pointer validity, buffer sizes, null checks) before calling safe Rust code.
- No shared mutable state across the FFI boundary. All state is owned by the Rust side and accessed through handles.
- Use `#[repr(C)]` for all types crossing the FFI boundary.

---

## Protobuf Contract Standards

### File Layout

```
contracts/protobuf/
├── identity/
│   ├── v1/
│   │   ├── identity.proto
│   │   └── identity_grpc.proto  (or use grpc service definitions)
│   └── v1alpha1/                 Pre-release schemas during development
├── auth/
│   └── v1/
│       └── auth.proto
└── google/                       Standard Google well-known types (imported)
```

### Conventions

- Package name: `fontis.<service>.v1`.
- File name: snake_case matching the service name.
- Every message and field has a comment explaining its purpose.
- Field numbers start at 1. Reserve field numbers for deleted fields using `reserved`.
- Enums start at 0 with an `UNSPECIFIED` value.
- Breaking changes (field removal, type change, number reassignment) require a new package version.
- Non-breaking changes (new field, new message, new enum value) are allowed within a version.

### Code Generation

- Generated Go code is committed to the service's `runtime/<service>/proto/` directory.
- CI verifies that generated code matches the Protobuf source (`make protoc-gen` and `git diff --exit-code`).
- Never manually edit generated files.

---

## Testing Standards

### Go Tests

- Unit tests use the standard `testing` package. Use `testify/assert` and `testify/require` for assertions.
- Integration tests are in `_integration_test.go` files and are excluded from `go test ./...` with `//go:build integration`.
- Integration tests can be run with `make test-integration` or `go test -tags=integration ./...`.
- Use `testing/sqlite` (or `modernc.org/sqlite`) for database-dependent tests. Never depend on an external database in unit tests.
- Mock gRPC dependencies using the generated mock interfaces or testify/mock.

### Rust Tests

- Unit tests are in `#[cfg(test)] mod tests {}` blocks within each source file.
- Integration tests are in `tests/` and use the public HAL API only.
- Tests that require hardware (TPM, block devices) are marked with `#[cfg(feature = "hwtest")]` and excluded from `cargo test`.

### Coverage

- Target: 80%+ coverage for Go services, 90%+ for Rust HAL.
- Coverage is tracked in CI. PRs must not decrease coverage without justification.
- Critical paths (auth, crypto, storage operations) require 100% coverage of error paths.

---

## Build System Standards

### Makefile

The root `Makefile` is the single entry point for all build operations. Standard targets:

| Target | Description |
| --- | --- |
| `build` | Build runtime services and HAL |
| `build-runtime` | Build all runtime Go services |
| `build-hal` | Build Rust HAL crate |
| `build-image` | Build Debian-based core OS image (planned, not yet implemented) |
| `fmt` | Format all code |
| `lint` | Lint all code |
| `typecheck` | Type-check all code |
| `test-unit` | Run all unit tests |
| `test-integration` | Run all integration tests |
| `security-scan` | Scan dependencies for vulnerabilities |
| `clean` | Clean all build artifacts |
| `qemu` | Boot OS image in QEMU |
| `protoc-gen` | Regenerate Protobuf Go code |

### CI

The CI pipeline runs on every PR push and on merge to main:

1. `git diff --check` (whitespace errors)
2. `make fmt` (format check, with `git diff --exit-code`)
3. `make lint`
4. `make typecheck`
5. `make build`
6. `make test-unit`
7. `make test-integration`
8. `make security-scan`
9. CodeQL analysis
10. Dependency review (if go.mod/go.sum or Cargo.toml/Cargo.lock changed)

All CI checks are required to pass before merge.
