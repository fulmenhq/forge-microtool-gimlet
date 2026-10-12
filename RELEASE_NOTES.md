# Release Notes

This document tracks release notes for forge-microtool-gimlet releases.

Earlier tags are not included in the initial commit.

## [0.2.0] - 2026-10-11

### Baseline refresh

**Release Type**: Minor
**Status**: Ready

#### Overview

Toolchain and dependency refresh for the public template baseline. Earlier tags are not included in the initial commit.

#### Key Changes

- Go `1.26.9`. `gofulmen` `v0.3.6`, `crucible` `v0.4.19`.
- Digest-pinned `goneat-tools-runner-glibc` `v0.5.8`. goneat bootstrap `v0.6.2`.
- Workflow actions are SHA-pinned. CI permissions are `contents: read`.
- Published archives: linux/amd64, linux/arm64, darwin/arm64, windows/amd64. Intel Macs build from source: `GOOS=darwin GOARCH=amd64 go build ./cmd/gimlet`.
- No signed-tag ceremony and no committed public-key pin. Optional local `SHA256SUMS` signing remains. See the CDRL guide.
- The GitHub release body is this release's notes file.
