# proxmox_virtual_environment_container (Data Source)

Retrieves information about a specific Proxmox VE LXC container.




## Schema

### Required

- `node_name` (String) The node name.
- `vmid` (Number) The container ID.

### Read-Only

- `cores` (Number) Number of CPU cores.
- `description` (String) The container description.
- `hostname` (String) The container hostname.
- `id` (String) The ID of this resource.
- `memory` (Number) Memory in MiB.
- `on_boot` (Boolean) Whether the container starts on host boot.
- `os_type` (String) The OS type.
- `status` (String) The container status.
- `swap` (Number) Swap in MiB.
- `tags` (String) The container tags.
- `template` (Boolean) Whether the container is a template.
- `unprivileged` (Boolean) Whether the container runs unprivileged.
