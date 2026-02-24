# proxmox_virtual_environment_hardware_mapping_pci_list (Data Source)

Retrieves the list of Proxmox VE cluster PCI hardware mapping names.

## Example Usage

```hcl
data "proxmox_virtual_environment_hardware_mapping_pci_list" "pve" {}

output "cpu_model" {
  value = data.proxmox_virtual_environment_hardware_mapping_pci_list.pve
}
```

## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (List of String) PCI hardware mapping names.
