# proxmox_virtual_environment_hosts (Data Source)

Retrieves the `/etc/hosts` configuration of a Proxmox VE node.




## Example Usage

```hcl
data "proxmox_virtual_environment_hosts" "example" {
  node_name = "pve"
}
```

## Schema

### Required

- `node_name` (String) The node to retrieve hosts configuration from.

### Read-Only

- `digest` (String) The SHA1 digest of the current hosts content.
- `entries` (Attributes List) The list of host entries parsed from `/etc/hosts`. (see [below for nested schema](#nestedatt--entries))
- `id` (String) The hosts datasource identifier.

<a id="nestedatt--entries"></a>
### Nested Schema for `entries`

Read-Only:

- `address` (String) The IP address of the host entry.
- `hostnames` (List of String) The hostnames associated with this address.
