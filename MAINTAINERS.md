# Forge-Microtool-Gimlet – Maintainers

**Project**: forge-microtool-gimlet
**Purpose**: Fulmen microtool forge template for single-purpose CLI tools with production-ready foundations
**Governance Model**: 3leaps Initiative

## Human Maintainers

### @3leapsdave (Dave Thompson)

- **Role**: Project Lead & Primary Maintainer
- **Responsibilities**: Template architecture, microtool standard compliance, CDRL workflow design, gofulmen integration, quality assurance
- **GitHub**: [@3leapsdave](https://github.com/3leapsdave)
- **X**: [@3leapsdave](https://x.com/3leapsdave)
- **Supervision**: All AI agent contributions

## Org contacts

- **General**: hello@3leaps.net
- **Security**: security@3leaps.net
- **Legal**: legal@3leaps.net

## Agentic Roles

Role prompts are the [3leaps/crucible role catalog](https://github.com/3leaps/crucible/tree/main/config/agentic/roles). Repo-local overlay notes are in [config/agentic/roles/README.md](config/agentic/roles/README.md).

**Supervision**: All agentic roles operate under @3leapsdave.

## Agent Attribution

- `Role:` bare role slug
- `Co-authored-by:` model name `<noreply@3leaps.net>`
- `Committer-of-Record: @3leapsdave`

See [commit style](https://github.com/3leaps/crucible/blob/main/docs/repository/commit-style.md).

## Governance Structure

- Human maintainers approve architecture, releases, and supervise AI agents
- AI co-maintainers execute tasks, maintain template quality, and uphold microtool standards under supervision
- Template changes prioritize ease of CDRL customization for downstream consumers

## Communication Channels

- **Primary**: GitHub Issues and Pull Requests
- **Escalation**: @3leapsdave through GitHub

## Contribution Guidelines

All contributors (human and AI) must:

- Follow Fulmen Forge Microtool Standard
- Follow [Go coding standards](https://github.com/3leaps/crucible/blob/main/docs/coding/go.md)
- Run `make check-all` before commits
- Maintain all REQUIRED commands (version --extended, envinfo, doctor)
- Maintain all REQUIRED modules (Pathfinder)
- Document all template features for CDRL users
- Ensure template remains buildable and testable
- Keep CDRL workflow simple (Clone → Degit → Refit → Launch)

## CDRL Template Notes

This is a **template repository** - consumers will:

1. Clone this repo
2. Degit (remove git history)
3. Refit (customize `.fulmen/app.yaml`, module path, binary name)
4. Launch (build and run their microtool)

All contributions must make this workflow easier, not harder. When in doubt, prioritize simplicity and clarity over complexity.

## Template Consumers

Template consumers may choose to:

- Simplify or remove `MAINTAINERS.md` and `AGENTS.md` for their project
- Keep the agent framework if they plan AI-assisted development
- Customize agent roles for their team's needs

These files are provided as examples of the Fulmen agentic development pattern, not as strict requirements for downstream projects.
