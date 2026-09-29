# Spec-driven workflow

## 0. Define the slice

Create a feature directory from `specs/templates/feature/`, or a bug record from `specs/templates/bug.md`. Give it a stable ID. Branch from the protected default branch. Keep the scope small enough to review in one sitting.

## 1. Research

Read the constitution, relevant decision records, interfaces, tests and official dependency documentation. List unknowns and distinguish observed facts from assumptions. Ask the product owner about missing behavior before coding. Do not paste large irrelevant context into the agent.

## 2. Specify

Write the actor, outcome, contract, invariants, acceptance examples, error responses, threat cases, idempotency, event contract and failure recovery in `spec.md`. State explicit non-goals. Approval of product intent happens here.

## 3. Plan

In `plan.md`, choose the smallest vertical slices, migrations, rollback, tests and operational changes. Record dependency rationale. Review database and security slices before execution.

## 4. Implement

Implement a slice, run focused tests, update docs and validation evidence, then continue. Keep commits reviewable. Any change to intent updates the spec first. Prefer ordinary CLI tools already in the repository over extra agent integrations.

## 5. Validate and release

Run `make check`, review the diff against acceptance criteria and inspect CI results. Exercise timeout, duplicate request, duplicate event and restart paths when relevant. A human reviews behavior and migration risk. Deploy with an explicit rollback and SLO. Replan after each feature: record a learning in an ADR or update the constitution in a separate change.

## Prompt examples

**Research:** “Read `AGENTS.md` and `docs/INDEX.md`, then the relevant links for feature `<id>`. Inspect the existing behavior. Produce only unknowns, risks and a proposed `spec.md`. Do not edit application code.”

**Plan:** “Using `specs/features/<id>/spec.md`, create a phased `plan.md`. For each phase list changed ports/adapters, tests, migrations, failure modes and rollback. Do not implement yet.”

**Implement:** “Implement phase 1 of `specs/features/<id>/plan.md`. Keep the public contract and architecture rules. Update `validation.md` with commands actually run and outcomes. Stop if a gate fails.”

**Review:** “Compare the diff to `spec.md` and `validation.md`. Focus on behavior, threat cases, retry/idempotency, migrations and unsupported claims. Report findings with file and line.”

These are starting prompts, not magic commands. A fresh agent conversation can reconstruct the state from versioned specs rather than chat history.
