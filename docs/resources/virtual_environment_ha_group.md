---
page_title: "proxmox_virtual_environment_ha_group Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE High Availability group.
---

# proxmox_virtual_environment_ha_group (Resource)

Manages a Proxmox VE High Availability group.

## Example Usage

### Basic HA Group

```terraform
resource "proxmox_virtual_environment_ha_group" "main_nodes" {
  group    = "main-nodes"
  nodes    = "pve1,pve2,pve3"
  comment  = "Main cluster HA group"
  restricted = false
}
```

### HA Group with Node Priority

```terraform
resource "proxmox_virtual_environment_ha_group" "priority_group" {
  group    = "priority-nodes"
  nodes    = "pve1:3,pve2:2,pve3:1"  # pve1 has highest priority
  comment  = "HA group with weighted priority"
  restricted = false
  no_failback = false
}
```

### Restricted HA Group

```terraform
resource "proxmox_virtual_environment_ha_group" "restricted_group" {
  group       = "restricted-services"
  nodes       = "pve1,pve2"
  comment     = "Restricted to high-performance nodes only"
  restricted  = true
  no_failback = false
}
```

### HA Group with Failback Disabled

```terraform
resource "proxmox_virtual_environment_ha_group" "failover_only" {
  group       = "failover-group"
  nodes       = "pve-primary:2,pve-backup:1"
  comment     = "Primary/backup without automatic failback"
  restricted  = true
  no_failback = true
}
```


## Schema

### Required

- `group` (String) The HA group name.
- `nodes` (String) Comma-separated list of nodes with optional priority (e.g., 'node1:2,node2:1').

### Optional

- `comment` (String) HA group description.
- `no_failback` (Boolean) Whether to disable automatic failback.
- `restricted` (Boolean) Whether the HA group is restricted to its members.

### Read-Only

- `id` (String) The ID of this resource.
