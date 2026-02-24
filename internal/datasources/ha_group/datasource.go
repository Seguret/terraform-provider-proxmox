package ha_group

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &HAGroupDataSource{}
var _ datasource.DataSourceWithConfigure = &HAGroupDataSource{}

type HAGroupDataSource struct {
	client *client.Client
}

type HAGroupDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Group      types.String `tfsdk:"group"`
	Comment    types.String `tfsdk:"comment"`
	Nodes      types.String `tfsdk:"nodes"`
	NoFailback types.Bool   `tfsdk:"nofailback"`
	Restricted types.Bool   `tfsdk:"restricted"`
}

func NewDataSource() datasource.DataSource {
	return &HAGroupDataSource{}
}

func (d *HAGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_ha_group"
}

func (d *HAGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a specific Proxmox VE High Availability group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The HA group identifier.",
				Computed:    true,
			},
			"group": schema.StringAttribute{
				Description: "The HA group name.",
				Required:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Description of the HA group.",
				Computed:    true,
			},
			"nodes": schema.StringAttribute{
				Description: "Comma-separated list of node:priority pairs.",
				Computed:    true,
			},
			"nofailback": schema.BoolAttribute{
				Description: "Whether to prevent failback.",
				Computed:    true,
			},
			"restricted": schema.BoolAttribute{
				Description: "Whether HA resources bound to this group may only run on the defined nodes.",
				Computed:    true,
			},
		},
	}
}

func (d *HAGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HAGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config HAGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID := config.Group.ValueString()
	tflog.Debug(ctx, "Reading HA group", map[string]any{"group": groupID})

	group, err := d.client.GetHAGroup(ctx, groupID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading HA group", err.Error())
		return
	}

	noFailback := types.BoolValue(false)
	if group.NoFailback != nil {
		noFailback = types.BoolValue(*group.NoFailback != 0)
	}

	restricted := types.BoolValue(false)
	if group.Restricted != nil {
		restricted = types.BoolValue(*group.Restricted != 0)
	}

	state := HAGroupDataSourceModel{
		ID:         types.StringValue(group.Group),
		Group:      types.StringValue(group.Group),
		Comment:    types.StringValue(group.Comment),
		Nodes:      types.StringValue(group.Nodes),
		NoFailback: noFailback,
		Restricted: restricted,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
