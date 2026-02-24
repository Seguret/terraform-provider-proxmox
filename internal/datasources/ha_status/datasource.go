package ha_status

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &HAStatusDataSource{}
var _ datasource.DataSourceWithConfigure = &HAStatusDataSource{}

type HAStatusDataSource struct {
	client *client.Client
}

type HAStatusEntryModel struct {
	SID         types.String `tfsdk:"sid"`
	State       types.String `tfsdk:"state"`
	Node        types.String `tfsdk:"node"`
	MaxRestart  types.Int64  `tfsdk:"max_restart"`
	MaxRelocate types.Int64  `tfsdk:"max_relocate"`
	CRMState    types.String `tfsdk:"crm_state"`
	Request     types.String `tfsdk:"request"`
}

type HAStatusDataSourceModel struct {
	ID      types.String         `tfsdk:"id"`
	Entries []HAStatusEntryModel `tfsdk:"entries"`
}

func NewDataSource() datasource.DataSource {
	return &HAStatusDataSource{}
}

func (d *HAStatusDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_ha_status"
}

func (d *HAStatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the current HA resource status from the Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"entries": schema.ListNestedAttribute{
				Description: "The list of HA resource status entries.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sid": schema.StringAttribute{
							Description: "The HA resource SID (e.g., 'vm:100').",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "The current HA state of the resource.",
							Computed:    true,
						},
						"node": schema.StringAttribute{
							Description: "The node the resource is currently running on.",
							Computed:    true,
						},
						"max_restart": schema.Int64Attribute{
							Description: "The maximum number of restarts allowed.",
							Computed:    true,
						},
						"max_relocate": schema.Int64Attribute{
							Description: "The maximum number of relocations allowed.",
							Computed:    true,
						},
						"crm_state": schema.StringAttribute{
							Description: "The CRM (Cluster Resource Manager) state.",
							Computed:    true,
						},
						"request": schema.StringAttribute{
							Description: "The current CRM request for this resource.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *HAStatusDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}
	d.client = c
}

func (d *HAStatusDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Proxmox VE HA status")

	entries, err := d.client.GetHAStatus(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading HA status", err.Error())
		return
	}

	state := HAStatusDataSourceModel{
		ID:      types.StringValue("ha_status"),
		Entries: make([]HAStatusEntryModel, len(entries)),
	}

	for i, e := range entries {
		state.Entries[i] = HAStatusEntryModel{
			SID:         types.StringValue(e.SID),
			State:       types.StringValue(e.State),
			Node:        types.StringValue(e.Node),
			MaxRestart:  types.Int64Value(int64(e.MaxRestart)),
			MaxRelocate: types.Int64Value(int64(e.MaxRelocate)),
			CRMState:    types.StringValue(e.CRMState),
			Request:     types.StringValue(e.Request),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
