#!/bin/sh
# Fontis initramfs init script
# Mounts root via dm-verity, sets up encrypted volumes, boots into platform runtime

export SHELL=/bin/sh

# Mount essential filesystems
mount -t proc proc /proc
mount -t sysfs sysfs /sys
mount -t devtmpfs devtmpfs /dev

# Wait for devices
udevadm trigger --action=add
udevadm settle --timeout=30

# Extract dm-verity root hash from kernel cmdline (set by dm-verity-image.bbclass)
ROOTHASH=$(sed -n 's/.*roothash=\([^ ]*\).*/\1/p' /proc/cmdline)

# TPM PCR extend for measured boot: measure kernel cmdline (includes roothash)
if [ -c /dev/tpm0 ]; then
    CMDLINE_HASH=$(cat /proc/cmdline | sha256sum | cut -d' ' -f1)
    tpm2_pcrextend 10:sha256="$CMDLINE_HASH" || {
        echo "ERROR: tpm2_pcrextend failed, measured boot compromised"
        exec /bin/sh
    }
fi

# Unlock root partition (LUKS with TPM sealing)
CRYPTROOT_DEV="/dev/disk/by-partlabel/cryptroot"
VERITY_DEV="/dev/disk/by-partlabel/verity"
LUKS_KEY_FILE=""

if [ -b "$CRYPTROOT_DEV" ]; then
    umask 077
    LUKS_KEY_FILE=$(mktemp -t luks-key.XXXXXXXX) || exit 1
    trap 'rm -f "$LUKS_KEY_FILE"' EXIT
    tpm2_unseal -c 0x81000000 -p 0 > "$LUKS_KEY_FILE" 2>/dev/null || true
    if [ -s "$LUKS_KEY_FILE" ]; then
        cryptsetup luksOpen --key-file="$LUKS_KEY_FILE" "$CRYPTROOT_DEV" cryptroot
    else
        echo "ERROR: TPM unseal failed, cannot unlock root unattended"
        exec /bin/sh
    fi
fi

# Open dm-verity device on top of unlocked root
if [ -b /dev/mapper/cryptroot ] && [ -b "$VERITY_DEV" ] && [ -n "$ROOTHASH" ]; then
    veritysetup open /dev/mapper/cryptroot verity-root "$VERITY_DEV" "$ROOTHASH"
fi

# Mount verified root filesystem
if [ -b /dev/mapper/verity-root ]; then
    mount -o ro /dev/mapper/verity-root /root || {
        echo "ERROR: failed to mount verity-root"
        rm -f "$LUKS_KEY_FILE"
        trap - EXIT
        exec /bin/sh
    }
else
    echo "ERROR: verity-root device not available, dropping to recovery shell"
    rm -f "$LUKS_KEY_FILE"
    trap - EXIT
    exec /bin/sh
fi

# Clean up
udevadm settle
rm -f "$LUKS_KEY_FILE"

# Switch to real root
exec switch_root /root /sbin/init
