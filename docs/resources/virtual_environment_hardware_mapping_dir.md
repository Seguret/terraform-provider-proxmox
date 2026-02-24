---
page_title: "proxmox_virtual_environment_hardware_mapping_dir Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE cluster-level directory hardware mapping.
---

# proxmox_virtual_environment_hardware_mapping_dir (Resource)

Manages a Proxmox VE cluster-level directory hardware mapping, which maps a logical name to per-node directory paths.

## Example Usage

```terraform
resource "proxmox_virtual_environment_hardware_mapping_dir" "example" {
  mapping_id = "shared-data"
  comment    = "Shared data directory"

  map = [
    {
      node = "pve1"
      path = "/mnt/shared"
    },
    {
      node = "pve2"
      path = "/mnt/shared"
    }
  ]
}
```

## Schema

### Required

- `map` (Attributes List) List of per-node directory path entries. (see [below for nested schema](#nestedatt--map))
- `mapping_id` (String) The hardware mapping identifier. Changing this forces a new resource.

### Optional

- `comment` (String) A human-readable description.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--map"></a>
### Nested Schema for `map`

Required:

- `node` (String) The Proxmox VE node name.
- `path` (String) The directory path on the node.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_hardware_mapping_dir.example shared-data
```
