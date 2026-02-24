terraform {
  required_providers {
    proxmox = {
      source  = "registry.terraform.io/Seguret/proxmox"
      version = "~> 0.1.0"
    }
  }
}

# Option 1: API Token (recommended)
provider "proxmox" {
  endpoint = "https://pve.example.com:8006"
  api_token = "user@pam!mytoken=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  insecure = true # Set to false in production with valid TLS certs
}

# Option 2: Username/Password (less recommended)
# provider "proxmox" {
#   endpoint = "https://pve.example.com:8006"
#   username = "root@pam"
#   password = "secret"
#   insecure = true
# }

# Smoke test: read the Proxmox VE version
data "proxmox_virtual_environment_version" "example" {}

output "proxmox_version" {
  value = data.proxmox_virtual_environment_version.example.version
}
