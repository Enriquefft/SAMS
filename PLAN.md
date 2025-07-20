## Summary

This document is the definitive, step-by-step implementation plan for the Self‑Hosted Automated Mail Server (SAMS) repository. It incorporates best practices, explicit security and compliance measures, clearly named Ansible roles (including Unattended Upgrades automation), and a split documentation structure to align precisely with the original requirements. Each component is described in dependency order, with validation steps baked in, ensuring a smooth end‑to‑end build, test, and release process.

---

## 1. Infrastructure Provisioning (`infra/`)

### 1.1 Initialize Root Module

* **Files:** `infra/main.tf`, `infra/variables.tf`, `infra/outputs.tf`
* Define the OCI provider and core network resources (VCN, subnets, security lists).
* Parameterize inputs: `region`, `compartment_ocid`, `ssh_public_key`, and free‑tier limits.
* Expose outputs: `vcn_id`, `subnet_ids`, `security_list_id`.
* **Validation:** `terraform fmt`, `terraform validate`, `terraform plan` against a free‑tier test account.

### 1.2 Develop `mail_server` Module

* **Module Path:** `infra/modules/mail_server/`
* **Resources:** compute instance(s), block storage volumes, network attachments.
* **Parameters:** `domain_names` (list), `instance_count`, `volume_size_gb`, `ssh_public_key`.
* **Outputs:** `instance_ids`, `volume_ids`, `public_ips`.
* Follow Terraform module best practices: input/output docs, tagging, version pinning.

### 1.3 Integration & Testing

* In root module, call the `mail_server` module for primary and optional secondary AZs.
* **Smoke Test:** Confirm VM reachability over SSH; ensure security lists allow SSH, HTTP, HTTPS.

---

## 2. Configuration Automation (`automation/`)

### 2.1 Inventory & Requirements

* **Files:** `automation/inventory.ini`, `automation/requirements.yml`.
* Define host groups: `mailservers`, `auxiliary`.
* Pin Galaxy roles (e.g., `geerlingguy.docker`, custom roles).

### 2.2 Ansible Roles

Each role lives under `automation/roles/` and is tested in isolation.

| Role Name             | Purpose                                                                |
| --------------------- | ---------------------------------------------------------------------- |
| `mailcow-install`     | Install Docker & Docker‑Compose; clone Mailcow at `mailcow_version`.   |
| `mailcow-config`      | Configure Mailcow settings: environment files, TLS certificates check. |
| `vault-integration`   | Retrieve Mailcow API key; store and rotate user passwords in Vault.    |
| `fail2ban`            | Install & configure Fail2Ban for SSH, Postfix, Dovecot.                |
| `unattended-upgrades` | Install and configure Ubuntu Unattended Upgrades with safe defaults.   |

### 2.3 Playbooks

* **`mailcow-users.yml`**: Loop through `mail_list.yml` for onboarding (create/mailbox & Vault KV).
* **`mailcow-offboard.yml`**: Remove mailboxes, revoke Vault tokens, delete KV entries.
* **`backup.yml`**: Invoke `backup/scripts/backup.sh` on primary node; report status.

### 2.4 Lint & Test

* Run `ansible-lint` on all roles and playbooks.
* Use Molecule or a disposable VM for role testing.

---

## 3. Secrets Management (`vault/`)

### 3.1 Vault Engine & Policies

* **README:** `vault/README.md` outlines KV path: `secret/data/mail/users/<DOMAIN>/`.
* **Policies:** `vault/policies/mail-policies.hcl` grants least-privilege access to Ansible.

### 3.2 Access Validation

* Use Vault CLI or API to test read/write for service accounts.

---

## 4. Documentation (`docs/`)

### 4.1 DNS Configuration Checklist

* **File:** `docs/dns_checklist_template.md`
* Placeholders for A, MX, SPF, DKIM (DNS record snippets).
* Step-by-step: add DNSSEC if desired, verify propagation.

### 4.2 Gmail Integration Guide

* **File:** `docs/gmail_setup_guide.md`
* Covers IMAP/POP3/SMTP settings, OAuth vs. App‑specific passwords, label/folder sync.

### 4.3 Vault Setup Guide

* **File:** `docs/vault_setup_guide.md`
* Walkthrough for initializing Vault, enabling KV engine, applying policies.

### 4.4 CI Rollback Procedure

* **File:** `docs/CI_rollback.md`
* Instructions to revert Terraform and Docker deployments in case of failure.

---

## 5. Monitoring & Alerting (`monitoring/`)

### 5.1 Prometheus Configuration

* **File:** `monitoring/prometheus/prometheus.yml`
* Scrape targets: Mailcow exporter, Node Exporter.

### 5.2 Alerting Rules

* **File:** `monitoring/prometheus/alerts/mail-alerts.yml`
* Alerts: queue length thresholds, disk usage above 80%, service down.

### 5.3 Optional Go Exporter

* **File:** `monitoring/exporters/mailcow_exporter.go`
* Scaffold in Go for custom Mailcow metrics (optional; use off‑the‑shelf if available).

---

## 6. Backup & Restore (`backup/`)

### 6.1 Backup Script

* **File:** `backup/scripts/backup.sh`
* Daily `rsync -a /opt/mailcow/maildata/ /backup/maildata/$(date +%F)`.
* Log rotation and error handling.

### 6.2 Restore Script

* **File:** `backup/scripts/restore.sh`
* Parameterized: select date snapshot, rsync back to `/opt/mailcow/maildata/`.

### 6.3 Retention & Lifecycle

* **Folder:** `backup/config/`
* Local retain 7 days; monthly sync to OCI Object Storage with lifecycle rules.

---

## 7. Security & Compliance (`security/`)

### 7.1 Fail2Ban Configuration

* Role `fail2ban` covers SSH, Postfix, Dovecot jails with safe defaults.
* Ensure logs are rotated and monitored.

### 7.2 Unattended Upgrades

* Role `unattended-upgrades` installs package, configures `/etc/apt/apt.conf.d/50unattended-upgrades`.
* Schedule via systemd timer; report via email on failure.

### 7.3 TLS & DNS Security

* Let’s Encrypt TLS automation in `mailcow-config` role.
* DNSSEC checklist in DNS documentation.

---

## 8. CI/CD & Version Control (`ci/`)

### 8.1 GitHub Actions Workflows

* **`terraform-plan.yml`**: PR-triggered plan with auto‑approval in dev branch.
* **`terraform-apply.yml`**: Manual approval step for prod environments.
* **`ansible-lint.yml`**: Lint roles/playbooks on push.

### 8.2 OCI Resource Manager

* **File:** `ci/oci-resource-manager.tf` (optional) for managed Terraform runs.

### 8.3 Drift Detection & Linting

* Nightly job: `tflint`, `terraform validate`, drift check via OCI CLI.

---

## 9. Custom Go Utilities (`tools/`)

### 9.1 Healthcheck CLI

* **Path:** `tools/go-mailcow-healthcheck/main.go`
* Poll Mailcow API endpoints; exit code non-zero on failures.

### 9.2 Password Rotator

* **Path:** `tools/go-password-rotator/main.go`
* CLI that rotates a single user or bulk from `mail_list.yml`, updates Vault and Mailcow.

### 9.3 Build & Release

* Go modules, Makefile targets: `make healthcheck`, `make rotator`, `make tools`.

---

## 10. Final Validation & Documentation

### 10.1 End-to-End Testing

* Deploy full stack in a test compartment; validate DNS, send/receive via Gmail.
* Automated smoke tests: send test email, verify inbox and logs.

### 10.2 Documentation Review

* Ensure `docs/*.md` updated with any parameter or workflow changes.

### 10.3 Version Tagging & Release

* Tag v1.0 in Git; update top‑level `README.md` with Quickstart snippet.

---

## Appendix: Repository Structure

```text
├── .gitignore
├── LICENSE
├── README.md
├── PLAN.md                        # This implementation plan
├── mail_list.yml                 # Central domains/users list
├── infra/                        # Terraform infra provisioning
│   ├── modules/mail_server/
│   ├── main.tf
│   ├── variables.tf
│   └── outputs.tf
├── automation/                   # Ansible roles & playbooks
│   ├── inventory.ini
│   ├── requirements.yml
│   ├── roles/
│   │   ├── mailcow-install/
│   │   ├── mailcow-config/
│   │   ├── vault-integration/
│   │   ├── fail2ban/
│   │   └── unattended-upgrades/
│   └── playbooks/
│       ├── mailcow-users.yml
│       ├── mailcow-offboard.yml
│       └── backup.yml
├── vault/                        # Vault setup & policies
│   ├── policies/mail-policies.hcl
│   └── README.md
├── docs/                         # User & operator documentation
│   ├── dns_checklist_template.md
│   ├── gmail_setup_guide.md
│   ├── vault_setup_guide.md
│   └── CI_rollback.md
├── monitoring/                   # Prometheus & Alertmanager
│   ├── prometheus/
│   │   ├── prometheus.yml
│   │   └── alerts/mail-alerts.yml
│   └── exporters/mailcow_exporter.go
├── backup/                       # Backup & restore
│   ├── scripts/{backup.sh,restore.sh}
│   └── config/                  # retention & lifecycle rules
├── security/                     # Security & compliance roles
│   └── (linked via automation/roles)
├── ci/                           # CI/CD workflows & configs
│   ├── .github/workflows/
│   └── oci-resource-manager.tf
└── tools/                        # Custom utilities in Go
    ├── go-mailcow-healthcheck/
    └── go-password-rotator/
```
