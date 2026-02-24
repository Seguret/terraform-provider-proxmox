# proxmox_virtual_environment_file (Data Source)

Retrieves information about a specific file stored in Proxmox VE storage.




## Example Usage

```hcl
data "proxmox_virtual_environment_file" "example" {
  node_name    = "pve"
  datastore_id = "local"
  volume_id    = "local:iso/ubuntu-22.04.iso"
}
```

## Schema

### Required

- `datastore_id` (String) The storage ID (e.g. `local`).
- `node_name` (String) The name of the Proxmox VE node.
- `volume_id` (String) The volume ID to look up (e.g. `local:iso/ubuntu.iso`).

### Read-Only

- `content_type` (String) The content type of the file (e.g. `iso`, `vztmpl`).
- `file_name` (String) The filename portion of the volume ID.
- `id` (String) The volume ID of the file (same as `volume_id`).
- `size` (Number) The file size in bytes as reported by Proxmox.
