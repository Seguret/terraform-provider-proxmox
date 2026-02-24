---
page_title: "proxmox_virtual_environment_storage_lvmthin Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages an LVM-thin storage in Proxmox VE.
---

# proxmox_virtual_environment_storage_lvmthin (Resource)

Manages an LVM-thin pool storage in Proxmox VE.

## Example Usage

```terraform
resource "proxmox_virtual_environment_storage_lvmthin" "example" {
  storage  = "local-lvm"
  vgname   = "pve"
  thinpool = "data"
  content  = "images,rootdir"
}
```

## Schema

### Required

- `storage` (String) The storage identifier/name. Changing this forces a new resource.
- `thinpool` (String) The LVM thin pool name within the volume group. Changing this forces a new resource.
- `vgname` (String) The LVM volume group that contains the thin pool. Changing this forces a new resource.

### Optional

- `content` (String) Comma-separated list of content types (e.g. `images,rootdir`).
- `disable` (Boolean) Whether to disable this storage. Defaults to `false`.
- `nodes` (String) Comma-separated list of cluster nodes where this storage is accessible.
- `shared` (Boolean) Whether the storage is shared across nodes.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_storage_lvmthin.example local-lvm
```
