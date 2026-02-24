---
page_title: "proxmox_virtual_environment_role Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE role.
---

# proxmox_virtual_environment_role (Resource)

Manages a Proxmox VE role.

## Example Usage

### VM Administrator Role

```terraform
resource "proxmox_virtual_environment_role" "vm_admin" {
  role_id    = "vm-admin"
  privileges = "VM.Allocate,VM.Audit,VM.Clone,VM.Console,VM.Config.CDROM,VM.Config.Cpu,VM.Config.Disk,VM.Config.HWType,VM.Config.Memory,VM.Config.Network,VM.Config.Options,VM.Monitor,VM.PowerMgmt,VM.Snapshot,VM.Snapshot.Rollback"
}
```

### Read-Only Viewer Role

```terraform
resource "proxmox_virtual_environment_role" "viewer" {
  role_id    = "viewer"
  privileges = "VM.Audit,Container.Audit,Datastore.Audit,Nodes.Audit,Cluster.Audit"
}
```

### Container Operator Role

```terraform
resource "proxmox_virtual_environment_role" "container_ops" {
  role_id    = "container-operator"
  privileges = "Container.Allocate,Container.Audit,Container.Console,Container.Config.Hostname,Container.Config.Memory,Container.Monitor,Container.PowerMgmt,Container.Snapshot"
}
```

### Storage Manager Role

```terraform
resource "proxmox_virtual_environment_role" "storage_manager" {
  role_id    = "storage-manager"
  privileges = "Datastore.Allocate,Datastore.AllocateSpace,Datastore.Audit,Datastore.Modify,Datastore.Read,Datastore.Report"
}
```

### Backup Operator Role

```terraform
resource "proxmox_virtual_environment_role" "backup_operator" {
  role_id    = "backup-operator"
  privileges = "VM.Audit,Container.Audit,Datastore.Audit,Datastore.Read,Datastore.AllocateSpace,Datastore.Modify"
}
```


## Schema

### Required

- `privileges` (String) Comma-separated list of privileges (e.g., 'VM.Audit,VM.Console').
- `role_id` (String) The role identifier.

### Read-Only

- `id` (String) The ID of this resource.
