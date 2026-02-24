package cluster_ceph_status

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &DataSource{}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *client.Client
}

type DataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	HealthStatus types.String `tfsdk:"health_status"`
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_ceph_status"
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Ceph cluster status (read-only).",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Data source identifier.",
				Computed:            true,
			},
			"health_status": schema.StringAttribute{
				MarkdownDescription: "Ceph cluster health status (HEALTH_OK, HEALTH_WARN, HEALTH_ERR).",
				Computed:            true,
			},
		},
	}
}

func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *client.Client, got something else.",
		)
		return
	}

	d.client = c
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	status, err := d.client.GetCephStatus(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", "Failed to retrieve Ceph status: "+err.Error())
		return
	}

	data.ID = types.StringValue("cluster_ceph_status")
	data.HealthStatus = types.StringValue(status.Health.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
