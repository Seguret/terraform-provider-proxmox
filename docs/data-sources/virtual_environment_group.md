# proxmox_virtual_environment_group (Data Source)

Retrieves information about a specific Proxmox VE access group.




## Example Usage

```hcl
data "proxmox_virtual_environment_group" "example" {
  group_id = "admins"
}
```

## Schema

### Required

- `group_id` (String) The group identifier to look up.

### Read-Only

- `comment` (String) The group description.
- `id` (String) The group identifier.
- `members` (List of String) The list of users that are members of this group.
