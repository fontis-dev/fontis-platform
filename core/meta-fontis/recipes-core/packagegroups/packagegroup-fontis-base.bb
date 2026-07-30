SUMMARY = "Fontis base system packages"
LICENSE = "MIT"

inherit packagegroup

PACKAGES = "\
    packagegroup-fontis-base \
    "

RDEPENDS:packagegroup-fontis-base = "\
    base-files \
    base-passwd \
    busybox \
    ca-certificates \
    cryptsetup \
    dbus \
    dm-verity \
    e2fsprogs \
    glibc \
    kmod \
    libgpiod \
    libpam \
    lm-sensors \
    openssl \
    systemd \
    tzdata \
    udev \
    util-linux \
    ${@bb.utils.contains('MACHINE_FEATURES', 'tpm2', 'tpm2-tss tpm2-abrmd tpm2-tools', '', d)} \
    "
