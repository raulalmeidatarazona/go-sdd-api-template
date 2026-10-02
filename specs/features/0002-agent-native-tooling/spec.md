# 0002 — Local agent capabilities and native gates

## Intent

Let a developer clone this Go template and use Codex, Claude Code, OpenCode or another AGENTS-compatible agent without inventing a repository layout. Keep repository checks in Go, and make the optional MCP server read-only.

## Acceptance criteria

1. One canonical skill source and two usable example skills are discoverable by Codex, Claude Code and OpenCode. Commands are usable where supported and documented elsewhere as prompt templates.
2. Each supported client can launch the same local MCP server after `npm ci --prefix agents/mcp`. It lists feature IDs, reads only approved feature documents, and refuses traversal or symlink escape.
3. The pre-commit, pre-push and CI checks use Go for architecture, SDD, secret and coverage policy. No Python runtime is required.
4. CI runs a real MCP stdio handshake test. Documentation explains setup and records versioned dependencies and trust boundary.

## Constraints and dependency rationale

The API has no MCP dependency. `agents/mcp` is a separate Node package because the supported clients use a local JavaScript stdio server. Its exact MCP SDK and Zod versions are locked. The client package is needed to verify the protocol in tests. Agent configuration is project-local and requires a trusted checkout; no tool may expose arbitrary file reads or shell execution. Review the nested lockfile and dependency audit with any update.
