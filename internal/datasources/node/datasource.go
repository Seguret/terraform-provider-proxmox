package node

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &NodeDataSource{}
var _ datasource.DataSourceWithConfigure = &NodeDataSource{}

type NodeDataSource struct {
	client *client.Client
}

type NodeDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	// Input
	NodeName types.String `tfsdk:"node_name"`

	// CPU
	CPUCores   types.Int64   `tfsdk:"cpu_cores"`
	CPUSockets types.Int64   `tfsdk:"cpu_sockets"`
	CPUThreads types.Int64   `tfsdk:"cpu_threads"`
	CPUModel   types.String  `tfsdk:"cpu_model"`
	CPUMHz     types.String  `tfsdk:"cpu_mhz"`
	CPUUsage   types.Float64 `tfsdk:"cpu_usage"`
	LoadAvg    []types.String `tfsdk:"load_average"`

	// Memory
	MemoryTotal types.Int64 `tfsdk:"memory_total"`
	MemoryUsed  types.Int64 `tfsdk:"memory_used"`
	MemoryFree  types.Int64 `tfsdk:"memory_free"`

	// Swap
	SwapTotal types.Int64 `tfsdk:"swap_total"`
	SwapUsed  types.Int64 `tfsdk:"swap_used"`
	SwapFree  types.Int64 `tfsdk:"swap_free"`

	// Root filesystem
	RootFSTotal types.Int64 `tfsdk:"rootfs_total"`
	RootFSUsed  types.Int64 `tfsdk:"rootfs_used"`
	RootFSFree  types.Int64 `tfsdk:"rootfs_free"`

	// System info
	Uptime     types.Int64  `tfsdk:"uptime"`
	KVersion   types.String `tfsdk:"kernel_version"`
	PVEVersion types.String `tfsdk:"pve_version"`

	// Boot info
	BootMode   types.String `tfsdk:"boot_mode"`
	SecureBoot types.Bool   `tfsdk:"secure_boot"`
}

func NewDataSource() datasource.DataSource {
	return &NodeDataSource{}
}

func (d *NodeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node"
}

func (d *NodeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the status and details of a specific Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the Proxmox VE node.",
				Required:    true,
			},

			// CPU
			"cpu_cores": schema.Int64Attribute{
				Description: "The number of CPU cores.",
				Computed:    true,
			},
			"cpu_sockets": schema.Int64Attribute{
				Description: "The number of CPU sockets.",
				Computed:    true,
			},
			"cpu_threads": schema.Int64Attribute{
				Description: "The number of CPU threads.",
				Computed:    true,
			},
			"cpu_model": schema.StringAttribute{
				Description: "The CPU model name.",
				Computed:    true,
			},
			"cpu_mhz": schema.StringAttribute{
				Description: "The CPU clock speed in MHz.",
				Computed:    true,
			},
			"cpu_usage": schema.Float64Attribute{
				Description: "The current CPU utilization (0.0-1.0).",
				Computed:    true,
			},
			"load_average": schema.ListAttribute{
				Description: "The system load averages (1, 5, 15 minutes).",
				Computed:    true,
				ElementType: types.StringType,
			},

			// Memory
			"memory_total": schema.Int64Attribute{
				Description: "Total memory in bytes.",
				Computed:    true,
			},
			"memory_used": schema.Int64Attribute{
				Description: "Used memory in bytes.",
				Computed:    true,
			},
			"memory_free": schema.Int64Attribute{
				Description: "Free memory in bytes.",
				Computed:    true,
			},

			// Swap
			"swap_total": schema.Int64Attribute{
				Description: "Total swap in bytes.",
				Computed:    true,
			},
			"swap_used": schema.Int64Attribute{
				Description: "Used swap in bytes.",
				Computed:    true,
			},
			"swap_free": schema.Int64Attribute{
				Description: "Free swap in bytes.",
				Computed:    true,
			},

			// Root filesystem
			"rootfs_total": schema.Int64Attribute{
				Description: "Total root filesystem size in bytes.",
				Computed:    true,
			},
			"rootfs_used": schema.Int64Attribute{
				Description: "Used root filesystem in bytes.",
				Computed:    true,
			},
			"rootfs_free": schema.Int64Attribute{
				Description: "Free root filesystem in bytes.",
				Computed:    true,
			},

			// System
			"uptime": schema.Int64Attribute{
				Description: "The node uptime in seconds.",
				Computed:    true,
			},
			"kernel_version": schema.StringAttribute{
				Description: "The kernel version string.",
				Computed:    true,
			},
			"pve_version": schema.StringAttribute{
				Description: "The Proxmox VE version string.",
				Computed:    true,
			},

			// Boot
			"boot_mode": schema.StringAttribute{
				Description: "The boot mode (e.g., 'efi' or 'bios').",
				Computed:    true,
			},
			"secure_boot": schema.BoolAttribute{
				Description: "Whether secure boot is enabled.",
				Computed:    true,
			},
		},
	}
}

func (d *NodeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config NodeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nodeName := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading Proxmox VE node status", map[string]any{"node": nodeName})

	status, err := d.client.GetNodeStatus(ctx, nodeName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Proxmox VE Node Status",
			fmt.Sprintf("An error occurred while reading node '%s': %s", nodeName, err.Error()),
		)
		return
	}

	state := NodeDataSourceModel{
		ID:       types.StringValue(nodeName),
		NodeName: types.StringValue(nodeName),
		CPUUsage: types.Float64Value(status.CPU),
		Uptime:   types.Int64Value(status.Uptime),
		KVersion: types.StringValue(status.KVersion),
		PVEVersion: types.StringValue(status.PVEVersion),
	}

	// CPU info
	if status.CPUInfo != nil {
		state.CPUCores = types.Int64Value(int64(status.CPUInfo.Cores))
		state.CPUSockets = types.Int64Value(int64(status.CPUInfo.Sockets))
		state.CPUThreads = types.Int64Value(int64(status.CPUInfo.Threads))
		state.CPUModel = types.StringValue(status.CPUInfo.Model)
		state.CPUMHz = types.StringValue(status.CPUInfo.MHz)
	} else {
		state.CPUCores = types.Int64Value(0)
		state.CPUSockets = types.Int64Value(0)
		state.CPUThreads = types.Int64Value(0)
		state.CPUModel = types.StringValue("")
		state.CPUMHz = types.StringValue("")
	}

	// Load average
	if len(status.LoadAvg) > 0 {
		state.LoadAvg = make([]types.String, len(status.LoadAvg))
		for i, v := range status.LoadAvg {
			state.LoadAvg[i] = types.StringValue(v)
		}
	} else {
		state.LoadAvg = []types.String{}
	}

	// Memory
	if status.Memory != nil {
		state.MemoryTotal = types.Int64Value(status.Memory.Total)
		state.MemoryUsed = types.Int64Value(status.Memory.Used)
		state.MemoryFree = types.Int64Value(status.Memory.Free)
	} else {
		state.MemoryTotal = types.Int64Value(0)
		state.MemoryUsed = types.Int64Value(0)
		state.MemoryFree = types.Int64Value(0)
	}

	// Swap
	if status.Swap != nil {
		state.SwapTotal = types.Int64Value(status.Swap.Total)
		state.SwapUsed = types.Int64Value(status.Swap.Used)
		state.SwapFree = types.Int64Value(status.Swap.Free)
	} else {
		state.SwapTotal = types.Int64Value(0)
		state.SwapUsed = types.Int64Value(0)
		state.SwapFree = types.Int64Value(0)
	}

	// Root filesystem
	if status.RootFS != nil {
		state.RootFSTotal = types.Int64Value(status.RootFS.Total)
		state.RootFSUsed = types.Int64Value(status.RootFS.Used)
		state.RootFSFree = types.Int64Value(status.RootFS.Free)
	} else {
		state.RootFSTotal = types.Int64Value(0)
		state.RootFSUsed = types.Int64Value(0)
		state.RootFSFree = types.Int64Value(0)
	}

	// Boot info
	if status.BootInfo != nil {
		state.BootMode = types.StringValue(status.BootInfo.Mode)
		state.SecureBoot = types.BoolValue(status.BootInfo.SecBoot == 1)
	} else {
		state.BootMode = types.StringValue("")
		state.SecureBoot = types.BoolValue(false)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
