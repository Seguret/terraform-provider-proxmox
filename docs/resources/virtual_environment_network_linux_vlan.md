---
page_title: "proxmox_virtual_environment_network_linux_vlan Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Linux VLAN network interface on a Proxmox VE node.
---

# proxmox_virtual_environment_network_linux_vlan (Resource)

Manages a Linux VLAN network interface on a Proxmox VE node.

## Example Usage

```terraform
resource "proxmox_virtual_environment_network_linux_vlan" "vlan100" {
  node_name       = "pve"
  name            = "ens18.100"
  vlan_id         = 100
  vlan_raw_device = "ens18"
  address         = "192.168.100.1/24"
  autostart       = true
}
```

## Schema

### Required

- `name` (String) The VLAN interface name (e.g. `ens18.100`). Changing this forces a new resource.
- `node_name` (String) The name of the Proxmox VE node. Changing this forces a new resource.

### Optional

- `address` (String) IPv4 address in CIDR notation (e.g. `192.168.100.1/24`).
- `autostart` (Boolean) Whether to bring up the interface at boot. Defaults to `false`.
- `comments` (String) Comments for the VLAN interface.
- `gateway` (String) IPv4 gateway address.
- `mtu` (Number) The interface MTU.
- `vlan_id` (Number) The VLAN tag (1-4094).
- `vlan_raw_device` (String) The parent interface for this VLAN.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_network_linux_vlan.example pve/ens18.100
```
