# proxmox_virtual_environment_pools (Data Source)

Retrieves the list of Proxmox VE resource pools.

## Example Usage

### List all pools

```terraform
data "proxmox_virtual_environment_pools" "all" {}

output "pools" {
  value = data.proxmox_virtual_environment_pools.all.pool_ids
}
```

### Get specific pool info

```terraform
data "proxmox_virtual_environment_pools" "all" {}

locals {
  pools_map = zipmap(
    data.proxmox_virtual_environment_pools.all.pool_ids,
    data.proxmox_virtual_environment_pools.all.comments
  )
}

output "pool_details" {
  value = local.pools_map
}
```
