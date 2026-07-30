SUMMARY = "Fontis display and graphics packages"
LICENSE = "MIT"

inherit packagegroup

RDEPENDS:packagegroup-fontis-graphics = "\
    libdrm \
    libegl \
    libgles2 \
    libgbm \
    libinput \
    libxkbcommon \
    mesa \
    mesa-driver-i915 \
    mesa-driver-radeon \
    pixman \
    wayland \
    wayland-protocols \
    wlroots \
    ${@bb.utils.contains('MACHINE_FEATURES', 'bluetooth', 'bluez5', '', d)} \
    "
