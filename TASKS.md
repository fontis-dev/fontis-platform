# Tasks

## Current phase: Platform Foundation

### Build system and toolchain

- [x] Create Makefile with targets: build, fmt, lint, typecheck, test-unit, test-integration, security-scan, clean.
- [ ] Set up Yocto layer structure for core OS image (requires Linux build host).
- [x] Configure Linux kernel with required drivers (storage, networking, TPM, UEFI) and security features.
- [x] Set up Rust toolchain and project structure for `core/hal/`.
- [x] Set up Go toolchain and project structure for `runtime/` services.
- [x] Create CI pipeline (`.github/workflows/ci.yml`) with build, lint, test jobs.
- [x] Configure Dependabot for Go modules, Rust crates, GitHub Actions.

### Core OS boot chain

- [x] Create UEFI Secure Boot configuration directory and setup docs.
- [ ] Create signed bootloader (systemd-boot or GRUB) configuration (requires Linux build host).
- [x] Create initramfs configuration directory and setup docs.
- [ ] Implement measured boot via TPM (event log, PCR extension) (requires Linux build host).
- [ ] Verify boot chain in QEMU (requires Linux build host).

### Contracts and scaffolding

- [x] Define Protobuf schemas for identity service.
- [x] Define Protobuf schemas for auth service.
- [x] Define Protobuf schemas for storage service.
- [x] Define Protobuf schemas for networking service.
- [x] Set up Protobuf code generation in the build system.
- [x] Implement identity service skeleton (gRPC server, config, store).
- [x] Implement auth service skeleton (gRPC server, config, store).
- [x] Implement storage service skeleton (gRPC server, config, store).
- [x] Implement networking service skeleton (gRPC server, config, store).
- [x] Implement logging service skeleton.
- [x] Implement module-lifecycle service skeleton.
- [x] Implement updates service skeleton.
- [x] Implement backup service skeleton.
- [x] Implement marketplace-client service skeleton.

### Testing

- [ ] Write integration test harness for runtime services.
- [ ] Write unit tests for Rust HAL crate.
- [ ] Write unit tests for Go service skeletons.
- [ ] Verify CI pipeline runs all tests (requires Linux/macOS build host).

### Documentation

- [x] Create AGENTS.md with platform-specific instructions.
- [x] Create SPEC.md with platform requirements.
- [x] Create ROADMAP.md with phased roadmap.
- [x] Create TASKS.md (this file).
- [x] Create docs/ARCHITECTURE.md.
- [x] Create docs/STANDARDS.md.
- [x] Create CONTRIBUTING.md, SECURITY.md, CODE_OF_CONDUCT.md.

## Next phase: Identity and Auth

- [ ] Implement identity service CRUD operations (households, profiles).
- [ ] Implement auth service (password hashing, session management, API tokens).
- [ ] Implement mTLS for inter-service communication.
- [ ] Write integration tests for identity + auth flow.
- [ ] Security review of auth implementation.

## Completion rule

Move tasks to completed only after their acceptance criteria and required validation pass. Do not mark implementation tasks done without running the corresponding tests or checks. Mark documentation tasks done only after the document is written and reviewed.
