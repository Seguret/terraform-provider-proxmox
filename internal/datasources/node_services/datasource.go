package node_services

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &NodeServicesDataSource{}
var _ datasource.DataSourceWithConfigure = &NodeServicesDataSource{}

// NodeServicesDataSource fetches the list of system services on a proxmox node.
type NodeServicesDataSource struct {
	client *client.Client
}

type NodeServicesDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	NodeName types.String `tfsdk:"node_name"`
	Services types.List   `tfsdk:"services"`
}

// serviceObjectType is the terraform object type for a single service entry.
var serviceObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":         types.StringType,
		"state":        types.StringType,
		"desc":         types.StringType,
		"active_state": types.StringType,
		"sub_state":    types.StringType,
	},
}

func NewDataSource() datasource.DataSource {
	return &NodeServicesDataSource{}
}

func (d *NodeServicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_services"
}

func (d *NodeServicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of system services on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the Proxmox VE node.",
				Required:    true,
			},
			"services": schema.ListNestedAttribute{
				Description: "The list of services on the node.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The service name.",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "The current run state of the service (e.g. 'running', 'stopped').",
							Computed:    true,
						},
						"desc": schema.StringAttribute{
							Description: "A short description of the service.",
							Computed:    true,
						},
						"active_state": schema.StringAttribute{
							Description: "The systemd active state (e.g. 'active', 'inactive').",
							Computed:    true,
						},
						"sub_state": schema.StringAttribute{
							Description: "The systemd sub-state (e.g. 'running', 'dead').",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *NodeServicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NodeServicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config NodeServicesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nodeName := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading Proxmox VE node services", map[string]any{"node": nodeName})

	services, err := d.client.ListNodeServices(ctx, nodeName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Proxmox VE Node Services",
			fmt.Sprintf("An error occurred while reading services for node '%s': %s", nodeName, err.Error()),
		)
		return
	}

	serviceObjects := make([]attr.Value, len(services))
	for i, svc := range services {
		obj, diags := types.ObjectValue(
			serviceObjectType.AttrTypes,
			map[string]attr.Value{
				"name":         types.StringValue(svc.Name),
				"state":        types.StringValue(svc.State),
				"desc":         types.StringValue(svc.Desc),
				"active_state": types.StringValue(svc.ActiveState),
				"sub_state":    types.StringValue(svc.SubState),
			},
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		serviceObjects[i] = obj
	}

	servicesList, diags := types.ListValue(serviceObjectType, serviceObjects)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := NodeServicesDataSourceModel{
		ID:       types.StringValue(nodeName),
		NodeName: types.StringValue(nodeName),
		Services: servicesList,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
