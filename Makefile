.PHONY: bootstrap fmt arch sdd secrets secrets-history deps test coverage security fast check run migrate integration

bootstrap:
	git config core.hooksPath .githooks

fmt:
	@test -z "$$(gofmt -l $$(find cmd internal -name '*.go'))" || (echo 'Run gofmt on Go files' && exit 1)

arch:
	python3 scripts/archcheck.py

sdd:
	python3 scripts/check_sdd.py --staged

secrets:
	python3 scripts/check_secrets.py

secrets-history:
	docker run --rm -v "$$(pwd):/repo" -w /repo zricethezav/gitleaks@sha256:cdbb7c955abce02001a9f6c9f602fb195b7fadc1e812065883f695d1eeaba854 git -v --no-banner --redact /repo

deps:
	go mod verify
	@test -z "$$(go mod tidy -diff)" || (echo 'go.mod/go.sum need tidy' && exit 1)

test:
	go test -race ./...

coverage:
	go test -coverprofile=coverage.out ./internal/domain ./internal/application ./internal/adapters/httpapi
	python3 scripts/check_coverage.py 80

security:
	go run golang.org/x/vuln/cmd/govulncheck@v1.1.4 ./...

fast: fmt arch sdd secrets deps test

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
