# Secure Boot

This directory contains the UEFI Secure Boot infrastructure for the Fontis device.

## Key hierarchy

```
PK (Platform Key)          — top-level key, signs KEK updates
 └── KEK (Key Exchange Key) — signs db/dbx updates
      └── db (Signature Database) — signs authorized EFI binaries
```

## Usage

### Generate development keys

```bash
./gen-keys.sh ./keys
```

This creates `PK`, `KEK`, and `db` key pairs, certificates, and enrollment files in the specified directory.

### Sign an EFI binary

```bash
./sign-efi.sh ./keys input.efi signed.efi
```

### Enroll keys on a device

Copy the `*.auth` files to the device's EFI system partition and enroll from the UEFI firmware setup menu, or use `sbkeysync` from an OS that supports it:

```bash
sbkeysync --pk=PK.auth --kek=KEK.auth --db=db.auth
```

## Yocto integration

The build system (meta-fontis layer) automatically signs the bootloader and kernel during `make build-core`. Key location is configurable via `SECURE_BOOT_KEY_DIR` in the machine config.

## Production vs development

| | Development | Production |
|---|---|---|
| Key generation | `gen-keys.sh` (in-repo) | Hardware Security Module (HSM) |
| Key storage | `./keys/` in repo | Offline, air-gapped |
| Certificate | Self-signed, 10-year | CA-issued, hardware-backed |
| Signing | Done by build system | Done in secure build environment |
| Enrollment | Manual via UEFI menu | Factory-enrolled |

## Dependencies

- `openssl` — key and certificate generation
- `sbsigntool` — EFI binary signing and verification (`sbsign`, `sbverify`, `sbkeysync`)
- `efitools` — certificate-to-ESL conversion (`cert-to-efi-sig-list`, `sign-efi-sig-list`)
- `uuidgen` — UUID generation for ESL files

Install on Debian/Ubuntu:
```bash
apt install openssl sbsigntool efitools uuid-runtime
```
