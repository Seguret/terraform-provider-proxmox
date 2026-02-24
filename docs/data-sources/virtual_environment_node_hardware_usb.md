# proxmox_virtual_environment_node_hardware_usb (Data Source)

Retrieves the list of USB devices on a Proxmox VE node.

## Example Usage

```hcl
data "proxmox_virtual_environment_node_hardware_usb" "pve" {
  node_name = "pve"
}

output "cpu_model" {
  value = data.proxmox_virtual_environment_node_hardware_usb.pve
}
```

## Schema

### Required

- `node_name` (String) The name of the Proxmox VE node.

### Read-Only

- `devices` (Attributes List) The list of USB devices on the node. (see [below for nested schema](#nestedatt--devices))
- `id` (String) Placeholder identifier.

<a id="nestedatt--devices"></a>
### Nested Schema for `devices`

Read-Only:

- `bus_num` (Number) The USB bus number.
- `dev_num` (Number) The USB device number.
- `manufacturer` (String) The device manufacturer name.
- `prod_id` (String) The USB product ID.
- `product` (String) The product name.
- `serialnumber` (String) The device serial number.
- `speed` (String) The USB device speed (e.g., '480', '5000').
- `vendor_id` (String) The USB vendor ID.
