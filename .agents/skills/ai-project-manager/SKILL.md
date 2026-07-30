---
name: ai-project-manager
description: Turn planning docs into implementation plans with approval checkpoints. Execute one reviewable phase at a time.
---

## ai-project-manager

### Workflow

1. **Read context** — Read `AGENTS.md`, `SPEC.md`, `ROADMAP.md`, `TASKS.md`, `docs/ARCHITECTURE.md`, `docs/STANDARDS.md`, and fontis-foundation `ENGINEERING_PRINCIPLES.md` and `AI_CONSTITUTION.md`.
2. **Assess current state** — Review what's done vs what's next. Check TASKS.md for in-progress and pending items.
3. **Produce a plan** — For the current phase, produce a requirement-linked plan:
   - What SPEC.md requirement does each task map to?
   - What's the implementation order (dependencies first)?
   - What automated validation exists or needs to be created?
   - What manual testing is required?
4. **Stop for approval** — Present the plan and pause. Do not implement until the user approves.
5. **Execute one phase** — Implement the smallest reviewable increment.
6. **Validate** — Run the full local gate. Update TASKS.md only after validation passes.
7. **Review** — Run `$pr-readiness` gate before opening a PR.

### Document templates

Project document templates are in `assets/project-docs/`. Use them for new projects or phases, not for routine updates to existing docs.
