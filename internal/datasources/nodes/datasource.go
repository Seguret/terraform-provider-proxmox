package nodes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &NodesDataSource{}
var _ datasource.DataSourceWithConfigure = &NodesDataSource{}

type NodesDataSource struct {
	client *client.Client
}

type NodesDataSourceModel struct {
	ID              types.String  `tfsdk:"id"`
	Names           []types.String `tfsdk:"names"`
	Online          []types.Bool   `tfsdk:"online"`
	CPUCount        []types.Int64  `tfsdk:"cpu_count"`
	CPUUtilization  []types.Float64 `tfsdk:"cpu_utilization"`
	MemoryUsed      []types.Int64  `tfsdk:"memory_used"`
	MemoryAvailable []types.Int64  `tfsdk:"memory_available"`
	Uptime          []types.Int64  `tfsdk:"uptime"`
	SSLFingerprints []types.String `tfsdk:"ssl_fingerprints"`
	SupportLevels   []types.String `tfsdk:"support_levels"`
}

func NewDataSource() datasource.DataSource {
	return &NodesDataSource{}
}

func (d *NodesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_nodes"
}

func (d *NodesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of nodes in the Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"names": schema.ListAttribute{
				Description: "The node names.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"online": schema.ListAttribute{
				Description: "Whether each node is online.",
				Computed:    true,
				ElementType: types.BoolType,
			},
			"cpu_count": schema.ListAttribute{
				Description: "The number of CPUs for each node.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"cpu_utilization": schema.ListAttribute{
				Description: "The CPU utilization (0.0-1.0) for each node.",
				Computed:    true,
				ElementType: types.Float64Type,
			},
			"memory_used": schema.ListAttribute{
				Description: "The used memory in bytes for each node.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"memory_available": schema.ListAttribute{
				Description: "The total available memory in bytes for each node.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"uptime": schema.ListAttribute{
				Description: "The uptime in seconds for each node.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"ssl_fingerprints": schema.ListAttribute{
				Description: "The SSL fingerprints for each node.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"support_levels": schema.ListAttribute{
				Description: "The support level for each node.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *NodesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NodesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Proxmox VE nodes list")

	nodesList, err := d.client.GetNodes(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Proxmox VE Nodes",
			"An error occurred while reading the Proxmox VE nodes: "+err.Error(),
		)
		return
	}

	state := NodesDataSourceModel{
		ID:              types.StringValue("nodes"),
		Names:           make([]types.String, len(nodesList)),
		Online:          make([]types.Bool, len(nodesList)),
		CPUCount:        make([]types.Int64, len(nodesList)),
		CPUUtilization:  make([]types.Float64, len(nodesList)),
		MemoryUsed:      make([]types.Int64, len(nodesList)),
		MemoryAvailable: make([]types.Int64, len(nodesList)),
		Uptime:          make([]types.Int64, len(nodesList)),
		SSLFingerprints: make([]types.String, len(nodesList)),
		SupportLevels:   make([]types.String, len(nodesList)),
	}

	for i, node := range nodesList {
		state.Names[i] = types.StringValue(node.Node)
		state.Online[i] = types.BoolValue(node.Status == "online")
		state.CPUCount[i] = types.Int64Value(node.MaxCPU)
		state.CPUUtilization[i] = types.Float64Value(node.CPU)
		state.MemoryUsed[i] = types.Int64Value(node.Mem)
		state.MemoryAvailable[i] = types.Int64Value(node.MaxMem)
		state.Uptime[i] = types.Int64Value(node.Uptime)
		state.SSLFingerprints[i] = types.StringValue(node.SSLFingerprint)
		state.SupportLevels[i] = types.StringValue(node.Level)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
