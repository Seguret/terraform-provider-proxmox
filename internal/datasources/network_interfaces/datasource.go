package network_interfaces

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &NetworkInterfacesDataSource{}
var _ datasource.DataSourceWithConfigure = &NetworkInterfacesDataSource{}

type NetworkInterfacesDataSource struct {
	client *client.Client
}

type NetworkInterfacesDataSourceModel struct {
	ID        types.String   `tfsdk:"id"`
	NodeName  types.String   `tfsdk:"node_name"`
	Names     []types.String `tfsdk:"names"`
	Types     []types.String `tfsdk:"types"`
	Addresses []types.String `tfsdk:"addresses"`
	CIDRs     []types.String `tfsdk:"cidrs"`
	Active    []types.Bool   `tfsdk:"active"`
}

func NewDataSource() datasource.DataSource {
	return &NetworkInterfacesDataSource{}
}

func (d *NetworkInterfacesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_network_interfaces"
}

func (d *NetworkInterfacesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of network interfaces on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{Computed: true},
			"node_name": schema.StringAttribute{Description: "The node name.", Required: true},
			"names":     schema.ListAttribute{Description: "Interface names.", Computed: true, ElementType: types.StringType},
			"types":     schema.ListAttribute{Description: "Interface types.", Computed: true, ElementType: types.StringType},
			"addresses": schema.ListAttribute{Description: "IPv4 addresses.", Computed: true, ElementType: types.StringType},
			"cidrs":     schema.ListAttribute{Description: "CIDR notation addresses.", Computed: true, ElementType: types.StringType},
			"active":    schema.ListAttribute{Description: "Whether each interface is active.", Computed: true, ElementType: types.BoolType},
		},
	}
}

func (d *NetworkInterfacesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkInterfacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config NetworkInterfacesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading network interfaces", map[string]any{"node": node})

	ifaces, err := d.client.GetNetworkInterfaces(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading network interfaces", err.Error())
		return
	}

	state := NetworkInterfacesDataSourceModel{
		ID:        types.StringValue(fmt.Sprintf("%s/network", node)),
		NodeName:  config.NodeName,
		Names:     make([]types.String, len(ifaces)),
		Types:     make([]types.String, len(ifaces)),
		Addresses: make([]types.String, len(ifaces)),
		CIDRs:     make([]types.String, len(ifaces)),
		Active:    make([]types.Bool, len(ifaces)),
	}

	for i, iface := range ifaces {
		state.Names[i] = types.StringValue(iface.Iface)
		state.Types[i] = types.StringValue(iface.Type)
		state.Addresses[i] = types.StringValue(iface.Address)
		state.CIDRs[i] = types.StringValue(iface.CIDR)
		state.Active[i] = types.BoolValue(iface.Active == 1)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
