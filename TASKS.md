# Tasks

## Current phase: Platform Foundation

### Build system and toolchain

- [x] Create Makefile with targets: build, fmt, lint, typecheck, test-unit, test-integration, security-scan, clean.
- [ ] Set up Yocto layer structure for core OS image (structure created; full build validation requires a Yocto build environment, see `make build-core`).
- [ ] Configure Linux kernel with required drivers (storage, networking, TPM, UEFI) and security features (SELinux/AppArmor, integrity subsystem, dm-crypt, dm-verity).
- [ ] Add GPU drivers to kernel config (i915, amdgpu, dummy) for DRM/KMS support.
- [ ] Add evdev and ALSA drivers to kernel config for input and audio.
- [x] Set up Rust toolchain and project structure for `core/hal/`.
- [x] Set up Go toolchain and project structure for `runtime/` services.
- [x] Create CI pipeline (`.github/workflows/ci.yml`) with build, lint, test jobs.
- [x] Create security pipeline (`.github/workflows/security.yml`) with gitleaks and trivy scans.
- [x] Configure Dependabot for Go modules, Rust crates, GitHub Actions.

### Core OS boot chain

- [ ] Implement UEFI Secure Boot configuration.
- [ ] Create signed bootloader (systemd-boot or GRUB) configuration.
- [ ] Create initramfs with minimal recovery shell.
- [ ] Implement measured boot via TPM (event log, PCR extension).
- [ ] Verify boot chain in QEMU.

### Display and input stack

- [ ] Evaluate and select Wayland compositor (wlroots-based vs Weston).
- [ ] Implement GPU/display HAL module (DRM connector enumeration, mode setting, framebuffer allocation).
- [ ] Implement audio HAL module (ALSA device enumeration, volume control, sink selection).
- [ ] Implement IR remote HAL module (gpio-ir or mceusb driver integration).
- [ ] Integrate compositor into core OS image.
- [ ] Implement display service (Go): compositor lifecycle, UI state coordinator, input dispatch.
- [ ] Implement display service gRPC Protobuf contracts (surface register, input routing, display info).
- [ ] Implement system UI compositor client (setup wizard, volume/input overlays).
- [ ] Implement HDMI hotplug detection and resolution auto-detection.
- [ ] Implement input capture pipeline (compositor → libinput → display service → module/system UI).
- [ ] Support game controller (Bluetooth) pairing and input.
- [ ] Support audio output switching (HDMI audio, analog 3.5mm).

### Contracts and scaffolding

- [x] Define Protobuf schemas for identity service.
- [x] Define Protobuf schemas for auth service.
- [x] Define Protobuf schemas for storage service.
- [x] Define Protobuf schemas for networking service.
- [x] Set up Protobuf code generation in the build system.
- [x] Implement identity service skeleton (gRPC server, empty handlers).
- [x] Implement auth service skeleton (gRPC server, empty handlers).
- [x] Implement storage service skeleton (gRPC server, empty handlers).
- [x] Implement networking service skeleton (gRPC server, empty handlers).

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
- [x] Create docs/ARCHITECTURE.md.
- [x] Create docs/STANDARDS.md.
- [x] Create CONTRIBUTING.md, SECURITY.md, CODE_OF_CONDUCT.md.
- [x] Create `.github/PULL_REQUEST_TEMPLATE.md` with mandatory checklists.
- [x] Create reusable skills in `.agents/skills/` (go-service, rust-hal, protobuf-contracts, pr-readiness, ai-project-manager).
- [x] Create execution policy rules (`.opencode/rules.jsonc`).
- [x] Add CodeRabbit local and PR-based review workflow to AGENTS.md.
- [x] Add RTK installation and usage guidance to AGENTS.md.

## Next phase: Identity and Auth

- [x] Implement identity service CRUD operations (households, profiles).
- [x] Implement auth service (password hashing, session management, API tokens).
- [x] Implement mTLS for inter-service communication.
- [ ] Write integration tests for identity + auth flow.
- [ ] Security review of auth implementation.

## Phase 3: Display and Audio

- [ ] Integrate compositor as system compositor (boot-time start, display output active at login prompt).
- [ ] Implement full display service gRPC API.
- [ ] Implement system UI compositor client (home screen, navigation grid).
- [ ] Implement input dispatch from compositor to display service to focused consumer.
- [ ] Support all four input types: IR remote, gamepad, keyboard, mouse.
- [ ] Implement audio service integration (ALSA, volume, HDMI/analog switching).
- [ ] Implement HDMI hotplug and resolution switching.
- [ ] Write integration tests for display service + compositor interaction.

## Phase 4: Module UI

- [ ] Design module UI rendering approach (Wayland surface per module vs shared rendering API).
- [ ] Implement module Wayland surface integration (compositor accepts surfaces from module containers).
- [ ] Implement display service module surface management (register, focus, fullscreen, overlay).
- [ ] Build module UI SDK helpers (Wayland protocol bindings, input event subscription).
- [ ] Build reference module with native UI (media library with remote/gamepad navigation).
- [ ] Write integration tests for module UI lifecycle.

## Completion rule

Move tasks to completed only after their acceptance criteria and required validation pass. Do not mark implementation tasks done without running the corresponding tests or checks. Mark documentation tasks done only after the document is written and reviewed.
