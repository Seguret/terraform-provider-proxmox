package acme_plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &ACMEPluginDataSource{}
var _ datasource.DataSourceWithConfigure = &ACMEPluginDataSource{}

type ACMEPluginDataSource struct {
	client *client.Client
}

type ACMEPluginDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	PluginID types.String `tfsdk:"plugin_id"`
	Type     types.String `tfsdk:"type"`
	API      types.String `tfsdk:"api"`
	Data     types.String `tfsdk:"data"`
}

func NewDataSource() datasource.DataSource {
	return &ACMEPluginDataSource{}
}

func (d *ACMEPluginDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_acme_plugin"
}

func (d *ACMEPluginDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a Proxmox VE ACME plugin.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this data source.",
				Computed:    true,
			},
			"plugin_id": schema.StringAttribute{
				Description: "The ACME plugin identifier.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The plugin type (e.g. 'dns' or 'standalone').",
				Computed:    true,
			},
			"api": schema.StringAttribute{
				Description: "The DNS provider API identifier.",
				Computed:    true,
			},
			"data": schema.StringAttribute{
				Description: "The raw plugin configuration data.",
				Computed:    true,
			},
		},
	}
}

func (d *ACMEPluginDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ACMEPluginDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ACMEPluginDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pluginID := config.PluginID.ValueString()
	tflog.Debug(ctx, "Reading ACME plugin", map[string]any{"plugin_id": pluginID})

	plugin, err := d.client.GetACMEPlugin(ctx, pluginID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ACME plugin", err.Error())
		return
	}

	state := ACMEPluginDataSourceModel{
		ID:       types.StringValue(pluginID),
		PluginID: types.StringValue(plugin.ID),
		Type:     types.StringValue(plugin.Type),
		API:      types.StringValue(plugin.API),
		Data:     types.StringValue(plugin.Data),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
