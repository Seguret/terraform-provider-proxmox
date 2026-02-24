# Provider Documentation

This directory contains the documentation for the Proxmox Terraform Provider in the format required by the Terraform Registry.

## Structure

```
docs/
├── index.md                    # Provider homepage
├── data-sources/               # Data source documentation
│   ├── vm.md
│   ├── vms.md
│   ├── container.md
│   ├── containers.md
│   ├── nodes.md
│   ├── datastores.md
│   └── ...
└── resources/                  # Resource documentation
    ├── vm.md
    ├── container.md
    ├── storage.md
    ├── pool.md
    ├── user.md
    ├── user_token.md
    ├── group.md
    ├── role.md
    └── ...
```

## Documentation Status

✅ **Completed:**
- index.md - Provider homepage with complete examples
- Core compute resources (vm, container)
- Access control resources (user, user_token, group, role)
- Storage resources (storage, pool)
- Core data sources (vms, vm, containers, nodes, datastores, cluster_status)

📝 **To Generate:**

The remaining 50+ resource and data source documentation files can be generated automatically using the `terraform-plugin-docs` tool.

## Generating Documentation

### Prerequisites

Install the tfplugindocs tool:

```bash
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
```

### Generate All Docs

From the provider root directory:

```bash
tfplugindocs generate
```

This will:
1. Parse resource and data source schemas from the Go code
2. Use example templates from `examples/` directory
3. Generate markdown files for all resources and data sources
4. Place them in the `docs/` directory

### Generate Specific Resource/Data Source

To generate docs for a specific resource:

```bash
tfplugindocs generate --provider-name proxmox --rendered-provider-name "Proxmox" --rendered-website-dir docs
```

## Adding Examples

Create example files in the `examples/` directory following this structure:

```
examples/
├── resources/
│   └── proxmox_vm/
│       ├── resource.tf          # Basic example
│       └── import.sh            # Import example
└── data-sources/
    └── proxmox_vms/
        └── data-source.tf       # Query example
```

Example for `examples/resources/proxmox_network_interface/resource.tf`:

```terraform
resource "proxmox_network_interface" "example" {
  node      = "pve"
  interface = "vmbr1"
  type      = "bridge"
  
  autostart = true
  bridge_ports = ["ens19"]
  
  address  = "192.168.100.1"
  netmask  = "255.255.255.0"
  gateway  = "192.168.100.254"
}
```

## Documentation Resources Needed

### High Priority Resources (Need Docs)

**Networking:**
- network_interface
- sdn_zone, sdn_vnet, sdn_subnet
- sdn_controller, sdn_dns, sdn_ipam

**Firewall:**
- firewall_rule, firewall_alias
- firewall_ipset, firewall_security_group
- firewall_options

**HA & Cluster:**
- ha_group, ha_resource
- replication, backup
- cluster_options, cluster_job_realm_sync

**Node Management:**
- node_config, node_service
- node_disk_* (directory, lvm, lvmthin, zfs)
- node_subscription, node_acme_certificate

**Certificates & ACME:**
- acme_account, acme_plugin
- certificate

**Notifications:**
- notification_endpoint_* (sendmail, smtp, gotify, webhook)
- notification_filter, notification_matcher

**Miscellaneous:**
- realm, acl
- dns, hosts, time
- apt_repository, download_file
- metrics_server
- hardware_mapping_pci, hardware_mapping_usb

### Data Sources (Need Docs)

40+ data sources including:
- Cluster: cluster_resources, cluster_tasks, cluster_config, cluster_nextid
- Node: node, node_disks, node_hardware_pci, node_hardware_usb, node_journal, node_services, node_smart, node_tasks
- Compute: vm_snapshots, container_snapshots
- Storage: storage_content, backups
- Access: users, groups, roles, user_permissions
- HA: ha_status, ha_groups, ha_resources
- SDN: sdn_vnets, sdn_zones
- Hardware: hardware_mapping_pci, hardware_mapping_usb
- ACME: acme_directories, openid_config
- APT: apt_changelog, apt_repositories, apt_versions
- Misc: pools, version

## Documentation Standards

Each resource/data source doc should include:

1. **Frontmatter**: page_title, subcategory, description
2. **Title and Description**: Clear explanation of what the resource does
3. **Example Usage**: At least one working example
4. **Schema Section**: 
   - Required arguments
   - Optional arguments
   - Read-only attributes
   - Nested block schemas
5. **Import Section**: How to import existing resources (if applicable)

## Contributing

When adding new resources or data sources:

1. Add schema descriptions in the Go code using `Description` fields
2. Create example files in `examples/`
3. Run `tfplugindocs generate` to create/update docs
4. Review generated docs for accuracy
5. Submit PR with both code and documentation
