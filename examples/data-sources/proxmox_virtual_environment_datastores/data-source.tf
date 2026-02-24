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
