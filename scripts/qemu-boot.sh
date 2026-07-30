#!/bin/bash
# Boot a Fontis platform image in QEMU with UEFI and TPM emulation.
# Supports both full WIC images and development kernel+initramfs boots.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# --- Configuration ---
MACHINE="${MACHINE:-fontis-dev}"
MEMORY="${MEMORY:-2048}"
SMP="${SMP:-2}"
ACCEL="${ACCEL:-kvm:tcg}"

BUILD_DIR="${BUILD_DIR:-$PROJECT_DIR/build}"
IMAGE_FILE="${IMAGE_FILE:-$BUILD_DIR/fontis-dev.wic}"
KERNEL="${KERNEL:-}"
INITRAMFS="${INITRAMFS:-}"
SECURE_BOOT="${SECURE_BOOT:-false}"
SECURE_BOOT_KEYS="${SECURE_BOOT_KEYS:-$HOME/.cache/fontis/secure-boot-keys}"

SWTPM_SOCKET="${SWTPM_SOCKET:-/tmp/fontis-swtpm.sock}"
SWTPM_STATE="${SWTPM_STATE:-/tmp/fontis-swtpm}"
SWTPM_PID_FILE="${SWTPM_PID_FILE:-/tmp/fontis-swtpm.pid}"

OVMF_CODE="${OVMF_CODE:-}"
OVMF_VARS="${OVMF_VARS:-}"

usage() {
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  --image <path>       Path to WIC image (default: $IMAGE_FILE)"
    echo "  --kernel <path>      Kernel bzImage (for direct kernel boot)"
    echo "  --initramfs <path>   Initramfs cpio.gz (for direct kernel boot)"
    echo "  --secure-boot        Enable UEFI Secure Boot"
    echo "  --mem <MB>           Memory in MB (default: $MEMORY)"
    echo "  --smp <n>            CPU count (default: $SMP)"
    echo "  --no-tpm             Disable TPM emulation"
    echo "  --help               Show this help"
    exit 0
}

# --- Parse arguments ---
USE_TPM=true

while [ $# -gt 0 ]; do
    case "$1" in
        --image) IMAGE_FILE="$2"; shift 2 ;;
        --kernel) KERNEL="$2"; shift 2 ;;
        --initramfs) INITRAMFS="$2"; shift 2 ;;
        --secure-boot) SECURE_BOOT=true; shift ;;
        --mem) MEMORY="$2"; shift 2 ;;
        --smp) SMP="$2"; shift 2 ;;
        --no-tpm) USE_TPM=false; shift ;;
        --help) usage ;;
        *) echo "Unknown option: $1"; usage ;;
    esac
done

echo "=== Fontis QEMU Boot ==="
echo "Machine:    $MACHINE"
echo "Memory:     ${MEMORY}MB"
echo "CPUs:       $SMP"
echo "SecureBoot: $SECURE_BOOT"
echo "TPM:        $USE_TPM"
echo ""

# --- Locate OVMF ---
locate_ovmf() {
    local paths=(
        /usr/share/ovmf/OVMF.fd
        /usr/share/edk2/x64/OVMF.fd
        /usr/share/edk2-ovmf/x64/OVMF.fd
        /usr/share/qemu/ovmf-x86_64.bin
        /usr/share/OVMF/OVMF_CODE.fd
    )
    for p in "${paths[@]}"; do
        if [ -f "$p" ]; then
            echo "$p"
            return 0
        fi
    done
    return 1
}

locate_ovmf_vars() {
    local code_path="$1"
    local dir
    dir=$(dirname "$code_path")
    local candidates=(
        "$dir/OVMF_VARS.fd"
        "$dir/OVMF_VARS_MS.fd"
    )
    for p in "${candidates[@]}"; do
        if [ -f "$p" ]; then
            echo "$p"
            return 0
        fi
    done
    return 1
}

# --- Determine boot mode ---
BOOT_MODE="image"

if [ -n "$KERNEL" ] && [ -n "$INITRAMFS" ]; then
    BOOT_MODE="direct"
elif [ -n "$KERNEL" ] || [ -n "$INITRAMFS" ]; then
    echo "Error: --kernel and --initramfs must be used together"
    exit 1
fi

# --- Verify image ---
if [ "$BOOT_MODE" = "image" ] && [ ! -f "$IMAGE_FILE" ]; then
    echo "Error: image not found: $IMAGE_FILE"
    echo "Build it first with: make build-core"
    echo "Or use --kernel and --initramfs for direct kernel boot."
    exit 1
fi

# --- Locate OVMF ---
if [ -z "$OVMF_CODE" ]; then
    OVMF_CODE=$(locate_ovmf || true)
fi

if [ -z "$OVMF_CODE" ]; then
    echo "Error: OVMF firmware not found."
    echo "Install ovmf package or set OVMF_CODE path."
    echo "  Debian/Ubuntu: apt install ovmf"
    echo "  Fedora:        dnf install edk2-ovmf"
    exit 1
fi

if [ -z "$OVMF_VARS" ]; then
    OVMF_VARS=$(locate_ovmf_vars "$OVMF_CODE" || true)
fi

echo "OVMF code: $OVMF_CODE"
echo "OVMF vars: ${OVMF_VARS:-not found (using temporary copy)}"

# --- Generate a writable vars file ---
TEMP_VARS=$(mktemp /tmp/fontis-ovmf-vars.XXXXXXXX.fd)
trap 'rm -f "$TEMP_VARS"' EXIT

if [ -n "$OVMF_VARS" ]; then
    cp "$OVMF_VARS" "$TEMP_VARS"
else
    # Create an empty 256KB vars file (standard OVMF vars size)
    dd if=/dev/zero bs=256K count=1 of="$TEMP_VARS" 2>/dev/null
fi

# --- Build QEMU arguments ---
QEMU_ARGS=(
    -machine q35,accel="$ACCEL",smm=on
    -cpu max
    -smp "$SMP"
    -m "$MEMORY"
    -global driver=cfi.pflash01,property=secure,value=on
    -drive if=pflash,format=raw,unit=0,file="$OVMF_CODE",readonly=on
    -drive if=pflash,format=raw,unit=1,file="$TEMP_VARS"
    -serial stdio
    -vga virtio
    -display none
    -device virtio-rng-pci
    -no-reboot
)

# Networking (user-mode slirp)
QEMU_ARGS+=(
    -nic user,model=virtio-net-pci
)

# --- Boot source ---
if [ "$BOOT_MODE" = "direct" ]; then
    echo "Boot mode: direct kernel boot"
    QEMU_ARGS+=(
        -kernel "$KERNEL"
        -initrd "$INITRAMFS"
        -append "console=ttyS0,115200 root=/dev/dm-0 ro rootfstype=squashfs fontis.devel"
    )
else
    echo "Boot mode: WIC image ($IMAGE_FILE)"
    QEMU_ARGS+=(
        -drive file="$IMAGE_FILE",format=raw,if=virtio,aio=threads,cache=writeback
    )
fi

# --- TPM emulation ---
if [ "$USE_TPM" = true ]; then
    # Check for swtpm
    if ! command -v swtpm &>/dev/null; then
        echo "Warning: swtpm not found, disabling TPM emulation"
        echo "Install: apt install swtpm    (Debian/Ubuntu)"
    else
        # Stop existing swtpm
        if [ -f "$SWTPM_PID_FILE" ]; then
            kill "$(cat "$SWTPM_PID_FILE")" 2>/dev/null || true
            rm -f "$SWTPM_PID_FILE"
        fi
        rm -f "$SWTPM_SOCKET"

        mkdir -p "$SWTPM_STATE"

        swtpm socket \
            --tpm2 \
            --tpmstate dir="$SWTPM_STATE" \
            --ctrl type=unixio,path="$SWTPM_SOCKET" \
            --pid file="$SWTPM_PID_FILE" \
            --daemon

        echo "swtpm started (socket: $SWTPM_SOCKET, pid: $(cat "$SWTPM_PID_FILE"))"

        # Ensure swtpm and OVMF vars are cleaned up on exit
        cleanup_swtpm() {
            if [ -f "$SWTPM_PID_FILE" ]; then
                kill "$(cat "$SWTPM_PID_FILE")" 2>/dev/null || true
                rm -f "$SWTPM_PID_FILE" "$SWTPM_SOCKET"
            fi
            rm -f "$TEMP_VARS"
        }
        trap cleanup_swtpm EXIT

        QEMU_ARGS+=(
            -chardev socket,id=chrtpm,path="$SWTPM_SOCKET"
            -tpmdev emulator,id=tpm0,chardev=chrtpm
            -device tpm-tis,tpmdev=tpm0
        )
    fi
fi

# --- Secure Boot ---
if [ "$SECURE_BOOT" = true ]; then
    if [ -f "$SECURE_BOOT_KEYS/db.key" ] && [ -f "$SECURE_BOOT_KEYS/db.crt" ]; then
        echo "Secure Boot enabled (keys: $SECURE_BOOT_KEYS)"
        echo "NOTE: Key enrollment must be performed manually in the UEFI menu"
        echo "      or by pre-configuring OVMF_VARS with enrolled keys."
    else
        echo "Warning: Secure Boot keys not found, booting without Secure Boot"
        SECURE_BOOT=false
    fi
fi

echo ""
echo "=== Starting QEMU ==="
echo "Arguments: ${QEMU_ARGS[*]}"
echo ""
echo "QEMU console output below (Ctrl+A X to exit):"
echo ""

# Execute QEMU
qemu-system-x86_64 "${QEMU_ARGS[@]}"

EXIT_CODE=$?
echo ""
echo "QEMU exited with code $EXIT_CODE"
exit $EXIT_CODE
