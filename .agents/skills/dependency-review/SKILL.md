---
name: dependency-review
description: Review a dependency or GitHub Actions update against compatibility, integrity, security and quality gates before recommending merge.
---

# Dependency review

Read the proposed diff and release notes from the official maintainer. Verify the target version or action commit against its upstream source. Check lockfile changes and the actual API usage in this repository. Run `make check` on the candidate branch and record any failing command and cause.

Recommend merge only when the update is compatible and all required gates pass. Do not approve, merge, disable a gate or rewrite an upstream PR without the current user's authorization. A green check on an outdated base is evidence for that commit, not for the current target branch.
