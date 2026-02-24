---
page_title: "proxmox_virtual_environment_sdn_zone_vlan Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a VLAN SDN zone in Proxmox VE.
---

# proxmox_virtual_environment_sdn_zone_vlan (Resource)

Manages a VLAN SDN zone in Proxmox VE.

## Example Usage

```terraform
resource "proxmox_virtual_environment_sdn_zone_vlan" "example" {
  zone   = "vlan-zone"
  bridge = "vmbr0"
  nodes  = "pve1,pve2"
}
```

## Schema

### Required

- `bridge` (String) The Linux bridge interface to use for VLAN tagging.
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
terraform import proxmox_virtual_environment_sdn_zone_vlan.example vlan-zone
```
