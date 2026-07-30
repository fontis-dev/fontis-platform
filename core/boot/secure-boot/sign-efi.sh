#!/bin/bash
# Sign an EFI binary with the Secure Boot db key
# Usage: ./sign-efi.sh <key-dir> <input-efi> [output-efi]
set -euo pipefail

usage() {
    echo "Usage: $(basename "$0") <key-dir> <input-efi> [output-efi]"
    echo ""
    echo "Signs an EFI binary with the db key for UEFI Secure Boot."
    echo "If output-efi is omitted, the input file is signed in-place."
    exit 1
}

KEYDIR="${1:?Missing key directory argument}"
INPUT="${2:?Missing input EFI binary argument}"
OUTPUT="${3:-$INPUT}"

if [ ! -f "$KEYDIR/db.key" ]; then
    echo "Error: db.key not found in $KEYDIR"
    echo "Run gen-keys.sh first to generate keys."
    exit 1
fi

if [ ! -f "$KEYDIR/db.crt" ]; then
    echo "Error: db.crt not found in $KEYDIR"
    exit 1
fi

# Check if already signed
if sbverify --cert "$KEYDIR/db.crt" "$INPUT" 2>/dev/null; then
    echo "Already signed with this certificate, skipping."
    exit 0
fi

sbsign --key "$KEYDIR/db.key" --cert "$KEYDIR/db.crt" --output "$OUTPUT" "$INPUT"
echo "Signed: $INPUT -> $OUTPUT"
