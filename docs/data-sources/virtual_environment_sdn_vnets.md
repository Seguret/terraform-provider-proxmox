# proxmox_virtual_environment_sdn_vnets (Data Source)

Retrieves the list of Proxmox VE SDN VNets.

## Example Usage

### Retrieve the list of Proxmox VE SDN VNets
```hcl
data "proxmox_virtual_environment_sdn_vnets" "test" {}

output "test" {
  value = data.proxmox_virtual_environment_sdn_vnets.test
}
```

## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `vnets` (List of String) VNet names.
- `zones` (List of String) Zone each VNet belongs to.
