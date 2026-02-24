# proxmox_virtual_environment_ha_resource (Data Source)

Retrieves information about a specific Proxmox VE High Availability resource.




## Example Usage

```hcl
data "proxmox_virtual_environment_ha_resource" "example" {
  sid = "vm:100"
}
```

## Schema

### Required

- `sid` (String) The HA resource SID (e.g. `vm:100`).

### Read-Only

- `comment` (String) Description of the HA resource.
- `group` (String) The HA group the resource belongs to.
- `id` (String) The HA resource SID.
- `max_relocate` (Number) Maximum number of relocation attempts.
- `max_restart` (Number) Maximum number of restart attempts.
- `state` (String) The desired state of the HA resource.
- `type` (String) The HA resource type.
