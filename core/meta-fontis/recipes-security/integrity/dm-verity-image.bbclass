# dm-verity image class
# Generates a hash tree for the root filesystem and attaches roothash
# to the kernel command line for verified boot.

inherit image

# Enable dm-verity for the root image
IMAGE_VERITY = "1"

# Verity hash format: sha256
VERITY_HASH_ALGORITHM ?= "sha256"

# Append roothash to kernel cmdline for initramfs verification
python __anonymous() {
    verity_hash = d.getVar('VERITY_HASH')
    if verity_hash:
        append = d.getVar('APPEND') or ''
        d.setVar('APPEND', '%s roothash=%s' % (append, verity_hash))
}
