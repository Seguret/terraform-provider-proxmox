---
page_title: "proxmox_virtual_environment_sdn_zone_qinq Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a QinQ SDN zone in Proxmox VE.
---

# proxmox_virtual_environment_sdn_zone_qinq (Resource)

Manages a QinQ (802.1ad double-tagging) SDN zone in Proxmox VE.

## Example Usage

```terraform
resource "proxmox_virtual_environment_sdn_zone_qinq" "example" {
  zone          = "qinq-zone"
  bridge        = "vmbr0"
  tag           = 100
  vlan_protocol = "802.1q"
}
```

## Schema

### Required

- `bridge` (String) The Linux bridge interface to use for QinQ double-tagging.
- `tag` (Number) Outer VLAN tag (1-4094) for the QinQ zone.
- `zone` (String) The SDN zone identifier. Changing this forces a new resource.

### Optional

- `dns` (String) DNS plugin name.
- `dns_zone` (String) DNS domain.
- `ipam` (String) IPAM plugin name.
- `nodes` (String) Comma-separated list of nodes where the zone is deployed.
- `reverse_dns` (String) Reverse DNS plugin name.
- `vlan_protocol` (String) VLAN protocol to use (`802.1q` or `802.1ad`). Defaults to `802.1q`.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_sdn_zone_qinq.example qinq-zone
```
