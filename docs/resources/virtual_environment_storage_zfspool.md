---
page_title: "proxmox_virtual_environment_storage_zfspool Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a ZFS pool storage in Proxmox VE.
---

# proxmox_virtual_environment_storage_zfspool (Resource)

Manages a ZFS pool or dataset storage in Proxmox VE.

## Example Usage

```terraform
resource "proxmox_virtual_environment_storage_zfspool" "example" {
  storage = "zfs-vm"
  pool    = "rpool/data"
  content = "images,rootdir"
}

resource "proxmox_virtual_environment_storage_zfspool" "with_blocksize" {
  storage   = "zfs-large"
  pool      = "tank/vm"
  blocksize = "16k"
  content   = "images"
}
```

## Schema

### Required

- `pool` (String) The ZFS pool or dataset name (e.g. `rpool` or `data/vm-store`). Changing this forces a new resource.
- `storage` (String) The storage identifier/name. Changing this forces a new resource.

### Optional

- `blocksize` (String) The ZFS block size (e.g. `8k`, `16k`, `32k`).
- `content` (String) Comma-separated list of content types (e.g. `images,rootdir`).
- `disable` (Boolean) Whether to disable this storage. Defaults to `false`.
- `nodes` (String) Comma-separated list of cluster nodes where this storage is accessible.
- `shared` (Boolean) Whether the storage is shared across nodes.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_storage_zfspool.example zfs-vm
```
