# SELinux policy for Fontis platform services
# Each runtime service runs as an isolated domain

FILESEXTRAPATHS:prepend := "${THISDIR}/files:"

SRC_URI += "file://fontis-platform.te"

do_install:append() {
    install -d ${D}${sysconfdir}/selinux/fontis
    install -m 644 ${WORKDIR}/fontis-platform.te ${D}${sysconfdir}/selinux/fontis/policy.te
}
