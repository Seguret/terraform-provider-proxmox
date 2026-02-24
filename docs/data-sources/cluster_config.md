# proxmox_cluster_config (Data Source)

Retrieves the cluster configuration including node information.




## Schema

### Read-Only

- `id` (String) Data source identifier.
- `nodes` (Attributes List) List of nodes in the cluster. (see [below for nested schema](#nestedatt--nodes))
- `totem_interface` (String) The totem interface configuration.

<a id="nestedatt--nodes"></a>
### Nested Schema for `nodes`

Read-Only:

- `ip` (String) Node IP address.
- `name` (String) Node name.
- `node_id` (Number) Node ID.
- `ring0_addr` (String) Ring 0 address.
- `ring1_addr` (String) Ring 1 address.
