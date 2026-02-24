package sdn_subnet

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &SDNSubnetDataSource{}
var _ datasource.DataSourceWithConfigure = &SDNSubnetDataSource{}

type SDNSubnetDataSource struct {
	client *client.Client
}

type SDNSubnetDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Vnet    types.String `tfsdk:"vnet"`
	Subnet  types.String `tfsdk:"subnet"`
	Gateway types.String `tfsdk:"gateway"`
	Type    types.String `tfsdk:"type"`
}

func NewDataSource() datasource.DataSource {
	return &SDNSubnetDataSource{}
}

func (d *SDNSubnetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_sdn_subnet"
}

func (d *SDNSubnetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a specific Proxmox VE SDN subnet.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The subnet identifier.",
				Computed:    true,
			},
			"vnet": schema.StringAttribute{
				Description: "The VNet this subnet belongs to.",
				Required:    true,
			},
			"subnet": schema.StringAttribute{
				Description: "The subnet CIDR (e.g. '10.0.0.0/24').",
				Required:    true,
			},
			"gateway": schema.StringAttribute{
				Description: "The subnet gateway IP address.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The subnet type.",
				Computed:    true,
			},
		},
	}
}

func (d *SDNSubnetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SDNSubnetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config SDNSubnetDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vnet := config.Vnet.ValueString()
	subnet := config.Subnet.ValueString()
	tflog.Debug(ctx, "Reading SDN subnet", map[string]any{"vnet": vnet, "subnet": subnet})

	result, err := d.client.GetSDNSubnet(ctx, vnet, subnet)
	if err != nil {
		resp.Diagnostics.AddError("Error reading SDN subnet", err.Error())
		return
	}

	state := SDNSubnetDataSourceModel{
		ID:      types.StringValue(fmt.Sprintf("%s/%s", vnet, result.Subnet)),
		Vnet:    types.StringValue(vnet),
		Subnet:  types.StringValue(result.Subnet),
		Gateway: types.StringValue(result.Gateway),
		Type:    types.StringValue(result.Type),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
