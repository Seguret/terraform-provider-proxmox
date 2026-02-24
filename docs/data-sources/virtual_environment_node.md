# proxmox_virtual_environment_node (Data Source)

Retrieves the status and details of a specific Proxmox VE node.

## Example Usage

```terraform
# Get detailed status of a specific node
data "proxmox_virtual_environment_node" "pve1" {
  node_name = "pve1"
}

output "cpu_model" {
  value = data.proxmox_virtual_environment_node.pve1.cpu_model
}

output "memory_total_gb" {
  value = data.proxmox_virtual_environment_node.pve1.memory_total / 1073741824
}

output "uptime_hours" {
  value = data.proxmox_virtual_environment_node.pve1.uptime / 3600
}

output "pve_version" {
  value = data.proxmox_virtual_environment_node.pve1.pve_version
}
```


## Schema

### Required

- `node_name` (String) The name of the Proxmox VE node.

### Read-Only

- `boot_mode` (String) The boot mode (e.g., 'efi' or 'bios').
- `cpu_cores` (Number) The number of CPU cores.
- `cpu_mhz` (String) The CPU clock speed in MHz.
- `cpu_model` (String) The CPU model name.
- `cpu_sockets` (Number) The number of CPU sockets.
- `cpu_threads` (Number) The number of CPU threads.
- `cpu_usage` (Number) The current CPU utilization (0.0-1.0).
- `id` (String) Placeholder identifier.
- `kernel_version` (String) The kernel version string.
- `load_average` (List of String) The system load averages (1, 5, 15 minutes).
- `memory_free` (Number) Free memory in bytes.
- `memory_total` (Number) Total memory in bytes.
- `memory_used` (Number) Used memory in bytes.
- `pve_version` (String) The Proxmox VE version string.
- `rootfs_free` (Number) Free root filesystem in bytes.
- `rootfs_total` (Number) Total root filesystem size in bytes.
- `rootfs_used` (Number) Used root filesystem in bytes.
- `secure_boot` (Boolean) Whether secure boot is enabled.
- `swap_free` (Number) Free swap in bytes.
- `swap_total` (Number) Total swap in bytes.
- `swap_used` (Number) Used swap in bytes.
- `uptime` (Number) The node uptime in seconds.
