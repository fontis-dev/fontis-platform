# Boot Chain

This directory contains boot chain configuration, initramfs scripts,
and secure boot documentation.

- `initramfs/` - Initramfs scripts and configuration
- `secure-boot/` - UEFI Secure Boot key management and signing scripts
- `systemd-boot/` - Bootloader configuration

Source of truth for boot recipes:
    core/meta-fontis/recipes-bsp/systemd-boot/
    core/meta-fontis/recipes-core/initrdscripts/fontis-initramfs/
