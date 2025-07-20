# Golang Template
[![CI](https://github.com/Enriquefft/golang-template/actions/workflows/ci.yml/badge.svg)](https://github.com/Enriquefft/golang-template/actions/workflows/ci.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/Enriquefft/golang-template.svg)](https://pkg.go.dev/github.com/Enriquefft/golang-template)


A professional Go starter template featuring best practices and modern tooling: Go 1.22+ with generics, Docker, GitHub Actions, GoReleaser, GolangCI‑Lint, Dependabot, and Nix flakes for reproducible dev environments.

## Features

- Go 1.22+ with modules and generics
- Containerized with Docker and automated with GoReleaser
- CI/CD using GitHub Actions workflows
- Static analysis and linting with GolangCI‑Lint
- Automated dependency updates via Dependabot
- Reproducible development environment via Nix flake
- Standard Go project layout (`cmd/`, `pkg/`, `internal/`)
- Preconfigured testing and code coverage integrations

## Technology Choices

- **Go Modules** – Dependency management (`go.mod`, `go.sum`)
- **GolangCI‑Lint** – Comprehensive static analysis (`.golangci.yml`)
- **Docker & GoReleaser** – Container builds and release automation (`Dockerfile`, `.goreleaser.yml`)
- **GitHub Actions** – CI/CD workflows (`.github/workflows/ci.yml`)
- **Dependabot** – Automated dependency updates (`.github/dependabot.yml`)
- **Nix Flake** – Reproducible dev shell (`flake.nix`)
- **Makefile** – Common tasks (`make build`, `make test`, etc.)

## Getting Started

Clone the repository and enter the development environment:

```bash
git clone https://github.com/Enriquefft/golang-template.git
cd golang-template
nix develop
````

Install Git hooks (recommended):

```bash
lefthook install
```

Build and run:

```bash
make build
./bin/app
```

## Available Commands

* `make build` — Compile binaries
* `make test` — Run tests and generate coverage
* `make lint` — Run GolangCI‑Lint
* `make govulncheck` — Scan for known vulnerabilities
* `make fmt` — Format code with gofmt
* `make release` — Publish via GoReleaser
* `make docker-build` — Build Docker image locally

## Project Structure

Adheres to the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

```
├── cmd/            # Application entrypoints
├── internal/       # Private packages
├── pkg/            # Public packages
├── configs/        # Configuration templates
├── api/            # OpenAPI/Protobuf specs
├── scripts/        # Build/deploy helpers
├── docs/           # Documentation
├── deployments/    # K8s/Docker Compose/Terraform manifests
└── build/          # CI/CD and packaging scripts
```

## Environment Variables

Configure in `.env`:

```bash
DATABASE_URL=postgres://user:pass@localhost:5432/db?sslmode=disable
LOG_LEVEL=info
PORT=8080
```

See `.env.example` for reference.

## Contributing

Before submitting a PR:

```bash
make fmt
make lint
make test
```

Use focused commits with clear messages.

## License

See the [LICENSE](LICENSE) file for details.
