# CDRL Guide: Refitting forge-microtool-gimlet

This guide explains how to customize the forge-microtool-gimlet template for your own microtool using the CDRL (Clone → Degit → Refit → Launch) workflow.

## CDRL Workflow Overview

1. **Clone**: `git clone` the template repository
2. **Degit**: Remove git history to start fresh
3. **Refit**: Customize the template for your microtool
4. **Launch**: Run your customized microtool

## Step-by-Step Refitting

### Step 1: Clone the Template

```bash
git clone https://github.com/fulmenhq/forge-microtool-gimlet.git my-tool
cd my-tool
```

### Step 2: Degit (Remove Git History)

```bash
rm -rf .git
git init
```

### Step 2.1 (Optional): Purge Template Release Artifacts

### Step 2.1.1 (Optional): Template-Only Operational Docs

Depending on your organization’s governance and process, you may want to delete or replace some template operational docs.

- Consider delete/replace if you’re not using Fulmen governance:
  - `AGENTS.md`
  - `REPOSITORY_SAFETY_PROTOCOLS.md`
  - `RELEASE_CHECKLIST.md`
  - `TESTING.md`
- Keep but refit:
  - `MAINTAINERS.md`
  - `LICENSE`
  - `README.md`
- Keep if you keep the release process:
  - `.github/workflows/*`
  - `scripts/*`
  - `Makefile` (ensure `BINARY_NAME` is refit)

Template release files are maintainer-owned and are usually not relevant to your newly refitted repository.

```bash
rm -f CHANGELOG.md RELEASE_NOTES.md
rm -rf docs/releases
```

### Step 2.1.2 (Optional): Template Residue Scan

After refitting, run a residue scan to ensure template identifiers and env prefixes are fully removed.

```bash
rg -n "\\bgimlet\\b|GIMLET_" \
  --glob '*.{go,md,yaml,yml,json,mod,sum}' \
  --glob 'Makefile' .
```

If you see matches, update or delete the references (common hotspots include CLI help text, envinfo headers, and tests).

### Step 2.2 Reset Your Project Version

During CDRL, reset `VERSION` to your project’s starting point (commonly `0.1.0`).

```bash
echo "0.1.0" > VERSION
```

### Step 3: Refit (Customize for Your Microtool)

#### 3.1 Update App Identity (MOST IMPORTANT - DO THIS FIRST)

Edit `.fulmen/app.yaml` to match your microtool:

```yaml
app:
  vendor: myorg # Your organization
  binary_name: my-tool # Your microtool name
  env_prefix: MYTOOL_ # Your env var prefix (MUST end with _)
  config_name: my-tool # Your config file name
  description: "My awesome CLI tool for X"
  version: "0.1.0"

metadata:
  repository_category: cli # Keep this (microtools are a CLI profile)
```

**Why this matters**: The app identity file controls:

- Environment variable prefixes (`MYTOOL_*` instead of `GIMLET_*`)
- Configuration file paths (`~/.config/myorg/my-tool.yaml`)
- Binary name in CLI help and version output
- Logger service name in structured logs
- Error messages and diagnostics

#### 3.2 Sync Embedded Identity Mirror

This template embeds a mirrored copy of the app identity into the compiled binary so commands like `--help` and `version` work even when the binary is copied outside the repo.

`internal/assets/**` is treated as generated/read-only. Edit `.fulmen/app.yaml` (SSOT) and sync the embedded mirror.

After updating `.fulmen/app.yaml`, refresh the embedded mirror:

```bash
make sync-embedded-identity
```

#### 3.3 Validate App Identity

Run validation to find hardcoded "gimlet" references:

```bash
make validate-app-identity
```

This will scan your codebase and report any files that still contain hardcoded "gimlet" references that need updating.

**Common residue hotspots** (teams often miss these):

| File                      | Location              | What to Verify                                                  |
| ------------------------- | --------------------- | --------------------------------------------------------------- |
| `internal/cmd/envinfo.go` | Line ~24              | Header uses `identity.BinaryName` (already dynamic in template) |
| `internal/cmd/root.go`    | `Short`/`Long` fields | Root command description in help output                         |
| `internal/cmd/*.go`       | Various               | Any hardcoded "gimlet" in command descriptions                  |
| `internal/cmd/*_test.go`  | Test assertions       | References to `GIMLET_*` env vars                               |
| `README.md`               | Examples and usage    | Update tool name in all examples                                |
| `docs/`                   | Documentation         | Update references throughout                                    |
| `.env.example`            | Env var examples      | Update prefix (e.g., `MYTOOL_*`)                                |
| `scripts/`                | Shell scripts         | Update any hardcoded references                                 |
| `Makefile`                | Help/output strings   | Update `BINARY_NAME` and messages                               |

**Pattern reference** (envinfo.go already demonstrates the correct approach):

```go
// Correct (dynamic via identity)
loggerInstance.Info(fmt.Sprintf("=== %s Environment Information ===", identity.BinaryName))

// Avoid (hardcoded)
loggerInstance.Info("=== Gimlet Environment Information ===")
```

If `make validate-app-identity` is too narrow for your refit, do a full scan:

```bash
rg -n "gimlet|GIMLET_|groningen|GRONINGEN_" \
  --glob '*.{go,md,yaml,yml,json,mod,sum}' \
  --glob 'Makefile' .
```

#### 3.4 Update Module Path

Edit `go.mod` to change the module path:

```go
module github.com/myorg/my-tool

go 1.25.1

require (
    github.com/fulmenhq/gofulmen v0.1.28
    // ... rest of dependencies
)
```

#### 3.5 Update Import Paths

Update all imports in your Go files to use the new module path:

```bash
# Find and replace throughout codebase (macOS)
find . -name "*.go" -type f -exec sed -i '' 's|github.com/fulmenhq/forge-microtool-gimlet|github.com/myorg/my-tool|g' {} +

# Linux version
find . -name "*.go" -type f -exec sed -i 's|github.com/fulmenhq/forge-microtool-gimlet|github.com/myorg/my-tool|g' {} +
```

Or manually update imports in:

- `cmd/gimlet/main.go`
- `internal/cmd/*.go`
- `internal/config/*.go`
- `internal/runtime/*.go`

#### 3.6 Rename Command Directory

Rename the command directory to match your binary name:

```bash
mv cmd/gimlet cmd/my-tool
```

Update Makefile references:

```makefile
# Update BINARY_NAME
BINARY_NAME := my-tool
```

#### 3.7 Update Documentation

Update `README.md`:

- Change "gimlet" to your tool name
- Update description and purpose
- Update usage examples
- Update installation instructions

Update CLI command descriptions in `internal/cmd/*.go`:

```go
// Before
Short: "Example command demonstrating pathfinder usage"

// After
Short: "Process files matching pattern"
```

#### 3.8 Customize Business Logic

**Keep the Required Commands** (don't modify these):

- `version` - Version information (REQUIRED by Fulmen standard)
- `envinfo` - Environment diagnostics (REQUIRED)
- `doctor` - Health checks (REQUIRED)

**Replace the Example Command** with your microtool's purpose:

```bash
# Remove example command (or rename and customize)
# internal/cmd/example.go

# Add your business logic in internal/
mkdir -p internal/processor  # Or whatever makes sense for your tool
# Add your implementation here
```

**Pathfinder Integration**: Keep pathfinder usage if your tool traverses filesystems. The `example.go` shows the REQUIRED pattern - copy it when implementing file operations.

#### 3.9 Update Environment Variables

Copy `.env.example` to `.env` and customize:

```bash
# Your app will now use MYTOOL_* prefix
MYTOOL_LOG_LEVEL=info
MYTOOL_CONFIG_PATH=/custom/path/config.yaml
```

#### 3.10 Update Metadata Files

- `LICENSE`: Choose your license (MIT recommended)
- `MAINTAINERS.md`: Add your maintainers and agent roles
- `AGENTS.md`: Update agent instructions if using AI co-maintainers

### Step 4: Launch Your Microtool

```bash
# Install dependencies and tools
make bootstrap  # installs goneat v0.6.2 via sfetch + hooks

# Run tests to verify everything works
make test

# Build your microtool
make build
# Optional: package release artifacts (archives + checksums)
make package              # No signing
SIGN=1 SIGNING_KEY_ID=<key-id> make package-sign  # Sign SHA256SUMS and export public key (requires gpg + your signing key)

# Standalone binary check (REQUIRED before first release)
make test-standalone-binary

# Try it out
./bin/my-tool --help
./bin/my-tool version --extended
./bin/my-tool doctor
```

### Signing Notes (optional)

- `make package-sign` can detach-sign `SHA256SUMS` with your own GPG key. It is not a signed-tag release ceremony.
- This template does not ship tag-signing scripts or a committed public-key pin. See the section "Release signing and provenance (not included)" below.
- Verify exported public key file before publishing: `./scripts/verify-public-key.sh dist/release/<binary>-release-signing-key.asc`.

## Verification Checklist

After refitting, verify these work:

- [ ] App identity loads from `.fulmen/app.yaml`
- [ ] `make validate-app-identity` shows no unexpected references
- [ ] Environment variables use your prefix (`MYTOOL_*`)
- [ ] CLI shows your tool name in help output
- [ ] `version --extended` shows correct tool name and versions
- [ ] `envinfo` displays correct app identity
- [ ] `doctor` runs health checks successfully
- [ ] Config file path uses your vendor/config_name
- [ ] Logs show your service name in structured output
- [ ] Tests pass with `make test`
- [ ] All imports use your new module path

## Common Refitting Patterns

### File Processing Microtool

```yaml
vendor: mycompany
binary_name: file-processor
env_prefix: FILEPROC_
config_name: file-processor
description: "Process and transform files matching patterns"
```

### Data Validation Microtool

```yaml
vendor: datacorp
binary_name: schema-validator
env_prefix: VALIDATOR_
config_name: schema-validator
description: "Validate data files against schemas"
```

### Deployment Automation Microtool

```yaml
vendor: devops
binary_name: fixture-deployer
env_prefix: DEPLOY_
config_name: fixture-deployer
description: "Deploy test fixtures to environments"
```

## Troubleshooting

### App Identity Not Loading

```bash
# Check file exists
ls -la .fulmen/app.yaml

# Validate YAML syntax
cat .fulmen/app.yaml | grep -E "vendor|binary_name|env_prefix"

# Test loading
go run cmd/my-tool/main.go envinfo
```

### Environment Variables Not Working

```bash
# Check prefix matches app.yaml
grep env_prefix .fulmen/app.yaml

# Test with explicit variable
export MYTOOL_LOG_LEVEL=debug
./bin/my-tool envinfo
```

### Import Paths Not Updated

```bash
# Find remaining old imports
grep -r "forge-microtool-gimlet" --include="*.go" .

# Re-run replacement command
find . -name "*.go" -type f -exec sed -i '' 's|github.com/fulmenhq/forge-microtool-gimlet|github.com/myorg/my-tool|g' {} +

# Re-tidy modules
go mod tidy
```

### Build Failures

```bash
# Clean and rebuild
make clean
go mod tidy
make build

# Check for syntax errors
go vet ./...
```

## Optional Integrations

### Crucible Shim (CONDITIONAL)

The template includes crucible shim integration for version info. This is **REQUIRED** if you need to:

- Validate schemas from Crucible SSOT
- Access taxonomy files
- Read standards documentation
- Access other Crucible assets

If your microtool doesn't need SSOT assets, you can **keep the existing usage** (it's lightweight and already integrated for version info).

### Pathfinder (REQUIRED)

**Do NOT remove** pathfinder integration - it's REQUIRED by the Fulmen Microtool Standard for any filesystem operations. The `example.go` command demonstrates proper usage. Copy this pattern when implementing your file operations.

### Schema Validation (OPTIONAL)

The template includes a **working schema validation example** showing both meta-validation and data validation patterns. This is helpful for microtools that:

- Process YAML/JSON configuration files
- Validate data structures against schemas
- Need to ensure data quality before processing

**What's included:**

- `internal/cmd/validate.go` - Command demonstrating validation patterns
- `tests/fixtures/schema-validation/` - Example schemas and data files
  - `config.schema.json` - JSON Schema with oneOf, enums, pattern validation
  - `valid-config.yaml` / `valid-config.json` - Valid data examples
  - `invalid-config.yaml` - Invalid data for testing
  - `valid-bluegreen.yaml` - Blue-green deployment example

**Two validation patterns demonstrated:**

1. **Meta-validation**: Verify schema is valid JSON Schema (Draft-07 or 2020-12)

   ```bash
   gimlet validate --schema config.schema.json --meta-only
   ```

2. **Data validation**: Verify data conforms to schema (supports YAML and JSON)
   ```bash
   gimlet validate --schema config.schema.json --data config.yaml
   ```

**Try the examples:**

```bash
# Meta-validate the schema
./bin/gimlet validate \
  --schema tests/fixtures/schema-validation/config.schema.json \
  --meta-only

# Validate valid data
./bin/gimlet validate \
  --schema tests/fixtures/schema-validation/config.schema.json \
  --data tests/fixtures/schema-validation/valid-config.yaml

# Test with invalid data (will fail with helpful errors)
./bin/gimlet validate \
  --schema tests/fixtures/schema-validation/config.schema.json \
  --data tests/fixtures/schema-validation/invalid-config.yaml
```

**When to keep this integration:**

- ✅ Your tool validates configuration files
- ✅ Your tool processes structured data with schemas
- ✅ You need data quality checks
- ✅ Users will provide YAML/JSON input

**When to remove:**

- ❌ Your tool doesn't process structured data
- ❌ You don't need schema validation
- ❌ Your data is unstructured (logs, text files, etc.)

**To remove schema validation:**

1. Delete `internal/cmd/validate.go`
2. Remove `rootCmd.AddCommand(newValidateCmd())` from `internal/cmd/root.go`
3. Delete `tests/fixtures/schema-validation/` directory
4. Remove schema validation examples from documentation

**To customize for your use case:**

1. Replace example schema in `tests/fixtures/schema-validation/config.schema.json`
2. Update validation logic in `internal/cmd/validate.go` if needed
3. Add your own data examples
4. Update command help text and examples

**Implementation notes:**

- Uses `gofulmen/schema` module (JSON Schema Draft-07 and 2020-12 support)
- Supports both YAML and JSON data files
- Provides clear error messages for validation failures
- Uses foundry exit codes (`ExitDataValidationError`) for automation
- Logs to stderr, validation results to stdout

### Config Files (OPTIONAL)

If your microtool doesn't need user config files, you can simplify:

1. Keep the simple config pattern in `internal/config/loader.go`
2. Remove config file checking logic if not needed
3. Update documentation to reflect configuration approach

Most microtools benefit from at least basic config support (defaults + env vars).

## Advanced Customization

### Adding New Commands

1. Create `internal/cmd/mycommand.go`:

   ```go
   package cmd

   import "github.com/spf13/cobra"

   func newMyCommandCmd() *cobra.Command {
       cmd := &cobra.Command{
           Use:   "mycommand",
           Short: "Description of my command",
           RunE: func(cmd *cobra.Command, args []string) error {
               // Your logic here
               return nil
           },
       }
       return cmd
   }
   ```

2. Register in `internal/cmd/root.go`:

   ```go
   rootCmd.AddCommand(newMyCommandCmd())
   ```

3. Add tests in `internal/cmd/mycommand_test.go`

### Custom Flags

Add flags to your commands:

```go
cmd.Flags().StringVarP(&myVar, "my-flag", "m", "default", "flag description")
cmd.Flags().BoolVarP(&myBool, "verbose", "v", false, "verbose output")
```

Map to environment variables in config loader if needed.

## Best Practices

1. **Keep Single Purpose**: If scope grows, promote to `cli` category (not microtool)
2. **Maintain Pathfinder**: Always use pathfinder for filesystem traversal
3. **Test Thoroughly**: Run `make check-all` before commits
4. **Document Changes**: Update README with your tool's usage
5. **Version Semantically**: Use semver in `.fulmen/app.yaml`
6. **Security First**: Validate all inputs, sanitize paths, no hardcoded credentials
7. **Log to Stderr**: Keep stdout clean for data output
8. **Exit Codes**: Use foundry exit codes for automation-friendly behavior

## Release signing and provenance (not included)

This template ships an optional local checksum signer (`make package-sign`):
it can GPG-detach-sign `SHA256SUMS` and, when `SIGNING_KEY_ID` is set, export
a public key for `make verify-release-key`. That helper is not a signed-tag
release ceremony and not a verified provenance chain. There are no tag-signing
scripts, no committed public-key pin, and no GitHub "Verified" tag gate. Tags
created from this template are whatever your release process makes them.

If your project needs GPG-signed tags and signed release manifests, adopt the
pattern from a public 3 Leaps Go repository, for example
[3leaps/spanwit](https://github.com/3leaps/spanwit) (`RELEASE_CHECKLIST.md`,
`scripts/release-*.sh`, and the `release-*` Makefile targets).

When you add release environment variables, derive their prefix from your app
identity (`env_prefix` in `.fulmen/app.yaml`) rather than hard-coding the
template's prefix, so the variables follow your project through CDRL. Do not
commit key material, tagger identities, or fingerprints from another project.

## Support

- **Template Issues**: Report to [forge-microtool-gimlet](https://github.com/fulmenhq/forge-microtool-gimlet)
- **Fulmen Ecosystem**: See [FulmenHQ Documentation](https://fulmenhq.dev)
- **gofulmen Library**: See [gofulmen](https://github.com/fulmenhq/gofulmen)
- **Microtool Standard**: See Crucible docs for Fulmen Forge Microtool Standard

---

**Remember**: Microtools must stay focused on a single purpose. The app identity system ensures customization happens primarily in `.fulmen/app.yaml` rather than scattered throughout the codebase. Most refit work should be updating identity, renaming directories, and replacing business logic in `internal/`.
