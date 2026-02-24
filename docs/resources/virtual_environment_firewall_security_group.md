---
page_title: "proxmox_virtual_environment_firewall_security_group Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE cluster firewall security group. Rules within the group are managed by a separate resource.
---

# proxmox_virtual_environment_firewall_security_group (Resource)

Manages a Proxmox VE cluster firewall security group. Rules within the group are managed by a separate resource.

## Example Usage

### Basic Security Group

```terraform
resource "proxmox_virtual_environment_firewall_security_group" "web_servers" {
  name    = "web-servers"
  comment = "Security group for web server resources"
}
```

### Multiple Security Groups

```terraform
resource "proxmox_virtual_environment_firewall_security_group" "database" {
  name    = "database-servers"
  comment = "Database server security rules"
}

resource "proxmox_virtual_environment_firewall_security_group" "application" {
  name    = "app-servers"
  comment = "Application server security rules"
}

resource "proxmox_virtual_environment_firewall_security_group" "cache" {
  name    = "cache-layer"
  comment = "Caching tier security rules"
}
```

### Security Group with Associated Rules

```terraform
resource "proxmox_virtual_environment_firewall_security_group" "api_tier" {
  name    = "api-servers"
  comment = "API backend security group"
}

resource "proxmox_virtual_environment_firewall_security_group_rule" "api_http" {
  group   = proxmox_virtual_environment_firewall_security_group.api_tier.name
  action  = "ACCEPT"
  type    = "in"
  dport   = "8080"
  protocol = "tcp"
  comment = "Allow API traffic"
}

resource "proxmox_virtual_environment_firewall_security_group_rule" "api_https" {
  group    = proxmox_virtual_environment_firewall_security_group.api_tier.name
  action   = "ACCEPT"
  type     = "in"
  dport    = "8443"
  protocol = "tcp"
  comment  = "Allow secure API traffic"
}
```

## Schema

### Required

- `name` (String) The security group name.

### Optional

- `comment` (String) A comment for the security group.

### Read-Only

- `id` (String) The ID of this resource.
