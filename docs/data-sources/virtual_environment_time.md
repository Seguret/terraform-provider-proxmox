# proxmox_virtual_environment_time (Data Source)

Retrieves the current time and timezone of a Proxmox VE node.




## Example Usage

```hcl
data "proxmox_virtual_environment_time" "example" {
  node_name = "pve"
}
```

## Schema

### Required

- `node_name` (String) The node to retrieve time information from.

### Read-Only

- `id` (String) The time datasource identifier.
- `local_time` (Number) The local time on the node as a Unix timestamp.
- `time` (Number) The UTC time on the node as a Unix timestamp.
- `timezone` (String) The timezone configured on the node (e.g. `Europe/Berlin`).
