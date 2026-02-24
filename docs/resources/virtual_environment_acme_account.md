---
page_title: "proxmox_virtual_environment_acme_account Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE ACME account registration (e.g. Let's Encrypt).
---

# proxmox_virtual_environment_acme_account (Resource)

Manages a Proxmox VE ACME account registration (e.g. Let's Encrypt).

## Example Usage

### Let's Encrypt Production Account

```terraform
resource "proxmox_virtual_environment_acme_account" "letsencrypt" {
  name    = "default"
  contact = "admin@example.com"
}
```

### Let's Encrypt Staging Account

```terraform
resource "proxmox_virtual_environment_acme_account" "letsencrypt_staging" {
  name      = "staging"
  contact   = "admin@example.com"
  directory = "https://acme-staging-v02.api.letsencrypt.org/directory"
}
```

### Custom ACME Server

```terraform
resource "proxmox_virtual_environment_acme_account" "custom_ca" {
  name      = "custom-ca"
  contact   = "admin@example.com"
  directory = "https://acme.custom-ca.com/directory"
}
```

### With TOS Acceptance

```terraform
resource "proxmox_virtual_environment_acme_account" "with_tos" {
  name    = "production"
  contact = "security@example.com"
  tos_url = "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf"
}
```


## Schema

### Required

- `contact` (String) Contact email address for the ACME account.

### Optional

- `directory` (String) ACME directory URL. Defaults to Let's Encrypt production.
- `name` (String) Account name (default: 'default').
- `tos_url` (String) URL of the ACME Terms of Service (must be accepted on account creation).

### Read-Only

- `account_url` (String) The registered ACME account URL (assigned by the CA).
- `id` (String) The ID of this resource.
