# fontis-platform Roadmap

## Phase 1: Platform Foundation

Status: Planned

### Outcome

Establish the build system, core OS image, and development toolchain. Produce a bootable image for Community Edition (x86-64, UEFI) with secure boot and a minimal initramfs.

### Included work

- Set up Yocto-based build system for the core OS image.
- Configure Linux kernel with required drivers and security features.
- Implement boot infrastructure: UEFI Secure Boot, signed bootloader, initramfs.
- Set up the Rust HAL toolchain and project structure under `core/hal/`.
- Set up the Go runtime project structure with build tooling.
- Create the `contracts/protobuf/` directory with initial identity and auth service schemas.
- Implement CI pipeline: build, lint, format check, unit test.
- Write integration test harness for runtime services.
- Document architecture (docs/ARCHITECTURE.md) and standards (docs/STANDARDS.md).

### Risks

- Yocto build time and complexity may slow initial iteration. Consider a simpler base image for early development.
- GPT-4o may hallucinate kernel configuration options. Verify every kernel config against kernel.org documentation.
- UEFI Secure Boot setup requires signing infrastructure and key management planning.

### Exit criteria

- `make build` produces a bootable x86-64 UEFI image from a clean checkout.
- The image boots on reference hardware (or QEMU) to a kernel panic with a visible initramfs shell.
- The Rust HAL crate compiles with no warnings.
- The Go runtime skeleton compiles with no errors.
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

### Risks

- Getting gRPC and mTLS right on Unix sockets requires careful implementation. Plan to start with plain gRPC and add TLS after basic connectivity works.
- Password hashing and session management must follow security best practices (Argon2id, secure session cookies/tokens).

### Exit criteria

- Identity service manages households and profiles through gRPC API.
- Auth service authenticates users and issues session tokens.
- Integration tests verify the full identity creation and authentication flow.
- Inter-service communication uses mTLS.
- No plaintext secrets in configuration or logs.

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

### Risks

- Storage encryption key management with TPM is complex. Start with passphrase-based LUKS and add TPM sealing later.
- WiFi management requires careful integration with NetworkManager or similar. Evaluate whether to use NetworkManager directly or implement a minimal WiFi daemon.

### Exit criteria

- Storage service can initialize a disk, create an encrypted volume, and mount it.
- Networking service can connect to a WiFi network and obtain an IP address via DHCP.
- Device is discoverable by name on the local network.
- All operations require valid authentication.

---

## Phase 4: Module System

Status: Planned

### Outcome

First-party module can be installed, started, stopped, and removed through the module lifecycle service. Module isolation is enforced through containers.

### Included work

- Implement module lifecycle service: module manifest parsing, container creation, resource limits.
- Create module SDK (in `fontis-modules` repo): types, contract helpers, testing utilities.
- Build one first-party reference module (e.g., a file browser or media viewer).
- Expose module management through Protobuf contracts.
- Integrate module lifecycle with the web UI (in `fontis-web` repo).
- Implement module update and rollback.

### Risks

- Container runtime (Docker/containerd/podman) selection affects module isolation guarantees. Evaluate for memory/disk overhead and security properties.
- Module manifest format needs to balance flexibility with simplicity. Start minimal and extend based on real module needs.

### Exit criteria

- A first-party module can be installed through the module lifecycle service.
- The module starts in an isolated container with resource limits.
- The module is accessible through the web UI.
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

### Risks

- A/B partition scheme doubles storage requirements for the root filesystem. Evaluate trade-offs against fallback initramfs approach.
- Update verification chain (signing keys, certificate management) must be designed carefully to avoid key compromise exposure.

### Exit criteria

- Update service can download, verify, and apply a signed update.
- System boots into the new update automatically.
- If new update fails to boot, system automatically rolls back to previous slot.
- Backup service can create and restore a full household backup.
- Restore preserves all profiles, settings, and module data.

---

## Phase 6: Polish and Stabilization

Status: Planned

### Outcome

Production-ready release with hardened security, performance optimization, complete documentation, and validated upgrade path.

### Included work

- Security audit and penetration testing.
- Performance benchmarking and optimization.
- Complete user-facing documentation (setup guide, FAQ, troubleshooting).
- Complete developer documentation (API reference, architecture, contribution guide).
- Hardened default configuration.
- Release automation and version tagging.
- Long-term support (LTS) commitment documentation.

### Risks

- Security audit may find issues requiring architectural changes. Budget time for remediation.
- Performance optimization for the storage and networking services may be non-trivial with realistic household data volumes.

### Exit criteria

- All acceptance criteria from SPEC.md are met.
- Security audit has no critical or high findings. Medium findings have remediation plans.
- Full backup and restore cycle is validated with realistic data volumes (hundreds of GB).
- Documentation is complete and reviewed.
- Release artifacts are reproducible and signed.
