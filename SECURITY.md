# Security Policy

## Supported Versions

| Version | Supported |
| --- | --- |
| Latest release | Yes |
| Development (main branch) | Yes (with caveats) |
| Older releases | No |

## Reporting a Vulnerability

Fontis takes security seriously. If you discover a security vulnerability, please report it privately.

**DO NOT** report security vulnerabilities through public GitHub issues.

Instead, email the security team at: **security@fontis.dev** (placeholder - create before release)

If you prefer, you may also report through GitHub's private vulnerability reporting mechanism once the repository is public.

### What to Include

- Description of the vulnerability.
- Steps to reproduce.
- Affected versions and components.
- Potential impact.
- Any suggested fix or mitigation.

### Response Timeline

- Within 48 hours: acknowledgement of receipt.
- Within 7 days: initial assessment and remediation plan.
- Within 30 days: fix released or remediation deployed, depending on severity.

## Security Features

fontis-platform includes security features by design:

- **Secure boot chain:** UEFI Secure Boot, signed bootloader, signed kernel, verified initramfs.
- **Full disk encryption:** LUKS with TPM-bound keys.
- **Measured boot:** TPM PCR attestation of boot integrity.
- **Service isolation:** Each runtime service runs as a unique Linux user with minimal capabilities.
- **Mandatory access control:** SELinux enforcing.
- **Inter-service mTLS:** All runtime service communication uses mutual TLS over Unix sockets.
- **Module sandboxing:** Modules run in OCI containers with no host access.
- **Dependency scanning:** Automated vulnerability scanning for Go modules, Rust crates, and system packages.
- **Signed updates:** All system updates are cryptographically signed and verified before installation.
- **Automatic rollback:** Failed updates automatically roll back to the previous known-good state.

## Responsible Disclosure

We ask that security researchers:

- Give us reasonable time to fix and release a fix before publishing details.
- Avoid exploiting vulnerabilities beyond what is necessary to demonstrate the issue.
- Do not access or modify user data without explicit permission.
- Follow applicable laws and regulations.

We will:

- Respond promptly to all reports.
- Keep you informed of progress.
- Credit you in release notes and changelogs (if desired).
- Handle the report confidentially until a fix is released.
