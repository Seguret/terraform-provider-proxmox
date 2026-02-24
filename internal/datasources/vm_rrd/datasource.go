package vm_rrd

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &VMRRDDataSource{}
var _ datasource.DataSourceWithConfigure = &VMRRDDataSource{}

type VMRRDDataSource struct {
	client *client.Client
}

type VMRRDDataPointModel struct {
	Time      types.Int64   `tfsdk:"time"`
	CPU       types.Float64 `tfsdk:"cpu"`
	MaxCPU    types.Float64 `tfsdk:"maxcpu"`
	Mem       types.Float64 `tfsdk:"mem"`
	MaxMem    types.Float64 `tfsdk:"maxmem"`
	NetIn     types.Float64 `tfsdk:"netin"`
	NetOut    types.Float64 `tfsdk:"netout"`
	DiskRead  types.Float64 `tfsdk:"diskread"`
	DiskWrite types.Float64 `tfsdk:"diskwrite"`
	LoadAvg   types.Float64 `tfsdk:"loadavg"`
	SwapTotal types.Float64 `tfsdk:"swaptotal"`
	SwapUsed  types.Float64 `tfsdk:"swapused"`
}

type VMRRDDataSourceModel struct {
	ID         types.String          `tfsdk:"id"`
	NodeName   types.String          `tfsdk:"node_name"`
	VMID       types.Int64           `tfsdk:"vmid"`
	Timeframe  types.String          `tfsdk:"timeframe"`
	DataPoints []VMRRDDataPointModel `tfsdk:"data_points"`
}

func NewDataSource() datasource.DataSource {
	return &VMRRDDataSource{}
}

func (d *VMRRDDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_vm_rrd"
}

func (d *VMRRDDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves RRD performance data for a Proxmox VE virtual machine.",
		Attributes: map[string]schema.Attribute{
			"id":        schema.StringAttribute{Computed: true},
			"node_name": schema.StringAttribute{Description: "The name of the node.", Required: true},
			"vmid":      schema.Int64Attribute{Description: "The VM ID.", Required: true},
			"timeframe": schema.StringAttribute{
				Description: "The timeframe for RRD data (hour, day, week, month, year).",
				Required:    true,
			},
			"data_points": schema.ListNestedAttribute{
				Description: "RRD data points.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"time":      schema.Int64Attribute{Computed: true, Description: "Unix timestamp."},
						"cpu":       schema.Float64Attribute{Computed: true, Description: "CPU usage (0-1)."},
						"maxcpu":    schema.Float64Attribute{Computed: true, Description: "Number of CPUs."},
						"mem":       schema.Float64Attribute{Computed: true, Description: "Memory used in bytes."},
						"maxmem":    schema.Float64Attribute{Computed: true, Description: "Total memory in bytes."},
						"netin":     schema.Float64Attribute{Computed: true, Description: "Network in (bytes/s)."},
						"netout":    schema.Float64Attribute{Computed: true, Description: "Network out (bytes/s)."},
						"diskread":  schema.Float64Attribute{Computed: true, Description: "Disk read (bytes/s)."},
						"diskwrite": schema.Float64Attribute{Computed: true, Description: "Disk write (bytes/s)."},
						"loadavg":   schema.Float64Attribute{Computed: true, Description: "1-minute load average."},
						"swaptotal": schema.Float64Attribute{Computed: true, Description: "Total swap in bytes."},
						"swapused":  schema.Float64Attribute{Computed: true, Description: "Used swap in bytes."},
					},
				},
			},
		},
	}
}

func (d *VMRRDDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VMRRDDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config VMRRDDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	vmid := int(config.VMID.ValueInt64())
	timeframe := config.Timeframe.ValueString()

	tflog.Debug(ctx, "Reading VM RRD data", map[string]any{"node": node, "vmid": vmid, "timeframe": timeframe})

	points, err := d.client.GetVMRRDData(ctx, node, vmid, timeframe)
	if err != nil {
		resp.Diagnostics.AddError("Error reading VM RRD data", err.Error())
		return
	}

	state := VMRRDDataSourceModel{
		ID:         types.StringValue(fmt.Sprintf("%s-%d-rrd-%s", node, vmid, timeframe)),
		NodeName:   config.NodeName,
		VMID:       config.VMID,
		Timeframe:  config.Timeframe,
		DataPoints: make([]VMRRDDataPointModel, len(points)),
	}

	for i, p := range points {
		state.DataPoints[i] = VMRRDDataPointModel{
			Time:      types.Int64Value(p.Time),
			CPU:       types.Float64Value(p.CPU),
			MaxCPU:    types.Float64Value(p.MaxCPU),
			Mem:       types.Float64Value(p.Mem),
			MaxMem:    types.Float64Value(p.MaxMem),
			NetIn:     types.Float64Value(p.NetIn),
			NetOut:    types.Float64Value(p.NetOut),
			DiskRead:  types.Float64Value(p.DiskRead),
			DiskWrite: types.Float64Value(p.DiskWrite),
			LoadAvg:   types.Float64Value(p.LoadAvg),
			SwapTotal: types.Float64Value(p.SwapTotal),
			SwapUsed:  types.Float64Value(p.SwapUsed),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
