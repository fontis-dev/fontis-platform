# Mesa configuration for Fontis: GPU-accelerated rendering, Vulkan, VA-API

PACKAGECONFIG:append = " \
    egl \
    gles2 \
    gbm \
    dri3 \
    ${@bb.utils.contains('MACHINE_FEATURES', 'x86-64', 'va', '', d)} \
    "
