.PHONY: bootstrap fmt arch sdd secrets secrets-history deps test coverage security fast check run migrate integration agent-check

bootstrap:
	git config core.hooksPath .githooks

fmt:
	@test -z "$$(gofmt -l $$(find cmd internal tools -name '*.go'))" || (echo 'Run gofmt on Go files' && exit 1)

arch:
	go run ./tools/repoguard arch

sdd:
	go run ./tools/repoguard sdd --staged

secrets:
	go run ./tools/repoguard secrets

secrets-history:
	docker run --rm -v "$$(pwd):/repo" -w /repo zricethezav/gitleaks@sha256:cdbb7c955abce02001a9f6c9f602fb195b7fadc1e812065883f695d1eeaba854 git -v --no-banner --redact /repo

deps:
	go mod verify
	@test -z "$$(go mod tidy -diff)" || (echo 'go.mod/go.sum need tidy' && exit 1)

test:
	go test -race ./...

coverage:
	go test -coverprofile=coverage.out ./internal/domain ./internal/application ./internal/adapters/httpapi
	go run ./tools/repoguard coverage 80

security:
	go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...

agent-check:
	npm ci --prefix agents/mcp --ignore-scripts
	npm test --prefix agents/mcp
	npm audit --prefix agents/mcp --audit-level high

fast: fmt arch sdd secrets deps test agent-check

integration:
	@set -e; set -a; if test -f .env; then . ./.env; fi; set +a; \
	: "$${TEST_DATABASE_URL:?TEST_DATABASE_URL is required}"; \
	if test "$${TEST_DATABASE_URL}" = "$${DATABASE_URL:-}"; then echo 'Test database must be separate'; exit 1; fi; \
	DATABASE_URL="$${TEST_DATABASE_URL}" sh scripts/migrate.sh; \
	TEST_DATABASE_URL="$${TEST_DATABASE_URL}" go test -tags=integration ./internal/adapters/postgres; \
	TEST_DATABASE_URL="$${TEST_DATABASE_URL}" go test -tags=integration ./internal/adapters/events

check: fast coverage security secrets-history integration

run:
	@set -a; . ./.env; set +a; go run ./cmd/api

migrate:
	@set -a; . ./.env; set +a; sh scripts/migrate.sh
