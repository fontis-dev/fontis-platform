# Repository instructions

## Scope

This repository is the fontis-dev device substrate: the Core OS layer and Platform Runtime that execute on the household Fontis Device. It contains kernel configuration, boot infrastructure, hardware abstraction, and all runtime services (identity, auth, networking, storage, logging, module lifecycle, updates, backup, marketplace client). It does not contain modules, cloud services, or UI applications.

For org-wide AI rules, read AI_CONSTITUTION.md from fontis-foundation. For engineering philosophy and decision precedence, read ENGINEERING_PRINCIPLES.md from fontis-foundation.

## Internal structure

```
├── core/                       Core OS layer
│   ├── boot/                   Secure boot and key management
│   ├── hal/                    Hardware abstraction layer (Rust)
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

## Workflow (follow in order)

1. **Inspect** — branch, worktree, validation state, TASKS.md status.
2. **Plan** — map the change to a SPEC.md requirement, break into reviewable items.
3. **Implement one phase** — smallest reviewable increment. Touch only what the task requires.
4. **Run focused checks** — `go vet ./...` and `go build ./...` and `go test -count=1 ./...` on affected packages.
5. **Run the full local gate** — `git diff --check && make fmt && make lint && make typecheck && make build && make test-unit && make test-integration && make security-scan`. Fix any failures.
6. **Update TASKS.md** — only after the gate passes.
7. **Commit** — one focused commit per item. No drive-by changes.
8. **Open PR** — fill in `.github/PULL_REQUEST_TEMPLATE.md` completely. CodeRabbit auto-reviews the PR (unlimited Pro+ on public repos). Fix its findings with follow-up commits.
9. **Independent human review** — not self-review. Resolve all threads.
10. **Manual testing** — on real target if HAL changes, screenshots if UI changes. QEMU or real-target validation required for boot and security changes (kernel, initramfs, bootloader/TPM, SELinux integration).
11. **Merge** — only when CI, CodeRabbit, independent review, security scans, and manual tests are clean.

## Command execution

- Prefer running code, tests, linters, and type checks over guessing.
- Read complete errors, logs, and stack traces before fixing them.
- Install RTK (`cargo install --git https://github.com/rtk-ai/rtk`) for command output compression. Use `rtk` when output is large or repetitive and a filtered summary is sufficient. Use raw commands when exact output matters.
- RTK reduces token consumption by 60-99% on build, test, and lint output. Prefer it for `make lint`, `make build`, `make test-unit`, and `make test-integration` runs.

## Build system

```bash
make build          # Build runtime services and HAL
make build-runtime  # Build all runtime Go services
make build-hal      # Build Rust HAL crate
make build-image    # Build Debian-based core OS image (planned, not yet implemented)
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
- Run `make build-image` after changing image configuration or base OS packages (once the target is implemented).
- For HAL changes affecting real hardware, manual testing on target is required.

## Pull requests

- Use `.github/PULL_REQUEST_TEMPLATE.md` for every PR. Fill in all sections.
- CodeRabbit is installed on this repo and auto-reviews every PR (unlimited Pro+ on public repos). Fix all actionable findings; document intentionally skipped items with a reason.
- After opening a PR, require CI validation, security checks, CodeRabbit review, and at least one independent human review on the latest commit before merging.
- Every review thread must be resolved or have a documented response.
- Record manual testing evidence in the PR body. For HAL changes, test on real target hardware. For UI changes, include screenshots.

## Skills

Reusable workflows are in `.agents/skills/`. Each skill has a `SKILL.md` with YAML front matter (`name`, `description`) for auto-discovery:

- `$go-service` — build, test, refactor Go runtime services (gRPC, SQLite, mTLS)
- `$rust-hal` — develop Rust HAL crate (safety, FFI, error conventions)
- `$protobuf-contracts` — define, review, regenerate Protobuf schemas and stubs
- `$pr-readiness` — validate PR merge readiness (CodeRabbit, CI, independent review, manual testing)
- `$ai-project-manager` — produce requirement-linked plans, pause for approval, execute one phase at a time

Invoke a skill explicitly when needed, or let the agent auto-select based on task description.

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
