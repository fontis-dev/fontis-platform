# fontis-platform Architecture

## Overview

fontis-platform is a purpose-built Linux-based operating system and runtime environment for the Fontis household device. It consists of two layers:

- **Core OS** (`core/`) — the kernel, boot chain, base system packages, and hardware abstraction. This layer is minimal and stable. It is updated infrequently and changes are high-risk.
- **Platform Runtime** (`runtime/`) — a set of independent Go services that provide the household-facing functionality. Each service is an isolated process communicating via gRPC.

---

## Architecture Diagram

```
                    ┌──────────────────────────────────────┐
                    │  Web UI (fontis-web)                  │
                    │  Mobile App (fontis-mobile)           │
                    └────────────┬─────────────────────────┘
                                 │ HTTPS/TLS 1.3
                                 ▼
┌──────────────────────────────────────────────────────────────────┐
│  API Gateway (envoy/traefik)                                      │
│  - TLS termination                                                │
│  - Authentication verification                                    │
│  - Rate limiting                                                   │
│  - Request routing to runtime services                            │
└────┬──────────┬──────────┬──────────┬──────────┬─────────────────┘
     │          │          │          │          │
     ▼          ▼          ▼          ▼          ▼
┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────────────┐
│ Identity│ │  Auth   │ │ Storage │ │Network  │ │ Module Lifecycle│
│ Service │ │ Service │ │ Service │ │ Service │ │ Service         │
│ (Go)    │ │ (Go)    │ │ (Go)    │ │ (Go)    │ │ (Go)            │
└────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────────┬────────┘
     │           │           │           │               │
     └───────────┴───────────┴───────────┴───────────────┘
                         │
                    gRPC over local Unix sockets
                    Mutual TLS
                    Protobuf contracts
                         │
                         ▼
┌──────────────────────────────────────────────────────────────────┐
│  Logging Service (Go)      │   Update Service (Go)               │
│  - Structured log collection│  - Update download and verification │
│  - Log rotation and query   │  - A/B update and rollback         │
│  - Module log aggregation  │  - Backup Service (Go)              │
└────────────────────────────┘  - Scheduled backups                │
                                - Cloud/network targets (opt-in)   │
                                └─────────────────────────────────┘
                         │
                         ▼
┌──────────────────────────────────────────────────────────────────┐
│  Hardware Abstraction Layer (Rust)                               │
│  core/hal/                                                        │
│  - Disk/block device access                                       │
│  - TPM operations                                                  │
│  - GPIO/sensor access                                              │
│  - Power management                                                │
└───────────────────────────┬──────────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────────────┐
│  Linux Kernel + Boot Chain                                       │
│  - UEFI Secure Boot → signed bootloader → signed kernel          │
│  - Measured boot (TPM PCRs)                                      │
│  - Full disk encryption (LUKS + TPM)                              │
│  - dm-verity for root filesystem verification                    │
│  - SELinux/AppArmor for mandatory access control                 │
└──────────────────────────────────────────────────────────────────┘
```

---

## Core OS Layer (`core/`)

### Kernel (`core/kernel/`)

- Linux kernel configured for the target architecture (x86-64, ARM64 future).
- Minimal configuration: only required drivers and subsystems.
- Security features enabled: SELinux, integrity subsystem (IMA/EVM), dm-crypt, dm-verity, TPM, audit.
- Kernel image is signed and verified during boot.

### Boot (`core/boot/`)

- **Bootloader:** systemd-boot (simple, supports UEFI Secure Boot, signed binaries).
- **Secure Boot:** UEFI Secure Boot with custom Platform Key (PK) and Key Exchange Key (KEK).
- **Measured Boot:** TPM PCRs record every stage of the boot chain.
- **Initramfs:** Minimal initramfs with cryptsetup, LVM (if applicable), and a recovery shell. The initramfs is signed and verified.
- **Root filesystem:** dm-verity protected. Read-only root with overlay for writable state.

### Hardware Abstraction (`core/hal/`)

The HAL is the only way runtime services access hardware. It is written in Rust for memory safety and zero-cost abstractions. It exposes a C-compatible FFI that the Go runtime services call through cgo.

**HAL responsibilities:**
- Block device enumeration and I/O.
- TPM 2.0 operations (PCR read/extend, key sealing/unsealing, attestation).
- Power management (shutdown, reboot, suspend, wake events).
- Thermal and sensor monitoring.
- LED indicators and physical button input.

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

#### Marketplace Client
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

### Data Security
1. Full disk encryption with LUKS. Keys protected by TPM (production) or passphrase (development).
2. Each service's database is encrypted at rest.
3. Network communication is encrypted with TLS 1.3.
4. Secrets (API keys, certificates) are stored in a dedicated encrypted store, never in configuration files or source code.

---

## Development Environment

### Building
```bash
make build                    # Build all targets
make build-core               # Build core OS image
make build-runtime            # Build runtime Go services
make build-hal                # Build Rust HAL crate
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
