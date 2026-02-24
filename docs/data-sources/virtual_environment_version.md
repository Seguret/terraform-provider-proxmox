# proxmox_virtual_environment_version (Data Source)

Retrieves the Proxmox VE version information.

## Example Usage

### Basic Example
```hcl
data "proxmox_virtual_environment_version" "vms" {}

output "vm_list" {
  value       = data.proxmox_virtual_environment_version.vms
  description = "The Proxmox VE version information."
}
```

### Output Example with Specific Attribute
```hcl
data "proxmox_virtual_environment_version" "vms" {}

output "vm_list" {
  value       = data.proxmox_virtual_environment_version.vms.release
  description = "The Proxmox VE release string (e.g., '8.1')."
}
```

## Schema

### Read-Only

- `id` (String) Placeholder identifier.
- `release` (String) The Proxmox VE release string (e.g., '8.1').
- `repo_id` (String) The Proxmox VE repository ID.
- `version` (String) The Proxmox VE version string (e.g., '8.1.3').
