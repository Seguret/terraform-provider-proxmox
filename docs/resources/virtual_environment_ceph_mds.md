---
page_title: "proxmox_virtual_environment_ceph_mds Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Ceph MDS daemon on a Proxmox VE node.
---

# proxmox_virtual_environment_ceph_mds (Resource)

Manages a Ceph MDS (Metadata Server) daemon on a Proxmox VE node.

## Example Usage

```terraform
resource "proxmox_virtual_environment_ceph_mds" "example" {
  node_name = "pve"
  name      = "pve"
}
```

## Schema

### Required

- `name` (String) The MDS daemon name. Changing this forces a new resource.
- `node_name` (String) The Proxmox node on which to create the MDS. Changing this forces a new resource.

### Read-Only

- `addr` (String) The MDS daemon address.
- `id` (String) The ID of this resource.
- `state` (String) The MDS daemon state.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_ceph_mds.example pve/pve
```
