# proxmox_virtual_environment_apt_repositories (Data Source)

Retrieves all APT repositories configured on a Proxmox VE node.




## Schema

### Required

- `node_name` (String) The node name.

### Read-Only

- `enabled` (List of Boolean) Whether each repository is enabled.
- `files` (List of String) Source file path for each repository.
- `id` (String) The ID of this resource.
- `suites` (List of String) Suite (e.g. bookworm) of each repository.
- `uris` (List of String) URI of each repository.
