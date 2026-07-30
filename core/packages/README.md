# Package and Image Definitions

This directory documents the Fontis image types and package groups.

Image recipes are in:
    core/meta-fontis/recipes-core/images/
        fontis-base-image.bb    - Minimal initramfs image
        fontis-full-image.bb    - Full platform image

Package groups are in:
    core/meta-fontis/recipes-core/packagegroups/
        packagegroup-fontis-base.bb      - System essentials
        packagegroup-fontis-graphics.bb  - Display/input stack
        packagegroup-fontis-security.bb  - Security tooling

WIC kickstart file:
    core/meta-fontis/wic/fontis-dev.wks
