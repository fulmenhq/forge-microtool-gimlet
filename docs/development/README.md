---
title: "Developer Guides Index"
description: "Entry point for Gimlet developer how-to guides"
last_updated: "2025-11-10"
---

# Developer Guides Index

This directory collects focused guides for gimlet microtool development. Start with the template overview for architecture and integration patterns, then use the CDRL guide for customization workflows.

## Available Guides

- [Gimlet Template Overview](../gimlet-overview.md) - Comprehensive overview of microtool architecture, components, quality gates, and hooks
- [CDRL Guide](fulmen_cdrl_guide.md) - Step-by-step template customization (Clone → Degit → Refit → Launch)

## Quick Reference

### Quality & Hooks

#### Generated Assets (Read-Only)

`internal/assets/**` is treated as generated/read-only. For embedded app identity, edit `.fulmen/app.yaml` then run `make sync-embedded-identity`.

Gimlet uses **goneat v0.6.2** (installed via `sfetch`) for unified code quality and automated git hooks:

**Configuration**: `.goneat/hooks.yaml`

```yaml
hooks:
  pre-commit:
    - command: "assess"
      args:
        [
          "--categories",
          "format,lint,security",
          "--fail-on",
          "critical",
          "--package-mode",
        ]
  pre-push:
    - command: "assess"
      args:
        [
          "--categories",
          "format,lint,security,dependencies,dates,tools,maturity,repo-status",
          "--fail-on",
          "high",
          "--package-mode",
        ]
```

**Testing hooks locally**:

```bash
goneat assess --hook pre-commit
goneat assess --hook pre-push

# Convenience aliases (optional)
make precommit
make prepush
```

**Regenerating hooks**:

```bash
goneat hooks generate    # Generate from .goneat/hooks.yaml
goneat hooks install     # Install to .git/hooks
```

### CI/CD

**GitHub Actions**: `.github/workflows/ci.yml`

- Runs on every push
- Executes: bootstrap → fmt → lint → test → build

**Makefile targets**:

```bash
make bootstrap           # Install goneat and dependencies
make check-all           # Run all quality checks (fmt, lint, test)
make build               # Build binary
make build-all           # Build multi-platform binaries
make doctor              # CDRL completeness check
make validate-app-identity  # Detect hardcoded references
```

### Code Quality Tools

**Format**:

```bash
make fmt                 # Format Go, YAML, Markdown (goneat)
```

**Lint**:

```bash
make lint                # Run go vet + goneat lint (includes actionlint/yamlfmt for workflows)
```

**Test**:

```bash
make test                # Run all tests (coverage report)
```

**SBOM** (Software Bill of Materials):

```bash
make dependencies        # Generate SBOM with goneat
```

**Security scanning**:

```bash
goneat assess --categories security    # Security vulnerability scan
```

### Version Management

**Using goneat**:

```bash
make version-bump-patch  # Increment patch version
make version-bump-minor  # Increment minor version
make version-bump-major  # Increment major version
```

**Manual**:

```bash
make version-set VERSION=1.2.3    # Set specific version
```

### CDRL Validation

**Template refit checks**:

```bash
make doctor                      # Comprehensive CDRL completeness
make validate-app-identity       # Find hardcoded "gimlet" references
```

See [CDRL Guide](fulmen_cdrl_guide.md) for full workflow.

## Standards References

External standards implemented by this template:

- [Fulmen Forge Microtool Standard](https://github.com/fulmenhq/crucible/blob/main/docs/architecture/fulmen-forge-microtool-standard.md) - Microtool requirements
- [Fulmen Template CDRL Standard](https://github.com/fulmenhq/crucible/blob/main/docs/architecture/fulmen-template-cdrl-standard.md) - CDRL workflow
- [Fulmen Makefile Standard](https://github.com/fulmenhq/crucible/blob/main/docs/standards/makefile-standard.md) - Required Makefile targets
- [App Identity Module](https://github.com/fulmenhq/crucible/blob/main/docs/standards/library/modules/app-identity.md) - `.fulmen/app.yaml` specification
- [Pathfinder Extension](https://github.com/fulmenhq/crucible/blob/main/docs/standards/library/extensions/pathfinder.md) - Safe filesystem operations

## Contributing

Follow the quality gates and hooks system for all contributions. See [MAINTAINERS.md](../../MAINTAINERS.md) for governance structure and AI co-maintainer information.

**Pre-commit checklist**:

- [ ] Code formatted (`make fmt`)
- [ ] Linting passes (`make lint`)
- [ ] Tests pass (`make test`)
- [ ] Hooks validate (`make precommit`)

**Pre-push checklist**:

- [ ] Security scan clean (`goneat assess --categories security`)
- [ ] CDRL validation passes (`make doctor`)
- [ ] Hooks validate (`make prepush`)

---

As new development patterns emerge, add focused guides here to keep this index as the single entry point for developer references.
