---
page_title: "proxmox_virtual_environment_storage_directory Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a directory storage in Proxmox VE.
---

# proxmox_virtual_environment_storage_directory (Resource)

Manages a directory (filesystem path) storage in Proxmox VE.

## Example Usage

```terraform
resource "proxmox_virtual_environment_storage_directory" "example" {
  storage = "local-extra"
  path    = "/mnt/extra"
  content = "images,rootdir,iso"
}
```

## Schema

### Required

- `path` (String) The filesystem path for the directory storage. Changing this forces a new resource.
- `storage` (String) The storage identifier/name. Changing this forces a new resource.

### Optional

- `content` (String) Comma-separated list of content types (e.g. `images,rootdir,vztmpl,iso,backup,snippets`).
- `disable` (Boolean) Whether to disable this storage. Defaults to `false`.
- `nodes` (String) Comma-separated list of cluster nodes where this storage is accessible.
- `prune_backups` (String) Backup retention policy (e.g. `keep-last=3,keep-weekly=2`).
- `shared` (Boolean) Whether the storage is shared across nodes.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_storage_directory.example local-extra
```
