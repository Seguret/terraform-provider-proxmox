package group

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &GroupDataSource{}
var _ datasource.DataSourceWithConfigure = &GroupDataSource{}

type GroupDataSource struct {
	client *client.Client
}

type GroupDataSourceModel struct {
	ID      types.String   `tfsdk:"id"`
	GroupID types.String   `tfsdk:"group_id"`
	Comment types.String   `tfsdk:"comment"`
	Members []types.String `tfsdk:"members"`
}

func NewDataSource() datasource.DataSource {
	return &GroupDataSource{}
}

func (d *GroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_group"
}

func (d *GroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a specific Proxmox VE access group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The group identifier.",
				Computed:    true,
			},
			"group_id": schema.StringAttribute{
				Description: "The group identifier to look up.",
				Required:    true,
			},
			"comment": schema.StringAttribute{
				Description: "The group description.",
				Computed:    true,
			},
			"members": schema.ListAttribute{
				Description: "The list of users that are members of this group.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *GroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *GroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config GroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupID := config.GroupID.ValueString()
	tflog.Debug(ctx, "Reading group", map[string]any{"group_id": groupID})

	group, err := d.client.GetGroup(ctx, groupID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading group", err.Error())
		return
	}

	// members can be in either field depending on API version — check both
	membersList := group.Members
	if len(membersList) == 0 {
		membersList = group.Users
	}

	members := make([]types.String, len(membersList))
	for i, m := range membersList {
		members[i] = types.StringValue(m)
	}

	state := GroupDataSourceModel{
		ID:      types.StringValue(group.GroupID),
		GroupID: types.StringValue(group.GroupID),
		Comment: types.StringValue(group.Comment),
		Members: members,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
