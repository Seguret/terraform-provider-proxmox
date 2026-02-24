---
page_title: "proxmox_virtual_environment_vm Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE virtual machine (QEMU/KVM).
---

# proxmox_virtual_environment_vm (Resource)

Manages a Proxmox VE virtual machine (QEMU/KVM).

## Example Usage

### Basic Ubuntu VM

```terraform
resource "proxmox_virtual_environment_vm" "ubuntu_vm" {
  name       = "ubuntu-server"
  node_name  = "pve"
  vmid       = 100
  memory     = 2048
  cpu_cores  = 2
  cpu_sockets = 1
  
  scsi0 = "local-lvm:32"
  net0  = "virtio,bridge=vmbr0"
  
  started = true
  on_boot = true
}
```

### Ubuntu VM with Cloud-Init

```terraform
resource "proxmox_virtual_environment_vm" "cloudinit_vm" {
  name        = "ubuntu-cloudinit"
  node_name   = "pve"
  vmid        = 101
  memory      = 2048
  cpu_cores   = 2
  cpu_sockets = 1
  description = "Ubuntu VM with Cloud-Init"
  tags        = "ubuntu;production"
  
  scsi0       = "local-lvm:32"
  net0        = "virtio,bridge=vmbr0"
  boot        = "order=scsi0"
  
  ci_user     = "ubuntu"
  ci_password = "MySecurePassword123!"
  ci_type     = "configdrive2"
  ipconfig0   = "ip=dhcp"
  
  ssh_keys    = file("${path.module}/authorized_keys")
  
  agent   = true
  started = true
}
```

### VM Cloned from Template

```terraform
resource "proxmox_virtual_environment_vm" "from_template" {
  name         = "web-server"
  node_name    = "pve"
  vmid         = 102
  clone_vmid   = 999  # Template VM ID
  full_clone   = true
  
  memory      = 4096
  cpu_cores   = 4
  cpu_sockets = 1
  
  on_boot = true
  started = true
}
```

### Windows VM with EFI

```terraform
resource "proxmox_virtual_environment_vm" "windows_vm" {
  name        = "windows-server"
  node_name   = "pve"
  vmid        = 103
  memory      = 8192
  cpu_cores   = 4
  cpu_sockets = 2
  
  bios    = "ovmf"
  machine = "q35"
  os_type = "win10"
  
  efidisk0 = "local-lvm:1"
  scsi0    = "local-lvm:50"
  scsi1    = "local-lvm:32"  # Secondary disk
  
  ide2  = "local:iso/windows-10.iso,media=cdrom"
  
  net0  = "virtio,bridge=vmbr0"
  vga   = "std"
  
  started = true
}
```

### Advanced Linux VM with Multiple Disks and Network

```terraform
resource "proxmox_virtual_environment_vm" "advanced_vm" {
  name         = "app-server"
  node_name    = "pve"
  vmid         = 104
  description  = "Advanced application server"
  tags         = "app;production;web"
  pool         = "production"
  
  # CPU and Memory
  memory       = 8192
  cpu_cores    = 8
  cpu_sockets  = 2
  cpu_type     = "host"
  balloon      = 2048
  
  # Storage
  scsi0        = "local-lvm:100"  # Root disk
  scsi1        = "local-lvm:200"  # Data disk
  scsi2        = "fast-ssd:150"   # Cache disk
  scsi_hw      = "virtio-scsi-pci"
  
  # Boot configuration
  boot         = "order=scsi0;net0"
  
  # Network interfaces
  net0         = "virtio,bridge=vmbr0"
  net1         = "virtio,bridge=vmbr1"
  net2         = "virtio,bridge=vmbr100,tag=200"
  
  ipconfig0    = "ip=192.168.1.100/24,gw=192.168.1.1"
  ipconfig1    = "ip=10.0.0.50/24"
  ipconfig2    = "ip=dhcp"
  
  # Serial port
  serial0      = "socket"
  
  # VGA
  vga          = "virtio"
  
  # Features
  agent        = true
  on_boot      = true
  protection   = true
  
  # Start configuration
  started      = true
}
```

### Snapshot Management

```terraform
# Create VM
resource "proxmox_virtual_environment_vm" "demo" {
  name      = "demo-vm"
  node_name = "pve"
  vmid      = 105
  memory    = 2048
  cpu_cores = 2
  
  scsi0 = "local-lvm:32"
  net0  = "virtio,bridge=vmbr0"
  
  started = true
}

# Create snapshot (requires VM to be running or stopped)
resource "proxmox_virtual_environment_vm_snapshot" "backup" {
  vm_id       = proxmox_virtual_environment_vm.demo.vmid
  node_name   = proxmox_virtual_environment_vm.demo.node_name
  snapshot_id = "pre-deployment"
  description = "Backup before deployment"
}
```


## Schema

### Required

- `node_name` (String) The node on which to create the VM.

### Optional

- `agent` (Boolean) Whether to enable the QEMU Guest Agent.
- `balloon` (Number) Balloon memory minimum in MiB. 0 to disable ballooning.
- `bios` (String) The BIOS type (seabios or ovmf).
- `boot` (String) Boot order (e.g., 'order=scsi0;ide2;net0').
- `ci_password` (String, Sensitive) Cloud-init password.
- `ci_type` (String) Cloud-init type (configdrive2 or nocloud).
- `ci_user` (String) Cloud-init user.
- `clone_vmid` (Number) VMID of the template/VM to clone from. If set, the VM is created as a clone.
- `cpu_cores` (Number) Number of CPU cores per socket.
- `cpu_sockets` (Number) Number of CPU sockets.
- `cpu_type` (String) The CPU type (e.g., 'host', 'kvm64', 'x86-64-v2-AES').
- `description` (String) The VM description.
- `efidisk0` (String) EFI disk configuration string.
- `full_clone` (Boolean) Whether to do a full clone (true) or linked clone (false).
- `ide0` (String) IDE disk 0 configuration string.
- `ide2` (String) IDE disk 2 configuration string (often used for cloud-init).
- `ipconfig0` (String) IP configuration for net0 (e.g., 'ip=dhcp' or 'ip=10.0.0.2/24,gw=10.0.0.1').
- `ipconfig1` (String) IP configuration for net1.
- `machine` (String) The machine type (e.g., q35, i440fx, or a specific version).
- `memory` (Number) Memory in MiB.
- `name` (String) The VM name.
- `nameserver` (String) Cloud-init DNS nameserver.
- `net0` (String) Network device 0 configuration (e.g., 'virtio=XX:XX:XX:XX:XX:XX,bridge=vmbr0').
- `net1` (String) Network device 1 configuration.
- `net2` (String) Network device 2 configuration.
- `net3` (String) Network device 3 configuration.
- `on_boot` (Boolean) Whether to start the VM on host boot.
- `os_type` (String) The OS type (l26, l24, win11, win10, win7, solaris, other).
- `pool` (String) The resource pool to add the VM to.
- `protection` (Boolean) Whether the VM is protected from removal.
- `scsi0` (String) SCSI disk 0 configuration string.
- `scsi1` (String) SCSI disk 1 configuration string.
- `scsi2` (String) SCSI disk 2 configuration string.
- `scsi3` (String) SCSI disk 3 configuration string.
- `scsi_hw` (String) The SCSI controller type (virtio-scsi-pci, virtio-scsi-single, lsi, megasas, pvscsi).
- `searchdomain` (String) Cloud-init DNS search domain.
- `serial0` (String) Serial device 0 (e.g., 'socket').
- `ssh_keys` (String) Cloud-init SSH public keys (URL-encoded, newline-separated).
- `started` (Boolean) Whether the VM should be started after creation.
- `tags` (String) Tags for the VM (semicolon-separated).
- `tpmstate0` (String) TPM state configuration string.
- `vga` (String) VGA configuration string (e.g., 'std', 'virtio', 'serial0').
- `virtio0` (String) VirtIO disk 0 configuration string.
- `virtio1` (String) VirtIO disk 1 configuration string.
- `vmid` (Number) The VM ID. If not set, the next available ID will be used.

### Read-Only

- `id` (String) The ID of this resource.
- `status` (String) The current VM status (running, stopped, etc.).
- `template` (Boolean) Whether the VM is a template.
