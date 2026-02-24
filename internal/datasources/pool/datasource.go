package pool

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &PoolDataSource{}
var _ datasource.DataSourceWithConfigure = &PoolDataSource{}

type PoolDataSource struct {
	client *client.Client
}

type PoolMemberModel struct {
	Type types.String `tfsdk:"type"`
	VMID types.Int64  `tfsdk:"vmid"`
	ID   types.String `tfsdk:"id"`
}

type PoolDataSourceModel struct {
	ID      types.String      `tfsdk:"id"`
	PoolID  types.String      `tfsdk:"pool_id"`
	Comment types.String      `tfsdk:"comment"`
	Members []PoolMemberModel `tfsdk:"members"`
}

func NewDataSource() datasource.DataSource {
	return &PoolDataSource{}
}

func (d *PoolDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_pool"
}

func (d *PoolDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a specific Proxmox VE resource pool.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The pool identifier.",
				Computed:    true,
			},
			"pool_id": schema.StringAttribute{
				Description: "The pool identifier to look up.",
				Required:    true,
			},
			"comment": schema.StringAttribute{
				Description: "The pool description.",
				Computed:    true,
			},
			"members": schema.ListNestedAttribute{
				Description: "The list of members in the pool.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "The member type (e.g. 'qemu', 'lxc', 'storage').",
							Computed:    true,
						},
						"vmid": schema.Int64Attribute{
							Description: "The VM/container ID (if applicable).",
							Computed:    true,
						},
						"id": schema.StringAttribute{
							Description: "The member resource identifier.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *PoolDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config PoolDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolID := config.PoolID.ValueString()
	tflog.Debug(ctx, "Reading pool", map[string]any{"pool_id": poolID})

	pool, err := d.client.GetPool(ctx, poolID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading pool", err.Error())
		return
	}

	members := make([]PoolMemberModel, len(pool.Members))
	for i, m := range pool.Members {
		members[i] = PoolMemberModel{
			Type: types.StringValue(m.Type),
			VMID: types.Int64Value(int64(m.VMID)),
			ID:   types.StringValue(m.ID),
		}
	}

	state := PoolDataSourceModel{
		ID:      types.StringValue(pool.PoolID),
		PoolID:  types.StringValue(pool.PoolID),
		Comment: types.StringValue(pool.Comment),
		Members: members,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
