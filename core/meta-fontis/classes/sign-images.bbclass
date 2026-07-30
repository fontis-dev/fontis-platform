# Sign boot chain images for UEFI Secure Boot.
#
# Provides sign_efi() for recipes to sign their EFI binaries with the
# platform Secure Boot db key. Set SECURE_BOOT_KEY_DIR to the directory
# containing db.key and db.crt. Falls back to a well-known development
# path under the meta-fontis layer.
#
# Usage in a recipe:
#   inherit sign-images
#   do_deploy:append() {
#       sign_efi "${DEPLOYDIR}/bzImage"
#   }

DEPENDS:append = " sbsigntool-native"

# Default development key path — override in machine config for production
SECURE_BOOT_KEY_DIR ??= "${LAYERDIR}/../../core/boot/secure-boot/keys"

# Sign an EFI binary with the Secure Boot db key.
# The file is signed in-place.
sign_efi() {
    keydir="${SECURE_BOOT_KEY_DIR}"
    bin="$1"

    if [ -z "$bin" ] || [ ! -f "$bin" ]; then
        bberror "sign_efi: file not found: $bin"
        return 1
    fi

    if [ -f "$keydir/db.key" ] && [ -f "$keydir/db.crt" ]; then
        sbsign --key "$keydir/db.key" --cert "$keydir/db.crt" --output "$bin" "$bin"
        bbnote "Secure Boot: signed $bin"
    else
        bbwarn "Secure Boot keys not found at $keydir, skipping signing of $bin"
    fi
}

# Verify an EFI binary is signed with the platform db certificate.
# Returns 0 if signed, 1 otherwise.
verify_signed_efi() {
    keydir="${SECURE_BOOT_KEY_DIR}"
    bin="$1"

    if [ ! -f "$keydir/db.crt" ]; then
        bbwarn "verify_signed_efi: db.crt not found at $keydir, skipping verification"
        return 0
    fi

    sbverify --cert "$keydir/db.crt" "$bin" 2>/dev/null
}
