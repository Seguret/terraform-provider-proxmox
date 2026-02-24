package container

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &ContainerResource{}
var _ resource.ResourceWithConfigure = &ContainerResource{}
var _ resource.ResourceWithImportState = &ContainerResource{}

type ContainerResource struct {
	client *client.Client
}

type ContainerResourceModel struct {
	ID          types.String `tfsdk:"id"`
	NodeName    types.String `tfsdk:"node_name"`
	VMID        types.Int64  `tfsdk:"vmid"`
	Hostname    types.String `tfsdk:"hostname"`
	Description types.String `tfsdk:"description"`
	Tags        types.String `tfsdk:"tags"`
	OSTemplate  types.String `tfsdk:"os_template"`
	OSType      types.String `tfsdk:"os_type"`
	OnBoot      types.Bool   `tfsdk:"on_boot"`
	Started     types.Bool   `tfsdk:"started"`
	Protection  types.Bool   `tfsdk:"protection"`
	Unprivileged types.Bool  `tfsdk:"unprivileged"`
	Pool        types.String `tfsdk:"pool"`
	Features    types.String `tfsdk:"features"`

	// CPU
	Cores    types.Int64 `tfsdk:"cpu_cores"`
	CPULimit types.Int64 `tfsdk:"cpu_limit"`
	CPUUnits types.Int64 `tfsdk:"cpu_units"`

	// Memory
	Memory types.Int64 `tfsdk:"memory"`
	Swap   types.Int64 `tfsdk:"swap"`

	// Root filesystem
	RootFS types.String `tfsdk:"rootfs"`

	// Network
	Net0 types.String `tfsdk:"net0"`
	Net1 types.String `tfsdk:"net1"`
	Net2 types.String `tfsdk:"net2"`
	Net3 types.String `tfsdk:"net3"`

	// Mount points
	MP0 types.String `tfsdk:"mp0"`
	MP1 types.String `tfsdk:"mp1"`
	MP2 types.String `tfsdk:"mp2"`

	// DNS
	Nameserver   types.String `tfsdk:"nameserver"`
	Searchdomain types.String `tfsdk:"searchdomain"`

	// Auth
	Password types.String `tfsdk:"password"`
	SSHKeys  types.String `tfsdk:"ssh_keys"`

	// Console
	Console types.Bool  `tfsdk:"console"`
	TTY     types.Int64 `tfsdk:"tty"`

	// Clone
	CloneVMID types.Int64 `tfsdk:"clone_vmid"`
	FullClone types.Bool  `tfsdk:"full_clone"`

	// Read-only
	Status   types.String `tfsdk:"status"`
	Template types.Bool   `tfsdk:"template"`
}

func NewResource() resource.Resource {
	return &ContainerResource{}
}

func (r *ContainerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_container"
}

func (r *ContainerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE LXC container.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"node_name": schema.StringAttribute{
				Description: "The node on which to create the container.",
				Required:    true,
			},
			"vmid": schema.Int64Attribute{
				Description: "The container ID. If not set, the next available ID will be used.",
				Optional:    true,
				Computed:    true,
			},
			"hostname": schema.StringAttribute{
				Description: "The container hostname.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The container description.",
				Optional:    true,
				Computed:    true,
			},
			"tags": schema.StringAttribute{
				Description: "Tags for the container (semicolon-separated).",
				Optional:    true,
				Computed:    true,
			},
			"os_template": schema.StringAttribute{
				Description: "The OS template to use (e.g., 'local:vztmpl/ubuntu-22.04-standard_22.04-1_amd64.tar.zst').",
				Optional:    true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"os_type": schema.StringAttribute{
				Description: "The OS type (debian, ubuntu, centos, fedora, opensuse, archlinux, alpine, gentoo, nixos, unmanaged).",
				Optional:    true,
				Computed:    true,
			},
			"on_boot": schema.BoolAttribute{
				Description: "Whether to start the container on host boot.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"started": schema.BoolAttribute{
				Description: "Whether the container should be started after creation.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"protection": schema.BoolAttribute{
				Description: "Whether the container is protected from removal.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"unprivileged": schema.BoolAttribute{
				Description: "Whether to create an unprivileged container.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"pool": schema.StringAttribute{
				Description: "The resource pool to add the container to.",
				Optional:    true,
			},
			"features": schema.StringAttribute{
				Description: "Container feature flags (e.g., 'nesting=1,keyctl=1').",
				Optional:    true,
				Computed:    true,
			},
			// CPU
			"cpu_cores": schema.Int64Attribute{
				Description: "Number of CPU cores.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
			},
			"cpu_limit": schema.Int64Attribute{
				Description: "CPU usage limit (0 = unlimited).",
				Optional:    true,
				Computed:    true,
			},
			"cpu_units": schema.Int64Attribute{
				Description: "CPU weight for a container (relative weight vs other containers).",
				Optional:    true,
				Computed:    true,
			},
			// Memory
			"memory": schema.Int64Attribute{
				Description: "Memory in MiB.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(512),
			},
			"swap": schema.Int64Attribute{
				Description: "Swap in MiB.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(512),
			},
			// Root filesystem
			"rootfs": schema.StringAttribute{
				Description: "Root filesystem configuration (e.g., 'local-lvm:8').",
				Optional:    true,
				Computed:    true,
			},
			// Network
			"net0": schema.StringAttribute{
				Description: "Network interface 0 (e.g., 'name=eth0,bridge=vmbr0,ip=dhcp').",
				Optional:    true,
				Computed:    true,
			},
			"net1": schema.StringAttribute{Description: "Network interface 1.", Optional: true, Computed: true},
			"net2": schema.StringAttribute{Description: "Network interface 2.", Optional: true, Computed: true},
			"net3": schema.StringAttribute{Description: "Network interface 3.", Optional: true, Computed: true},
			// Mount points
			"mp0": schema.StringAttribute{Description: "Mount point 0.", Optional: true, Computed: true},
			"mp1": schema.StringAttribute{Description: "Mount point 1.", Optional: true, Computed: true},
			"mp2": schema.StringAttribute{Description: "Mount point 2.", Optional: true, Computed: true},
			// DNS
			"nameserver": schema.StringAttribute{
				Description: "DNS nameserver.",
				Optional:    true,
				Computed:    true,
			},
			"searchdomain": schema.StringAttribute{
				Description: "DNS search domain.",
				Optional:    true,
				Computed:    true,
			},
			// Auth
			"password": schema.StringAttribute{
				Description: "Root password for the container.",
				Optional:    true,
				Sensitive:   true,
			},
			"ssh_keys": schema.StringAttribute{
				Description: "SSH public keys.",
				Optional:    true,
				Computed:    true,
			},
			// Console
			"console": schema.BoolAttribute{
				Description: "Whether to attach a console device.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"tty": schema.Int64Attribute{
				Description: "Number of TTY devices (0-6).",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(2),
			},
			// Clone
			"clone_vmid": schema.Int64Attribute{
				Description: "VMID of the container to clone from.",
				Optional:    true,
			},
			"full_clone": schema.BoolAttribute{
				Description: "Whether to do a full clone.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			// Read-only
			"status":   schema.StringAttribute{Computed: true, Description: "The container current status."},
			"template": schema.BoolAttribute{Computed: true, Description: "Whether the container is a template."},
		},
	}
}

func (r *ContainerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *ContainerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ContainerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	vmid := int(plan.VMID.ValueInt64())
	if plan.VMID.IsNull() || plan.VMID.IsUnknown() || vmid == 0 {
		nextID, err := r.client.GetNextVMID(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Error getting next VMID", err.Error())
			return
		}
		vmid = nextID
	}

	onboot := boolToInt(plan.OnBoot.ValueBool())
	protection := boolToInt(plan.Protection.ValueBool())
	unprivileged := boolToInt(plan.Unprivileged.ValueBool())
	console := boolToInt(plan.Console.ValueBool())

	createReq := &models.ContainerCreateRequest{
		VMID:         vmid,
		Hostname:     plan.Hostname.ValueString(),
		Description:  plan.Description.ValueString(),
		Tags:         plan.Tags.ValueString(),
		OSTemplate:   plan.OSTemplate.ValueString(),
		OSType:       plan.OSType.ValueString(),
		OnBoot:       &onboot,
		Protection:   &protection,
		Unprivileged: &unprivileged,
		Pool:         plan.Pool.ValueString(),
		Features:     plan.Features.ValueString(),
		Cores:        int(plan.Cores.ValueInt64()),
		CPULimit:     int(plan.CPULimit.ValueInt64()),
		CPUUnits:     int(plan.CPUUnits.ValueInt64()),
		Memory:       int(plan.Memory.ValueInt64()),
		Swap:         int(plan.Swap.ValueInt64()),
		RootFS:       plan.RootFS.ValueString(),
		Net0:         plan.Net0.ValueString(),
		Net1:         plan.Net1.ValueString(),
		Net2:         plan.Net2.ValueString(),
		Net3:         plan.Net3.ValueString(),
		MP0:          plan.MP0.ValueString(),
		MP1:          plan.MP1.ValueString(),
		MP2:          plan.MP2.ValueString(),
		Nameserver:   plan.Nameserver.ValueString(),
		Searchdomain: plan.Searchdomain.ValueString(),
		Password:     plan.Password.ValueString(),
		SSHKeys:      plan.SSHKeys.ValueString(),
		Console:      &console,
		TTY:          int(plan.TTY.ValueInt64()),
	}

	// wire up clone settings if this is a clone operation
	if !plan.CloneVMID.IsNull() && plan.CloneVMID.ValueInt64() > 0 {
		cloneSource := int(plan.CloneVMID.ValueInt64())
		full := boolToInt(plan.FullClone.ValueBool())
		createReq.Clone = &cloneSource
		createReq.Full = &full
	}

	tflog.Debug(ctx, "Creating container", map[string]any{"node": node, "vmid": vmid})

	upid, err := r.client.CreateContainer(ctx, node, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating container", err.Error())
		return
	}
	if upid != "" {
		if err := r.client.WaitForTask(ctx, node, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for container creation", err.Error())
			return
		}
	}

	plan.VMID = types.Int64Value(int64(vmid))
	plan.ID = types.StringValue(fmt.Sprintf("%s/%d", node, vmid))

	if plan.Started.ValueBool() {
		status, err := r.client.GetContainerStatus(ctx, node, vmid)
		if err == nil && status.Status != "running" {
			upid, err := r.client.StartContainer(ctx, node, vmid)
			if err != nil {
				resp.Diagnostics.AddError("Error starting container", err.Error())
				return
			}
			if err := r.client.WaitForTask(ctx, node, upid); err != nil {
				resp.Diagnostics.AddError("Error waiting for container start", err.Error())
				return
			}
		}
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ContainerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ContainerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ContainerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ContainerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	vmid := int(plan.VMID.ValueInt64())

	configMap := r.buildConfigMap(&plan)
	if len(configMap) > 0 {
		if err := r.client.UpdateContainerConfig(ctx, node, vmid, configMap); err != nil {
			resp.Diagnostics.AddError("Error updating container", err.Error())
			return
		}
	}

	status, err := r.client.GetContainerStatus(ctx, node, vmid)
	if err != nil {
		resp.Diagnostics.AddError("Error reading container status", err.Error())
		return
	}

	if plan.Started.ValueBool() && status.Status != "running" {
		upid, err := r.client.StartContainer(ctx, node, vmid)
		if err != nil {
			resp.Diagnostics.AddError("Error starting container", err.Error())
			return
		}
		if err := r.client.WaitForTask(ctx, node, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for container start", err.Error())
			return
		}
	} else if !plan.Started.ValueBool() && status.Status == "running" {
		upid, err := r.client.ShutdownContainer(ctx, node, vmid)
		if err != nil {
			resp.Diagnostics.AddError("Error shutting down container", err.Error())
			return
		}
		if err := r.client.WaitForTask(ctx, node, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for container shutdown", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%d", node, vmid))
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ContainerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ContainerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := state.NodeName.ValueString()
	vmid := int(state.VMID.ValueInt64())

	status, err := r.client.GetContainerStatus(ctx, node, vmid)
	if err != nil {
		resp.Diagnostics.AddError("Error reading container status", err.Error())
		return
	}
	if status.Status == "running" {
		upid, err := r.client.StopContainer(ctx, node, vmid)
		if err != nil {
			resp.Diagnostics.AddError("Error stopping container", err.Error())
			return
		}
		if err := r.client.WaitForTask(ctx, node, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for container stop", err.Error())
			return
		}
	}

	upid, err := r.client.DeleteContainer(ctx, node, vmid)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting container", err.Error())
		return
	}
	if upid != "" {
		if err := r.client.WaitForTask(ctx, node, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for container deletion", err.Error())
		}
	}
}

func (r *ContainerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Import ID must be in format 'node_name/vmid'")
		return
	}
	vmid, err := strconv.Atoi(parts[1])
	if err != nil {
		resp.Diagnostics.AddError("Invalid VMID", err.Error())
		return
	}
	state := ContainerResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(parts[0]),
		VMID:     types.Int64Value(int64(vmid)),
		Started:  types.BoolValue(true),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ContainerResource) readIntoModel(ctx context.Context, model *ContainerResourceModel, diagnostics *diag.Diagnostics) {
	node := model.NodeName.ValueString()
	vmid := int(model.VMID.ValueInt64())

	cfg, err := r.client.GetContainerConfig(ctx, node, vmid)
	if err != nil {
		diagnostics.AddError("Error reading container config", err.Error())
		return
	}
	status, err := r.client.GetContainerStatus(ctx, node, vmid)
	if err != nil {
		diagnostics.AddError("Error reading container status", err.Error())
		return
	}

	model.Hostname = types.StringValue(cfg.Hostname)
	model.Description = types.StringValue(cfg.Description)
	model.Tags = types.StringValue(cfg.Tags)
	model.OSType = types.StringValue(cfg.OSType)
	model.Features = types.StringValue(cfg.Features)
	model.RootFS = types.StringValue(cfg.RootFS)
	model.Net0 = types.StringValue(cfg.Net0)
	model.Net1 = types.StringValue(cfg.Net1)
	model.Net2 = types.StringValue(cfg.Net2)
	model.Net3 = types.StringValue(cfg.Net3)
	model.MP0 = types.StringValue(cfg.MP0)
	model.MP1 = types.StringValue(cfg.MP1)
	model.MP2 = types.StringValue(cfg.MP2)
	model.Nameserver = types.StringValue(cfg.Nameserver)
	model.Searchdomain = types.StringValue(cfg.Searchdomain)
	model.SSHKeys = types.StringValue(cfg.SSHKeys)

	if cfg.OnBoot != nil {
		model.OnBoot = types.BoolValue(*cfg.OnBoot == 1)
	}
	if cfg.Protection != nil {
		model.Protection = types.BoolValue(*cfg.Protection == 1)
	}
	if cfg.Unprivileged != nil {
		model.Unprivileged = types.BoolValue(*cfg.Unprivileged == 1)
	}
	if cfg.Template != nil {
		model.Template = types.BoolValue(*cfg.Template == 1)
	} else {
		model.Template = types.BoolValue(false)
	}
	if cfg.Console != nil {
		model.Console = types.BoolValue(*cfg.Console == 1)
	}
	if cfg.Cores > 0 {
		model.Cores = types.Int64Value(int64(cfg.Cores))
	}
	if cfg.Memory > 0 {
		model.Memory = types.Int64Value(int64(cfg.Memory))
	}
	model.Swap = types.Int64Value(int64(cfg.Swap))
	model.CPULimit = types.Int64Value(int64(cfg.CPULimit))
	model.CPUUnits = types.Int64Value(int64(cfg.CPUUnits))
	model.TTY = types.Int64Value(int64(cfg.TTY))
	model.Status = types.StringValue(status.Status)
}

func (r *ContainerResource) buildConfigMap(plan *ContainerResourceModel) map[string]interface{} {
	m := make(map[string]interface{})
	set := func(k, v string) {
		if v != "" {
			m[k] = v
		}
	}
	set("hostname", plan.Hostname.ValueString())
	set("description", plan.Description.ValueString())
	set("tags", plan.Tags.ValueString())
	set("ostype", plan.OSType.ValueString())
	set("features", plan.Features.ValueString())
	set("nameserver", plan.Nameserver.ValueString())
	set("searchdomain", plan.Searchdomain.ValueString())
	set("net0", plan.Net0.ValueString())
	set("net1", plan.Net1.ValueString())
	set("net2", plan.Net2.ValueString())
	set("net3", plan.Net3.ValueString())
	set("mp0", plan.MP0.ValueString())
	set("mp1", plan.MP1.ValueString())
	set("mp2", plan.MP2.ValueString())

	if !plan.OnBoot.IsNull() {
		m["onboot"] = boolToInt(plan.OnBoot.ValueBool())
	}
	if !plan.Protection.IsNull() {
		m["protection"] = boolToInt(plan.Protection.ValueBool())
	}
	if !plan.Console.IsNull() {
		m["console"] = boolToInt(plan.Console.ValueBool())
	}
	if v := plan.Cores.ValueInt64(); v > 0 {
		m["cores"] = v
	}
	if v := plan.Memory.ValueInt64(); v > 0 {
		m["memory"] = v
	}
	if !plan.Swap.IsNull() {
		m["swap"] = plan.Swap.ValueInt64()
	}
	if !plan.CPULimit.IsNull() {
		m["cpulimit"] = plan.CPULimit.ValueInt64()
	}
	if !plan.CPUUnits.IsNull() {
		m["cpuunits"] = plan.CPUUnits.ValueInt64()
	}
	if !plan.TTY.IsNull() {
		m["tty"] = plan.TTY.ValueInt64()
	}
	return m
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
