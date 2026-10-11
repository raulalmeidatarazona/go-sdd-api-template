# Local agent capabilities

`AGENTS.md` is the shared repository contract. `.agents/skills` is the canonical skill directory; `agents/skills` is a discoverable alias. Each skill is short and points to the relevant specification or runbook only when needed. `agents/commands` contains reusable prompts. `agents/mcp` contains a read-only local MCP server for feature documents.

| Client | Skills | Prompts | MCP registration |
| --- | --- | --- | --- |
| Codex | `.agents/skills` | Paste `agents/commands/*.md` with `$ARGUMENTS` replaced; project custom slash commands are not supported | `.codex/config.toml` |
| Claude Code | `.claude/skills` links to the canonical skills | `/sdd-start <request>`, `/review-gate <scope>` | `.mcp.json` |
| OpenCode | `.agents/skills` | `/sdd-start <request>`, `/review-gate <scope>` via `.opencode/commands` links | `opencode.json` |
| Other agents | `AGENTS.md`, then `.agents/skills` | Read `agents/commands` as prompt templates | Any stdio MCP client may launch `node agents/mcp/mcp.js` |

## Setup and verify

The Go API and repository gates need Go. The optional MCP extension needs Node.js 24+ and its own pinned lockfile:

```sh
npm ci --prefix agents/mcp
make agent-check
```

`make check` and CI also run the MCP test. The server exposes only `list_features` and `read_feature_document`. It accepts a feature ID and one of `spec.md`, `plan.md`, `validation.md`; it cannot run commands, write files, read arbitrary paths or access environment secrets. It resolves the real file path before reading to reject symlink escapes. MCP clients should trust this project before launching its local server; Codex loads `.codex/config.toml` only for trusted projects. Each newly cloned repository must run `npm ci --prefix agents/mcp` once before connecting a client.

To start a feature, ask the agent to use `sdd-slice` or invoke `/sdd-start` where supported. The agent should agree on observable criteria, implement one slice, run `make check`, and write actual results in `validation.md`. Use `dependency-review` for update PRs.

## Format references

- [Codex skills](https://learn.chatgpt.com/docs/build-skills) and [project MCP](https://learn.chatgpt.com/docs/extend/mcp?surface=cli)
- [Claude skills](https://code.claude.com/docs/en/skills) and [project MCP](https://code.claude.com/docs/en/mcp)
- [OpenCode skills](https://opencode.ai/v2/docs/skills), [commands](https://opencode.ai/v2/docs/commands), and [MCP](https://opencode.ai/v2/docs/mcp-servers)
- [MCP TypeScript SDK](https://github.com/modelcontextprotocol/typescript-sdk)
