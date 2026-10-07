# AGENTS.md — AI Agent Instructions for openshift/restic

## Project Overview
This is the OpenShift fork of [restic](https://restic.net), a fast, secure, and efficient backup program. Restic provides encrypted and deduplicated backups to various storage backends. This fork (`openshift/restic`) is maintained on the `oadp-dev` branch for use by the OADP (OpenShift API for Data Protection) operator for file-level backup and restore operations.

- **Primary Language**: Go
- **Module**: `github.com/restic/restic`
- **Default Branch**: `oadp-dev`

## Build Instructions
```bash
# Build the restic binary
make restic
# or equivalently:
make all

# Direct Go build
go build -o restic ./cmd/restic
```

## Test Instructions
```bash
# Run all tests
make test
# or directly:
go test ./...

# Run a specific package's tests
go test ./internal/backend/... -run TestName

# Vet code
go vet ./...
```

## Linting
```bash
# Run golangci-lint (configuration in .golangci.yml)
golangci-lint run ./...
```

Configuration: `.golangci.yml`

## Code Conventions
- Standard Go project layout
- Main binary in `cmd/restic/`
- Internal packages in `internal/` (backend, cache, crypto, repository, etc.)
- Error handling: use `errors.Wrap` / `errors.Errorf` for context
- Documentation in `doc/` directory
- Changelog entries in `changelog/`

## Project Structure
```
cmd/restic/    - Main restic binary and CLI commands
internal/      - Private packages
  backend/     - Storage backends (local, sftp, rest, s3, azure, gs, etc.)
  cache/       - Local cache management
  crypto/      - Encryption primitives
  repository/  - Repository format and operations
  restic/      - Core data types (snapshot, tree, node, etc.)
  walker/      - Tree walking utilities
doc/           - Documentation
docker/        - Docker build files
helpers/       - Helper scripts
changelog/     - Changelog entries for releases
```

## CI/CD
- GitHub Actions workflows in `.github/workflows/`:
  - `tests.yml` — Test suite
  - `bz-pr-action.yml` — Bugzilla PR integration
  - `pr-merge.yml` — PR merge automation
- Reproduce CI locally:
  ```bash
  make test
  golangci-lint run ./...
  ```

## Common Tasks

### Adding a new storage backend
1. Create a new package under `internal/backend/`
2. Implement the `backend.Backend` interface
3. Register in the backend factory
4. Add CLI flags in `cmd/restic/`

### Updating OADP-specific patches
- OADP patches live on the `oadp-dev` branch
- Keep patches minimal and rebasing-friendly against upstream `restic/restic`

### Working with the repository format
- Repository format code is in `internal/repository/`
- Encryption in `internal/crypto/`
- Pack file handling in `internal/repository/pack/`
