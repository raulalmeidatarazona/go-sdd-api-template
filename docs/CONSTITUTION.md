# Constitution

## Mission

Build one small, dependable backend API at a time. Optimize for correct outcomes that humans can review. Agents accelerate implementation inside explicit boundaries; they do not decide product intent or silently accept risk.

## Technical choices

See `README.md` for the exact runtime and commands. The architecture is ports and adapters with separate commands and queries. Domain events are written to a durable outbox in the same database transaction as state changes. Delivery is at least once; consumers must be idempotent. HTTP, storage and event transport are adapters. Business rules remain independent of them.

## Roadmap

1. Replace the sample Job context with a real bounded context through a feature spec.
2. Add authentication and authorization appropriate to the product; the sample API key is only a bootstrap guard.
3. Add production integrations one at a time with explicit failure behavior and consumer contracts.
4. Define an SLO and run a restore, retry and duplicate-delivery exercise before production.

## Invariants

- No write without validation, an idempotency key and a transaction that records its event.
- No side effect may be acknowledged before the durable state and outbox commit.
- No external dependency in domain logic; no adapter imported by application logic.
- Failed dependencies make readiness fail; liveness stays independent of dependencies.
- Every code change has a linked feature or bug spec and repeatable validation.
- Never claim exactly once delivery or a numeric uptime guarantee from a repository template.

Amend this document in a dedicated PR with an ADR, the trade-off, migration plan and gate updates. Do not drift by accumulating exceptions.
