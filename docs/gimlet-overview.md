---
title: "Gimlet Template Overview"
description: "Comprehensive overview of the Gimlet microtool template"
author: "infoarch"
date: "2025-11-10"
last_updated: "2025-11-10"
status: "active"
tags: ["overview", "template", "microtool", "go", "fulmen"]
---

# Gimlet Template Overview

## Purpose & Scope

Gimlet is a production-ready microtool template from the FulmenHQ ecosystem, providing foundational patterns for **single-purpose CLI tools** with minimal complexity. It serves as the canonical reference implementation of the Fulmen Forge Microtool Standard, demonstrating best practices for focused command-line utilities, safe filesystem operations, and operational excellence.

**Target Audience**: Go developers building single-purpose deployment tools, automation scripts, and utility commands that need reliability, structured logging, and graceful shutdown without the overhead of a full workhorse application.

**Design Philosophy**: **Ultra-narrow focus** with production-ready foundations. Microtools solve one problem well, staying under 2000 lines of business logic. When scope grows, promote to CLI category. Template provides standard quality gates, hooks, and observability—then gets out of your way.

## What Are Microtools?

Microtools are **single-purpose CLI tools** distinct from workhorse applications:

| Aspect        | Microtool (Gimlet)             | Workhorse (Groningen)         |
| ------------- | ------------------------------ | ----------------------------- |
| **Purpose**   | One narrow task                | Multi-purpose service         |
| **Interface** | CLI only (stdin/stdout/stderr) | HTTP endpoints + CLI          |
| **Scope**     | <2000 lines business logic     | Full application              |
| **Config**    | Simple (defaults + env vars)   | Enterprise three-layer        |
| **Examples**  | Fixture deployer, asset sync   | API server, background worker |
| **Category**  | `cli` (microtool profile)      | `workhorse`                   |

**When to use gimlet**: Building tools like "deploy fixtures to S3", "sync config files", "validate manifests"

**When to use groningen**: Building HTTP APIs, background workers, multi-command services

## Template Architecture

### Core Components

| Component           | Status            | Purpose                                                  | Key Features                      |
| ------------------- | ----------------- | -------------------------------------------------------- | --------------------------------- |
| **CLI Framework**   | ✅ Complete       | Cobra-based commands (version, envinfo, doctor, example) | REQUIRED standard commands        |
| **App Identity**    | ✅ CDRL-ready     | `.fulmen/app.yaml` for template customization            | Single source of truth            |
| **Logging**         | ✅ Structured     | Stderr logging with configurable levels                  | SIMPLE profile                    |
| **Signal Handling** | ✅ Cross-platform | Graceful shutdown, double-tap Ctrl+C                     | Context cancellation              |
| **Exit Codes**      | ✅ Standardized   | Foundry exit codes for automation                        | 54 semantic codes                 |
| **Pathfinder**      | ✅ REQUIRED       | Safe filesystem traversal                                | Path protection, pattern matching |

### Standard Commands (Per Microtool Standard)

Every microtool MUST implement:

| Command   | Purpose              | Key Features                                            |
| --------- | -------------------- | ------------------------------------------------------- |
| `version` | Display version info | `--extended` shows commit, build date, gofulmen version |
| `envinfo` | Show environment     | App identity, runtime info, config, env vars            |
| `doctor`  | System diagnostics   | Go version check, gofulmen access, config health        |
| `example` | Business logic demo  | Pathfinder usage for safe filesystem operations         |

**Customization**: Replace `example` command with your tool's single purpose (e.g., `deploy`, `sync`, `validate`)

### CLI Usage Patterns

```bash
# Help and version
gimlet --help                     # Show help
gimlet version                    # Show version
gimlet version --extended         # Show extended version info

# Diagnostics (REQUIRED commands)
gimlet doctor                     # Run environment diagnostics
gimlet envinfo                    # Show environment information

# Business logic (customize this)
gimlet example .                  # Find all files in current directory
gimlet example . --include "*.go" # Find all Go files
gimlet example . --exclude "vendor/" --max-depth 3  # Custom filters
```

### Environment Variables

All env vars use the `GIMLET_` prefix (customize via `.fulmen/app.yaml`):

```bash
GIMLET_LOG_LEVEL=debug            # Log level (trace|debug|info|warn|error)
GIMLET_CONFIG_PATH=./config.yaml  # Config file path override
```

See `.env.example` for complete list.

## Integration Stack

### Fulmen Ecosystem Integration

| Module                   | Integration       | Benefits                           |
| ------------------------ | ----------------- | ---------------------------------- |
| **gofulmen/appidentity** | ✅ Complete       | Dynamic app name, env var prefixes |
| **gofulmen/logging**     | ✅ SIMPLE profile | Structured stderr logging          |
| **gofulmen/signals**     | ✅ Cross-platform | Graceful shutdown, double-tap      |
| **gofulmen/foundry**     | ✅ Exit codes     | 54 standardized codes              |
| **gofulmen/pathfinder**  | ✅ REQUIRED       | Safe filesystem operations         |
| **gofulmen/config**      | ✅ Simple pattern | Defaults + env vars                |

### External Dependencies

| Dependency                 | Version | Purpose            |
| -------------------------- | ------- | ------------------ |
| **github.com/spf13/cobra** | v1.8.0  | CLI framework      |
| **go.uber.org/zap**        | v1.27.0 | Structured logging |
| **gopkg.in/yaml.v3**       | v3.0.1  | Config parsing     |

**Minimal dependencies**: Microtools avoid heavy dependencies to stay lightweight and focused.

## Configuration Management

### Simple Config Pattern

Microtools use **two-layer config** (NOT enterprise three-layer):

1. **Defaults**: Built-in sensible defaults
2. **Overrides**: Environment variables and optional config file

### Config File Discovery

```
~/.config/{vendor}/{config_name}.yaml
```

Example: `~/.config/fulmenhq/gimlet.yaml`

### Precedence

CLI flags → Env vars → Config file → Defaults

## Quality Gates & Development Workflow

### Code Quality Tools

**Goneat v0.6.2** provides unified quality tooling (installed via `sfetch` with hooks):

```bash
make bootstrap    # Install goneat v0.6.2 via sfetch + hooks
make fmt          # Format code (Go, YAML, Markdown)
make lint         # Lint checks (go vet + goneat assess: shfmt, actionlint, checkmake, yamlfmt lint)
make test         # Run tests
make check-all    # Run all quality checks
```

### Git Hooks

**Automated validation** via `.goneat/hooks.yaml` (rules configured in `.goneat/assess.yaml`):

- **Pre-commit**: Format + lint + date checks (fail-on: critical)
- **Pre-push**: Security + maturity + repo-status (fail-on: high)

Hooks run `goneat assess` directly (avoids Makefile recursion/loops):

```yaml
# .goneat/hooks.yaml
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

**Test hooks manually**:

```bash
goneat assess --hook pre-commit
goneat assess --hook pre-push

# Convenience aliases (optional)
make precommit
make prepush
```

**Regenerate hooks**:

```bash
goneat hooks generate    # Generate from .goneat/hooks.yaml
goneat hooks install     # Install to .git/hooks
```

### CDRL Validation Targets

**Template-specific quality gates**:

```bash
make doctor                  # CDRL completeness check
make validate-app-identity   # Detect hardcoded "gimlet" references
```

See `docs/development/fulmen_cdrl_guide.md` for CDRL workflow.

## CI/CD Integration

### GitHub Actions Workflow

`.github/workflows/ci.yml` runs on every push:

1. **Checkout** code
2. **Run** in `goneat-tools-runner-glibc` v0.5.8 (the image is the toolchain; Go 1.26.9)
3. **Verify** the bundled toolchain and the embedded identity mirror
4. **Download** modules
5. **Quality** checks (format, lint, test)
6. **Build** the binary and run the standalone-binary check

### Release Workflow

```bash
# Version bump
make version-bump-patch    # Increment patch version (goneat)
make version-bump-minor    # Increment minor version
make version-bump-major    # Increment major version

# Release preparation
make release-check         # Validate release readiness
make release-build         # Build multi-platform binaries
```

Binaries created in `bin/` with SHA256 checksums.

## Project Structure

```
forge-microtool-gimlet/
├── .fulmen/
│   └── app.yaml              # App identity (CDRL starting point)
├── .goneat/
│   └── hooks.yaml            # Git hooks configuration
├── cmd/
│   └── gimlet/
│       └── main.go           # Entry point
├── internal/
│   ├── cmd/                  # Cobra command implementations
│   │   ├── root.go           # Root command setup
│   │   ├── version.go        # Version command (REQUIRED)
│   │   ├── envinfo.go        # Environment info (REQUIRED)
│   │   ├── doctor.go         # Diagnostics (REQUIRED)
│   │   └── example.go        # Pathfinder demo (customize for your tool)
│   ├── config/               # Config types and loader
│   └── runtime/              # Runtime setup (logging, etc.)
├── docs/
│   ├── gimlet-overview.md    # This document
│   └── development/
│       ├── README.md         # Developer guides index
│       └── fulmen_cdrl_guide.md  # CDRL workflow
├── .env.example              # Environment variable template
├── .gitignore
├── Makefile                  # Standard targets
├── VERSION                   # Semantic version
├── go.mod
├── go.sum
└── README.md
```

## Customization After CDRL

1. **Replace example command** - Modify `internal/cmd/example.go` for your tool's purpose
   - Keep pathfinder usage (REQUIRED module for filesystem operations)
   - Update flags and logic for your use case
2. **Add commands** - Create new files in `internal/cmd/` following the pattern
3. **Keep doctor and envinfo** - These are REQUIRED by the microtool standard
4. **Extend config** - Update `internal/config/types.go` for tool-specific settings
5. **Add tests** - Follow `*_test.go` pattern
6. **Update docs** - Customize README with your tool's specific usage

See `docs/development/fulmen_cdrl_guide.md` for step-by-step instructions.

## Standards Compliance

This template implements:

- [Fulmen Forge Microtool Standard](https://github.com/fulmenhq/crucible/blob/main/docs/architecture/fulmen-forge-microtool-standard.md)
- [Fulmen Template CDRL Standard](https://github.com/fulmenhq/crucible/blob/main/docs/architecture/fulmen-template-cdrl-standard.md)
- [App Identity Module](https://github.com/fulmenhq/crucible/blob/main/docs/standards/library/modules/app-identity.md)
- [Pathfinder Extension](https://github.com/fulmenhq/crucible/blob/main/docs/standards/library/extensions/pathfinder.md)
- [Fulmen Makefile Standard](https://github.com/fulmenhq/crucible/blob/main/docs/standards/makefile-standard.md)

## Microtool vs Workhorse: When to Promote

**Promote to CLI category** if your microtool:

- Grows beyond single purpose
- Needs HTTP endpoints
- Requires enterprise three-layer config
- Exceeds 2000 lines of business logic
- Becomes multi-command Swiss Army knife

**Example**: If gimlet adds schema validation, code generation, AND docs building → promote to `cli` category

## Testing

```bash
make test                 # Run all tests
go test ./internal/...    # Test internal packages
go test ./cmd/...         # Test command wiring
```

Coverage target: 70% minimum

## Building

```bash
make build                # Build for current platform
make build-all            # Build for multiple platforms (Linux, macOS, Windows)
make package              # Package archives + checksums (dist/release)
SIGN=1 SIGNING_KEY_ID=<key-id> make package-sign  # Sign SHA256SUMS and export public key (requires gpg)
```

Binaries will be in `bin/` with SHA256 checksums.

## License

MIT

## Maintainers

See [MAINTAINERS.md](../MAINTAINERS.md) for governance structure and AI co-maintainer information.

## Additional Resources

- [CDRL Guide](development/fulmen_cdrl_guide.md) - Step-by-step customization
- [Microtool Standard](https://github.com/fulmenhq/crucible/blob/main/docs/architecture/fulmen-forge-microtool-standard.md) - Requirements and patterns
- [Groningen Overview](https://github.com/fulmenhq/forge-workhorse-groningen/blob/main/docs/groningen-overview.md) - Workhorse template comparison

---

**Document Status**: Active
**Last Updated**: 2025-11-10
**Maintained By**: Microtool Steward
