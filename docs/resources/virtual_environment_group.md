---
page_title: "proxmox_virtual_environment_group Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE group.
---

# proxmox_virtual_environment_group (Resource)

Manages a Proxmox VE group.

## Example Usage

### Basic Group

```terraform
resource "proxmox_virtual_environment_group" "developers" {
  group_id = "developers"
  comment  = "Development team members"
}
```

### Multiple Groups with Different Roles

```terraform
resource "proxmox_virtual_environment_group" "admins" {
  group_id = "admins"
  comment  = "System administrators"
}

resource "proxmox_virtual_environment_group" "operators" {
  group_id = "operators"
  comment  = "Infrastructure operators"
}

resource "proxmox_virtual_environment_group" "viewers" {
  group_id = "viewers"
  comment  = "Read-only access for monitoring"
}
```

### Groups for Department Access

```terraform
resource "proxmox_virtual_environment_group" "finance_team" {
  group_id = "finance-team"
  comment  = "Finance department VMs and resources"
}

resource "proxmox_virtual_environment_group" "hr_team" {
  group_id = "hr-team"
  comment  = "HR department resources"
}

resource "proxmox_virtual_environment_group" "it_team" {
  group_id = "it-team"
  comment  = "IT infrastructure team"
}
```


## Schema

### Required

- `group_id` (String) The group identifier.

### Optional

- `comment` (String) A comment for the group.

### Read-Only

- `id` (String) The ID of this resource.
- `members` (String) The group members (comma-separated user IDs). Read-only.
