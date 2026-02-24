# proxmox_virtual_environment_groups (Data Source)

Retrieves the list of Proxmox VE groups.

## Example Usage

### List all groups

```terraform
data "proxmox_virtual_environment_groups" "all" {}

output "groups" {
  value = data.proxmox_virtual_environment_groups.all.group_ids
}
```

### Create a map of group IDs to comments

```terraform
data "proxmox_virtual_environment_groups" "all" {}

locals {
  groups = zipmap(
    data.proxmox_virtual_environment_groups.all.group_ids,
    data.proxmox_virtual_environment_groups.all.comments
  )
}

output "group_map" {
  value = local.groups
}
```
