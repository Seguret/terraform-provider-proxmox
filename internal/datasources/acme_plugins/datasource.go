package acme_plugins

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &ACMEPluginsDataSource{}
var _ datasource.DataSourceWithConfigure = &ACMEPluginsDataSource{}

type ACMEPluginsDataSource struct {
	client *client.Client
}

type ACMEPluginEntryModel struct {
	PluginID types.String `tfsdk:"plugin_id"`
	Type     types.String `tfsdk:"type"`
}

type ACMEPluginsDataSourceModel struct {
	ID      types.String           `tfsdk:"id"`
	Plugins []ACMEPluginEntryModel `tfsdk:"plugins"`
}

func NewDataSource() datasource.DataSource {
	return &ACMEPluginsDataSource{}
}

func (d *ACMEPluginsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_acme_plugins"
}

func (d *ACMEPluginsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of Proxmox VE ACME plugins.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this data source.",
				Computed:    true,
			},
			"plugins": schema.ListNestedAttribute{
				Description: "The list of ACME plugins.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"plugin_id": schema.StringAttribute{
							Description: "The ACME plugin identifier.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The plugin type (e.g. 'dns' or 'standalone').",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *ACMEPluginsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ACMEPluginsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading ACME plugins list")

	plugins, err := d.client.GetACMEPlugins(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ACME plugins", err.Error())
		return
	}

	state := ACMEPluginsDataSourceModel{
		ID:      types.StringValue("acme_plugins"),
		Plugins: make([]ACMEPluginEntryModel, len(plugins)),
	}

	for i, p := range plugins {
		state.Plugins[i] = ACMEPluginEntryModel{
			PluginID: types.StringValue(p.ID),
			Type:     types.StringValue(p.Type),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
