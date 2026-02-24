# proxmox_virtual_environment_sdn_subnet (Data Source)

Retrieves information about a specific Proxmox VE SDN subnet.




## Example Usage

```hcl
data "proxmox_virtual_environment_sdn_subnet" "example" {
  vnet   = "myvnet"
  subnet = "10.0.0.0/24"
}
```

## Schema

### Required

- `subnet` (String) The subnet CIDR (e.g. `10.0.0.0/24`).
- `vnet` (String) The VNet this subnet belongs to.

### Read-Only

- `gateway` (String) The subnet gateway IP address.
- `id` (String) The subnet identifier.
- `type` (String) The subnet type.
