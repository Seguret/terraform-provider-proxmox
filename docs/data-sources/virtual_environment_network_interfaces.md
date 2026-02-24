# proxmox_virtual_environment_network_interfaces (Data Source)

Retrieves the list of network interfaces on a Proxmox VE node.

## Example

```hcl
data "proxmox_virtual_environment_network_interfaces" "pve" {
  node_name = "pve"
}

output "interfaces" {
  value = data.proxmox_virtual_environment_network_interfaces.pve
}
```

## Schema

### Required

- `node_name` (String) The node name.

### Read-Only

- `active` (List of Boolean) Whether each interface is active.
- `addresses` (List of String) IPv4 addresses.
- `cidrs` (List of String) CIDR notation addresses.
- `id` (String) The ID of this resource.
- `names` (List of String) Interface names.
- `types` (List of String) Interface types.
