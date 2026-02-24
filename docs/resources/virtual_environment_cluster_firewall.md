---
page_title: "proxmox_virtual_environment_cluster_firewall Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages cluster-level firewall options for Proxmox VE.
---

# proxmox_virtual_environment_cluster_firewall (Resource)

Manages cluster-level firewall options for Proxmox VE. This is a singleton resource — only one instance exists per cluster. Destroying this resource removes it from Terraform state only; the underlying configuration is not deleted.

## Example Usage

```terraform
resource "proxmox_virtual_environment_cluster_firewall" "example" {
  enable     = true
  policy_in  = "DROP"
  policy_out = "ACCEPT"

  log_ratelimit_enable = true
  log_ratelimit_rate   = "1/second"
  log_ratelimit_burst  = 5
}
```

## Schema

### Optional

- `enable` (Boolean) Whether the cluster firewall is enabled.
- `log_ratelimit_burst` (Number) Initial burst value for log rate limiting.
- `log_ratelimit_enable` (Boolean) Whether log rate limiting is enabled.
- `log_ratelimit_rate` (String) Rate for log rate limiting (e.g. `1/second`).
- `policy_in` (String) Default inbound policy (`ACCEPT`, `DROP`, `REJECT`).
- `policy_out` (String) Default outbound policy (`ACCEPT`, `DROP`, `REJECT`).

### Read-Only

- `id` (String) The ID of this resource.
