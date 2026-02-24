---
page_title: "proxmox_virtual_environment_network_linux_bridge Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Linux bridge network interface on a Proxmox VE node.
---

# proxmox_virtual_environment_network_linux_bridge (Resource)

Manages a Linux bridge network interface on a Proxmox VE node.

## Example Usage

```terraform
resource "proxmox_virtual_environment_network_linux_bridge" "vmbr1" {
  node_name        = "pve"
  name             = "vmbr1"
  address          = "10.0.0.1/24"
  autostart        = true
  bridge_ports     = "ens19"
  bridge_vlan_aware = true
  comments         = "Internal VM network"
}
```

## Schema

### Required

- `name` (String) The bridge interface name (e.g. `vmbr0`). Changing this forces a new resource.
- `node_name` (String) The name of the Proxmox VE node. Changing this forces a new resource.

### Optional

- `address` (String) IPv4 address in CIDR notation (e.g. `192.168.1.1/24`).
- `address6` (String) IPv6 address in CIDR notation.
- `autostart` (Boolean) Whether to bring up the interface at boot. Defaults to `false`.
- `bridge_ports` (String) Space-separated list of ports to add to the bridge.
- `bridge_vlan_aware` (Boolean) Whether the bridge is VLAN aware. Defaults to `false`.
- `comments` (String) Comments for the bridge interface.
- `gateway` (String) IPv4 gateway address.
- `gateway6` (String) IPv6 gateway address.
- `mtu` (Number) The interface MTU.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_network_linux_bridge.example pve/vmbr1
```
