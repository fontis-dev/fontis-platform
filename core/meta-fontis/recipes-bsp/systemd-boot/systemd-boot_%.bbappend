# Fontis systemd-boot configuration
# - Custom timeout
# - Default entry to Fontis platform

FILESEXTRAPATHS:prepend := "${THISDIR}/files:"

SRC_URI += "file://fontis-boot.conf"

do_install:append() {
    install -d ${D}/loader/entries
    install -m 644 ${WORKDIR}/fontis-boot.conf ${D}/loader/entries/fontis.conf
}

FILES:${PN} += "/loader/entries/fontis.conf"
