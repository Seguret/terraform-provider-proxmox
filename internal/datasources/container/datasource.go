package container

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

var _ datasource.DataSource = &ContainerDataSource{}
var _ datasource.DataSourceWithConfigure = &ContainerDataSource{}

type ContainerDataSource struct {
	client *client.Client
}

type ContainerDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	NodeName    types.String `tfsdk:"node_name"`
	VMID        types.Int64  `tfsdk:"vmid"`
	Hostname    types.String `tfsdk:"hostname"`
	Description types.String `tfsdk:"description"`
	Tags        types.String `tfsdk:"tags"`
	Status      types.String `tfsdk:"status"`
	OSType      types.String `tfsdk:"os_type"`
	Cores       types.Int64  `tfsdk:"cores"`
	Memory      types.Int64  `tfsdk:"memory"`
	Swap        types.Int64  `tfsdk:"swap"`
	OnBoot      types.Bool   `tfsdk:"on_boot"`
	Template    types.Bool   `tfsdk:"template"`
	Unprivileged types.Bool  `tfsdk:"unprivileged"`
}

func NewDataSource() datasource.DataSource {
	return &ContainerDataSource{}
}

func (d *ContainerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_container"
}

func (d *ContainerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a specific Proxmox VE LXC container.",
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true},
			"node_name":   schema.StringAttribute{Description: "The node name.", Required: true},
			"vmid":        schema.Int64Attribute{Description: "The container ID.", Required: true},
			"hostname":    schema.StringAttribute{Description: "The container hostname.", Computed: true},
			"description": schema.StringAttribute{Description: "The container description.", Computed: true},
			"tags":        schema.StringAttribute{Description: "The container tags.", Computed: true},
			"status":      schema.StringAttribute{Description: "The container status.", Computed: true},
			"os_type":     schema.StringAttribute{Description: "The OS type.", Computed: true},
			"cores":       schema.Int64Attribute{Description: "Number of CPU cores.", Computed: true},
			"memory":      schema.Int64Attribute{Description: "Memory in MiB.", Computed: true},
			"swap":        schema.Int64Attribute{Description: "Swap in MiB.", Computed: true},
			"on_boot":     schema.BoolAttribute{Description: "Whether the container starts on host boot.", Computed: true},
			"template":    schema.BoolAttribute{Description: "Whether the container is a template.", Computed: true},
			"unprivileged": schema.BoolAttribute{Description: "Whether the container runs unprivileged.", Computed: true},
		},
	}
}

func (d *ContainerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ContainerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ContainerDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	vmid := int(config.VMID.ValueInt64())

	tflog.Debug(ctx, "Reading container", map[string]any{"node": node, "vmid": vmid})

	cfg, err := d.client.GetContainerConfig(ctx, node, vmid)
	if err != nil {
		resp.Diagnostics.AddError("Error reading container config",
			fmt.Sprintf("Unable to read container %d on node '%s': %s", vmid, node, err.Error()))
		return
	}

	status, err := d.client.GetContainerStatus(ctx, node, vmid)
	if err != nil {
		resp.Diagnostics.AddError("Error reading container status",
			fmt.Sprintf("Unable to read container status %d on node '%s': %s", vmid, node, err.Error()))
		return
	}

	state := ContainerDataSourceModel{
		ID:          types.StringValue(fmt.Sprintf("%s/%s", node, strconv.Itoa(vmid))),
		NodeName:    types.StringValue(node),
		VMID:        types.Int64Value(int64(vmid)),
		Hostname:    types.StringValue(cfg.Hostname),
		Description: types.StringValue(cfg.Description),
		Tags:        types.StringValue(cfg.Tags),
		Status:      types.StringValue(status.Status),
		OSType:      types.StringValue(cfg.OSType),
		Cores:       types.Int64Value(int64(cfg.Cores)),
		Memory:      types.Int64Value(int64(cfg.Memory)),
		Swap:        types.Int64Value(int64(cfg.Swap)),
	}

	if cfg.OnBoot != nil {
		state.OnBoot = types.BoolValue(*cfg.OnBoot == 1)
	} else {
		state.OnBoot = types.BoolValue(false)
	}
	if cfg.Template != nil {
		state.Template = types.BoolValue(*cfg.Template == 1)
	} else {
		state.Template = types.BoolValue(false)
	}
	if cfg.Unprivileged != nil {
		state.Unprivileged = types.BoolValue(*cfg.Unprivileged == 1)
	} else {
		state.Unprivileged = types.BoolValue(true)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
