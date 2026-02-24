package time

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &TimeDataSource{}
var _ datasource.DataSourceWithConfigure = &TimeDataSource{}

type TimeDataSource struct {
	client *client.Client
}

type TimeDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	NodeName  types.String `tfsdk:"node_name"`
	Timezone  types.String `tfsdk:"timezone"`
	LocalTime types.Int64  `tfsdk:"local_time"`
	Time      types.Int64  `tfsdk:"time"`
}

func NewDataSource() datasource.DataSource {
	return &TimeDataSource{}
}

func (d *TimeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_time"
}

func (d *TimeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the current time and timezone of a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The time datasource identifier.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The node to retrieve time information from.",
				Required:    true,
			},
			"timezone": schema.StringAttribute{
				Description: "The timezone configured on the node (e.g. 'Europe/Berlin').",
				Computed:    true,
			},
			"local_time": schema.Int64Attribute{
				Description: "The local time on the node as a Unix timestamp.",
				Computed:    true,
			},
			"time": schema.Int64Attribute{
				Description: "The UTC time on the node as a Unix timestamp.",
				Computed:    true,
			},
		},
	}
}

func (d *TimeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TimeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TimeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading node time", map[string]any{"node": node})

	result, err := d.client.GetNodeTime(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading node time", err.Error())
		return
	}

	timezone, _ := result["timezone"].(string)

	var localTime int64
	if v, ok := result["localtime"].(float64); ok {
		localTime = int64(v)
	}

	var utcTime int64
	if v, ok := result["time"].(float64); ok {
		utcTime = int64(v)
	}

	state := TimeDataSourceModel{
		ID:        types.StringValue(fmt.Sprintf("%s/time", node)),
		NodeName:  types.StringValue(node),
		Timezone:  types.StringValue(timezone),
		LocalTime: types.Int64Value(localTime),
		Time:      types.Int64Value(utcTime),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
