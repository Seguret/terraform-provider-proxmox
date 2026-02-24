# proxmox_virtual_environment_node_hardware_pci (Data Source)

Retrieves the list of PCI devices on a Proxmox VE node.

### Example Usage

```hcl
data "proxmox_virtual_environment_node_hardware_pci" "pve" {
  node_name = "pve"
}

output "cpu_model" {
  value = data.proxmox_virtual_environment_node_hardware_pci.pve
}
```

## Schema

### Required

- `node_name` (String) The name of the Proxmox VE node.

### Read-Only

- `devices` (Attributes List) The list of PCI devices on the node. (see [below for nested schema](#nestedatt--devices))
- `id` (String) Placeholder identifier.

<a id="nestedatt--devices"></a>
### Nested Schema for `devices`

Read-Only:

- `class` (String) The PCI device class.
- `device` (String) The PCI device name.
- `device_id` (String) The PCI device ID.
- `id` (String) The PCI device identifier.
- `iommu_group` (Number) The IOMMU group the device belongs to.
- `vendor` (String) The PCI device vendor name.
- `vendor_id` (String) The PCI vendor ID.
