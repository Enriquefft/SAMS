# Self-Hosted Automated Mail Server (SAMS)

SAMS provides infrastructure-as-code and automation for running a production-grade mail server on free-tier cloud resources. It uses Terraform to provision Oracle Cloud infrastructure and Ansible to deploy and manage [Mailcow](https://mailcow.email/). Passwords are stored in HashiCorp Vault and users access their mailboxes from Gmail via IMAP or POP3.

## Features

- Multi-domain support with automated onboarding and offboarding
- Infrastructure provisioning with Terraform modules in `infra/`
- Configuration management using Ansible roles in `automation/`
- Backup scripts and monitoring alerts
- Custom Go utilities in `tools/`
- CI workflows for linting, testing and release

## Getting Started

Clone the repository and initialize the development environment:

```bash
git clone https://github.com/Enriquefft/SAMS.git
cd SAMS
nix develop
```

Install Git hooks (optional):

```bash
lefthook install
```

Run the example program:

```bash
make build
./bin/sams
```

## Makefile Commands

- `make build` – compile the main binary to `bin/sams`
- `make test` – run unit tests
- `make lint` – run GolangCI-Lint
- `make vet` – type-check with go vet
- `make healthcheck` – build the Mailcow healthcheck CLI
- `make rotator` – build the password rotator CLI
- `make tools` – build all CLI tools

## Project Structure

```
├── cmd/            # Application entrypoints
├── internal/       # Private packages
├── pkg/            # Public packages (e.g. version info)
├── infra/          # Terraform modules
├── automation/     # Ansible roles & playbooks
├── docs/           # User and operator guides
├── tools/          # Go utilities
└── build/          # Release scripts
```

## Contributing

Before submitting a pull request:

```bash
make fmt
make lint
make vet
make test
```

Use focused commits with clear messages.

## License

See the [LICENSE](LICENSE) file for details.
