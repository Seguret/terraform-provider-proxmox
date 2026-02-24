package user_permissions

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &UserPermissionsDataSource{}

type UserPermissionsDataSource struct {
	client *client.Client
}

type UserPermissionsDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	UserID      types.String `tfsdk:"user_id"`
	Path        types.String `tfsdk:"path"`
	Permissions types.Map    `tfsdk:"permissions"`
}

func NewDataSource() datasource.DataSource {
	return &UserPermissionsDataSource{}
}

func (d *UserPermissionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_permissions"
}

func (d *UserPermissionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves permissions for a specific Proxmox VE user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"user_id": schema.StringAttribute{
				Description: "The user ID to query permissions for (e.g., 'root@pam').",
				Required:    true,
			},
			"path": schema.StringAttribute{
				Description: "Optional path to filter permissions (e.g., '/vms/100'). If empty, returns all permissions.",
				Optional:    true,
			},
			"permissions": schema.MapAttribute{
				Description: "Map of permissions by path. Each path contains a map of privilege to boolean.",
				Computed:    true,
				ElementType: types.MapType{ElemType: types.Int64Type},
			},
		},
	}
}

func (d *UserPermissionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cl, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = cl
}

func (d *UserPermissionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config UserPermissionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userID := config.UserID.ValueString()
	path := config.Path.ValueString()

	permissions, err := d.client.GetUserPermissions(ctx, userID, path)
	if err != nil {
		resp.Diagnostics.AddError("Error reading user permissions", err.Error())
		return
	}

	// build the nested permissions map from the API response
	permMap := make(map[string]types.Map)
	for permPath, privs := range permissions.Permissions {
		privMap := make(map[string]types.Int64)
		for priv, val := range privs {
			privMap[priv] = types.Int64Value(int64(val))
		}
		permMapVal, diags := types.MapValueFrom(ctx, types.Int64Type, privMap)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		permMap[permPath] = permMapVal
	}

	permissionsMap, diags := types.MapValueFrom(ctx, types.MapType{ElemType: types.Int64Type}, permMap)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.ID = types.StringValue(fmt.Sprintf("%s/%s", userID, path))
	config.Permissions = permissionsMap

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
