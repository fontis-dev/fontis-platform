# Tasks

## Current phase: Platform Foundation

### Build system and toolchain

- [ ] Create Makefile with targets: build, fmt, lint, typecheck, test-unit, test-integration, security-scan, clean.
- [ ] Set up Yocto layer structure for core OS image.
- [ ] Configure Linux kernel with required drivers (storage, networking, TPM, UEFI) and security features (SELinux/AppArmor, integrity subsystem, dm-crypt, dm-verity).
- [ ] Set up Rust toolchain and project structure for `core/hal/`.
- [ ] Set up Go toolchain and project structure for `runtime/` services.
- [ ] Create CI pipeline (`.github/workflows/ci.yml`) with build, lint, test jobs.
- [ ] Configure Dependabot for Go modules, Rust crates, GitHub Actions.

### Core OS boot chain

- [ ] Implement UEFI Secure Boot configuration.
- [ ] Create signed bootloader (systemd-boot or GRUB) configuration.
- [ ] Create initramfs with minimal recovery shell.
- [ ] Implement measured boot via TPM (event log, PCR extension).
- [ ] Verify boot chain in QEMU.

### Contracts and scaffolding

- [ ] Define Protobuf schemas for identity service.
- [ ] Define Protobuf schemas for auth service.
- [ ] Set up Protobuf code generation in the build system.
- [ ] Implement identity service skeleton (gRPC server, empty handlers).
- [ ] Implement auth service skeleton (gRPC server, empty handlers).
- [ ] Implement storage service skeleton (gRPC server, empty handlers).
- [ ] Implement networking service skeleton (gRPC server, empty handlers).

### Testing

- [ ] Write integration test harness for runtime services.
- [ ] Write unit tests for Rust HAL crate.
- [ ] Write unit tests for Go service skeletons.
- [ ] Verify CI pipeline runs all tests.

### Documentation

- [x] Create AGENTS.md with platform-specific instructions.
- [x] Create SPEC.md with platform requirements.
- [x] Create ROADMAP.md with phased roadmap.
- [x] Create TASKS.md (this file).
- [ ] Create docs/ARCHITECTURE.md.
- [ ] Create docs/STANDARDS.md.
- [ ] Create CONTRIBUTING.md, SECURITY.md, CODE_OF_CONDUCT.md.

## Next phase: Identity and Auth

- [ ] Implement identity service CRUD operations (households, profiles).
- [ ] Implement auth service (password hashing, session management, API tokens).
- [ ] Implement mTLS for inter-service communication.
- [ ] Write integration tests for identity + auth flow.
- [ ] Security review of auth implementation.

## Completion rule

Move tasks to completed only after their acceptance criteria and required validation pass. Do not mark implementation tasks done without running the corresponding tests or checks. Mark documentation tasks done only after the document is written and reviewed.
