---
page_title: "proxmox_virtual_environment_user Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE user.
---

# proxmox_virtual_environment_user (Resource)

Manages a Proxmox VE user.

## Example Usage

### Local User (PVE Realm)

```terraform
resource "proxmox_virtual_environment_user" "admin_user" {
  user_id    = "admin@pve"
  comment    = "Local Administrator"
  email      = "admin@example.com"
  first_name = "Admin"
  last_name  = "User"
  enabled    = true
  expire     = 0  # Never expires
}
```

### Automation User with Token

```terraform
resource "proxmox_virtual_environment_user" "automation" {
  user_id  = "terraform@pve"
  comment  = "Terraform automation account"
  email    = "terraform@example.com"
  enabled  = true
  expire   = 0
}

resource "proxmox_virtual_environment_user_token" "terraform_token" {
  user_id  = proxmox_virtual_environment_user.automation.user_id
  token_id = "terraform-token"
  comment  = "Token for Terraform provider"
  expire   = 0  # Never expires
}
```

### Users with Groups

```terraform
resource "proxmox_virtual_environment_group" "developers" {
  group_id = "developers"
  comment  = "Development team"
}

resource "proxmox_virtual_environment_user" "dev_user_1" {
  user_id    = "alice@pve"
  comment    = "Alice Developer"
  email      = "alice@example.com"
  first_name = "Alice"
  last_name  = "Developer"
  groups     = proxmox_virtual_environment_group.developers.group_id
  enabled    = true
}

resource "proxmox_virtual_environment_user" "dev_user_2" {
  user_id    = "bob@pve"
  comment    = "Bob Developer"
  email      = "bob@example.com"
  first_name = "Bob"
  last_name  = "Developer"
  groups     = proxmox_virtual_environment_group.developers.group_id
  enabled    = true
}
```

### User with Expiration

```terraform
resource "proxmox_virtual_environment_user" "contractor" {
  user_id    = "contractor@pve"
  comment    = "External contractor - temporary access"
  email      = "contractor@external.com"
  first_name = "Contractor"
  last_name  = "Name"
  enabled    = true
  expire     = 1735689600  # 2025-01-01 expiration
}
```



## Schema

### Required

- `user_id` (String) The user identifier (e.g., 'user@pam' or 'user@pve').

### Optional

- `comment` (String) A comment for the user.
- `email` (String) The user's email address.
- `enabled` (Boolean) Whether the user account is enabled.
- `expire` (Number) Account expiration date (UNIX epoch). 0 means no expiration.
- `first_name` (String) The user's first name.
- `groups` (String) Comma-separated list of groups.
- `keys` (String) Keys for two-factor authentication.
- `last_name` (String) The user's last name.
- `password` (String, Sensitive) The user password. Only used during creation for PVE realm users.

### Read-Only

- `id` (String) The ID of this resource.
