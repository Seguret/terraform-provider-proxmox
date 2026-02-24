package node_scan

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &NodeScanDataSource{}
var _ datasource.DataSourceWithConfigure = &NodeScanDataSource{}

type NodeScanDataSource struct {
	client *client.Client
}

type NodeScanDataSourceModel struct {
	ID       types.String   `tfsdk:"id"`
	NodeName types.String   `tfsdk:"node_name"`
	ScanType types.String   `tfsdk:"scan_type"`
	Server   types.String   `tfsdk:"server"`
	Portal   types.String   `tfsdk:"portal"`
	VG       types.String   `tfsdk:"vg"`
	Keys     []types.String `tfsdk:"keys"`
	Values   []types.String `tfsdk:"values"`
}

func NewDataSource() datasource.DataSource {
	return &NodeScanDataSource{}
}

func (d *NodeScanDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_scan"
}

func (d *NodeScanDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Scans for available storage of a given type on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{Computed: true},
			"node_name": schema.StringAttribute{Description: "The name of the node.", Required: true},
			"scan_type": schema.StringAttribute{
				Description: "The type of storage to scan (nfs, iscsi, cifs, lvm, lvmthin, zfs, pbs).",
				Required:    true,
			},
			"server": schema.StringAttribute{
				Description: "Server address (for nfs, cifs, pbs scan types).",
				Optional:    true,
			},
			"portal": schema.StringAttribute{
				Description: "iSCSI portal address (for iscsi scan type).",
				Optional:    true,
			},
			"vg": schema.StringAttribute{
				Description: "LVM volume group name (for lvmthin scan type).",
				Optional:    true,
			},
			"keys": schema.ListAttribute{
				Description: "Result field keys (flattened from each result map).",
				Computed:    true,
				ElementType: types.StringType,
			},
			"values": schema.ListAttribute{
				Description: "Result field values corresponding to each key.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *NodeScanDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NodeScanDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config NodeScanDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	scanType := config.ScanType.ValueString()

	params := map[string]string{}
	if !config.Server.IsNull() && !config.Server.IsUnknown() {
		params["server"] = config.Server.ValueString()
	}
	if !config.Portal.IsNull() && !config.Portal.IsUnknown() {
		params["portal"] = config.Portal.ValueString()
	}
	if !config.VG.IsNull() && !config.VG.IsUnknown() {
		params["vg"] = config.VG.ValueString()
	}

	tflog.Debug(ctx, "Scanning node storage", map[string]any{"node": node, "type": scanType})

	results, err := d.client.ScanNode(ctx, node, scanType, params)
	if err != nil {
		resp.Diagnostics.AddError("Error scanning node storage", err.Error())
		return
	}

	state := NodeScanDataSourceModel{
		ID:       types.StringValue(fmt.Sprintf("%s-scan-%s", node, scanType)),
		NodeName: config.NodeName,
		ScanType: config.ScanType,
		Server:   config.Server,
		Portal:   config.Portal,
		VG:       config.VG,
		Keys:     []types.String{},
		Values:   []types.String{},
	}

	for _, entry := range results {
		for k, v := range entry {
			state.Keys = append(state.Keys, types.StringValue(k))
			state.Values = append(state.Values, types.StringValue(fmt.Sprintf("%v", v)))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
