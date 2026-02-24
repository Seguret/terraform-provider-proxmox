# proxmox_virtual_environment_node_journal (Data Source)

Retrieves the systemd journal entries from a Proxmox VE node.

## Example Usage

### Fetching the journal entries for a Proxmox VE node named "pve"
```hcl
data "proxmox_virtual_environment_node_journal" "pve" {
  node_name = "pve"
}

output "cpu_model" {
  value = data.proxmox_virtual_environment_node_journal.pve
}
```

## Schema

### Required

- `node_name` (String) The name of the Proxmox VE node.

### Read-Only

- `entries` (List of String) The journal log entries.
- `id` (String) Placeholder identifier.
