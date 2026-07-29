#!/usr/bin/env bash
set -euo pipefail

echo "=== fontis-platform validation ==="

errors=0

# Check required files exist
required_files=(
    "AGENTS.md"
    "SPEC.md"
    "ROADMAP.md"
    "TASKS.md"
    "CLAUDE.md"
    "README.md"
    "CONTRIBUTING.md"
    "SECURITY.md"
    "CODE_OF_CONDUCT.md"
    ".gitignore"
    "Makefile"
    "docs/ARCHITECTURE.md"
    "docs/STANDARDS.md"
)

for f in "${required_files[@]}"; do
    if [[ ! -f "$f" ]]; then
        echo "ERROR: Missing required file: $f"
        errors=$((errors + 1))
    fi
done

# Verify core/ and runtime/ directory structure exists
required_dirs=(
    "core/kernel"
    "core/boot"
    "core/packages"
    "core/hal"
    "runtime/identity"
    "runtime/auth"
    "runtime/networking"
    "runtime/storage"
    "runtime/logging"
    "runtime/module-lifecycle"
    "runtime/updates"
    "runtime/backup"
    "runtime/marketplace-client"
    "contracts/protobuf"
)

for d in "${required_dirs[@]}"; do
    if [[ ! -d "$d" ]]; then
        echo "WARNING: Missing directory (create when implementing): $d"
    fi
done

# Check for forbidden patterns
if find . -name "*.pem" -o -name "*.key" -o -name "*.crt" 2>/dev/null | grep -q .; then
    echo "ERROR: Private key or certificate file found in repository"
    errors=$((errors + 1))
fi

if find runtime/ -name "*.sqlite" 2>/dev/null | grep -q .; then
    echo "ERROR: SQLite database file found in runtime/ (runtime data, not source)"
    errors=$((errors + 1))
fi

if [[ $errors -eq 0 ]]; then
    echo "=== validation passed ==="
else
    echo "=== validation failed with $errors error(s) ==="
    exit 1
fi
