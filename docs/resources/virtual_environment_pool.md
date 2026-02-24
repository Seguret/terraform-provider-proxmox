---
page_title: "proxmox_virtual_environment_pool Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE resource pool.
---

# proxmox_virtual_environment_pool (Resource)

Manages a Proxmox VE resource pool.

## Example Usage

### Basic Pool

```terraform
resource "proxmox_virtual_environment_pool" "production" {
  pool_id = "production"
  comment = "Production environment resources"
}
```

### Multiple Pools

```terraform
resource "proxmox_virtual_environment_pool" "development" {
  pool_id = "dev"
  comment = "Development and testing environments"
}

resource "proxmox_virtual_environment_pool" "staging" {
  pool_id = "staging"
  comment = "Staging environment for pre-production testing"
}

resource "proxmox_virtual_environment_pool" "internal" {
  pool_id = "internal"
  comment = "Internal tools and monitoring"
}
```

### Pool with Resource Assignment

```terraform
resource "proxmox_virtual_environment_pool" "app_pool" {
  pool_id = "app-services"
  comment = "Pool for application services"
}

resource "proxmox_virtual_environment_vm" "web_server" {
  name      = "web-01"
  node_name = "pve"
  vmid      = 100
  pool      = proxmox_virtual_environment_pool.app_pool.pool_id
  
  memory    = 2048
  cpu_cores = 2
  
  scsi0 = "local-lvm:32"
  net0  = "virtio,bridge=vmbr0"
  
  started = true
}

resource "proxmox_virtual_environment_container" "app_container" {
  vmid      = 200
  node_name = "pve"
  hostname  = "app-ct"
  pool      = proxmox_virtual_environment_pool.app_pool.pool_id
  
  os_template = "local:vztmpl/ubuntu-22.04-standard_22.04-1_amd64.tar.zst"
  
  memory  = 1024
  cores   = 2
  
  rootfs  = "local-lvm:8"
  net0    = "name=eth0,bridge=vmbr0,ip=dhcp"
  
  started = true
}
```



## Schema

### Required

- `pool_id` (String) The pool identifier.

### Optional

- `comment` (String) A comment for the pool.

### Read-Only

- `id` (String) The ID of this resource.
