package provider

import (
	"context"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"

	// Data sources
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/acme_directories"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/apt_changelog"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/apt_repositories"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/apt_versions"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/backups"
	cluster_ceph_status_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/cluster_ceph_status"
	cluster_config_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/cluster_config"
	cluster_mapping_pci_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/cluster_mapping_pci"
	cluster_mapping_usb_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/cluster_mapping_usb"
	cluster_nextid_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/cluster_nextid"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/cluster_resources"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/cluster_status"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/cluster_tasks"
	container_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/container"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/containers"
	container_snapshots_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/container_snapshots"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/datastores"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/groups"
	ha_groups_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/ha_groups"
	ha_resources_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/ha_resources"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/ha_status"
	hw_pci_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/hardware_mapping_pci"
	hw_usb_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/hardware_mapping_usb"
	network_interfaces_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/network_interfaces"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/node"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/node_disks"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/node_hardware_pci"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/node_hardware_usb"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/node_journal"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/node_services"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/node_smart"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/node_tasks"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/nodes"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/pools"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/roles"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/storage_content"
	sdn_vnets_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/sdn_vnets"
	sdn_zones_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/sdn_zones"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/users"
	user_permissions_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/user_permissions"
	openid_config_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/openid_config"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/version"
	vm_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/vm"
	vm_snapshots_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/vm_snapshots"
	vm_rrd_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/vm_rrd"
	"github.com/Seguret/terraform-provider-proxmox/internal/datasources/vms"
	node_syslog_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/node_syslog"
	node_rrd_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/node_rrd"
	node_netstat_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/node_netstat"
	node_scan_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/node_scan"
	cluster_log_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/cluster_log"
	cluster_backup_info_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/cluster_backup_info"
	container_rrd_ds "github.com/Seguret/terraform-provider-proxmox/internal/datasources/container_rrd"

	// Resources
	"github.com/Seguret/terraform-provider-proxmox/internal/resources/acl"
	acme_account "github.com/Seguret/terraform-provider-proxmox/internal/resources/acme_account"
	acme_plugin "github.com/Seguret/terraform-provider-proxmox/internal/resources/acme_plugin"
	apt_repository "github.com/Seguret/terraform-provider-proxmox/internal/resources/apt_repository"
	"github.com/Seguret/terraform-provider-proxmox/internal/resources/backup"
	"github.com/Seguret/terraform-provider-proxmox/internal/resources/certificate"
	cluster_job_realm_sync "github.com/Seguret/terraform-provider-proxmox/internal/resources/cluster_job_realm_sync"
	cluster_options "github.com/Seguret/terraform-provider-proxmox/internal/resources/cluster_options"
	container_res "github.com/Seguret/terraform-provider-proxmox/internal/resources/container"
	container_snapshot "github.com/Seguret/terraform-provider-proxmox/internal/resources/container_snapshot"
	"github.com/Seguret/terraform-provider-proxmox/internal/resources/dns"
	download_file "github.com/Seguret/terraform-provider-proxmox/internal/resources/download_file"
	firewall_alias "github.com/Seguret/terraform-provider-proxmox/internal/resources/firewall_alias"
	firewall_ipset "github.com/Seguret/terraform-provider-proxmox/internal/resources/firewall_ipset"
	firewall_options "github.com/Seguret/terraform-provider-proxmox/internal/resources/firewall_options"
	firewall_rule "github.com/Seguret/terraform-provider-proxmox/internal/resources/firewall_rule"
	firewall_security_group "github.com/Seguret/terraform-provider-proxmox/internal/resources/firewall_security_group"
	firewall_security_group_rule "github.com/Seguret/terraform-provider-proxmox/internal/resources/firewall_security_group_rule"
	"github.com/Seguret/terraform-provider-proxmox/internal/resources/group"
	ha_group "github.com/Seguret/terraform-provider-proxmox/internal/resources/ha_group"
	ha_resource "github.com/Seguret/terraform-provider-proxmox/internal/resources/ha_resource"
	hardware_mapping_pci "github.com/Seguret/terraform-provider-proxmox/internal/resources/hardware_mapping_pci"
	hardware_mapping_usb "github.com/Seguret/terraform-provider-proxmox/internal/resources/hardware_mapping_usb"
	"github.com/Seguret/terraform-provider-proxmox/internal/resources/hosts"
	metrics_server "github.com/Seguret/terraform-provider-proxmox/internal/resources/metrics_server"
	network_interface "github.com/Seguret/terraform-provider-proxmox/internal/resources/network_interface"
	node_acme_certificate "github.com/Seguret/terraform-provider-proxmox/internal/resources/node_acme_certificate"
	node_config "github.com/Seguret/terraform-provider-proxmox/internal/resources/node_config"
	node_disk_directory "github.com/Seguret/terraform-provider-proxmox/internal/resources/node_disk_directory"
	node_disk_lvm "github.com/Seguret/terraform-provider-proxmox/internal/resources/node_disk_lvm"
	node_disk_lvmthin "github.com/Seguret/terraform-provider-proxmox/internal/resources/node_disk_lvmthin"
	node_disk_zfs "github.com/Seguret/terraform-provider-proxmox/internal/resources/node_disk_zfs"
	node_service "github.com/Seguret/terraform-provider-proxmox/internal/resources/node_service"
	node_subscription "github.com/Seguret/terraform-provider-proxmox/internal/resources/node_subscription"
	notification_endpoint_gotify "github.com/Seguret/terraform-provider-proxmox/internal/resources/notification_endpoint_gotify"
	notification_endpoint_sendmail "github.com/Seguret/terraform-provider-proxmox/internal/resources/notification_endpoint_sendmail"
	notification_endpoint_smtp "github.com/Seguret/terraform-provider-proxmox/internal/resources/notification_endpoint_smtp"
	notification_endpoint_webhook "github.com/Seguret/terraform-provider-proxmox/internal/resources/notification_endpoint_webhook"
	notification_filter "github.com/Seguret/terraform-provider-proxmox/internal/resources/notification_filter"
	notification_matcher "github.com/Seguret/terraform-provider-proxmox/internal/resources/notification_matcher"
	"github.com/Seguret/terraform-provider-proxmox/internal/resources/pool"
	"github.com/Seguret/terraform-provider-proxmox/internal/resources/realm"
	"github.com/Seguret/terraform-provider-proxmox/internal/resources/replication"
	"github.com/Seguret/terraform-provider-proxmox/internal/resources/role"
	sdn_controller "github.com/Seguret/terraform-provider-proxmox/internal/resources/sdn_controller"
	sdn_dns "github.com/Seguret/terraform-provider-proxmox/internal/resources/sdn_dns"
	sdn_ipam "github.com/Seguret/terraform-provider-proxmox/internal/resources/sdn_ipam"
	sdn_subnet "github.com/Seguret/terraform-provider-proxmox/internal/resources/sdn_subnet"
	sdn_vnet "github.com/Seguret/terraform-provider-proxmox/internal/resources/sdn_vnet"
	sdn_zone "github.com/Seguret/terraform-provider-proxmox/internal/resources/sdn_zone"
	"github.com/Seguret/terraform-provider-proxmox/internal/resources/storage"
	node_time "github.com/Seguret/terraform-provider-proxmox/internal/resources/time"
	"github.com/Seguret/terraform-provider-proxmox/internal/resources/user"
	user_password "github.com/Seguret/terraform-provider-proxmox/internal/resources/user_password"
	user_tfa "github.com/Seguret/terraform-provider-proxmox/internal/resources/user_tfa"
	user_token "github.com/Seguret/terraform-provider-proxmox/internal/resources/user_token"
	"github.com/Seguret/terraform-provider-proxmox/internal/resources/vm"
	vm_snapshot "github.com/Seguret/terraform-provider-proxmox/internal/resources/vm_snapshot"

	// Ceph
	ceph_fs   "github.com/Seguret/terraform-provider-proxmox/internal/resources/ceph_fs"
	ceph_mds  "github.com/Seguret/terraform-provider-proxmox/internal/resources/ceph_mds"
	ceph_mgr  "github.com/Seguret/terraform-provider-proxmox/internal/resources/ceph_mgr"
	ceph_mon  "github.com/Seguret/terraform-provider-proxmox/internal/resources/ceph_mon"
	ceph_osd  "github.com/Seguret/terraform-provider-proxmox/internal/resources/ceph_osd"
	ceph_pool "github.com/Seguret/terraform-provider-proxmox/internal/resources/ceph_pool"

	// Firewall (cluster/node level)
	cluster_firewall "github.com/Seguret/terraform-provider-proxmox/internal/resources/cluster_firewall"
	node_firewall    "github.com/Seguret/terraform-provider-proxmox/internal/resources/node_firewall"

	// Pool membership
	pool_membership "github.com/Seguret/terraform-provider-proxmox/internal/resources/pool_membership"

	// SDN applier + type-specific zones
	sdn_applier    "github.com/Seguret/terraform-provider-proxmox/internal/resources/sdn_applier"
	sdn_zone_evpn  "github.com/Seguret/terraform-provider-proxmox/internal/resources/sdn_zone_evpn"
	sdn_zone_qinq  "github.com/Seguret/terraform-provider-proxmox/internal/resources/sdn_zone_qinq"
	sdn_zone_simple "github.com/Seguret/terraform-provider-proxmox/internal/resources/sdn_zone_simple"
	sdn_zone_vlan  "github.com/Seguret/terraform-provider-proxmox/internal/resources/sdn_zone_vlan"
	sdn_zone_vxlan "github.com/Seguret/terraform-provider-proxmox/internal/resources/sdn_zone_vxlan"

	// Storage type-specific
	storage_cifs      "github.com/Seguret/terraform-provider-proxmox/internal/resources/storage_cifs"
	storage_directory "github.com/Seguret/terraform-provider-proxmox/internal/resources/storage_directory"
	storage_lvm       "github.com/Seguret/terraform-provider-proxmox/internal/resources/storage_lvm"
	storage_lvmthin   "github.com/Seguret/terraform-provider-proxmox/internal/resources/storage_lvmthin"
	storage_nfs       "github.com/Seguret/terraform-provider-proxmox/internal/resources/storage_nfs"
	storage_pbs       "github.com/Seguret/terraform-provider-proxmox/internal/resources/storage_pbs"
	storage_zfspool   "github.com/Seguret/terraform-provider-proxmox/internal/resources/storage_zfspool"

	// Network type-specific
	network_linux_bridge "github.com/Seguret/terraform-provider-proxmox/internal/resources/network_linux_bridge"
	network_linux_vlan   "github.com/Seguret/terraform-provider-proxmox/internal/resources/network_linux_vlan"

	// APT standard repositories
	apt_standard_repo "github.com/Seguret/terraform-provider-proxmox/internal/resources/apt_standard_repository"

	// LDAP realm
	realm_ldap "github.com/Seguret/terraform-provider-proxmox/internal/resources/realm_ldap"

	// Hardware mapping dir
	hardware_mapping_dir "github.com/Seguret/terraform-provider-proxmox/internal/resources/hardware_mapping_dir"

	// File resource
	file_res "github.com/Seguret/terraform-provider-proxmox/internal/resources/file"

	// OCI image
	oci_image "github.com/Seguret/terraform-provider-proxmox/internal/resources/oci_image"

	// Cloned VM
	cloned_vm "github.com/Seguret/terraform-provider-proxmox/internal/resources/cloned_vm"

	// New datasources
	acme_account_ds       "github.com/Seguret/terraform-provider-proxmox/internal/datasources/acme_account"
	acme_accounts_ds      "github.com/Seguret/terraform-provider-proxmox/internal/datasources/acme_accounts"
	acme_plugin_ds        "github.com/Seguret/terraform-provider-proxmox/internal/datasources/acme_plugin"
	acme_plugins_ds       "github.com/Seguret/terraform-provider-proxmox/internal/datasources/acme_plugins"
	apt_std_repo_ds       "github.com/Seguret/terraform-provider-proxmox/internal/datasources/apt_standard_repository"
	dns_ds                "github.com/Seguret/terraform-provider-proxmox/internal/datasources/dns"
	file_ds               "github.com/Seguret/terraform-provider-proxmox/internal/datasources/file"
	group_ds              "github.com/Seguret/terraform-provider-proxmox/internal/datasources/group"
	ha_group_ds           "github.com/Seguret/terraform-provider-proxmox/internal/datasources/ha_group"
	ha_resource_ds        "github.com/Seguret/terraform-provider-proxmox/internal/datasources/ha_resource"
	hosts_ds              "github.com/Seguret/terraform-provider-proxmox/internal/datasources/hosts"
	metrics_server_ds     "github.com/Seguret/terraform-provider-proxmox/internal/datasources/metrics_server"
	pool_ds               "github.com/Seguret/terraform-provider-proxmox/internal/datasources/pool"
	role_ds               "github.com/Seguret/terraform-provider-proxmox/internal/datasources/role"
	sdn_vnet_ds           "github.com/Seguret/terraform-provider-proxmox/internal/datasources/sdn_vnet"
	sdn_subnet_ds         "github.com/Seguret/terraform-provider-proxmox/internal/datasources/sdn_subnet"
	time_ds               "github.com/Seguret/terraform-provider-proxmox/internal/datasources/time"
	user_ds               "github.com/Seguret/terraform-provider-proxmox/internal/datasources/user"
)

var _ provider.Provider = &ProxmoxProvider{}

// ProxmoxProvider is the root terraform provider for Proxmox VE.
type ProxmoxProvider struct {
	version string
}

// ProxmoxProviderModel holds the provider-level configuration values.
type ProxmoxProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIToken types.String `tfsdk:"api_token"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

func New(v string) func() provider.Provider {
	return func() provider.Provider {
		return &ProxmoxProvider{
			version: v,
		}
	}
}

func (p *ProxmoxProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "proxmox"
	resp.Version = p.version
}

func (p *ProxmoxProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for managing Proxmox Virtual Environment resources.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "The Proxmox VE API endpoint URL (e.g., https://pve.example.com:8006). " +
					"Can also be set with the PROXMOX_VE_ENDPOINT environment variable.",
				Optional: true,
			},
			"api_token": schema.StringAttribute{
				Description: "The Proxmox VE API token (e.g., user@pam!tokenid=uuid). " +
					"Can also be set with the PROXMOX_VE_API_TOKEN environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"username": schema.StringAttribute{
				Description: "The Proxmox VE username for ticket-based auth (e.g., root@pam). " +
					"Can also be set with the PROXMOX_VE_USERNAME environment variable.",
				Optional: true,
			},
			"password": schema.StringAttribute{
				Description: "The Proxmox VE password for ticket-based auth. " +
					"Can also be set with the PROXMOX_VE_PASSWORD environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"insecure": schema.BoolAttribute{
				Description: "Whether to skip TLS certificate verification. Defaults to false. " +
					"Can also be set with the PROXMOX_VE_INSECURE environment variable.",
				Optional: true,
			},
		},
	}
}

func (p *ProxmoxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Proxmox VE client")

	var config ProxmoxProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// config values take precedence over environment variables
	endpoint := resolveString(config.Endpoint, "PROXMOX_VE_ENDPOINT")
	apiToken := resolveString(config.APIToken, "PROXMOX_VE_API_TOKEN")
	username := resolveString(config.Username, "PROXMOX_VE_USERNAME")
	password := resolveString(config.Password, "PROXMOX_VE_PASSWORD")
	insecure := resolveBool(config.Insecure, "PROXMOX_VE_INSECURE")

	if endpoint == "" {
		resp.Diagnostics.AddError(
			"Missing Proxmox VE Endpoint",
			"The provider requires the 'endpoint' attribute or the PROXMOX_VE_ENDPOINT environment variable.",
		)
		return
	}

	if apiToken == "" && (username == "" || password == "") {
		resp.Diagnostics.AddError(
			"Missing Authentication",
			"Either 'api_token' or both 'username' and 'password' must be provided "+
				"(via attributes or PROXMOX_VE_* environment variables).",
		)
		return
	}

	cfg := client.Config{
		Endpoint: endpoint,
		APIToken: apiToken,
		Username: username,
		Password: password,
		Insecure: insecure,
		Timeout:  30 * time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Proxmox VE Client",
			"An unexpected error occurred when creating the Proxmox VE API client: "+err.Error(),
		)
		return
	}

	// hand the client off to resources and datasources
	resp.DataSourceData = c
	resp.ResourceData = c

	tflog.Info(ctx, "Proxmox VE client configured", map[string]any{"endpoint": endpoint})
}

func (p *ProxmoxProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Access management
		user.NewResource,
		group.NewResource,
		role.NewResource,
		acl.NewResource,
		pool.NewResource,

		// Storage
		storage.NewResource,

		// Compute
		vm.NewResource,
		container_res.NewResource,

		// Snapshots
		vm_snapshot.NewResource,
		container_snapshot.NewResource,

		// Networking
		network_interface.NewResource,

		// Firewall
		firewall_rule.NewResource,
		firewall_options.NewResource,
		firewall_ipset.NewResource,
		firewall_alias.NewResource,
		firewall_security_group.NewResource,
		firewall_security_group_rule.NewResource,

		// Node management
		dns.NewResource,
		hosts.NewResource,
		certificate.NewResource,

		// High Availability
		ha_resource.NewResource,
		ha_group.NewResource,

		// Cluster
		cluster_options.NewResource,

		// Authentication
		realm.NewResource,

		// User tokens
		user_token.NewResource,
		
		// User password management
		user_password.NewResource,

		// Replication
		replication.NewResource,

		// Storage content
		download_file.NewResource,

		// Node management (time/timezone)
		node_time.NewResource,

		// Hardware mappings
		hardware_mapping_pci.NewResource,
		hardware_mapping_usb.NewResource,

		// APT repositories
		apt_repository.NewResource,

		// Backup schedules
		backup.NewResource,

		// SDN
		sdn_zone.NewResource,
		sdn_vnet.NewResource,
		sdn_subnet.NewResource,
		sdn_controller.NewResource,
		sdn_dns.NewResource,
		sdn_ipam.NewResource,

		// ACME
		acme_account.NewResource,
		acme_plugin.NewResource,

		// Metrics
		metrics_server.NewResource,

		// Notifications
		notification_endpoint_sendmail.NewResource,
		notification_endpoint_gotify.NewResource,
		notification_endpoint_smtp.NewResource,
		notification_endpoint_webhook.NewResource,
		notification_filter.NewResource,
		notification_matcher.NewResource,

		// Cluster jobs
		cluster_job_realm_sync.NewResource,

		// Node disks
		node_disk_directory.NewResource,
		node_disk_lvm.NewResource,
		node_disk_lvmthin.NewResource,
		node_disk_zfs.NewResource,

		// Node management (extended)
		node_config.NewResource,
		node_subscription.NewResource,
		node_service.NewResource,
		node_acme_certificate.NewResource,

		// User TFA
		user_tfa.NewResource,

		// Ceph
		ceph_pool.NewResource,
		ceph_osd.NewResource,
		ceph_mon.NewResource,
		ceph_mds.NewResource,
		ceph_mgr.NewResource,
		ceph_fs.NewResource,

		// Firewall (cluster/node level)
		cluster_firewall.NewResource,
		node_firewall.NewResource,

		// Pool membership
		pool_membership.NewResource,

		// SDN applier + type-specific zones
		sdn_applier.NewResource,
		sdn_zone_simple.NewResource,
		sdn_zone_vlan.NewResource,
		sdn_zone_vxlan.NewResource,
		sdn_zone_evpn.NewResource,
		sdn_zone_qinq.NewResource,

		// Storage type-specific
		storage_nfs.NewResource,
		storage_cifs.NewResource,
		storage_pbs.NewResource,
		storage_directory.NewResource,
		storage_lvm.NewResource,
		storage_lvmthin.NewResource,
		storage_zfspool.NewResource,

		// Network type-specific
		network_linux_bridge.NewResource,
		network_linux_vlan.NewResource,

		// APT standard repositories
		apt_standard_repo.NewResource,

		// LDAP realm
		realm_ldap.NewResource,

		// Hardware mapping dir
		hardware_mapping_dir.NewResource,

		// File upload/download
		file_res.NewResource,

		// OCI image
		oci_image.NewResource,

		// Cloned VM
		cloned_vm.NewResource,
	}
}

func (p *ProxmoxProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Cluster info
		version.NewDataSource,
		nodes.NewDataSource,
		node.NewDataSource,
		cluster_nextid_ds.NewDataSource,
		cluster_config_ds.NewDataSource,
		cluster_mapping_pci_ds.NewDataSource,
		cluster_mapping_usb_ds.NewDataSource,
		cluster_ceph_status_ds.NewDataSource,

		// Storage
		datastores.NewDataSource,
		storage_content.NewDataSource,

		// VMs
		vms.NewDataSource,
		vm_ds.NewDataSource,

		// Containers
		containers.NewDataSource,
		container_ds.NewDataSource,

		// High Availability
		ha_resources_ds.NewDataSource,
		ha_groups_ds.NewDataSource,

		// Access management
		users.NewDataSource,
		groups.NewDataSource,
		roles.NewDataSource,
		pools.NewDataSource,
		user_permissions_ds.NewDataSource,
		openid_config_ds.NewDataSource,

		// Snapshots
		vm_snapshots_ds.NewDataSource,
		container_snapshots_ds.NewDataSource,

		// Networking
		network_interfaces_ds.NewDataSource,

		// Hardware mappings
		hw_pci_ds.NewDataSource,
		hw_usb_ds.NewDataSource,

		// SDN
		sdn_zones_ds.NewDataSource,
		sdn_vnets_ds.NewDataSource,

		// APT
		apt_repositories.NewDataSource,
		apt_changelog.NewDataSource,
		apt_versions.NewDataSource,

		// Backup jobs
		backups.NewDataSource,

		// Cluster info (extended)
		cluster_status.NewDataSource,
		cluster_resources.NewDataSource,
		cluster_tasks.NewDataSource,

		// High Availability (extended)
		ha_status.NewDataSource,

		// Node hardware
		node_hardware_pci.NewDataSource,
		node_hardware_usb.NewDataSource,

		// Node disks
		node_disks.NewDataSource,
		node_smart.NewDataSource,

		// Node services
		node_services.NewDataSource,

		// Node tasks and journal
		node_tasks.NewDataSource,
		node_journal.NewDataSource,

		// ACME
		acme_directories.NewDataSource,

		// Node monitoring
		node_syslog_ds.NewDataSource,
		node_rrd_ds.NewDataSource,
		node_netstat_ds.NewDataSource,
		node_scan_ds.NewDataSource,

		// Cluster monitoring
		cluster_log_ds.NewDataSource,
		cluster_backup_info_ds.NewDataSource,

		// VM/Container RRD
		vm_rrd_ds.NewDataSource,
		container_rrd_ds.NewDataSource,

		// Single-item datasources (Access)
		ha_group_ds.NewDataSource,
		ha_resource_ds.NewDataSource,
		group_ds.NewDataSource,
		role_ds.NewDataSource,
		user_ds.NewDataSource,
		pool_ds.NewDataSource,

		// Single-item datasources (SDN)
		sdn_vnet_ds.NewDataSource,
		sdn_subnet_ds.NewDataSource,

		// Node read-only datasources
		dns_ds.NewDataSource,
		hosts_ds.NewDataSource,
		time_ds.NewDataSource,
		metrics_server_ds.NewDataSource,

		// ACME datasources
		acme_account_ds.NewDataSource,
		acme_accounts_ds.NewDataSource,
		acme_plugin_ds.NewDataSource,
		acme_plugins_ds.NewDataSource,

		// APT standard repository datasource
		apt_std_repo_ds.NewDataSource,

		// File datasource
		file_ds.NewDataSource,
	}
}

func resolveString(val types.String, envVar string) string {
	if !val.IsNull() && !val.IsUnknown() {
		return val.ValueString()
	}
	return os.Getenv(envVar)
}

func resolveBool(val types.Bool, envVar string) bool {
	if !val.IsNull() && !val.IsUnknown() {
		return val.ValueBool()
	}
	v := os.Getenv(envVar)
	return v == "true" || v == "1" || v == "yes"
}
