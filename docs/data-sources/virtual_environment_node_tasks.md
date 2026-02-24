# proxmox_virtual_environment_node_tasks (Data Source)

Retrieves the list of recent tasks on a Proxmox VE node.


## Example Usage

```hcl
data "proxmox_virtual_environment_node_tasks" "pve" {
  node_name = "pve"
}

output "cpu_model" {
  value = data.proxmox_virtual_environment_node_tasks.pve
}
```

## Schema

### Required

- `node_name` (String) The name of the Proxmox VE node.

### Read-Only

- `id` (String) Placeholder identifier.
- `tasks` (Attributes List) The list of tasks on the node. (see [below for nested schema](#nestedatt--tasks))

<a id="nestedatt--tasks"></a>
### Nested Schema for `tasks`

Read-Only:

- `end_time` (Number) The task end time as a Unix timestamp.
- `start_time` (Number) The task start time as a Unix timestamp.
- `status` (String) The task status (e.g., 'OK', error message).
- `task_id` (String) The task-specific identifier (e.g., VM ID).
- `type` (String) The task type.
- `upid` (String) The unique task process ID.
- `user` (String) The user who initiated the task.
