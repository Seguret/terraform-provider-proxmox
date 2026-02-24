# proxmox_virtual_environment_node_scan (Data Source)

Scans for available storage of a given type on a Proxmox VE node.




## Example Usage

```hcl
data "proxmox_virtual_environment_node_scan" "nfs" {
  node_name = "pve"
  scan_type = "nfs"
  server    = "nas.example.com"
}
```

## Schema

### Required

- `node_name` (String) The name of the node.
- `scan_type` (String) The type of storage to scan (`nfs`, `iscsi`, `cifs`, `lvm`, `lvmthin`, `zfs`, `pbs`).

### Optional

- `portal` (String) iSCSI portal address (for `iscsi` scan type).
- `server` (String) Server address (for `nfs`, `cifs`, `pbs` scan types).
- `vg` (String) LVM volume group name (for `lvmthin` scan type).

### Read-Only

- `id` (String) The datasource identifier.
- `keys` (List of String) Result field keys (flattened from each result map).
- `values` (List of String) Result field values corresponding to each key.
