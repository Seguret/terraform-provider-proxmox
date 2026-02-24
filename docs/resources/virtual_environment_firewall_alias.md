---
page_title: "proxmox_virtual_environment_firewall_alias Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE firewall alias (named IP/CIDR).
---

# proxmox_virtual_environment_firewall_alias (Resource)

Manages a Proxmox VE firewall alias (named IP/CIDR).

## Example Usage

### Office Network Alias

```terraform
resource "proxmox_virtual_environment_firewall_alias" "office_network" {
  scope   = "cluster"
  name    = "office-net"
  cidr    = "203.0.113.0/24"
  comment = "Office network"
}
```

### Multiple Network Aliases

```terraform
resource "proxmox_virtual_environment_firewall_alias" "vpn_clients" {
  scope   = "cluster"
  name    = "vpn-clients"
  cidr    = "192.168.100.0/24"
  comment = "Remote VPN clients"
}

resource "proxmox_virtual_environment_firewall_alias" "management_subnet" {
  scope   = "cluster"
  name    = "mgmt-subnet"
  cidr    = "10.0.0.0/24"
  comment = "Management network"
}

resource "proxmox_virtual_environment_firewall_alias" "backup_server" {
  scope   = "cluster"
  name    = "backup-srv"
  cidr    = "10.0.5.100/32"
  comment = "Backup server IP"
}
```

### Node-Specific Aliases

```terraform
resource "proxmox_virtual_environment_firewall_alias" "node_trusted" {
  scope   = "node/pve"
  name    = "trusted-nodes"
  cidr    = "10.0.0.0/23"
  comment = "Trusted Proxmox nodes"
}
```


## Schema

### Required

- `cidr` (String) The IP or CIDR that the alias represents.
- `name` (String) The alias name.
- `scope` (String) The firewall scope: 'cluster', 'node/<node>', 'vm/<node>/<vmid>', 'ct/<node>/<vmid>'.

### Optional

- `comment` (String) Alias description.

### Read-Only

- `id` (String) The ID of this resource.
