package role

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

var _ datasource.DataSource = &RoleDataSource{}
var _ datasource.DataSourceWithConfigure = &RoleDataSource{}

type RoleDataSource struct {
	client *client.Client
}

type RoleDataSourceModel struct {
	ID         types.String   `tfsdk:"id"`
	RoleID     types.String   `tfsdk:"role_id"`
	Privileges []types.String `tfsdk:"privileges"`
}

func NewDataSource() datasource.DataSource {
	return &RoleDataSource{}
}

func (d *RoleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_role"
}

func (d *RoleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a specific Proxmox VE access role.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The role identifier.",
				Computed:    true,
			},
			"role_id": schema.StringAttribute{
				Description: "The role identifier to look up.",
				Required:    true,
			},
			"privileges": schema.ListAttribute{
				Description: "The list of privileges assigned to the role.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *RoleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config RoleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleID := config.RoleID.ValueString()
	tflog.Debug(ctx, "Reading role", map[string]any{"role_id": roleID})

	role, err := d.client.GetRole(ctx, roleID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading role", err.Error())
		return
	}

	// proxmox returns privs as a comma-separated string — split it out
	var privileges []types.String
	if role.Privs != "" {
		parts := strings.Split(role.Privs, ",")
		privileges = make([]types.String, len(parts))
		for i, p := range parts {
			privileges[i] = types.StringValue(strings.TrimSpace(p))
		}
	} else {
		privileges = []types.String{}
	}

	state := RoleDataSourceModel{
		ID:         types.StringValue(role.RoleID),
		RoleID:     types.StringValue(role.RoleID),
		Privileges: privileges,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
