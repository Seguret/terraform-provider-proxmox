# proxmox_virtual_environment_hardware_mapping_usb_list (Data Source)

Retrieves the list of Proxmox VE cluster USB hardware mapping names.

## Example Usage

```hcl
data "proxmox_virtual_environment_hardware_mapping_usb_list" "pve" {}

output "cpu_model" {
  value = data.proxmox_virtual_environment_hardware_mapping_usb_list.pve
}
```

## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (List of String) USB hardware mapping names.
