SUMMARY = "Fontis security packages: SELinux, audit, TPM tools"
LICENSE = "MIT"

inherit packagegroup

RDEPENDS:packagegroup-fontis-security = "\
    auditd \
    checksec \
    ${@bb.utils.contains('MACHINE_FEATURES', 'tpm2', 'tpm2-tools tpm2-abrmd', '', d)} \
    "
