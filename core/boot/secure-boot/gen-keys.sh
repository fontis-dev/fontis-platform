#!/bin/bash
# Generate UEFI Secure Boot key hierarchy: PK, KEK, db
# Usage: ./gen-keys.sh [output-dir]
set -euo pipefail

KEYDIR="${1:-$HOME/.config/fontis/secure-boot-keys}"
echo "Generating UEFI Secure Boot keys in: $KEYDIR"
mkdir -p "$KEYDIR"

# Platform Key (PK) — top of the trust chain, controls KEK updates
openssl req -new -x509 -newkey rsa:2048 -subj "/CN=Fontis Platform Key/" \
    -keyout "$KEYDIR/PK.key" -out "$KEYDIR/PK.crt" \
    -days 3650 -nodes -sha256

# Key Exchange Key (KEK) — controls db/dbx updates
openssl req -new -x509 -newkey rsa:2048 -subj "/CN=Fontis Key Exchange Key/" \
    -keyout "$KEYDIR/KEK.key" -out "$KEYDIR/KEK.crt" \
    -days 3650 -nodes -sha256

# Signature Database key (db) — signs authorized EFI binaries
openssl req -new -x509 -newkey rsa:2048 -subj "/CN=Fontis Signature Database Key/" \
    -keyout "$KEYDIR/db.key" -out "$KEYDIR/db.crt" \
    -days 3650 -nodes -sha256

# Convert certificates to EFI Signature List format for enrollment
cert-to-efi-sig-list -g "$(uuidgen)" "$KEYDIR/PK.crt" "$KEYDIR/PK.esl"
cert-to-efi-sig-list -g "$(uuidgen)" "$KEYDIR/KEK.crt" "$KEYDIR/KEK.esl"
cert-to-efi-sig-list -g "$(uuidgen)" "$KEYDIR/db.crt" "$KEYDIR/db.esl"

# Create authenticated EFI variable payloads for enrollment
# PK is self-signed. KEK is signed by PK. db is signed by KEK.
sign-efi-sig-list -k "$KEYDIR/PK.key" -c "$KEYDIR/PK.crt" PK "$KEYDIR/PK.esl" "$KEYDIR/PK.auth"
sign-efi-sig-list -k "$KEYDIR/PK.key" -c "$KEYDIR/PK.crt" KEK "$KEYDIR/KEK.esl" "$KEYDIR/KEK.auth"
sign-efi-sig-list -k "$KEYDIR/KEK.key" -c "$KEYDIR/KEK.crt" db "$KEYDIR/db.esl" "$KEYDIR/db.auth"

echo "Keys generated successfully."
echo "  PK:  $KEYDIR/PK.key, $KEYDIR/PK.crt"
echo "  KEK: $KEYDIR/KEK.key, $KEYDIR/KEK.crt"
echo "  db:  $KEYDIR/db.key, $KEYDIR/db.crt"
echo ""
echo "Enrollment files:"
echo "  $KEYDIR/PK.auth   (enroll Platform Key)"
echo "  $KEYDIR/KEK.auth  (enroll Key Exchange Key)"
echo "  $KEYDIR/db.auth   (enroll Signature Database key)"
echo ""
echo "WARNING: Production keys must be stored offline and secured."
echo "These development keys are suitable for testing only."
