# NetworkManager configuration for Fontis: WiFi, ethernet, VPN

PACKAGECONFIG:append = " wifi"
PACKAGECONFIG:append = " ${@bb.utils.contains('DISTRO_FEATURES', 'bluetooth', 'bluetooth', '', d)}"

FILESEXTRAPATHS:prepend := "${THISDIR}/files:"

SRC_URI += "file://fontis-nm-config.conf"

do_install:append() {
    install -d ${D}${sysconfdir}/NetworkManager/conf.d
    install -m 644 ${WORKDIR}/fontis-nm-config.conf ${D}${sysconfdir}/NetworkManager/conf.d/99-fontis.conf
}
