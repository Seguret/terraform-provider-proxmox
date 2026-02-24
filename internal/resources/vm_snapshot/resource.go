package vm_snapshot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &VMSnapshotResource{}
var _ resource.ResourceWithConfigure = &VMSnapshotResource{}
var _ resource.ResourceWithImportState = &VMSnapshotResource{}

type VMSnapshotResource struct {
	client *client.Client
}

type VMSnapshotResourceModel struct {
	ID          types.String `tfsdk:"id"`
	NodeName    types.String `tfsdk:"node_name"`
	VMID        types.Int64  `tfsdk:"vmid"`
	SnapName    types.String `tfsdk:"snap_name"`
	Description types.String `tfsdk:"description"`
	VMState     types.Bool   `tfsdk:"vm_state"`
	Snaptime    types.Int64  `tfsdk:"snaptime"`
}

func NewResource() resource.Resource {
	return &VMSnapshotResource{}
}

func (r *VMSnapshotResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_vm_snapshot"
}

func (r *VMSnapshotResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE virtual machine snapshot.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The node where the VM resides.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vmid": schema.Int64Attribute{
				Description: "The VM ID.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"snap_name": schema.StringAttribute{
				Description: "The snapshot name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The snapshot description.",
				Optional:    true,
				Computed:    true,
			},
			"vm_state": schema.BoolAttribute{
				Description: "Whether to include the VM memory state in the snapshot.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"snaptime": schema.Int64Attribute{
				Description: "The snapshot creation time (Unix timestamp).",
				Computed:    true,
			},
		},
	}
}

func (r *VMSnapshotResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VMSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VMSnapshotResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	vmid := int(plan.VMID.ValueInt64())

	vmstate := 0
	if plan.VMState.ValueBool() {
		vmstate = 1
	}

	createReq := &models.VMSnapshotCreateRequest{
		Snapname:    plan.SnapName.ValueString(),
		Description: plan.Description.ValueString(),
		VMSTATE:     &vmstate,
	}

	tflog.Debug(ctx, "Creating VM snapshot", map[string]any{
		"node": node, "vmid": vmid, "snapname": plan.SnapName.ValueString(),
	})

	upid, err := r.client.CreateVMSnapshot(ctx, node, vmid, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VM snapshot", err.Error())
		return
	}

	if upid != "" {
		if err := r.client.WaitForTask(ctx, node, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for snapshot creation", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%d/%s", node, vmid, plan.SnapName.ValueString()))

	snap, err := r.client.GetVMSnapshot(ctx, node, vmid, plan.SnapName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading VM snapshot", err.Error())
		return
	}
	plan.Description = types.StringValue(snap.Description)
	plan.Snaptime = types.Int64Value(snap.Snaptime)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VMSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VMSnapshotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := state.NodeName.ValueString()
	vmid := int(state.VMID.ValueInt64())

	snap, err := r.client.GetVMSnapshot(ctx, node, vmid, state.SnapName.ValueString())
	if err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading VM snapshot", err.Error())
		return
	}

	state.Description = types.StringValue(snap.Description)
	state.Snaptime = types.Int64Value(snap.Snaptime)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VMSnapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VMSnapshotResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	vmid := int(plan.VMID.ValueInt64())

	updateReq := &models.VMSnapshotUpdateRequest{
		Description: plan.Description.ValueString(),
	}

	if err := r.client.UpdateVMSnapshot(ctx, node, vmid, plan.SnapName.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating VM snapshot", err.Error())
		return
	}

	snap, err := r.client.GetVMSnapshot(ctx, node, vmid, plan.SnapName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading VM snapshot", err.Error())
		return
	}
	plan.Description = types.StringValue(snap.Description)
	plan.Snaptime = types.Int64Value(snap.Snaptime)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VMSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VMSnapshotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := state.NodeName.ValueString()
	vmid := int(state.VMID.ValueInt64())

	upid, err := r.client.DeleteVMSnapshot(ctx, node, vmid, state.SnapName.ValueString())
	if err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting VM snapshot", err.Error())
		return
	}

	if upid != "" {
		if err := r.client.WaitForTask(ctx, node, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for snapshot deletion", err.Error())
		}
	}
}

func (r *VMSnapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// format: node/vmid/snapname
	parts := strings.SplitN(req.ID, "/", 3)
	if len(parts) != 3 {
		resp.Diagnostics.AddError("Invalid import ID", "Format: node_name/vmid/snap_name")
		return
	}
	vmid, err := strconv.Atoi(parts[1])
	if err != nil {
		resp.Diagnostics.AddError("Invalid VMID", err.Error())
		return
	}
	state := VMSnapshotResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(parts[0]),
		VMID:     types.Int64Value(int64(vmid)),
		SnapName: types.StringValue(parts[2]),
	}
	snap, err := r.client.GetVMSnapshot(ctx, parts[0], vmid, parts[2])
	if err != nil {
		resp.Diagnostics.AddError("Error reading VM snapshot", err.Error())
		return
	}
	state.Description = types.StringValue(snap.Description)
	state.Snaptime = types.Int64Value(snap.Snaptime)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.StatusCode == 404
	}
	return false
}
