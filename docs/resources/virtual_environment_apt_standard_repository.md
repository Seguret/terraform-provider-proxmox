---
page_title: "proxmox_virtual_environment_apt_standard_repository Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages the enabled state of a standard (built-in) Proxmox VE APT repository on a node.
---

# proxmox_virtual_environment_apt_standard_repository (Resource)

Manages the enabled state of a standard (built-in) Proxmox VE APT repository on a node.

## Example Usage

```terraform
resource "proxmox_virtual_environment_apt_standard_repository" "no_subscription" {
  node_name = "pve"
  handle    = "pve-no-subscription"
  enabled   = true
}

resource "proxmox_virtual_environment_apt_standard_repository" "enterprise" {
  node_name = "pve"
  handle    = "pve-enterprise"
  enabled   = false
}
```

## Schema

### Required

- `handle` (String) The standard repository handle (e.g. `pve-no-subscription`, `pve-enterprise`). Changing this forces a new resource.
- `node_name` (String) The name of the Proxmox VE node. Changing this forces a new resource.

### Optional

- `enabled` (Boolean) Whether the repository is enabled. Defaults to `true`.

### Read-Only

- `file_path` (String) The path to the APT sources file containing this repository.
- `id` (String) The ID of this resource.
- `index` (Number) The index of the repository entry within its file.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_apt_standard_repository.example pve/pve-no-subscription
```
