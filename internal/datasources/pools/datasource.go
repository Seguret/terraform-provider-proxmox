package pools

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &PoolsDataSource{}
var _ datasource.DataSourceWithConfigure = &PoolsDataSource{}

type PoolsDataSource struct {
	client *client.Client
}

type PoolsDataSourceModel struct {
	ID       types.String   `tfsdk:"id"`
	PoolIDs  []types.String `tfsdk:"pool_ids"`
	Comments []types.String `tfsdk:"comments"`
}

func NewDataSource() datasource.DataSource {
	return &PoolsDataSource{}
}

func (d *PoolsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_pools"
}

func (d *PoolsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of Proxmox VE resource pools.",
		Attributes: map[string]schema.Attribute{
			"id":       schema.StringAttribute{Computed: true},
			"pool_ids": schema.ListAttribute{Description: "The pool identifiers.", Computed: true, ElementType: types.StringType},
			"comments": schema.ListAttribute{Description: "Comments for each pool.", Computed: true, ElementType: types.StringType},
		},
	}
}

func (d *PoolsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PoolsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Proxmox VE pools list")

	poolsList, err := d.client.GetPools(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading pools", err.Error())
		return
	}

	state := PoolsDataSourceModel{
		ID:       types.StringValue("pools"),
		PoolIDs:  make([]types.String, len(poolsList)),
		Comments: make([]types.String, len(poolsList)),
	}

	for i, p := range poolsList {
		state.PoolIDs[i] = types.StringValue(p.PoolID)
		state.Comments[i] = types.StringValue(p.Comment)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
