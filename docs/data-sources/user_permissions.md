# proxmox_user_permissions (Data Source)

Retrieves permissions for a specific Proxmox VE user.




## Schema

### Required

- `user_id` (String) The user ID to query permissions for (e.g., 'root@pam').

### Optional

- `path` (String) Optional path to filter permissions (e.g., '/vms/100'). If empty, returns all permissions.

### Read-Only

- `id` (String) The ID of this resource.
- `permissions` (Map of Map of Number) Map of permissions by path. Each path contains a map of privilege to boolean.
