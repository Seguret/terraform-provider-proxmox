# proxmox_virtual_environment_sdn_vnet (Data Source)

Retrieves information about a specific Proxmox VE SDN VNet.




## Example Usage

```hcl
data "proxmox_virtual_environment_sdn_vnet" "example" {
  vnet = "myvnet"
}
```

## Schema

### Required

- `vnet` (String) The VNet name to look up.

### Read-Only

- `alias` (String) An alias for the VNet.
- `id` (String) The VNet identifier.
- `tag` (Number) The VLAN tag assigned to the VNet.
- `zone` (String) The SDN zone the VNet belongs to.
