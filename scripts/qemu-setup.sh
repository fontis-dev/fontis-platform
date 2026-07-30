#!/bin/bash
# QEMU setup checker for Fontis platform boot chain verification
# Verifies dependencies and creates a test configuration.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "=== Fontis QEMU Boot Environment Setup ==="
echo ""

# --- Dependency checks ---
MISSING=""

check_cmd() {
    if ! command -v "$1" &>/dev/null; then
        echo "  MISSING: $1${2:+ ($2)}"
        MISSING="$MISSING $1"
    else
        echo "  FOUND:   $1${2:+ ($2)}"
    fi
}

echo "Checking required tools..."
check_cmd qemu-system-x86_64 "QEMU x86-64 system emulator"
check_cmd swtpm "Software TPM 2.0 emulator"
check_cmd openssl "Cryptographic key generation"

echo ""
echo "Checking optional tools..."
check_cmd sbsign "EFI binary signing (Secure Boot)"
check_cmd sbverify "EFI signature verification (Secure Boot)"
check_cmd cert-to-efi-sig-list "EFI signature list conversion (Secure Boot)"
check_cmd sign-efi-sig-list "EFI signature list signing (Secure Boot)"

echo ""

# --- Locate OVMF ---
OVMF_SEARCH_PATHS=(
    /usr/share/ovmf/OVMF.fd
    /usr/share/edk2/x64/OVMF.fd
    /usr/share/edk2-ovmf/x64/OVMF.fd
    /usr/share/qemu/ovmf-x86_64.bin
    /usr/share/OVMF/OVMF_CODE.fd
)

OVMF_CODE=""
OVMF_VARS=""

for path in "${OVMF_SEARCH_PATHS[@]}"; do
    if [ -f "$path" ]; then
        OVMF_CODE="$path"
        # Look for vars file alongside code
        vars_candidates=(
            "${path%/*}/OVMF_VARS.fd"
            "${path%/*}/OVMF_VARS_MS.fd"
        )
        for vf in "${vars_candidates[@]}"; do
            if [ -f "$vf" ]; then
                OVMF_VARS="$vf"
                break
            fi
        done
        echo "Found OVMF: $path${OVMF_VARS:+ (vars: $OVMF_VARS)}"
        break
    fi
done

if [ -z "$OVMF_CODE" ]; then
    echo "  NOT FOUND: OVMF UEFI firmware"
    echo "  Install: apt install ovmf       (Debian/Ubuntu)"
    echo "           dnf install edk2-ovmf  (Fedora)"
    MISSING="$MISSING ovmf"
fi

# --- Locate swtpm ---
SWTPM_SOCKET="${SWTPM_SOCKET:-/tmp/fontis-swtpm.sock}"
SWTPM_STATE="${SWTPM_STATE:-/tmp/fontis-swtpm}"
echo ""
echo "swtpm socket: $SWTPM_SOCKET"
echo "swtpm state:  $SWTPM_STATE"

# --- Generate test keys if missing and Secure Boot mode is requested ---
SECURE_BOOT="${SECURE_BOOT:-false}"
SECURE_BOOT_KEYS="${SECURE_BOOT_KEYS:-$HOME/.cache/fontis/secure-boot-keys}"

if [ "$SECURE_BOOT" = "true" ]; then
    if [ -f "$SECURE_BOOT_KEYS/db.key" ] && [ -f "$SECURE_BOOT_KEYS/db.crt" ]; then
        echo ""
        echo "Secure Boot keys: $SECURE_BOOT_KEYS (existing)"
    elif [ -f "$PROJECT_DIR/core/boot/secure-boot/gen-keys.sh" ]; then
        echo ""
        echo "Secure Boot keys not found. Generating development keys..."
        mkdir -p "$SECURE_BOOT_KEYS"
        "$PROJECT_DIR/core/boot/secure-boot/gen-keys.sh" "$SECURE_BOOT_KEYS"
    else
        echo ""
        echo "ERROR: Secure Boot key generator not found at $PROJECT_DIR/core/boot/secure-boot/gen-keys.sh"
        echo "Cannot generate required keys for --secure-boot mode."
        exit 1
    fi
fi

echo ""
echo "=== Summary ==="
if [ -n "$MISSING" ]; then
    echo "Missing dependencies:$MISSING"
    echo "Install them before running qemu-boot.sh."
    exit 1
else
    echo "All dependencies satisfied."
fi
