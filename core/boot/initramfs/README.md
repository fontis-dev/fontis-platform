# Initramfs

This directory contains the initramfs configuration for the Fontis platform.

## Purpose

The initramfs provides:

- LUKS volume unlock (passphrase or TPM)
- Root filesystem verification (dm-verity)
- Writable overlay setup (dm-crypt)
- Network setup for remote unlock (optional)
- Recovery shell fallback

## Build

The initramfs is built as part of the Yocto image build. The init script is at:

```
initramfs/
├── init                        Main init script
├── init.d/                     Modular init stages
│   ├── 01-unlock-luks         Unlock encrypted root
│   ├── 02-verify-root         Verify dm-verity hash tree
│   ├── 03-setup-overlay       Setup writable overlay
│   └── 04-start-init          Transition to system init
└── bin/                        Minimal busybox binary
```
