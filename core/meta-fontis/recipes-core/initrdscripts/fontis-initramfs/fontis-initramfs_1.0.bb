SUMMARY = "Fontis initramfs: minimal recovery and boot shell"
LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://${COMMON_LICENSE_DIR}/MIT;md5=0835ade698e0bcf8506ecda2f7b4f302"

SRC_URI = "file://init.sh"

S = "${WORKDIR}"

do_install() {
    install -d ${D}/init.d
    install -m 0755 ${WORKDIR}/init.sh ${D}/init.d/init.sh
}

FILES:${PN} = "/init.d/init.sh"
