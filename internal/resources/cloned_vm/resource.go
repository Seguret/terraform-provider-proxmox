package cloned_vm

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

var _ resource.Resource = &ClonedVMResource{}
var _ resource.ResourceWithConfigure = &ClonedVMResource{}
var _ resource.ResourceWithImportState = &ClonedVMResource{}

type ClonedVMResource struct {
	client *client.Client
}

type ClonedVMResourceModel struct {
	ID            types.String `tfsdk:"id"`
	NodeName      types.String `tfsdk:"node_name"`
	SourceNode    types.String `tfsdk:"source_node"`
	SourceVMID    types.Int64  `tfsdk:"source_vmid"`
	VMID          types.Int64  `tfsdk:"vm_id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	FullClone     types.Bool   `tfsdk:"full_clone"`
	TargetNode    types.String `tfsdk:"target_node"`
	TargetStorage types.String `tfsdk:"target_storage"`
	Tags          types.String `tfsdk:"tags"`
	// Read-only
	Status types.String `tfsdk:"status"`
}

func NewResource() resource.Resource {
	return &ClonedVMResource{}
}

func (r *ClonedVMResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_cloned_vm"
}

func (r *ClonedVMResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages a clone of an existing Proxmox VE virtual machine.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The resource identifier in the form '{node_name}/{vmid}'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the Proxmox VE node that will host the cloned VM.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_node": schema.StringAttribute{
				Description: "The Proxmox VE node where the source VM resides.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_vmid": schema.Int64Attribute{
				Description: "The VMID of the VM to clone.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"vm_id": schema.Int64Attribute{
				Description: "The VMID for the cloned VM. If omitted, Proxmox auto-assigns the next available ID.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the cloned VM.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for the cloned VM.",
				Optional:    true,
			},
			"full_clone": schema.BoolAttribute{
				Description: "Whether to perform a full clone (true) or a linked clone (false). Defaults to true.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					// full_clone cant be changed once the VM is created
				},
			},
			"target_node": schema.StringAttribute{
				Description: "The target node for the cloned VM (for cross-node clones). Defaults to node_name.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_storage": schema.StringAttribute{
				Description: "The storage ID where the full clone's disks should be placed.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tags": schema.StringAttribute{
				Description: "Semicolon-separated tags for the cloned VM.",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "The current status of the VM (e.g. 'running', 'stopped').",
				Computed:    true,
			},
		},
	}
}

func (r *ClonedVMResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ClonedVMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ClonedVMResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceNode := plan.SourceNode.ValueString()
	sourceVMID := int(plan.SourceVMID.ValueInt64())

	// if VMID wasnt specified, ask proxmox for the next available one
	newVMID := int(plan.VMID.ValueInt64())
	if newVMID == 0 {
		id, err := r.client.GetNextVMID(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Error fetching next available VMID", err.Error())
			return
		}
		newVMID = id
	}

	fullClone := boolToInt(plan.FullClone.ValueBool())
	cloneReq := &models.VMCloneRequest{
		NewID:       newVMID,
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Full:        &fullClone,
	}

	// use target_node to place the clone on a specific node if provided
	if !plan.TargetNode.IsNull() && !plan.TargetNode.IsUnknown() && plan.TargetNode.ValueString() != "" {
		cloneReq.Target = plan.TargetNode.ValueString()
	} else if plan.NodeName.ValueString() != sourceNode {
		// node_name differs from source_node — use node_name as the destination
		cloneReq.Target = plan.NodeName.ValueString()
	}

	if !plan.TargetStorage.IsNull() && !plan.TargetStorage.IsUnknown() {
		cloneReq.Storage = plan.TargetStorage.ValueString()
	}

	tflog.Debug(ctx, "Cloning Proxmox VE VM", map[string]any{
		"source_node": sourceNode,
		"source_vmid": sourceVMID,
		"new_vmid":    newVMID,
		"full_clone":  fullClone,
	})

	upid, err := r.client.CloneVM(ctx, sourceNode, sourceVMID, cloneReq)
	if err != nil {
		resp.Diagnostics.AddError("Error cloning VM", err.Error())
		return
	}

	if err := r.client.WaitForUPID(ctx, upid); err != nil {
		resp.Diagnostics.AddError("Error waiting for VM clone task", err.Error())
		return
	}

	// figure out which node the clone ended up on
	destNode := plan.NodeName.ValueString()
	if !plan.TargetNode.IsNull() && plan.TargetNode.ValueString() != "" {
		destNode = plan.TargetNode.ValueString()
	}

	plan.VMID = types.Int64Value(int64(newVMID))
	plan.ID = types.StringValue(fmt.Sprintf("%s/%d", destNode, newVMID))

	// push name, description, and tags onto the newly cloned VM
	if err := r.applyMutableConfig(ctx, destNode, newVMID, &plan); err != nil {
		resp.Diagnostics.AddError("Error applying VM config after clone", err.Error())
		return
	}

	r.readIntoModel(ctx, &plan, destNode, newVMID, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ClonedVMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ClonedVMResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node, vmid, err := parseID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid resource ID", err.Error())
		return
	}

	r.readIntoModel(ctx, &state, node, vmid, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ClonedVMResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ClonedVMResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node, vmid, err := parseID(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid resource ID", err.Error())
		return
	}

	if err := r.applyMutableConfig(ctx, node, vmid, &plan); err != nil {
		resp.Diagnostics.AddError("Error updating VM config", err.Error())
		return
	}

	r.readIntoModel(ctx, &plan, node, vmid, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ClonedVMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ClonedVMResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node, vmid, err := parseID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid resource ID", err.Error())
		return
	}

	tflog.Debug(ctx, "Deleting cloned Proxmox VE VM", map[string]any{
		"node": node,
		"vmid": vmid,
	})

	upid, err := r.client.DeleteVM(ctx, node, vmid)
	if err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting VM", err.Error())
		return
	}

	if upid != "" {
		if err := r.client.WaitForUPID(ctx, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for VM delete task", err.Error())
		}
	}
}

func (r *ClonedVMResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// format: {node}/{vmid}
	node, vmid, err := parseID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID",
			"Expected format: <node>/<vmid>. "+err.Error())
		return
	}

	state := ClonedVMResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(node),
		VMID:     types.Int64Value(int64(vmid)),
	}

	r.readIntoModel(ctx, &state, node, vmid, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// readIntoModel fetches the VM config and status and fills in the model.
// Adds a warning (not an error) on 404 so terraform can detect drift gracefully.
func (r *ClonedVMResource) readIntoModel(
	ctx context.Context,
	model *ClonedVMResourceModel,
	node string,
	vmid int,
	diagnostics interface {
		AddError(string, string)
		AddWarning(string, string)
	},
) {
	cfg, err := r.client.GetVMConfig(ctx, node, vmid)
	if err != nil {
		if isNotFound(err) {
			diagnostics.AddWarning("VM not found",
				fmt.Sprintf("VM %d on node '%s' was not found. It may have been deleted outside Terraform.", vmid, node))
			return
		}
		diagnostics.AddError("Error reading VM config",
			fmt.Sprintf("Could not read VM %d on node '%s': %s", vmid, node, err.Error()))
		return
	}

	status, err := r.client.GetVMStatus(ctx, node, vmid)
	if err != nil {
		diagnostics.AddError("Error reading VM status",
			fmt.Sprintf("Could not read status of VM %d on node '%s': %s", vmid, node, err.Error()))
		return
	}

	model.NodeName = types.StringValue(node)
	model.VMID = types.Int64Value(int64(vmid))
	model.ID = types.StringValue(fmt.Sprintf("%s/%d", node, vmid))

	if cfg.Name != "" {
		model.Name = types.StringValue(cfg.Name)
	}
	if cfg.Description != "" {
		model.Description = types.StringValue(cfg.Description)
	}
	if cfg.Tags != "" {
		model.Tags = types.StringValue(cfg.Tags)
	}
	model.Status = types.StringValue(status.Status)
}

// applyMutableConfig pushes name, description, and tags to the VM config API.
func (r *ClonedVMResource) applyMutableConfig(ctx context.Context, node string, vmid int, model *ClonedVMResourceModel) error {
	configMap := map[string]interface{}{}

	if !model.Name.IsNull() && !model.Name.IsUnknown() {
		configMap["name"] = model.Name.ValueString()
	}
	if !model.Description.IsNull() && !model.Description.IsUnknown() {
		configMap["description"] = model.Description.ValueString()
	}
	if !model.Tags.IsNull() && !model.Tags.IsUnknown() {
		configMap["tags"] = model.Tags.ValueString()
	}

	if len(configMap) == 0 {
		return nil
	}

	return r.client.UpdateVMConfig(ctx, node, vmid, configMap)
}

// parseID breaks a "{node}/{vmid}" resource ID into node name and VMID.
func parseID(id string) (string, int, error) {
	idx := strings.LastIndex(id, "/")
	if idx < 0 {
		return "", 0, fmt.Errorf("expected format '<node>/<vmid>', got '%s'", id)
	}
	node := id[:idx]
	vmidStr := id[idx+1:]
	vmid, err := strconv.Atoi(vmidStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid VMID '%s': %w", vmidStr, err)
	}
	return node, vmid, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}
