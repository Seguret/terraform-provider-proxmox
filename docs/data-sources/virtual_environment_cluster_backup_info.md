# proxmox_virtual_environment_cluster_backup_info (Data Source)

Retrieves the list of VMs and containers not covered by any backup job in the Proxmox VE cluster.




## Schema

### Read-Only

- `id` (String) The datasource identifier.
- `vms` (Attributes List) VMs and containers without backup coverage. (see [below for nested schema](#nestedatt--vms))

<a id="nestedatt--vms"></a>
### Nested Schema for `vms`

Read-Only:

- `name` (String) Guest name.
- `type` (String) Guest type (`qemu` or `lxc`).
- `vmid` (Number) VM/CT ID.
