---
page_title: "proxmox_virtual_environment_sdn_zone_vxlan Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a VXLAN SDN zone in Proxmox VE.
---

# proxmox_virtual_environment_sdn_zone_vxlan (Resource)

Manages a VXLAN SDN zone in Proxmox VE.

## Example Usage

```terraform
resource "proxmox_virtual_environment_sdn_zone_vxlan" "example" {
  zone  = "vxlan-zone"
  peers = "10.0.0.1,10.0.0.2,10.0.0.3"
}
```

## Schema

### Required

- `peers` (String) Comma-separated list of VXLAN peer IP addresses.
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
terraform import proxmox_virtual_environment_sdn_zone_vxlan.example vxlan-zone
```
