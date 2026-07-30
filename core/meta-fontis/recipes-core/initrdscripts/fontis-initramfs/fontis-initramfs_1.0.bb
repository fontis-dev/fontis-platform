SUMMARY = "Fontis initramfs: minimal recovery and boot shell"
LICENSE = "MIT"
LIC_FILES_CHKSUM = "file://${COMMON_LICENSE_DIR}/MIT;md5=0835ade698e0bcf8506ecda2f7b4f302"

SRC_URI = "file://init.sh"

S = "${WORKDIR}"

# This is an initramfs image recipe
inherit core-image

# Minimal base packages for initramfs environment
IMAGE_INSTALL = "busybox tpm2-tools cryptsetup dm-verity kmod"

# Generate cpio archive for initramfs
IMAGE_FSTYPES = "cpio.gz"

# No package management in initramfs
IMAGE_FEATURES = ""

# Install init script into the image
ROOTFS_POSTPROCESS_COMMAND += "install_fontis_init; "

install_fontis_init() {
    install -m 0755 ${WORKDIR}/init.sh ${IMAGE_ROOTFS}/init
}

# Deploy the initramfs artifact as initrd.img for bootloader consumption
do_image_complete[postfuncs] += "deploy_initramfs_artifact"

deploy_initramfs_artifact() {
    install -d ${DEPLOY_DIR_IMAGE}
    if [ -f ${IMGDEPLOYDIR}/${IMAGE_NAME}${IMAGE_NAME_SUFFIX}.cpio.gz ]; then
        ln -sf ${IMAGE_NAME}${IMAGE_NAME_SUFFIX}.cpio.gz ${DEPLOY_DIR_IMAGE}/initrd.img
    fi
}
