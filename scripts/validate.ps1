param([switch]$PassThru)

$ErrorActionPreference = "Stop"
$errors = 0

Write-Host "=== fontis-platform validation ==="

$requiredFiles = @(
    "AGENTS.md",
    "SPEC.md",
    "ROADMAP.md",
    "TASKS.md",
    "CLAUDE.md",
    "README.md",
    "CONTRIBUTING.md",
    "SECURITY.md",
    "CODE_OF_CONDUCT.md",
    ".gitignore",
    "Makefile",
    "docs/ARCHITECTURE.md",
    "docs/STANDARDS.md"
)

foreach ($f in $requiredFiles) {
    if (-not (Test-Path -LiteralPath $f -PathType Leaf)) {
        Write-Host "ERROR: Missing required file: $f"
        $errors++
    }
}

# Check for forbidden files
$forbiddenPatterns = @("*.pem", "*.key", "*.crt")
foreach ($pattern in $forbiddenPatterns) {
    $matches = Get-ChildItem -Recurse -Filter $pattern -ErrorAction SilentlyContinue |
        Where-Object { $_.FullName -notmatch '\\.git\\' }
    if ($matches) {
        Write-Host "ERROR: Private key or certificate file found: $pattern"
        $errors++
    }
}

$dbFiles = Get-ChildItem -Recurse -Filter "*.sqlite" -Path "runtime" -ErrorAction SilentlyContinue
if ($dbFiles) {
    Write-Host "ERROR: SQLite database file found in runtime/ (runtime data, not source)"
    $errors++
}

if ($errors -eq 0) {
    Write-Host "=== validation passed ==="
} else {
    Write-Host "=== validation failed with $errors error(s) ==="
    if (-not $PassThru) { exit 1 }
}

if ($PassThru) { return $errors }
