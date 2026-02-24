---
page_title: "proxmox_virtual_environment_container Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE LXC container.
---

# proxmox_virtual_environment_container (Resource)

Manages a Proxmox VE LXC container.

## Example Usage

### Basic Ubuntu Container

```terraform
resource "proxmox_virtual_environment_container" "ubuntu_container" {
  vmid         = 200
  node_name    = "pve"
  hostname     = "ubuntu-ct"
  os_template  = "local:vztmpl/ubuntu-22.04-standard_22.04-1_amd64.tar.zst"
  
  memory  = 1024
  cores   = 2
  
  rootfs  = "local-lvm:8"
  net0    = "name=eth0,bridge=vmbr0,ip=dhcp"
  
  started = true
  on_boot = true
}
```

### Unprivileged Container with SSH

```terraform
resource "proxmox_virtual_environment_container" "app_container" {
  vmid            = 201
  node_name       = "pve"
  hostname        = "app-server"
  description     = "Application container"
  tags            = "app;production"
  os_template     = "local:vztmpl/debian-12-standard_12.0-1_amd64.tar.zst"
  
  memory          = 2048
  swap            = 512
  cpu_cores       = 4
  cpu_units       = 1024
  
  unprivileged    = true
  
  rootfs          = "local-lvm:16"
  net0            = "name=eth0,bridge=vmbr0,ip=192.168.1.50/24,gw=192.168.1.1"
  nameserver      = "8.8.8.8 8.8.4.4"
  searchdomain    = "example.com"
  
  ssh_keys        = <<-EOT
    ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC... user@host
  EOT
  
  password        = "SecurePassword123!"
  
  on_boot         = true
  started         = true
}
```

### Container with Multiple Mount Points

```terraform
resource "proxmox_virtual_environment_container" "storage_container" {
  vmid           = 202
  node_name      = "pve"
  hostname       = "storage-ct"
  os_template    = "local:vztmpl/ubuntu-22.04-standard_22.04-1_amd64.tar.zst"
  
  memory  = 1024
  cores   = 2
  
  # Root filesystem
  rootfs  = "local-lvm:10"
  
  # Additional mount points
  mp0     = "fast-ssd:50,mp=/data"
  mp1     = "backup-storage:100,mp=/backups"
  mp2     = "local-lvm:20,mp=/var/log"
  
  # Network
  net0    = "name=eth0,bridge=vmbr0,ip=dhcp"
  
  started = true
}
```

### Privileged Container (Docker Host)

```terraform
resource "proxmox_virtual_environment_container" "docker_container" {
  vmid            = 203
  node_name       = "pve"
  hostname        = "docker-host"
  description     = "Docker host container"
  os_template     = "local:vztmpl/ubuntu-22.04-standard_22.04-1_amd64.tar.zst"
  
  memory          = 4096
  cpu_cores       = 4
  
  unprivileged    = false  # Privileged container needed for Docker
  console         = true
  
  features        = "nesting=1,keyctl=1"
  
  rootfs          = "local-lvm:50"
  
  net0            = "name=eth0,bridge=vmbr0,ip=dhcp"
  
  password        = "MySecurePass123!"
  
  on_boot         = true
  started         = true
}
```

### Container with Multiple Networks

```terraform
resource "proxmox_virtual_environment_container" "multinetwork_container" {
  vmid            = 204
  node_name       = "pve"
  hostname        = "router-ct"
  description     = "Multi-network container"
  os_template     = "local:vztmpl/debian-12-standard_12.0-1_amd64.tar.zst"
  
  memory  = 2048
  cores   = 4
  
  unprivileged = false
  
  rootfs  = "local-lvm:16"
  
  # Primary network (management)
  net0    = "name=eth0,bridge=vmbr0,ip=192.168.1.100/24,gw=192.168.1.1"
  
  # Secondary network (production)
  net1    = "name=eth1,bridge=vmbr1,ip=10.0.0.100/24"
  
  # Tertiary network (isolated)
  net2    = "name=eth2,bridge=vmbr2,ip=172.16.0.100/24"
  
  nameserver   = "8.8.8.8"
  searchdomain = "example.com"
  
  password = "RoutingPassword123!"
  
  on_boot = true
  started = true
}
```

### Container Cloned from Existing

```terraform
resource "proxmox_virtual_environment_container" "cloned_container" {
  vmid      = 205
  node_name = "pve"
  
  clone_vmid = 200  # Clone from existing container 200
  full_clone = true
  
  hostname = "cloned-container"
  
  memory = 1024
  cores  = 2
  
  net0   = "name=eth0,bridge=vmbr0,ip=dhcp"
  
  started = true
}
```

### Container as a Template

```terraform
resource "proxmox_virtual_environment_container" "template_container" {
  vmid           = 999
  node_name      = "pve"
  hostname       = "ubuntu-template"
  description    = "Ubuntu template for cloning"
  os_template    = "local:vztmpl/ubuntu-22.04-standard_22.04-1_amd64.tar.zst"
  
  memory = 1024
  cores  = 2
  
  rootfs = "local-lvm:10"
  net0   = "name=eth0,bridge=vmbr0"
  
  # Container will be converted to template after creation
  started = false
  
  protection = true
}

# Create a VM from this container template
resource "proxmox_virtual_environment_container" "from_template" {
  vmid           = 206
  node_name      = "pve"
  hostname       = "production-ct"
  clone_vmid     = proxmox_virtual_environment_container.template_container.vmid
  full_clone     = true
  
  memory = 2048
  cores  = 4
  
  net0   = "name=eth0,bridge=vmbr0,ip=dhcp"
  
  started = true
}
```


## Schema

### Required

- `node_name` (String) The node on which to create the container.

### Optional

- `clone_vmid` (Number) VMID of the container to clone from.
- `console` (Boolean) Whether to attach a console device.
- `cpu_cores` (Number) Number of CPU cores.
- `cpu_limit` (Number) CPU usage limit (0 = unlimited).
- `cpu_units` (Number) CPU weight for a container (relative weight vs other containers).
- `description` (String) The container description.
- `features` (String) Container feature flags (e.g., 'nesting=1,keyctl=1').
- `full_clone` (Boolean) Whether to do a full clone.
- `hostname` (String) The container hostname.
- `memory` (Number) Memory in MiB.
- `mp0` (String) Mount point 0.
- `mp1` (String) Mount point 1.
- `mp2` (String) Mount point 2.
- `nameserver` (String) DNS nameserver.
- `net0` (String) Network interface 0 (e.g., 'name=eth0,bridge=vmbr0,ip=dhcp').
- `net1` (String) Network interface 1.
- `net2` (String) Network interface 2.
- `net3` (String) Network interface 3.
- `on_boot` (Boolean) Whether to start the container on host boot.
- `os_template` (String) The OS template to use (e.g., 'local:vztmpl/ubuntu-22.04-standard_22.04-1_amd64.tar.zst').
- `os_type` (String) The OS type (debian, ubuntu, centos, fedora, opensuse, archlinux, alpine, gentoo, nixos, unmanaged).
- `password` (String, Sensitive) Root password for the container.
- `pool` (String) The resource pool to add the container to.
- `protection` (Boolean) Whether the container is protected from removal.
- `rootfs` (String) Root filesystem configuration (e.g., 'local-lvm:8').
- `searchdomain` (String) DNS search domain.
- `ssh_keys` (String) SSH public keys.
- `started` (Boolean) Whether the container should be started after creation.
- `swap` (Number) Swap in MiB.
- `tags` (String) Tags for the container (semicolon-separated).
- `tty` (Number) Number of TTY devices (0-6).
- `unprivileged` (Boolean) Whether to create an unprivileged container.
- `vmid` (Number) The container ID. If not set, the next available ID will be used.

### Read-Only

- `id` (String) The ID of this resource.
- `status` (String) The container current status.
- `template` (Boolean) Whether the container is a template.
