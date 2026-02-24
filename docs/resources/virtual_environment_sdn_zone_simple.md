---
page_title: "proxmox_virtual_environment_sdn_zone_simple Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Simple SDN zone in Proxmox VE.
---

# proxmox_virtual_environment_sdn_zone_simple (Resource)

Manages a Simple SDN zone in Proxmox VE. The Simple zone is the most basic SDN zone type — it provides isolated layer-2 segments without encapsulation.

## Example Usage

```terraform
resource "proxmox_virtual_environment_sdn_zone_simple" "example" {
  zone  = "simple-zone"
  nodes = "pve1,pve2"
}
```

## Schema

### Required

- `zone` (String) The SDN zone identifier. Changing this forces a new resource.

### Optional

- `dns` (String) DNS plugin name.
- `dns_zone` (String) DNS domain.
- `ipam` (String) IPAM plugin name.
- `nodes` (String) Comma-separated list of nodes where the zone is deployed.
- `reverse_dns` (String) Reverse DNS plugin name.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_sdn_zone_simple.example simple-zone
```
