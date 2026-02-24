# proxmox_virtual_environment_vms (Data Source)

Retrieves the list of VMs on a Proxmox VE node.

## Example Usage

### List all VMs on a node

```terraform
data "proxmox_virtual_environment_vms" "all" {
  node_name = "pve"
}

output "vm_names" {
  value = data.proxmox_virtual_environment_vms.all.names
}

output "vm_statuses" {
  value = data.proxmox_virtual_environment_vms.all.statuses
}
```

### Get specific VM information

```terraform
data "proxmox_virtual_environment_vms" "pve" {
  node_name = "pve"
}

locals {
  vm_info = zipmap(
    data.proxmox_virtual_environment_vms.pve.names,
    data.proxmox_virtual_environment_vms.pve.vmids
  )
}

output "all_vms" {
  value = local.vm_info
}
```



## Schema

### Required

- `node_name` (String) The node name.

### Read-Only

- `cpus` (List of Number) The number of CPUs for each VM.
- `id` (String) The ID of this resource.
- `max_disk` (List of Number) The maximum disk size in bytes for each VM.
- `max_memory` (List of Number) The maximum memory in bytes for each VM.
- `names` (List of String) The VM names.
- `statuses` (List of String) The VM statuses.
- `tags` (List of String) The VM tags.
- `template` (List of Boolean) Whether each VM is a template.
- `uptime` (List of Number) The uptime in seconds for each VM.
- `vmids` (List of Number) The VM IDs.
