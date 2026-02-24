---
page_title: "proxmox_virtual_environment_ceph_mon Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Ceph monitor (MON) on a Proxmox VE node.
---

# proxmox_virtual_environment_ceph_mon (Resource)

Manages a Ceph Monitor (MON) on a Proxmox VE node.

## Example Usage

```terraform
resource "proxmox_virtual_environment_ceph_mon" "example" {
  node_name = "pve"
  mon_name  = "pve"
}
```

## Schema

### Required

- `mon_name` (String) The MON name (typically the node name). Changing this forces a new resource.
- `node_name` (String) The Proxmox node on which to create the MON. Changing this forces a new resource.

### Read-Only

- `addr` (String) The MON bind address.
- `host` (String) The MON host address.
- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_ceph_mon.example pve/pve
```
