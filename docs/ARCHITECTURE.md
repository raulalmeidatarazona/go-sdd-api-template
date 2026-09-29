# Architecture

The sample Job flow demonstrates one command (`CreateJob`), one query (`GetJob`), a domain event (`JobCreated`) and a transactional outbox. Do not add a bus, generic base class or broker merely to satisfy an architecture diagram. Introduce new abstraction when a second real use case needs it.

```text
HTTP request -> HTTP adapter -> application command/query -> port -> PostgreSQL adapter
                                                      \-> domain model
PostgreSQL transaction: Job row + outbox row
Outbox worker -> event publisher port -> external receiver (at least once)
```

Allowed direction: domain has no outward imports. Application may depend on domain. Adapters may depend on application and domain. Composition root wires adapters. `make arch` checks imports; changing this rule requires an ADR and corresponding gate change.

Commands and queries are separate operations and tests. A command owns transaction boundaries through a persistence port. Queries have no write side effects. Events are facts named in past tense and versioned in their payload. The outbox record includes a stable event ID; receivers deduplicate it. A failed publish remains pending for retry. A process crash after publish and before acknowledgement may deliver twice.

The example is intentionally modest: one bounded context and one relational store. CQRS does not imply two databases. Add projections only when a measured need exists. Avoid an in-process event handler for critical external effects because it cannot recover after a crash.
