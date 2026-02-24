# proxmox_virtual_environment_apt_changelog (Data Source)

Retrieves the APT changelog for a package on a Proxmox VE node.




## Schema

### Required

- `node_name` (String) The name of the Proxmox VE node.
- `package` (String) The name of the APT package to retrieve the changelog for.

### Read-Only

- `changelog` (String) The changelog content for the package.
- `id` (String) Placeholder identifier.
