# proxmox_virtual_environment_users (Data Source)

Retrieves the list of Proxmox VE users.

## Example Usage

### Get all users

```terraform
data "proxmox_virtual_environment_users" "all" {}

output "user_list" {
  value = data.proxmox_virtual_environment_users.all.user_ids
}

output "user_emails" {
  value = data.proxmox_virtual_environment_users.all.emails
}
```

### Find enabled users

```terraform
data "proxmox_virtual_environment_users" "all" {}

locals {
  enabled_users = [
    for i, user_id in data.proxmox_virtual_environment_users.all.user_ids :
    user_id
    if data.proxmox_virtual_environment_users.all.enabled[i]
  ]
}

output "active_users" {
  value = local.enabled_users
}
```
