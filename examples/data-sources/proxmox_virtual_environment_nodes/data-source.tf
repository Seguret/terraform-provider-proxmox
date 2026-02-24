# List all nodes in the Proxmox VE cluster
data "proxmox_virtual_environment_nodes" "all" {}

output "node_names" {
  value = data.proxmox_virtual_environment_nodes.all.names
}

output "node_online" {
  value = data.proxmox_virtual_environment_nodes.all.online
}
