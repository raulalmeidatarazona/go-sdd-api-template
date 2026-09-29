# Validation: 0001 — template baseline

Validated on 2026-09-29 with the pinned runtimes in this repository. This record documents the template build only; every clone must rerun `make check` and record its own results after changing behavior.

## Acceptance evidence

- POST with a new idempotency key returned 201; retry with the same key and payload returned 200 with the same Job ID; reusing the key for a different payload returned 409.
- `/live` returned 200 independently; `/ready` returned 503 before the database migration and 200 after it.
- PostgreSQL integration tests confirmed one business row and one outbox event in the same transaction, durable retry state, and readiness recovery after successful delivery.

## Commands actually run

- `make check`: passed. Core coverage **95.6%**, minimum 80%; race tests, architecture, database integration and event retry tests passed; `govulncheck` reported no reachable vulnerabilities after upgrading `pgx` to 5.9.2.
- `make migrate` and live HTTP smoke requests: passed.

## Manual review and remaining risk

The example webhook publisher was exercised through a fake receiver in integration tests. Run an end-to-end delivery, duplicate-event and power-loss drill with the real receiver before production. Replace the bootstrap API key with product-specific auth. Repository template settings and branch protection must be enabled on GitHub.
