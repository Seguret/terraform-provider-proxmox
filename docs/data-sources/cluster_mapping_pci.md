# proxmox_cluster_mapping_pci (Data Source)

Retrieves PCI hardware mappings from the cluster.




## Schema

### Read-Only

- `id` (String) Data source identifier.
- `mappings` (Attributes List) List of PCI hardware mappings. (see [below for nested schema](#nestedatt--mappings))

<a id="nestedatt--mappings"></a>
### Nested Schema for `mappings`

Read-Only:

- `description` (String) Mapping description.
- `id` (String) Mapping ID.
- `map` (List of String) List of PCI device mappings.
