package node_netstat

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &NodeNetstatDataSource{}
var _ datasource.DataSourceWithConfigure = &NodeNetstatDataSource{}

type NodeNetstatDataSource struct {
	client *client.Client
}

type NodeNetstatEntryModel struct {
	Iface   types.String `tfsdk:"iface"`
	RxPkts  types.Int64  `tfsdk:"rx_pkts"`
	TxPkts  types.Int64  `tfsdk:"tx_pkts"`
	RxBytes types.Int64  `tfsdk:"rx_bytes"`
	TxBytes types.Int64  `tfsdk:"tx_bytes"`
	RxErr   types.Int64  `tfsdk:"rx_err"`
	TxErr   types.Int64  `tfsdk:"tx_err"`
	RxDrop  types.Int64  `tfsdk:"rx_drop"`
	TxDrop  types.Int64  `tfsdk:"tx_drop"`
}

type NodeNetstatDataSourceModel struct {
	ID       types.String            `tfsdk:"id"`
	NodeName types.String            `tfsdk:"node_name"`
	Entries  []NodeNetstatEntryModel `tfsdk:"entries"`
}

func NewDataSource() datasource.DataSource {
	return &NodeNetstatDataSource{}
}

func (d *NodeNetstatDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_netstat"
}

func (d *NodeNetstatDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves network interface statistics for a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{Computed: true},
			"node_name": schema.StringAttribute{Description: "The name of the node.", Required: true},
			"entries": schema.ListNestedAttribute{
				Description: "Network interface statistics entries.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"iface":    schema.StringAttribute{Computed: true, Description: "Interface name."},
						"rx_pkts":  schema.Int64Attribute{Computed: true, Description: "Received packets."},
						"tx_pkts":  schema.Int64Attribute{Computed: true, Description: "Transmitted packets."},
						"rx_bytes": schema.Int64Attribute{Computed: true, Description: "Received bytes."},
						"tx_bytes": schema.Int64Attribute{Computed: true, Description: "Transmitted bytes."},
						"rx_err":   schema.Int64Attribute{Computed: true, Description: "Receive errors."},
						"tx_err":   schema.Int64Attribute{Computed: true, Description: "Transmit errors."},
						"rx_drop":  schema.Int64Attribute{Computed: true, Description: "Received packets dropped."},
						"tx_drop":  schema.Int64Attribute{Computed: true, Description: "Transmitted packets dropped."},
					},
				},
			},
		},
	}
}

func (d *NodeNetstatDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NodeNetstatDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config NodeNetstatDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading node netstat", map[string]any{"node": node})

	entries, err := d.client.GetNodeNetstat(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading node netstat", err.Error())
		return
	}

	state := NodeNetstatDataSourceModel{
		ID:       types.StringValue(fmt.Sprintf("%s-netstat", node)),
		NodeName: config.NodeName,
		Entries:  make([]NodeNetstatEntryModel, len(entries)),
	}

	for i, e := range entries {
		state.Entries[i] = NodeNetstatEntryModel{
			Iface:   types.StringValue(e.Iface),
			RxPkts:  types.Int64Value(e.RxPkts),
			TxPkts:  types.Int64Value(e.TxPkts),
			RxBytes: types.Int64Value(e.RxBytes),
			TxBytes: types.Int64Value(e.TxBytes),
			RxErr:   types.Int64Value(e.RxErr),
			TxErr:   types.Int64Value(e.TxErr),
			RxDrop:  types.Int64Value(e.RxDrop),
			TxDrop:  types.Int64Value(e.TxDrop),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
