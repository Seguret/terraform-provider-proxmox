# proxmox_virtual_environment_dns (Data Source)

Retrieves the DNS configuration of a Proxmox VE node.




## Example Usage

```hcl
data "proxmox_virtual_environment_dns" "example" {
  node_name = "pve"
}
```

## Schema

### Required

- `node_name` (String) The node to retrieve DNS configuration from.

### Read-Only

- `domain` (String) The search domain configured on the node.
- `id` (String) The DNS datasource identifier.
- `servers` (List of String) The list of DNS server addresses configured on the node.
