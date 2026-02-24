---
page_title: "proxmox_virtual_environment_firewall_rule Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE firewall rule. Scope determines whether the rule applies to cluster, node, VM, or container.
---

# proxmox_virtual_environment_firewall_rule (Resource)

Manages a Proxmox VE firewall rule. Scope determines whether the rule applies to cluster, node, VM, or container.

## Example Usage

### Cluster-Wide SSH Access Rule

```terraform
resource "proxmox_virtual_environment_firewall_rule" "allow_ssh" {
  scope    = "cluster"
  type     = "in"
  action   = "ACCEPT"
  protocol = "tcp"
  dport    = "22"
  comment  = "Allow SSH access"
  enabled  = true
}
```

### VM-Specific Firewall Rules

```terraform
resource "proxmox_virtual_environment_firewall_rule" "vm_http" {
  scope    = "vm/pve/100"
  type     = "in"
  action   = "ACCEPT"
  protocol = "tcp"
  dport    = "80"
  comment  = "Allow HTTP traffic to VM"
  enabled  = true
}

resource "proxmox_virtual_environment_firewall_rule" "vm_https" {
  scope    = "vm/pve/100"
  type     = "in"
  action   = "ACCEPT"
  protocol = "tcp"
  dport    = "443"
  comment  = "Allow HTTPS traffic to VM"
  enabled  = true
}
```

### Node-Level Rules with Interface

```terraform
resource "proxmox_virtual_environment_firewall_rule" "node_rule_in" {
  scope    = "node/pve"
  type     = "in"
  action   = "ACCEPT"
  protocol = "tcp"
  dport    = "8006"
  iface    = "vmbr0"
  comment  = "Allow Proxmox web UI on vmbr0"
  enabled  = true
}
```

### Container Rules with Logging

```terraform
resource "proxmox_virtual_environment_firewall_rule" "container_dns" {
  scope    = "ct/pve/200"
  type     = "in"
  action   = "ACCEPT"
  protocol = "tcp"
  dport    = "53"
  comment  = "Allow DNS to container"
  log      = "notice"
  enabled  = true
}

resource "proxmox_virtual_environment_firewall_rule" "container_reject" {
  scope    = "ct/pve/200"
  type     = "in"
  action   = "REJECT"
  protocol = "tcp"
  dport    = "3306"
  comment  = "Reject MySQL from outside"
  log      = "warning"
  enabled  = true
}
```


## Schema

### Required

- `action` (String) The rule action (ACCEPT, DROP, REJECT).
- `scope` (String) The firewall scope: 'cluster', 'node/<node>', 'vm/<node>/<vmid>', 'ct/<node>/<vmid>'.
- `type` (String) The rule direction (in or out).

### Optional

- `comment` (String) Rule comment.
- `dest` (String) Destination address/CIDR/IPset.
- `dport` (String) Destination port(s).
- `enabled` (Boolean) Whether the rule is enabled.
- `iface` (String) Network interface.
- `log` (String) Log level (emerg, alert, crit, err, warning, notice, info, debug, nolog).
- `macro` (String) Macro name (e.g., 'SSH', 'HTTP').
- `proto` (String) Protocol (tcp, udp, icmp, etc.).
- `source` (String) Source address/CIDR/IPset.
- `sport` (String) Source port(s).

### Read-Only

- `id` (String) The ID of this resource.
- `pos` (Number) The rule position. Computed after creation.
