package storage_lvmthin

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

const storageType = "lvmthin"

var _ resource.Resource = &StorageLVMThinResource{}
var _ resource.ResourceWithConfigure = &StorageLVMThinResource{}
var _ resource.ResourceWithImportState = &StorageLVMThinResource{}

type StorageLVMThinResource struct {
	client *client.Client
}

type StorageLVMThinModel struct {
	ID       types.String `tfsdk:"id"`
	Storage  types.String `tfsdk:"storage"`
	VGName   types.String `tfsdk:"vgname"`
	Thinpool types.String `tfsdk:"thinpool"`
	Content  types.String `tfsdk:"content"`
	Nodes    types.String `tfsdk:"nodes"`
	Disable  types.Bool   `tfsdk:"disable"`
	Shared   types.Bool   `tfsdk:"shared"`
}

func NewResource() resource.Resource {
	return &StorageLVMThinResource{}
}

func (r *StorageLVMThinResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_storage_lvmthin"
}

func (r *StorageLVMThinResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE LVM-thin storage definition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"storage": schema.StringAttribute{
				Description: "The storage identifier/name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vgname": schema.StringAttribute{
				Description: "The LVM volume group that contains the thin pool.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"thinpool": schema.StringAttribute{
				Description: "The LVM thin pool name within the volume group.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Description: "Comma-separated list of content types (images, rootdir).",
				Optional:    true,
				Computed:    true,
			},
			"nodes": schema.StringAttribute{
				Description: "Comma-separated list of cluster nodes where this storage is accessible. Empty means all nodes.",
				Optional:    true,
				Computed:    true,
			},
			"disable": schema.BoolAttribute{
				Description: "Whether to disable this storage.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"shared": schema.BoolAttribute{
				Description: "Whether the storage is shared across nodes.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *StorageLVMThinResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StorageLVMThinResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan StorageLVMThinModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// the thinpool name lives in the Pool field of the API model
	createReq := &models.StorageCreateRequest{
		Storage: plan.Storage.ValueString(),
		Type:    storageType,
		VGName:  plan.VGName.ValueString(),
		Pool:    plan.Thinpool.ValueString(),
		Content: plan.Content.ValueString(),
		Nodes:   plan.Nodes.ValueString(),
		Disable: boolToIntPtr(plan.Disable.ValueBool()),
	}

	if !plan.Shared.IsNull() && !plan.Shared.IsUnknown() {
		createReq.Shared = boolToIntPtr(plan.Shared.ValueBool())
	}

	tflog.Debug(ctx, "Creating LVM-thin storage", map[string]any{"storage": createReq.Storage})

	if err := r.client.CreateStorage(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating LVM-thin storage", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Storage.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageLVMThinResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StorageLVMThinModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StorageLVMThinResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan StorageLVMThinModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := plan.Content.ValueString()
	nodes := plan.Nodes.ValueString()
	disableInt := boolToIntPtr(plan.Disable.ValueBool())

	updateReq := &models.StorageUpdateRequest{
		Content: &content,
		Nodes:   &nodes,
		Disable: disableInt,
	}

	if !plan.Shared.IsNull() && !plan.Shared.IsUnknown() {
		updateReq.Shared = boolToIntPtr(plan.Shared.ValueBool())
	}

	if err := r.client.UpdateStorage(ctx, plan.Storage.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating LVM-thin storage", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Storage.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageLVMThinResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state StorageLVMThinModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteStorage(ctx, state.Storage.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting LVM-thin storage", err.Error())
	}
}

func (r *StorageLVMThinResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := StorageLVMThinModel{
		ID:      types.StringValue(req.ID),
		Storage: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StorageLVMThinResource) readIntoModel(ctx context.Context, model *StorageLVMThinModel, diagnostics *diag.Diagnostics) {
	cfg, err := r.client.GetStorageConfig(ctx, model.Storage.ValueString())
	if err != nil {
		if isNotFound(err) {
			diagnostics.AddWarning("Storage not found", fmt.Sprintf("LVM-thin storage %q no longer exists, removing from state.", model.Storage.ValueString()))
			return
		}
		diagnostics.AddError("Error reading LVM-thin storage", err.Error())
		return
	}

	if cfg.Type != storageType {
		diagnostics.AddError("Storage type mismatch",
			fmt.Sprintf("Expected storage type %q but got %q for storage %q.", storageType, cfg.Type, model.Storage.ValueString()))
		return
	}

	model.VGName = types.StringValue(cfg.VGName)
	// API stores thinpool name in the Pool field
	model.Thinpool = types.StringValue(cfg.Pool)
	model.Content = types.StringValue(cfg.Content)
	model.Nodes = types.StringValue(cfg.Nodes)

	if cfg.Disable != nil {
		model.Disable = types.BoolValue(*cfg.Disable == 1)
	} else {
		model.Disable = types.BoolValue(false)
	}

	if cfg.Shared != nil {
		model.Shared = types.BoolValue(*cfg.Shared == 1)
	} else {
		model.Shared = types.BoolValue(false)
	}
}

func boolToIntPtr(b bool) *int {
	v := 0
	if b {
		v = 1
	}
	return &v
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}
