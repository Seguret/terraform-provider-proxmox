package containers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &ContainersDataSource{}
var _ datasource.DataSourceWithConfigure = &ContainersDataSource{}

type ContainersDataSource struct {
	client *client.Client
}

type ContainersDataSourceModel struct {
	ID       types.String   `tfsdk:"id"`
	NodeName types.String   `tfsdk:"node_name"`
	VMIDs    []types.Int64  `tfsdk:"vmids"`
	Names    []types.String `tfsdk:"names"`
	Statuses []types.String `tfsdk:"statuses"`
	Tags     []types.String `tfsdk:"tags"`
	CPUs     []types.Int64  `tfsdk:"cpus"`
	MaxMem   []types.Int64  `tfsdk:"max_memory"`
	MaxDisk  []types.Int64  `tfsdk:"max_disk"`
	Uptime   []types.Int64  `tfsdk:"uptime"`
	Template []types.Bool   `tfsdk:"template"`
}

func NewDataSource() datasource.DataSource {
	return &ContainersDataSource{}
}

func (d *ContainersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_containers"
}

func (d *ContainersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of LXC containers on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{Computed: true},
			"node_name": schema.StringAttribute{Description: "The node name.", Required: true},
			"vmids":     schema.ListAttribute{Description: "The container IDs.", Computed: true, ElementType: types.Int64Type},
			"names":     schema.ListAttribute{Description: "The container names.", Computed: true, ElementType: types.StringType},
			"statuses":  schema.ListAttribute{Description: "The container statuses.", Computed: true, ElementType: types.StringType},
			"tags":      schema.ListAttribute{Description: "The container tags.", Computed: true, ElementType: types.StringType},
			"cpus":      schema.ListAttribute{Description: "The number of CPUs.", Computed: true, ElementType: types.Int64Type},
			"max_memory": schema.ListAttribute{Description: "The maximum memory in bytes.", Computed: true, ElementType: types.Int64Type},
			"max_disk":  schema.ListAttribute{Description: "The maximum disk size in bytes.", Computed: true, ElementType: types.Int64Type},
			"uptime":    schema.ListAttribute{Description: "The uptime in seconds.", Computed: true, ElementType: types.Int64Type},
			"template":  schema.ListAttribute{Description: "Whether each container is a template.", Computed: true, ElementType: types.BoolType},
		},
	}
}

func (d *ContainersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	d.client = c
}

func (d *ContainersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ContainersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nodeName := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading containers list", map[string]any{"node": nodeName})

	ctList, err := d.client.GetContainers(ctx, nodeName)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Containers",
			fmt.Sprintf("Error reading containers on node '%s': %s", nodeName, err.Error()))
		return
	}

	state := ContainersDataSourceModel{
		ID:       types.StringValue(fmt.Sprintf("containers/%s", nodeName)),
		NodeName: types.StringValue(nodeName),
		VMIDs:    make([]types.Int64, len(ctList)),
		Names:    make([]types.String, len(ctList)),
		Statuses: make([]types.String, len(ctList)),
		Tags:     make([]types.String, len(ctList)),
		CPUs:     make([]types.Int64, len(ctList)),
		MaxMem:   make([]types.Int64, len(ctList)),
		MaxDisk:  make([]types.Int64, len(ctList)),
		Uptime:   make([]types.Int64, len(ctList)),
		Template: make([]types.Bool, len(ctList)),
	}

	for i, ct := range ctList {
		state.VMIDs[i] = types.Int64Value(int64(ct.VMID))
		state.Names[i] = types.StringValue(ct.Name)
		state.Statuses[i] = types.StringValue(ct.Status)
		state.Tags[i] = types.StringValue(ct.Tags)
		state.CPUs[i] = types.Int64Value(int64(ct.CPUs))
		state.MaxMem[i] = types.Int64Value(ct.MaxMem)
		state.MaxDisk[i] = types.Int64Value(ct.MaxDisk)
		state.Uptime[i] = types.Int64Value(ct.Uptime)
		state.Template[i] = types.BoolValue(ct.Type == "template")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
