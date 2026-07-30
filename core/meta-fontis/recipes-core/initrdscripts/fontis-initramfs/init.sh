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

# TPM PCR extend for measured boot
if [ -c /dev/tpm0 ]; then
    tpm2_pcrextend 10:sha256=0000000000000000000000000000000000000000000000000000000000000000
fi

# Unlock root partition (LUKS with TPM sealing)
if [ -b /dev/sda2 ]; then
    tpm2_unseal -c 0x81000000 -p 0 > /tmp/luks-key 2>/dev/null || true
    if [ -s /tmp/luks-key ]; then
        cryptsetup luksOpen --key-file=/tmp/luks-key /dev/sda2 cryptroot
        rm -f /tmp/luks-key
    else
        cryptsetup luksOpen /dev/sda2 cryptroot
    fi
fi

# Open dm-verity device on top of unlocked root
if [ -b /dev/mapper/cryptroot ] && [ -b /dev/sda3 ]; then
    veritysetup open /dev/mapper/cryptroot verity-root /dev/sda3 /etc/roothash.pem
fi

# Mount verified root filesystem
if [ -b /dev/mapper/verity-root ]; then
    mount -o ro /dev/mapper/verity-root /root
else
    echo "ERROR: verity-root device not available, dropping to recovery shell"
    exec /bin/sh
fi

# Clean up
udevadm settle

# Switch to real root
exec switch_root /root /sbin/init
