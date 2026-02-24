---
page_title: "proxmox_virtual_environment_storage_cifs Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a CIFS/SMB storage in Proxmox VE.
---

# proxmox_virtual_environment_storage_cifs (Resource)

Manages a CIFS/SMB network storage in Proxmox VE.

## Example Usage

```terraform
resource "proxmox_virtual_environment_storage_cifs" "example" {
  storage  = "cifs-backup"
  server   = "nas.example.com"
  share    = "backups"
  username = "proxmox"
  content  = "backup,iso"
}
```

## Schema

### Required

- `server` (String) The CIFS/SMB server address. Changing this forces a new resource.
- `share` (String) The CIFS/SMB share name. Changing this forces a new resource.
- `storage` (String) The storage identifier/name. Changing this forces a new resource.

### Optional

- `content` (String) Comma-separated list of content types (e.g. `images,rootdir,vztmpl,iso,backup,snippets`).
- `disable` (Boolean) Whether to disable this storage. Defaults to `false`.
- `domain` (String) The domain for CIFS authentication.
- `nodes` (String) Comma-separated list of cluster nodes where this storage is accessible.
- `shared` (Boolean) Whether the storage is shared across nodes.
- `username` (String) The username for CIFS authentication.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_storage_cifs.example cifs-backup
```
