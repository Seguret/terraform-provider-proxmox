---
page_title: "proxmox_virtual_environment_node_config Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages the configuration of a Proxmox VE node.
---

# proxmox_virtual_environment_node_config (Resource)

Manages the configuration of a Proxmox VE node.

## Example Usage

### Basic Node Configuration

```terraform
resource "proxmox_virtual_environment_node_config" "pve1" {
  node_name   = "pve1"
  description = "Primary Proxmox node"
}
```

### Node with Wake-on-LAN

```terraform
resource "proxmox_virtual_environment_node_config" "pve_with_wol" {
  node_name   = "pve"
  description = "Production Proxmox node with WoL support"
  wakeonlan   = "aa:bb:cc:dd:ee:ff"
}
```

### Node with Boot Delay

```terraform
resource "proxmox_virtual_environment_node_config" "pve_delayed_boot" {
  node_name                 = "pve"
  description               = "Node with delayed VM startup"
  startall_onboot_delay     = 120  # 2 minute delay before starting VMs
}
```

### Complete Node Configuration

```terraform
resource "proxmox_virtual_environment_node_config" "pve_complete" {
  node_name                 = "pve"
  description               = "Fully configured node - production"
  wakeonlan                 = "00:11:22:33:44:55"
  startall_onboot_delay     = 30
}
```


## Schema

### Required

- `node_name` (String) The node name.

### Optional

- `description` (String) A description for the node.
- `startall_onboot_delay` (Number) Initial delay in seconds, before starting all the Virtual Guests with on-boot enabled.
- `wakeonlan` (String) The Wake-on-LAN MAC address for this node.

### Read-Only

- `id` (String) The ID of this resource.
