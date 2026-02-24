package vm

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &VMDataSource{}
var _ datasource.DataSourceWithConfigure = &VMDataSource{}

type VMDataSource struct {
	client *client.Client
}

type VMDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	NodeName    types.String `tfsdk:"node_name"`
	VMID        types.Int64  `tfsdk:"vmid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Tags        types.String `tfsdk:"tags"`
	Status      types.String `tfsdk:"status"`
	OSType      types.String `tfsdk:"os_type"`
	BIOS        types.String `tfsdk:"bios"`
	Machine     types.String `tfsdk:"machine"`
	Sockets     types.Int64  `tfsdk:"cpu_sockets"`
	Cores       types.Int64  `tfsdk:"cpu_cores"`
	CPUType     types.String `tfsdk:"cpu_type"`
	Memory      types.Int64  `tfsdk:"memory"`
	Template    types.Bool   `tfsdk:"template"`
	OnBoot      types.Bool   `tfsdk:"on_boot"`
}

func NewDataSource() datasource.DataSource {
	return &VMDataSource{}
}

func (d *VMDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_vm"
}

func (d *VMDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a specific Proxmox VE virtual machine.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"node_name": schema.StringAttribute{
				Description: "The node name.",
				Required:    true,
			},
			"vmid": schema.Int64Attribute{
				Description: "The VM ID.",
				Required:    true,
			},
			"name":        schema.StringAttribute{Description: "The VM name.", Computed: true},
			"description": schema.StringAttribute{Description: "The VM description.", Computed: true},
			"tags":        schema.StringAttribute{Description: "The VM tags.", Computed: true},
			"status":      schema.StringAttribute{Description: "The VM status.", Computed: true},
			"os_type":     schema.StringAttribute{Description: "The OS type.", Computed: true},
			"bios":        schema.StringAttribute{Description: "The BIOS type.", Computed: true},
			"machine":     schema.StringAttribute{Description: "The machine type.", Computed: true},
			"cpu_sockets": schema.Int64Attribute{Description: "Number of CPU sockets.", Computed: true},
			"cpu_cores":   schema.Int64Attribute{Description: "Number of CPU cores per socket.", Computed: true},
			"cpu_type":    schema.StringAttribute{Description: "The CPU type.", Computed: true},
			"memory":      schema.Int64Attribute{Description: "Memory in MiB.", Computed: true},
			"template":    schema.BoolAttribute{Description: "Whether the VM is a template.", Computed: true},
			"on_boot":     schema.BoolAttribute{Description: "Whether the VM starts on host boot.", Computed: true},
		},
	}
}

func (d *VMDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *VMDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config VMDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	vmid := int(config.VMID.ValueInt64())

	tflog.Debug(ctx, "Reading VM", map[string]any{"node": node, "vmid": vmid})

	cfg, err := d.client.GetVMConfig(ctx, node, vmid)
	if err != nil {
		resp.Diagnostics.AddError("Error reading VM config",
			fmt.Sprintf("Unable to read VM %d on node '%s': %s", vmid, node, err.Error()))
		return
	}

	status, err := d.client.GetVMStatus(ctx, node, vmid)
	if err != nil {
		resp.Diagnostics.AddError("Error reading VM status",
			fmt.Sprintf("Unable to read VM status %d on node '%s': %s", vmid, node, err.Error()))
		return
	}

	state := VMDataSourceModel{
		ID:          types.StringValue(fmt.Sprintf("%s/%s", node, strconv.Itoa(vmid))),
		NodeName:    types.StringValue(node),
		VMID:        types.Int64Value(int64(vmid)),
		Name:        types.StringValue(cfg.Name),
		Description: types.StringValue(cfg.Description),
		Tags:        types.StringValue(cfg.Tags),
		Status:      types.StringValue(status.Status),
		OSType:      types.StringValue(cfg.OSType),
		BIOS:        types.StringValue(cfg.BIOS),
		Machine:     types.StringValue(cfg.Machine),
		Sockets:     types.Int64Value(int64(cfg.Sockets)),
		Cores:       types.Int64Value(int64(cfg.Cores)),
		CPUType:     types.StringValue(cfg.CPUType),
		Memory:      types.Int64Value(int64(cfg.Memory)),
	}

	if cfg.Template != nil {
		state.Template = types.BoolValue(*cfg.Template == 1)
	} else {
		state.Template = types.BoolValue(false)
	}

	if cfg.OnBoot != nil {
		state.OnBoot = types.BoolValue(*cfg.OnBoot == 1)
	} else {
		state.OnBoot = types.BoolValue(false)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
