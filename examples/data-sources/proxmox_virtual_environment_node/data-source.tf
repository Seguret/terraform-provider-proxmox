# Get detailed status of a specific node
data "proxmox_virtual_environment_node" "pve1" {
  node_name = "pve1"
}

output "cpu_model" {
  value = data.proxmox_virtual_environment_node.pve1.cpu_model
}

output "memory_total_gb" {
  value = data.proxmox_virtual_environment_node.pve1.memory_total / 1073741824
}

output "uptime_hours" {
  value = data.proxmox_virtual_environment_node.pve1.uptime / 3600
}

output "pve_version" {
  value = data.proxmox_virtual_environment_node.pve1.pve_version
}
