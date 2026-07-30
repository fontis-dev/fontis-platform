# fontis-platform Architecture

## Overview

fontis-platform is a Linux-based appliance operating system for the Fontis household device. It is built on Debian Stable with custom Fontis platform services. It consists of two layers:

- **Core OS** — Debian Stable base (minimal install, hardened kernel, boot chain, initramfs) plus the Rust Hardware Abstraction Layer (`core/hal/`). The base OS is updated via atomic A/B image updates; the HAL is specific to Fontis hardware.
- **Platform Runtime** (`runtime/`) — a set of independent Go services that provide the household-facing functionality. Each service is an isolated process communicating via gRPC.

---

## Architecture Diagram

```
                    ┌──────────────────────────────────────┐   ┌───────────────────────┐
                    │  Web UI (fontis-web)                  │   │  Native Display        │
                    │  Mobile App (fontis-mobile)           │   │  (HDMI/DP → TV/Monitor)│
                    └────────────┬─────────────────────────┘   └───────────┬───────────┘
                                 │ HTTPS/TLS 1.3                          │ Wayland protocol
                                 ▼                                         ▼
┌──────────────────────────────────────────────────────┐   ┌───────────────────────────┐
│  API Gateway (envoy/traefik)                          │   │  Wayland Compositor       │
│  - TLS termination                                    │   │  (wlroots-based)          │
│  - Authentication verification                        │   │  - DRM/KMS output         │
│  - Rate limiting                                       │   │  - GPU-accelerated render │
│  - Request routing to runtime services                │   │  - Input routing          │
└────┬──────────┬──────────┬──────────┬──────────┬──────┘   │  - Module surfaces        │
     │          │          │          │          │          └───────────┬───────────────┘
     ▼          ▼          ▼          ▼          ▼                      │ gRPC / IPC
┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────────────┐   │
│ Identity│ │  Auth   │ │ Storage │ │Network  │ │ Module Lifecycle│   │
│ Service │ │ Service │ │ Service │ │ Service │ │ Service         │   │
│ (Go)    │ │ (Go)    │ │ (Go)    │ │ (Go)    │ │ (Go)            │   │
└────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────────┬────────┘   │
     │           │           │           │               │            │
     └───────────┴───────────┴───────────┴───────────────┘────────────┘
                          │
                     gRPC over local Unix sockets
                     Mutual TLS
                     Protobuf contracts
                          │
                          ▼
┌────────────────────────────────────────────────────────────────────────────────────┐
│  Logging Service (Go)      │   Display Service (Go)     │  Update Service (Go)     │
│  - Structured log collection│  - Compositor lifecycle    │  - Update download/verify│
│  - Log rotation and query   │  - UI state coordinator    │  - A/B update and rollback│
│  - Module log aggregation  │  - Input dispatch routing   │  - Backup Service (Go)   │
└────────────────────────────┘  - System UI (setup overlay)│  - Scheduled backups     │
                                └──────────────────────────┘  - Cloud/network targets │
                                                              └───────────────────────┘
                          │
                          ▼
┌────────────────────────────────────────────────────────────────────────────────────┐
│  Hardware Abstraction Layer (Rust)                                                   │
│  core/hal/                                                                            │
│  - Disk/block device access                                                           │
│  - TPM operations                                                                      │
│  - GPIO/sensor access, IR receiver                                                    │
│  - Power management                                                                    │
│  - GPU/DRM access (via Mesa)                                                          │
│  - Audio (ALSA/PipeWire)                                                              │
└───────────────────────────────────────┬────────────────────────────────────────────────┘
                                        │
                                        ▼
┌────────────────────────────────────────────────────────────────────────────────────┐
│  Linux Kernel + Boot Chain                                                         │
│  - UEFI Secure Boot → signed bootloader → signed kernel                            │
│  - Measured boot (TPM PCRs)                                                        │
│  - Full disk encryption (LUKS + TPM)                                                │
│  - dm-verity for root filesystem verification                                      │
│  - SELinux/AppArmor for mandatory access control                                   │
│  - DRM/KMS, evdev, libinput, ALSA/snd-hda-intel drivers                           │
└────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Core OS Layer

The Core OS is Debian Stable with a minimal, hardened package selection. The kernel is the Debian Stable kernel with required hardware drivers (storage, networking, TPM, UEFI, GPU, input, audio). Custom kernel config fragments from the Yocto layer informed the Debian kernel module selection.

### Boot Chain

- **Bootloader:** systemd-boot (UEFI) — included with Debian, configured for Fontis.
- **Secure Boot:** UEFI Secure Boot with custom Platform Key (PK), Key Exchange Key (KEK), and Signature Database (db). Keys are generated with `core/boot/secure-boot/gen-keys.sh`. EFI binaries (bootloader, kernel) are signed with the db key.
- **Measured Boot:** TPM 2.0 measures each boot stage. The initramfs (Debian initramfs-tools with Fontis hooks) records an event log and extends PCRs for kernel, cmdline, root hash, and boot state.
- **Root filesystem:** dm-verity protected (via Debian's veritysetup). Read-only root with encrypted overlay for writable state.
- **Initramfs:** Built with Debian's initramfs-tools, customized with Fontis hooks for TPM event logging, LUKS unlock (TPM-sealed key with passphrase fallback), and dm-verity verification.

### Hardware Abstraction (`core/hal/`)

The HAL is the only way runtime services access hardware. It is written in Rust for memory safety and zero-cost abstractions. It exposes a C-compatible FFI that the Go runtime services call through cgo.

**HAL responsibilities:**
- Block device enumeration and I/O.
- TPM 2.0 operations (PCR read/extend, key sealing/unsealing, attestation).
- Power management (shutdown, reboot, suspend, wake events).
- Thermal and sensor monitoring.
- LED indicators and physical button input.
- GPU abstraction (DRM connector enumeration, mode setting, framebuffer allocation — consumed by compositor, not runtime services).
- Audio abstraction (ALSA/PipeWire device enumeration, volume control, sink selection).

---

## Display and Input Stack

The native display and input stack is a new subsystem that renders the Fontis user interface directly on the device's HDMI/DisplayPort output. It runs alongside the existing web/mobile UI paths and shares the same runtime services.

### Layers

1. **Linux Kernel** — DRM/KMS drivers (i915, amdgpu, etc.), evdev input subsystem, ALSA audio drivers.
2. **Wayland Compositor** (wlroots-based) — manages display outputs (HDMI/DP), GPU-accelerated rendering via Mesa, input device capture via libinput, exposes Wayland protocol for client surfaces.
3. **Display Service** (Go runtime service) — owns compositor lifecycle, coordinates UI state across modules, routes input events, renders system overlays (setup wizard, volume indicator, system menu).
4. **Module UI** — each module that provides a native UI runs as a Wayland client connected to the compositor. Modules declare surfaces, receive input events, and render their content through the compositor. Module isolation: Wayland surfaces are sandboxed by the compositor (a misbehaving module cannot read or corrupt another module's surface).

### Input Devices

| Device | Connection | Protocol | Driver |
|--------|-----------|----------|--------|
| IR remote (included) | Built-in IR receiver | LIRC / input-event | kernel gpio-ir / mceusb |
| Game controller | Bluetooth | HID over GATT | kernel hid-generic |
| Keyboard | USB | HID | kernel usbhid |
| Mouse | USB | HID | kernel usbhid |

All input devices are captured by the Wayland compositor via libinput. The compositor routes input events to the display service, which dispatches them to the appropriate consumer (focused module, system UI, or a display service action).

### Display Resolutions

| Mode | Native resolution | Refresh | Notes |
|------|------------------|---------|-------|
| 720p | 1280×720 | 60 Hz | Minimum supported |
| 1080p | 1920×1080 | 60 Hz | Target for TV output |
| 1440p | 2560×1440 | 60 Hz | Monitor use |
| 4K | 3840×2160 | 30/60 Hz | Max supported (HDMI 2.0) |

Multi-monitor support is deferred to a future version.

---

## Platform Runtime (`runtime/`)

### Design Principles

1. **Service isolation.** Each service runs as a separate Linux user with its own filesystem namespace, network namespace, and capability set. No service has access to another service's data store.
2. **Communication through contracts.** All inter-service communication uses gRPC with Protobuf schemas defined in `contracts/protobuf/`. A schema change is a breaking change.
3. **Local Unix sockets.** Services communicate over gRPC on local Unix sockets with mutual TLS. There is no TCP networking for inter-service communication.
4. **Stateless services.** Services store state in their own SQLite databases (managed by the storage service or directly). No service relies on in-memory state that cannot be recovered.
5. **Health and readiness.** Every service exposes a gRPC health check endpoint and a readiness endpoint.

### Service Descriptions

#### Identity Service
- **Owner:** tbd
- **Data:** SQLite database of households, profiles, identity metadata.
- **API:** Create/read/update/delete households and profiles. Verify identity claims.
- **Dependencies:** Storage service (for database volume).

#### Auth Service
- **Owner:** tbd
- **Data:** SQLite database of password hashes, sessions, API tokens.
- **API:** Authenticate (password), create session, validate session, revoke session, create API token, validate API token.
- **Password hashing:** Argon2id.
- **Sessions:** Short-lived (15 min) JWT bearer tokens with refresh token rotation.
- **Dependencies:** Identity service (for user lookup).

#### Storage Service
- **Owner:** tbd
- **Data:** Manages block devices, volume groups, encrypted volumes, filesystem pools.
- **API:** List devices, create volume, mount volume, extend volume, snapshot, create filesystem.
- **Key management:** LUKS passphrase (initial) → TPM-sealed key (production).
- **Dependencies:** HAL (block device access, TPM operations).

#### Networking Service
- **Owner:** tbd
- **API:** List interfaces, scan WiFi, connect to WiFi, configure ethernet, configure VPN, enable/disable DHCP, set static IP, DNS configuration.
- **Dependencies:** HAL (if needed for RF control).

#### Logging Service
- **Owner:** tbd
- **Data:** Structured logs stored in a dedicated SQLite database or log files with rotation.
- **API:** Submit log entry, query logs (by service, level, time range, text search), configure log level.
- **Dependencies:** Storage service (for log volume).

#### Module Lifecycle Service
- **Owner:** tbd
- **Data:** SQLite database of installed modules, their manifests, versions, and states.
- **API:** Install module, uninstall module, start module, stop module, update module, list modules, get module logs.
- **Module isolation:** Each module runs in an OCI container with resource limits (CPU, memory, disk), read-only root filesystem, and no direct hardware access.
- **Dependencies:** Container runtime, storage service (for module data volumes), marketplace client (for downloads).

#### Update Service
- **Owner:** tbd
- **API:** Check for updates, download update, apply update, get update status, rollback.
- **Update mechanism:** A/B partition scheme. The system boots from one slot (A or B). Updates are applied to the inactive slot. On success, the bootloader switches to the updated slot. On failure, the bootloader falls back to the known-good slot.
- **Update verification:** Updates are signed with the fontis release signing key. The signature is verified before installation.
- **Dependencies:** Boot HAL (for slot management), networking service (for download), storage service.

#### Backup Service
- **Owner:** tbd
- **Data:** Backup manifests, schedules, targets.
- **API:** Create backup, list backups, restore from backup, configure backup schedule, configure backup target (local USB, network share, cloud target).
- **Backup content:** Profiles, settings, module data, selected storage volumes.
- **Dependencies:** Storage service, networking service.

#### Display Service
- **Owner:** tbd
- **API:** Register module UI surface, request fullscreen/overlay, route input event, get display info (resolution, refresh rate), configure audio output.
- **Compositor lifecycle:** Starts and manages the Wayland compositor process. Restarts compositor on crash. Handles mode setting.
- **UI state coordinator:** Tracks which module has focus, manages overlay stack (system notifications over module UI), handles transitions (setup wizard → home screen → module).
- **Input dispatch:** Routes input events (remote, gamepad, keyboard, mouse) to the focused module, system UI, or display service actions (e.g., long-pause → system menu).
- **System UI:** The display service owns a system Wayland surface for overlays (volume indicator, input source indicator, system menu) and the setup wizard.
- **Dependencies:** HAL (GPU/DRM via compositor, audio), compositor process, wayland-protocols.

### Marketplace Client
- **Owner:** tbd
- **API:** Browse modules, search modules, get module details, download module, check for module updates.
- **Graceful degradation:** If the marketplace is unavailable (no internet, opt-out), the API returns cached data and does not block any other functionality.
- **Dependencies:** Networking service (for HTTP access to marketplace API).

---

## Security Architecture

### Boot Security
1. UEFI Secure Boot verifies the bootloader signature against the enrolled PK.
2. The bootloader verifies the kernel and initramfs signatures.
3. The kernel measures the initramfs and root filesystem into TPM PCRs.
4. The root filesystem is mounted as dm-verity (read-only, verified at block level).
5. A writable overlay (dm-crypt) provides persistent state.

### Service Security
1. Each runtime service runs as a unique Linux user with minimal capabilities.
2. Services use file-system-level permissions to protect their data stores.
3. Inter-service communication uses mTLS with a per-device CA.
4. The HAL is the only component with direct hardware access. Runtime services call HAL through its FFI.

### Module Security
1. Modules run in OCI containers with no access to host namespaces.
2. Containers have resource limits (CPU, memory, disk IOPS).
3. Containers have read-only root filesystems.
4. Module manifests declare required permissions. The user approves permissions during installation.
5. Modules have no direct hardware access. All hardware access must go through the public runtime APIs.
6. Modules that render a native UI connect to the Wayland compositor as sandboxed clients. The compositor enforces surface isolation: one module cannot read or draw on another module's surface. Module UI rendering is routed through the Wayland protocol, not through direct DRM access.

### Display and Input Security
1. The Wayland compositor runs as an unprivileged user with only the capability to access DRM devices (`CAP_SYS_ADMIN` for KMS is delegated through logind/seatd).
2. Input devices are captured by the compositor. Runtime services and modules never read `/dev/input/*` directly.
3. The display service communicates with the compositor over a dedicated Unix socket. Module management commands (surface register, input route) are authenticated by the compositor.
4. The compositor isolates module Wayland surfaces: surface buffers are not shared between clients. Key event injection between surfaces is prevented.

### Data Security
1. Full disk encryption with LUKS. Keys protected by TPM (production) or passphrase (development).
2. Each service's database is encrypted at rest.
3. Network communication is encrypted with TLS 1.3.
4. Secrets (API keys, certificates) are stored in a dedicated encrypted store, never in configuration files or source code.

---

## Development Environment

### Building
```bash
make build                    # Build runtime services and HAL
make build-runtime            # Build runtime Go services
make build-hal                # Build Rust HAL crate
make build-image              # Build Debian-based core OS image
```

### Running in Emulation
```bash
make qemu                     # Boot the OS image in QEMU
make qemu-with-hal            # Boot with HAL emulation (fake TPM, loopback block devices)
```

### Testing
```bash
make test-unit                # Run all unit tests
make test-integration         # Run integration tests
make test-hal                 # Run HAL tests (requires HAL emulation)
```

### Code Quality
```bash
make fmt                      # Format all code (gofmt, rustfmt, prettier for docs)
make lint                     # Lint all code (golangci-lint, clippy, markdownlint)
make typecheck                # Type-check (go vet, cargo check)
```
