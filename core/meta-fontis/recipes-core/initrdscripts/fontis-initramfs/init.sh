#!/bin/sh
# Fontis initramfs init script
# - TPM 2.0 measured boot with event log
# - LUKS root unlock (TPM-sealed key with passphrase fallback for development)
# - dm-verity root filesystem verification
# - Boot counter for A/B update tracking
# - Recovery shell with diagnostic information

export SHELL=/bin/sh
export PATH=/sbin:/usr/sbin:/bin:/usr/bin

mount -t proc proc /proc
mount -t sysfs sysfs /sys
mount -t devtmpfs devtmpfs /dev
mount -t tmpfs tmpfs /var
mkdir -p /var/log

echo "Fontis initramfs: starting boot sequence"

TPM_DEV="/dev/tpm0"
CRYPTROOT_DEV="/dev/disk/by-partlabel/cryptroot"
VERITY_DEV="/dev/disk/by-partlabel/verity"

recovery_shell() {
    local reason="$1"
    echo ""
    echo "=============================================="
    echo "Fontis Recovery Shell"
    echo "Reason: $reason"
    echo "=============================================="
    echo "Diagnostics:"
    echo "  cat /proc/cmdline     - kernel command line"
    echo "  ls -la /dev/disk/by-partlabel/ - partition labels"
    echo "  dmesg                 - kernel log"
    echo "  tpm2_pcrread          - TPM PCR values"
    echo "  cat /var/log/tpm_event_log - TPM event log"
    echo ""
    export PS1="(fontis-recovery) # "
    exec /bin/sh
}

udevadm trigger --action=add
udevadm settle --timeout=30

EVENTLOG="/var/log/tpm_event_log"

log_tpm_event() {
    local pcr="$1" event_type="$2" description="$3" hash="$4"
    local timestamp
    timestamp=$(date -u '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null || echo "unknown")
    echo "$timestamp pcr:$pcr type:$event_type hash:$hash desc:$description" >> "$EVENTLOG"
}

pcr_extend_and_log() {
    local pcr="$1" event_type="$2" description="$3" data="$4"

    if [ -z "$data" ]; then
        echo "Warning: no data to measure for $description, skipping"
        return 1
    fi

    local hash
    hash=$(echo "$data" | sha256sum | cut -d' ' -f1)

    if [ -c "$TPM_DEV" ]; then
        tpm2_pcrextend "$pcr:sha256=$hash" 2>/dev/null || {
            echo "ERROR: tpm2_pcrextend(pcr $pcr) failed for $description"
            recovery_shell "TPM PCR extend failed for $description"
        }
        log_tpm_event "$pcr" "$event_type" "$description" "$hash"
    fi

    echo "Measured: pcr $pcr $description"
}

wait_for_tpm() {
    local timeout=10
    while [ $timeout -gt 0 ] && [ ! -c "$TPM_DEV" ]; do
        sleep 1
        timeout=$((timeout - 1))
    done
    if [ ! -c "$TPM_DEV" ]; then
        echo "Warning: TPM not found, continuing without TPM"
        return 1
    fi
    return 0
}

# Boot counter for A/B update tracking
BOOT_COUNTER_NV_INDEX="0x01C10101"

update_boot_counter() {
    if [ ! -c "$TPM_DEV" ]; then
        return
    fi
    local count=0
    if command -v tpm2_nvread >/dev/null 2>&1; then
        count=$(tpm2_nvread "$BOOT_COUNTER_NV_INDEX" 2>/dev/null | od -An -td4 | tr -d ' ') || true
    fi
    if [ -z "$count" ] || [ "$count" -eq 0 ] 2>/dev/null; then
        # NVRAM index doesn't exist yet; define a 4-byte counter with default auth
        tpm2_nvdefine -C o "$BOOT_COUNTER_NV_INDEX" -s 4 -a "ownerread|ownerwrite|nt=1" 2>/dev/null || true
        count=0
    fi
    count=$((count + 1))
    echo "$count" | tpm2_nvwrite "$BOOT_COUNTER_NV_INDEX" 2>/dev/null || true
    echo "Boot counter: $count"
    pcr_extend_and_log 11 "BOOT_COUNTER" "boot counter increment" "$count"
}

# Record initial PCR state for diagnosis
if [ -c "$TPM_DEV" ]; then
    tpm2_pcrread sha256:0,4,8,10,11 2>/dev/null > /var/log/tpm_initial_pcrs.txt || true
fi

# Wait for TPM before first PCR measurement
wait_for_tpm || true

# Stage 1: measure kernel version into PCR 10
KERNEL_INFO=$(cat /proc/version 2>/dev/null || echo "unknown")
pcr_extend_and_log 10 "KERNEL_VERSION" "kernel version" "$KERNEL_INFO"

# Stage 2: measure kernel cmdline into PCR 10
CMDLINE=$(cat /proc/cmdline)
pcr_extend_and_log 10 "CMDLINE" "kernel command line" "$CMDLINE"

update_boot_counter

# Extract dm-verity root hash from kernel cmdline
ROOTHASH=$(sed -n 's/.*roothash=\([^ ]*\).*/\1/p' /proc/cmdline)

if [ -n "$ROOTHASH" ]; then
    pcr_extend_and_log 11 "ROOTHASH" "dm-verity root hash" "$ROOTHASH"
fi

# ---- Root unlock ----

LUKS_KEY_FILE=""

unlock_root_tpm() {
    umask 077
    LUKS_KEY_FILE=$(mktemp -t luks-key.XXXXXXXX) || return 1
    trap 'rm -f "$LUKS_KEY_FILE"' EXIT
    tpm2_unseal -c 0x81000000 -p 0 > "$LUKS_KEY_FILE" 2>/dev/null || return 1
    if [ -s "$LUKS_KEY_FILE" ]; then
        cryptsetup luksOpen --key-file="$LUKS_KEY_FILE" "$CRYPTROOT_DEV" cryptroot && return 0
    fi
    return 1
}

unlock_root_passphrase() {
    echo ""
    echo "TPM unlock failed. Enter LUKS passphrase for $CRYPTROOT_DEV:"
    echo "(Ctrl+D to abort and drop to recovery shell)"
    cryptsetup luksOpen "$CRYPTROOT_DEV" cryptroot
}

if [ -b "$CRYPTROOT_DEV" ]; then
    if ! unlock_root_tpm; then
        echo "TPM key unseal failed."
        if grep -q "fontis.devel" /proc/cmdline 2>/dev/null; then
            unlock_root_passphrase || recovery_shell "LUKS unlock failed"
        else
            sleep 1
            if ! unlock_root_tpm; then
                recovery_shell "TPM LUKS unlock failed (production mode, use fontis.devel for passphrase fallback)"
            fi
        fi
    fi

    if [ -b /dev/mapper/cryptroot ]; then
        ROOT_UUID=$(blkid /dev/mapper/cryptroot -s UUID -o value 2>/dev/null || echo "unknown")
        pcr_extend_and_log 11 "CRYPTROOT" "LUKS root volume UUID" "$ROOT_UUID"
    fi
fi

# ---- dm-verity verification ----

if [ ! -b /dev/mapper/cryptroot ] || [ ! -b "$VERITY_DEV" ] || [ -z "$ROOTHASH" ]; then
    recovery_shell "dm-verity prerequisites missing (cryptroot, verity device, or roothash)"
fi

veritysetup open /dev/mapper/cryptroot verity-root "$VERITY_DEV" "$ROOTHASH" || {
    recovery_shell "dm-verity verification FAILED (root hash mismatch)"
}

pcr_extend_and_log 11 "VERITY" "dm-verity verification passed" "$ROOTHASH"

# ---- Mount verified root ----

if [ -b /dev/mapper/verity-root ]; then
    mount -o ro /dev/mapper/verity-root /root || recovery_shell "Failed to mount verified root"
else
    recovery_shell "verity-root device not available"
fi

pcr_extend_and_log 11 "BOOT_COMPLETE" "initramfs boot completed" "ok"

udevadm settle --timeout=5
rm -f "$LUKS_KEY_FILE"
trap - EXIT

echo "Fontis initramfs: switching to platform runtime"
exec switch_root /root /sbin/init
