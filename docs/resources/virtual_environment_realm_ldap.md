---
page_title: "proxmox_virtual_environment_realm_ldap Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages an LDAP/Active Directory authentication realm in Proxmox VE.
---

# proxmox_virtual_environment_realm_ldap (Resource)

Manages an LDAP or Active Directory authentication realm in Proxmox VE.

## Example Usage

```terraform
resource "proxmox_virtual_environment_realm_ldap" "example" {
  realm    = "company-ldap"
  server1  = "ldap.example.com"
  base_dn  = "dc=example,dc=com"
  user_attr = "uid"
  comment  = "Company LDAP"
}

resource "proxmox_virtual_environment_realm_ldap" "active_directory" {
  realm     = "company-ad"
  server1   = "ad.example.com"
  port      = 636
  base_dn   = "dc=example,dc=com"
  user_attr = "sAMAccountName"
  domain    = "EXAMPLE"
  secure    = true
  comment   = "Company Active Directory"
}
```

## Schema

### Required

- `base_dn` (String) The LDAP base distinguished name for user searches.
- `realm` (String) The realm identifier (e.g. `my-ldap`). Changing this forces a new resource.
- `server1` (String) The primary LDAP server hostname or IP address.

### Optional

- `comment` (String) A comment for the realm.
- `default` (Boolean) Whether this realm is the default login realm. Defaults to `false`.
- `domain` (String) The optional AD domain name.
- `port` (Number) The LDAP server port. Defaults to `389`.
- `secure` (Boolean) Whether to use LDAPS (TLS/SSL). Defaults to `false`.
- `user_attr` (String) The LDAP attribute used to identify users (e.g. `uid` or `sAMAccountName`).

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import proxmox_virtual_environment_realm_ldap.example company-ldap
```
