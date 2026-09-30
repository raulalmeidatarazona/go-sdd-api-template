# Go backend API template

An opinionated, small starting point for a Go service: spec-driven development, ports and adapters, explicit command/query separation, PostgreSQL transaction + outbox, strict quality gates and operational runbooks. The sample Job context is replaceable. Read [the map](docs/INDEX.md) and [API contract](docs/openapi.yaml) before expanding the template.

## Start a new service

Requires Go 1.27+, Python 3, PostgreSQL 16+ (Docker Compose is included), `psql` for remote migrations (local Docker includes it), and Git. Install `govulncheck` through the pinned `make security` command. No hidden global agent setup is required.

```sh
# First use GitHub's "Use this template" button to create YOUR_ORG/my-service.
git clone git@github.com:YOUR_ORG/my-service.git
cd my-service
python3 scripts/rename_module.py github.com/your-org/my-service
cp .env.example .env
make bootstrap
docker compose up -d db
make migrate
make check
make run
```

Replace the module name, sample API key and constitution before creating business features. `make bootstrap` configures versioned hooks in this clone. If you clone this source repository directly, create a new GitHub repository and change `origin` to its URL before pushing.

```sh
set -a; . ./.env; set +a
curl -sS -H "Authorization: Bearer $API_KEY" \
  -H 'Idempotency-Key: demo-1' -H 'Content-Type: application/json' \
  -d '{"name":"example"}' http://localhost:8080/v1/jobs
curl -sS http://localhost:8080/live
curl -sS http://localhost:8080/ready
```

The API reads settings from the environment; `make run` loads `.env` locally. In development, the outbox persists events even without a receiver. Set `EVENT_WEBHOOK_URL` and `EVENT_WEBHOOK_TOKEN` to start delivery. In production, both are required and the URL must use HTTPS. A receiver must deduplicate `X-Event-Id` and return 2xx only after processing durably. Never treat a webhook as guaranteed delivery without that contract.

## Work with an agent

1. Read `AGENTS.md` and `docs/INDEX.md`.
2. Create `specs/features/<id>/` from `specs/templates/feature/` and agree on `spec.md`.
3. Plan slices in `plan.md`, implement one slice, record actual checks in `validation.md`.
4. Run `make check`, review the diff and open a PR. See `docs/WORKFLOW.md` for ready-to-use prompts.

## Gates and repository setup

- `make fast`: format, architecture, SDD, secrets and focused tests.
- `make check`: full tests/coverage, isolated PostgreSQL integration, dependency integrity and vulnerability scan.
- `.githooks/pre-commit` runs the fast gate; `.githooks/pre-push` runs the full gate; CI repeats it.
- Protect `main` in each new GitHub repository: require `quality` status, one reviewer, no force push or direct push; enable Dependabot and secret scanning. Add a collaborator who can review PRs. Admin settings are not inherited by repositories created from a template and cannot be encoded by Git hooks alone.

See `docs/OPERATIONS.md` before any production deployment. The sample API key is a bootstrap guard, not a product authorization system.
