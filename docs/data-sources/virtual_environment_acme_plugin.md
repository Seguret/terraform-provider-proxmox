# proxmox_virtual_environment_acme_plugin (Data Source)

Retrieves a Proxmox VE ACME plugin.

## Example Usage

```terraform
data "proxmox_virtual_environment_acme_plugin" "plugin" {
  plugin_id = "my-dns-plugin"
}

output "acme_plugin_type" {
  value = data.proxmox_virtual_environment_acme_plugin.plugin.type
}
```


## Schema

### Required

- `plugin_id` (String) The ACME plugin identifier.

### Read-Only

- `api` (String) The DNS provider API identifier.
- `data` (String) The raw plugin configuration data.
- `id` (String) The unique identifier for this data source.
- `type` (String) The plugin type (e.g. 'dns' or 'standalone').
