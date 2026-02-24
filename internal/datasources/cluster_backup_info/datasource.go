package cluster_backup_info

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &ClusterBackupInfoDataSource{}
var _ datasource.DataSourceWithConfigure = &ClusterBackupInfoDataSource{}

type ClusterBackupInfoDataSource struct {
	client *client.Client
}

type ClusterBackupInfoEntryModel struct {
	VMID types.Int64  `tfsdk:"vmid"`
	Type types.String `tfsdk:"type"`
	Name types.String `tfsdk:"name"`
}

type ClusterBackupInfoDataSourceModel struct {
	ID  types.String                  `tfsdk:"id"`
	VMs []ClusterBackupInfoEntryModel `tfsdk:"vms"`
}

func NewDataSource() datasource.DataSource {
	return &ClusterBackupInfoDataSource{}
}

func (d *ClusterBackupInfoDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_cluster_backup_info"
}

func (d *ClusterBackupInfoDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of VMs not covered by any backup job in the Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"vms": schema.ListNestedAttribute{
				Description: "VMs without backup coverage.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vmid": schema.Int64Attribute{Computed: true, Description: "VM/CT ID."},
						"type": schema.StringAttribute{Computed: true, Description: "Guest type (qemu or lxc)."},
						"name": schema.StringAttribute{Computed: true, Description: "Guest name."},
					},
				},
			},
		},
	}
}

func (d *ClusterBackupInfoDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	d.client = c
}

func (d *ClusterBackupInfoDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading cluster backup info")

	entries, err := d.client.GetClusterBackupInfo(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading cluster backup info", err.Error())
		return
	}

	state := ClusterBackupInfoDataSourceModel{
		ID:  types.StringValue("cluster-backup-info"),
		VMs: make([]ClusterBackupInfoEntryModel, len(entries)),
	}

	for i, e := range entries {
		state.VMs[i] = ClusterBackupInfoEntryModel{
			VMID: types.Int64Value(int64(e.VMID)),
			Type: types.StringValue(e.Type),
			Name: types.StringValue(e.Name),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
