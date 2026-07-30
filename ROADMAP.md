# fontis-platform Roadmap

## Phase 1: Platform Foundation

Status: Planned

### Outcome

Establish the build system, core OS image, and development toolchain. Produce a bootable image for Community Edition (x86-64, UEFI) with secure boot and a minimal initramfs.

### Included work

- Set up Debian Stable base with custom image builder (debootstrap + hardening + Fontis packages).
- Select Linux kernel (Debian Stable kernel or backports) with required drivers (storage, networking, TPM, UEFI, GPU, input, audio).
- Evaluate Wayland compositor options (wlroots-based, Weston) and select target.
- Implement boot chain hardening: UEFI Secure Boot key generation, initramfs TPM integration.
- Set up the Rust HAL toolchain and project structure under `core/hal/` including GPU and audio HAL stubs.
- Set up the Go runtime project structure with build tooling.
- Create the `contracts/protobuf/` directory with initial identity and auth service schemas.
- Implement CI pipeline: build, lint, format check, unit test.
- Write integration test harness for runtime services.
- Document architecture (docs/ARCHITECTURE.md) including display/input stack and standards (docs/STANDARDS.md).

### Risks

- UEFI Secure Boot setup requires signing infrastructure and key management planning.
- Debian package selection must be minimal; avoid unnecessary services and attack surface.

### Exit criteria

- `make build-image` produces a bootable x86-64 UEFI image from a clean checkout.
- The image boots on reference hardware (or QEMU) to a working Debian base with Fontis runtime services.
- The Rust HAL crate compiles with no warnings and includes GPU/audio stubs.
- The Go runtime skeleton compiles with no errors.
- Wayland compositor candidate is selected and evaluated on target hardware (renders a test surface through DRM/KMS).
- CI pipeline runs and passes on every PR push.
- Architecture and standards documents are reviewed and merged.
- All AI-generated code has been reviewed by a human.

---

## Phase 2: Identity and Auth

Status: Planned

### Outcome

First two runtime services operational: identity management and authentication. A user can create a household, add profiles, authenticate, and obtain session tokens.

### Included work

- Implement identity service: user profiles, household CRUD, profile management.
- Implement auth service: password-based auth, session management, API token generation.
- Finalize Protobuf contracts for identity and auth.
- Implement inter-service gRPC communication over local Unix sockets.
- Implement TLS for inter-service communication (self-signed CA for development).
- Write unit and integration tests for both services.
- Storage service stub: provide local SQLite database for service metadata.
- Implement display service skeleton (gRPC server, compositor lifecycle management stubs).
- Implement HAL audio stub and GPU stub.
- Integrate IR remote receiver driver and test basic input capture.

### Risks

- Getting gRPC and mTLS right on Unix sockets requires careful implementation. Plan to start with plain gRPC and add TLS after basic connectivity works.
- Password hashing and session management must follow security best practices (Argon2id, secure session cookies/tokens).
- IR remote driver support varies by hardware. Plan to support standard consumer IR (RC-6) via mceusb or gpio-ir initially.

### Exit criteria

- Identity service manages households and profiles through gRPC API.
- Auth service authenticates users and issues session tokens.
- Integration tests verify the full identity creation and authentication flow.
- Inter-service communication uses mTLS.
- No plaintext secrets in configuration or logs.
- Display service compiles and connects to the compositor (compositor runs as separate process, display service manages its lifecycle).
- IR remote input is captured and routed through the display service.

---

## Phase 3: Storage and Networking

Status: Planned

### Outcome

Storage management (pool setup, volume creation, encryption) and networking (WiFi, ethernet, local discovery) are operational.

### Included work

- Implement storage service: detect storage devices, create/mount encrypted volumes, manage filesystem pools.
- Implement networking service: scan WiFi networks, connect to WiFi/ethernet, manage VPN connections.
- Implement local device discovery (mDNS/DNS-SD).
- Expose storage and networking APIs through Protobuf contracts.
- Integrate with identity/auth for authenticated access.
- Write comprehensive tests.
- Integrate Wayland compositor (wlroots-based): GPU buffer allocation, DRM/KMS mode setting, HDMI audio output.
- Implement full display service: UI state coordination, input dispatch to modules/system UI, system overlay surfaces.
- Implement system UI compositor client: setup wizard screen, volume/HUD overlay.
- Support game controller (Bluetooth) input capture via compositor.
- Support HDMI hotplug detection and resolution switching (720p/1080p/1440p/4K auto-detection).
- Implement audio service integration: ALSA device enumeration, volume control, HDMI/analog output switching.

### Risks

- Storage encryption key management with TPM is complex. Start with passphrase-based LUKS and add TPM sealing later.
- WiFi management requires careful integration with NetworkManager or similar. Evaluate whether to use NetworkManager directly or implement a minimal WiFi daemon.
- Wayland compositor GPU requirements may limit hardware support. Mesa drivers for Intel/AMD GPUs are mature; ensure kernel DRM driver coverage for target hardware.
- HDMI audio (via GPU) requires correct ALSA device routing. Implementation complexity depends on GPU vendor.

### Exit criteria

- Storage service can initialize a disk, create an encrypted volume, and mount it.
- Networking service can connect to a WiFi network and obtain an IP address via DHCP.
- Device is discoverable by name on the local network.
- All operations require valid authentication.
- Display service renders a native setup wizard on HDMI output at up to 4K resolution.
- Compositor handles mode setting, GPU-accelerated rendering, and HDMI audio output.
- Input from IR remote, keyboard, and mouse is functional through the compositor.
- Game controller can be paired via Bluetooth and used for input.

---

## Phase 4: Module System

Status: Planned

### Outcome

First-party module can be installed, started, stopped, and removed through the module lifecycle service. Module isolation is enforced through containers.

### Included work

- Implement module lifecycle service: module manifest parsing, container creation, resource limits.
- Create module SDK (in `fontis-modules` repo): types, contract helpers, testing utilities.
- Build one first-party reference module (e.g., a media library or game launcher).
- Expose module management through Protobuf contracts.
- Integrate module lifecycle with the web UI (in `fontis-web` repo).
- Implement module update and rollback.
- Implement Wayland surface support for module containers: modules expose a Wayland surface through the compositor that the display service manages (focus, fullscreen, overlay).
- Build a module UI SDK: Wayland protocol surface helpers, input event subscription, display service API bindings.
- Implement module-to-compositor surface isolation: each module renders in its own Wayland buffer, compositor enforces surface boundaries.
- Build reference module with native UI (e.g., media library with thumbnails, navigation via remote/gamepad).

### Risks

- Container runtime (Docker/containerd/podman) selection affects module isolation guarantees. Evaluate for memory/disk overhead and security properties.
- Module manifest format needs to balance flexibility with simplicity. Start minimal and extend based on real module needs.
- Module Wayland surface integration adds complexity: the compositor needs to discover and manage surfaces from module containers. Evaluate approaches: (a) module runs a Wayland client inside the container and connects to compositor via socket mount, (b) display service creates surfaces on behalf of modules.

### Exit criteria

- A first-party module can be installed through the module lifecycle service.
- The module starts in an isolated container with resource limits.
- The module is accessible through the web UI and through the native display UI.
- A reference module renders a native UI surface managed by the compositor.
- Input events (remote, gamepad, keyboard, mouse) route correctly to the focused module surface.
- The module can be stopped and removed without leaving residual state.
- Module update with automatic rollback on failure works.

---

## Phase 5: Updates and Backup

Status: Planned

### Outcome

System update mechanism works with atomic upgrades and rollback. Backup and restore of household data is functional.

### Included work

- Implement update service: check for updates, download signed update bundles, verify signatures, apply updates atomically (A/B partition scheme).
- Implement automatic rollback on update failure.
- Implement backup service: scheduled backups, full and incremental backup modes.
- Implement restore: full restore from backup.
- (Optional, opt-in) Cloud backup target for backup redundancy.
- Display update progress and backup status on native display (system overlay).
- Display rollback status and error recovery screens on native display.

### Risks

- A/B partition scheme doubles storage requirements for the root filesystem. Evaluate trade-offs against fallback initramfs approach.
- Update verification chain (signing keys, certificate management) must be designed carefully to avoid key compromise exposure.

### Exit criteria

- Update service can download, verify, and apply a signed update.
- System boots into the new update automatically.
- If new update fails to boot, system automatically rolls back to previous slot.
- Native display shows update progress and rollback status.
- Backup service can create and restore a full household backup.
- Restore preserves all profiles, settings, and module data.

---

## Phase 6: Polish and Stabilization

Status: Planned

### Outcome

Production-ready release with hardened security, performance optimization, complete documentation, and validated upgrade path.

### Included work

- Security audit and penetration testing.
- Performance benchmarking and optimization (compositor GPU usage, display service latency, input response time).
- Complete user-facing documentation (setup guide including display setup, FAQ, troubleshooting).
- Complete developer documentation (API reference, architecture, contribution guide, module UI SDK docs).
- Hardened default configuration.
- Release automation and version tagging.
- Long-term support (LTS) commitment documentation.
- GUI polish: animation smoothness (60 fps), input latency < 16 ms, overlay transition consistency.
- Test display compatibility across target resolutions (720p, 1080p, 1440p, 4K) and common TV/monitor models.

### Risks

- Security audit may find issues requiring architectural changes. Budget time for remediation.
- Performance optimization for the storage and networking services may be non-trivial with realistic household data volumes.

### Exit criteria

- All acceptance criteria from SPEC.md are met (including native display and input).
- Security audit has no critical or high findings. Medium findings have remediation plans.
- Full backup and restore cycle is validated with realistic data volumes (hundreds of GB).
- GUI output is smooth (60 fps) with input latency under 16 ms.
- Display output validated across target resolutions and common TV/monitor models.
- Documentation is complete and reviewed.
- Release artifacts are reproducible and signed.
