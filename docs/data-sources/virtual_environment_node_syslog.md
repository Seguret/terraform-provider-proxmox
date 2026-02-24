# proxmox_virtual_environment_node_syslog (Data Source)

Retrieves syslog entries from a Proxmox VE node.




## Example Usage

```hcl
data "proxmox_virtual_environment_node_syslog" "example" {
  node_name = "pve"
  limit     = 100
}
```

## Schema

### Required

- `node_name` (String) The name of the node.

### Optional

- `limit` (Number) Maximum number of syslog entries to return (default `500`).
- `start` (Number) Start index for syslog entries (default `0`).

### Read-Only

- `id` (String) The datasource identifier.
- `lines` (List of Number) Line numbers of the syslog entries.
- `texts` (List of String) Text content of the syslog entries.
