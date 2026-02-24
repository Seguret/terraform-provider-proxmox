# proxmox_virtual_environment_acme_directories (Data Source)

Retrieves the list of ACME directory endpoints available in Proxmox VE.




## Schema

### Read-Only

- `directories` (Attributes List) The list of ACME directory entries. (see [below for nested schema](#nestedatt--directories))
- `id` (String) Placeholder identifier.

<a id="nestedatt--directories"></a>
### Nested Schema for `directories`

Read-Only:

- `name` (String) The human-readable name of the ACME directory.
- `url` (String) The ACME directory URL.
