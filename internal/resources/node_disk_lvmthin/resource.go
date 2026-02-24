package node_disk_lvmthin

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

var _ resource.Resource = &NodeDiskLVMThinResource{}
var _ resource.ResourceWithConfigure = &NodeDiskLVMThinResource{}
var _ resource.ResourceWithImportState = &NodeDiskLVMThinResource{}

type NodeDiskLVMThinResource struct {
	client *client.Client
}

type NodeDiskLVMThinResourceModel struct {
	ID          types.String `tfsdk:"id"`
	NodeName    types.String `tfsdk:"node_name"`
	Device      types.String `tfsdk:"device"`
	Name        types.String `tfsdk:"name"`
	VolumeGroup types.String `tfsdk:"volume_group"`
	AddStorage  types.Bool   `tfsdk:"add_storage"`
}

func NewResource() resource.Resource {
	return &NodeDiskLVMThinResource{}
}

func (r *NodeDiskLVMThinResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_disk_lvmthin"
}

func (r *NodeDiskLVMThinResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an LVM-thin pool on a Proxmox VE node disk.",
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
				Description: "The LVM-thin pool name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"volume_group": schema.StringAttribute{
				Description: "The LVM volume group that will contain the thin pool.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"add_storage": schema.BoolAttribute{
				Description: "Whether to automatically add the created LVM-thin pool as a Proxmox storage.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *NodeDiskLVMThinResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NodeDiskLVMThinResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NodeDiskLVMThinResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	createReq := &models.NodeDiskLVMThinCreateRequest{
		Dev:         plan.Device.ValueString(),
		Name:        plan.Name.ValueString(),
		VolumeGroup: plan.VolumeGroup.ValueString(),
		AddStorage:  plan.AddStorage.ValueBool(),
	}

	tflog.Debug(ctx, "Creating node disk LVM-thin pool", map[string]any{
		"node":         node,
		"name":         createReq.Name,
		"dev":          createReq.Dev,
		"volume_group": createReq.VolumeGroup,
	})

	upid, err := r.client.CreateNodeDiskLVMThin(ctx, node, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating node disk LVM-thin pool", err.Error())
		return
	}

	if err := r.client.WaitForUPID(ctx, upid); err != nil {
		resp.Diagnostics.AddError("Error waiting for node disk LVM-thin pool creation", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s:%s", node, plan.VolumeGroup.ValueString(), plan.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NodeDiskLVMThinResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// no single-item GET for LVM-thin pools — keep state as-is
}

func (r *NodeDiskLVMThinResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// all attributes are ForceNew so Update is never called
}

func (r *NodeDiskLVMThinResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NodeDiskLVMThinResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := state.NodeName.ValueString()
	name := state.Name.ValueString()
	volumeGroup := state.VolumeGroup.ValueString()

	tflog.Debug(ctx, "Deleting node disk LVM-thin pool", map[string]any{
		"node":         node,
		"name":         name,
		"volume_group": volumeGroup,
	})

	upid, err := r.client.DeleteNodeDiskLVMThin(ctx, node, name, volumeGroup)
	if err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting node disk LVM-thin pool", err.Error())
		return
	}

	if err := r.client.WaitForUPID(ctx, upid); err != nil {
		resp.Diagnostics.AddError("Error waiting for node disk LVM-thin pool deletion", err.Error())
	}
}

func (r *NodeDiskLVMThinResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// import format: node_name:volume_group:name
	parts := strings.SplitN(req.ID, ":", 3)
	if len(parts) != 3 {
		resp.Diagnostics.AddError("Invalid import ID",
			"Format: <node_name>:<volume_group>:<name> (e.g. 'pve:pve-vg:data')")
		return
	}

	state := NodeDiskLVMThinResourceModel{
		ID:          types.StringValue(req.ID),
		NodeName:    types.StringValue(parts[0]),
		VolumeGroup: types.StringValue(parts[1]),
		Name:        types.StringValue(parts[2]),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}
