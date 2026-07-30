# SELinux policy for Fontis platform services
# Each runtime service runs as an isolated domain

FILESEXTRAPATHS:prepend := "${THISDIR}/files:"

SRC_URI += "file://fontis-platform.te"

DEPENDS += "checkpolicy-native policycoreutils-native"

do_compile() {
    checkmodule -M -m ${WORKDIR}/fontis-platform.te -o ${WORKDIR}/fontis-platform.mod
    semodule_package -m ${WORKDIR}/fontis-platform.mod -o ${WORKDIR}/fontis-platform.pp
}

do_install:append() {
    install -d ${D}${sysconfdir}/selinux/fontis
    install -m 644 ${WORKDIR}/fontis-platform.te ${D}${sysconfdir}/selinux/fontis/policy.te
    install -m 644 ${WORKDIR}/fontis-platform.pp ${D}${sysconfdir}/selinux/fontis/policy.pp
}
