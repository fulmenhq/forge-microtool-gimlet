# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

## [0.2.0] - 2026-10-11

Earlier tags are not included in the initial commit.

### Changed

- Go directive is `1.26.9`. `gofulmen` is `v0.3.6` (`crucible` `v0.4.19`).
- CI and release use digest-pinned `goneat-tools-runner-glibc` `v0.5.8`. The goneat bootstrap pin is `v0.6.2`.
- Workflow actions are SHA-pinned. CI permissions are `contents: read`.
- Release archives are linux/amd64, linux/arm64, darwin/arm64, and windows/amd64. darwin/amd64 is not built.
- The GitHub release body is `docs/releases/v0.2.0.md`.
- The trademark notice in `LICENSE` names a North Carolina LLC.
