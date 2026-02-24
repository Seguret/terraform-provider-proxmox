---
page_title: "proxmox_virtual_environment_acl Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE access control list entry.
---

# proxmox_virtual_environment_acl (Resource)

Manages a Proxmox VE access control list entry.

## Example Usage

### User ACL on Specific VM

```terraform
resource "proxmox_virtual_environment_user" "operator" {
  user_id  = "operator@pve"
  comment  = "VM operator"
  enabled  = true
}

resource "proxmox_virtual_environment_role" "vm_operator" {
  role_id    = "vm-operator"
  privileges = "VM.Audit,VM.Console,VM.Monitor,VM.PowerMgmt"
}

resource "proxmox_virtual_environment_acl" "operator_vm_100" {
  user_id   = proxmox_virtual_environment_user.operator.user_id
  role_id   = proxmox_virtual_environment_role.vm_operator.role_id
  path      = "/vms/100"
  propagate = true
}
```

### Group ACL on Storage

```terraform
resource "proxmox_virtual_environment_group" "storage_admins" {
  group_id = "storage-admins"
  comment  = "Storage management team"
}

resource "proxmox_virtual_environment_acl" "storage_admins_local" {
  group_id  = proxmox_virtual_environment_group.storage_admins.group_id
  role_id   = "storage-manager"
  path      = "/storage/local-lvm"
  propagate = false
}
```

### Full Cluster Admin

```terraform
resource "proxmox_virtual_environment_user" "cluster_admin" {
  user_id  = "admin@pve"
  comment  = "Cluster administrator"
  enabled  = true
}

resource "proxmox_virtual_environment_acl" "cluster_admin" {
  user_id   = proxmox_virtual_environment_user.cluster_admin.user_id
  role_id   = "Administrator"
  path      = "/"
  propagate = true
}
```

### Department Access Hierarchy

```terraform
resource "proxmox_virtual_environment_acl" "finance_team_vms" {
  group_id  = "finance-team"
  role_id   = "vm-operator"
  path      = "/vms"
  propagate = true
}

resource "proxmox_virtual_environment_acl" "hr_team_storage" {
  group_id  = "hr-team"
  role_id   = "datastore-user"
  path      = "/storage"
  propagate = true
}
```


## Schema

### Required

- `path` (String) The access control path (e.g., '/', '/vms/100', '/storage/local').
- `role_id` (String) The role to assign.

### Optional

- `group_id` (String) The group to assign the role to. Mutually exclusive with user_id.
- `propagate` (Boolean) Whether to propagate the ACL to child objects.
- `user_id` (String) The user to assign the role to. Mutually exclusive with group_id.

### Read-Only

- `id` (String) The ID of this resource.
