package node_disks

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &NodeDisksDataSource{}
var _ datasource.DataSourceWithConfigure = &NodeDisksDataSource{}

type NodeDisksDataSource struct {
	client *client.Client
}

type NodeDiskModel struct {
	Dev    types.String `tfsdk:"dev"`
	Type   types.String `tfsdk:"type"`
	Size   types.Int64  `tfsdk:"size"`
	Model  types.String `tfsdk:"model"`
	Serial types.String `tfsdk:"serial"`
	Vendor types.String `tfsdk:"vendor"`
	WWN    types.String `tfsdk:"wwn"`
	Health types.String `tfsdk:"health"`
	Used   types.String `tfsdk:"used"`
	GPT    types.Bool   `tfsdk:"gpt"`
}

type NodeDisksDataSourceModel struct {
	ID       types.String    `tfsdk:"id"`
	NodeName types.String    `tfsdk:"node_name"`
	Disks    []NodeDiskModel `tfsdk:"disks"`
}

func NewDataSource() datasource.DataSource {
	return &NodeDisksDataSource{}
}

func (d *NodeDisksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_disks"
}

func (d *NodeDisksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of disks on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The node name.",
				Required:    true,
			},
			"disks": schema.ListNestedAttribute{
				Description: "The list of disks on the node.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"dev": schema.StringAttribute{
							Description: "The device name (e.g. /dev/sda).",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The disk type (e.g. hdd, ssd).",
							Computed:    true,
						},
						"size": schema.Int64Attribute{
							Description: "The disk size in bytes.",
							Computed:    true,
						},
						"model": schema.StringAttribute{
							Description: "The disk model name.",
							Computed:    true,
						},
						"serial": schema.StringAttribute{
							Description: "The disk serial number.",
							Computed:    true,
						},
						"vendor": schema.StringAttribute{
							Description: "The disk vendor.",
							Computed:    true,
						},
						"wwn": schema.StringAttribute{
							Description: "The disk World Wide Name.",
							Computed:    true,
						},
						"health": schema.StringAttribute{
							Description: "The SMART health status.",
							Computed:    true,
						},
						"used": schema.StringAttribute{
							Description: "How the disk is currently used (e.g. LVM, ZFS).",
							Computed:    true,
						},
						"gpt": schema.BoolAttribute{
							Description: "Whether the disk has a GPT partition table.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *NodeDisksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NodeDisksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config NodeDisksDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading node disks", map[string]any{"node": node})

	disks, err := d.client.ListNodeDisks(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading node disks", err.Error())
		return
	}

	state := NodeDisksDataSourceModel{
		ID:       types.StringValue(fmt.Sprintf("%s/disks", node)),
		NodeName: config.NodeName,
		Disks:    make([]NodeDiskModel, len(disks)),
	}

	for i, disk := range disks {
		gpt := types.BoolNull()
		if disk.GPT != nil {
			gpt = types.BoolValue(*disk.GPT != 0)
		}
		
		state.Disks[i] = NodeDiskModel{
			Dev:    types.StringValue(disk.Dev),
			Type:   types.StringValue(disk.Type),
			Size:   types.Int64Value(disk.Size),
			Model:  types.StringValue(disk.Model),
			Serial: types.StringValue(disk.Serial),
			Vendor: types.StringValue(disk.Vendor),
			WWN:    types.StringValue(disk.WWN),
			Health: types.StringValue(disk.Health),
			Used:   types.StringValue(disk.Used),
			GPT:    gpt,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
