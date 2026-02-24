# proxmox_virtual_environment_acme_accounts (Data Source)

Retrieves the list of Proxmox VE ACME accounts.

## Example Usage

```terraform
data "proxmox_virtual_environment_acme_accounts" "all" {}

output "acme_account_names" {
  value = [for a in data.proxmox_virtual_environment_acme_accounts.all.accounts : a.name]
}
```


## Schema

### Read-Only

- `accounts` (List of Object) The list of ACME accounts.
  - `name` (String) The ACME account name.
- `id` (String) The unique identifier for this data source.
