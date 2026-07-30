# fontis-platform Specification

## Problem

Households lack a trusted, private, local-first digital home platform. Their digital lives are fragmented across dozens of services and devices. They cannot confidently answer "where is everything that matters?"

fontis-platform is the device substrate that solves this: a secure, maintainable, simple OS and runtime that provides the foundational services every household needs, while remaining open for extension through modules.

## Users

- **Household members** — non-technical users who interact through the device's native display output, web UI, or mobile app. They need reliability, privacy, and simplicity.
- **Module developers** — technical users who build and publish modules extending platform functionality. They need stable APIs, clear contracts, and an SDK.
- **Device owners** — technology-interested users who may install Community Edition on their own hardware. They need documentation and a working reference implementation.

## Required behavior

### Core OS

- Boot securely with measured boot and verified root filesystem.
- Provide a stable hardware abstraction layer that runtime services consume.
- Support unattended updates with atomic upgrade and rollback.
- Enforce process isolation: each runtime service runs as a separate Linux user with minimal capabilities.
- Provide a native display stack with compositor, GPU acceleration, and multi-resolution output (TV to monitor, 720p–4K).
- Support input from remote control, game controller, keyboard, and mouse.

### Platform Runtime

- **Identity service** — manage user profiles, household membership, and identity verification.
- **Auth service** — authenticate users, authorize actions, manage sessions and API tokens.
- **Networking service** — manage network interfaces, WiFi, ethernet, VPN. Provide local device discovery.
- **Storage service** — manage local storage pools, filesystem layout, volume encryption. Expose storage to modules through a documented API.
- **Logging service** — collect, store, and query system and module logs. Support structured logging and log rotation.
- **Module lifecycle service** — install, update, rollback, and remove modules. Isolate modules in containers with resource limits. Modules that provide a native UI render through the platform display service.
- **Update service** — check for, download, verify, and apply system updates. Support automatic rollback on failure.
- **Backup service** — schedule and manage backups to local or network destinations. Cloud backup is opt-in and non-destructive.
- **Marketplace client** — discover, browse, and download modules from the module marketplace (optional, graceful degradation).

### Security

- Full disk encryption (LUKS + TPM).
- Mutual TLS for all inter-service communication.
- Secure boot chain with signed bootloader, kernel, and initramfs.
- Module sandboxing with no direct hardware access.
- Secrets never stored in source code, logs, or unencrypted configuration.

## Architecture

See `docs/ARCHITECTURE.md` for the complete architecture.

## Technology stack

| Layer | Technology | Rationale |
| --- | --- | --- |
| OS substrate | Linux (Debian Stable) | Proven stability, long-term security support, widely supported appliance OS |
| Hardware abstraction | Rust | Memory safety, zero-cost abstractions, security-critical |
| Runtime services | Go | Simple, explicit, excellent tooling, strong standard library |
| Storage metadata | SQLite | Embedded, reliable, zero-administration |
| Bulk storage | ZFS or Btrfs | Snapshots, checksums, send/receive |
| Inter-service | gRPC + Protobuf | Typed contracts, code generation, efficient |
| Display server | Wayland compositor | Modern, secure, GPU-accelerated display protocol |
| UI framework | TBD (Qt / GTK / Flutter / custom) | Native rendering for TV/monitor output |
| GPU abstraction | Mesa + DRM/KMS | Open-source GPU driver stack |
| Input handling | libinput + evdev | Unified input for remote, gamepad, keyboard, mouse |
| Module isolation | OCI containers | Industry standard, strong isolation |
| Lightweight modules | WebAssembly | Sandboxed, portable, safe |
| Build system | GNU Make | Simple, explicit, universal |

## Compatibility

- **Official Fontis Device** — fully supported hardware with guaranteed compatibility.
- **Community hardware** — x86-64 systems with UEFI, TPM 2.0, and compatible storage. Community Edition only.
- **Supported architectures** — x86-64 (primary), ARM64 (future).
- **Boot modes** — UEFI only (no legacy BIOS support).
- **Display output** — HDMI/DisplayPort. Native resolutions 720p, 1080p, 1440p, 4K. Variable refresh rate preferred. Multi-monitor optional (future).
- **Audio output** — HDMI audio, analog 3.5mm, Bluetooth (optional, future).
- **Input** — IR remote control (included), Bluetooth game controller, USB keyboard and mouse.

## Non-goals

- Being a general-purpose operating system. Fontis runs only the platform services and modules.
- Being a cloud service. Fontis is local-first. Cloud is always optional and opt-in.
- Being a developer workstation. There is no shell, compiler toolchain, or development environment on the device.
- Supporting legacy BIOS boot.
- Providing a general-purpose graphical desktop environment (file manager, arbitrary app launcher, window manager). The native UI is purpose-built for platform functionality and modules.
- Replacing a NAS, media server, or home automation hub. Those are modules.

## Acceptance criteria

- A fresh platform build completes from a clean checkout with `make build`.
- The built image boots on reference hardware and presents a setup wizard through the native display UI and/or web UI.
- A user can complete initial setup (create household, configure networking, set up storage).
- A first-party module can be installed, configured, and removed through the native UI and/or web UI.
- The update service can download, verify, and apply a signed update with automatic rollback on failure.
- All runtime services pass their unit and integration test suites.
- Security scan reports zero critical or high vulnerabilities in production dependencies.
- The platform functions fully without network access or cloud enrollment. Cloud features gracefully degrade.
- A full backup can be created and restored without data loss.

## Unresolved questions

- Whether to use ZFS or Btrfs as the primary filesystem. Decision deferred to Phase 1 implementation based on stability, performance, and snapshot support on reference hardware.
- Whether the first release targets the Official Device or Community Edition first. Likely Community Edition for developer iteration.
- The exact kernel version and configuration baseline. Will be pinned during Phase 1 build system setup.
- Native UI framework choice: Qt, GTK, Flutter, or custom lightweight toolkit. Decision deferred to Phase 2 based on GPU requirements, module rendering needs, and development velocity.
- Wayland compositor selection: wlroots-based (simple, flexible) or Weston (reference implementation). Choice affects GPU driver compatibility and module rendering surface support.
- Whether native UI modules render as Wayland surfaces (separate windows) or through a shared rendering API (single composited surface) -- affects module isolation and security model.
