# proxmox_virtual_environment_cluster_tasks (Data Source)

Retrieves the list of recent tasks across the Proxmox VE cluster.




## Schema

### Read-Only

- `id` (String) Placeholder identifier.
- `tasks` (Attributes List) The list of cluster tasks. (see [below for nested schema](#nestedatt--tasks))

<a id="nestedatt--tasks"></a>
### Nested Schema for `tasks`

Read-Only:

- `end_time` (Number) The task end time as a Unix timestamp.
- `node` (String) The node on which the task ran.
- `pid` (Number) The process ID of the task.
- `start_time` (Number) The task start time as a Unix timestamp.
- `status` (String) The task status (e.g., 'OK', error message).
- `task_id` (String) The task-specific identifier (e.g., VM ID).
- `type` (String) The task type.
- `upid` (String) The unique task process ID.
- `user` (String) The user who initiated the task.
