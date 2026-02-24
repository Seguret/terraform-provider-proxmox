---
page_title: "proxmox_virtual_environment_realm Resource - terraform-provider-proxmox"
subcategory: ""
description: |-
  Manages a Proxmox VE authentication realm.
---

# proxmox_virtual_environment_realm (Resource)

Manages a Proxmox VE authentication realm.

## Example Usage

### LDAP Realm

```terraform
resource "proxmox_virtual_environment_realm" "ldap" {
  realm      = "company-ldap"
  type       = "ldap"
  comment    = "Company LDAP directory"
  server1    = "ldap.example.com"
  port       = 389
  base_dn    = "ou=users,dc=example,dc=com"
  user_attr  = "uid"
  group_attr = "cn"
}
```

### Active Directory Realm

```terraform
resource "proxmox_virtual_environment_realm" "active_directory" {
  realm      = "company-ad"
  type       = "ad"
  comment    = "Company Active Directory"
  server1    = "dc1.example.com"
  domain     = "example.com"
  port       = 389
  base_dn    = "cn=users,dc=example,dc=com"
  auto_create = true
}
```

### OpenID Connect Realm

```terraform
resource "proxmox_virtual_environment_realm" "oidc" {
  realm       = "keycloak-oidc"
  type        = "openid"
  comment     = "Keycloak OIDC authentication"
  issuer_url  = "https://keycloak.example.com/auth/realms/master"
  client_id   = "proxmox-client"
  client_key  = "your-secret-key"
  username_claim = "preferred_username"
  auto_create = true
}
```

### Default Realm Configuration

```terraform
resource "proxmox_virtual_environment_realm" "ldap_default" {
  realm       = "primary-ldap"
  type        = "ldap"
  comment     = "Primary LDAP realm"
  server1     = "ldap1.example.com"
  server2     = "ldap2.example.com"
  port        = 389
  base_dn     = "ou=users,dc=example,dc=com"
  default     = true
  auto_create = true
}
```


## Schema

### Required

- `realm` (String) The realm identifier (e.g., 'my-ldap').
- `type` (String) The realm type (pam, pve, ad, ldap, openid).

### Optional

- `auto_create` (Boolean) Automatically create users on first login.
- `base_dn` (String) LDAP base distinguished name.
- `bind_dn` (String) LDAP bind distinguished name.
- `client_id` (String) OpenID Connect client ID.
- `client_key` (String, Sensitive) OpenID Connect client secret.
- `comment` (String) A comment for the realm.
- `default` (Boolean) Whether this is the default realm.
- `issuer_url` (String) OpenID Connect issuer URL.
- `password` (String, Sensitive) LDAP bind password.
- `port` (Number) Server port (for ldap/ad).
- `secure` (Boolean) Use LDAPS/TLS.
- `server1` (String) Primary server address (for ldap/ad).
- `server2` (String) Secondary server address (for ldap/ad).
- `tfa` (String) Two-factor authentication provider.
- `user_attr` (String) LDAP user attribute name.
- `username_claim` (String) OpenID Connect claim used as the username.
- `verify` (Boolean) Verify server TLS certificate.

### Read-Only

- `id` (String) The ID of this resource.
