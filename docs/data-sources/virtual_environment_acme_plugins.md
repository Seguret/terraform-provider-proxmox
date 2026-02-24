# proxmox_virtual_environment_acme_plugins (Data Source)

Retrieves the list of Proxmox VE ACME plugins.

## Example Usage

```terraform
data "proxmox_virtual_environment_acme_plugins" "all" {}

output "acme_plugin_ids" {
  value = [for p in data.proxmox_virtual_environment_acme_plugins.all.plugins : p.plugin_id]
}
```


## Schema

### Read-Only

- `id` (String) The unique identifier for this data source.
- `plugins` (List of Object) The list of ACME plugins.
  - `plugin_id` (String) The ACME plugin identifier.
  - `type` (String) The plugin type (e.g. 'dns' or 'standalone').
