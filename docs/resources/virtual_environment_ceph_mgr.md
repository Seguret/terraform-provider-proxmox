---
page_title: "proxmox_virtual_environment_ceph_mgr Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Ceph MGR daemon on a Proxmox VE node.
---

# proxmox_virtual_environment_ceph_mgr (Resource)

Manages a Ceph Manager (MGR) daemon on a Proxmox VE node.

## Example Usage

```terraform
resource "proxmox_virtual_environment_ceph_mgr" "example" {
  node_name = "pve"
  mgr_id    = "pve"
}
```

## Schema

### Required

- `mgr_id` (String) The MGR daemon ID (typically the node name). Changing this forces a new resource.
- `node_name` (String) The Proxmox node on which to create the MGR. Changing this forces a new resource.

### Read-Only

- `addr` (String) The MGR daemon address.
- `id` (String) The ID of this resource.
- `state` (String) The MGR daemon state (`active` or `standby`).

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_ceph_mgr.example pve/pve
```
