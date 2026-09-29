# Security controls

## Threat model

Untrusted inputs include HTTP bodies, headers, webhook receivers, dependencies, issue text and agent tool output. Protect credentials, customer data and the integrity of durable state. Assume agents can make mistakes and external text can contain prompt injection.

## Enforced now

- Bounded JSON body, strict input validation and parameterized SQL.
- API key required for business routes; compare secrets without leaking them. Replace with product-specific authentication and authorization before production.
- Explicit server timeouts, graceful shutdown and no stack traces in HTTP responses.
- `.env` and credentials ignored by Git; secret scanning in local gates and CI.
- Locked dependencies and vulnerability audit in the `make security` gate.
- Architecture and coverage gates in hooks and CI.

## Before production

Use a secret manager, TLS ingress, per-principal authorization, request limits at the edge, audit logs, network restrictions and a formal key rotation plan. Give the database user minimum privileges. Pin and review GitHub Actions and enable branch protection, code scanning and repository secret scanning where available. Never expose credentials to an agent unless needed for a particular action; review external content as data.

No scanner proves absence of vulnerabilities. Triage findings with a dated decision record; do not silently suppress them. The percentage in the user's motivating text is a concern, not a measured fact for these repositories.
