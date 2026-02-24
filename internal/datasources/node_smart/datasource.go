package node_smart

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &NodeSmartDataSource{}
var _ datasource.DataSourceWithConfigure = &NodeSmartDataSource{}

type NodeSmartDataSource struct {
	client *client.Client
}

type NodeSmartDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	NodeName types.String `tfsdk:"node_name"`
	Disk     types.String `tfsdk:"disk"`
	Type     types.String `tfsdk:"type"`
	Health   types.String `tfsdk:"health"`
	Wearout  types.Int64  `tfsdk:"wearout"`
	Text     types.String `tfsdk:"text"`
}

func NewDataSource() datasource.DataSource {
	return &NodeSmartDataSource{}
}

func (d *NodeSmartDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_smart"
}

func (d *NodeSmartDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves SMART data for a disk on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The node name.",
				Required:    true,
			},
			"disk": schema.StringAttribute{
				Description: "The device path of the disk (e.g. /dev/sda).",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The SMART type (e.g. SATA, NVMe).",
				Computed:    true,
			},
			"health": schema.StringAttribute{
				Description: "The SMART health status.",
				Computed:    true,
			},
			"wearout": schema.Int64Attribute{
				Description: "The wearout indicator for SSDs (percentage remaining).",
				Computed:    true,
			},
			"text": schema.StringAttribute{
				Description: "The raw SMART text output.",
				Computed:    true,
			},
		},
	}
}

func (d *NodeSmartDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NodeSmartDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config NodeSmartDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	disk := config.Disk.ValueString()
	tflog.Debug(ctx, "Reading disk SMART data", map[string]any{"node": node, "disk": disk})

	smart, err := d.client.GetNodeDiskSmart(ctx, node, disk)
	if err != nil {
		resp.Diagnostics.AddError("Error reading disk SMART data", err.Error())
		return
	}

	state := NodeSmartDataSourceModel{
		ID:       types.StringValue(fmt.Sprintf("%s/disks/smart/%s", node, disk)),
		NodeName: config.NodeName,
		Disk:     config.Disk,
		Type:     types.StringValue(smart.Type),
		Health:   types.StringValue(smart.Health),
		Wearout:  types.Int64Value(int64(smart.Wearout)),
		Text:     types.StringValue(smart.Text),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
