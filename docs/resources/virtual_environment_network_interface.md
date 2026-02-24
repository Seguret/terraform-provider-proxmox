---
page_title: "proxmox_virtual_environment_network_interface Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a network interface on a Proxmox VE node.
---

# proxmox_virtual_environment_network_interface (Resource)

Manages a network interface on a Proxmox VE node.

## Example Usage

### Basic Bridge Interface

```terraform
resource "proxmox_virtual_environment_network_interface" "vmbr0" {
  node_name = "pve"
  iface     = "vmbr0"
  type      = "bridge"
  address   = "192.168.1.100"
  netmask   = "255.255.255.0"
  gateway   = "192.168.1.1"
  autostart = true
}
```

### VLAN Tagged Interface

```terraform
resource "proxmox_virtual_environment_network_interface" "vlan100" {
  node_name = "pve"
  iface     = "eth0.100"
  type      = "vlan"
  address   = "10.0.100.1"
  netmask   = "255.255.255.0"
  vlan_raw_device = "eth0"
  autostart = true
}
```

### Bonded Interface (Active-Backup)

```terraform
resource "proxmox_virtual_environment_network_interface" "bond0" {
  node_name     = "pve"
  iface         = "bond0"
  type          = "bond"
  bond_mode     = "active-backup"
  bond_miimon   = 100
  bond_slaves   = "eth1,eth2"
  address       = "192.168.1.50"
  netmask       = "255.255.255.0"
  gateway       = "192.168.1.1"
  autostart     = true
  apply_config  = true
}
```

### Managed Bridge with VLAN

```terraform
resource "proxmox_virtual_environment_network_interface" "vmbr_vlan" {
  node_name       = "pve"
  iface           = "vmbr100"
  type            = "bridge"
  bridge_ports    = "eth0"
  bridge_vlan_aware = true
  address         = "10.0.100.10"
  netmask         = "255.255.255.0"
  gateway         = "10.0.100.1"
  autostart       = true
}
```

### IPv6 Enabled Interface

```terraform
resource "proxmox_virtual_environment_network_interface" "ipv6_bridge" {
  node_name = "pve"
  iface     = "vmbr1"
  type      = "bridge"
  address   = "192.168.50.1"
  netmask   = "255.255.255.0"
  address6  = "2001:db8:50::1/64"
  autostart = true
}
```


## Schema

### Required

- `iface` (String) The interface name (e.g., 'vmbr0', 'bond0', 'eth0.100').
- `node_name` (String) The node name.
- `type` (String) The interface type (bridge, bond, eth, vlan, OVSBridge, OVSBond, OVSPort, OVSIntPort).

### Optional

- `address` (String) IPv4 address.
- `address6` (String) IPv6 address.
- `apply_config` (Boolean) Whether to apply the network configuration immediately after changes.
- `autostart` (Boolean) Whether to bring up the interface at boot.
- `bond_mode` (String) Bond mode (balance-rr, active-backup, balance-xor, etc.).
- `bond_primary` (String) Primary bond interface.
- `bond_xmit_hash_policy` (String) Bond transmit hash policy.
- `bridge_fd` (Number) Bridge forward delay.
- `bridge_ports` (String) Bridge ports (space-separated interface names).
- `bridge_stp` (String) STP mode (on or off).
- `bridge_vlan_aware` (Boolean) Whether to enable VLAN-aware bridge.
- `cidr` (String) IPv4 CIDR (e.g., '192.168.1.1/24').
- `cidr6` (String) IPv6 CIDR.
- `comments` (String) Comments for the interface.
- `comments6` (String) IPv6 comments.
- `gateway` (String) IPv4 gateway.
- `gateway6` (String) IPv6 gateway.
- `method` (String) IPv4 method (static, dhcp, manual).
- `method6` (String) IPv6 method (static, dhcp, manual).
- `mtu` (Number) Interface MTU.
- `netmask` (String) IPv4 netmask.
- `netmask6` (Number) IPv6 prefix length.
- `slaves` (String) Slave interfaces for bond.
- `vlan_id` (Number) VLAN ID.
- `vlan_raw_device` (String) The underlying device for the VLAN interface.

### Read-Only

- `id` (String) The ID of this resource.
