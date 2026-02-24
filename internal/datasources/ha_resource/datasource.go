package ha_resource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &HAResourceDataSource{}
var _ datasource.DataSourceWithConfigure = &HAResourceDataSource{}

type HAResourceDataSource struct {
	client *client.Client
}

type HAResourceDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	SID         types.String `tfsdk:"sid"`
	Type        types.String `tfsdk:"type"`
	State       types.String `tfsdk:"state"`
	Group       types.String `tfsdk:"group"`
	MaxRestart  types.Int64  `tfsdk:"max_restart"`
	MaxRelocate types.Int64  `tfsdk:"max_relocate"`
	Comment     types.String `tfsdk:"comment"`
}

func NewDataSource() datasource.DataSource {
	return &HAResourceDataSource{}
}

func (d *HAResourceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_ha_resource"
}

func (d *HAResourceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a specific Proxmox VE High Availability resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The HA resource SID.",
				Computed:    true,
			},
			"sid": schema.StringAttribute{
				Description: "The HA resource SID (e.g. 'vm:100').",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The HA resource type.",
				Computed:    true,
			},
			"state": schema.StringAttribute{
				Description: "The desired state of the HA resource.",
				Computed:    true,
			},
			"group": schema.StringAttribute{
				Description: "The HA group the resource belongs to.",
				Computed:    true,
			},
			"max_restart": schema.Int64Attribute{
				Description: "Maximum number of restart attempts.",
				Computed:    true,
			},
			"max_relocate": schema.Int64Attribute{
				Description: "Maximum number of relocation attempts.",
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Description of the HA resource.",
				Computed:    true,
			},
		},
	}
}

func (d *HAResourceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *HAResourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config HAResourceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sid := config.SID.ValueString()
	tflog.Debug(ctx, "Reading HA resource", map[string]any{"sid": sid})

	resource, err := d.client.GetHAResource(ctx, sid)
	if err != nil {
		resp.Diagnostics.AddError("Error reading HA resource", err.Error())
		return
	}

	state := HAResourceDataSourceModel{
		ID:          types.StringValue(resource.SID),
		SID:         types.StringValue(resource.SID),
		Type:        types.StringValue(resource.Type),
		State:       types.StringValue(resource.State),
		Group:       types.StringValue(resource.Group),
		MaxRestart:  types.Int64Value(int64(resource.MaxRestart)),
		MaxRelocate: types.Int64Value(int64(resource.MaxRelocate)),
		Comment:     types.StringValue(resource.Comment),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
