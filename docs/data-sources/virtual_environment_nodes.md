# proxmox_virtual_environment_nodes (Data Source)

Retrieves the list of nodes in the Proxmox VE cluster.

## Example Usage

```terraform
# List all nodes in the Proxmox VE cluster
data "proxmox_virtual_environment_nodes" "all" {}

output "node_names" {
  value = data.proxmox_virtual_environment_nodes.all.names
}

output "node_online" {
  value = data.proxmox_virtual_environment_nodes.all.online
}
```


## Schema

### Read-Only

- `cpu_count` (List of Number) The number of CPUs for each node.
- `cpu_utilization` (List of Number) The CPU utilization (0.0-1.0) for each node.
- `id` (String) Placeholder identifier.
- `memory_available` (List of Number) The total available memory in bytes for each node.
- `memory_used` (List of Number) The used memory in bytes for each node.
- `names` (List of String) The node names.
- `online` (List of Boolean) Whether each node is online.
- `ssl_fingerprints` (List of String) The SSL fingerprints for each node.
- `support_levels` (List of String) The support level for each node.
- `uptime` (List of Number) The uptime in seconds for each node.
