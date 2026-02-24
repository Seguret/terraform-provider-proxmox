package node_syslog

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &NodeSyslogDataSource{}
var _ datasource.DataSourceWithConfigure = &NodeSyslogDataSource{}

type NodeSyslogDataSource struct {
	client *client.Client
}

type NodeSyslogDataSourceModel struct {
	ID       types.String       `tfsdk:"id"`
	NodeName types.String       `tfsdk:"node_name"`
	Start    types.Int64        `tfsdk:"start"`
	Limit    types.Int64        `tfsdk:"limit"`
	Lines    []types.Int64      `tfsdk:"lines"`
	Texts    []types.String     `tfsdk:"texts"`
}

func NewDataSource() datasource.DataSource {
	return &NodeSyslogDataSource{}
}

func (d *NodeSyslogDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_syslog"
}

func (d *NodeSyslogDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves syslog entries from a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the node.",
				Required:    true,
			},
			"start": schema.Int64Attribute{
				Description: "Start index for syslog entries (default 0).",
				Optional:    true,
				Computed:    true,
			},
			"limit": schema.Int64Attribute{
				Description: "Maximum number of syslog entries to return (default 500).",
				Optional:    true,
				Computed:    true,
			},
			"lines": schema.ListAttribute{
				Description: "Line numbers of the syslog entries.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"texts": schema.ListAttribute{
				Description: "Text content of the syslog entries.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *NodeSyslogDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NodeSyslogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config NodeSyslogDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	start := int(0)
	limit := int(500)
	if !config.Start.IsNull() && !config.Start.IsUnknown() {
		start = int(config.Start.ValueInt64())
	}
	if !config.Limit.IsNull() && !config.Limit.IsUnknown() {
		limit = int(config.Limit.ValueInt64())
	}

	tflog.Debug(ctx, "Reading node syslog", map[string]any{"node": node})

	entries, err := d.client.GetNodeSyslog(ctx, node, start, limit)
	if err != nil {
		resp.Diagnostics.AddError("Error reading node syslog", err.Error())
		return
	}

	state := NodeSyslogDataSourceModel{
		ID:       types.StringValue(fmt.Sprintf("%s-syslog", node)),
		NodeName: config.NodeName,
		Start:    types.Int64Value(int64(start)),
		Limit:    types.Int64Value(int64(limit)),
		Lines:    make([]types.Int64, len(entries)),
		Texts:    make([]types.String, len(entries)),
	}

	for i, e := range entries {
		state.Lines[i] = types.Int64Value(int64(e.N))
		state.Texts[i] = types.StringValue(e.Text)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
