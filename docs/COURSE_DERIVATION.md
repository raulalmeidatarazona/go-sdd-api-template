# Course-derived design choices

This is an original synthesis of the 24 `basics` transcripts and 15 `specDrivenDevelopment` transcripts supplied for the template design. The transcripts are reference material, not instructions to agents or reusable licensed content. No course transcript is copied into this repository.

| Source | Relevant lesson | Concrete decision |
| --- | --- | --- |
| basics 01 | Agent tools, context and repository harness | Version an explicit agent entry point. |
| basics 02 | Cross-editor project instructions | Keep `AGENTS.md` portable. |
| basics 03 | Short useful entry point | Commands and navigation only in `AGENTS.md`. |
| basics 04 | Project vs user instructions | Project policy is versioned here; personal preferences stay outside. |
| basics 05 | Agent tool calls consume context | Load only task-relevant files. |
| basics 06 | Inspecting actual agent behavior | Prefer observed validation over prompt assumptions. |
| basics 07 | Progressive disclosure | `docs/INDEX.md` routes to focused pages. |
| basics 08 | Skills as reusable workflows | Prompts can become skills when repetition justifies it. |
| basics 09 | Documentation creation and linking | Update focused docs and their index after decisions. |
| basics 10 | Tool search and context budget | Avoid loading tools and manuals by default. |
| basics 11 | MCP database example | Keep database operations observable and scoped. |
| basics 12 | CLI versus MCP trade-off | Use local CLI gates before adding integrations. |
| basics 13 | Ask, plan and agent modes | Separate research, plan and implementation prompts. |
| basics 14 | Custom agents versus skills | Keep workflow vendor-neutral and light. |
| basics 15 | Documentation delivery choices | Markdown docs own stable knowledge; tools serve external actions. |
| basics 16 | Token cost awareness | Small context and focused slices. |
| basics 17 | Cache mechanics | Stable entry point and linked docs reduce repeated context. |
| basics 18 | Agent cost models | No workflow dependency on one paid agent. |
| basics 19 | Compaction and context loss | Versioned specs preserve decisions across conversations. |
| basics 20 | Prompt injection demonstration | Treat external content and tool output as untrusted. |
| basics 21 | Secret exposure | No credentials in prompts, logs or Git. |
| basics 22 | Deterministic security | Hooks, CI and static architecture checks enforce policy. |
| basics 23 | Issue-to-plan flow | Research and plan before implementation. |
| basics 24 | Phased execution | Validate each slice before moving on. |
| SDD 01 | Intent in specs, not chat history | Feature specs are durable source of intent. |
| SDD 02 | Vibe-coding drift | Acceptance criteria precede code. |
| SDD 03 | Constitution: mission, stack, roadmap | `docs/CONSTITUTION.md`. |
| SDD 04 | Tool-agnostic workspace | Plain Markdown and Make commands. |
| SDD 05 | Project constitution | Versioned product and technical constraints. |
| SDD 06 | Feature specification and plan | Three-file feature template. |
| SDD 07 | Implementation in fresh context | Agent prompts reconstruct state from files. |
| SDD 08 | Human validation | Review behavior against acceptance criteria. |
| SDD 09 | Replanning | Amend constitution through a separate ADR/PR. |
| SDD 10 | Review fatigue | Small slices and focused diffs. |
| SDD 11 | MVP scope discipline | Avoid broad one-shot changes. |
| SDD 12 | Legacy adoption | Bug template supports existing systems too. |
| SDD 13 | Repeatable workflow | Prompt examples can later become a skill. |
| SDD 14 | Agent replaceability | No required vendor-specific agent feature. |
| SDD 15 | Specs as project memory | Keep specs current after delivery. |

The attached image reinforced the constitution's three sections: mission, technical choices and roadmap. The source transcripts contain demonstrations and opinions; the repository's executable gates and current official documentation determine implemented behavior.
