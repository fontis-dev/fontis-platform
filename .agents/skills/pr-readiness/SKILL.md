---
name: pr-readiness
description: Validate PR merge readiness — local CodeRabbit review, CI status, independent review, manual testing evidence.
---

## pr-readiness

### Gate checklist

Before marking a PR ready for merge, verify each item:

1. **Local CodeRabbit review**
   ```bash
   coderabbit review --agent --uncommitted --include-untracked
   ```
   Fix all actionable findings. Document any intentionally skipped items with a reason.

2. **Full local gate**
   ```bash
   git diff --check && make fmt && make lint && make typecheck && make build && make test-unit && make test-integration && make security-scan
   ```
   Every step must pass.

3. **CI status**
   - All CI jobs pass on the latest commit.
   - Security scans (gitleaks, trivy, CodeQL) report zero new findings.

4. **Independent review**
   - At least one reviewer has approved the latest commit.
   - All review threads are resolved or have a documented response.

5. **Manual testing**
   - Record the test environment (OS, architecture, hardware/QEMU, runtime versions).
   - For HAL changes: tested on real target hardware.
   - For UI changes: screenshot included.
   - Describe what was tested, expected vs actual results.

6. **PR template sections**
   - Problem, Approach, Implementation, Impact on Contracts
   - AI Contribution (if applicable)
   - Verification checklist checked
   - Changelog entry
   - Follow-up / known limitations

### Commands

```bash
# Local code review
coderabbit review --agent --uncommitted --include-untracked

# Full validation gate
git diff --check && make fmt && make lint && make typecheck && make build && make test-unit && make test-integration && make security-scan
```
