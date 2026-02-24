# proxmox_virtual_environment_user (Data Source)

Retrieves information about a specific Proxmox VE user.




## Example Usage

```hcl
data "proxmox_virtual_environment_user" "example" {
  user_id = "admin@pam"
}
```

## Schema

### Required

- `user_id` (String) The user identifier to look up (e.g. `user@realm`).

### Read-Only

- `comment` (String) The user description.
- `email` (String) The user's email address.
- `enabled` (Boolean) Whether the user account is enabled.
- `firstname` (String) The user's first name.
- `groups` (List of String) The list of groups the user belongs to.
- `id` (String) The user identifier.
- `lastname` (String) The user's last name.
