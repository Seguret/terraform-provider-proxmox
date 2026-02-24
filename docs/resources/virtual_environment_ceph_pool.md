---
page_title: "proxmox_virtual_environment_ceph_pool Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Ceph pool on a Proxmox VE node.
---

# proxmox_virtual_environment_ceph_pool (Resource)

Manages a Ceph pool on a Proxmox VE node.

## Example Usage

```terraform
resource "proxmox_virtual_environment_ceph_pool" "example" {
  node_name        = "pve"
  name             = "vm-storage"
  size             = 3
  min_size         = 2
  pg_num           = 128
  pg_autoscale_mode = "on"
  application      = "rbd"
}
```

## Schema

### Required

- `name` (String) The name of the Ceph pool. Changing this forces a new resource.
- `node_name` (String) The Proxmox node on which to manage the Ceph pool. Changing this forces a new resource.

### Optional

- `add_storages` (Boolean) Whether to also add a Proxmox storage definition for this pool on creation.
- `application` (String) Pool application (`rbd`, `cephfs`, `rgw`).
- `crush_rule` (String) The CRUSH rule to use for this pool.
- `min_size` (Number) Minimum number of replicas for I/O. Defaults to `2`.
- `pg_autoscale_mode` (String) PG autoscale mode (`on`, `off`, `warn`).
- `pg_num` (Number) Number of placement groups. Defaults to `128`.
- `size` (Number) Number of replicas. Defaults to `3`.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_ceph_pool.example pve/vm-storage
```
