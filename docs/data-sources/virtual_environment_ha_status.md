# proxmox_virtual_environment_ha_status (Data Source)

Retrieves the current HA resource status from the Proxmox VE cluster.




## Schema

### Read-Only

- `entries` (Attributes List) The list of HA resource status entries. (see [below for nested schema](#nestedatt--entries))
- `id` (String) Placeholder identifier.

<a id="nestedatt--entries"></a>
### Nested Schema for `entries`

Read-Only:

- `crm_state` (String) The CRM (Cluster Resource Manager) state.
- `max_relocate` (Number) The maximum number of relocations allowed.
- `max_restart` (Number) The maximum number of restarts allowed.
- `node` (String) The node the resource is currently running on.
- `request` (String) The current CRM request for this resource.
- `sid` (String) The HA resource SID (e.g., 'vm:100').
- `state` (String) The current HA state of the resource.
