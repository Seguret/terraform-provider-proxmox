# proxmox_virtual_environment_roles (Data Source)

Retrieves the list of Proxmox VE roles.

## Example Usage

### List all roles

```terraform
data "proxmox_virtual_environment_roles" "all" {}

output "roles" {
  value = data.proxmox_virtual_environment_roles.all.role_ids
}
```

### Create a map of roles to their privileges

```terraform
data "proxmox_virtual_environment_roles" "all" {}

locals {
  roles = zipmap(
    data.proxmox_virtual_environment_roles.all.role_ids,
    data.proxmox_virtual_environment_roles.all.privileges
  )
}

output "role_privileges" {
  value = local.roles
}
```
