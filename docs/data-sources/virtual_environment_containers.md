# proxmox_virtual_environment_containers (Data Source)

Retrieves the list of LXC containers on a Proxmox VE node.

## Example Usage

### List all containers on a node

```terraform
data "proxmox_virtual_environment_containers" "all" {
  node_name = "pve"
}

output "container_names" {
  value = data.proxmox_virtual_environment_containers.all.names
}

output "container_count" {
  value = length(data.proxmox_virtual_environment_containers.all.names)
}
```

### Filter running containers

```terraform
data "proxmox_virtual_environment_containers" "all" {
  node_name = "pve"
}

locals {
  running_containers = [
    for i, name in data.proxmox_virtual_environment_containers.all.names :
    name
    if data.proxmox_virtual_environment_containers.all.statuses[i] == "running"
  ]
}

output "running" {
  value = local.running_containers
}
```



## Schema

### Required

- `node_name` (String) The node name.

### Read-Only

- `cpus` (List of Number) The number of CPUs.
- `id` (String) The ID of this resource.
- `max_disk` (List of Number) The maximum disk size in bytes.
- `max_memory` (List of Number) The maximum memory in bytes.
- `names` (List of String) The container names.
- `statuses` (List of String) The container statuses.
- `tags` (List of String) The container tags.
- `template` (List of Boolean) Whether each container is a template.
- `uptime` (List of Number) The uptime in seconds.
- `vmids` (List of Number) The container IDs.
