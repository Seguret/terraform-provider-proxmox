---
page_title: "proxmox_virtual_environment_ceph_osd Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Ceph OSD on a Proxmox VE node.
---

# proxmox_virtual_environment_ceph_osd (Resource)

Manages a Ceph Object Storage Daemon (OSD) on a Proxmox VE node. All configuration fields require replacement — the OSD must be destroyed and recreated if any of them change.

## Example Usage

```terraform
resource "proxmox_virtual_environment_ceph_osd" "example" {
  node_name = "pve"
  dev       = "/dev/sdb"
}

resource "proxmox_virtual_environment_ceph_osd" "encrypted" {
  node_name = "pve"
  dev       = "/dev/sdc"
  encrypted = true
  db_dev    = "/dev/nvme0n1p1"
  wal_dev   = "/dev/nvme0n1p2"
}
```

## Schema

### Required

- `dev` (String) The block device path (e.g. `/dev/sdb`). Changing this forces a new resource.
- `node_name` (String) The Proxmox node on which to manage the OSD. Changing this forces a new resource.

### Optional

- `db_dev` (String) Block device for the OSD DB. Changing this forces a new resource.
- `encrypted` (Boolean) Whether to encrypt the OSD. Changing this forces a new resource.
- `wal_dev` (String) Block device for the OSD WAL journal. Changing this forces a new resource.

### Read-Only

- `id` (String) The ID of this resource.
- `osd_id` (Number) The OSD ID assigned by Ceph.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_ceph_osd.example pve/0
```
