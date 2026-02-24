---
page_title: "proxmox_virtual_environment_storage_pbs Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox Backup Server storage in Proxmox VE.
---

# proxmox_virtual_environment_storage_pbs (Resource)

Manages a Proxmox Backup Server (PBS) storage in Proxmox VE.

## Example Usage

```terraform
resource "proxmox_virtual_environment_storage_pbs" "example" {
  storage   = "pbs-backup"
  server    = "pbs.example.com"
  datastore = "main"
  username  = "backup@pbs"
}

resource "proxmox_virtual_environment_storage_pbs" "secure" {
  storage      = "pbs-secure"
  server       = "pbs.example.com"
  datastore    = "offsite"
  username     = "backup@pbs"
  namespace    = "proxmox"
  fingerprint  = "AA:BB:CC:..."
}
```

## Schema

### Required

- `datastore` (String) The PBS datastore name. Changing this forces a new resource.
- `server` (String) The Proxmox Backup Server address. Changing this forces a new resource.
- `storage` (String) The storage identifier/name. Changing this forces a new resource.
- `username` (String) The PBS username (e.g. `user@pbs`). Changing this forces a new resource.

### Optional

- `content` (String) Comma-separated list of content types (e.g. `backup`).
- `disable` (Boolean) Whether to disable this storage. Defaults to `false`.
- `fingerprint` (String) The PBS server TLS certificate fingerprint.
- `namespace` (String) The PBS namespace within the datastore.
- `nodes` (String) Comma-separated list of cluster nodes where this storage is accessible.
- `shared` (Boolean) Whether the storage is shared across nodes.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_storage_pbs.example pbs-backup
```
