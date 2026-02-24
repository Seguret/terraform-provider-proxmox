package cluster_status

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &ClusterStatusDataSource{}
var _ datasource.DataSourceWithConfigure = &ClusterStatusDataSource{}

type ClusterStatusDataSource struct {
	client *client.Client
}

type ClusterStatusEntryModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Type    types.String `tfsdk:"type"`
	IP      types.String `tfsdk:"ip"`
	Online  types.Int64  `tfsdk:"online"`
	Local   types.Int64  `tfsdk:"local"`
	NodeID  types.Int64  `tfsdk:"node_id"`
	Version types.Int64  `tfsdk:"version"`
	Quorate types.Int64  `tfsdk:"quorate"`
	Nodes   types.Int64  `tfsdk:"nodes"`
}

type ClusterStatusDataSourceModel struct {
	ID      types.String              `tfsdk:"id"`
	Entries []ClusterStatusEntryModel `tfsdk:"entries"`
}

func NewDataSource() datasource.DataSource {
	return &ClusterStatusDataSource{}
}

func (d *ClusterStatusDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_cluster_status"
}

func (d *ClusterStatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the current status of the Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"entries": schema.ListNestedAttribute{
				Description: "The list of cluster status entries.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The entry identifier.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the node or cluster.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The entry type (e.g., 'node', 'cluster').",
							Computed:    true,
						},
						"ip": schema.StringAttribute{
							Description: "The IP address of the node.",
							Computed:    true,
						},
						"online": schema.Int64Attribute{
							Description: "Whether the node is online (1) or offline (0).",
							Computed:    true,
						},
						"local": schema.Int64Attribute{
							Description: "Whether this is the local node (1) or not (0).",
							Computed:    true,
						},
						"node_id": schema.Int64Attribute{
							Description: "The numeric node ID.",
							Computed:    true,
						},
						"version": schema.Int64Attribute{
							Description: "The configuration version.",
							Computed:    true,
						},
						"quorate": schema.Int64Attribute{
							Description: "Whether the cluster has quorum (1) or not (0).",
							Computed:    true,
						},
						"nodes": schema.Int64Attribute{
							Description: "The total number of nodes in the cluster.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *ClusterStatusDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ClusterStatusDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Proxmox VE cluster status")

	entries, err := d.client.GetClusterStatus(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading cluster status", err.Error())
		return
	}

	state := ClusterStatusDataSourceModel{
		ID:      types.StringValue("cluster"),
		Entries: make([]ClusterStatusEntryModel, len(entries)),
	}

	for i, e := range entries {
		state.Entries[i] = ClusterStatusEntryModel{
			ID:      types.StringValue(e.ID),
			Name:    types.StringValue(e.Name),
			Type:    types.StringValue(e.Type),
			IP:      types.StringValue(e.IP),
			Online:  types.Int64Value(int64(e.Online)),
			Local:   types.Int64Value(int64(e.Local)),
			NodeID:  types.Int64Value(int64(e.NodeID)),
			Version: types.Int64Value(int64(e.Version)),
			Quorate: types.Int64Value(int64(e.Quorate)),
			Nodes:   types.Int64Value(int64(e.Nodes)),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
