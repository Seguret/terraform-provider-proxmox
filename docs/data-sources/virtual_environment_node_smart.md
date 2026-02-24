# proxmox_virtual_environment_node_smart (Data Source)

Retrieves SMART data for a disk on a Proxmox VE node.

## Example Usage

### Basic Example
```hcl
data "proxmox_virtual_environment_node_smart" "pve" {
  node_name = "pve"
  disk = "/dev/nvme0n1"
}

output "cpu_model" {
  value = data.proxmox_virtual_environment_node_smart.pve
}
```

## Schema

### Required

- `disk` (String) The device path of the disk (e.g. /dev/sda).
- `node_name` (String) The node name.

### Read-Only

- `health` (String) The SMART health status.
- `id` (String) Placeholder identifier.
- `text` (String) The raw SMART text output.
- `type` (String) The SMART type (e.g. SATA, NVMe).
- `wearout` (Number) The wearout indicator for SSDs (percentage remaining).
