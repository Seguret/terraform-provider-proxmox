# proxmox_virtual_environment_datastores (Data Source)

Retrieves the list of datastores (storage) available on a Proxmox VE node.

## Example Usage

```terraform
# List all storage available on a specific node
data "proxmox_virtual_environment_datastores" "pve1" {
  node_name = "pve1"
}

output "storage_names" {
  value = data.proxmox_virtual_environment_datastores.pve1.names
}

output "storage_types" {
  value = data.proxmox_virtual_environment_datastores.pve1.types
}

output "storage_available_gb" {
  value = [for v in data.proxmox_virtual_environment_datastores.pve1.available : v / 1073741824]
}
```


## Schema

### Required

- `node_name` (String) The name of the Proxmox VE node.

### Read-Only

- `active` (List of Boolean) Whether each storage is active.
- `available` (List of Number) Available size in bytes for each storage.
- `content_types` (List of String) The content types supported by each storage (comma-separated, e.g., 'images,rootdir,vztmpl,iso,backup').
- `enabled` (List of Boolean) Whether each storage is enabled.
- `id` (String) Placeholder identifier.
- `names` (List of String) The storage names.
- `shared` (List of Boolean) Whether each storage is shared across nodes.
- `total` (List of Number) Total size in bytes for each storage.
- `types` (List of String) The storage types (e.g., 'dir', 'lvm', 'zfspool', 'nfs', 'cifs').
- `used` (List of Number) Used size in bytes for each storage.
- `used_fraction` (List of Number) The used fraction (0.0-1.0) for each storage.
