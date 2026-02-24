# proxmox_virtual_environment_node_services (Data Source)

Retrieves the list of system services on a Proxmox VE node.

## Example Usage

```hcl
data "proxmox_virtual_environment_node_services" "pve" {
  node_name = "pve"
}

output "cpu_model" {
  value = data.proxmox_virtual_environment_node_services.pve
}
```

## Schema

### Required

- `node_name` (String) The name of the Proxmox VE node.

### Read-Only

- `id` (String) Placeholder identifier.
- `services` (Attributes List) The list of services on the node. (see [below for nested schema](#nestedatt--services))

<a id="nestedatt--services"></a>
### Nested Schema for `services`

Read-Only:

- `active_state` (String) The systemd active state (e.g. 'active', 'inactive').
- `desc` (String) A short description of the service.
- `name` (String) The service name.
- `state` (String) The current run state of the service (e.g. 'running', 'stopped').
- `sub_state` (String) The systemd sub-state (e.g. 'running', 'dead').
