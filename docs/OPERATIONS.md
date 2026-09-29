# Operations and failure design

## Availability model

Set a measurable SLO for each real service; the template does not promise “nine nines.” Liveness means the process runs. Readiness requires database access and no outbox event older than the configured lag threshold. Traffic must stop on failed readiness. Use multiple replicas and a managed HA database when the SLO requires them.

## Power loss and retries

Write the business row and outbox row in one transaction. If power fails before commit, neither exists and the caller retries with the same idempotency key. If it fails after commit but before response, retry returns the existing result. If it fails after event delivery but before outbox acknowledgement, the worker sends again; consumers deduplicate by event ID. The worker backs off after failures and never discards an undelivered event automatically.

## Runbook

1. Check `/live` and `/ready` separately, then database health and oldest pending outbox age.
2. If the publisher is down, keep writes within an agreed backlog budget or fail readiness; do not delete pending rows.
3. Restore database backups to a separate environment, verify row and outbox consistency, then resume the worker.
4. For a poison event, preserve its ID and payload, investigate the receiver contract, and replay after correction. Never mark delivered to silence an alert.
5. For deploys, use backward-compatible expand/migrate/contract schema changes and a tested rollback.

## Production checklist

- Product-specific auth, rate limits and TLS configured.
- SLO, error budget, alert thresholds, logs, traces and backup retention owned.
- Restore drill, duplicate-delivery drill and DB outage drill passed.
- Database migration reviewed, backup taken and rollback rehearsed.
- CI required on protected branch; security alerts and dependency updates routed to an owner.
