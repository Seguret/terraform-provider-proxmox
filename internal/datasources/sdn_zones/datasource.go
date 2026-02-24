package sdn_zones

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &SDNZonesDataSource{}
var _ datasource.DataSourceWithConfigure = &SDNZonesDataSource{}

type SDNZonesDataSource struct {
	client *client.Client
}

type SDNZonesDataSourceModel struct {
	ID    types.String   `tfsdk:"id"`
	Zones []types.String `tfsdk:"zones"`
	Types []types.String `tfsdk:"types"`
}

func NewDataSource() datasource.DataSource {
	return &SDNZonesDataSource{}
}

func (d *SDNZonesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_sdn_zones"
}

func (d *SDNZonesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of Proxmox VE SDN zones.",
		Attributes: map[string]schema.Attribute{
			"id":    schema.StringAttribute{Computed: true},
			"zones": schema.ListAttribute{Description: "SDN zone names.", Computed: true, ElementType: types.StringType},
			"types": schema.ListAttribute{Description: "SDN zone types.", Computed: true, ElementType: types.StringType},
		},
	}
}

func (d *SDNZonesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SDNZonesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading SDN zones list")

	zones, err := d.client.GetSDNZones(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading SDN zones", err.Error())
		return
	}

	state := SDNZonesDataSourceModel{
		ID:    types.StringValue("sdn_zones"),
		Zones: make([]types.String, len(zones)),
		Types: make([]types.String, len(zones)),
	}
	for i, z := range zones {
		state.Zones[i] = types.StringValue(z.Zone)
		state.Types[i] = types.StringValue(z.Type)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
