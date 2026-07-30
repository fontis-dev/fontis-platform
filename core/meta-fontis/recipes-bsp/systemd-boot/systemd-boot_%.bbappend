# Fontis systemd-boot configuration
# - Custom timeout
# - Default entry to Fontis platform
# - Signed EFI binary for UEFI Secure Boot

inherit sign-images

FILESEXTRAPATHS:prepend := "${THISDIR}/files:"

SRC_URI += "file://fontis-boot.conf"

# Register the boot entry in EFI deploy metadata for bootimg-efi/systemd-boot.bbclass
SYSTEMD_BOOT_ENTRIES = "fontis-boot.conf"

do_install:append() {
    install -d ${D}/loader/entries
    install -m 644 ${WORKDIR}/fontis-boot.conf ${D}/loader/entries/fontis.conf
    # Deploy to EFI provider path for image inclusion
    install -d ${D}${EFI_PREFIX}/loader/entries
    install -m 644 ${WORKDIR}/fontis-boot.conf ${D}${EFI_PREFIX}/loader/entries/fontis.conf
}

# Sign the bootloader binary after installation
do_deploy:append() {
    bootloader="${DEPLOYDIR}/${BOOTLOADER_EFI_BINARY}"
    if [ -f "$bootloader" ]; then
        sign_efi "$bootloader"
    fi
}

FILES:${PN} += "/loader/entries/fontis.conf ${EFI_PREFIX}/loader/entries/fontis.conf"
