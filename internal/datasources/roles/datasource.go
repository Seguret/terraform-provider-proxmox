package roles

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &RolesDataSource{}
var _ datasource.DataSourceWithConfigure = &RolesDataSource{}

type RolesDataSource struct {
	client *client.Client
}

type RolesDataSourceModel struct {
	ID         types.String   `tfsdk:"id"`
	RoleIDs    []types.String `tfsdk:"role_ids"`
	Privileges []types.String `tfsdk:"privileges"`
}

func NewDataSource() datasource.DataSource {
	return &RolesDataSource{}
}

func (d *RolesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_roles"
}

func (d *RolesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of Proxmox VE roles.",
		Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{Computed: true},
			"role_ids":   schema.ListAttribute{Description: "The role identifiers.", Computed: true, ElementType: types.StringType},
			"privileges": schema.ListAttribute{Description: "Comma-separated privileges for each role.", Computed: true, ElementType: types.StringType},
		},
	}
}

func (d *RolesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RolesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Proxmox VE roles list")

	rolesList, err := d.client.GetRoles(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading roles", err.Error())
		return
	}

	state := RolesDataSourceModel{
		ID:         types.StringValue("roles"),
		RoleIDs:    make([]types.String, len(rolesList)),
		Privileges: make([]types.String, len(rolesList)),
	}

	for i, r := range rolesList {
		state.RoleIDs[i] = types.StringValue(r.RoleID)
		state.Privileges[i] = types.StringValue(r.Privs)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
