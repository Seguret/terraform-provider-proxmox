---
page_title: "proxmox_virtual_environment_notification_endpoint_smtp Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE SMTP notification endpoint.
---

# proxmox_virtual_environment_notification_endpoint_smtp (Resource)

Manages a Proxmox VE SMTP notification endpoint.

## Example Usage

### Basic SMTP Endpoint

```terraform
resource "proxmox_virtual_environment_notification_endpoint_smtp" "basic" {
  name   = "smtp-basic"
  server = "mail.example.com"
  from   = "proxmox@example.com"
}
```

### SMTP with TLS

```terraform
resource "proxmox_virtual_environment_notification_endpoint_smtp" "secure" {
  name     = "smtp-secure"
  server   = "smtp.gmail.com"
  from     = "alerts@example.com"
  mode     = "tls"
  port     = 465
  username = "alerts@gmail.com"
  password = "app-specific-password"
  comment  = "Gmail SMTP endpoint"
}
```

### SMTP with User Notifications

```terraform
resource "proxmox_virtual_environment_notification_endpoint_smtp" "users" {
  name          = "smtp-users"
  server        = "mail.example.com"
  from          = "proxmox-alerts@example.com"
  mode          = "starttls"
  port          = 587
  mailto_user   = ["admin@pve", "operator@pve"]
  comment       = "Notify specific users"
}
```

### SMTP with Email List

```terraform
resource "proxmox_virtual_environment_notification_endpoint_smtp" "recipients" {
  name    = "smtp-alerts"
  server  = "mail.example.com"
  from    = "system@example.com"
  mode    = "starttls"
  port    = 587
  mailto  = ["ops-team@example.com", "admin@example.com", "security@example.com"]
  comment = "Multi-recipient alert notifications"
}
```

## Schema

### Required

- `from` (String) The sender email address.
- `name` (String) The name of the SMTP endpoint.
- `server` (String) The SMTP server hostname or IP address.

### Optional

- `comment` (String) Comment for the SMTP endpoint.
- `disable` (Boolean) Whether the SMTP endpoint is disabled.
- `mailto` (List of String) List of email addresses to send notifications to.
- `mailto_user` (List of String) List of users to send notifications to (by Proxmox user ID).
- `mode` (String) The SMTP encryption mode: insecure, starttls, or tls.
- `password` (String, Sensitive) The SMTP password for authentication.
- `port` (Number) The SMTP server port.
- `username` (String) The SMTP username for authentication.

### Read-Only

- `id` (String) The ID of this resource.
