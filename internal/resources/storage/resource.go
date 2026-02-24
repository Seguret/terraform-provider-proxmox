package storage

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

var _ resource.Resource = &StorageResource{}
var _ resource.ResourceWithConfigure = &StorageResource{}
var _ resource.ResourceWithImportState = &StorageResource{}

type StorageResource struct {
	client *client.Client
}

type StorageResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Storage  types.String `tfsdk:"storage"`
	Type     types.String `tfsdk:"type"`
	Content  types.String `tfsdk:"content"`
	Enabled  types.Bool   `tfsdk:"enabled"`
	Shared   types.Bool   `tfsdk:"shared"`
	Path     types.String `tfsdk:"path"`
	Pool     types.String `tfsdk:"pool"`
	VGName   types.String `tfsdk:"vgname"`
	Nodes    types.String `tfsdk:"nodes"`
	Server   types.String `tfsdk:"server"`
	Export   types.String `tfsdk:"export"`
	Share    types.String `tfsdk:"share"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Domain   types.String `tfsdk:"domain"`
	Datastore types.String `tfsdk:"datastore"`
	Namespace types.String `tfsdk:"namespace"`
	Fingerprint types.String `tfsdk:"fingerprint"`
	PruneBackups types.String `tfsdk:"prune_backups"`
}

func NewResource() resource.Resource {
	return &StorageResource{}
}

func (r *StorageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_storage"
}

func (r *StorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE storage definition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"storage": schema.StringAttribute{
				Description: "The storage identifier.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The storage type (dir, lvm, lvmthin, zfspool, nfs, cifs, glusterfs, iscsi, iscsidirect, rbd, cephfs, pbs).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Description: "Comma-separated list of content types (images, rootdir, vztmpl, iso, backup, snippets).",
				Optional:    true,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the storage is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"shared": schema.BoolAttribute{
				Description: "Whether the storage is shared across nodes.",
				Optional:    true,
				Computed:    true,
			},
			"path": schema.StringAttribute{
				Description: "The filesystem path (for 'dir' type).",
				Optional:    true,
				Computed:    true,
			},
			"pool": schema.StringAttribute{
				Description: "The ZFS/Ceph pool name.",
				Optional:    true,
				Computed:    true,
			},
			"vgname": schema.StringAttribute{
				Description: "The LVM volume group name.",
				Optional:    true,
				Computed:    true,
			},
			"nodes": schema.StringAttribute{
				Description: "Comma-separated list of nodes where this storage is available.",
				Optional:    true,
				Computed:    true,
			},
			"server": schema.StringAttribute{
				Description: "The server address (for NFS, CIFS, iSCSI, PBS).",
				Optional:    true,
				Computed:    true,
			},
			"export": schema.StringAttribute{
				Description: "The NFS export path.",
				Optional:    true,
				Computed:    true,
			},
			"share": schema.StringAttribute{
				Description: "The CIFS share name.",
				Optional:    true,
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username (for CIFS, PBS).",
				Optional:    true,
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password (for CIFS, PBS). Write-only.",
				Optional:    true,
				Sensitive:   true,
			},
			"domain": schema.StringAttribute{
				Description: "The domain (for CIFS).",
				Optional:    true,
				Computed:    true,
			},
			"datastore": schema.StringAttribute{
				Description: "The PBS datastore name.",
				Optional:    true,
				Computed:    true,
			},
			"namespace": schema.StringAttribute{
				Description: "The PBS namespace.",
				Optional:    true,
				Computed:    true,
			},
			"fingerprint": schema.StringAttribute{
				Description: "The PBS server fingerprint.",
				Optional:    true,
				Computed:    true,
			},
			"prune_backups": schema.StringAttribute{
				Description: "Backup retention policy.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *StorageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *StorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan StorageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	disableInt := 0
	if !plan.Enabled.IsNull() && !plan.Enabled.ValueBool() {
		disableInt = 1
	}

	createReq := &models.StorageCreateRequest{
		Storage:      plan.Storage.ValueString(),
		Type:         plan.Type.ValueString(),
		Content:      plan.Content.ValueString(),
		Disable:      &disableInt,
		Path:         plan.Path.ValueString(),
		Pool:         plan.Pool.ValueString(),
		VGName:       plan.VGName.ValueString(),
		Nodes:        plan.Nodes.ValueString(),
		Server:       plan.Server.ValueString(),
		Export:       plan.Export.ValueString(),
		Share:        plan.Share.ValueString(),
		Username:     plan.Username.ValueString(),
		Password:     plan.Password.ValueString(),
		Domain:       plan.Domain.ValueString(),
		Datastore:    plan.Datastore.ValueString(),
		Namespace:    plan.Namespace.ValueString(),
		Fingerprint:  plan.Fingerprint.ValueString(),
		PruneBackups: plan.PruneBackups.ValueString(),
	}

	if !plan.Shared.IsNull() {
		sharedInt := 0
		if plan.Shared.ValueBool() {
			sharedInt = 1
		}
		createReq.Shared = &sharedInt
	}

	tflog.Debug(ctx, "Creating storage", map[string]any{"storage": createReq.Storage, "type": createReq.Type})

	if err := r.client.CreateStorage(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating storage", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Storage.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StorageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan StorageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := plan.Content.ValueString()
	nodes := plan.Nodes.ValueString()
	pruneBackups := plan.PruneBackups.ValueString()

	disableInt := 0
	if !plan.Enabled.IsNull() && !plan.Enabled.ValueBool() {
		disableInt = 1
	}

	updateReq := &models.StorageUpdateRequest{
		Content:      &content,
		Disable:      &disableInt,
		Nodes:        &nodes,
		PruneBackups: &pruneBackups,
	}

	if !plan.Shared.IsNull() {
		sharedInt := 0
		if plan.Shared.ValueBool() {
			sharedInt = 1
		}
		updateReq.Shared = &sharedInt
	}

	if err := r.client.UpdateStorage(ctx, plan.Storage.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating storage", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Storage.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state StorageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteStorage(ctx, state.Storage.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting storage", err.Error())
	}
}

func (r *StorageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := StorageResourceModel{
		ID:      types.StringValue(req.ID),
		Storage: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StorageResource) readIntoModel(ctx context.Context, model *StorageResourceModel, diagnostics *diag.Diagnostics) {
	cfg, err := r.client.GetStorageConfig(ctx, model.Storage.ValueString())
	if err != nil {
		diagnostics.AddError("Error reading storage", err.Error())
		return
	}

	model.Type = types.StringValue(cfg.Type)
	model.Content = types.StringValue(cfg.Content)

	if cfg.Disable != nil {
		model.Enabled = types.BoolValue(*cfg.Disable == 0)
	} else {
		model.Enabled = types.BoolValue(true)
	}

	if cfg.Shared != nil {
		model.Shared = types.BoolValue(*cfg.Shared == 1)
	} else {
		model.Shared = types.BoolValue(false)
	}

	model.Path = types.StringValue(cfg.Path)
	model.Pool = types.StringValue(cfg.Pool)
	model.VGName = types.StringValue(cfg.VGName)
	model.Nodes = types.StringValue(cfg.Nodes)
	model.Server = types.StringValue(cfg.Server)
	model.Export = types.StringValue(cfg.Export)
	model.Share = types.StringValue(cfg.Share)
	model.Username = types.StringValue(cfg.Username)
	model.Domain = types.StringValue(cfg.Domain)
	model.Datastore = types.StringValue(cfg.Datastore)
	model.Namespace = types.StringValue(cfg.Namespace)
	model.Fingerprint = types.StringValue(cfg.Fingerprint)
	model.PruneBackups = types.StringValue(cfg.PruneBackups)
}
