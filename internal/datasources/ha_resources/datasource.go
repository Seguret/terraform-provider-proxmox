package ha_resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &HAResourcesDataSource{}
var _ datasource.DataSourceWithConfigure = &HAResourcesDataSource{}

type HAResourcesDataSource struct {
	client *client.Client
}

type HAResourcesDataSourceModel struct {
	ID      types.String   `tfsdk:"id"`
	SIDs    []types.String `tfsdk:"sids"`
	States  []types.String `tfsdk:"states"`
	Groups  []types.String `tfsdk:"groups"`
}

func NewDataSource() datasource.DataSource {
	return &HAResourcesDataSource{}
}

func (d *HAResourcesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_ha_resources"
}

func (d *HAResourcesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of Proxmox VE High Availability resources.",
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true},
			"sids":   schema.ListAttribute{Description: "HA resource SIDs.", Computed: true, ElementType: types.StringType},
			"states": schema.ListAttribute{Description: "Desired states of the HA resources.", Computed: true, ElementType: types.StringType},
			"groups": schema.ListAttribute{Description: "HA group names.", Computed: true, ElementType: types.StringType},
		},
	}
}

func (d *HAResourcesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HAResourcesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading HA resources list")

	resources, err := d.client.GetHAResources(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading HA resources", err.Error())
		return
	}

	state := HAResourcesDataSourceModel{
		ID:     types.StringValue("ha_resources"),
		SIDs:   make([]types.String, len(resources)),
		States: make([]types.String, len(resources)),
		Groups: make([]types.String, len(resources)),
	}

	for i, r := range resources {
		state.SIDs[i] = types.StringValue(r.SID)
		state.States[i] = types.StringValue(r.State)
		state.Groups[i] = types.StringValue(r.Group)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
