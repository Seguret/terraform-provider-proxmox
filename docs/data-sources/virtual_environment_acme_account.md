# proxmox_virtual_environment_acme_account (Data Source)

Retrieves a Proxmox VE ACME account.

## Example Usage

```terraform
data "proxmox_virtual_environment_acme_account" "account" {
  name = "my-acme-account"
}

output "acme_account_emails" {
  value = data.proxmox_virtual_environment_acme_account.account.email
}
```


## Schema

### Required

- `name` (String) The ACME account name.

### Read-Only

- `created_at` (String) The account creation timestamp.
- `directory` (String) The ACME directory URL.
- `email` (List of String) The contact email addresses for this ACME account.
- `id` (String) The unique identifier for this data source.
- `tos` (String) The Terms of Service URL.
