package sdn_vnets

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &SDNVnetsDataSource{}
var _ datasource.DataSourceWithConfigure = &SDNVnetsDataSource{}

type SDNVnetsDataSource struct {
	client *client.Client
}

type SDNVnetsDataSourceModel struct {
	ID    types.String   `tfsdk:"id"`
	Vnets []types.String `tfsdk:"vnets"`
	Zones []types.String `tfsdk:"zones"`
}

func NewDataSource() datasource.DataSource {
	return &SDNVnetsDataSource{}
}

func (d *SDNVnetsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_sdn_vnets"
}

func (d *SDNVnetsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of Proxmox VE SDN VNets.",
		Attributes: map[string]schema.Attribute{
			"id":    schema.StringAttribute{Computed: true},
			"vnets": schema.ListAttribute{Description: "VNet names.", Computed: true, ElementType: types.StringType},
			"zones": schema.ListAttribute{Description: "Zone each VNet belongs to.", Computed: true, ElementType: types.StringType},
		},
	}
}

func (d *SDNVnetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SDNVnetsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading SDN vnets list")

	vnets, err := d.client.GetSDNVnets(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading SDN vnets", err.Error())
		return
	}

	state := SDNVnetsDataSourceModel{
		ID:    types.StringValue("sdn_vnets"),
		Vnets: make([]types.String, len(vnets)),
		Zones: make([]types.String, len(vnets)),
	}
	for i, v := range vnets {
		state.Vnets[i] = types.StringValue(v.Vnet)
		state.Zones[i] = types.StringValue(v.Zone)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
