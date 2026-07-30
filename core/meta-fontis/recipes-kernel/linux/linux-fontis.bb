LINUX_VERSION = "6.6.y"
LINUX_VERSION_EXTENSION = "-fontis"

SRCREV = "${AUTOREV}"

SRC_URI = "git://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git;branch=linux-${LINUX_VERSION};protocol=https"

require recipes-kernel/linux/linux-yocto.inc

KCONFIG_MODE = "--alldefconfig"

SRC_URI += "file://fontis.cfg \
            file://security.cfg \
            file://graphics.cfg \
            file://storage-net.cfg \
            file://audio-input.cfg \
            "

COMPATIBLE_MACHINE = "fontis-dev"
