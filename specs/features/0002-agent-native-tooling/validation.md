# Validation: 0002

Validated locally on 2026-10-02 with Go 1.27, Node.js 24.14 and PostgreSQL 16.6.

## Evidence

- `make check`: passed, including race tests, 95.6% core statement coverage against the 80% minimum, architecture, staged SDD and secret gates, `govulncheck`, Git history secret scan, PostgreSQL integration, MCP stdio tests and `npm audit` with no reported vulnerabilities.
- `npm test --prefix agents/mcp`: three tests passed; a real MCP client connected over stdio, listed tools, and traversal/symlink reads were refused.
- In a disposable copy, `go run ./tools/repoguard rename-module github.com/example/new-api` followed by `go test ./...`: passed.
- `opencode mcp list`: project server connected. Claude Code is not installed on this machine. Codex's project MCP config is loaded only after the checkout is trusted; its runtime connection was not exercised here.

## Remaining work for a new service

Run `make check` in the clone, review the new service's threat model, and enable branch protection and a separate human reviewer on its GitHub repository. The example MCP tools are intentionally read-only.
