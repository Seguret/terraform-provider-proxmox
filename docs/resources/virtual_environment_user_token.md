---
page_title: "proxmox_virtual_environment_user_token Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE user API token.
---

# proxmox_virtual_environment_user_token (Resource)

Manages a Proxmox VE user API token.

## Example Usage

### API Token for Terraform Provider

```terraform
resource "proxmox_virtual_environment_user" "terraform" {
  user_id  = "terraform@pve"
  comment  = "Terraform automation account"
  email    = "terraform@example.com"
  enabled  = true
}

resource "proxmox_virtual_environment_user_token" "terraform_token" {
  user_id              = proxmox_virtual_environment_user.terraform.user_id
  token_id             = "terraform-token"
  comment              = "Token for Terraform provider"
  expire               = 0  # Never expires
  privileges_separation = true
}
```

### Token with Expiration

```terraform
resource "proxmox_virtual_environment_user_token" "ci_token" {
  user_id   = "ci-user@pve"
  token_id  = "ci-pipeline"
  comment   = "CI/CD pipeline token - expires in 1 year"
  expire    = 1735689600  # 2025-01-01
  privileges_separation = true
}
```

### Multiple Tokens for Same User

```terraform
resource "proxmox_virtual_environment_user" "automation" {
  user_id  = "automation@pve"
  comment  = "Automation service account"
  enabled  = true
}

# Token for regular operations
resource "proxmox_virtual_environment_user_token" "api_token" {
  user_id  = proxmox_virtual_environment_user.automation.user_id
  token_id = "api-operations"
  comment  = "Standard operations token"
  expire   = 0
}

# Token for backups (different scope)
resource "proxmox_virtual_environment_user_token" "backup_token" {
  user_id  = proxmox_virtual_environment_user.automation.user_id
  token_id = "backup-manager"
  comment  = "Backup management token"
  expire   = 0
}
```


## Schema

### Required

- `token_id` (String) The token name.
- `user_id` (String) The user ID (e.g., 'root@pam').

### Optional

- `comment` (String) A comment for the token.
- `expire` (Number) Token expiration date as a UNIX timestamp. 0 means no expiration.
- `privileges_separation` (Boolean) Whether the token has privilege separation (tokens cannot exceed user privileges when true).

### Read-Only

- `id` (String) The full token identifier (userid/tokenid).
- `value` (String, Sensitive) The token secret. Only available after creation; stored in state.
