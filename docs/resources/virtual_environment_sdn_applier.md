---
page_title: "proxmox_virtual_environment_sdn_applier Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Applies pending Proxmox VE SDN configuration changes to the cluster.
---

# proxmox_virtual_environment_sdn_applier (Resource)

Applies pending Proxmox VE SDN (Software-Defined Networking) configuration changes to the cluster. Use this resource to trigger `apply` after making SDN zone, VNet, or subnet changes.

## Example Usage

```terraform
resource "proxmox_virtual_environment_sdn_zone_simple" "internal" {
  zone = "internal"
}

resource "proxmox_virtual_environment_sdn_applier" "apply" {
  keep_up_to_date = true

  depends_on = [proxmox_virtual_environment_sdn_zone_simple.internal]
}
```

## Schema

### Optional

- `keep_up_to_date` (Boolean) If `true` (default), call ApplySDN on every Create and Update to keep SDN applied. Defaults to `true`.

### Read-Only

- `id` (String) The ID of this resource.
