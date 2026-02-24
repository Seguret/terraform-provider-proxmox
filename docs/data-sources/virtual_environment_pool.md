# proxmox_virtual_environment_pool (Data Source)

Retrieves information about a specific Proxmox VE resource pool.




## Example Usage

```hcl
data "proxmox_virtual_environment_pool" "example" {
  pool_id = "dev"
}
```

## Schema

### Required

- `pool_id` (String) The pool identifier to look up.

### Read-Only

- `comment` (String) The pool description.
- `id` (String) The pool identifier.
- `members` (Attributes List) The list of members in the pool. (see [below for nested schema](#nestedatt--members))

<a id="nestedatt--members"></a>
### Nested Schema for `members`

Read-Only:

- `id` (String) The member resource identifier.
- `type` (String) The member type (e.g. `qemu`, `lxc`, `storage`).
- `vmid` (Number) The VM/container ID (if applicable).
