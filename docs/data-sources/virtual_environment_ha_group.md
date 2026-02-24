# proxmox_virtual_environment_ha_group (Data Source)

Retrieves information about a specific Proxmox VE High Availability group.




## Example Usage

```hcl
data "proxmox_virtual_environment_ha_group" "example" {
  group = "production"
}
```

## Schema

### Required

- `group` (String) The HA group name.

### Read-Only

- `comment` (String) Description of the HA group.
- `id` (String) The HA group identifier.
- `nofailback` (Boolean) Whether to prevent failback.
- `nodes` (String) Comma-separated list of node:priority pairs.
- `restricted` (Boolean) Whether HA resources bound to this group may only run on the defined nodes.
