# Documentation map

Start here, then open only what you need.

| Task | Read |
| --- | --- |
| Start a project | `README.md`, `docs/CONSTITUTION.md`, `docs/WORKFLOW.md` |
| Add an endpoint or behavior | `docs/ARCHITECTURE.md`, feature spec and nearby code/tests |
| Change database or events | `docs/OPERATIONS.md`, `docs/ARCHITECTURE.md`, migration and feature spec |
| Handle credentials or input | `docs/SECURITY.md`, feature spec |
| API contract | `docs/openapi.yaml`, HTTP tests and feature spec |
| Release or incident | `docs/OPERATIONS.md`, `docs/QUALITY_GATES.md` |
| Understand course influences | `docs/COURSE_DERIVATION.md` |

The map is stable; details live next to their owners. Update an existing page when a fact changes. Add a new page only when a new decision cannot fit its owner. Every feature spec links to relevant decisions and validation evidence. The executable gates, lockfile and migrations are the authoritative record of implemented behavior.
