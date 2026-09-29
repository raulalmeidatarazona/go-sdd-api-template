# Quality gates

`make bootstrap` installs repository hooks with `core.hooksPath`. `pre-commit` checks formatting, architecture, specs, secrets and fast tests. `pre-push` runs the full `make check`. CI repeats the full gate because local hooks can be bypassed. If any required tool or scan is unavailable, the gate fails closed.

`make check` runs format verification, architecture, tests with a minimum coverage threshold, an isolated database integration test, dependency verification/audit, staged-file secret scanning and pinned Gitleaks history scanning (Docker required). Lockfiles must be committed. New dependency changes require a spec rationale and review. The gate reports what ran; it cannot prove a package is legitimate or that software is free of flaws.

Protect the default branch in GitHub with required CI status, pull request review and disallowed force pushes. GitHub settings are a separate administrative control; see repository setup in `README.md`. No individual local hook can enforce remote branch protection.
