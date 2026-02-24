---
page_title: "proxmox_virtual_environment_storage Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE storage definition.
---

# proxmox_virtual_environment_storage (Resource)

Manages a Proxmox VE storage definition.

## Example Usage

### Directory Storage

```terraform
resource "proxmox_virtual_environment_storage" "local_dir" {
  storage = "custom-storage"
  type    = "dir"
  path    = "/mnt/custom-storage"
  content = "images,iso,vztmpl,backup"
  enabled = true
}
```

### NFS Storage

```terraform
resource "proxmox_virtual_environment_storage" "nfs_backup" {
  storage = "nfs-backup"
  type    = "nfs"
  server  = "192.168.1.10"
  export  = "/export/proxmox-backups"
  content = "backup,vztmpl"
  enabled = true
  nodes   = "pve,pve2"
}
```

### LVM Storage

```terraform
resource "proxmox_virtual_environment_storage" "lvm_vms" {
  storage = "vm-storage"
  type    = "lvm"
  vgname  = "vg_vms"
  content = "images"
  enabled = true
}
```

### CIFS/SMB Storage

```terraform
resource "proxmox_virtual_environment_storage" "windows_share" {
  storage  = "windows-backup"
  type     = "cifs"
  server   = "windows-server.local"
  share    = "proxmox-backups"
  username = "proxmox-user"
  password = "SecurePassword123!"
  content  = "backup,iso"
  domain   = "WORKGROUP"
  enabled  = true
}
```

### ZFS Pool Storage

```terraform
resource "proxmox_virtual_environment_storage" "zfs_pool" {
  storage = "zfs-storage"
  type    = "zfspool"
  pool    = "tank"
  content = "images,rootdir"
  enabled = true
}
```

### PBS (Proxmox Backup Server) Storage

```terraform
resource "proxmox_virtual_environment_storage" "pbs_datastore" {
  storage     = "pbs-backups"
  type        = "pbs"
  server      = "backup-server.local"
  datastore   = "proxmox-backups"
  username    = "proxmox@pam"
  password    = "BackupPassword123!"
  fingerprint = "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0"
  content     = "backup"
  enabled     = true
}
```



## Schema

### Required

- `storage` (String) The storage identifier.
- `type` (String) The storage type (dir, lvm, lvmthin, zfspool, nfs, cifs, glusterfs, iscsi, iscsidirect, rbd, cephfs, pbs).

### Optional

- `content` (String) Comma-separated list of content types (images, rootdir, vztmpl, iso, backup, snippets).
- `datastore` (String) The PBS datastore name.
- `domain` (String) The domain (for CIFS).
- `enabled` (Boolean) Whether the storage is enabled.
- `export` (String) The NFS export path.
- `fingerprint` (String) The PBS server fingerprint.
- `namespace` (String) The PBS namespace.
- `nodes` (String) Comma-separated list of nodes where this storage is available.
- `password` (String, Sensitive) The password (for CIFS, PBS). Write-only.
- `path` (String) The filesystem path (for 'dir' type).
- `pool` (String) The ZFS/Ceph pool name.
- `prune_backups` (String) Backup retention policy.
- `server` (String) The server address (for NFS, CIFS, iSCSI, PBS).
- `share` (String) The CIFS share name.
- `shared` (Boolean) Whether the storage is shared across nodes.
- `username` (String) The username (for CIFS, PBS).
- `vgname` (String) The LVM volume group name.

### Read-Only

- `id` (String) The ID of this resource.
