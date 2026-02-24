# proxmox_virtual_environment_container_snapshots (Data Source)

Retrieves the list of snapshots for a Proxmox VE container.




## Schema

### Required

- `node_name` (String) The node name.
- `vmid` (Number) The container ID.

### Read-Only

- `descriptions` (List of String) Snapshot descriptions.
- `id` (String) The ID of this resource.
- `snap_names` (List of String) Snapshot names.
- `snaptimes` (List of Number) Snapshot creation times as UNIX timestamps.
