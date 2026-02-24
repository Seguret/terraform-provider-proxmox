# proxmox Provider

Terraform provider for managing Proxmox Virtual Environment resources.

## Example Usage

```terraform
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
```


## Schema

### Optional

- `api_token` (String, Sensitive) The Proxmox VE API token (e.g., user@pam!tokenid=uuid). Can also be set with the PROXMOX_VE_API_TOKEN environment variable.
- `endpoint` (String) The Proxmox VE API endpoint URL (e.g., https://pve.example.com:8006). Can also be set with the PROXMOX_VE_ENDPOINT environment variable.
- `insecure` (Boolean) Whether to skip TLS certificate verification. Defaults to false. Can also be set with the PROXMOX_VE_INSECURE environment variable.
- `password` (String, Sensitive) The Proxmox VE password for ticket-based auth. Can also be set with the PROXMOX_VE_PASSWORD environment variable.
- `username` (String) The Proxmox VE username for ticket-based auth (e.g., root@pam). Can also be set with the PROXMOX_VE_USERNAME environment variable.
