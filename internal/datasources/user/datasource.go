package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &UserDataSource{}
var _ datasource.DataSourceWithConfigure = &UserDataSource{}

type UserDataSource struct {
	client *client.Client
}

type UserDataSourceModel struct {
	ID        types.String   `tfsdk:"id"`
	UserID    types.String   `tfsdk:"user_id"`
	Email     types.String   `tfsdk:"email"`
	Enabled   types.Bool     `tfsdk:"enabled"`
	Comment   types.String   `tfsdk:"comment"`
	Groups    []types.String `tfsdk:"groups"`
	FirstName types.String   `tfsdk:"firstname"`
	LastName  types.String   `tfsdk:"lastname"`
}

func NewDataSource() datasource.DataSource {
	return &UserDataSource{}
}

func (d *UserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_user"
}

func (d *UserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a specific Proxmox VE user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The user identifier.",
				Computed:    true,
			},
			"user_id": schema.StringAttribute{
				Description: "The user identifier to look up (e.g. 'user@realm').",
				Required:    true,
			},
			"email": schema.StringAttribute{
				Description: "The user's email address.",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the user account is enabled.",
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "The user description.",
				Computed:    true,
			},
			"groups": schema.ListAttribute{
				Description: "The list of groups the user belongs to.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"firstname": schema.StringAttribute{
				Description: "The user's first name.",
				Computed:    true,
			},
			"lastname": schema.StringAttribute{
				Description: "The user's last name.",
				Computed:    true,
			},
		},
	}
}

func (d *UserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config UserDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userID := config.UserID.ValueString()
	tflog.Debug(ctx, "Reading user", map[string]any{"user_id": userID})

	u, err := d.client.GetUser(ctx, userID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading user", err.Error())
		return
	}

	// Enable is a 0/1 int pointer — treat nil as enabled
	enabled := types.BoolValue(true)
	if u.Enable != nil {
		enabled = types.BoolValue(*u.Enable != 0)
	}

	// Groups comes back as a comma-separated string — split it out
	var groups []types.String
	if u.Groups != "" {
		parts := strings.Split(u.Groups, ",")
		groups = make([]types.String, len(parts))
		for i, p := range parts {
			groups[i] = types.StringValue(strings.TrimSpace(p))
		}
	} else {
		groups = []types.String{}
	}

	state := UserDataSourceModel{
		ID:        types.StringValue(u.UserID),
		UserID:    types.StringValue(u.UserID),
		Email:     types.StringValue(u.Email),
		Enabled:   enabled,
		Comment:   types.StringValue(u.Comment),
		Groups:    groups,
		FirstName: types.StringValue(u.FirstName),
		LastName:  types.StringValue(u.LastName),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
