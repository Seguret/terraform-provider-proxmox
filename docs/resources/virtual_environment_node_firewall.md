---
page_title: "proxmox_virtual_environment_node_firewall Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages node-level firewall options for a Proxmox VE node.
---

# proxmox_virtual_environment_node_firewall (Resource)

Manages node-level firewall options for a Proxmox VE node. This is a singleton resource per node — only one instance should exist per node. Destroying this resource removes it from Terraform state only.

## Example Usage

```terraform
resource "proxmox_virtual_environment_node_firewall" "example" {
  node_name  = "pve"
  enable     = true
  policy_in  = "DROP"
  policy_out = "ACCEPT"
}
```

## Schema

### Required

- `node_name` (String) The node name.

### Optional

- `enable` (Boolean) Whether the node firewall is enabled.
- `log_ratelimit_burst` (Number) Initial burst value for log rate limiting.
- `log_ratelimit_enable` (Boolean) Whether log rate limiting is enabled.
- `log_ratelimit_rate` (String) Rate for log rate limiting (e.g. `1/second`).
- `policy_in` (String) Default inbound policy (`ACCEPT`, `DROP`, `REJECT`).
- `policy_out` (String) Default outbound policy (`ACCEPT`, `DROP`, `REJECT`).

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_node_firewall.example pve
```
