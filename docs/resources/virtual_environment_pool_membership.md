---
page_title: "proxmox_virtual_environment_pool_membership Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages VM and container membership in a Proxmox VE resource pool.
---

# proxmox_virtual_environment_pool_membership (Resource)

Manages the membership of VMs and containers in a Proxmox VE resource pool.

## Example Usage

```terraform
resource "proxmox_virtual_environment_pool" "dev" {
  pool_id = "dev"
  comment = "Development pool"
}

resource "proxmox_virtual_environment_pool_membership" "dev" {
  pool_id    = proxmox_virtual_environment_pool.dev.pool_id
  vms        = [100, 101, 102]
  containers = [200]
}
```

## Schema

### Required

- `pool_id` (String) The pool identifier.

### Optional

- `containers` (List of Number) The container IDs to include in the pool.
- `vms` (List of Number) The VM IDs to include in the pool.

### Read-Only

- `id` (String) The ID of this resource.
