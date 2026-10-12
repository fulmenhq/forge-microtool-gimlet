# forge-microtool-gimlet

Fulmen Microtools: Forge Precision CLI Instruments in Go to Thrive On Scale

Gimlet is the reference implementation for Fulmen Microtools, empowering developers to craft specialized, enterprise-grade CLI utilities that excel in focused automation tasks. Built on the Fulmen ecosystem's single-source-of-truth standards, it provides production-ready foundations for single-purpose tools that scale seamlessly.

**Template Category:** Microtool (G-instrument - Go implementation)
**Status:** Reference Implementation
**Standard:** [Fulmen Forge Microtool Standard](https://github.com/fulmenhq/crucible/blob/main/docs/architecture/fulmen-forge-microtool-standard.md)

## What is This?

This is a **working reference implementation** of a Fulmen microtool forge. It demonstrates:

- ✅ **App Identity Module** - Centralized identity in `.fulmen/app.yaml`
- ✅ **gofulmen Integration** - Logging, signals, config, error handling, exit codes
- ✅ **Pathfinder Integration** - Safe filesystem traversal (ready for extension)
- ✅ **CDRL Workflow** - Clone → Degit → Refit → Launch customization
- ✅ **Quality Gates** - `make validate-app-identity`, `make doctor`, `make check-all`
- ✅ **Production Patterns** - Graceful shutdown, structured logging, error envelopes

## Quick Start (CDRL Workflow)

### 1. Clone

```bash
git clone https://github.com/fulmenhq/forge-microtool-gimlet.git my-tool
cd my-tool
```

### 2. Degit

```bash
rm -rf .git
git init
```

### 3. Refit

**MOST IMPORTANT:** Edit `.fulmen/app.yaml` FIRST:

```yaml
vendor: myorg # Your organization
binary_name: my-tool # Your tool name
env_prefix: MYTOOL_ # Your env prefix
config_name: my-tool # Your config name
description: "My awesome tool"
```

Then validate and find hardcoded "gimlet" references:

```bash
make validate-app-identity
```

Update module path in `go.mod`:

```go
module github.com/myorg/my-tool
```

Rename command directory:

```bash
mv cmd/gimlet cmd/my-tool
```

Update imports in `cmd/my-tool/main.go`:

```go
import (
    "github.com/myorg/my-tool/internal/cmd"
    "github.com/myorg/my-tool/internal/config"
    "github.com/myorg/my-tool/internal/runtime"
)
```

### 4. Launch

```bash
make bootstrap  # installs goneat v0.6.2 via sfetch into user-space bin + hooks
go build ./cmd/my-tool
./my-tool --help
```

## Features

### Built-in Gofulmen Modules

- **App Identity** - Single source of truth for binary name, env prefix, vendor
- **Logging** - Structured stderr logging with configurable levels
- **Signals** - Graceful shutdown, double-tap Ctrl+C, context cancellation
- **Exit Codes** - Standardized exit codes from Foundry catalog
- **Pathfinder** - Safe filesystem traversal with pattern matching (REQUIRED for microtools)
- **Config** - Simple config loading with env var overrides

### Standard Commands (Per Microtool Standard)

- **`version --extended`** - Shows version, commit, build date, Go version, gofulmen version
- **`envinfo`** - Displays app identity, runtime info, config, environment variables
- **`doctor`** - Runs diagnostic checks (Go version, gofulmen access, config health)
- **`example`** - Demonstrates pathfinder usage for safe filesystem operations
- **`validate`** - Demonstrates schema validation (meta-validation + data validation)

### CLI Commands

```bash
# Help and version
gimlet --help                     # Show help
gimlet version                    # Show version
gimlet version --extended         # Show extended version info

# Diagnostics
gimlet doctor                     # Run environment diagnostics
gimlet envinfo                    # Show environment information

# Example (demonstrates pathfinder)
gimlet example .                  # Find all files in current directory
gimlet example . --include "*.go" # Find all Go files
gimlet example . --exclude "vendor/" --max-depth 3  # Custom filters

# Schema validation (demonstrates optional integration)
gimlet validate --schema config.schema.json --meta-only  # Validate schema
gimlet validate --schema config.schema.json --data config.yaml  # Validate data
```

### Environment Variables

All env vars use the `GIMLET_` prefix (customize via `.fulmen/app.yaml`):

```bash
GIMLET_LOG_LEVEL=debug            # Log level (trace|debug|info|warn|error)
GIMLET_CONFIG_PATH=./config.yaml  # Config file path override
```

See `.env.example` for a complete list of available environment variables.

## Development

### Note on Release Docs (Template-Owned)

`CHANGELOG.md`, `RELEASE_NOTES.md`, and `docs/releases/` are maintained for the **template itself**. During CDRL, delete or rewrite them for your own project.

### Prerequisites

- Go 1.26.9 (module target; the goneat-tools-runner image is the toolchain)
- make

### Setup

```bash
make bootstrap  # Install goneat v0.6.2 via sfetch and dependencies (installs hooks)
make test       # Run tests
make lint       # Run lint checks (go vet + goneat assess, incl. workflow linting)
make build      # Build binary
```

### Quality Gates

```bash
make validate-app-identity  # Check for hardcoded "gimlet" references
make doctor                 # Run health checks
make check-all              # Run all quality checks (fmt, lint, test)
```

### Project Structure

```
forge-microtool-gimlet/
├── .fulmen/
│   └── app.yaml           # App identity (CDRL starting point)
├── cmd/
│   └── gimlet/
│       └── main.go        # Entry point
├── internal/
│   ├── cmd/              # Cobra command implementations
│   │   ├── root.go       # Root command setup
│   │   ├── version.go    # Version command (REQUIRED)
│   │   ├── envinfo.go    # Environment info (REQUIRED)
│   │   ├── doctor.go     # Diagnostics (REQUIRED)
│   │   └── example.go    # Pathfinder demo (customize for your tool)
│   ├── config/           # Config types and loader
│   └── runtime/          # Runtime setup (logging, etc.)
├── docs/
│   └── development/
│       └── fulmen_cdrl_guide.md  # Detailed CDRL instructions
├── .env.example          # Environment variable template
├── Makefile              # Standard targets
├── go.mod
└── README.md
```

## Customization After CDRL

1. **Replace example command** - Modify `internal/cmd/example.go` for your tool's purpose
   - Keep pathfinder usage (REQUIRED module for filesystem operations)
   - Update flags and logic for your use case
2. **Add commands** - Create new files in `internal/cmd/` following the pattern
3. **Keep doctor and envinfo** - These are REQUIRED by the microtool standard
4. **Extend config** - Update `internal/config/types.go` for tool-specific settings
5. **Add tests** - Follow `*_test.go` pattern (see removed example for reference)
6. **Update docs** - Customize README with your tool's specific usage

## Testing

```bash
make test                 # Run all tests
go test ./internal/...    # Test internal packages
go test ./cmd/...         # Test command wiring
```

## Building

```bash
make build                # Build for current platform
make build-all            # linux/amd64, linux/arm64, darwin/arm64, windows/amd64
# Intel Macs: GOOS=darwin GOARCH=amd64 go build ./cmd/gimlet
make package              # Package archives + checksums (dist/release)
SIGN=1 SIGNING_KEY_ID=<key-id> make package-sign  # Sign SHA256SUMS and export public key (requires gpg)
```

Binaries will be in `bin/` with SHA256 checksums.

### Signing (optional)

- `make package-sign` can detach-sign `SHA256SUMS` with your own GPG key. It is not a signed-tag release ceremony.
- This template does not ship tag-signing scripts, a committed public-key pin, or a GitHub "Verified" tag gate.
- For a signed-tag pattern, see a public 3 Leaps Go repository such as [spanwit](https://github.com/3leaps/spanwit), and the note in the [CDRL guide](docs/development/fulmen_cdrl_guide.md).

## Standards Compliance

This template implements:

- [Fulmen Forge Microtool Standard](https://github.com/fulmenhq/crucible/blob/main/docs/architecture/fulmen-forge-microtool-standard.md)
- [Fulmen Template CDRL Standard](https://github.com/fulmenhq/crucible/blob/main/docs/architecture/fulmen-template-cdrl-standard.md)
- [App Identity Module](https://github.com/fulmenhq/crucible/blob/main/docs/standards/library/modules/app-identity.md)
- [Pathfinder Extension](https://github.com/fulmenhq/crucible/blob/main/docs/standards/library/extensions/pathfinder.md)

## Documentation

- [CDRL Guide](docs/development/fulmen_cdrl_guide.md) - Step-by-step customization
- [Microtool Standard](https://github.com/fulmenhq/crucible/blob/main/docs/architecture/fulmen-forge-microtool-standard.md) - Requirements and patterns

## Contributing

See [`AGENTS.md`](AGENTS.md) for the role model, commit-attribution format, and quality gates. This template conforms to the [3 Leaps OSS policies](https://github.com/3leaps/oss-policies): keep the repository free of proprietary or client-specific material, and keep machine-local or sensitive data outside the repository tree (e.g. `../AGENTS.<reposlug>.local.md`) rather than behind `.gitignore`.

## License

MIT, with a trademark notice in `LICENSE`.

## Maintainers

See [MAINTAINERS.md](MAINTAINERS.md)
