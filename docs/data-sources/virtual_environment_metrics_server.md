# proxmox_virtual_environment_metrics_server (Data Source)

Retrieves information about a specific Proxmox VE external metrics server.




## Example Usage

```hcl
data "proxmox_virtual_environment_metrics_server" "example" {
  server_id = "influx"
}
```

## Schema

### Required

- `server_id` (String) The metrics server identifier to look up.

### Read-Only

- `disable` (Boolean) Whether the metrics server integration is disabled.
- `id` (String) The metrics server identifier.
- `port` (Number) The port the metrics server listens on.
- `server` (String) The hostname or IP address of the metrics server.
- `type` (String) The metrics server type (e.g. `influxdb`, `graphite`).
