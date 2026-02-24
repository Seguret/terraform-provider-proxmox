package node_disk_lvm

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &NodeDiskLVMResource{}
var _ resource.ResourceWithConfigure = &NodeDiskLVMResource{}
var _ resource.ResourceWithImportState = &NodeDiskLVMResource{}

type NodeDiskLVMResource struct {
	client *client.Client
}

type NodeDiskLVMResourceModel struct {
	ID         types.String `tfsdk:"id"`
	NodeName   types.String `tfsdk:"node_name"`
	Device     types.String `tfsdk:"device"`
	Name       types.String `tfsdk:"name"`
	AddStorage types.Bool   `tfsdk:"add_storage"`
}

func NewResource() resource.Resource {
	return &NodeDiskLVMResource{}
}

func (r *NodeDiskLVMResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_disk_lvm"
}

func (r *NodeDiskLVMResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an LVM volume group on a Proxmox VE node disk.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The node name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"device": schema.StringAttribute{
				Description: "The device path (e.g. /dev/sdb).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The LVM volume group name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"add_storage": schema.BoolAttribute{
				Description: "Whether to automatically add the created LVM VG as a Proxmox storage.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *NodeDiskLVMResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NodeDiskLVMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NodeDiskLVMResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	createReq := &models.NodeDiskLVMCreateRequest{
		Dev:        plan.Device.ValueString(),
		Name:       plan.Name.ValueString(),
		AddStorage: plan.AddStorage.ValueBool(),
	}

	tflog.Debug(ctx, "Creating node disk LVM", map[string]any{
		"node": node,
		"name": createReq.Name,
		"dev":  createReq.Dev,
	})

	upid, err := r.client.CreateNodeDiskLVM(ctx, node, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating node disk LVM", err.Error())
		return
	}

	if err := r.client.WaitForUPID(ctx, upid); err != nil {
		resp.Diagnostics.AddError("Error waiting for node disk LVM creation", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", node, plan.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NodeDiskLVMResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// no single-item GET for LVM VGs — just keep state as-is
}

func (r *NodeDiskLVMResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// all attributes are ForceNew so Update is never called
}

func (r *NodeDiskLVMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NodeDiskLVMResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := state.NodeName.ValueString()
	name := state.Name.ValueString()

	tflog.Debug(ctx, "Deleting node disk LVM", map[string]any{"node": node, "name": name})

	upid, err := r.client.DeleteNodeDiskLVM(ctx, node, name)
	if err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting node disk LVM", err.Error())
		return
	}

	if err := r.client.WaitForUPID(ctx, upid); err != nil {
		resp.Diagnostics.AddError("Error waiting for node disk LVM deletion", err.Error())
	}
}

func (r *NodeDiskLVMResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// import format: node_name:name
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID",
			"Format: <node_name>:<name> (e.g. 'pve:pve-vg')")
		return
	}

	state := NodeDiskLVMResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(parts[0]),
		Name:     types.StringValue(parts[1]),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}
