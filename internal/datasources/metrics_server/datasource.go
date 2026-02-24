package metrics_server

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &MetricsServerDataSource{}
var _ datasource.DataSourceWithConfigure = &MetricsServerDataSource{}

type MetricsServerDataSource struct {
	client *client.Client
}

type MetricsServerDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	ServerID types.String `tfsdk:"server_id"`
	Type     types.String `tfsdk:"type"`
	Port     types.Int64  `tfsdk:"port"`
	Server   types.String `tfsdk:"server"`
	Disable  types.Bool   `tfsdk:"disable"`
}

func NewDataSource() datasource.DataSource {
	return &MetricsServerDataSource{}
}

func (d *MetricsServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_metrics_server"
}

func (d *MetricsServerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a specific Proxmox VE external metrics server.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The metrics server identifier.",
				Computed:    true,
			},
			"server_id": schema.StringAttribute{
				Description: "The metrics server identifier to look up.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The metrics server type (e.g. 'influxdb', 'graphite').",
				Computed:    true,
			},
			"port": schema.Int64Attribute{
				Description: "The port the metrics server listens on.",
				Computed:    true,
			},
			"server": schema.StringAttribute{
				Description: "The hostname or IP address of the metrics server.",
				Computed:    true,
			},
			"disable": schema.BoolAttribute{
				Description: "Whether the metrics server integration is disabled.",
				Computed:    true,
			},
		},
	}
}

func (d *MetricsServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MetricsServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config MetricsServerDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverID := config.ServerID.ValueString()
	tflog.Debug(ctx, "Reading metrics server", map[string]any{"server_id": serverID})

	ms, err := d.client.GetMetricsServer(ctx, serverID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading metrics server", err.Error())
		return
	}

	disable := types.BoolValue(false)
	if ms.Disable != nil {
		disable = types.BoolValue(*ms.Disable != 0)
	}

	state := MetricsServerDataSourceModel{
		ID:       types.StringValue(ms.ID),
		ServerID: types.StringValue(ms.ID),
		Type:     types.StringValue(ms.Type),
		Port:     types.Int64Value(int64(ms.Port)),
		Server:   types.StringValue(ms.Server),
		Disable:  disable,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
