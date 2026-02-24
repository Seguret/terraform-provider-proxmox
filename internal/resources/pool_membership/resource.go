package pool_membership

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &PoolMembershipResource{}
var _ resource.ResourceWithConfigure = &PoolMembershipResource{}
var _ resource.ResourceWithImportState = &PoolMembershipResource{}

type PoolMembershipResource struct {
	client *client.Client
}

type PoolMembershipResourceModel struct {
	ID         types.String  `tfsdk:"id"`
	PoolID     types.String  `tfsdk:"pool_id"`
	VMs        []types.Int64 `tfsdk:"vms"`
	Containers []types.Int64 `tfsdk:"containers"`
}

func NewResource() resource.Resource {
	return &PoolMembershipResource{}
}

func (r *PoolMembershipResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_pool_membership"
}

func (r *PoolMembershipResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages membership of VMs and containers in a Proxmox VE resource pool.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pool_id": schema.StringAttribute{
				Description: "The pool identifier.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vms": schema.ListAttribute{
				Description: "The VM IDs to include in the pool.",
				Optional:    true,
				ElementType: types.Int64Type,
			},
			"containers": schema.ListAttribute{
				Description: "The container IDs to include in the pool.",
				Optional:    true,
				ElementType: types.Int64Type,
			},
		},
	}
}

func (r *PoolMembershipResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PoolMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PoolMembershipResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolID := plan.PoolID.ValueString()
	tflog.Debug(ctx, "Creating pool membership", map[string]any{"pool_id": poolID})

	// proxmox adds listed members when delete flag is not set
	if err := r.setMembers(ctx, poolID, plan.VMs, plan.Containers, false); err != nil {
		resp.Diagnostics.AddError("Error setting pool membership", err.Error())
		return
	}

	plan.ID = types.StringValue(poolID)

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading pool membership", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PoolMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PoolMembershipResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading pool membership", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PoolMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PoolMembershipResourceModel
	var state PoolMembershipResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolID := plan.PoolID.ValueString()
	tflog.Debug(ctx, "Updating pool membership", map[string]any{"pool_id": poolID})

	// figure out what to remove (in state but not in plan)
	toRemoveVMs := diffInt64Lists(state.VMs, plan.VMs)
	toRemoveCTs := diffInt64Lists(state.Containers, plan.Containers)

	if len(toRemoveVMs) > 0 || len(toRemoveCTs) > 0 {
		if err := r.setMembers(ctx, poolID, toRemoveVMs, toRemoveCTs, true); err != nil {
			resp.Diagnostics.AddError("Error removing pool members", err.Error())
			return
		}
	}

	// figure out what to add (in plan but not yet in state)
	toAddVMs := diffInt64Lists(plan.VMs, state.VMs)
	toAddCTs := diffInt64Lists(plan.Containers, state.Containers)

	if len(toAddVMs) > 0 || len(toAddCTs) > 0 {
		if err := r.setMembers(ctx, poolID, toAddVMs, toAddCTs, false); err != nil {
			resp.Diagnostics.AddError("Error adding pool members", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(poolID)

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading pool membership", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PoolMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PoolMembershipResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolID := state.PoolID.ValueString()
	tflog.Debug(ctx, "Deleting pool membership", map[string]any{"pool_id": poolID})

	// remove all members that are still tracked in state
	if len(state.VMs) > 0 || len(state.Containers) > 0 {
		if err := r.setMembers(ctx, poolID, state.VMs, state.Containers, true); err != nil {
			if !isNotFound(err) {
				resp.Diagnostics.AddError("Error removing pool members", err.Error())
			}
		}
	}
}

func (r *PoolMembershipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := PoolMembershipResourceModel{
		ID:     types.StringValue(req.ID),
		PoolID: types.StringValue(req.ID),
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing pool membership", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// readIntoModel fetches the pool from the API and splits members into VMs and containers.
func (r *PoolMembershipResource) readIntoModel(ctx context.Context, model *PoolMembershipResourceModel) error {
	pool, err := r.client.GetPool(ctx, model.PoolID.ValueString())
	if err != nil {
		return err
	}

	var vms []types.Int64
	var containers []types.Int64

	for _, member := range pool.Members {
		switch member.Type {
		case "qemu":
			vms = append(vms, types.Int64Value(int64(member.VMID)))
		case "lxc":
			containers = append(containers, types.Int64Value(int64(member.VMID)))
		}
	}

	if vms == nil {
		vms = []types.Int64{}
	}
	if containers == nil {
		containers = []types.Int64{}
	}

	model.VMs = vms
	model.Containers = containers
	return nil
}

// setMembers adds or removes pool members depending on the delete flag.
func (r *PoolMembershipResource) setMembers(ctx context.Context, poolID string, vms []types.Int64, containers []types.Int64, delete bool) error {
	vmStr := int64ListToCSV(vms)
	ctStr := int64ListToCSV(containers)

	// VMs and containers share the same "vms" param — proxmox uses VMID space for both
	allIDs := mergeCSV(vmStr, ctStr)

	updateReq := &models.PoolUpdateRequest{}
	if allIDs != "" {
		updateReq.VMs = allIDs
	}
	if delete {
		one := 1
		updateReq.Delete = &one
	}

	return r.client.UpdatePool(ctx, poolID, updateReq)
}

// int64ListToCSV turns a list of int64 values into a comma-separated string.
func int64ListToCSV(ids []types.Int64) string {
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, strconv.FormatInt(id.ValueInt64(), 10))
	}
	return strings.Join(parts, ",")
}

// mergeCSV joins two comma-separated strings, skipping any empty ones.
func mergeCSV(a, b string) string {
	switch {
	case a == "":
		return b
	case b == "":
		return a
	default:
		return a + "," + b
	}
}

// diffInt64Lists returns elements that are in a but not in b.
func diffInt64Lists(a, b []types.Int64) []types.Int64 {
	bSet := make(map[int64]struct{}, len(b))
	for _, v := range b {
		bSet[v.ValueInt64()] = struct{}{}
	}
	var result []types.Int64
	for _, v := range a {
		if _, found := bSet[v.ValueInt64()]; !found {
			result = append(result, v)
		}
	}
	return result
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}
