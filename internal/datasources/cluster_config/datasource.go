package cluster_config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	ID             types.String `tfsdk:"id"`
	TotemInterface types.String `tfsdk:"totem_interface"`
	Nodes          types.List   `tfsdk:"nodes"`
}

type NodeModel struct {
	Name      types.String `tfsdk:"name"`
	NodeID    types.Int64  `tfsdk:"node_id"`
	IP        types.String `tfsdk:"ip"`
	Ring0Addr types.String `tfsdk:"ring0_addr"`
	Ring1Addr types.String `tfsdk:"ring1_addr"`
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_config"
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves the cluster configuration including node information.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Data source identifier.",
				Computed:            true,
			},
			"totem_interface": schema.StringAttribute{
				MarkdownDescription: "The totem interface configuration.",
				Computed:            true,
			},
			"nodes": schema.ListNestedAttribute{
				MarkdownDescription: "List of nodes in the cluster.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Node name.",
							Computed:            true,
						},
						"node_id": schema.Int64Attribute{
							MarkdownDescription: "Node ID.",
							Computed:            true,
						},
						"ip": schema.StringAttribute{
							MarkdownDescription: "Node IP address.",
							Computed:            true,
						},
						"ring0_addr": schema.StringAttribute{
							MarkdownDescription: "Ring 0 address.",
							Computed:            true,
						},
						"ring1_addr": schema.StringAttribute{
							MarkdownDescription: "Ring 1 address.",
							Computed:            true,
						},
					},
				},
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

	config, err := d.client.GetClusterConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", "Failed to retrieve cluster config: "+err.Error())
		return
	}

	data.ID = types.StringValue("cluster_config")
	data.TotemInterface = types.StringValue(config.TotemInterface)

	// map cluster nodes into the schema model
	var nodes []NodeModel
	for _, node := range config.Nodes {
		nodes = append(nodes, NodeModel{
			Name:      types.StringValue(node.Name),
			NodeID:    types.Int64Value(int64(node.NodeID)),
			IP:        types.StringValue(node.IP),
			Ring0Addr: types.StringValue(node.Ring0Addr),
			Ring1Addr: types.StringValue(node.Ring1Addr),
		})
	}

	nodesList, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":       types.StringType,
			"node_id":    types.Int64Type,
			"ip":         types.StringType,
			"ring0_addr": types.StringType,
			"ring1_addr": types.StringType,
		},
	}, nodes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Nodes = nodesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
