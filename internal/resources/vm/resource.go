package vm

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &VMResource{}
var _ resource.ResourceWithConfigure = &VMResource{}
var _ resource.ResourceWithImportState = &VMResource{}

type VMResource struct {
	client *client.Client
}

type VMResourceModel struct {
	ID          types.String `tfsdk:"id"`
	NodeName    types.String `tfsdk:"node_name"`
	VMID        types.Int64  `tfsdk:"vmid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Tags        types.String `tfsdk:"tags"`
	OnBoot      types.Bool   `tfsdk:"on_boot"`
	Started     types.Bool   `tfsdk:"started"`
	Protection  types.Bool   `tfsdk:"protection"`
	Agent       types.Bool   `tfsdk:"agent"`
	OSType      types.String `tfsdk:"os_type"`
	BIOS        types.String `tfsdk:"bios"`
	Machine     types.String `tfsdk:"machine"`
	SCSIHw      types.String `tfsdk:"scsi_hw"`
	Boot        types.String `tfsdk:"boot"`
	Pool        types.String `tfsdk:"pool"`

	// CPU
	Sockets types.Int64  `tfsdk:"cpu_sockets"`
	Cores   types.Int64  `tfsdk:"cpu_cores"`
	CPUType types.String `tfsdk:"cpu_type"`

	// Memory
	Memory  types.Int64 `tfsdk:"memory"`
	Balloon types.Int64 `tfsdk:"balloon"`

	// VGA
	VGA types.String `tfsdk:"vga"`

	// Disks
	SCSI0   types.String `tfsdk:"scsi0"`
	SCSI1   types.String `tfsdk:"scsi1"`
	SCSI2   types.String `tfsdk:"scsi2"`
	SCSI3   types.String `tfsdk:"scsi3"`
	VirtIO0 types.String `tfsdk:"virtio0"`
	VirtIO1 types.String `tfsdk:"virtio1"`
	IDE0    types.String `tfsdk:"ide0"`
	IDE2    types.String `tfsdk:"ide2"`
	EFIDisk0 types.String `tfsdk:"efidisk0"`
	TPMState0 types.String `tfsdk:"tpmstate0"`

	// Network
	Net0 types.String `tfsdk:"net0"`
	Net1 types.String `tfsdk:"net1"`
	Net2 types.String `tfsdk:"net2"`
	Net3 types.String `tfsdk:"net3"`

	// Cloud-init
	CIUser       types.String `tfsdk:"ci_user"`
	CIPassword   types.String `tfsdk:"ci_password"`
	CIType       types.String `tfsdk:"ci_type"`
	IPConfig0    types.String `tfsdk:"ipconfig0"`
	IPConfig1    types.String `tfsdk:"ipconfig1"`
	Nameserver   types.String `tfsdk:"nameserver"`
	Searchdomain types.String `tfsdk:"searchdomain"`
	SSHKeys      types.String `tfsdk:"ssh_keys"`

	// Serial
	Serial0 types.String `tfsdk:"serial0"`

	// Clone source
	CloneVMID types.Int64 `tfsdk:"clone_vmid"`
	FullClone types.Bool  `tfsdk:"full_clone"`

	// Read-only
	Status   types.String `tfsdk:"status"`
	Template types.Bool   `tfsdk:"template"`
}

func NewResource() resource.Resource {
	return &VMResource{}
}

func (r *VMResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_vm"
}

func (r *VMResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE virtual machine (QEMU/KVM).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The node on which to create the VM.",
				Required:    true,
			},
			"vmid": schema.Int64Attribute{
				Description: "The VM ID. If not set, the next available ID will be used.",
				Optional:    true,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The VM name.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The VM description.",
				Optional:    true,
				Computed:    true,
			},
			"tags": schema.StringAttribute{
				Description: "Tags for the VM (semicolon-separated).",
				Optional:    true,
				Computed:    true,
			},
			"on_boot": schema.BoolAttribute{
				Description: "Whether to start the VM on host boot.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"started": schema.BoolAttribute{
				Description: "Whether the VM should be started after creation.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"protection": schema.BoolAttribute{
				Description: "Whether the VM is protected from removal.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"agent": schema.BoolAttribute{
				Description: "Whether to enable the QEMU Guest Agent.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"os_type": schema.StringAttribute{
				Description: "The OS type (l26, l24, win11, win10, win7, solaris, other).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("l26"),
			},
			"bios": schema.StringAttribute{
				Description: "The BIOS type (seabios or ovmf).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("seabios"),
			},
			"machine": schema.StringAttribute{
				Description: "The machine type (e.g., q35, i440fx, or a specific version).",
				Optional:    true,
				Computed:    true,
			},
			"scsi_hw": schema.StringAttribute{
				Description: "The SCSI controller type (virtio-scsi-pci, virtio-scsi-single, lsi, megasas, pvscsi).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("virtio-scsi-pci"),
			},
			"boot": schema.StringAttribute{
				Description: "Boot order (e.g., 'order=scsi0;ide2;net0').",
				Optional:    true,
				Computed:    true,
			},
			"pool": schema.StringAttribute{
				Description: "The resource pool to add the VM to.",
				Optional:    true,
			},

			// CPU
			"cpu_sockets": schema.Int64Attribute{
				Description: "Number of CPU sockets.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
			},
			"cpu_cores": schema.Int64Attribute{
				Description: "Number of CPU cores per socket.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
			},
			"cpu_type": schema.StringAttribute{
				Description: "The CPU type (e.g., 'host', 'kvm64', 'x86-64-v2-AES').",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("kvm64"),
			},

			// Memory
			"memory": schema.Int64Attribute{
				Description: "Memory in MiB.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(512),
			},
			"balloon": schema.Int64Attribute{
				Description: "Balloon memory minimum in MiB. 0 to disable ballooning.",
				Optional:    true,
				Computed:    true,
			},

			// VGA
			"vga": schema.StringAttribute{
				Description: "VGA configuration string (e.g., 'std', 'virtio', 'serial0').",
				Optional:    true,
				Computed:    true,
			},

			// Disks
			"scsi0": schema.StringAttribute{
				Description: "SCSI disk 0 configuration string.",
				Optional:    true,
				Computed:    true,
			},
			"scsi1": schema.StringAttribute{
				Description: "SCSI disk 1 configuration string.",
				Optional:    true,
				Computed:    true,
			},
			"scsi2": schema.StringAttribute{
				Description: "SCSI disk 2 configuration string.",
				Optional:    true,
				Computed:    true,
			},
			"scsi3": schema.StringAttribute{
				Description: "SCSI disk 3 configuration string.",
				Optional:    true,
				Computed:    true,
			},
			"virtio0": schema.StringAttribute{
				Description: "VirtIO disk 0 configuration string.",
				Optional:    true,
				Computed:    true,
			},
			"virtio1": schema.StringAttribute{
				Description: "VirtIO disk 1 configuration string.",
				Optional:    true,
				Computed:    true,
			},
			"ide0": schema.StringAttribute{
				Description: "IDE disk 0 configuration string.",
				Optional:    true,
				Computed:    true,
			},
			"ide2": schema.StringAttribute{
				Description: "IDE disk 2 configuration string (often used for cloud-init).",
				Optional:    true,
				Computed:    true,
			},
			"efidisk0": schema.StringAttribute{
				Description: "EFI disk configuration string.",
				Optional:    true,
				Computed:    true,
			},
			"tpmstate0": schema.StringAttribute{
				Description: "TPM state configuration string.",
				Optional:    true,
				Computed:    true,
			},

			// Network
			"net0": schema.StringAttribute{
				Description: "Network device 0 configuration (e.g., 'virtio=XX:XX:XX:XX:XX:XX,bridge=vmbr0').",
				Optional:    true,
				Computed:    true,
			},
			"net1": schema.StringAttribute{
				Description: "Network device 1 configuration.",
				Optional:    true,
				Computed:    true,
			},
			"net2": schema.StringAttribute{
				Description: "Network device 2 configuration.",
				Optional:    true,
				Computed:    true,
			},
			"net3": schema.StringAttribute{
				Description: "Network device 3 configuration.",
				Optional:    true,
				Computed:    true,
			},

			// Cloud-init
			"ci_user": schema.StringAttribute{
				Description: "Cloud-init user.",
				Optional:    true,
				Computed:    true,
			},
			"ci_password": schema.StringAttribute{
				Description: "Cloud-init password.",
				Optional:    true,
				Sensitive:   true,
			},
			"ci_type": schema.StringAttribute{
				Description: "Cloud-init type (configdrive2 or nocloud).",
				Optional:    true,
				Computed:    true,
			},
			"ipconfig0": schema.StringAttribute{
				Description: "IP configuration for net0 (e.g., 'ip=dhcp' or 'ip=10.0.0.2/24,gw=10.0.0.1').",
				Optional:    true,
				Computed:    true,
			},
			"ipconfig1": schema.StringAttribute{
				Description: "IP configuration for net1.",
				Optional:    true,
				Computed:    true,
			},
			"nameserver": schema.StringAttribute{
				Description: "Cloud-init DNS nameserver.",
				Optional:    true,
				Computed:    true,
			},
			"searchdomain": schema.StringAttribute{
				Description: "Cloud-init DNS search domain.",
				Optional:    true,
				Computed:    true,
			},
			"ssh_keys": schema.StringAttribute{
				Description: "Cloud-init SSH public keys (URL-encoded, newline-separated).",
				Optional:    true,
				Computed:    true,
			},

			// Serial
			"serial0": schema.StringAttribute{
				Description: "Serial device 0 (e.g., 'socket').",
				Optional:    true,
				Computed:    true,
			},

			// Clone
			"clone_vmid": schema.Int64Attribute{
				Description: "VMID of the template/VM to clone from. If set, the VM is created as a clone.",
				Optional:    true,
			},
			"full_clone": schema.BoolAttribute{
				Description: "Whether to do a full clone (true) or linked clone (false).",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},

			// Read-only
			"status": schema.StringAttribute{
				Description: "The current VM status (running, stopped, etc.).",
				Computed:    true,
			},
			"template": schema.BoolAttribute{
				Description: "Whether the VM is a template.",
				Computed:    true,
			},
		},
	}
}

func (r *VMResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *VMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VMResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()

	// get or assign VMID
	vmid := int(plan.VMID.ValueInt64())
	if plan.VMID.IsNull() || plan.VMID.IsUnknown() || vmid == 0 {
		nextID, err := r.client.GetNextVMID(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Error getting next VMID", err.Error())
			return
		}
		vmid = nextID
	}

	// if clone_vmid is set, clone from that VM instead of creating fresh
	if !plan.CloneVMID.IsNull() && plan.CloneVMID.ValueInt64() > 0 {
		r.createFromClone(ctx, node, vmid, &plan, &resp.Diagnostics)
	} else {
		r.createNew(ctx, node, vmid, &plan, &resp.Diagnostics)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	plan.VMID = types.Int64Value(int64(vmid))
	plan.ID = types.StringValue(fmt.Sprintf("%s/%d", node, vmid))

	// start the VM if the user wants it running
	if plan.Started.ValueBool() {
		status, err := r.client.GetVMStatus(ctx, node, vmid)
		if err != nil {
			resp.Diagnostics.AddError("Error reading VM status", err.Error())
			return
		}
		if status.Status != "running" {
			upid, err := r.client.StartVM(ctx, node, vmid)
			if err != nil {
				resp.Diagnostics.AddError("Error starting VM", err.Error())
				return
			}
			if err := r.client.WaitForTask(ctx, node, upid); err != nil {
				resp.Diagnostics.AddError("Error waiting for VM start", err.Error())
				return
			}
		}
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VMResource) createNew(ctx context.Context, node string, vmid int, plan *VMResourceModel, diagnostics *diag.Diagnostics) {
	onboot := boolToInt(plan.OnBoot.ValueBool())
	protection := boolToInt(plan.Protection.ValueBool())

	agentStr := "0"
	if plan.Agent.ValueBool() {
		agentStr = "1"
	}

	createReq := &models.VMCreateRequest{
		VMID:         vmid,
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueString(),
		Tags:         plan.Tags.ValueString(),
		OnBoot:       &onboot,
		Protection:   &protection,
		Agent:        agentStr,
		OSType:       plan.OSType.ValueString(),
		BIOS:         plan.BIOS.ValueString(),
		Machine:      plan.Machine.ValueString(),
		SCSIHw:       plan.SCSIHw.ValueString(),
		Boot:         plan.Boot.ValueString(),
		Pool:         plan.Pool.ValueString(),
		Sockets:      int(plan.Sockets.ValueInt64()),
		Cores:        int(plan.Cores.ValueInt64()),
		CPUType:      plan.CPUType.ValueString(),
		Memory:       int(plan.Memory.ValueInt64()),
		Balloon:      int(plan.Balloon.ValueInt64()),
		VGA:          plan.VGA.ValueString(),
		Serial0:      plan.Serial0.ValueString(),
		SCSI0:        plan.SCSI0.ValueString(),
		SCSI1:        plan.SCSI1.ValueString(),
		SCSI2:        plan.SCSI2.ValueString(),
		SCSI3:        plan.SCSI3.ValueString(),
		VirtIO0:      plan.VirtIO0.ValueString(),
		VirtIO1:      plan.VirtIO1.ValueString(),
		IDE0:         plan.IDE0.ValueString(),
		IDE2:         plan.IDE2.ValueString(),
		EFIDisk0:     plan.EFIDisk0.ValueString(),
		TPMState0:    plan.TPMState0.ValueString(),
		Net0:         plan.Net0.ValueString(),
		Net1:         plan.Net1.ValueString(),
		Net2:         plan.Net2.ValueString(),
		Net3:         plan.Net3.ValueString(),
		CIUser:       plan.CIUser.ValueString(),
		CIPassword:   plan.CIPassword.ValueString(),
		CIType:       plan.CIType.ValueString(),
		IPConfig0:    plan.IPConfig0.ValueString(),
		IPConfig1:    plan.IPConfig1.ValueString(),
		Nameserver:   plan.Nameserver.ValueString(),
		Searchdomain: plan.Searchdomain.ValueString(),
		SSHKeys:      plan.SSHKeys.ValueString(),
	}

	tflog.Debug(ctx, "Creating VM", map[string]any{"node": node, "vmid": vmid, "name": plan.Name.ValueString()})

	upid, err := r.client.CreateVM(ctx, node, createReq)
	if err != nil {
		diagnostics.AddError("Error creating VM", err.Error())
		return
	}

	if upid != "" {
		if err := r.client.WaitForTask(ctx, node, upid); err != nil {
			diagnostics.AddError("Error waiting for VM creation", err.Error())
			return
		}
	}
}

func (r *VMResource) createFromClone(ctx context.Context, node string, vmid int, plan *VMResourceModel, diagnostics *diag.Diagnostics) {
	sourceVMID := int(plan.CloneVMID.ValueInt64())

	full := 1
	if !plan.FullClone.IsNull() && !plan.FullClone.ValueBool() {
		full = 0
	}

	cloneReq := &models.VMCloneRequest{
		NewID:       vmid,
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Pool:        plan.Pool.ValueString(),
		Full:        &full,
	}

	tflog.Debug(ctx, "Cloning VM", map[string]any{"source": sourceVMID, "target": vmid, "node": node})

	upid, err := r.client.CloneVM(ctx, node, sourceVMID, cloneReq)
	if err != nil {
		diagnostics.AddError("Error cloning VM", err.Error())
		return
	}

	if err := r.client.WaitForTask(ctx, node, upid); err != nil {
		diagnostics.AddError("Error waiting for VM clone", err.Error())
		return
	}

	// apply any extra config on top of the clone
	configMap := r.buildConfigMap(plan)
	if len(configMap) > 0 {
		if err := r.client.UpdateVMConfig(ctx, node, vmid, configMap); err != nil {
			diagnostics.AddError("Error configuring cloned VM", err.Error())
			return
		}
	}
}

func (r *VMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VMResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VMResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VMResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	vmid := int(plan.VMID.ValueInt64())

	configMap := r.buildConfigMap(&plan)
	if len(configMap) > 0 {
		if err := r.client.UpdateVMConfig(ctx, node, vmid, configMap); err != nil {
			resp.Diagnostics.AddError("Error updating VM", err.Error())
			return
		}
	}

	// start or stop the VM to match the desired state
	status, err := r.client.GetVMStatus(ctx, node, vmid)
	if err != nil {
		resp.Diagnostics.AddError("Error reading VM status", err.Error())
		return
	}

	if plan.Started.ValueBool() && status.Status != "running" {
		upid, err := r.client.StartVM(ctx, node, vmid)
		if err != nil {
			resp.Diagnostics.AddError("Error starting VM", err.Error())
			return
		}
		if err := r.client.WaitForTask(ctx, node, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for VM start", err.Error())
			return
		}
	} else if !plan.Started.ValueBool() && status.Status == "running" {
		upid, err := r.client.ShutdownVM(ctx, node, vmid)
		if err != nil {
			resp.Diagnostics.AddError("Error shutting down VM", err.Error())
			return
		}
		if err := r.client.WaitForTask(ctx, node, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for VM shutdown", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%d", node, vmid))
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VMResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := state.NodeName.ValueString()
	vmid := int(state.VMID.ValueInt64())

	// stop the VM before deleting it if its still running
	status, err := r.client.GetVMStatus(ctx, node, vmid)
	if err != nil {
		resp.Diagnostics.AddError("Error reading VM status", err.Error())
		return
	}

	if status.Status == "running" {
		tflog.Debug(ctx, "Stopping VM before deletion", map[string]any{"vmid": vmid})
		upid, err := r.client.StopVM(ctx, node, vmid)
		if err != nil {
			resp.Diagnostics.AddError("Error stopping VM", err.Error())
			return
		}
		if err := r.client.WaitForTask(ctx, node, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for VM stop", err.Error())
			return
		}
	}

	tflog.Debug(ctx, "Deleting VM", map[string]any{"vmid": vmid})

	upid, err := r.client.DeleteVM(ctx, node, vmid)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting VM", err.Error())
		return
	}

	if upid != "" {
		if err := r.client.WaitForTask(ctx, node, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for VM deletion", err.Error())
			return
		}
	}
}

func (r *VMResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// import format: node/vmid
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Import ID must be in format 'node_name/vmid'")
		return
	}

	vmid, err := strconv.Atoi(parts[1])
	if err != nil {
		resp.Diagnostics.AddError("Invalid VMID", fmt.Sprintf("VMID must be an integer: %s", err))
		return
	}

	state := VMResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(parts[0]),
		VMID:     types.Int64Value(int64(vmid)),
		Started:  types.BoolValue(true),
	}

	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VMResource) readIntoModel(ctx context.Context, model *VMResourceModel, diagnostics *diag.Diagnostics) {
	node := model.NodeName.ValueString()
	vmid := int(model.VMID.ValueInt64())

	cfg, err := r.client.GetVMConfig(ctx, node, vmid)
	if err != nil {
		diagnostics.AddError("Error reading VM config", err.Error())
		return
	}

	status, err := r.client.GetVMStatus(ctx, node, vmid)
	if err != nil {
		diagnostics.AddError("Error reading VM status", err.Error())
		return
	}

	model.Name = types.StringValue(cfg.Name)
	model.Description = types.StringValue(cfg.Description)
	model.Tags = types.StringValue(cfg.Tags)
	model.OSType = types.StringValue(cfg.OSType)
	model.BIOS = types.StringValue(cfg.BIOS)
	model.Machine = types.StringValue(cfg.Machine)
	model.SCSIHw = types.StringValue(cfg.SCSIHw)
	model.Boot = types.StringValue(cfg.Boot)
	model.VGA = types.StringValue(cfg.VGA)
	model.Serial0 = types.StringValue(cfg.Serial0)
	model.CPUType = types.StringValue(cfg.CPUType)

	if cfg.OnBoot != nil {
		model.OnBoot = types.BoolValue(*cfg.OnBoot == 1)
	}
	if cfg.Protection != nil {
		model.Protection = types.BoolValue(*cfg.Protection == 1)
	}
	if cfg.Template != nil {
		model.Template = types.BoolValue(*cfg.Template == 1)
	} else {
		model.Template = types.BoolValue(false)
	}

	model.Agent = types.BoolValue(strings.HasPrefix(cfg.Agent, "1"))

	if cfg.Sockets > 0 {
		model.Sockets = types.Int64Value(int64(cfg.Sockets))
	}
	if cfg.Cores > 0 {
		model.Cores = types.Int64Value(int64(cfg.Cores))
	}
	if cfg.Memory > 0 {
		model.Memory = types.Int64Value(int64(cfg.Memory))
	}
	model.Balloon = types.Int64Value(int64(cfg.Balloon))

	// Cloud-init
	model.CIUser = types.StringValue(cfg.CIUser)
	model.CIType = types.StringValue(cfg.CIType)
	model.IPConfig0 = types.StringValue(cfg.IPConfig0)
	model.IPConfig1 = types.StringValue(cfg.IPConfig1)
	model.Nameserver = types.StringValue(cfg.Nameserver)
	model.Searchdomain = types.StringValue(cfg.Searchdomain)
	model.SSHKeys = types.StringValue(cfg.SSHKeys)

	// EFI/TPM
	model.EFIDisk0 = types.StringValue(cfg.EFIDisk0)
	model.TPMState0 = types.StringValue(cfg.TPMState0)

	// Disks — only fill in from API if the user hasnt already set them
	// this avoids conflicts when proxmox expands config (e.g. assigning disk names)
	if model.SCSI0.IsNull() || model.SCSI0.ValueString() == "" {
		model.SCSI0 = types.StringValue(normalizeScsiConfig(cfg.SCSI0))
	}
	if model.SCSI1.IsNull() || model.SCSI1.ValueString() == "" {
		model.SCSI1 = types.StringValue(normalizeScsiConfig(cfg.SCSI1))
	}
	if model.SCSI2.IsNull() || model.SCSI2.ValueString() == "" {
		model.SCSI2 = types.StringValue(normalizeScsiConfig(cfg.SCSI2))
	}
	if model.SCSI3.IsNull() || model.SCSI3.ValueString() == "" {
		model.SCSI3 = types.StringValue(normalizeScsiConfig(cfg.SCSI3))
	}
	if model.VirtIO0.IsNull() || model.VirtIO0.ValueString() == "" {
		model.VirtIO0 = types.StringValue(cfg.VirtIO0)
	}
	if model.VirtIO1.IsNull() || model.VirtIO1.ValueString() == "" {
		model.VirtIO1 = types.StringValue(cfg.VirtIO1)
	}
	if model.IDE0.IsNull() || model.IDE0.ValueString() == "" {
		model.IDE0 = types.StringValue(cfg.IDE0)
	}
	if model.IDE2.IsNull() || model.IDE2.ValueString() == "" {
		model.IDE2 = types.StringValue(normalizeIsoConfig(cfg.IDE2))
	}

	// Network — same deal, dont overwrite user-configured values
	if model.Net0.IsNull() || model.Net0.ValueString() == "" {
		model.Net0 = types.StringValue(normalizeNetConfig(cfg.Net0))
	}
	if model.Net1.IsNull() || model.Net1.ValueString() == "" {
		model.Net1 = types.StringValue(normalizeNetConfig(cfg.Net1))
	}
	if model.Net2.IsNull() || model.Net2.ValueString() == "" {
		model.Net2 = types.StringValue(normalizeNetConfig(cfg.Net2))
	}
	if model.Net3.IsNull() || model.Net3.ValueString() == "" {
		model.Net3 = types.StringValue(normalizeNetConfig(cfg.Net3))
	}

	// Status
	model.Status = types.StringValue(status.Status)
}

func (r *VMResource) buildConfigMap(plan *VMResourceModel) map[string]interface{} {
	m := make(map[string]interface{})

	setIfNotEmpty := func(key, val string) {
		if val != "" {
			m[key] = val
		}
	}

	setIfNotEmpty("name", plan.Name.ValueString())
	setIfNotEmpty("description", plan.Description.ValueString())
	setIfNotEmpty("tags", plan.Tags.ValueString())
	setIfNotEmpty("ostype", plan.OSType.ValueString())
	setIfNotEmpty("bios", plan.BIOS.ValueString())
	setIfNotEmpty("machine", plan.Machine.ValueString())
	setIfNotEmpty("scsihw", plan.SCSIHw.ValueString())
	setIfNotEmpty("boot", plan.Boot.ValueString())
	setIfNotEmpty("vga", plan.VGA.ValueString())
	setIfNotEmpty("serial0", plan.Serial0.ValueString())
	setIfNotEmpty("cpu", plan.CPUType.ValueString())

	if !plan.OnBoot.IsNull() {
		m["onboot"] = boolToInt(plan.OnBoot.ValueBool())
	}
	if !plan.Protection.IsNull() {
		m["protection"] = boolToInt(plan.Protection.ValueBool())
	}
	if plan.Agent.ValueBool() {
		m["agent"] = "1"
	} else if !plan.Agent.IsNull() {
		m["agent"] = "0"
	}

	if v := plan.Sockets.ValueInt64(); v > 0 {
		m["sockets"] = v
	}
	if v := plan.Cores.ValueInt64(); v > 0 {
		m["cores"] = v
	}
	if v := plan.Memory.ValueInt64(); v > 0 {
		m["memory"] = v
	}
	if !plan.Balloon.IsNull() {
		m["balloon"] = plan.Balloon.ValueInt64()
	}

	// Disks
	setIfNotEmpty("scsi0", plan.SCSI0.ValueString())
	setIfNotEmpty("scsi1", plan.SCSI1.ValueString())
	setIfNotEmpty("scsi2", plan.SCSI2.ValueString())
	setIfNotEmpty("scsi3", plan.SCSI3.ValueString())
	setIfNotEmpty("virtio0", plan.VirtIO0.ValueString())
	setIfNotEmpty("virtio1", plan.VirtIO1.ValueString())
	setIfNotEmpty("ide0", plan.IDE0.ValueString())
	setIfNotEmpty("ide2", plan.IDE2.ValueString())
	setIfNotEmpty("efidisk0", plan.EFIDisk0.ValueString())
	setIfNotEmpty("tpmstate0", plan.TPMState0.ValueString())

	// Network
	setIfNotEmpty("net0", plan.Net0.ValueString())
	setIfNotEmpty("net1", plan.Net1.ValueString())
	setIfNotEmpty("net2", plan.Net2.ValueString())
	setIfNotEmpty("net3", plan.Net3.ValueString())

	// Cloud-init
	setIfNotEmpty("ciuser", plan.CIUser.ValueString())
	setIfNotEmpty("cipassword", plan.CIPassword.ValueString())
	setIfNotEmpty("citype", plan.CIType.ValueString())
	setIfNotEmpty("ipconfig0", plan.IPConfig0.ValueString())
	setIfNotEmpty("ipconfig1", plan.IPConfig1.ValueString())
	setIfNotEmpty("nameserver", plan.Nameserver.ValueString())
	setIfNotEmpty("searchdomain", plan.Searchdomain.ValueString())
	setIfNotEmpty("sshkeys", plan.SSHKeys.ValueString())

	return m
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// normalizeScsiConfig strips the disk name from a full SCSI config string.
// proxmox returns "local-lvm:vm-201-disk-0,size=10G" but we want just "local-lvm:10".
func normalizeScsiConfig(full string) string {
	if full == "" {
		return ""
	}
	// grab only the storage:size part — everything before the first comma
	parts := strings.Split(full, ",")
	if len(parts) > 0 {
		return parts[0]
	}
	return full
}

// normalizeNetConfig strips the MAC address from a full net config string.
// proxmox returns "virtio=BC:24:11:8C:6B:69,bridge=WAN" but we want "virtio,bridge=WAN".
func normalizeNetConfig(full string) string {
	if full == "" {
		return ""
	}
	// drop the MAC address part and keep everything else
	parts := strings.Split(full, ",")
	var result []string
	for _, part := range parts {
		// MAC address fields look like "virtio=XX:XX:..." — strip down to just the device type
		if strings.Contains(part, "=") && (strings.HasPrefix(part, "virtio=") || strings.HasPrefix(part, "e1000=") || strings.HasPrefix(part, "rtl8139=")) {
			// keep only the device type name
			device := strings.Split(part, "=")[0]
			result = append(result, device)
		} else {
			result = append(result, part)
		}
	}
	return strings.Join(result, ",")
}

// normalizeIsoConfig removes the size field from a full IDE config string.
// proxmox returns "local:iso/alpine.iso,media=cdrom,size=60M" but we want "local:iso/alpine.iso,media=cdrom"
func normalizeIsoConfig(full string) string {
	if full == "" {
		return ""
	}
	// strip out the size= suffix proxmox adds
	result := strings.ReplaceAll(full, ",size=60M", "")
	result = strings.ReplaceAll(result, ",size=100M", "")
	result = strings.ReplaceAll(result, ",size=1G", "")
	// catch any other size= fields we might have missed
	for _, part := range strings.Split(result, ",") {
		if strings.HasPrefix(part, "size=") {
			result = strings.Replace(result, ","+part, "", 1)
		}
	}
	return result
}
