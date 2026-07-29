# Repository instructions

## Scope

This repository is the fontis-dev device substrate: the Core OS layer and Platform Runtime that execute on the household Fontis Device. It contains kernel configuration, boot infrastructure, hardware abstraction, and all runtime services (identity, auth, networking, storage, logging, module lifecycle, updates, backup, marketplace client). It does not contain modules, cloud services, or UI applications.

For org-wide AI rules, read AI_CONSTITUTION.md from fontis-foundation. For engineering philosophy and decision precedence, read ENGINEERING_PRINCIPLES.md from fontis-foundation.

## Internal structure

```
├── core/                       Core OS layer
│   ├── kernel/                 Linux kernel configuration
│   ├── boot/                   Bootloader, initramfs, secure boot
│   ├── packages/               Base system packages and image
│   └── hal/                    Hardware abstraction layer (Rust)
├── runtime/                    Platform Runtime services (Go)
│   ├── identity/               User identity, profiles, households
│   ├── auth/                   Authentication and authorization
│   ├── networking/             Network management
│   ├── storage/                Local storage and filesystem
│   ├── logging/                System logging and telemetry
│   ├── module-lifecycle/       Module install, update, remove
│   ├── updates/                System updates and rollback
│   ├── backup/                 Backup and restore
│   └── marketplace-client/     Module marketplace discovery
├── contracts/                  Shared interface definitions
│   ├── protobuf/               gRPC service schemas
│   ├── api/                    REST API specifications
│   └── events/                 Event type definitions
├── tools/                      Platform build tooling
├── tests/                      Cross-cutting tests
├── docs/                       Architecture and standards
├── scripts/                    Validation and build scripts
└── .github/                    CI/CD configuration
```

## Operating principles

- Working code only. Plausibility is not correctness; verify before reporting done.
- Never fabricate file paths, APIs, commit hashes, command output, or test results.
- Touch only what the task requires. Avoid drive-by refactors, formatting, or cleanup.
- Keep communication direct and concise. Skip flattery, filler, ceremonial openings, and emoji.
- Read the ENGINEERING_PRINCIPLES.md precedence: Security > Correctness > Maintainability > Simplicity > Reliability > Performance > Developer Productivity > Feature Velocity.

## Command execution

- Prefer running code, tests, linters, and type checks over guessing.
- Read complete errors, logs, and stack traces before fixing them.
- Use `rtk` when command output is large or repetitive and a filtered summary is sufficient. Use raw commands when exact output matters.

## Build system

```bash
make build          # Build all targets
make fmt            # Format all code
make lint           # Lint all code
make typecheck      # Type-check (Go, Rust)
make test-unit      # Run unit tests
make test-integration  # Run integration tests
make security-scan  # Scan dependencies for vulnerabilities
make clean          # Clean build artifacts
```

## Architecture rules

- `runtime/` services must never import from `core/` directly. All hardware access goes through `core/hal/` interfaces.
- Each runtime service is an independent Go process. Inter-service communication is through gRPC over local Unix sockets.
- Protobuf schemas in `contracts/protobuf/` are the source of truth for service interfaces. Generated code must never be manually edited.
- Every runtime service must have its own test suite. Cross-service integration tests go in `tests/integration/`.
- No service has access to another service's data store. All access is through the service's public API.
- Security-critical code (crypto, auth, HAL) must be written in Rust. Service-level code should be Go.

## Validation

- Run the full local gate before committing: `git diff --check && make fmt && make lint && make typecheck && make build && make test-unit && make test-integration && make security-scan`
- Run `scripts/validate.sh` after changing project configuration, build scripts, or CI.
- For HAL changes affecting real hardware, manual testing on target is required.

## Documentation routing

- This file for platform-specific agent rules.
- `SPEC.md` for platform requirements and acceptance criteria.
- `ROADMAP.md` for phased outcomes and exit criteria.
- `TASKS.md` for current work and validation status.
- `docs/ARCHITECTURE.md` for platform internal architecture.
- `docs/STANDARDS.md` for platform coding standards.
- fontis-foundation `AI_CONSTITUTION.md` for org-wide AI rules.
- fontis-foundation `ENGINEERING_PRINCIPLES.md` for decision precedence.

## Maintenance

- Keep this file concise. Add rules only when they prevent a real repeat mistake.
- When the AI makes a repeat error, tighten the relevant rule rather than appending warnings.
- This file is reviewed with every architecture-significant PR.
