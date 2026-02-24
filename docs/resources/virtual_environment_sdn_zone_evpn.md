---
page_title: "proxmox_virtual_environment_sdn_zone_evpn Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages an EVPN SDN zone in Proxmox VE.
---

# proxmox_virtual_environment_sdn_zone_evpn (Resource)

Manages an EVPN (Ethernet VPN) SDN zone in Proxmox VE.

## Example Usage

```terraform
resource "proxmox_virtual_environment_sdn_zone_evpn" "example" {
  zone       = "evpn-zone"
  controller = "mycontroller"
  vrf_vxlan  = 4000
  nodes      = "pve1,pve2"
}
```

## Schema

### Required

- `controller` (String) EVPN controller name.
- `vrf_vxlan` (Number) VRF VxLAN tag number.
- `zone` (String) The SDN zone identifier. Changing this forces a new resource.

### Optional

- `advertise_subnets` (Boolean) Advertise subnets via EVPN.
- `dns` (String) DNS plugin name.
- `dns_zone` (String) DNS domain.
- `exit_nodes` (String) Comma-separated list of nodes acting as exit nodes for the zone.
- `exit_nodes_local_routing` (Boolean) Enable local routing on exit nodes.
- `ipam` (String) IPAM plugin name.
- `nodes` (String) Comma-separated list of nodes where the zone is deployed.
- `reverse_dns` (String) Reverse DNS plugin name.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_sdn_zone_evpn.example evpn-zone
```
