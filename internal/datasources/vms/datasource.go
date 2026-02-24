package vms

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &VMsDataSource{}
var _ datasource.DataSourceWithConfigure = &VMsDataSource{}

type VMsDataSource struct {
	client *client.Client
}

type VMsDataSourceModel struct {
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
	return &VMsDataSource{}
}

func (d *VMsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_vms"
}

func (d *VMsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of VMs on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"node_name": schema.StringAttribute{
				Description: "The node name.",
				Required:    true,
			},
			"vmids": schema.ListAttribute{
				Description: "The VM IDs.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"names": schema.ListAttribute{
				Description: "The VM names.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"statuses": schema.ListAttribute{
				Description: "The VM statuses.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"tags": schema.ListAttribute{
				Description: "The VM tags.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"cpus": schema.ListAttribute{
				Description: "The number of CPUs for each VM.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"max_memory": schema.ListAttribute{
				Description: "The maximum memory in bytes for each VM.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"max_disk": schema.ListAttribute{
				Description: "The maximum disk size in bytes for each VM.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"uptime": schema.ListAttribute{
				Description: "The uptime in seconds for each VM.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"template": schema.ListAttribute{
				Description: "Whether each VM is a template.",
				Computed:    true,
				ElementType: types.BoolType,
			},
		},
	}
}

func (d *VMsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VMsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config VMsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nodeName := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading VMs list", map[string]any{"node": nodeName})

	vmList, err := d.client.GetVMs(ctx, nodeName)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read VMs",
			fmt.Sprintf("Error reading VMs on node '%s': %s", nodeName, err.Error()))
		return
	}

	state := VMsDataSourceModel{
		ID:       types.StringValue(fmt.Sprintf("vms/%s", nodeName)),
		NodeName: types.StringValue(nodeName),
		VMIDs:    make([]types.Int64, len(vmList)),
		Names:    make([]types.String, len(vmList)),
		Statuses: make([]types.String, len(vmList)),
		Tags:     make([]types.String, len(vmList)),
		CPUs:     make([]types.Int64, len(vmList)),
		MaxMem:   make([]types.Int64, len(vmList)),
		MaxDisk:  make([]types.Int64, len(vmList)),
		Uptime:   make([]types.Int64, len(vmList)),
		Template: make([]types.Bool, len(vmList)),
	}

	for i, vm := range vmList {
		state.VMIDs[i] = types.Int64Value(int64(vm.VMID))
		state.Names[i] = types.StringValue(vm.Name)
		state.Statuses[i] = types.StringValue(vm.Status)
		state.Tags[i] = types.StringValue(vm.Tags)
		state.CPUs[i] = types.Int64Value(int64(vm.CPUs))
		state.MaxMem[i] = types.Int64Value(vm.MaxMem)
		state.MaxDisk[i] = types.Int64Value(vm.MaxDisk)
		state.Uptime[i] = types.Int64Value(vm.Uptime)
		state.Template[i] = types.BoolValue(vm.Template == 1)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
