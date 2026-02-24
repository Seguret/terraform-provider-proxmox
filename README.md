# Terraform Provider for Proxmox VE

![Go Version](https://img.shields.io/badge/go-1.24%2B-blue)
![Terraform](https://img.shields.io/badge/terraform-1.x-623CE4)
![Plugin Framework](https://img.shields.io/badge/plugin--framework-1.17.0-623CE4)
![License](https://img.shields.io/badge/license-MIT-green)

---

## Summary

Terraform provider for comprehensive resource management of the [Proxmox Virtual Environment](https://www.proxmox.com/en/proxmox-virtual-environment/overview).Built with the [Terraform Plugin Framework v1.17](https://github.com/hashicorp/terraform-plugin-framework), it exposes over 40 resources and 20 data sources covering virtual machines, LXC containers, storage, firewall, networking, SDN, high availability, replication, certificates, ACME, and more.

**Quick setup:**

```hcl
terraform {
  required_providers {
    proxmox = {
      source  = "registry.terraform.io/Seguret/proxmox"
      version = "~> 0.1"
    }
  }
}

provider "proxmox" {
  endpoint  = "https://pve.example.com:8006"
  api_token = "terraform@pam!deploy=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

---

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Provider Configuration](#provider-configuration)
- [Authentication](#authentication)
- [Environment Variables](#environment-variables)
- [Resources](#resources)
  - [Access Management](#access-management)
  - [Compute](#compute)
  - [Storage](#storage)
  - [Networking](#networking)
  - [Firewall](#firewall)
  - [High Availability](#high-availability)
  - [Node Management](#node-management)
  - [Cluster](#cluster)
  - [Replication](#replication)
  - [SDN](#sdn)
  - [ACME](#acme)
  - [Metrics](#metrics)
  - [Hardware Mappings](#hardware-mappings)
  - [APT](#apt)
  - [Backup](#backup)
- [Data Sources](#data-sources)
- [Resource Examples](#resource-examples)
- [Import Reference](#import-reference)
- [Development](#development)
- [License](#license)

---

## Features

### Resources

| Category | Resource Type | Description |
|---|---|---|
| **Access Management** | `proxmox_virtual_environment_user` | Proxmox VE user account |
| | `proxmox_virtual_environment_group` | User group |
| | `proxmox_virtual_environment_role` | Custom role with privileges |
| | `proxmox_virtual_environment_acl` | Access control list entry |
| | `proxmox_virtual_environment_pool` | Resource pool |
| | `proxmox_virtual_environment_realm` | Authentication realm (PAM, LDAP, AD, OpenID) |
| | `proxmox_virtual_environment_user_token` | User API token |
| **Compute** | `proxmox_virtual_environment_vm` | QEMU/KVM virtual machine |
| | `proxmox_virtual_environment_container` | LXC container |
| | `proxmox_virtual_environment_vm_snapshot` | VM snapshot |
| | `proxmox_virtual_environment_container_snapshot` | Container snapshot |
| **Storage** | `proxmox_virtual_environment_storage` | Cluster-level storage definition |
| | `proxmox_virtual_environment_download_file` | Download ISO or template to node storage |
| **Networking** | `proxmox_virtual_environment_network_interface` | Node network interface (bridge, bond, VLAN) |
| **Firewall** | `proxmox_virtual_environment_firewall_rule` | Firewall rule (cluster/node/VM/CT scope) |
| | `proxmox_virtual_environment_firewall_options` | Firewall options per scope |
| | `proxmox_virtual_environment_firewall_ipset` | IP set with CIDR entries |
| | `proxmox_virtual_environment_firewall_alias` | Named IP/CIDR alias |
| | `proxmox_virtual_environment_firewall_security_group` | Cluster firewall security group |
| | `proxmox_virtual_environment_firewall_security_group_rule` | Rule inside a security group |
| **High Availability** | `proxmox_virtual_environment_ha_resource` | HA resource (VM or CT) |
| | `proxmox_virtual_environment_ha_group` | HA group |
| **Node Management** | `proxmox_virtual_environment_dns` | Node DNS configuration |
| | `proxmox_virtual_environment_hosts` | Node /etc/hosts file |
| | `proxmox_virtual_environment_certificate` | Custom TLS certificate on a node |
| | `proxmox_virtual_environment_time` | Node timezone |
| **Cluster** | `proxmox_virtual_environment_cluster_options` | Cluster-wide options (singleton) |
| **Replication** | `proxmox_virtual_environment_replication` | VM/CT replication job |
| **Backup** | `proxmox_virtual_environment_backup` | vzdump backup schedule |
| **SDN** | `proxmox_virtual_environment_sdn_zone` | SDN zone |
| | `proxmox_virtual_environment_sdn_vnet` | SDN VNet |
| | `proxmox_virtual_environment_sdn_subnet` | SDN subnet within a VNet |
| **ACME** | `proxmox_virtual_environment_acme_account` | ACME account (Let's Encrypt) |
| | `proxmox_virtual_environment_acme_plugin` | ACME DNS challenge plugin |
| **Metrics** | `proxmox_virtual_environment_metrics_server` | External metrics server (Graphite/InfluxDB) |
| **Hardware Mappings** | `proxmox_virtual_environment_hardware_mapping_pci` | Cluster PCI hardware mapping |
| | `proxmox_virtual_environment_hardware_mapping_usb` | Cluster USB hardware mapping |
| **APT** | `proxmox_virtual_environment_apt_repository` | Standard APT repository on a node |

### Data Sources

| Data Source Type | Description |
|---|---|
| `proxmox_virtual_environment_version` | Proxmox VE API version info |
| `proxmox_virtual_environment_nodes` | List of all cluster nodes |
| `proxmox_virtual_environment_node` | Detailed status of a specific node |
| `proxmox_virtual_environment_datastores` | Storage visible on a node |
| `proxmox_virtual_environment_vms` | Summary list of VMs on a node |
| `proxmox_virtual_environment_vm` | Details of a specific VM |
| `proxmox_virtual_environment_containers` | Summary list of LXC containers on a node |
| `proxmox_virtual_environment_container` | Details of a specific container |
| `proxmox_virtual_environment_vm_snapshots` | List of snapshots for a VM |
| `proxmox_virtual_environment_container_snapshots` | List of snapshots for a container |
| `proxmox_virtual_environment_ha_resources` | List of all HA resources |
| `proxmox_virtual_environment_ha_groups` | List of all HA groups |
| `proxmox_virtual_environment_users` | List of all users |
| `proxmox_virtual_environment_groups` | List of all groups |
| `proxmox_virtual_environment_roles` | List of all roles |
| `proxmox_virtual_environment_pools` | List of all resource pools |
| `proxmox_virtual_environment_network_interfaces` | Network interfaces on a node |
| `proxmox_virtual_environment_hardware_mapping_pci_list` | List of all PCI hardware mappings |
| `proxmox_virtual_environment_hardware_mapping_usb_list` | List of all USB hardware mappings |
| `proxmox_virtual_environment_sdn_zones` | List of SDN zones |
| `proxmox_virtual_environment_sdn_vnets` | List of SDN VNets |
| `proxmox_virtual_environment_apt_repositories` | APT repositories on a node |
| `proxmox_virtual_environment_backups` | List of backup jobs |

---

## Requirements

| Component | Minimum Version |
|---|---|
| [Go](https://go.dev/dl/) | 1.24 |
| [Terraform](https://developer.hashicorp.com/terraform/downloads) | 1.x |
| Proxmox VE | 7.x or 8.x |
| GNU Make | any |

---

## Provider Configuration

```hcl
terraform {
  required_providers {
    proxmox = {
      source  = "registry.terraform.io/Seguret/proxmox"
      version = "~> 0.1"
    }
  }
}

# API Token authentication (recommended)
provider "proxmox" {
  endpoint  = "https://pve.example.com:8006"
  api_token = "terraform@pam!deploy=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  insecure  = false
}

# Username/password authentication (ticket-based)
provider "proxmox" {
  endpoint = "https://pve.example.com:8006"
  username = "root@pam"
  password = "your_password"
  insecure = false
}
```

### Provider Attributes

| Attribute | Type | Required | Description |
|---|---|---|---|
| `endpoint` | string | Yes | Base URL of the Proxmox VE API, e.g. `https://pve.example.com:8006`. Trailing slashes are stripped automatically. |
| `api_token` | string | One of | API token string in the format `user@realm!tokenid=uuid`. Marked as sensitive. |
| `username` | string | One of | Username for ticket-based auth, e.g. `root@pam`. |
| `password` | string | One of | Password for ticket-based auth. Marked as sensitive. |
| `insecure` | boolean | No | Skip TLS certificate verification. Default: `false`. Use only in lab environments. |

Either `api_token` or both `username` and `password` must be provided (via attributes or corresponding environment variables).

---

## Authentication

### API Token (Recommended)

Create an API token in the Proxmox web UI under **Datacenter > Permissions > API Tokens**. The format is:

```
<user>@<realm>!<tokenid>=<uuid>
```

Example: `terraform@pam!deploy=6a4f8b2c-1234-5678-abcd-ef0123456789`

Provide the token via the `api_token` attribute or the `PROXMOX_VE_API_TOKEN` environment variable.

**Minimum privileges required for the token role:**

| Operation | Required Privileges |
|---|---|
| VM management | `VM.Allocate`, `VM.Config.*`, `VM.PowerMgmt` |
| Container management | `VM.Allocate`, `VM.Config.*`, `VM.PowerMgmt` |
| Storage | `Datastore.Allocate`, `Datastore.AllocateSpace` |
| Access control | `Sys.Audit`, `User.Modify` |
| Pool | `Pool.Allocate` |
| Networking | `Sys.Modify` |
| Firewall | `Sys.Modify` |
| HA | `Sys.Modify` |
| SDN | `SDN.Allocate` |

### Username and Password (Ticket-based)

The provider calls `POST /api2/json/access/ticket` with the supplied credentials and stores the resulting ticket and CSRF token for the duration of the session. Tickets are not automatically renewed during long-running applies.

> **Note:** For production use, API tokens are strongly preferred over username/password because they support privilege separation, can be revoked individually, and do not require transmitting a password.

---

## Environment Variables

All provider attributes can be set via environment variables. Explicit values in the `provider` block take precedence over environment variables.

| Environment Variable | Corresponding Attribute | Description |
|---|---|---|
| `PROXMOX_VE_ENDPOINT` | `endpoint` | Proxmox VE API URL |
| `PROXMOX_VE_API_TOKEN` | `api_token` | API token string |
| `PROXMOX_VE_USERNAME` | `username` | Username for ticket auth |
| `PROXMOX_VE_PASSWORD` | `password` | Password for ticket auth |
| `PROXMOX_VE_INSECURE` | `insecure` | Set to `true`, `1`, or `yes` to disable TLS verification |

Example using environment variables:

```shell
export PROXMOX_VE_ENDPOINT="https://pve.example.com:8006"
export PROXMOX_VE_API_TOKEN="terraform@pam!deploy=6a4f8b2c-..."
terraform apply
```

---

## Resources

### Access Management

#### `proxmox_virtual_environment_user`

Manages a Proxmox VE user account.

| Attribute | Type | Default | Description |
|---|---|---|---|
| `user_id` | string | - | **Required.** User identifier (e.g. `user@pam` or `user@pve`). Forces replacement. |
| `password` | string | - | Password. Used only at creation for PVE realm users. Sensitive. |
| `email` | string | - | User email address. |
| `enabled` | boolean | `true` | Enable the user account. |
| `expire` | integer | `0` | Expiration date as Unix epoch. `0` means no expiration. |
| `first_name` | string | - | First name. |
| `last_name` | string | - | Last name. |
| `comment` | string | - | Comment. |
| `groups` | string | - | Comma-separated list of group memberships. |
| `keys` | string | - | Two-factor authentication keys. |

**Import:** `terraform import proxmox_virtual_environment_user.example user@pve`

---

#### `proxmox_virtual_environment_group`

Manages a Proxmox VE group.

| Attribute | Type | Description |
|---|---|---|
| `group_id` | string | **Required.** Group identifier. Forces replacement. |
| `comment` | string | Comment. |
| `members` | string | **Read-only.** Comma-separated user IDs in the group. |

**Import:** `terraform import proxmox_virtual_environment_group.example mygroup`

---

#### `proxmox_virtual_environment_role`

Manages a custom Proxmox VE role with a specific set of privileges.

| Attribute | Type | Description |
|---|---|---|
| `role_id` | string | **Required.** Role identifier. Forces replacement. |
| `privileges` | string | **Required.** Comma-separated privileges (e.g. `VM.Audit,VM.Console,VM.PowerMgmt`). |

**Import:** `terraform import proxmox_virtual_environment_role.example MyRole`

---

#### `proxmox_virtual_environment_acl`

Manages an Access Control List entry. Assigns a role to a user or group on a specific path.

| Attribute | Type | Default | Description |
|---|---|---|---|
| `path` | string | - | **Required.** ACL path (e.g. `/`, `/vms/100`, `/storage/local`). Forces replacement. |
| `role_id` | string | - | **Required.** Role to assign. Forces replacement. |
| `user_id` | string | - | User to assign the role to. Mutually exclusive with `group_id`. |
| `group_id` | string | - | Group to assign the role to. Mutually exclusive with `user_id`. |
| `propagate` | boolean | `true` | Propagate the ACL to child objects. |

> **Note:** ACL resources do not support import because the Proxmox API returns the full ACL list rather than individual entries.

---

#### `proxmox_virtual_environment_pool`

Manages a Proxmox VE resource pool.

| Attribute | Type | Description |
|---|---|---|
| `pool_id` | string | **Required.** Pool identifier. Forces replacement. |
| `comment` | string | Comment. |

**Import:** `terraform import proxmox_virtual_environment_pool.example production`

---

#### `proxmox_virtual_environment_realm`

Manages a Proxmox VE authentication realm. Supports PAM, PVE, Active Directory, LDAP, and OpenID Connect.

| Attribute | Type | Description |
|---|---|---|
| `realm` | string | **Required.** Realm identifier (e.g. `my-ldap`). Forces replacement. |
| `type` | string | **Required.** Realm type: `pam`, `pve`, `ad`, `ldap`, `openid`. Forces replacement. |
| `comment` | string | Comment. |
| `default` | boolean | Whether this is the default realm. |
| `server1` | string | Primary server address (for `ldap`/`ad`). |
| `server2` | string | Secondary server address (for `ldap`/`ad`). |
| `port` | integer | Server port (for `ldap`/`ad`). |
| `base_dn` | string | LDAP base distinguished name. |
| `bind_dn` | string | LDAP bind distinguished name. |
| `password` | string | LDAP bind password. Sensitive. |
| `user_attr` | string | LDAP user attribute name. |
| `secure` | boolean | Use LDAPS/TLS. |
| `verify` | boolean | Verify server TLS certificate. |
| `issuer_url` | string | OpenID Connect issuer URL. |
| `client_id` | string | OpenID Connect client ID. |
| `client_key` | string | OpenID Connect client secret. Sensitive. |
| `username_claim` | string | OpenID Connect claim used as username. |
| `auto_create` | boolean | Automatically create users on first login. |
| `tfa` | string | Two-factor authentication provider. |

**Import:** `terraform import proxmox_virtual_environment_realm.example my-ldap`

---

#### `proxmox_virtual_environment_user_token`

Manages a Proxmox VE user API token. The token secret is only available after creation and is stored in state.

| Attribute | Type | Default | Description |
|---|---|---|---|
| `user_id` | string | - | **Required.** User ID (e.g. `root@pam`). Forces replacement. |
| `token_id` | string | - | **Required.** Token name. Forces replacement. |
| `comment` | string | - | Comment. |
| `expire` | integer | `0` | Expiration as Unix timestamp. `0` means no expiration. |
| `privileges_separation` | boolean | `true` | Enable privilege separation (token cannot exceed user's privileges). |
| `value` | string | - | **Read-only, Sensitive.** The token secret (only available after creation). |

**Import:** `terraform import proxmox_virtual_environment_user_token.example root@pam/my-token`

---

### Compute

#### `proxmox_virtual_environment_vm`

Manages a QEMU/KVM virtual machine. Supports creation of new VMs and cloning from existing templates.

| Attribute | Type | Default | Description |
|---|---|---|---|
| `node_name` | string | - | **Required.** Proxmox node to create the VM on. |
| `vmid` | integer | auto | VM ID. If omitted, the next available ID is assigned automatically. |
| `name` | string | - | VM name. |
| `description` | string | - | VM description. |
| `tags` | string | - | Semicolon-separated tags (e.g. `web;prod`). |
| `on_boot` | boolean | `false` | Start VM on host boot. |
| `started` | boolean | `true` | Start VM after creation. |
| `protection` | boolean | `false` | Protect VM from deletion. |
| `agent` | boolean | `false` | Enable QEMU Guest Agent. |
| `os_type` | string | `l26` | OS type: `l26`, `l24`, `win11`, `win10`, `win7`, `solaris`, `other`. |
| `bios` | string | `seabios` | BIOS type: `seabios` or `ovmf` (UEFI). |
| `machine` | string | - | Machine type: `q35`, `i440fx`, or specific version. |
| `scsi_hw` | string | `virtio-scsi-pci` | SCSI controller: `virtio-scsi-pci`, `virtio-scsi-single`, `lsi`, `megasas`, `pvscsi`. |
| `boot` | string | - | Boot order (e.g. `order=scsi0;ide2;net0`). |
| `pool` | string | - | Resource pool to add the VM to. |
| `cpu_sockets` | integer | `1` | Number of CPU sockets. |
| `cpu_cores` | integer | `1` | Number of CPU cores per socket. |
| `cpu_type` | string | `kvm64` | CPU type: `host`, `kvm64`, `x86-64-v2-AES`, etc. |
| `memory` | integer | `512` | Memory in MiB. |
| `balloon` | integer | - | Minimum balloon memory in MiB. `0` disables ballooning. |
| `vga` | string | - | VGA configuration: `std`, `virtio`, `serial0`, etc. |
| `scsi0`..`scsi3` | string | - | SCSI disk configuration (e.g. `local-lvm:32,iothread=1`). |
| `virtio0`, `virtio1` | string | - | VirtIO disk configuration. |
| `ide0`, `ide2` | string | - | IDE disk configuration. `ide2` is often used for cloud-init. |
| `efidisk0` | string | - | EFI disk configuration (for OVMF BIOS). |
| `tpmstate0` | string | - | TPM state configuration. |
| `net0`..`net3` | string | - | Network interface configuration (e.g. `virtio,bridge=vmbr0,firewall=1`). |
| `ci_user` | string | - | Cloud-init username. |
| `ci_password` | string | - | Cloud-init password. Sensitive. |
| `ci_type` | string | - | Cloud-init type: `configdrive2` or `nocloud`. |
| `ipconfig0`, `ipconfig1` | string | - | IP config for net0/net1 (e.g. `ip=dhcp` or `ip=10.0.0.2/24,gw=10.0.0.1`). |
| `nameserver` | string | - | DNS nameserver for cloud-init. |
| `searchdomain` | string | - | DNS search domain for cloud-init. |
| `ssh_keys` | string | - | Public SSH keys for cloud-init (newline-separated). |
| `serial0` | string | - | Serial device 0 (e.g. `socket`). |
| `clone_vmid` | integer | - | VMID of the template to clone. If set, VM is created as a clone. |
| `full_clone` | boolean | `true` | Full clone (`true`) or linked clone (`false`). |
| `status` | string | - | **Read-only.** Current VM status: `running`, `stopped`, etc. |
| `template` | boolean | - | **Read-only.** Whether the VM is a template. |

**Import:** `terraform import proxmox_virtual_environment_vm.example pve/100`

The import ID format is `<node_name>/<vmid>`.

---

#### `proxmox_virtual_environment_container`

Manages an LXC container. Supports creation of new containers and cloning from existing templates.

| Attribute | Type | Default | Description |
|---|---|---|---|
| `node_name` | string | - | **Required.** Proxmox node to create the container on. |
| `vmid` | integer | auto | Container ID. If omitted, assigned automatically. |
| `hostname` | string | - | Container hostname. |
| `description` | string | - | Container description. |
| `tags` | string | - | Semicolon-separated tags. |
| `os_template` | string | - | OS template to use (e.g. `local:vztmpl/ubuntu-22.04-standard_22.04-1_amd64.tar.zst`). Forces replacement if changed. |
| `os_type` | string | - | OS type: `debian`, `ubuntu`, `centos`, `fedora`, `opensuse`, `archlinux`, `alpine`, `gentoo`, `nixos`, `unmanaged`. |
| `on_boot` | boolean | `false` | Start container on host boot. |
| `started` | boolean | `true` | Start container after creation. |
| `protection` | boolean | `false` | Protect container from deletion. |
| `unprivileged` | boolean | `true` | Create unprivileged container (recommended for security). |
| `pool` | string | - | Resource pool. |
| `features` | string | - | Container feature flags (e.g. `nesting=1,keyctl=1`). |
| `cpu_cores` | integer | `1` | Number of CPU cores. |
| `cpu_limit` | integer | `0` | CPU usage limit (0 = unlimited). |
| `cpu_units` | integer | - | Relative CPU weight. |
| `memory` | integer | `512` | Memory in MiB. |
| `swap` | integer | `512` | Swap in MiB. |
| `rootfs` | string | - | Root filesystem configuration (e.g. `local-lvm:8`). |
| `net0`..`net3` | string | - | Network interface (e.g. `name=eth0,bridge=vmbr0,ip=dhcp`). |
| `mp0`..`mp2` | string | - | Additional mount point. |
| `nameserver` | string | - | DNS nameserver. |
| `searchdomain` | string | - | DNS search domain. |
| `password` | string | - | Root password. Sensitive. |
| `ssh_keys` | string | - | Public SSH keys. |
| `console` | boolean | `true` | Attach a console device. |
| `tty` | integer | `2` | Number of TTY devices (0-6). |
| `clone_vmid` | integer | - | VMID of the container to clone. |
| `full_clone` | boolean | `true` | Full or linked clone. |
| `status` | string | - | **Read-only.** Current status. |
| `template` | boolean | - | **Read-only.** Whether this is a template. |

**Import:** `terraform import proxmox_virtual_environment_container.example pve/300`

The import ID format is `<node_name>/<vmid>`.

---

#### `proxmox_virtual_environment_vm_snapshot`

Manages a snapshot of a VM.

| Attribute | Type | Description |
|---|---|---|
| `node_name` | string | **Required.** Node name. Forces replacement. |
| `vmid` | integer | **Required.** VM ID. Forces replacement. |
| `snap_name` | string | **Required.** Snapshot name. Forces replacement. |
| `description` | string | Snapshot description. |
| `include_ram` | boolean | Include RAM in the snapshot. |

**Import:** `terraform import proxmox_virtual_environment_vm_snapshot.example pve/100/mysnapshot`

---

#### `proxmox_virtual_environment_container_snapshot`

Manages a snapshot of an LXC container.

| Attribute | Type | Description |
|---|---|---|
| `node_name` | string | **Required.** Node name. Forces replacement. |
| `vmid` | integer | **Required.** Container ID. Forces replacement. |
| `snap_name` | string | **Required.** Snapshot name. Forces replacement. |
| `description` | string | Snapshot description. |

**Import:** `terraform import proxmox_virtual_environment_container_snapshot.example pve/300/mysnapshot`

---

### Storage

#### `proxmox_virtual_environment_storage`

Manages a global storage definition visible at **Datacenter > Storage**. Supports all Proxmox storage types.

**Supported types:** `dir`, `lvm`, `lvmthin`, `zfspool`, `nfs`, `cifs`, `glusterfs`, `iscsi`, `iscsidirect`, `rbd`, `cephfs`, `pbs`.

| Attribute | Type | Description |
|---|---|---|
| `storage` | string | **Required.** Storage identifier. Forces replacement. |
| `type` | string | **Required.** Storage type. Forces replacement. |
| `content` | string | Comma-separated content types: `images`, `rootdir`, `vztmpl`, `iso`, `backup`, `snippets`. |
| `enabled` | boolean | Enable the storage (default: `true`). |
| `shared` | boolean | Shared across all cluster nodes. |
| `nodes` | string | Comma-separated nodes on which the storage is available. |
| `path` | string | Filesystem path (for `dir` type). |
| `pool` | string | ZFS or Ceph pool name. |
| `vgname` | string | LVM volume group name. |
| `server` | string | Server address (for NFS, CIFS, iSCSI, PBS). |
| `export` | string | NFS export path. |
| `share` | string | CIFS share name. |
| `username` | string | Username (for CIFS, PBS). |
| `password` | string | Password (for CIFS, PBS). Sensitive. |
| `domain` | string | Windows domain (for CIFS). |
| `datastore` | string | PBS datastore name. |
| `namespace` | string | PBS namespace. |
| `fingerprint` | string | PBS server certificate fingerprint. |
| `prune_backups` | string | Backup retention policy (e.g. `keep-last=7,keep-weekly=4`). |

**Import:** `terraform import proxmox_virtual_environment_storage.example local-backup`

---

#### `proxmox_virtual_environment_download_file`

Downloads a file from a URL to a Proxmox VE node storage. Uses the Proxmox download API which runs the transfer on the node side.

| Attribute | Type | Description |
|---|---|---|
| `node_name` | string | **Required.** Node name. Forces replacement. |
| `storage` | string | **Required.** Storage name. Forces replacement. |
| `url` | string | **Required.** Source URL to download from. Forces replacement. |
| `file_name` | string | **Required.** Target filename on the storage. Forces replacement. |
| `content_type` | string | **Required.** Content type: `iso` or `vztmpl`. Forces replacement. |
| `checksum` | string | Expected checksum for verification. Forces replacement. |
| `checksum_algorithm` | string | Checksum algorithm (e.g. `md5`, `sha1`, `sha256`). Forces replacement. |
| `verify_tls` | boolean | Verify TLS certificate of the download URL. Default: `true`. |
| `size` | integer | **Read-only.** File size in bytes. |
| `file_id` | string | **Read-only.** Full volume ID of the file in storage. |

**Import:** `terraform import proxmox_virtual_environment_download_file.example pve/local/iso/ubuntu-24.04-live-server-amd64.iso`

---

### Networking

#### `proxmox_virtual_environment_network_interface`

Manages a network interface on a Proxmox VE node. Supports bridge, bond, VLAN, and physical interface types.

| Attribute | Type | Default | Description |
|---|---|---|---|
| `node_name` | string | - | **Required.** Node name. Forces replacement. |
| `iface` | string | - | **Required.** Interface name (e.g. `vmbr0`, `bond0`, `eth0.100`). Forces replacement. |
| `type` | string | - | **Required.** Type: `bridge`, `bond`, `eth`, `vlan`, `OVSBridge`, `OVSBond`, `OVSPort`, `OVSIntPort`. Forces replacement. |
| `method` | string | - | IPv4 method: `static`, `dhcp`, `manual`. |
| `method6` | string | - | IPv6 method: `static`, `dhcp`, `manual`. |
| `address` | string | - | IPv4 address. |
| `netmask` | string | - | IPv4 netmask. |
| `gateway` | string | - | IPv4 gateway. |
| `cidr` | string | - | IPv4 in CIDR format (e.g. `192.168.1.1/24`). |
| `address6` | string | - | IPv6 address. |
| `netmask6` | integer | - | IPv6 prefix length. |
| `gateway6` | string | - | IPv6 gateway. |
| `bridge_ports` | string | - | Bridge ports (space-separated interface names). |
| `bridge_stp` | string | - | STP mode: `on` or `off`. |
| `bridge_fd` | integer | `0` | Bridge forward delay. |
| `bridge_vlan_aware` | boolean | `false` | Enable VLAN-aware bridge. |
| `bond_mode` | string | - | Bond mode: `balance-rr`, `active-backup`, `balance-xor`, etc. |
| `slaves` | string | - | Slave interfaces for the bond. |
| `vlan_raw_device` | string | - | Underlying device for the VLAN interface. |
| `vlan_id` | integer | - | VLAN ID. |
| `autostart` | boolean | `true` | Bring up the interface at boot. |
| `mtu` | integer | - | Interface MTU. |
| `comments` | string | - | Interface comments. |
| `apply_config` | boolean | `false` | Apply network configuration immediately after changes. |

**Import:** `terraform import proxmox_virtual_environment_network_interface.example pve/vmbr0`

---

### Firewall

All firewall resources that operate on a scope use the following scope format:

| Scope | Format | Example |
|---|---|---|
| Cluster | `cluster` | `cluster` |
| Node | `node/<node_name>` | `node/pve` |
| VM | `vm/<node_name>/<vmid>` | `vm/pve/100` |
| Container | `ct/<node_name>/<vmid>` | `ct/pve/300` |

#### `proxmox_virtual_environment_firewall_rule`

Manages a firewall rule. The `scope` attribute determines whether it applies at cluster, node, VM, or container level.

| Attribute | Type | Default | Description |
|---|---|---|---|
| `scope` | string | - | **Required.** Rule scope. Forces replacement. |
| `type` | string | - | **Required.** Direction: `in` or `out`. |
| `action` | string | - | **Required.** Action: `ACCEPT`, `DROP`, `REJECT`. |
| `enabled` | boolean | `true` | Enable the rule. |
| `pos` | integer | - | **Read-only.** Position in the rule list. |
| `macro` | string | - | Predefined macro name (e.g. `SSH`, `HTTP`, `HTTPS`). |
| `proto` | string | - | Protocol: `tcp`, `udp`, `icmp`, etc. |
| `source` | string | - | Source address/CIDR/IPset. |
| `dest` | string | - | Destination address/CIDR/IPset. |
| `dport` | string | - | Destination port(s). |
| `sport` | string | - | Source port(s). |
| `iface` | string | - | Network interface. |
| `log` | string | - | Log level: `emerg`, `alert`, `crit`, `err`, `warning`, `notice`, `info`, `debug`, `nolog`. |
| `comment` | string | - | Rule comment. |

**Import:** `terraform import proxmox_virtual_environment_firewall_rule.example cluster/0`

The import ID format is `<scope>/<position>` (e.g. `cluster/0`, `vm/pve/200/3`).

---

#### `proxmox_virtual_environment_firewall_options`

Manages firewall options for a given scope (cluster, node, VM, or container). Deleting this resource disables the firewall without removing it (it is a singleton per scope).

| Attribute | Type | Default | Description |
|---|---|---|---|
| `scope` | string | - | **Required.** Firewall scope. Forces replacement. |
| `enabled` | boolean | `false` | Enable the firewall. |
| `policy_in` | string | - | Default inbound policy: `ACCEPT`, `DROP`, `REJECT`. |
| `policy_out` | string | - | Default outbound policy: `ACCEPT`, `DROP`, `REJECT`. |
| `ebtables` | boolean | `true` | Enable Ethernet bridge filtering. |
| `ip_filter` | boolean | `false` | Enable IP filter (block IPs not associated with the VM). |
| `mac_filter` | boolean | `true` | Enable MAC filter. |
| `ndp` | boolean | `false` | Enable NDP (Neighbor Discovery Protocol) for IPv6. |
| `dhcp` | boolean | `false` | Allow DHCP traffic. |
| `log_ratelimit` | string | - | Log rate limit (e.g. `enable=1,rate=1/second,burst=5`). |

**Import:** `terraform import proxmox_virtual_environment_firewall_options.example cluster`

---

#### `proxmox_virtual_environment_firewall_ipset`

Manages a firewall IP set with CIDR entries.

| Attribute | Type | Description |
|---|---|---|
| `scope` | string | **Required.** Firewall scope. Forces replacement. |
| `name` | string | **Required.** IP set name. Forces replacement. |
| `comment` | string | IP set description. |
| `cidrs` | list | List of CIDR entries. Each entry has `cidr`, `comment`, and `no_match` attributes. |

**Import:** `terraform import proxmox_virtual_environment_firewall_ipset.example cluster/my-ipset`

The import ID format is `<scope>/<name>`.

---

#### `proxmox_virtual_environment_firewall_alias`

Manages a named IP/CIDR alias.

| Attribute | Type | Description |
|---|---|---|
| `scope` | string | **Required.** Firewall scope. Forces replacement. |
| `name` | string | **Required.** Alias name. Forces replacement. |
| `cidr` | string | **Required.** IP or CIDR that the alias represents. |
| `comment` | string | Alias description. |

**Import:** `terraform import proxmox_virtual_environment_firewall_alias.example cluster/my-alias`

The import ID format is `<scope>/<name>`.

---

#### `proxmox_virtual_environment_firewall_security_group`

Manages a cluster-level firewall security group. Rules within the group are managed by the separate `firewall_security_group_rule` resource.

| Attribute | Type | Description |
|---|---|---|
| `name` | string | **Required.** Security group name. Forces replacement. |
| `comment` | string | Comment. |

**Import:** `terraform import proxmox_virtual_environment_firewall_security_group.example my-sg`

---

#### `proxmox_virtual_environment_firewall_security_group_rule`

Manages a rule inside a cluster firewall security group.

| Attribute | Type | Default | Description |
|---|---|---|---|
| `security_group` | string | - | **Required.** Security group name. Forces replacement. |
| `pos` | integer | - | **Read-only.** Rule position. |
| `type` | string | - | **Required.** Direction: `in` or `out`. |
| `action` | string | - | **Required.** Action: `ACCEPT`, `DROP`, `REJECT`. |
| `enabled` | boolean | `true` | Enable the rule. |
| `macro` | string | - | Macro name (e.g. `SSH`, `HTTP`). |
| `proto` | string | - | Protocol. |
| `source` | string | - | Source address/CIDR/IPset. |
| `dest` | string | - | Destination address/CIDR/IPset. |
| `dport` | string | - | Destination port(s). |
| `sport` | string | - | Source port(s). |
| `iface` | string | - | Network interface. |
| `log` | string | - | Log level. |
| `comment` | string | - | Rule comment. |

**Import:** `terraform import proxmox_virtual_environment_firewall_security_group_rule.example my-sg/0`

The import ID format is `<security_group>/<position>`.

---

### High Availability

#### `proxmox_virtual_environment_ha_resource`

Manages a Proxmox VE High Availability resource.

| Attribute | Type | Default | Description |
|---|---|---|---|
| `sid` | string | - | **Required.** HA resource SID (e.g. `vm:100` or `ct:200`). Forces replacement. |
| `type` | string | - | Resource type: `vm` or `ct`. |
| `state` | string | - | Desired state: `started`, `stopped`, `enabled`, `disabled`, `ignored`. |
| `group` | string | - | HA group name. |
| `max_restart` | integer | `1` | Maximum restart attempts. |
| `max_relocate` | integer | `1` | Maximum relocation attempts. |
| `comment` | string | - | Description. |

**Import:** `terraform import proxmox_virtual_environment_ha_resource.example vm:100`

---

#### `proxmox_virtual_environment_ha_group`

Manages a Proxmox VE High Availability group.

| Attribute | Type | Default | Description |
|---|---|---|---|
| `group` | string | - | **Required.** HA group name. Forces replacement. |
| `nodes` | string | - | **Required.** Comma-separated nodes with optional priority (e.g. `node1:2,node2:1`). |
| `restricted` | boolean | `false` | Restrict HA group to its members. |
| `no_failback` | boolean | `false` | Disable automatic failback. |
| `comment` | string | - | Description. |

**Import:** `terraform import proxmox_virtual_environment_ha_group.example my-ha-group`

---

### Node Management

#### `proxmox_virtual_environment_dns`

Manages the DNS configuration of a Proxmox VE node.

| Attribute | Type | Description |
|---|---|---|
| `node_name` | string | **Required.** Node name. Forces replacement. |
| `search` | string | **Required.** DNS search domain. |
| `dns1` | string | First DNS server. |
| `dns2` | string | Second DNS server. |
| `dns3` | string | Third DNS server. |

**Import:** `terraform import proxmox_virtual_environment_dns.example pve`

---

#### `proxmox_virtual_environment_hosts`

Manages the `/etc/hosts` file of a Proxmox VE node.

| Attribute | Type | Description |
|---|---|---|
| `node_name` | string | **Required.** Node name. Forces replacement. |
| `entries` | string | **Required.** Full content of `/etc/hosts`. |
| `digest` | string | **Read-only.** Hosts file digest (for conflict detection). |

**Import:** `terraform import proxmox_virtual_environment_hosts.example pve`

---

#### `proxmox_virtual_environment_certificate`

Manages a custom TLS certificate on a Proxmox VE node.

| Attribute | Type | Default | Description |
|---|---|---|---|
| `node_name` | string | - | **Required.** Node name. Forces replacement. |
| `certificate` | string | - | **Required.** PEM-encoded certificate (and chain). |
| `private_key` | string | - | PEM-encoded private key. Sensitive. |
| `force` | boolean | `true` | Overwrite existing certificate. |
| `restart` | boolean | `true` | Restart pveproxy after uploading. |
| `fingerprint` | string | - | **Read-only.** Certificate SHA-256 fingerprint. |
| `issuer` | string | - | **Read-only.** Certificate issuer. |
| `subject` | string | - | **Read-only.** Certificate subject. |
| `not_before` | integer | - | **Read-only.** Validity start (Unix timestamp). |
| `not_after` | integer | - | **Read-only.** Validity end (Unix timestamp). |

**Import:** `terraform import proxmox_virtual_environment_certificate.example pve`

---

#### `proxmox_virtual_environment_time`

Manages the timezone configuration of a Proxmox VE node.

| Attribute | Type | Description |
|---|---|---|
| `node_name` | string | **Required.** Node name. Forces replacement. |
| `timezone` | string | **Required.** Timezone (e.g. `Europe/Rome`, `UTC`). |

**Import:** `terraform import proxmox_virtual_environment_time.example pve`

---

### Cluster

#### `proxmox_virtual_environment_cluster_options`

Manages cluster-wide Proxmox VE options. This is a singleton resource; there is one per cluster. The import ID is always `cluster`.

| Attribute | Type | Description |
|---|---|---|
| `keyboard` | string | Default keyboard layout (e.g. `en-us`, `de`, `fr`). |
| `language` | string | Default web UI language. |
| `email_from` | string | Email address used as sender for notifications. |
| `http_proxy` | string | HTTP proxy URL for the cluster (used for apt, etc.). |
| `max_workers` | integer | Maximum number of workers for bulk operations. |
| `migration_unsecure` | boolean | Allow unsecured (non-TLS) migrations. |
| `migration_type` | string | Migration type: `secure`, `insecure`, or `websocket`. |
| `ha_shutdown_policy` | string | HA shutdown policy: `freeze`, `failover`, `conditional`, or `migrate`. |

**Import:** `terraform import proxmox_virtual_environment_cluster_options.example cluster`

---

### Replication

#### `proxmox_virtual_environment_replication`

Manages a Proxmox VE replication job. The job ID in state is formatted as `<vmid>-<job_id>`.

| Attribute | Type | Default | Description |
|---|---|---|---|
| `job_id` | integer | - | **Required.** Numeric job ID within the VM (e.g. `0`, `1`). |
| `source_vm_id` | integer | - | **Required.** VMID to replicate. |
| `target_node` | string | - | **Required.** Destination node name. Forces replacement. |
| `schedule` | string | - | Replication schedule (e.g. `*/15`). |
| `rate` | float | - | Bandwidth limit in MB/s. |
| `comment` | string | - | Comment. |
| `enabled` | boolean | `true` | Enable the replication job. |

**Import:** `terraform import proxmox_virtual_environment_replication.example 100-0`

The import ID format is `<vmid>-<job_id>` (e.g. `100-0`).

---

### SDN

#### `proxmox_virtual_environment_sdn_zone`

Manages a Proxmox VE SDN zone.

| Attribute | Type | Description |
|---|---|---|
| `zone` | string | **Required.** SDN zone name. Forces replacement. |
| `type` | string | **Required.** Zone type: `simple`, `vlan`, `qinq`, `vxlan`, `evpn`. Forces replacement. |
| `comment` | string | Description. |
| `bridge` | string | Bridge interface (for `vlan`/`qinq` zones). |
| `tag` | integer | VLAN tag (for `vlan`/`qinq` zones). |
| `peers` | string | Comma-separated VXLAN peer addresses (for `vxlan` zones). |
| `vrf_vxlan` | integer | VRF VXLAN tag (for `evpn` zones). |
| `controller` | string | EVPN controller name (for `evpn` zones). |
| `mtu` | integer | MTU value for the zone. |
| `dns` | string | DNS plugin name. |
| `dns_zone` | string | DNS domain. |
| `reverse_dns` | string | Reverse DNS plugin name. |
| `ipam` | string | IPAM plugin name. |

**Import:** `terraform import proxmox_virtual_environment_sdn_zone.example my-zone`

---

#### `proxmox_virtual_environment_sdn_vnet`

Manages a Proxmox VE SDN VNet.

| Attribute | Type | Description |
|---|---|---|
| `vnet` | string | **Required.** VNet name (max 8 characters). Forces replacement. |
| `zone` | string | **Required.** SDN zone this VNet belongs to. |
| `alias` | string | Alias/description for the VNet. |
| `tag` | integer | VLAN tag (for VLAN-aware zones). |
| `vlan_aware` | boolean | Whether the VNet is VLAN-aware. |
| `comment` | string | Description. |

**Import:** `terraform import proxmox_virtual_environment_sdn_vnet.example myvnet`

---

#### `proxmox_virtual_environment_sdn_subnet`

Manages a Proxmox VE SDN subnet within a VNet.

| Attribute | Type | Default | Description |
|---|---|---|---|
| `vnet` | string | - | **Required.** Parent VNet name. Forces replacement. |
| `subnet` | string | - | **Required.** Subnet CIDR (e.g. `10.0.0.0/24`). Forces replacement. |
| `gateway` | string | - | Subnet gateway IP address. |
| `snat` | boolean | `false` | Enable SNAT for outbound traffic. |
| `dhcp_dns_server` | string | - | DNS server pushed to DHCP clients. |
| `dns_zone_prefix` | string | - | DNS zone prefix for PTR records. |

**Import:** `terraform import proxmox_virtual_environment_sdn_subnet.example myvnet/10.0.0.0/24`

The import ID format is `<vnet>/<subnet-cidr>`.

---

### ACME

#### `proxmox_virtual_environment_acme_account`

Manages a Proxmox VE ACME account registration (e.g. Let's Encrypt).

| Attribute | Type | Default | Description |
|---|---|---|---|
| `name` | string | `default` | Account name. Forces replacement. |
| `contact` | string | - | **Required.** Contact email address. |
| `directory` | string | - | ACME directory URL. Defaults to Let's Encrypt production. Forces replacement. |
| `tos_url` | string | - | URL of the ACME Terms of Service (must be accepted at account creation). |
| `account_url` | string | - | **Read-only.** Registered ACME account URL (assigned by the CA). |

**Import:** `terraform import proxmox_virtual_environment_acme_account.example default`

---

#### `proxmox_virtual_environment_acme_plugin`

Manages a Proxmox VE ACME DNS challenge plugin.

| Attribute | Type | Description |
|---|---|---|
| `plugin_id` | string | **Required.** Plugin identifier. Forces replacement. |
| `type` | string | **Required.** Plugin type: `dns` or `standalone`. Forces replacement. |
| `api` | string | DNS API name (e.g. `cf` for Cloudflare, `aws` for Route 53). |
| `data` | string | DNS API credentials in `key=value` format (one per line). Sensitive. |
| `nodes` | string | Comma-separated list of nodes to restrict the plugin to. |
| `validation_delay` | integer | Delay in seconds to wait for DNS propagation. |

**Import:** `terraform import proxmox_virtual_environment_acme_plugin.example my-dns-plugin`

---

### Metrics

#### `proxmox_virtual_environment_metrics_server`

Manages an external metrics server (Graphite or InfluxDB).

| Attribute | Type | Default | Description |
|---|---|---|---|
| `name` | string | - | **Required.** Metrics server identifier. Forces replacement. |
| `type` | string | - | **Required.** Server type: `graphite` or `influxdb`. Forces replacement. |
| `server` | string | - | **Required.** Hostname or IP address. |
| `port` | integer | - | **Required.** Server port. |
| `enabled` | boolean | `true` | Enable the metrics server. |
| `mtu` | integer | - | MTU (for InfluxDB UDP). |
| `path` | string | - | Root Graphite path (for Graphite type). |
| `proto` | string | - | Protocol: `udp` or `tcp` (for Graphite type). |
| `timeout` | integer | - | TCP socket connection timeout in seconds. |
| `bucket` | string | - | InfluxDB bucket/database name. |
| `influxdb_proto` | string | - | InfluxDB protocol: `udp`, `http`, or `https`. |
| `organization` | string | - | InfluxDB organization name. |
| `token` | string | - | InfluxDB access token. Sensitive. |
| `max_body_size` | integer | - | Maximum body size in bytes for InfluxDB HTTP(S). |

**Import:** `terraform import proxmox_virtual_environment_metrics_server.example my-influxdb`

---

### Hardware Mappings

#### `proxmox_virtual_environment_hardware_mapping_pci`

Manages a cluster-level PCI hardware mapping.

| Attribute | Type | Description |
|---|---|---|
| `name` | string | **Required.** Mapping name. Forces replacement. |
| `comment` | string | Description. |
| `map` | list(string) | **Required.** List of per-node PCI device entries in Proxmox format: `node=<node>,path=<pci-path>,id=<vendor>:<device>[,iommu-group=<n>]`. |
| `mdevs` | string | Mediated device types. |

**Import:** `terraform import proxmox_virtual_environment_hardware_mapping_pci.example my-gpu`

---

#### `proxmox_virtual_environment_hardware_mapping_usb`

Manages a cluster-level USB hardware mapping.

| Attribute | Type | Description |
|---|---|---|
| `name` | string | **Required.** Mapping name. Forces replacement. |
| `comment` | string | Description. |
| `map` | list(string) | **Required.** List of per-node USB device entries in Proxmox format. |

**Import:** `terraform import proxmox_virtual_environment_hardware_mapping_usb.example my-usb-device`

---

### APT

#### `proxmox_virtual_environment_apt_repository`

Manages a standard Proxmox VE APT repository on a node. Only standard (built-in) repository handles are supported.

| Attribute | Type | Default | Description |
|---|---|---|---|
| `node_name` | string | - | **Required.** Node name. Forces replacement. |
| `handle` | string | - | **Required.** Standard repository handle (e.g. `pve-no-subscription`, `pve-enterprise`). Forces replacement. |
| `enabled` | boolean | `true` | Enable the repository. |
| `file_path` | string | - | **Read-only.** Sources file path. |
| `index` | integer | - | **Read-only.** Index within the sources file. |

**Import:** `terraform import proxmox_virtual_environment_apt_repository.example pve/pve-no-subscription`

---

### Backup

#### `proxmox_virtual_environment_backup`

Manages a vzdump backup schedule.

| Attribute | Type | Default | Description |
|---|---|---|---|
| `storage` | string | - | **Required.** Target storage for backups. |
| `schedule` | string | - | **Required.** Schedule in systemd calendar format (e.g. `daily`, `Mon,Tue 02:00`). |
| `vmids` | string | - | Comma-separated VMID(s) to back up. Leave empty when `all` is `true`. |
| `nodes` | string | - | Comma-separated nodes to run the job on. Empty means all nodes. |
| `all` | boolean | `false` | Back up all VMs and containers. |
| `compress` | string | - | Compression algorithm: `0`, `1`, `gzip`, `lzo`, `zstd`. |
| `mode` | string | - | Backup mode: `snapshot`, `suspend`, `stop`. |
| `comment` | string | - | Job description. |
| `mailto` | string | - | Comma-separated email addresses to notify. |
| `mail_notification` | string | - | Email notification mode: `always` or `failure`. |
| `max_files` | integer | - | Maximum number of backups to keep (deprecated; prefer prune settings). |
| `enabled` | boolean | `true` | Enable the backup job. |
| `bw_limit` | float | - | Bandwidth limit in KiB/s. |
| `notes_template` | string | - | Template for backup notes. |

**Import:** `terraform import proxmox_virtual_environment_backup.example <job-id>`

The import ID is the backup job ID assigned by Proxmox.

---

## Data Sources

### `proxmox_virtual_environment_version`

Returns Proxmox VE API version information.

```hcl
data "proxmox_virtual_environment_version" "current" {}

output "pve_version" {
  value = data.proxmox_virtual_environment_version.current.version
}
```

**Exported attributes:** `release`, `repo_id`, `version`.

---

### `proxmox_virtual_environment_nodes`

Returns a list of all nodes in the cluster. All exported attributes are parallel lists indexed by node position.

```hcl
data "proxmox_virtual_environment_nodes" "cluster" {}

output "online_nodes" {
  value = [
    for i, online in data.proxmox_virtual_environment_nodes.cluster.online :
    data.proxmox_virtual_environment_nodes.cluster.names[i]
    if online
  ]
}
```

**Exported attributes (parallel lists):** `names`, `online`, `cpu_count`, `cpu_utilization`, `memory_used`, `memory_available`, `uptime`, `ssl_fingerprints`, `support_levels`.

---

### `proxmox_virtual_environment_node`

Returns detailed status of a specific node, including CPU, memory, swap, filesystem, and boot details.

**Input:** `node_name`.

**Exported attributes:** `cpu_cores`, `cpu_sockets`, `cpu_threads`, `cpu_model`, `cpu_mhz`, `cpu_usage`, `load_average`, `memory_total`, `memory_used`, `memory_free`, `swap_total`, `swap_used`, `swap_free`, `rootfs_total`, `rootfs_used`, `rootfs_free`, `uptime`, `kernel_version`, `pve_version`, `boot_mode`, `secure_boot`.

---

### `proxmox_virtual_environment_datastores`

Returns the list of storage visible on a specific node. All attributes are parallel lists.

**Input:** `node_name`.

**Exported attributes (parallel lists):** `names`, `types`, `content_types`, `active`, `enabled`, `shared`, `total`, `used`, `available`, `used_fraction`.

---

### `proxmox_virtual_environment_vms`

Returns summary information for all VMs on a node. All attributes are parallel lists.

**Input:** `node_name`.

**Exported attributes (parallel lists):** `vmids`, `names`, `statuses`, `tags`, `cpus`, `max_memory`, `max_disk`, `uptime`, `template`.

---

### `proxmox_virtual_environment_vm`

Returns detailed configuration of a specific VM.

**Input:** `node_name`, `vmid`.

---

### `proxmox_virtual_environment_containers`

Returns summary information for all LXC containers on a node. Structure is analogous to `proxmox_virtual_environment_vms`.

**Input:** `node_name`.

**Exported attributes (parallel lists):** `vmids`, `names`, `statuses`, `tags`, `cpus`, `max_memory`, `max_disk`, `uptime`, `template`.

---

### `proxmox_virtual_environment_container`

Returns detailed configuration of a specific container.

**Input:** `node_name`, `vmid`.

---

### `proxmox_virtual_environment_vm_snapshots`

Lists all snapshots for a VM.

**Input:** `node_name`, `vmid`.

---

### `proxmox_virtual_environment_container_snapshots`

Lists all snapshots for an LXC container.

**Input:** `node_name`, `vmid`.

---

### `proxmox_virtual_environment_ha_resources`

Lists all HA resources.

---

### `proxmox_virtual_environment_ha_groups`

Lists all HA groups.

---

### `proxmox_virtual_environment_users`

Lists all users.

---

### `proxmox_virtual_environment_groups`

Lists all groups.

---

### `proxmox_virtual_environment_roles`

Lists all roles.

---

### `proxmox_virtual_environment_pools`

Lists all resource pools.

---

### `proxmox_virtual_environment_network_interfaces`

Lists network interfaces on a specific node.

**Input:** `node_name`.

---

### `proxmox_virtual_environment_hardware_mapping_pci_list`

Lists all PCI hardware mappings.

---

### `proxmox_virtual_environment_hardware_mapping_usb_list`

Lists all USB hardware mappings.

---

### `proxmox_virtual_environment_sdn_zones`

Lists all SDN zones.

---

### `proxmox_virtual_environment_sdn_vnets`

Lists all SDN VNets.

---

### `proxmox_virtual_environment_apt_repositories`

Lists APT repositories on a specific node.

**Input:** `node_name`.

---

### `proxmox_virtual_environment_backups`

Lists all backup jobs.

---

## Resource Examples

### Access Management

```hcl
# Create a user, group, role, and ACL entry
resource "proxmox_virtual_environment_group" "operators" {
  group_id = "operators"
  comment  = "Infrastructure operators"
}

resource "proxmox_virtual_environment_user" "alice" {
  user_id    = "alice@pve"
  password   = var.alice_password
  first_name = "Alice"
  email      = "alice@example.com"
  enabled    = true
  groups     = proxmox_virtual_environment_group.operators.group_id
}

resource "proxmox_virtual_environment_role" "vm_operator" {
  role_id    = "VMOperator"
  privileges = "VM.Audit,VM.Console,VM.PowerMgmt,VM.Config.Disk"
}

resource "proxmox_virtual_environment_acl" "operators_vms" {
  path      = "/vms"
  role_id   = proxmox_virtual_environment_role.vm_operator.role_id
  group_id  = proxmox_virtual_environment_group.operators.group_id
  propagate = true
}

resource "proxmox_virtual_environment_user_token" "alice_terraform" {
  user_id              = proxmox_virtual_environment_user.alice.user_id
  token_id             = "terraform"
  comment              = "Terraform automation token"
  privileges_separation = true
}

output "alice_token" {
  value     = proxmox_virtual_environment_user_token.alice_terraform.value
  sensitive = true
}
```

### Compute — VM with Cloud-Init

```hcl
resource "proxmox_virtual_environment_vm" "ubuntu_web" {
  node_name = "pve"
  name      = "ubuntu-web"
  vmid      = 200

  bios    = "ovmf"
  machine = "q35"
  os_type = "l26"
  scsi_hw = "virtio-scsi-pci"

  cpu_cores   = 4
  cpu_sockets = 1
  cpu_type    = "host"
  memory      = 4096

  efidisk0 = "local-lvm:1,efitype=4m,pre-enrolled-keys=0"
  scsi0    = "local-lvm:32,iothread=1"
  ide2     = "local:cloudinit"

  net0 = "virtio,bridge=vmbr0,firewall=1"

  ci_user      = "admin"
  ci_password  = var.vm_password
  ipconfig0    = "ip=10.0.0.10/24,gw=10.0.0.1"
  nameserver   = "1.1.1.1"
  searchdomain = "example.com"
  ssh_keys     = file("~/.ssh/id_ed25519.pub")

  agent      = true
  on_boot    = true
  started    = true
  protection = false
  tags       = "ubuntu;web;prod"
  pool       = "production"
}
```

### Compute — LXC Container

```hcl
resource "proxmox_virtual_environment_container" "app" {
  node_name   = "pve"
  hostname    = "app-01"
  vmid        = 300
  os_template = "local:vztmpl/ubuntu-22.04-standard_22.04-1_amd64.tar.zst"
  os_type     = "ubuntu"

  cpu_cores = 2
  memory    = 2048
  swap      = 1024

  rootfs = "local-lvm:20"

  net0 = "name=eth0,bridge=vmbr0,ip=10.0.0.30/24,gw=10.0.0.1,firewall=1"

  nameserver   = "1.1.1.1"
  searchdomain = "example.com"

  ssh_keys     = file("~/.ssh/id_ed25519.pub")
  unprivileged = true
  features     = "nesting=1"

  on_boot = true
  started = true
  tags    = "ubuntu;app;prod"
  pool    = "production"
}
```

### Storage

```hcl
# Local directory storage
resource "proxmox_virtual_environment_storage" "local_backup" {
  storage = "local-backup"
  type    = "dir"
  path    = "/mnt/backup"
  content = "backup,iso"
  enabled = true
}

# NFS shared storage
resource "proxmox_virtual_environment_storage" "nfs_iso" {
  storage = "nfs-iso"
  type    = "nfs"
  server  = "nas.example.com"
  export  = "/exports/iso"
  content = "iso,vztmpl"
  shared  = true
}

# Proxmox Backup Server
resource "proxmox_virtual_environment_storage" "pbs_main" {
  storage       = "pbs-main"
  type          = "pbs"
  server        = "pbs.example.com"
  datastore     = "main"
  username      = "backup@pbs"
  password      = var.pbs_password
  fingerprint   = "AB:CD:EF:..."
  content       = "backup"
  prune_backups = "keep-last=7,keep-weekly=4,keep-monthly=3"
}
```

### Networking

```hcl
# VLAN-aware bridge with static IP
resource "proxmox_virtual_environment_network_interface" "vmbr0" {
  node_name         = "pve"
  iface             = "vmbr0"
  type              = "bridge"
  method            = "static"
  cidr              = "192.168.1.1/24"
  gateway           = "192.168.1.254"
  bridge_ports      = "eth0"
  bridge_stp        = "off"
  bridge_fd         = 0
  bridge_vlan_aware = true
  autostart         = true
  apply_config      = true
}
```

### Firewall

```hcl
# Enable firewall at cluster level
resource "proxmox_virtual_environment_firewall_options" "cluster" {
  scope      = "cluster"
  enabled    = true
  policy_in  = "DROP"
  policy_out = "ACCEPT"
  ebtables   = true
}

# Allow SSH from management network
resource "proxmox_virtual_environment_firewall_rule" "ssh_mgmt" {
  scope   = "cluster"
  type    = "in"
  action  = "ACCEPT"
  macro   = "SSH"
  source  = "10.0.0.0/8"
  comment = "SSH from management"
  enabled = true
}

# IP set for trusted hosts
resource "proxmox_virtual_environment_firewall_ipset" "trusted" {
  scope   = "cluster"
  name    = "trusted-hosts"
  comment = "Management and monitoring hosts"

  cidrs = [
    { cidr = "10.0.0.0/24", comment = "Management network" },
    { cidr = "10.0.1.50",   comment = "Monitoring server" },
  ]
}

# Security group for web servers
resource "proxmox_virtual_environment_firewall_security_group" "web" {
  name    = "web-servers"
  comment = "Rules for web-facing VMs"
}

resource "proxmox_virtual_environment_firewall_security_group_rule" "web_http" {
  security_group = proxmox_virtual_environment_firewall_security_group.web.name
  type           = "in"
  action         = "ACCEPT"
  proto          = "tcp"
  dport          = "80"
  comment        = "HTTP"
}

resource "proxmox_virtual_environment_firewall_security_group_rule" "web_https" {
  security_group = proxmox_virtual_environment_firewall_security_group.web.name
  type           = "in"
  action         = "ACCEPT"
  proto          = "tcp"
  dport          = "443"
  comment        = "HTTPS"
}
```

### High Availability

```hcl
resource "proxmox_virtual_environment_ha_group" "production" {
  group       = "production"
  nodes       = "pve1:2,pve2:1"
  restricted  = false
  no_failback = false
  comment     = "Production HA group"
}

resource "proxmox_virtual_environment_ha_resource" "web_vm" {
  sid          = "vm:200"
  state        = "started"
  group        = proxmox_virtual_environment_ha_group.production.group
  max_restart  = 3
  max_relocate = 2
  comment      = "Web server VM"
}
```

### SDN

```hcl
resource "proxmox_virtual_environment_sdn_zone" "vxlan_zone" {
  zone    = "vxlan1"
  type    = "vxlan"
  peers   = "10.0.0.1,10.0.0.2,10.0.0.3"
  comment = "VXLAN overlay zone"
}

resource "proxmox_virtual_environment_sdn_vnet" "app_vnet" {
  vnet    = "appvnet"
  zone    = proxmox_virtual_environment_sdn_zone.vxlan_zone.zone
  comment = "Application VNet"
}

resource "proxmox_virtual_environment_sdn_subnet" "app_subnet" {
  vnet    = proxmox_virtual_environment_sdn_vnet.app_vnet.vnet
  subnet  = "10.100.0.0/24"
  gateway = "10.100.0.1"
  snat    = true
}
```

### ACME and Certificates

```hcl
resource "proxmox_virtual_environment_acme_account" "letsencrypt" {
  name    = "default"
  contact = "admin@example.com"
}

resource "proxmox_virtual_environment_acme_plugin" "cloudflare" {
  plugin_id = "cf-dns"
  type      = "dns"
  api       = "cf"
  data      = "CF_Token=${var.cf_api_token}\nCF_Account_ID=${var.cf_account_id}"
}
```

### Metrics

```hcl
resource "proxmox_virtual_environment_metrics_server" "influxdb" {
  name           = "influxdb-prod"
  type           = "influxdb"
  server         = "influxdb.example.com"
  port           = 8086
  influxdb_proto = "https"
  bucket         = "proxmox"
  organization   = "myorg"
  token          = var.influxdb_token
  enabled        = true
}
```

### Backup Schedule

```hcl
resource "proxmox_virtual_environment_backup" "nightly" {
  storage          = "pbs-main"
  schedule         = "daily"
  all              = true
  compress         = "zstd"
  mode             = "snapshot"
  comment          = "Nightly backup of all VMs"
  mailto           = "ops@example.com"
  mail_notification = "failure"
  enabled          = true
}
```

---

## Import Reference

The table below summarizes the import ID format for every resource that supports `terraform import`.

| Resource Type | Import ID Format | Example |
|---|---|---|
| `proxmox_virtual_environment_acl` | Not supported | — |
| `proxmox_virtual_environment_acme_account` | `<account_name>` | `default` |
| `proxmox_virtual_environment_acme_plugin` | `<plugin_id>` | `cf-dns` |
| `proxmox_virtual_environment_apt_repository` | `<node>/<handle>` | `pve/pve-no-subscription` |
| `proxmox_virtual_environment_backup` | `<job_id>` | `backup-123456` |
| `proxmox_virtual_environment_certificate` | `<node_name>` | `pve` |
| `proxmox_virtual_environment_cluster_options` | `cluster` | `cluster` |
| `proxmox_virtual_environment_container` | `<node_name>/<vmid>` | `pve/300` |
| `proxmox_virtual_environment_container_snapshot` | `<node_name>/<vmid>/<snap_name>` | `pve/300/mysnapshot` |
| `proxmox_virtual_environment_dns` | `<node_name>` | `pve` |
| `proxmox_virtual_environment_download_file` | `<node>/<storage>/<volid>` | `pve/local/iso/ubuntu.iso` |
| `proxmox_virtual_environment_firewall_alias` | `<scope>/<name>` | `cluster/my-alias` |
| `proxmox_virtual_environment_firewall_ipset` | `<scope>/<name>` | `cluster/my-ipset` |
| `proxmox_virtual_environment_firewall_options` | `<scope>` | `cluster` or `vm/pve/100` |
| `proxmox_virtual_environment_firewall_rule` | `<scope>/<pos>` | `cluster/0` or `vm/pve/200/3` |
| `proxmox_virtual_environment_firewall_security_group` | `<name>` | `my-sg` |
| `proxmox_virtual_environment_firewall_security_group_rule` | `<group>/<pos>` | `my-sg/0` |
| `proxmox_virtual_environment_group` | `<group_id>` | `operators` |
| `proxmox_virtual_environment_ha_group` | `<group_name>` | `production` |
| `proxmox_virtual_environment_ha_resource` | `<sid>` | `vm:100` |
| `proxmox_virtual_environment_hardware_mapping_pci` | `<name>` | `my-gpu` |
| `proxmox_virtual_environment_hardware_mapping_usb` | `<name>` | `my-usb` |
| `proxmox_virtual_environment_hosts` | `<node_name>` | `pve` |
| `proxmox_virtual_environment_metrics_server` | `<name>` | `influxdb-prod` |
| `proxmox_virtual_environment_network_interface` | `<node_name>/<iface>` | `pve/vmbr0` |
| `proxmox_virtual_environment_pool` | `<pool_id>` | `production` |
| `proxmox_virtual_environment_realm` | `<realm>` | `my-ldap` |
| `proxmox_virtual_environment_replication` | `<vmid>-<job_id>` | `100-0` |
| `proxmox_virtual_environment_role` | `<role_id>` | `VMOperator` |
| `proxmox_virtual_environment_sdn_subnet` | `<vnet>/<subnet-cidr>` | `appvnet/10.100.0.0/24` |
| `proxmox_virtual_environment_sdn_vnet` | `<vnet>` | `appvnet` |
| `proxmox_virtual_environment_sdn_zone` | `<zone>` | `vxlan1` |
| `proxmox_virtual_environment_storage` | `<storage_id>` | `local-backup` |
| `proxmox_virtual_environment_time` | `<node_name>` | `pve` |
| `proxmox_virtual_environment_user` | `<user_id>` | `alice@pve` |
| `proxmox_virtual_environment_user_token` | `<user_id>/<token_id>` | `root@pam/my-token` |
| `proxmox_virtual_environment_vm` | `<node_name>/<vmid>` | `pve/200` |
| `proxmox_virtual_environment_vm_snapshot` | `<node_name>/<vmid>/<snap_name>` | `pve/200/mysnapshot` |

---

## Development

### Prerequisites

- Go 1.24 or later
- GNU Make
- [golangci-lint](https://golangci-lint.run/usage/install/) (for linting)

### Building

```shell
git clone https://github.com/Seguret/terraform-provider-proxmox.git
cd terraform-provider-proxmox
make build
```

This produces a `terraform-provider-proxmox` binary in the project directory.

### Local Installation

The `install` target compiles the provider and copies it to the local Terraform plugin cache so it can be used without a registry:

```shell
make install
```

The binary is installed to:

```
~/.terraform.d/plugins/registry.terraform.io/Seguret/proxmox/0.1.0/<OS>_<ARCH>/
```

To use the locally installed provider, add a `dev_overrides` block to `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/Seguret/proxmox" = "/home/<user>/.terraform.d/plugins/registry.terraform.io/Seguret/proxmox/0.1.0/<OS>_<ARCH>"
  }
  direct {}
}
```

### Running Unit Tests

```shell
make test
# or directly:
go test ./... -v
```

### Running Acceptance Tests

Acceptance tests create real resources on a live Proxmox VE instance. Set the required environment variables before running:

```shell
export TF_ACC=1
export PROXMOX_VE_ENDPOINT="https://pve.example.com:8006"
export PROXMOX_VE_API_TOKEN="root@pam!terraform=<uuid>"
export PROXMOX_VE_INSECURE="true"   # only if using a self-signed certificate

make testacc
```

> **Warning:** Acceptance tests create and destroy real resources. Use a dedicated test environment, not a production cluster.

### Available Make Targets

| Target | Description |
|---|---|
| `make build` | Compile the provider binary |
| `make install` | Compile and install to the local plugin cache |
| `make test` | Run unit tests |
| `make testacc` | Run acceptance tests (requires `TF_ACC=1`) |
| `make lint` | Run golangci-lint |
| `make fmt` | Format all Go source files |
| `make clean` | Remove the compiled binary |

### Debug Mode

Start the provider in debug mode to attach a debugger (e.g. Delve):

```shell
./terraform-provider-proxmox -debug
```

The provider prints a `TF_REATTACH_PROVIDERS` value. Export it in your shell, then run `terraform apply` normally and the provider process will be reused.

### Project Structure

```
internal/
  client/          # HTTP client, API models, and endpoint methods
  provider/        # Provider registration (resources and data sources)
  resources/       # One directory per resource type
    acl/
    acme_account/
    acme_plugin/
    apt_repository/
    backup/
    certificate/
    cluster_options/
    container/
    container_snapshot/
    dns/
    download_file/
    firewall_alias/
    firewall_ipset/
    firewall_options/
    firewall_rule/
    firewall_security_group/
    firewall_security_group_rule/
    group/
    ha_group/
    ha_resource/
    hardware_mapping_pci/
    hardware_mapping_usb/
    hosts/
    metrics_server/
    network_interface/
    pool/
    realm/
    replication/
    role/
    sdn_subnet/
    sdn_vnet/
    sdn_zone/
    storage/
    time/
    user/
    user_token/
    vm/
    vm_snapshot/
  datasources/     # One directory per data source type
```

Each resource follows the same pattern:
- `resource.go` — schema definition, CRUD operations, and `ImportState`
- A single `readIntoModel` function called from `Create`, `Read`, `Update`, and `ImportState`

### Contributing

Contributions are welcome. Follow these steps:

1. Fork the repository and create a branch from `main`.
2. Write or update tests for any new or changed behavior.
3. Ensure `make lint` and `make test` pass without errors.
4. Open a pull request with a clear description of the change and its motivation.

For significant changes, open an issue first to discuss the approach before investing time in the implementation.

---

## License

This project is distributed under the MIT License. See the [LICENSE](LICENSE) file for details.
