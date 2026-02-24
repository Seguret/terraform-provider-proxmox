# proxmox_virtual_environment_apt_standard_repository (Data Source)

Reads the state of a standard (built-in) Proxmox VE APT repository on a node.

## Example Usage

```terraform
data "proxmox_virtual_environment_apt_standard_repository" "pve_no_subscription" {
  node_name = "pve"
  handle    = "pve-no-subscription"
}

output "apt_repo_enabled" {
  value = data.proxmox_virtual_environment_apt_standard_repository.pve_no_subscription.enabled
}
```


## Schema

### Required

- `handle` (String) The standard repository handle (e.g. 'pve-no-subscription', 'pve-enterprise').
- `node_name` (String) The name of the Proxmox VE node.

### Read-Only

- `description` (String) A description of the repository.
- `enabled` (Boolean) Whether the repository is currently enabled.
- `id` (String) Unique identifier for this data source (node/handle).
- `name` (String) The repository name / handle alias.
