package storage_zfspool

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

const storageType = "zfspool"

var _ resource.Resource = &StorageZFSPoolResource{}
var _ resource.ResourceWithConfigure = &StorageZFSPoolResource{}
var _ resource.ResourceWithImportState = &StorageZFSPoolResource{}

type StorageZFSPoolResource struct {
	client *client.Client
}

type StorageZFSPoolModel struct {
	ID        types.String `tfsdk:"id"`
	Storage   types.String `tfsdk:"storage"`
	Pool      types.String `tfsdk:"pool"`
	Blocksize types.String `tfsdk:"blocksize"`
	Content   types.String `tfsdk:"content"`
	Nodes     types.String `tfsdk:"nodes"`
	Disable   types.Bool   `tfsdk:"disable"`
	Shared    types.Bool   `tfsdk:"shared"`
}

func NewResource() resource.Resource {
	return &StorageZFSPoolResource{}
}

func (r *StorageZFSPoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_storage_zfspool"
}

func (r *StorageZFSPoolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE ZFS pool storage definition.",
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
			"pool": schema.StringAttribute{
				Description: "The ZFS pool or dataset name (e.g. 'rpool' or 'data/vm-store').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"blocksize": schema.StringAttribute{
				Description: "The ZFS block size (e.g. '8k', '16k', '32k'). Stored as preallocation hint.",
				Optional:    true,
				Computed:    true,
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

func (r *StorageZFSPoolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StorageZFSPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan StorageZFSPoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &models.StorageCreateRequest{
		Storage:       plan.Storage.ValueString(),
		Type:          storageType,
		Pool:          plan.Pool.ValueString(),
		Preallocation: plan.Blocksize.ValueString(),
		Content:       plan.Content.ValueString(),
		Nodes:         plan.Nodes.ValueString(),
		Disable:       boolToIntPtr(plan.Disable.ValueBool()),
	}

	if !plan.Shared.IsNull() && !plan.Shared.IsUnknown() {
		createReq.Shared = boolToIntPtr(plan.Shared.ValueBool())
	}

	tflog.Debug(ctx, "Creating ZFS pool storage", map[string]any{"storage": createReq.Storage})

	if err := r.client.CreateStorage(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating ZFS pool storage", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Storage.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageZFSPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StorageZFSPoolModel
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

func (r *StorageZFSPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan StorageZFSPoolModel
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
		resp.Diagnostics.AddError("Error updating ZFS pool storage", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Storage.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageZFSPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state StorageZFSPoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteStorage(ctx, state.Storage.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting ZFS pool storage", err.Error())
	}
}

func (r *StorageZFSPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := StorageZFSPoolModel{
		ID:      types.StringValue(req.ID),
		Storage: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StorageZFSPoolResource) readIntoModel(ctx context.Context, model *StorageZFSPoolModel, diagnostics *diag.Diagnostics) {
	cfg, err := r.client.GetStorageConfig(ctx, model.Storage.ValueString())
	if err != nil {
		if isNotFound(err) {
			diagnostics.AddWarning("Storage not found", fmt.Sprintf("ZFS pool storage %q no longer exists, removing from state.", model.Storage.ValueString()))
			return
		}
		diagnostics.AddError("Error reading ZFS pool storage", err.Error())
		return
	}

	if cfg.Type != storageType {
		diagnostics.AddError("Storage type mismatch",
			fmt.Sprintf("Expected storage type %q but got %q for storage %q.", storageType, cfg.Type, model.Storage.ValueString()))
		return
	}

	model.Pool = types.StringValue(cfg.Pool)
	model.Blocksize = types.StringValue(cfg.Preallocation)
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
