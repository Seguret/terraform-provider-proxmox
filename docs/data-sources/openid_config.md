# proxmox_openid_config (Data Source)

Retrieves OpenID Connect configuration for a Proxmox VE authentication realm.




## Schema

### Required

- `realm` (String) The authentication realm name (e.g., 'pve', 'openid-realm').

### Read-Only

- `authorization_endpoint` (String) The authorization endpoint URL.
- `id` (String) The ID of this resource.
- `issuer` (String) The OpenID issuer URL.
- `jwks_uri` (String) The JWKS (JSON Web Key Set) URI.
- `response_types_supported` (List of String) List of supported OAuth2 response types.
- `scopes_supported` (List of String) List of supported OAuth2 scopes.
- `token_endpoint` (String) The token endpoint URL.
- `userinfo_endpoint` (String) The userinfo endpoint URL.
