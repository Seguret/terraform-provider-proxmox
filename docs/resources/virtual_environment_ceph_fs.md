---
page_title: "proxmox_virtual_environment_ceph_fs Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a CephFS filesystem on a Proxmox VE node.
---

# proxmox_virtual_environment_ceph_fs (Resource)

Manages a CephFS filesystem on a Proxmox VE node.

## Example Usage

```terraform
resource "proxmox_virtual_environment_ceph_fs" "example" {
  node_name = "pve"
  name      = "cephfs"
  pg_num    = 32
}
```

## Schema

### Required

- `name` (String) The CephFS filesystem name. Changing this forces a new resource.
- `node_name` (String) The Proxmox node on which to create the CephFS. Changing this forces a new resource.

### Optional

- `pg_num` (Number) Number of placement groups for the filesystem pools. Changing this forces a new resource.

### Read-Only

- `data_pool` (String) The data pool name.
- `id` (String) The ID of this resource.
- `metadata_pool` (String) The metadata pool name.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_ceph_fs.example pve/cephfs
```
