package users

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &UsersDataSource{}
var _ datasource.DataSourceWithConfigure = &UsersDataSource{}

type UsersDataSource struct {
	client *client.Client
}

type UsersDataSourceModel struct {
	ID       types.String   `tfsdk:"id"`
	UserIDs  []types.String `tfsdk:"user_ids"`
	Emails   []types.String `tfsdk:"emails"`
	Enabled  []types.Bool   `tfsdk:"enabled"`
	Comments []types.String `tfsdk:"comments"`
}

func NewDataSource() datasource.DataSource {
	return &UsersDataSource{}
}

func (d *UsersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_users"
}

func (d *UsersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of Proxmox VE users.",
		Attributes: map[string]schema.Attribute{
			"id":       schema.StringAttribute{Computed: true},
			"user_ids": schema.ListAttribute{Description: "The user identifiers.", Computed: true, ElementType: types.StringType},
			"emails":   schema.ListAttribute{Description: "The email addresses of each user.", Computed: true, ElementType: types.StringType},
			"enabled":  schema.ListAttribute{Description: "Whether each user account is enabled.", Computed: true, ElementType: types.BoolType},
			"comments": schema.ListAttribute{Description: "Comments for each user.", Computed: true, ElementType: types.StringType},
		},
	}
}

func (d *UsersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UsersDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Proxmox VE users list")

	usersList, err := d.client.GetUsers(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading users", err.Error())
		return
	}

	state := UsersDataSourceModel{
		ID:       types.StringValue("users"),
		UserIDs:  make([]types.String, len(usersList)),
		Emails:   make([]types.String, len(usersList)),
		Enabled:  make([]types.Bool, len(usersList)),
		Comments: make([]types.String, len(usersList)),
	}

	for i, u := range usersList {
		state.UserIDs[i] = types.StringValue(u.UserID)
		state.Emails[i] = types.StringValue(u.Email)
		if u.Enable != nil {
			state.Enabled[i] = types.BoolValue(*u.Enable == 1)
		} else {
			state.Enabled[i] = types.BoolValue(true)
		}
		state.Comments[i] = types.StringValue(u.Comment)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
