# fontis-platform

The device substrate repository for the [Fontis](https://github.com/fontis-dev) ecosystem.

This repository contains the Core OS layer and Platform Runtime that execute on the household Fontis Device. It is the largest and most security-sensitive repository in the Fontis ecosystem.

## Repository structure

```
core/               Core OS layer
├── kernel/         Linux kernel configuration
├── boot/           Bootloader, initramfs, secure boot
├── packages/       Base system packages and OS image
└── hal/            Hardware abstraction layer (Rust)

runtime/            Platform Runtime services (Go)
├── identity/       User identity, profiles, households
├── auth/           Authentication and authorization
├── networking/     Network management (WiFi, ethernet, VPN)
├── storage/        Local storage and filesystem management
├── logging/        System logging and telemetry
├── module-lifecycle/   Module install, update, remove
├── updates/        System updates and rollback
├── backup/         Backup and restore
└── marketplace-client/  Module marketplace discovery

contracts/          Shared interface definitions
├── protobuf/       gRPC service schemas
├── api/            REST API specifications
└── events/         Event type definitions

tools/              Platform build tooling
tests/              Cross-cutting tests
docs/               Architecture and standards documentation
```

## Quickstart

```bash
# Build all targets
make build

# Run unit tests
make test-unit

# Boot in QEMU (when available)
make qemu
```

## Documentation

- `AGENTS.md` — AI agent instructions for this repository
- `SPEC.md` — Product requirements and acceptance criteria
- `ROADMAP.md` — Development phases and exit criteria
- `TASKS.md` — Current work items
- `docs/ARCHITECTURE.md` — Platform internal architecture
- `docs/STANDARDS.md` — Platform-specific engineering standards
- `fontis-foundation` — Organization-wide governance and AI constitution

## License

MIT
