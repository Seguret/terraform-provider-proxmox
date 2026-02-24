# proxmox_virtual_environment_apt_versions (Data Source)

Retrieves the list of installed APT package versions on a Proxmox VE node.




## Schema

### Required

- `node_name` (String) The name of the Proxmox VE node.

### Read-Only

- `id` (String) Placeholder identifier.
- `packages` (Attributes List) The list of installed packages and their versions. (see [below for nested schema](#nestedatt--packages))

<a id="nestedatt--packages"></a>
### Nested Schema for `packages`

Read-Only:

- `old_version` (String) The previous installed version (if an upgrade is pending).
- `package` (String) The package name.
- `priority` (String) The package priority.
- `section` (String) The package section.
- `title` (String) The package short description.
- `version` (String) The currently installed version.
