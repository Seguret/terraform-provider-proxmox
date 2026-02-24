# proxmox_virtual_environment_vm (Data Source)

Retrieves information about a specific Proxmox VE virtual machine.

## Example Usage

### Get VM Information

```terraform
data "proxmox_virtual_environment_vm" "vms" {
  node_name = "pve"
  vmid = 104
}

output "vm_list" {
  value       = data.proxmox_virtual_environment_vm.vms
  description = ""
}
```

### Outputting the VM status

```terraform
data "proxmox_virtual_environment_vm" "vms" {
  node_name = "pve"
  vmid = 104
}

output "vm_list" {
  value       = data.proxmox_virtual_environment_vm.vms.status
  description = "List of VMs in Proxmox VE"
}
```

## Schema

### Required

- `node_name` (String) The node name.
- `vmid` (Number) The VM ID.

### Read-Only

- `bios` (String) The BIOS type.
- `cpu_cores` (Number) Number of CPU cores per socket.
- `cpu_sockets` (Number) Number of CPU sockets.
- `cpu_type` (String) The CPU type.
- `description` (String) The VM description.
- `id` (String) The ID of this resource.
- `machine` (String) The machine type.
- `memory` (Number) Memory in MiB.
- `name` (String) The VM name.
- `on_boot` (Boolean) Whether the VM starts on host boot.
- `os_type` (String) The OS type.
- `status` (String) The VM status.
- `tags` (String) The VM tags.
- `template` (Boolean) Whether the VM is a template.
