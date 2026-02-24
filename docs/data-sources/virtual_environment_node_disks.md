# proxmox_virtual_environment_node_disks (Data Source)

Retrieves the list of disks on a Proxmox VE node.

## Example Usage

```hcl
data "proxmox_virtual_environment_node_disks" "pve" {
  node_name = "pve"
}

output "interfaces" {
  value = data.proxmox_virtual_environment_node_disks.pve
}
```

## Schema

### Required

- `node_name` (String) The node name.

### Read-Only

- `disks` (Attributes List) The list of disks on the node. (see [below for nested schema](#nestedatt--disks))
- `id` (String) Placeholder identifier.

<a id="nestedatt--disks"></a>
### Nested Schema for `disks`

Read-Only:

- `dev` (String) The device name (e.g. /dev/sda).
- `gpt` (Boolean) Whether the disk has a GPT partition table.
- `health` (String) The SMART health status.
- `model` (String) The disk model name.
- `serial` (String) The disk serial number.
- `size` (Number) The disk size in bytes.
- `type` (String) The disk type (e.g. hdd, ssd).
- `used` (String) How the disk is currently used (e.g. LVM, ZFS).
- `vendor` (String) The disk vendor.
- `wwn` (String) The disk World Wide Name.
