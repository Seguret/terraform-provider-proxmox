---
page_title: "proxmox_virtual_environment_cloned_vm Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Creates and manages a clone of an existing Proxmox VE virtual machine.
---

# proxmox_virtual_environment_cloned_vm (Resource)

Creates and manages a clone of an existing Proxmox VE virtual machine.

## Example Usage

```terraform
resource "proxmox_virtual_environment_cloned_vm" "example" {
  node_name   = "pve"
  source_node = "pve"
  source_vmid = 9000

  vm_id       = 100
  name        = "my-clone"
  description = "Cloned from template"
  full_clone  = true
  tags        = "production;web"
}
```

## Schema

### Required

- `node_name` (String) The name of the Proxmox VE node that will host the cloned VM. Changing this forces a new resource.
- `source_node` (String) The Proxmox VE node where the source VM resides. Changing this forces a new resource.
- `source_vmid` (Number) The VMID of the VM to clone. Changing this forces a new resource.

### Optional

- `description` (String) A description for the cloned VM.
- `full_clone` (Boolean) Whether to perform a full clone (`true`) or a linked clone (`false`). Defaults to `true`.
- `name` (String) The name of the cloned VM.
- `tags` (String) Semicolon-separated tags for the cloned VM.
- `target_node` (String) The target node for the cloned VM (for cross-node clones). Changing this forces a new resource.
- `target_storage` (String) The storage ID where the full clone's disks should be placed. Changing this forces a new resource.
- `vm_id` (Number) The VMID for the cloned VM. If omitted, Proxmox auto-assigns the next available ID.

### Read-Only

- `id` (String) The resource identifier in the form `{node_name}/{vmid}`.
- `status` (String) The current status of the VM (e.g. `running`, `stopped`).

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_cloned_vm.example pve/100
```
