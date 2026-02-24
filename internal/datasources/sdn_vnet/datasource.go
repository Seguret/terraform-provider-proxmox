package sdn_vnet

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &SDNVnetDataSource{}
var _ datasource.DataSourceWithConfigure = &SDNVnetDataSource{}

type SDNVnetDataSource struct {
	client *client.Client
}

type SDNVnetDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Vnet  types.String `tfsdk:"vnet"`
	Zone  types.String `tfsdk:"zone"`
	Tag   types.Int64  `tfsdk:"tag"`
	Alias types.String `tfsdk:"alias"`
}

func NewDataSource() datasource.DataSource {
	return &SDNVnetDataSource{}
}

func (d *SDNVnetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_sdn_vnet"
}

func (d *SDNVnetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a specific Proxmox VE SDN VNet.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The VNet identifier.",
				Computed:    true,
			},
			"vnet": schema.StringAttribute{
				Description: "The VNet name to look up.",
				Required:    true,
			},
			"zone": schema.StringAttribute{
				Description: "The SDN zone the VNet belongs to.",
				Computed:    true,
			},
			"tag": schema.Int64Attribute{
				Description: "The VLAN tag assigned to the VNet.",
				Computed:    true,
			},
			"alias": schema.StringAttribute{
				Description: "An alias for the VNet.",
				Computed:    true,
			},
		},
	}
}

func (d *SDNVnetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SDNVnetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config SDNVnetDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vnetID := config.Vnet.ValueString()
	tflog.Debug(ctx, "Reading SDN VNet", map[string]any{"vnet": vnetID})

	vnet, err := d.client.GetSDNVnet(ctx, vnetID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading SDN VNet", err.Error())
		return
	}

	state := SDNVnetDataSourceModel{
		ID:    types.StringValue(vnet.Vnet),
		Vnet:  types.StringValue(vnet.Vnet),
		Zone:  types.StringValue(vnet.Zone),
		Tag:   types.Int64Value(int64(vnet.Tag)),
		Alias: types.StringValue(vnet.Alias),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
