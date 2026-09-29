# Plan: 0001

1. Model Job, command and query with ports; unit test validation and idempotency behavior.
2. Add the PostgreSQL adapter, migration and transactional outbox; test database behavior.
3. Add HTTP adapter, key guard, health probes and worker; test request and failure paths.
4. Wire hooks and CI; record commands and results in `validation.md`.
