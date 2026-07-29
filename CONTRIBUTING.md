# Contributing to fontis-platform

## Code of Conduct

This project follows the Contributor Covenant Code of Conduct. See CODE_OF_CONDUCT.md.

## Development Setup

### Prerequisites

- Go 1.23+
- Rust 1.80+
- GCC or Clang (for C dependencies)
- QEMU (for emulation)
- GNU Make

### Quickstart

```bash
# Clone the repository
git clone https://github.com/fontis-dev/fontis-platform
cd fontis-platform

# Build everything
make build

# Run tests
make test-unit

# Boot in QEMU (when available)
make qemu
```

## How to Contribute

### Reporting Bugs

Open an issue with:
- A clear title and description.
- Steps to reproduce.
- Expected vs actual behavior.
- Environment details (hardware, platform version if applicable).

### Suggesting Features

Open an issue with:
- The problem you're solving.
- The proposed solution.
- Alternatives considered.
- How this aligns with the Fontis vision (local-first, private, simple).

### Pull Requests

1. Read AI_CONSTITUTION.md from fontis-foundation for org-wide AI rules.
2. Read AGENTS.md for platform-specific agent instructions.
3. Read SPEC.md for requirements and acceptance criteria.
4. Read ROADMAP.md for the current phase and exit criteria.
5. Read TASKS.md for current actionable work.
6. Make your changes in a feature branch.
7. Run the full local validation gate before committing.
8. Open a PR with the PR template. Keep it draft until all gates pass.
9. Address review feedback. Re-run validation after changes.
10. Request re-review. Merge only after approval.

### PR Guidelines

- Keep PRs small and focused. One logical change per PR.
- Every PR must include tests.
- Every PR must update documentation if behavior changes.
- Breaking changes require a MAJOR version bump and must be coordinated with consumers.

## Local Validation Gate

Before committing or opening a PR, run:

```bash
make fmt
make lint
make typecheck
make build
make test-unit
make test-integration
make security-scan
```

## Architecture Decisions

Significant architecture decisions are recorded as ADRs in the fontis-foundation repository. If your change affects cross-service contracts, system architecture, or security boundaries, write an ADR before implementing.
