# proxmox_virtual_environment_cluster_resources (Data Source)

Retrieves resources across the Proxmox VE cluster.




## Schema

### Optional

- `resource_type` (String) Optional filter for resource type (e.g., 'vm', 'storage', 'node', 'sdn').

### Read-Only

- `id` (String) Placeholder identifier.
- `resources` (Attributes List) The list of cluster resources. (see [below for nested schema](#nestedatt--resources))

<a id="nestedatt--resources"></a>
### Nested Schema for `resources`

Read-Only:

- `cpu` (Number) The current CPU utilization (0.0-1.0).
- `disk` (Number) The current disk usage in bytes.
- `id` (String) The resource identifier.
- `max_cpu` (Number) The maximum number of CPUs allocated.
- `max_disk` (Number) The maximum disk space allocated in bytes.
- `max_mem` (Number) The maximum memory allocated in bytes.
- `mem` (Number) The current memory usage in bytes.
- `name` (String) The resource name.
- `node` (String) The node the resource resides on.
- `pool` (String) The resource pool the resource belongs to.
- `status` (String) The current status of the resource.
- `storage` (String) The storage identifier (for storage resources).
- `type` (String) The resource type.
- `uptime` (Number) The uptime in seconds.
- `vmid` (Number) The VM or container ID.
