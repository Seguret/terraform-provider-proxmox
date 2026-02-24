---
page_title: "proxmox_virtual_environment_storage_nfs Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages an NFS storage in Proxmox VE.
---

# proxmox_virtual_environment_storage_nfs (Resource)

Manages an NFS network storage in Proxmox VE.

## Example Usage

```terraform
resource "proxmox_virtual_environment_storage_nfs" "example" {
  storage = "nfs-backup"
  server  = "nas.example.com"
  export  = "/export/proxmox"
  content = "backup,iso"
}
```

## Schema

### Required

- `export` (String) The NFS export path on the server. Changing this forces a new resource.
- `server` (String) The NFS server address. Changing this forces a new resource.
- `storage` (String) The storage identifier/name. Changing this forces a new resource.

### Optional

- `content` (String) Comma-separated list of content types (e.g. `images,rootdir,vztmpl,iso,backup,snippets`).
- `disable` (Boolean) Whether to disable this storage. Defaults to `false`.
- `nodes` (String) Comma-separated list of cluster nodes where this storage is accessible.
- `shared` (Boolean) Whether the storage is shared across nodes.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_storage_nfs.example nfs-backup
```
