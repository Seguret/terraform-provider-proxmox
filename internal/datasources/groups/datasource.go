package groups

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &GroupsDataSource{}
var _ datasource.DataSourceWithConfigure = &GroupsDataSource{}

type GroupsDataSource struct {
	client *client.Client
}

type GroupsDataSourceModel struct {
	ID       types.String   `tfsdk:"id"`
	GroupIDs []types.String `tfsdk:"group_ids"`
	Comments []types.String `tfsdk:"comments"`
}

func NewDataSource() datasource.DataSource {
	return &GroupsDataSource{}
}

func (d *GroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_groups"
}

func (d *GroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of Proxmox VE groups.",
		Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{Computed: true},
			"group_ids": schema.ListAttribute{Description: "The group identifiers.", Computed: true, ElementType: types.StringType},
			"comments":  schema.ListAttribute{Description: "Comments for each group.", Computed: true, ElementType: types.StringType},
		},
	}
}

func (d *GroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *GroupsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Proxmox VE groups list")

	groupsList, err := d.client.GetGroups(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading groups", err.Error())
		return
	}

	state := GroupsDataSourceModel{
		ID:       types.StringValue("groups"),
		GroupIDs: make([]types.String, len(groupsList)),
		Comments: make([]types.String, len(groupsList)),
	}

	for i, g := range groupsList {
		state.GroupIDs[i] = types.StringValue(g.GroupID)
		state.Comments[i] = types.StringValue(g.Comment)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
