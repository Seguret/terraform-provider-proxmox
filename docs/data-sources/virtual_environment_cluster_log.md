# proxmox_virtual_environment_cluster_log (Data Source)

Retrieves the Proxmox VE cluster event log.




## Example Usage

```hcl
data "proxmox_virtual_environment_cluster_log" "example" {
  max = 100
}
```

## Schema

### Optional

- `max` (Number) Maximum number of log entries to return (default `50`).

### Read-Only

- `entries` (Attributes List) Cluster log entries. (see [below for nested schema](#nestedatt--entries))
- `id` (String) The datasource identifier.

<a id="nestedatt--entries"></a>
### Nested Schema for `entries`

Read-Only:

- `gid` (Number) Group ID.
- `msg` (String) Log message.
- `node` (String) Node name.
- `pid` (Number) Process ID.
- `severity` (String) Severity level.
- `tag` (String) Log tag/service.
- `time` (Number) Unix timestamp.
- `uid` (Number) User ID.
- `user_id` (String) User who triggered the event.
