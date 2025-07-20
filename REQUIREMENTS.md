# Requirements for Self‑Hosted Automated Mail Server

This document defines the full set of requirements to build and automate a generic, self-hosted mail server solution—called **Self‑Hosted Automated Mail Server**—capable of serving multiple domains for small startups. It leverages Oracle Cloud Always‑Free infrastructure, Mailcow (containerized mail suite), Terraform, Ansible, and HashiCorp Vault. End users access mailboxes (e.g., `user@<DOMAIN>`) via Gmail (IMAP/POP3).

---

## 1. Overview & Goals

- **Objective:** Provision and manage a production-grade, self-hosted mail service for arbitrary domains, with fully automated infrastructure, mail-server stack, user lifecycle, and integrations.
- **Key Requirements:**
  - Automate infra provisioning on Oracle Cloud Always‑Free (or equivalent) using Terraform.
  - Deploy Mailcow via Ansible for optimal developer experience.
  - Generate initial passwords stored in Vault; allow users to change passwords via webmail UI.
  - Automate user on‑boarding/off‑boarding from a central `mail_list.yml`.
  - Provide step‑by‑step Gmail setup documentation (IMAP & POP3) for end users.
  - Implement monitoring, alerting, and scheduled backups for mail data and configurations.

---

## 2. Assumptions & One‑Time Manual Steps

### 2.1. DNS Configuration (Manual)

For each domain `<DOMAIN>`, configure DNS records via your DNS provider:

```txt
# A record for mail host
mail.<DOMAIN>.          A   <STATIC_PUBLIC_IP>

# MX record
<DOMAIN>.               MX  10 mail.<DOMAIN>.

# SPF record
<DOMAIN>.               TXT "v=spf1 mx ip4:<STATIC_PUBLIC_IP> -all"

# DKIM (selector s1)
s1._domainkey.<DOMAIN>. TXT "v=DKIM1; k=rsa; p=<PUBLIC_DKIM_KEY>"

# DMARC policy
_dmarc.<DOMAIN>.        TXT "v=DMARC1; p=quarantine; rua=mailto:postmaster@<DOMAIN>"
```

> **Deliverable:** `dns_checklist_template.md` (Markdown) with placeholders for `<DOMAIN>`, `<STATIC_PUBLIC_IP>`, `<PUBLIC_DKIM_KEY>`, and support for dynamic sub-domain entries.

---

## 3. Infrastructure Provisioning

### 3.1. Cloud Provider & Free Tier

- **Oracle Cloud Always‑Free Example:**
  - 2× `VM.Standard.E2.1.Micro` (1 GB RAM, 1/8 OCPU each)
  - Ampere Flex pool: 3 000 OCPU‑h & 18 000 GB‑h/month
  - 2× 50 GB block volumes
  - Ingress ports: 25 (SMTP), 587 (SMTP), 143 (IMAP), 993 (IMAPS), 995 (POP3S)
  - Egress: all ports
- **Scaling:** To scale beyond Always‑Free limits, adjust module variables (`instance_shape`, `block_volume_sizes`) to larger VM shapes or add additional volumes; see §10 for CI/CD–driven scaling.

### 3.2. IaC Tooling

- **Terraform (v1.5+)** with the OCI provider.
- Optionally use **OCI Resource Manager** or Terraform Cloud for managed runs.

### 3.3. Terraform Resources

Define in `main.tf` (replace variables as needed):

```hcl
provider "oci" { region = var.region }

resource "oci_core_virtual_network" "vcn" { /* … */ }
resource "oci_core_subnet"           "subnet" { /* … */ }

resource "oci_core_security_list"    "mail_sec" {
  ingress_security_rules = [
    { protocol = "6"; source = "0.0.0.0/0"; tcp_options = { min = 25;  max = 25  } },
    { protocol = "6"; source = "0.0.0.0/0"; tcp_options = { min = 587; max = 587 } },
    { protocol = "6"; source = "0.0.0.0/0"; tcp_options = { min = 143; max = 143 } },
    { protocol = "6"; source = "0.0.0.0/0"; tcp_options = { min = 993; max = 993 } },
    { protocol = "6"; source = "0.0.0.0/0"; tcp_options = { min = 995; max = 995 } },
  ]
  egress_security_rules = [{ protocol = "all"; destination = "0.0.0.0/0" }]
}

resource "oci_core_instance" "mail_server" {
  metadata = { user_data = base64encode(file("cloud-init.sh")) }
  /* … */
}

resource "oci_core_volume"             "mail_data"   { size_in_gbs = 50  /* … */ }
resource "oci_core_volume_attachment"  "attach_data" { instance_id = oci_core_instance.mail_server.id; volume_id = oci_core_volume.mail_data.id }
```

> **Module Interface:**
>
> - **Inputs:**
>   - `domain_names` (list of strings)
>   - `instance_count` (number)
>   - `region`, `compartment_ocid`, `ssh_public_key`
> - **Outputs:**
>   - `vcn_id`, `subnet_ids`, `instance_ids`, `security_list_id`

---

## 4. Mail Server Stack Deployment

### 4.1. Mailcow (Containerized)

- **Why Mailcow:** Web UI, REST API, integrated password management, Docker‑Compose, active community support.
- **Components:** Postfix, Dovecot, Rspamd, MariaDB, Redis, Nginx, Roundcube.

### 4.2. Ansible Roles

- **Role: **``

  1. Install Docker & Docker‑Compose (version locked via `docker_compose_version` var) on Ubuntu 22.04.
  2. Clone Mailcow repo at tag `mailcow_version`.
  3. Configure `mailcow.conf`:
     ```ini
     MAILCOW_HOSTNAME=mail.<DOMAIN>
     MAILCOW_GENERATE_LETSENCRYPT_CERTS=1
     ```
  4. Create Docker volumes (`maildata`, `mailstate`, `mysql`, `redis`).
  5. Run `docker-compose pull && docker-compose up -d`.

- **Role: **``

  1. Wait for all containers to report “healthy.”
  2. Verify Let’s Encrypt certificate issuance/renewal.
  3. Retrieve and store Mailcow admin API key for automation.

> **Deliverable:** Idempotent, parameterized playbooks with version locking and upgrade-testing procedures.

---

## 5. Secrets Management

### 5.1. Vault Integration

- Use Vault KV engine at `secret/mail/users/<DOMAIN>/`.
- Ansible generates a unique, random 16‑char password per user and stores it at `secret/mail/users/<DOMAIN>/<user>`.
- **Audit Note:** A future playbook can fetch current Mailcow DB passwords and sync them back into Vault for complete auditability (TBD).

### 5.2. Password Self‑Service

1. Users visit `https://mail.<DOMAIN>/` and log in with initial credentials.
2. In webmail Settings, users change their password.
3. Vault retains only the initial credential for audit; current passwords live in Mailcow DB.

> **Deliverable:** Vault setup and policy guide; rotation best practices.

---

## 6. User Lifecycle Automation

### 6.1. Input File: `mail_list.yml`

```yaml
domains:
  - name: example.com
    users:
      - alice
      - bob
  - name: startup.io
    users:
      - carol
      - dave
```

### 6.2. Provisioning Playbook: `mailcow-users.yml`

- Loops over each domain & user:
  1. Generate & store password in Vault.
  2. Fetch password from Vault.
  3. Create mailbox via Mailcow API.

### 6.3. Offboarding Playbook: `mailcow-offboard.yml`

- Deletes mailbox via API and Vault secret metadata.

> **Deliverable:** Playbooks capable of multi-domain provisioning and offboarding.

---

## 7. Gmail Integration Documentation

Provide `gmail_setup_guide.md` with placeholders:

1. **IMAP (recommended):**
   - Incoming server: `mail.<DOMAIN>`
   - Port: 993 (SSL/TLS)
   - Username: `user@<DOMAIN>`
   - Password: initial or user-set credential
2. **POP3:**
   - Incoming server: `mail.<DOMAIN>`
   - Port: 995 (SSL/TLS)
   - Username/password as above
3. **Send mail as (SMTP):**
   - Server: `mail.<DOMAIN>`
   - Port: 587 (STARTTLS)
   - Authentication: required

> **Deliverable:** User-friendly guide covering both IMAP & POP3, with port/security-list updates.

---

## 8. Monitoring & Backup

### 8.1. Monitoring

- Mailcow exporter & Docker healthchecks.
- Prometheus Node Exporter on each VM.
- **Alert Rules:** Stored in `monitoring/alerts/*.yml`:
  - Mail queue length > 50 (per domain)
  - Disk usage > 80% on mail data volume
- **Notifications:** Alertmanager → Slack `#mailops` & email `ops@example.com`.

### 8.2. Backups

- **Daily Cron Job:**
  ```bash
  rsync -a /opt/mailcow/maildata/ /backup/mailcow/<DOMAIN>/$(date +"%Y%m%d")/
  ```
- **Config Backup:** `mailcow.conf`, `/etc/postfix`, `/etc/dovecot`.
- **Retention & Rotation:**
  - Local: 7 days, `gzip`-compressed, AES-256 encrypted
  - Monthly: copy to OCI Object Storage; lifecycle rule deletes objects > 365 days old
- **Restore Procedure:**
  - Scripts in `backup/restore.sh` support restoring specific snapshots.

> **Deliverable:** Backup scripts and restore documentation parameterized by domain.

---

## 9. Security & Compliance

- **TLS:** Let’s Encrypt via Mailcow.
- **Firewall:** OCI Security Lists per subnet.
- **Authentication:** Dovecot SASL via Mailcow DB.
- **Anti-spam/virus:** Rspamd included.
- **Hardening:** Fail2ban on SSH; Unattended Upgrades for OS.

> **Deliverable:** Security checklist template for all deployments.

---

## 10. CI/CD & Version Control

- **Repositories:**
  - `infra/` (Terraform modules)
  - `automation/` (Ansible roles & playbooks)
- **Branches & Triggers:**
  - `dev`: On PR, run `terraform plan`, `ansible-lint`, `yamllint`.
  - `main`: On merge, run `terraform plan` then `terraform apply` via OCI Resource Manager.
- **Drift Detection:**
  - Nightly job runs `terraform plan -detailed-exitcode` and alerts on drift.
- **Rollback Strategy:**
  - Revert merge commit → re-run `terraform apply` (documented in `CI/rollback.md`).
- **Linting & Testing:**
  - `tflint`, `terraform validate`
  - `ansible-lint`

> **Deliverable:** CI workflows with domain-parameterized pipelines and rollback docs.

---

## 11. Appendix

### A. Sample `mail_list.yml`

```yaml
domains:
  - name: example.com
    users: [alice, bob]
  - name: startup.io
    users: [carol, dave]
```

### B. Terraform Variables (`variables.tf`)

```hcl
variable "region"           { type = string }
variable "compartment_ocid" { type = string }
variable "ssh_public_key"   { type = string }
```

### C. Ansible Inventory (`inventory.ini`)

```ini
[mailcow]
mailcow ansible_host=<STATIC_PUBLIC_IP> ansible_user=opc ansible_ssh_private_key_file=~/.ssh/id_rsa
```

