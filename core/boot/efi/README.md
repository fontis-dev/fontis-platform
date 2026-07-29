# UEFI Secure Boot configuration

This directory will contain the UEFI Secure Boot configuration for the Fontis platform.

## Components

- Platform Key (PK) and Key Exchange Key (KEK) certificates and keys
- Signature database (db) entries for authorized bootloaders and kernels
- Forbidden signature database (dbx) entries for revoked binaries

## Initial setup (development)

For development, self-signed keys are used:

```bash
# Generate PK
openssl req -new -x509 -newkey rsa:2048 -subj "/CN=Fontis Platform Key/" \
    -keyout PK.key -out PK.crt -days 3650 -nodes

# Generate KEK
openssl req -new -x509 -newkey rsa:2048 -subj "/CN=Fontis Key Exchange Key/" \
    -keyout KEK.key -out KEK.crt -days 3650 -nodes

# Generate db certificate
openssl req -new -x509 -newkey rsa:2048 -subj "/CN=Fontis Signature Database/" \
    -keyout db.key -out db.crt -days 3650 -nodes
```

Production signing will use a hardware security module (HSM) with the same key hierarchy.
