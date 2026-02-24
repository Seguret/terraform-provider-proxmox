# proxmox_virtual_environment_vm_snapshots (Data Source)

Retrieves the list of snapshots for a Proxmox VE virtual machine.

## Example Usage

```hcl
data "proxmox_virtual_environment_vm_snapshots" "vms" {
  node_name = "pve"
  vmid = 104
}

output "vm_list" {
  value       = data.proxmox_virtual_environment_vm_snapshots.vms
  description = "List of VMs with their snapshots"
} 
```
### Example Usage with Snapshot Descriptions

```hcl
data "proxmox_virtual_environment_vm_snapshots" "vms" {
  node_name = "pve"
  vmid = 104
}

output "vm_list" {
  value       = data.proxmox_virtual_environment_vm_snapshots.vms.descriptions
  description = "List of snapshot descriptions for the specified VM"
} 
```

## Schema

### Required

- `node_name` (String) The node name.
- `vmid` (Number) The VM ID.

### Read-Only

- `descriptions` (List of String) Snapshot descriptions.
- `id` (String) The ID of this resource.
- `snap_names` (List of String) Snapshot names.
- `snaptimes` (List of Number) Snapshot creation times as UNIX timestamps.
