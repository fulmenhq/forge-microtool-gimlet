# forge-microtool-gimlet - AI Agent Guide

## Read First

1. Check `../AGENTS.forge-microtool-gimlet.local.md` (one level above the repo root) if present. It holds machine-specific instructions and lives outside the repository tree on purpose. A `.gitignore` entry is a convenience filter, not a security boundary. This repo declares conformance with the [3 Leaps OSS Sensitive Local Data Policy](https://github.com/3leaps/oss-policies/blob/main/SENSITIVE-LOCAL-DATA.md). The in-tree `AGENTS.local.md` gitignore entry is a backstop only. It is not the out-of-band pointer.
2. Read `MAINTAINERS.md` for contacts and governance.
3. Understand project scope before making changes. This repository is a template. Consumers clone it, remove its history, refit `.fulmen/app.yaml` and the module path, and launch their own tool.

## Operating Model

| Aspect   | Setting                                  |
| -------- | ---------------------------------------- |
| Mode     | Supervised (human reviews before commit) |
| Role     | Development assistant                    |
| Identity | Per session (no persistent memory)       |

`CODEX.md`, `CLAUDE.md`, and `.cline/rules/PROJECT.md` are optional adapters. They are not in this tree. Do not add them unless a checkout needs one. This file is the rules file.

## Quick Reference

| Task           | Command                      |
| -------------- | ---------------------------- |
| Quality checks | `make check-all`             |
| Tests          | `make test`                  |
| Build          | `make build`                 |
| Format         | `make fmt`                   |
| Doctor         | `make doctor`                |
| Identity check | `make validate-app-identity` |

## Session Protocol

### Before Changes

- Read relevant code first
- Understand the scope
- Keep changes minimal and focused
- Treat `internal/assets/**` as generated. Edit `.fulmen/app.yaml` and run `make sync-embedded-identity`

### Before Committing

- Run `make check-all`
- Verify tests pass
- Use the three trailers below
- `Committer-of-Record` is the last trailer

```
feat: add a doctor check for the identity file

Role: devlead
Co-authored-by: Example Model <noreply@3leaps.net>
Committer-of-Record: @3leapsdave
```

## DO / DO NOT

### DO

- Run quality gates before commits
- Read files before editing
- Keep changes focused
- Use `make` targets
- Use app identity for the application name, environment prefix, and config paths
- Use gofulmen for logging, signals, and exit codes
- Keep the pathfinder examples
- Document template behavior for CDRL users

### DO NOT

- Push without maintainer approval
- Skip quality gates
- Commit secrets
- Touch code outside task scope
- Commit `AGENTS.local.md` or files under `.plans/`
- Hardcode the application name
- Edit generated files under `internal/assets/**` except through the sync target

## Standards Reference

- Canonical: https://github.com/3leaps/crucible
- Fallback: clone https://github.com/3leaps/crucible as a sibling directory
- Commit style: https://github.com/3leaps/crucible/blob/main/docs/repository/commit-style.md
- Role catalog: https://github.com/3leaps/crucible/tree/main/config/agentic/roles
- Repo-local role notes: `config/agentic/roles/README.md`
