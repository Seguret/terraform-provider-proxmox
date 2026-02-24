# proxmox_virtual_environment_node_netstat (Data Source)

Retrieves network interface statistics for a Proxmox VE node.




## Example Usage

```hcl
data "proxmox_virtual_environment_node_netstat" "example" {
  node_name = "pve"
}
```

## Schema

### Required

- `node_name` (String) The name of the node.

### Read-Only

- `entries` (Attributes List) Network interface statistics entries. (see [below for nested schema](#nestedatt--entries))
- `id` (String) The datasource identifier.

<a id="nestedatt--entries"></a>
### Nested Schema for `entries`

Read-Only:

- `iface` (String) Interface name.
- `rx_bytes` (Number) Received bytes.
- `rx_drop` (Number) Received packets dropped.
- `rx_err` (Number) Receive errors.
- `rx_pkts` (Number) Received packets.
- `tx_bytes` (Number) Transmitted bytes.
- `tx_drop` (Number) Transmitted packets dropped.
- `tx_err` (Number) Transmit errors.
- `tx_pkts` (Number) Transmitted packets.
