# proxmox_virtual_environment_cluster_status (Data Source)

Retrieves the current status of the Proxmox VE cluster.




## Schema

### Read-Only

- `entries` (Attributes List) The list of cluster status entries. (see [below for nested schema](#nestedatt--entries))
- `id` (String) Placeholder identifier.

<a id="nestedatt--entries"></a>
### Nested Schema for `entries`

Read-Only:

- `id` (String) The entry identifier.
- `ip` (String) The IP address of the node.
- `local` (Number) Whether this is the local node (1) or not (0).
- `name` (String) The name of the node or cluster.
- `node_id` (Number) The numeric node ID.
- `nodes` (Number) The total number of nodes in the cluster.
- `online` (Number) Whether the node is online (1) or offline (0).
- `quorate` (Number) Whether the cluster has quorum (1) or not (0).
- `type` (String) The entry type (e.g., 'node', 'cluster').
- `version` (Number) The configuration version.
