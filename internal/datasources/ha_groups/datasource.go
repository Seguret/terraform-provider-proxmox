package ha_groups

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &HAGroupsDataSource{}
var _ datasource.DataSourceWithConfigure = &HAGroupsDataSource{}

type HAGroupsDataSource struct {
	client *client.Client
}

type HAGroupsDataSourceModel struct {
	ID     types.String   `tfsdk:"id"`
	Groups []types.String `tfsdk:"groups"`
	Nodes  []types.String `tfsdk:"nodes"`
}

func NewDataSource() datasource.DataSource {
	return &HAGroupsDataSource{}
}

func (d *HAGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_ha_groups"
}

func (d *HAGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of Proxmox VE High Availability groups.",
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true},
			"groups": schema.ListAttribute{Description: "HA group names.", Computed: true, ElementType: types.StringType},
			"nodes":  schema.ListAttribute{Description: "Node members of each HA group.", Computed: true, ElementType: types.StringType},
		},
	}
}

func (d *HAGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HAGroupsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading HA groups list")

	groups, err := d.client.GetHAGroups(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading HA groups", err.Error())
		return
	}

	state := HAGroupsDataSourceModel{
		ID:     types.StringValue("ha_groups"),
		Groups: make([]types.String, len(groups)),
		Nodes:  make([]types.String, len(groups)),
	}

	for i, g := range groups {
		state.Groups[i] = types.StringValue(g.Group)
		state.Nodes[i] = types.StringValue(g.Nodes)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
