# proxmox_virtual_environment_role (Data Source)

Retrieves information about a specific Proxmox VE access role.




## Example Usage

```hcl
data "proxmox_virtual_environment_role" "example" {
  role_id = "PVEVMAdmin"
}
```

## Schema

### Required

- `role_id` (String) The role identifier to look up.

### Read-Only

- `id` (String) The role identifier.
- `privileges` (List of String) The list of privileges assigned to the role.
