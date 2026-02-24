# proxmox_virtual_environment_storage_content (Data Source)

Retrieves the list of files stored in a Proxmox VE node storage.

## Example Usage

```hcl
data "proxmox_virtual_environment_storage_content" "test" {
  node_name = "pve"
  storage = "local"
}

output "test" {
  value = data.proxmox_virtual_environment_storage_content.test
} 

output "volids" {
  value = data.proxmox_virtual_environment_storage_content.test.volids
} 
```

## Schema

### Required

- `node_name` (String) The name of the Proxmox VE node.
- `storage` (String) The storage name (e.g., 'local', 'local-lvm').

### Optional

- `content_type` (String) Optional filter for content type (e.g., 'iso', 'vztmpl', 'backup').

### Read-Only

- `content` (List of String) The content type for each file.
- `ctime` (List of Number) The creation time (unix timestamp) for each file (if provided).
- `format` (List of String) The format for each file (if provided by Proxmox).
- `id` (String) Placeholder identifier.
- `notes` (List of String) Notes for each file (if provided).
- `size` (List of Number) The size in bytes for each file.
- `used` (List of Number) The used size in bytes for each file (if provided).
- `volids` (List of String) The volume IDs for each file.
