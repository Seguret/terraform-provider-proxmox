---
page_title: "proxmox_virtual_environment_cluster_options Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages cluster-wide Proxmox VE options.
---

# proxmox_virtual_environment_cluster_options (Resource)

Manages cluster-wide Proxmox VE options.

## Example Usage

### Basic Cluster Options

```terraform
resource "proxmox_virtual_environment_cluster_options" "cluster" {
  email_from = "proxmox@example.com"
  keyboard   = "en-us"
  language   = "en"
}
```

### HA Configuration

```terraform
resource "proxmox_virtual_environment_cluster_options" "ha_cluster" {
  email_from          = "proxmox-alerts@example.com"
  ha_shutdown_policy  = "failover"
  keyboard            = "en-us"
  language            = "en"
  max_workers         = 4
}
```

### Complete Cluster Options

```terraform
resource "proxmox_virtual_environment_cluster_options" "production" {
  email_from           = "proxmox@example.com"
  ha_shutdown_policy   = "conditional"
  http_proxy           = "http://proxy.example.com:3128"
  keyboard             = "en-us"
  language             = "en"
  max_workers          = 8
  migration_type       = "secure"
  migration_unsecure   = false
}
```

### Development Cluster

```terraform
resource "proxmox_virtual_environment_cluster_options" "development" {
  email_from           = "dev-proxmox@example.com"
  ha_shutdown_policy   = "freeze"
  keyboard             = "en-us"
  language             = "en"
  max_workers          = 2
  migration_type       = "websocket"
}
```


## Schema

### Optional

- `email_from` (String) Email address used as sender for notifications.
- `ha_shutdown_policy` (String) HA shutdown policy: 'freeze', 'failover', 'conditional', or 'migrate'.
- `http_proxy` (String) HTTP proxy URL for the cluster (used for apt, etc.).
- `keyboard` (String) Default keyboard layout (e.g., 'en-us', 'de', 'fr').
- `language` (String) Default web UI language.
- `max_workers` (Number) Maximum number of workers for bulk operations.
- `migration_type` (String) Migration type: 'secure' or 'insecure' or 'websocket'.
- `migration_unsecure` (Boolean) Whether to allow unsecured migrations (non-TLS).

### Read-Only

- `id` (String) The ID of this resource.
