# proxmox_virtual_environment_container_rrd (Data Source)

Retrieves RRD performance data for a Proxmox VE LXC container.




## Example Usage

```hcl
data "proxmox_virtual_environment_container_rrd" "example" {
  node_name = "pve"
  vmid      = 100
  timeframe = "hour"
}
```

## Schema

### Required

- `node_name` (String) The name of the node.
- `timeframe` (String) The timeframe for RRD data (`hour`, `day`, `week`, `month`, `year`).
- `vmid` (Number) The container ID.

### Read-Only

- `data_points` (Attributes List) RRD data points. (see [below for nested schema](#nestedatt--data_points))
- `id` (String) The datasource identifier.

<a id="nestedatt--data_points"></a>
### Nested Schema for `data_points`

Read-Only:

- `cpu` (Number) CPU usage (0-1).
- `diskread` (Number) Disk read (bytes/s).
- `diskwrite` (Number) Disk write (bytes/s).
- `loadavg` (Number) 1-minute load average.
- `maxcpu` (Number) Number of CPUs.
- `maxmem` (Number) Total memory in bytes.
- `mem` (Number) Memory used in bytes.
- `netin` (Number) Network in (bytes/s).
- `netout` (Number) Network out (bytes/s).
- `swaptotal` (Number) Total swap in bytes.
- `swapused` (Number) Used swap in bytes.
- `time` (Number) Unix timestamp.
