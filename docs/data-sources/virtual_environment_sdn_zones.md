# proxmox_virtual_environment_sdn_zones (Data Source)

Retrieves the list of Proxmox VE SDN zones.

## Example

### Basic Usage
```hcl
data "proxmox_virtual_environment_sdn_zones" "test" {
}

output "test" {
  value = data.proxmox_virtual_environment_sdn_zones.test
}
```

## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `types` (List of String) SDN zone types.
- `zones` (List of String) SDN zone names.
