---
name: sdd-slice
description: Implement a backend behavior from a feature spec in one verifiable vertical slice. Use when adding or changing API behavior in this repository.
---

# SDD slice

Read `AGENTS.md` and `docs/INDEX.md`. Locate the feature's `spec.md`, `plan.md` and `validation.md`; if absent, create them from `specs/templates/feature/` before changing behavior. Confirm observable acceptance criteria, failure paths, idempotency, authorization and outbox effects where relevant.

Implement the smallest slice that meets one criterion. Keep domain rules independent of transport and persistence. Update the API contract and focused tests with the behavior. Record exact commands and results in `validation.md`; run `make check` before handoff. Report any failed gate and leave it visible.

For the stepwise prompt, see `agents/commands/sdd-start.md`. For architecture boundaries, read `docs/ARCHITECTURE.md` only when the change reaches those boundaries.
