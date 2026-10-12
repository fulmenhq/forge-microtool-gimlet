# CI/CD Configuration

This document explains the CI/CD setup for this repository.

## Image is the toolchain

CI and release jobs run in `ghcr.io/fulmenhq/goneat-tools-runner-glibc:v0.5.8`, digest-pinned in the workflows. The glibc image is the default. A `-musl` variant exists for Alpine. The image supplies Go, goneat, and the formatters. Jobs do not install those tools.

The module `go` directive is `1.26.9`.

### Container permissions

Jobs set `options: --user root` so checkout can write the runner temp directory, then normalize `/__w/_temp` before later steps.

### What CI runs

`.github/workflows/ci.yml` has one `build-test` job: verify the image toolchain, verify the embedded identity mirror, download modules, format check, lint, test, build, and the standalone-binary check. Workflow permissions are `contents: read`.

The release workflow packages archives and opens a draft GitHub Release. The release body is `docs/releases/v0.2.0.md`. Generated release notes are off.

### Local development

```bash
make bootstrap   # goneat v0.6.2 via sfetch, plus hooks
make fmt-check
make lint
make test
```

## References

- [Fulmen toolbox containers](https://github.com/fulmenhq/fulmen-toolbox)
- [goneat documentation](https://github.com/fulmenhq/goneat)
- [GitHub Actions container jobs](https://docs.github.com/en/actions/using-jobs/running-jobs-in-a-container)
