# Agent Instructions

This repository provides a **Go** starter template with Docker, GoReleaser, GolangCI-Lint and Nix flake support. It is meant to be cloned and adapted for your own projects. Update names, documentation and configuration files once you copy this template.

## Workflow

1. **Format** Go files using `make fmt` which runs `gofmt`.
2. **Type-check** the code with `go vet ./...` and ensure it builds.
3. **Lint** with `make lint` (GolangCI-Lint).
4. **Run tests** with `make test` and include coverage. Use Arrange‑Act‑Assert patterns for non-trivial features and add integration tests when packages interact.
5. Keep commits small and focused with clear messages prefixed by `feat:`, `fix:`, `style:`, or `chore:`.
6. Pull request descriptions must summarize the changes, cite any updated files, and show test output. Note if commands fail because of environment issues.

## Repository Structure

- `cmd/` – application entry points
- `internal/` – private packages
- `pkg/` – public packages
- `configs/` – configuration templates
- `scripts/` – build and deploy helpers
- `docs/` – documentation
- `deployments/` – Docker Compose and other manifests
- `build/` – release and CI scripts
- `testdata/` – sample data for tests

## Code Style

- Four space indentation and LF line endings, as enforced by `.editorconfig`.
- Format code with `gofmt` via `make fmt`.
- Package names are short and lowercase. Exported types and functions use `PascalCase`.
- Keep functions small and single purpose; related code should appear close together.
- Ensure type safety with compile‑time checks and generics where appropriate.
- Provide robust error handling and input validation.
- New features must include unit tests and, when multiple components interact, integration tests.
- Remove unused code and avoid duplication.

## AGENTS.md Inheritance

These instructions apply to the entire repository. If another `AGENTS.md` exists in a subdirectory, it overrides these rules for files within that folder.
