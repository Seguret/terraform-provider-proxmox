---
page_title: "proxmox_virtual_environment_ha_resource Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE High Availability resource.
---

# proxmox_virtual_environment_ha_resource (Resource)

Manages a Proxmox VE High Availability resource.

## Example Usage

### HA-Enabled Virtual Machine

```terraform
resource "proxmox_virtual_environment_ha_resource" "vm_ha" {
  sid       = "vm:100"
  group     = "production"
  comment   = "HA-managed production web server"
  state     = "enabled"
  max_restart = 3
  max_relocate = 1
}
```

### HA-Enabled Container

```terraform
resource "proxmox_virtual_environment_ha_resource" "container_ha" {
  sid       = "ct:200"
  group     = "app-servers"
  comment   = "HA-managed application container"
  state     = "enabled"
  max_restart = 5
  max_relocate = 2
}
```

### Multiple HA Resources

```terraform
resource "proxmox_virtual_environment_ha_resource" "web_server" {
  sid       = "vm:100"
  group     = "production"
  comment   = "Primary web server"
  state     = "enabled"
  max_restart = 3
}

resource "proxmox_virtual_environment_ha_resource" "app_server" {
  sid       = "vm:101"
  group     = "production"
  comment   = "Application backend"
  state     = "enabled"
  max_restart = 3
}

resource "proxmox_virtual_environment_ha_resource" "database" {
  sid       = "vm:102"
  group     = "production-db"
  comment   = "Primary database server"
  state     = "enabled"
  max_restart = 2  # More conservative for database
  max_relocate = 1
}
```

### Disabled HA Resource

```terraform
resource "proxmox_virtual_environment_ha_resource" "maintenance" {
  sid     = "vm:103"
  group   = "maintenance"
  comment = "Under maintenance - HA disabled"
  state   = "disabled"
}
```

## Schema

### Required

- `sid` (String) The HA resource SID (e.g., 'vm:100' or 'ct:200').

### Optional

- `comment` (String) HA resource description.
- `group` (String) The HA group name.
- `max_relocate` (Number) Maximum number of relocation attempts.
- `max_restart` (Number) Maximum number of restart attempts.
- `state` (String) The desired state (started, stopped, enabled, disabled, ignored).
- `type` (String) The resource type (vm or ct).

### Read-Only

- `id` (String) The ID of this resource.
